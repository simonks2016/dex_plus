package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type WsClient struct {
	url    string
	header http.Header
	dialer *websocket.Dialer
	logger *log.Logger

	// lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	closed atomic.Bool
	cfg    *Config

	// conn owned by (re)connect loop; read/write loops use current conn
	conn atomic.Pointer[websocket.Conn]

	// signals
	reconnectCh chan string
	outCh       chan []byte // optional: send business msg

	handler         func(reader []byte) error
	notifyReconnect func(reason string)
	beforeReconnect func()
}

func NewWsClient(backgroundCtx context.Context, cfg *Config) *WsClient {
	ctx, cancel := context.WithCancel(backgroundCtx)

	d := cfg.Dialer
	if d == nil {
		d = websocket.DefaultDialer
	}
	if cfg.HandshakeTimeout > 0 {
		d.HandshakeTimeout = cfg.HandshakeTimeout
	}

	c := &WsClient{
		url:         cfg.URL,
		header:      cfg.Header,
		dialer:      d,
		logger:      cfg.Logger,
		ctx:         ctx,
		cancel:      cancel,
		reconnectCh: make(chan string, 1),
		outCh:       make(chan []byte, 1024),
		cfg:         cfg,
	}
	if c.logger == nil {
		c.logger = log.Default()
	}
	return c
}

func (ws *WsClient) Connect() {

	// run loops
	go ws.reconnectLoop()
	go ws.writeLoop()
	go ws.listen()

	// trigger first connect
	ws.signalReconnect("init connect")
}

func (ws *WsClient) IsConnecting() bool {
	return ws.closed.Load() == false
}

func (c *WsClient) Close() {
	if c.closed.Swap(true) {
		return
	}
	c.cancel()

	if conn := c.conn.Load(); conn != nil {
		_ = conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "client close"),
			time.Now().Add(1*time.Second),
		)
		_ = conn.Close()
	}
}

func (c *WsClient) Send(msg []byte) error {
	if c.closed.Load() {
		return errors.New("client closed")
	}
	select {
	case c.outCh <- msg:
		return nil
	case <-c.ctx.Done():
		return c.ctx.Err()
	default:
		// 避免 outCh 堵死把系统拖垮：宁可丢/报错也别无限堆积
		return errors.New("send queue full")
	}
}

func (c *WsClient) signalReconnect(reason string) {
	if c.closed.Load() {
		return
	}
	select {
	case c.reconnectCh <- reason:
	default:
		c.logger.Printf("It is already reconnecting: %s", reason)
	}
}

// --- loops ---
func (c *WsClient) reconnectLoop() {
	backoff := c.cfg.ReconnectBackoffMin

	for {
		select {
		case <-c.ctx.Done():
			return
		case reason := <-c.reconnectCh:
			if c.closed.Load() {
				return
			}
			// 执行重新连接前工作
			c.beforeReconnect()
			// 打印日志
			c.logger.Printf("[ws] reconnect requested: %s", reason)
			// close old conn
			if old := c.conn.Load(); old != nil {
				_ = old.Close()
			}

			for {
				if c.closed.Load() {
					return
				}
				conn, _, err := c.dialer.DialContext(c.ctx, c.url, c.header)
				if err == nil {
					c.setupConn(conn, *c.cfg)
					c.conn.Store(conn)
					c.notifyReconnect(reason)
					c.logger.Printf("[ws] connected")
					backoff = c.cfg.ReconnectBackoffMin
					break
				}

				c.logger.Printf("[ws] dial failed: %v; backoff=%v", err, backoff)
				select {
				case <-time.After(backoff):
				case <-c.ctx.Done():
					return
				}

				// exponential backoff with cap
				backoff *= 2
				if backoff > c.cfg.ReconnectBackoffMax {
					backoff = c.cfg.ReconnectBackoffMax
				}
			}
		}
	}
}

func (c *WsClient) setupConn(conn *websocket.Conn, cfg Config) {
	// 关键：读 deadline + pong handler 续命
	_ = conn.SetReadDeadline(time.Now().Add(cfg.PongWait))

	// 应对Pong的处理
	conn.SetPongHandler(func(appData string) error {
		_ = conn.SetReadDeadline(time.Now().Add(cfg.PongWait))
		return nil
	})

	// 如果你希望对方 ping 来时也续命（通常可以）
	conn.SetPingHandler(func(appData string) error {
		_ = conn.SetReadDeadline(time.Now().Add(cfg.PongWait))
		// gorilla 默认不会自动 pong，这里主动回
		deadline := time.Now().Add(cfg.WriteTimeout)
		return conn.WriteControl(websocket.PongMessage, []byte(appData), deadline)
	})
}

func (c *WsClient) listen() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		conn := c.conn.Load()
		if conn == nil {
			time.Sleep(50 * time.Millisecond)
			continue
		}
		messageType, reader, err := conn.NextReader()
		if err != nil {
			// 读超时、断链、协议错误 → 统一触发重连
			c.signalReconnect("read error: " + err.Error())
			time.Sleep(50 * time.Millisecond)
			continue
		}
		// 收到任何消息，都算“活着”，续一下 deadline（可选，但很实用）
		_ = conn.SetReadDeadline(time.Now().Add(c.cfg.PongWait))

		// 假如发送过来是text message
		if messageType == websocket.TextMessage {
			data, err := io.ReadAll(io.LimitReader(reader, c.cfg.maxMessageSize))
			if err != nil {
				c.logger.Printf("[ws] handler error: %s", err)
			}
			if c.cfg.payloadType == JSONType && json.Valid(data) {
				if err = c.handler(data); err != nil {
					c.logger.Printf("[ws] handler error: %s", err)
				}
			} else {
				c.logger.Printf("[ws] It is not JSON: %s", string(data))
			}
		} else if messageType == websocket.BinaryMessage {
			c.logger.Printf("[ws] binary message received")
			continue
		} else if messageType == websocket.CloseMessage {
			c.Close()
			c.logger.Printf("[ws] connection closed")
			return
		} else {
			c.logger.Printf("[ws] revicer message type: %v", messageType)
			return
		}

	}
}

func (c *WsClient) writeLoop() {
	ticker := time.NewTicker(c.cfg.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return

		case <-ticker.C:
			conn := c.conn.Load()
			if conn == nil {
				continue
			}
			// 定时 ping
			_ = conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteTimeout))
			// 定时发送PING 数据
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.signalReconnect("ping failed: " + err.Error())
			}

		case msg := <-c.outCh:
			conn := c.conn.Load()
			if conn == nil {
				// 还没连上：这里你可以选择丢弃或等待
				continue
			}
			_ = conn.SetWriteDeadline(time.Now().Add(c.cfg.WriteTimeout))
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				c.signalReconnect("write failed: " + err.Error())
			}
		}
	}
}

func (c *WsClient) WriteMessage(data []byte) {
	c.outCh <- data
}

func (c *WsClient) SetHandler(handler func(data []byte) error) *WsClient {
	c.handler = handler
	return c
}

func (c *WsClient) SetReconnectNotify(notify func(string), before func()) *WsClient {
	c.notifyReconnect = notify
	c.beforeReconnect = before
	return c
}
