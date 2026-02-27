package kraken

import (
	"context"
	"log"
	"time"

	"github.com/simonks2016/dex_plus/internal/client"
	"github.com/simonks2016/dex_plus/kraken/internal"
	"github.com/simonks2016/dex_plus/kraken/payload"
)

type Public struct {
	client  *internal.KrakenClient
	symbols []string
	logger  *log.Logger
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

	cli := internal.NewKrakenClient(ctx, cfg)

	p1 := Public{
		client:  cli,
		symbols: symbols,
		logger:  nil,
	}
	return &p1
}

func (p *Public) SubscribeTrade(callback func(trades []payload.Trade) error) {

	// 订阅Trade频道
	p.client.Subscribe(internal.SubscribeChannel{
		Channel: "trade",
		Symbols: p.symbols,
		Caller: []internal.Caller{
			func(envelope *internal.KrakenEnvelope) error {
				data, err := payload.ParseData[payload.Trade](envelope)
				if err != nil {
					return err
				}
				return callback(data)
			},
		},
	})
}

func (p *Public) SubscribeOrderBook(callback func(ob []payload.OrderBook) error) {
	// 订阅盘口数据
	p.client.Subscribe(internal.SubscribeChannel{
		Channel: "book",
		Symbols: p.symbols,
		Caller: []internal.Caller{
			func(envelope *internal.KrakenEnvelope) error {
				data, err := payload.ParseData[payload.OrderBook](envelope)
				if err != nil {
					return err
				}
				return callback(data)
			},
		},
	})
}

// Connect 连接
func (p *Public) Connect() {
	p.client.Connect()
}

// Close 关闭连接
func (p *Public) Close() {
	p.client.Close()
}

// ExchangeName 交易所名字
func (p *Public) ExchangeName() string {
	return "kraken"
}

// SetSymbols 设置品种
func (p *Public) SetSymbols(symbols ...string) {
	p.symbols = symbols
}
