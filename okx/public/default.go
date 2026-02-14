package public

import (
	"context"
	"log"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/simonks2016/dex_plus/internal/client"
	"github.com/simonks2016/dex_plus/okx"
	"github.com/simonks2016/dex_plus/okx/internal"
)

type Public struct {
	client     *internal.OKXClient
	logger     *log.Logger
	instId     []string
	instFamily []string
	ctx        context.Context
}

func NewPublic(bg context.Context, pool *ants.Pool, opts ...client.Option) OKXPublic {

	cfg := client.NewConfig()
	cfg.SetWriteBufferSize(4000)
	cfg.SetReadBufferSize(4000)
	cfg.SetReadWorkerNum(100)
	cfg.SetReadTimeout(time.Second * time.Duration(10))
	cfg.SetWriteTimeout(time.Second * time.Duration(10))
	cfg.SendTimeout = time.Minute * time.Duration(10)
	cfg.WithURL(okx.PublicURL(true))
	cfg.IsNeedAuth = false
	cfg.IsForbidIPV6 = false

	for _, opt := range opts {
		opt(cfg)
	}

	if pool != nil {
		pool, _ = ants.NewPool(ants.DefaultAntsPoolSize, ants.WithNonblocking(true))
	}

	// 创建一个新的客户端
	cli := internal.NewOKXClient(bg, nil, cfg)
	cli.SetThreadPool(pool)

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
	SubscribeTrade(callback func(trade []okx.AggregatedTrades) error)
	SubscribeTradeAll(callback func(trade []okx.RawTrades) error)
	SubscribeBook(channel string, callback func(books []okx.OrderBook) error)
}
