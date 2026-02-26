package internal

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/goccy/go-json"
	"github.com/panjf2000/ants/v2"
	"github.com/simonks2016/dex_plus/internal/client"
)

type CoinbaseClient struct {
	ctx           context.Context
	client        *client.WsClient
	logger        *log.Logger
	pool          *ants.Pool
	cfg           *client.Config
	authDone      atomic.Bool
	isConnected   atomic.Bool
	isRequireAuth bool
	handler       map[string][]Caller
	channels      []string
	symbols       []string
}

func NewCoinbaseClient(ctx context.Context, cfg *client.Config) *CoinbaseClient {

	pool, _ := ants.NewPool(ants.DefaultAntsPoolSize, ants.WithNonblocking(true))

	cli := CoinbaseClient{
		ctx:           ctx,
		client:        client.NewWsClient(ctx, cfg),
		logger:        cfg.Logger,
		pool:          pool,
		cfg:           cfg,
		handler:       make(map[string][]Caller),
		channels:      make([]string, 0),
		symbols:       make([]string, 0),
		isRequireAuth: cfg.IsNeedAuth,
	}
	cli.client.SetObserver(&cli)

	cli.handler["error"] = append(cli.handler["error"], func(data []byte) error {

		var d = make(map[string]any)

		if err := json.Unmarshal(data, &d); err != nil {
			return err
		} else {
			if cli.logger != nil {
				cli.logger.Printf("[error]received a error message from Coinbase,%s", d["message"])
			}
		}

		return nil
	})

	return &cli
}

func (cli *CoinbaseClient) Connect() {
	cli.client.Start()
}

func (cli *CoinbaseClient) Close() {
	cli.client.Close()
}

func (cli *CoinbaseClient) Send(dataBytes []byte) error {

	if !cli.isConnected.Load() {
		return fmt.Errorf("client is not connected")
	}
	ctx, cancel := context.WithTimeout(cli.ctx, time.Second*time.Duration(5))
	defer cancel()
	// 发送数据
	return cli.client.Send(ctx, dataBytes)
}

func (cli *CoinbaseClient) Subscribe(channels ...string) {

	seenChannels := make(map[string]bool)

	// 记录当前已有的值，防止重复追加（如果 cli.channels 初始不为空）
	for _, v := range cli.channels {
		seenChannels[v] = true
	}

	for _, channel := range channels {
		// 1. Channels 去重追加
		if !seenChannels[channel] {
			cli.channels = append(cli.channels, channel)
			seenChannels[channel] = true
		}
	}
}

func (cli *CoinbaseClient) SetSymbols(symbols ...string) {

	seenSymbols := make(map[string]bool)
	for _, v := range cli.symbols {
		seenSymbols[v] = true
	}
	for _, symbol := range symbols {
		// 1. Channels 去重追加
		if !seenSymbols[symbol] {
			cli.symbols = append(cli.symbols, symbol)
			seenSymbols[symbol] = true
		}
	}
}

func (cli *CoinbaseClient) SetHandler(name string, caller ...Caller) {
	cli.handler[name] = append(cli.handler[name], caller...)
}

type SubscribeParams struct {
	Channel    string
	Symbol     []string
	Caller     []Caller
	ReturnName string
}
