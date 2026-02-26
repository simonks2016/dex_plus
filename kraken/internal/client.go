package internal

import (
	"context"
	"errors"
	"log"
	"sync/atomic"

	"github.com/panjf2000/ants/v2"
	"github.com/simonks2016/dex_plus/internal/client"
)

type KrakenClient struct {
	ctx              context.Context
	client           *client.WsClient
	logger           *log.Logger
	pool             *ants.Pool
	cfg              *client.Config
	isConnected      atomic.Bool
	isAuthDone       atomic.Bool
	isRequireAuth    bool
	handler          map[string][]Caller
	subscribeRequest map[string][]string
}

func NewKrakenClient(ctx context.Context, cfg *client.Config) *KrakenClient {

	pool, _ := ants.NewPool(ants.DefaultAntsPoolSize, ants.WithNonblocking(true))

	krakenClient := &KrakenClient{
		ctx:              ctx,
		client:           client.NewWsClient(ctx, cfg),
		logger:           cfg.Logger,
		pool:             pool,
		cfg:              cfg,
		isRequireAuth:    cfg.IsNeedAuth,
		handler:          make(map[string][]Caller),
		subscribeRequest: make(map[string][]string),
	}
	krakenClient.client.SetObserver(krakenClient)

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
	// 输入
	for _, channel := range channels {
		k.subscribeRequest[channel.Channel] = append(k.subscribeRequest[channel.Channel], channel.Symbols...)
		k.handler[channel.Channel] = append(k.handler[channel.Channel], channel.Caller...)
	}
}

type SubscribeChannel struct {
	Channel string   `json:"channel"`
	Symbols []string `json:"symbols"`
	Caller  []Caller `json:"caller"`
}
