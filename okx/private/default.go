package private

import (
	"context"
	"log"
	"time"

	"github.com/simonks2016/dex_plus/okx"
	"github.com/simonks2016/dex_plus/okx/internal"
	"github.com/simonks2016/dex_plus/okx/param"
	"github.com/simonks2016/dex_plus/option"
	"github.com/simonks2016/dex_plus/websocket"
)

type Private struct {
	client *internal.OKXClient
	logger *log.Logger
}

func NewPrivate(apiKey, secretKey, passphrase string, bg context.Context, options ...option.Option) OKXPrivate {

	var cfg = websocket.NewConfig()
	cfg.SetWriteBufferSize(4000)
	cfg.SetReadBufferSize(4000)
	cfg.SetReadWorkerNum(100)
	cfg.SetReadTimeout(time.Second * time.Duration(10))
	cfg.SetWriteTimeout(time.Second * time.Duration(10))
	cfg.WithURL(okx.PrivateURL(true))

	if runMode := option.GetOption("is_sandbox_environment", options...); runMode != nil {
		if isSandBox, ok := runMode.(bool); ok {
			cfg.WithURL(okx.PrivateURL(!isSandBox))
		}
	}
	options = append(options, option.WithURL(cfg.URL))
	// 创建新的
	cli := internal.NewOKXClient(bg,
		internal.NewAuth(apiKey, passphrase, secretKey),
		cfg,
		options...)

	return &Private{client: cli}
}

type OKXPrivate interface {
	SetLogger(logger *log.Logger) OKXPrivate

	SubscribePositionAndBalance(func(posAndBala ...okx.PositionAndBalance) error)
	SubscribeTrade(func(trade ...okx.TradeFill) error)

	PlaceOrder(...param.PlaceOrderParams) error
	AmendOrder(...param.AmendOrder) error
	CancelOrder(...param.CancelOrder) error
	Connect()
	Close()
	Reconnect()
}
