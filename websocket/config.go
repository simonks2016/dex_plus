package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type PayloadType int

const (
	JSONType PayloadType = 1
	ProtobufType
)

type Config struct {
	URL    string
	Header http.Header
	Dialer *websocket.Dialer
	Logger *log.Logger

	// timeouts
	HandshakeTimeout time.Duration

	// keepalive
	PingInterval        time.Duration // how often we send ping
	PongWait            time.Duration // max time to wait pong/any msg
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration // write deadline
	ReconnectBackoffMin time.Duration
	ReconnectBackoffMax time.Duration

	maxMessageSize int64

	WriteBufferSize int
	ReadBufferSize  int
	ReadWorkerNum   int
	IsForbidIPV6    bool
}

func NewConfig() *Config {
	return &Config{
		PingInterval:        15 * time.Second,
		PongWait:            45 * time.Second,
		ReadTimeout:         5 * time.Second,
		WriteTimeout:        5 * time.Second,
		ReconnectBackoffMin: 500 * time.Millisecond,
		ReconnectBackoffMax: 10 * time.Second,
		maxMessageSize:      4 << 20,
		ReadBufferSize:      4000,
		WriteBufferSize:     4000,
		ReadWorkerNum:       10,
		IsForbidIPV6:        false,
	}
}
func (c *Config) WithURL(url string) *Config {
	c.URL = url
	return c
}
func (c *Config) WithHeader(header http.Header) *Config {
	c.Header = header
	return c
}
func (c *Config) WithLogger(logger *log.Logger) *Config {
	c.Logger = logger
	return c
}
func (c *Config) SetHandshakeTimeout(timeout time.Duration) *Config {
	c.HandshakeTimeout = timeout
	return c
}
func (c *Config) SetPingInterval(timeout time.Duration) *Config {
	c.PingInterval = timeout
	return c
}
func (c *Config) SetPongWait(timeout time.Duration) *Config {
	c.PongWait = timeout
	return c
}
func (c *Config) SetWriteTimeout(timeout time.Duration) *Config {
	c.WriteTimeout = timeout
	return c
}
func (c *Config) SetReadTimeout(timeout time.Duration) *Config {
	c.ReadTimeout = timeout
	return c
}

func (c *Config) ForbidIPV6() *Config {
	c.IsForbidIPV6 = true
	return c
}

func (c *Config) SetMaxMessageSize(size int) *Config {
	c.maxMessageSize = int64(size)
	return c
}
func (c *Config) SetReadBufferSize(size int) *Config {
	c.ReadBufferSize = size
	return c
}
func (c *Config) SetWriteBufferSize(size int) *Config {
	c.WriteBufferSize = size
	return c
}
func (c *Config) SetReadWorkerNum(num int) *Config {
	c.ReadWorkerNum = num
	return c
}
