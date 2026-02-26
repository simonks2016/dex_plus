package bitstamp

import (
	"context"
	"log"

	"github.com/simonks2016/dex_plus/bitstamp/internal"
	"github.com/simonks2016/dex_plus/bitstamp/payload"
	"github.com/simonks2016/dex_plus/internal/client"
)

type Public struct {
	client  *internal.BitstampClient
	logger  *log.Logger
	ctx     context.Context
	symbols []string
}

func NewPublic(ctx context.Context, symbols ...string) *Public {

	cfg := client.NewConfig()
	cfg.WithURL(internal.WsURL)
	cfg.IsNeedAuth = false

	return &Public{
		client:  internal.NewBitstampClient(ctx, cfg),
		logger:  cfg.Logger,
		ctx:     ctx,
		symbols: symbols,
	}
}

func (p *Public) Connect() { p.client.Connect() }
func (p *Public) Close()   { p.client.Close() }

// SubscribeTrades 订阅成交数据
func (p *Public) SubscribeTrades(callback func(string, payload.Trade) error) {
	p.client.Subscribe("live_trades", func(env *internal.Envelope) error {
		data, err := payload.ParseData[payload.Trade](env)
		if err != nil {
			return err
		}
		return callback(env.GetSymbol(), data)
	}, p.symbols...)
}

// SubscribeOrderBook 订阅盘口数据
func (p *Public) SubscribeOrderBook(callback func(string, payload.OrderBook) error) {
	p.client.Subscribe("diff_order_book_", func(env *internal.Envelope) error {
		data, err := payload.ParseData[payload.OrderBook](env)
		if err != nil {
			return err
		}
		return callback(env.GetSymbol(), data)
	}, p.symbols...)
}
