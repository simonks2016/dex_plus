package websocket

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
	url    string
	header http.Header
	dialer *websocket.Dialer
	logger *log.Logger
	cfg    *Config

	ctx    context.Context
	cancel context.CancelFunc
	closed atomic.Bool

	// 当前连接（只作为读/写使用；切换由 connLoop 控制）
	conn atomic.Pointer[websocket.Conn]

	// 控制信号
	reconnectCh chan string

	// 业务通道
	readCh      chan []byte
	writeCh     chan []byte
	writeAuthCh chan []byte

	// “就绪门闩”：连接+验权成功后 close
	readyMu sync.Mutex
	readyCh chan struct{}

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
	if cfg.HandshakeTimeout > 0 {
		d.HandshakeTimeout = cfg.HandshakeTimeout
	}

	if cfg.IsForbidIPV6 {
		d.NetDialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
			nd := &net.Dialer{
				Timeout:   10 * time.Second, // 你也可以用 cfg.DialTimeout
				KeepAlive: 30 * time.Second,
			}
			// 忽略传进来的 network（通常是 "tcp"），强制 "tcp4"
			return nd.DialContext(ctx, "tcp4", address)
		}
	}

	c := &WsClient{
		url:         cfg.URL,
		header:      cfg.Header,
		dialer:      d,
		logger:      cfg.Logger,
		ctx:         ctx,
		cancel:      cancel,
		reconnectCh: make(chan string, 1),
		writeCh:     make(chan []byte, cfg.WriteBufferSize),
		readCh:      make(chan []byte, cfg.ReadBufferSize),
		writeAuthCh: make(chan []byte, 64),
		readyCh:     make(chan struct{}),
		cfg:         cfg,
	}
	if c.logger == nil {
		c.logger = log.Default()
	}
	return c
}

func (c *WsClient) SetObserver(ob ConnectionObserver) *WsClient {
	c.ob = ob
	return c
}

func (c *WsClient) Start() {
	if c.ob == nil {
		panic("websocket start failed: websocket observer is nil")
	}
	go c.connLoop()
	go c.writePump()
	c.startWorkers()
	c.signalReconnect("init")
}

func (c *WsClient) connLoop() {
	backoffMin := c.cfg.ReconnectBackoffMin
	backoffMax := c.cfg.ReconnectBackoffMax

	for {
		select {
		case <-c.ctx.Done():
			return
		case reason := <-c.reconnectCh:
			if c.closed.Load() {
				return
			}
			if c.ob != nil {
				c.ob.OnConnecting(reason)
			}

			// 重置 ready gate（因为要重新登录）
			c.resetReady()

			// 关闭旧连接
			c.closeAndClearConn()

			backoff := backoffMin
			for {
				conn, _, err := c.dialer.DialContext(c.ctx, c.url, c.header)
				if err == nil {
					c.setupConn(conn)
					c.conn.Store(conn)

					if c.logger != nil {
						c.logger.Printf("[okx] successful connected to %s from okx", c.url)
					}

					// 启动 readPump（每个连接一个）
					go c.readPump(conn)
					if c.ob != nil {
						// 发送已连接
						c.ob.OnConnected()
						// 连接成功 → 开始鉴权流程
						c.callAuthHook(conn)
					} else {
						c.markReady()
					}
					break
				}

				c.logger.Printf("[ws] dial failed: %v, retrying in %v", err, backoff)

				timer := time.NewTimer(backoff)
				select {
				case <-timer.C:
					timer.Stop()
					backoff = min(backoff*2, backoffMax)
				case <-c.ctx.Done():
					timer.Stop()
					return
				}
			}
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
			c.writeControl(websocket.PingMessage, nil)
		case msg := <-c.writeAuthCh:
			c.writeMessage(websocket.TextMessage, msg)
		case msg := <-c.writeCh:
			c.writeMessage(websocket.TextMessage, msg)
		}
	}
}

func (c *WsClient) writeMessage(mt int, data []byte) {
	conn := c.conn.Load()
	if conn == nil {
		return
	}

	_ = conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteTimeout))
	if err := conn.WriteMessage(mt, data); err != nil {
		c.signalReconnect("write failed: " + err.Error())
	}
}

func (c *WsClient) writeControl(mt int, data []byte) {
	conn := c.conn.Load()
	if conn == nil {
		return
	}

	deadline := time.Now().Add(c.cfg.WriteTimeout)
	_ = conn.SetWriteDeadline(deadline)
	if err := conn.WriteControl(mt, data, deadline); err != nil {
		c.signalReconnect("control failed: " + err.Error())
	}
}

func (c *WsClient) readPump(conn *websocket.Conn) {
	defer func() {
		// 只有“当前 conn”断了才触发重连：避免旧 conn 的 readPump 误触发
		if c.conn.Load() == conn {
			c.signalReconnect("read loop exited")
		}
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		msgType, r, err := conn.NextReader()
		if err != nil {
			if c.conn.Load() == conn {
				if c.ob != nil {
					c.ob.OnError(err)
				}
			}
			return
		}

		if msgType != websocket.TextMessage && msgType != websocket.BinaryMessage {
			continue
		}

		limited := io.LimitReader(r, c.cfg.maxMessageSize)
		data, err := io.ReadAll(limited)
		if err != nil {
			if c.conn.Load() == conn {
				if c.ob != nil {
					c.ob.OnError(err)
				}
			}
			return
		}

		select {
		case c.readCh <- data:
		default:
			// read 队列满：你可以选择丢弃 / 触发重连 / 记录告警
			c.logger.Printf("[ws] readCh full, dropping message")
		}
	}
}

func (c *WsClient) resetReady() {
	c.readyMu.Lock()
	defer c.readyMu.Unlock()
	c.readyCh = make(chan struct{})
}

func (c *WsClient) markReady() {
	c.readyMu.Lock()
	ch := c.readyCh
	c.readyMu.Unlock()

	select {
	case <-ch:
		// 已经 ready
	default:
		close(ch)
	}
}

func (c *WsClient) WaitReady(ctx context.Context) error {
	c.readyMu.Lock()
	ch := c.readyCh
	c.readyMu.Unlock()

	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-c.ctx.Done():
		return c.ctx.Err()
	}
}

func (c *WsClient) Send(ctx context.Context, data []byte) error {
	if c.closed.Load() {
		return fmt.Errorf("client closed")
	}

	if err := c.WaitReady(ctx); err != nil {
		return err
	}

	select {
	case c.writeCh <- data:
		return nil
	default:
		return fmt.Errorf("writeCh full")
	}
}

func (c *WsClient) Close() {
	if c.ob != nil {
		c.ob.OnDisconnecting()
	}
	if c.closed.Swap(true) {
		return
	}
	c.cancel()
	c.closeAndClearConn()
}

func (c *WsClient) Reconnect(reason string) {
	c.signalReconnect(reason)
}

func (c *WsClient) setupConn(conn *websocket.Conn) {
	conn.SetReadLimit(c.cfg.maxMessageSize)
	_ = conn.SetReadDeadline(time.Now().Add(c.cfg.PongWait))

	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(c.cfg.PongWait))
		return nil
	})
}

// sendAuth 发送鉴权信息
// parameters:
// @ctx context.Context
// @data []byte
// response:
// @err error
func (c *WsClient) sendAuth(ctx context.Context, data []byte) error {
	if c.closed.Load() {
		return fmt.Errorf("client closed")
	}
	if c.conn.Load() == nil {
		return fmt.Errorf("not connected")
	}

	if c.logger != nil {
		c.logger.Printf("[ws] sending auth data")
	}

	select {
	case c.writeAuthCh <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("write queue full")
	}
}

// sendData
// parameters:
// @msg []byte 信息二进制
func (c *WsClient) sendData(msg []byte) {
	if conn := c.conn.Load(); conn != nil {
		// 设置写入超时时间
		_ = conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteTimeout))
		// 发送Text Message
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			c.signalReconnect("write failed: " + err.Error())
		}
	}
}

func (c *WsClient) callAuthHook(conn *websocket.Conn) {
	// 每次重连前你已经 resetReady()，这里按鉴权结果决定是否 markReady
	if c.ob == nil {
		c.markReady()
		return
	}

	var once sync.Once

	send := func(b []byte) error {

		if b == nil || len(b) <= 0 {
			c.markReady()
			return nil
		}
		// 防止旧连接还在回调时误发
		if c.conn.Load() != conn {
			return fmt.Errorf("stale connection")
		}
		// 给鉴权发送一个独立超时（避免卡死）
		timeout := 15 * time.Second
		// 设置超时时间
		ctx, cancel := context.WithTimeout(c.ctx, timeout)
		defer cancel()
		// 发送“高优先级”信息
		return c.sendAuth(ctx, b) // 内部高优先级
	}

	done := func(err error) {
		once.Do(func() {
			// 旧连接 done 不能影响新连接
			if c.conn.Load() != conn {
				return
			}
			if err == nil {
				c.markReady()
				return
			}
			// 鉴权失败：保持未 ready，并触发重连/告警
			if c.ob != nil {
				c.ob.OnError(fmt.Errorf("auth failed: %w", err))
			}
			c.signalReconnect("auth failed: " + err.Error())
		})
	}
	// 发 login + 保存 done，待 login ack 后回调 done
	c.ob.OnAuth(send, done)
}

func (c *WsClient) signalReconnect(reason string) {
	select {
	case c.reconnectCh <- reason:
	default:
		// 已经在尝试重连中，无需重复发送信号
		c.logger.Printf("[ws] reconnect channel full")
	}
}

func (c *WsClient) closeAndClearConn() {
	conn := c.conn.Swap(nil)
	if conn == nil {
		return
	}
	_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "client exit"), time.Now().Add(time.Second))
	_ = conn.Close()

	if c.ob != nil {
		c.ob.OnDisconnected()
	}
}

func (c *WsClient) startWorkers() {
	workerNum := c.cfg.ReadWorkerNum
	if workerNum <= 0 {
		workerNum = 10
	}

	for i := 0; i < workerNum; i++ {
		go func(id int) {
			for {
				select {
				case <-c.ctx.Done():
					return
				case msg := <-c.readCh:
					if c.ob != nil {
						if err := c.ob.OnMessage(msg); err != nil {
							if c.logger != nil {
								c.logger.Printf("[error] on message failed: " + err.Error())
							}
						}
					}
					continue
				}
			}
		}(i)
	}
}

// IsConnected 外部调用：判断当前是否处于可用连接状态
func (c *WsClient) IsConnected() bool {
	return !c.closed.Load() && c.conn.Load() != nil
}
