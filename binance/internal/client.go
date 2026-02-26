package internal

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/simonks2016/dex_plus/internal/client"
)

type SubscribeParams struct {
	Channel           string `json:"channel"`
	Is100Ms           bool   `json:"is100ms"`
	Symbol            string `json:"symbol"`
	ReturnChannelName string `json:"return_channel_name"`
}

type BinanceClient struct {
	ctx              context.Context
	client           *client.WsClient
	logger           *log.Logger
	pool             *ants.Pool
	cfg              *client.Config
	auth             *Auth
	IsRequireAuth    bool
	url              string
	authDone         atomic.Bool
	isConnected      atomic.Bool
	handlerMap       map[string][]Caller
	subscribedParams *BinanceParams
	deadQueue        [][]byte
}

func NewBinanceClient(ctx context.Context, auth *Auth, cfg *client.Config) *BinanceClient {

	pool, _ := ants.NewPool(ants.DefaultAntsPoolSize, ants.WithNonblocking(true))

	cli := &BinanceClient{
		client:           client.NewWsClient(ctx, cfg),
		auth:             auth,
		ctx:              ctx,
		handlerMap:       make(map[string][]Caller),
		IsRequireAuth:    cfg.IsNeedAuth,
		logger:           cfg.Logger,
		url:              cfg.URL,
		cfg:              cfg,
		pool:             pool,
		deadQueue:        make([][]byte, 0),
		subscribedParams: NewBinanceParams(SubscribeMethod),
	}
	cli.client.SetObserver(cli)

	return cli
}

func (b *BinanceClient) Connect() {
	b.client.Start()
}

func (b *BinanceClient) Close() {
	b.client.Close()
}
func (b *BinanceClient) Subscribe(params *SubscribeParams, caller ...Caller) {

	if b.handlerMap == nil {
		b.handlerMap = make(map[string][]Caller)
	}
	// 添加处理函数
	b.handlerMap[params.ReturnChannelName] = append(b.handlerMap[params.ReturnChannelName], caller...)
	// 创建订阅参数

	if b.isConnected.Load() {
		// 添加到已有订阅参数
		b.subscribedParams.Add(params.Symbol, params.Channel, params.Is100Ms)
		// 创建一个新的立刻发送给币安
		dataBytes := NewBinanceParams(SubscribeMethod, params.Channel).Add(
			params.Symbol,
			params.Channel,
			params.Is100Ms).Json()
		// 发送订阅信息
		if err := b.Send(dataBytes); err != nil {
			b.deadQueue = append(b.deadQueue, dataBytes)
			return
		}
	} else {
		b.subscribedParams.Add(params.Symbol, params.Channel, params.Is100Ms)
	}
}

func (b *BinanceClient) Unsubscribe() {

	if b.subscribedParams != nil {
		// 复制一个新的
		p1 := b.subscribedParams.CopyNew(UnsubscribeMethod)
		// 发送取消订阅信息
		if err := b.Send(p1.Json()); err != nil {
			if b.logger != nil {
				b.logger.Println("Failed to Unsubscribe Channel:", err)
			}
			return
		}
	}
}

func (b *BinanceClient) Send(dataBytes []byte) error {
	if b.isConnected.Load() {
		if b.authDone.Load() {
			ctx, cancel := context.WithTimeout(b.ctx, time.Second*time.Duration(5))
			defer cancel()
			// 发送数据
			return b.client.Send(ctx, dataBytes)
		}
	} else {
		return errors.New("the client has been disconnected")
	}
	return nil
}

func (b *BinanceClient) clearDeadMessage() {

	t1 := time.NewTicker(time.Second * time.Duration(20))
	defer t1.Stop()

	for {
		select {
		case <-b.ctx.Done():
			return
		case <-t1.C:
			if err := b.replySendDeadMessage(); err != nil {
				if b.logger != nil {
					b.logger.Println("failed to reply send dead message error:", err)
				}
				continue
			}
		}
	}
}

func (b *BinanceClient) replySendDeadMessage() error {

	if len(b.deadQueue) > 0 {
		// 1. 暂存旧切片并立即重置原队列，避免逻辑处理期间原队列被修改
		queueToProcess := b.deadQueue
		b.deadQueue = nil // 使用 nil 比 [][]byte{} 更节省一次分配

		for i, data := range queueToProcess {
			if err := b.Send(data); err != nil {
				if b.logger != nil {
					b.logger.Println("[error] Failed to Subscribe Binance Channel:", err.Error())
				}
				// 失败处理：如果发送失败，建议将剩余未发送的重新放回队列，防止数据丢失
				b.deadQueue = append(queueToProcess[i:], b.deadQueue...)
				return err
			}
		}
	}
	return nil
}
