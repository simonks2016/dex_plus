package public

import (
	"context"
	"log"
	"time"

	"github.com/simonks2016/dex_plus/okx"
	"github.com/simonks2016/dex_plus/okx/internal"
	"github.com/simonks2016/dex_plus/option"
	"github.com/simonks2016/dex_plus/websocket"
)

type Public struct {
	client     *internal.OKXClient
	logger     *log.Logger
	instId     []string
	instFamily []string
	ctx        context.Context
}

func NewPublic(bg context.Context, opts ...option.Option) OKXPublic {

	cfg := websocket.NewConfig()
	cfg.SetWriteBufferSize(4000)
	cfg.SetReadBufferSize(4000)
	cfg.SetReadWorkerNum(100)
	cfg.SetReadTimeout(time.Second * time.Duration(10))
	cfg.SetWriteTimeout(time.Second * time.Duration(10))
	cfg.WithURL("wss://ws.okx.com:8443/ws/v5/public")
	// 创建一个新的客户端
	cli := internal.NewOKXClient(
		bg,
		nil,
		cfg,
		opts...)

	return &Public{
		ctx:    bg,
		client: cli,
	}
}

type OKXPublic interface {
	SetLogger(logger *log.Logger) OKXPublic
	SetInstId(id ...string) OKXPublic
	SetInstFamily(id ...string) OKXPublic
	Connect()
	Reconnect()
	Close()
	SubscribeTicker(callback func(tickers []okx.Ticker) error)
	SubscribeTrade(callback func(trade []okx.Trade) error)
	SubscribeBook(channel string, callback func(books []okx.OrderBook) error)
}
