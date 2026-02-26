package internal

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/simonks2016/dex_plus/bitstamp/params"
	"github.com/simonks2016/dex_plus/internal/client"
)

type BitstampClient struct {
	ctx           context.Context
	client        *client.WsClient
	logger        *log.Logger
	pool          *ants.Pool
	cfg           *client.Config
	authDone      atomic.Bool
	isConnected   atomic.Bool
	isRequireAuth bool
	handler       map[string][]Caller
}

func NewBitstampClient(ctx context.Context, cfg *client.Config) *BitstampClient {

	pool, _ := ants.NewPool(ants.DefaultAntsPoolSize, ants.WithNonblocking(true))

	cli := BitstampClient{
		ctx:           ctx,
		client:        client.NewWsClient(ctx, cfg),
		logger:        cfg.Logger,
		pool:          pool,
		cfg:           cfg,
		isRequireAuth: cfg.IsNeedAuth,
		handler:       make(map[string][]Caller),
	}
	cli.client.SetObserver(&cli)
	return &cli
}

// Connect 连接
func (cli *BitstampClient) Connect() {
	cli.client.Start()
}

// Close 关闭连接
func (cli *BitstampClient) Close() {
	cli.client.Close()
}

// Send 发送信息
func (cli *BitstampClient) Send(dataBytes []byte) error {

	if !cli.isConnected.Load() {
		return errors.New("the client has been disconnected")
	}
	ctx, cancel := context.WithTimeout(cli.ctx, time.Second*time.Duration(5))
	defer cancel()
	// 发送信息
	return cli.client.Send(ctx, dataBytes)
}

// Subscribe 订阅
func (cli *BitstampClient) Subscribe(channel string, call func(*Envelope) error, symbols ...string) {

	channels := make([]string, 0, len(symbols))

	for _, symbol := range symbols {
		ch := channel + "_" + symbol
		channels = append(channels, ch)
		cli.handler[ch] = append(cli.handler[ch], call)
	}
}

// Unsubscribe 取消订阅
func (cli *BitstampClient) Unsubscribe() {

	for channelName, _ := range cli.handler {
		p1 := params.NewUnsubscribeParams(channelName)
		// 发送信息
		if err := cli.Send(p1.Json()); err != nil {
			if cli.logger != nil {
				cli.logger.Printf("[error] failed to unsubscribe channel,%s", err.Error())
			}
			return
		}
	}
}
