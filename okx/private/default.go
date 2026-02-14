package private

import (
	"context"
	"log"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/simonks2016/dex_plus/internal/client"
	"github.com/simonks2016/dex_plus/okx"
	"github.com/simonks2016/dex_plus/okx/internal"
	"github.com/simonks2016/dex_plus/okx/param"
)

type Private struct {
	client *internal.OKXClient
	logger *log.Logger
}

func NewPrivate(apiKey, secretKey, passphrase string, bg context.Context, pool *ants.Pool, opts ...client.Option) OKXPrivate {

	var cfg = client.NewConfig()
	cfg.SetWriteBufferSize(4000)
	cfg.SetReadBufferSize(4000)
	cfg.SetReadWorkerNum(100)
	cfg.SetReadTimeout(time.Second * time.Duration(10))
	cfg.SetWriteTimeout(time.Second * time.Duration(10))
	cfg.WithURL(okx.PrivateURL(true))
	cfg.IsNeedAuth = true
	cfg.SendTimeout = time.Minute * time.Duration(5)

	for _, opt := range opts {
		opt(cfg)
	}
	cfg.IsNeedAuth = true

	if pool != nil {
		pool, _ = ants.NewPool(ants.DefaultAntsPoolSize, ants.WithNonblocking(true))
	}
	// 创建新的
	cli := internal.NewOKXClient(bg,
		internal.NewAuth(apiKey, passphrase, secretKey),
		cfg)
	cli.SetThreadPool(pool)

	return &Private{client: cli}
}

type OKXPrivate interface {
	SetLogger(logger *log.Logger) OKXPrivate

	SubscribePosition(func(pos ...okx.Position) error, *int64)
	SubscribePositionAndBalance(func(posAndBala ...okx.PositionAndBalance) error)
	SubscribeTrade(func(trade ...okx.TradeFill) error)
	SubscribeOrderFilled(func(orders ...okx.OrderState) error)

	PlaceOrder(...param.PlaceOrderParams) error
	AmendOrder(...param.AmendOrder) error
	CancelOrder(...param.CancelOrder) error
	Connect()
	Close()
	Reconnect()
}
