package coinbase

import (
	"context"
	"log"
	"time"

	"github.com/goccy/go-json"
	"github.com/simonks2016/dex_plus/coinbase/internal"
	"github.com/simonks2016/dex_plus/coinbase/payload"
	"github.com/simonks2016/dex_plus/internal/client"
)

type Public struct {
	client *internal.CoinbaseClient
	cfg    *client.Config
	logger *log.Logger
}

func NewPublic(ctx context.Context, symbols ...string) *Public {

	cfg := client.NewConfig()
	cfg.WithURL(internal.WsURL)
	cfg.SetReadTimeout(time.Minute)
	cfg.SetReadWorkerNum(10)
	cfg.SetWriteTimeout(time.Minute)
	cfg.SetWriteBufferSize(50)
	cfg.SendTimeout = time.Minute
	cfg.SetReadBufferSize(5000)
	cfg.ForbidIPV6()

	cli := internal.NewCoinbaseClient(ctx, cfg)
	cli.SetSymbols(symbols...)

	return &Public{
		client: cli,
		cfg:    cfg,
		logger: cfg.Logger,
	}
}

func (p *Public) Connect() { p.client.Connect() }
func (p *Public) Close()   { p.client.Close() }
func (p *Public) SubscribeTrade(callback func(trades payload.MatchedTrade) error) {

	p.client.Subscribe("matches")
	p.client.SetHandler("match", func(data []byte) error {
		var t1 payload.MatchedTrade
		if err := json.Unmarshal(data, &t1); err != nil {
			return err
		}
		return callback(t1)
	})
}
func (p *Public) SubscribeOrderBook(callbackSnapshot func(payload.OrderBook) error, callbackUpdate func(update payload.OrderBookUpdate) error) {

	p.client.Subscribe("level2")
	p.client.SetHandler("snapshot", func(data []byte) error {
		var t1 payload.OrderBook
		if err := json.Unmarshal(data, &t1); err != nil {
			return err
		}
		return callbackSnapshot(t1)
	})
	p.client.SetHandler("l2update", func(data []byte) error {
		var t1 payload.OrderBookUpdate
		if err := json.Unmarshal(data, &t1); err != nil {
			return err
		}
		return callbackUpdate(t1)
	})
}

// ExchangeName 交易所名字
func (p *Public) ExchangeName() string { return "coinbase" }

// SetSymbols 设置品种
func (p *Public) SetSymbols(symbols ...string) { p.client.SetSymbols(symbols...) }
