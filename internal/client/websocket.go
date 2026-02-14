package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type WsClient struct {
	cfg    *Config
	logger *log.Logger
	dialer *websocket.Dialer

	ctx        context.Context
	cancelFunc context.CancelFunc

	closed atomic.Bool
	conn   atomic.Pointer[websocket.Conn]
	mu     sync.Mutex // 用于保护连接切换时的原子性，避免重复重连

	reconnectCh chan string
	writeCh     chan []byte
	readCh      chan []byte

	ob ConnectionObserver
}

func NewWsClient(ctx context.Context, cfg *Config) *WsClient {
	ctx, cancel := context.WithCancel(ctx)

	d := cfg.Dialer
	if d == nil {
		d = &websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: 10 * time.Second,
		}
	}

	if cfg.IsForbidIPV6 {
		d.NetDialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			nd := &net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}
			return nd.DialContext(ctx, "tcp4", address)
		}
	}

	w := &WsClient{
		logger:     cfg.Logger,
		cfg:        cfg,
		ctx:        ctx,
		dialer:     d,
		cancelFunc: cancel,

		reconnectCh: make(chan string, 1),
		writeCh:     make(chan []byte, cfg.WriteBufferSize),
		readCh:      make(chan []byte, cfg.ReadBufferSize),
	}
	return w
}

func (c *WsClient) Start() {
	if c.ob == nil {
		log.Fatal("websocket observer is nil")
	}

	go c.connLoop()
	go c.writePump()
	c.startWorkers()

	// 首次启动信号
	c.signalReconnect("initial_connect")
}

// connLoop 负责管理生命周期：拨号、重连、清理
func (c *WsClient) connLoop() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case reason := <-c.reconnectCh:
			if c.closed.Load() {
				return
			}
			c.handleConnect(reason)
		}
	}
}

func (c *WsClient) handleConnect(reason string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 1. 清理旧连接
	c.closeAndClearConn()

	c.ob.OnConnecting(reason)

	// 2. 指数退避重连
	backoff := c.cfg.ReconnectBackoffMin
	for {
		conn, _, err := c.dialer.DialContext(c.ctx, c.cfg.URL, c.cfg.Header)
		if err == nil {
			c.setupConn(conn)
			c.conn.Store(conn)
			go c.readPump(conn) // 为每个新连接开启独立的 readPump
			c.ob.OnConnected()
			return
		}

		c.logger.Printf("[ws] dial failed: %v, retry in %v", err, backoff)

		timer := time.NewTimer(backoff)
		select {
		case <-c.ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
			backoff = min(backoff*2, c.cfg.ReconnectBackoffMax)
		}
	}
}

func (c *WsClient) writePump() {
	ticker := time.NewTicker(c.cfg.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.doWrite(websocket.PingMessage, nil)
		case msg := <-c.writeCh:
			c.doWrite(websocket.TextMessage, msg)
		}
	}
}

// 统一写入方法，解决并发安全和超时问题
func (c *WsClient) doWrite(mt int, data []byte) {
	conn := c.conn.Load()
	if conn == nil {
		return
	}

	_ = conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteTimeout))
	var err error
	if mt == websocket.PingMessage {
		err = conn.WriteControl(mt, data, time.Now().Add(c.cfg.WriteTimeout))
	} else {
		err = conn.WriteMessage(mt, data)
	}

	if err != nil {
		c.signalReconnect(fmt.Sprintf("write_err: %v", err))
	}
}

func (c *WsClient) readPump(conn *websocket.Conn) {
	// 确保 readPump 退出时，如果是当前连接则触发重连
	defer func() {
		if c.conn.Load() == conn {
			c.signalReconnect("read_pump_exit")
		}
	}()

	for {
		msgType, r, err := conn.NextReader()
		if err != nil {
			return
		}

		if msgType != websocket.TextMessage && msgType != websocket.BinaryMessage {
			continue
		}

		data, err := io.ReadAll(io.LimitReader(r, c.cfg.maxMessageSize))
		if err != nil {
			return
		}

		select {
		case c.readCh <- data:
		default:
			c.logger.Printf("[ws] readCh full, drop msg")
		}
	}
}

// Send 业务层调用的发送方法
func (c *WsClient) Send(ctx context.Context, data []byte) error {
	if c.closed.Load() {
		return fmt.Errorf("client closed")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case c.writeCh <- data:
		return nil
	case <-time.After(time.Second): // 避免 writeCh 满时永久阻塞业务协程
		return fmt.Errorf("write channel busy")
	}
}

func (c *WsClient) signalReconnect(reason string) {
	select {
	case c.reconnectCh <- reason:
	default:
		// 已经在排队重连了
	}
}

func (c *WsClient) closeAndClearConn() {
	conn := c.conn.Swap(nil)
	if conn != nil {
		_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bye"))
		_ = conn.Close()
		c.ob.OnDisconnected()
	}
}

func (c *WsClient) Close() {
	if !c.closed.Swap(true) {
		c.cancelFunc()
		c.closeAndClearConn()
	}
}

// setupConn 初始化连接
func (c *WsClient) setupConn(conn *websocket.Conn) {
	// 1. 设置读取限制，防止大包攻击
	if c.cfg.maxMessageSize > 0 {
		conn.SetReadLimit(c.cfg.maxMessageSize)
	} else {
		conn.SetReadLimit(1024 * 1024) // 默认 1MB
	}

	// 2. 设置初始读取超时 (心跳检测的核心)
	// 如果在 PongWait 时间内没收到 Pong，NextReader 会返回错误
	_ = conn.SetReadDeadline(time.Now().Add(c.cfg.PongWait))

	// 3. 注册 Pong 处理回调
	// 每次收到服务端的 Pong，就顺延下一次的读取超时时间
	conn.SetPongHandler(func(appData string) error {
		_ = conn.SetReadDeadline(time.Now().Add(c.cfg.PongWait))
		return nil
	})

	// 4. (可选) 如果服务端发 Ping，我们也回 Pong
	conn.SetPingHandler(func(appData string) error {
		_ = conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteTimeout))
		return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(c.cfg.WriteTimeout))
	})
}

// startWorkers 启动共工作节点
func (c *WsClient) startWorkers() {
	workerNum := c.cfg.ReadWorkerNum
	if workerNum <= 0 {
		workerNum = 5 // 默认 5 个协程并行处理业务逻辑
	}

	for i := 0; i < workerNum; i++ {
		go func(id int) {
			for {
				select {
				case <-c.ctx.Done():
					return
				case msg, ok := <-c.readCh:
					if !ok {
						return
					}
					// 执行业务回调
					if c.ob != nil {
						if err := c.ob.OnMessage(msg); err != nil {
							if c.logger != nil {
								c.logger.Printf("[worker-%d] OnMessage error: %v", id, err)
							}
						}
					}
				}
			}
		}(i)
	}
}

// SetObserver 设置事件监听器
func (cli *WsClient) SetObserver(ob ConnectionObserver) Client {
	cli.ob = ob
	return cli
}

// Reconnect 重启
func (cli *WsClient) Reconnect(reason string) {
	cli.signalReconnect(reason)
}
