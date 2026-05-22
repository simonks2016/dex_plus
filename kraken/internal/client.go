package internal

import (
	"context"
	"errors"
	"log"
	"sync/atomic"

	"github.com/panjf2000/ants/v2"
	"github.com/simonks2016/dex_plus/internal/client"
	"github.com/simonks2016/dex_plus/kraken/params"
	"github.com/simonks2016/dex_plus/kraken/payload"
)

type KrakenClient struct {
	ctx               context.Context
	client            *client.WsClient
	logger            *log.Logger
	pool              *ants.Pool
	cfg               *client.Config
	isConnected       atomic.Bool
	isAuthDone        atomic.Bool
	isRequireAuth     bool
	handler           map[string][]Caller
	subscribeRequest  map[string][]string
	instrumentService *InstrumentService
	channelState      *SubscribeChannelState
}

func NewKrakenClient(ctx context.Context, cfg *client.Config) *KrakenClient {

	pool, _ := ants.NewPool(ants.DefaultAntsPoolSize, ants.WithNonblocking(true))

	krakenClient := &KrakenClient{
		ctx:               ctx,
		client:            client.NewWsClient(ctx, cfg),
		logger:            cfg.Logger,
		pool:              pool,
		cfg:               cfg,
		isRequireAuth:     cfg.IsNeedAuth,
		handler:           make(map[string][]Caller),
		subscribeRequest:  make(map[string][]string),
		instrumentService: NewInstrumentService(),
		channelState:      NewSubscribeChannelState(),
	}
	krakenClient.client.SetObserver(krakenClient)
	// 添加处理instrument
	krakenClient.handler["instrument"] = append(krakenClient.handler["instrument"], krakenClient.onInstrument)

	return krakenClient
}

func (k *KrakenClient) Send(data []byte) error {

	if !k.isConnected.Load() {
		return errors.New("the client is not connected")
	}

	ctx, cancel := context.WithTimeout(k.ctx, k.cfg.SendTimeout)
	defer cancel()

	//
	return k.client.Send(ctx, data)
}

func (k *KrakenClient) Connect() {
	k.client.Start()
}

func (k *KrakenClient) Close() {
	k.client.Close()
}

func (k *KrakenClient) Subscribe(channels ...SubscribeChannel) {
	for _, channel := range channels {
		k.subscribeRequest[channel.Channel] = append(k.subscribeRequest[channel.Channel], channel.Symbols...)
		k.handler[channel.Channel] = append(k.handler[channel.Channel], channel.Caller...)

		for _, symbol := range channel.Symbols {
			k.channelState.Switch(channel.Channel, symbol, Subscribing)
		}
	}
}

type SubscribeChannel struct {
	Channel string   `json:"channel"`
	Symbols []string `json:"symbols"`
	Caller  []Caller `json:"caller"`
}

func (k *KrakenClient) GetTradingPair(symbol string) (payload.Pair, bool) {
	return k.instrumentService.GetTradingPair(symbol)
}

func (k *KrakenClient) Resubscribe(channel string, symbols ...string) error {
	var newSymbols []string

	for _, symbol := range symbols {
		s, ex := k.channelState.Get(channel, symbol)
		if !ex {
			continue
		}
		// 只允许已订阅的进入重订阅，避免重复 resubscribe
		if s == Subscribed || s == SubscribeFailed {
			k.channelState.Switch(channel, symbol, Resubscribing)
			newSymbols = append(newSymbols, symbol)
		}
	}

	// 没有需要重订阅的，直接返回，避免空 unsubscribe/subscribe
	if len(newSymbols) == 0 {
		return nil
	}

	unsubscribeParam := params.NewKrakenParams(params.Unsubscribe, channel, newSymbols...)
	if err := k.Send(unsubscribeParam.Json()); err != nil {
		// 发送失败，建议恢复状态
		for _, symbol := range newSymbols {
			k.channelState.Switch(channel, symbol, Subscribed)
		}
		return err
	}

	return nil
}
