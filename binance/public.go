package binance

import (
	"context"
	"log"
	"time"

	"github.com/simonks2016/dex_plus/binance/internal"
	"github.com/simonks2016/dex_plus/binance/payload"
	"github.com/simonks2016/dex_plus/internal/client"
)

type Public struct {
	client  *internal.BinanceClient
	logger  *log.Logger
	symbols []string
}

func NewPublic(ctx context.Context, symbol ...string) *Public {

	cfg := client.NewConfig()
	cfg.WithURL(internal.WsURL)
	cfg.SetReadTimeout(time.Minute)
	cfg.SetReadWorkerNum(10)
	cfg.SetWriteTimeout(time.Minute)
	cfg.SetWriteBufferSize(50)
	cfg.SendTimeout = time.Minute
	cfg.SetReadBufferSize(5000)
	cfg.ForbidIPV6()

	cli := internal.NewBinanceClient(ctx, nil, cfg)

	return &Public{
		client:  cli,
		symbols: symbol,
	}
}

func subscribeChannel[T payload.BinancePayloadType](p *Public, channel string, callback func(T) error, opts ...SubscribeOption) {

	caller := func(data map[string]any) error {

		if d, err := payload.DecodeBinanceMap[T](data); err != nil {
			if p.logger != nil {
				p.logger.Printf("[error] failed to decode agg trade data: %v", err.Error())
			}
			return err
		} else {
			return callback(d)
		}
	}

	for _, symbol := range p.symbols {
		pa := internal.SubscribeParams{
			Channel:           channel,
			Is100Ms:           false,
			Symbol:            symbol,
			ReturnChannelName: channel,
		}
		for _, opt := range opts {
			opt(&pa)
		}
		p.client.Subscribe(&pa, caller)
	}
}

// SubscribeAggTrade 订阅归集交易
func (p *Public) SubscribeAggTrade(callback func(payload.AggTrade) error) {
	subscribeChannel[payload.AggTrade](p, "aggTrade", callback,
		WithReturnChannelName("aggTrade"),
		WithIs100Ms(),
	)
}

// SubscribeOrderBookDelta 订阅增量盘口深度数据
func (p *Public) SubscribeOrderBookDelta(callback func(payload.OrderBookDelta) error) {
	//
	subscribeChannel[payload.OrderBookDelta](p, "depth", callback,
		WithReturnChannelName("depthUpdate"),
		WithIs100Ms(),
	)
}

func (p *Public) Connect() {
	p.client.Connect()
}

func (p *Public) Close() {
	p.client.Close()
}

type SubscribeOption func(params *internal.SubscribeParams)

func WithReturnChannelName(name string) SubscribeOption {

	return func(params *internal.SubscribeParams) {
		params.ReturnChannelName = name
	}
}

func WithIs100Ms() SubscribeOption {
	return func(params *internal.SubscribeParams) {
		params.Is100Ms = true
	}
}

func (p *Public) ExchangeName() string { return "binance" }
