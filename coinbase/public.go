package coinbase

import (
	"context"
	"log"
	"time"

	"github.com/goccy/go-json"
	"github.com/simonks2016/dex_plus/coinbase/internal"
	"github.com/simonks2016/dex_plus/coinbase/payload"
	"github.com/simonks2016/dex_plus/internal/bookManager"
	"github.com/simonks2016/dex_plus/internal/client"
)

type Public struct {
	client      *internal.CoinbaseClient
	cfg         *client.Config
	logger      *log.Logger
	ctx         context.Context
	bookManager *bookManager.BookManager
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
		client:      cli,
		cfg:         cfg,
		logger:      cfg.Logger,
		ctx:         ctx,
		bookManager: bookManager.NewBookManager(),
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

// SubscribeOrderBook 优化：支持流式设置和统一管理
func (p *Public) SubscribeOrderBook(interval time.Duration, callback func([]payload.OrderBook) error) {
	p.client.Subscribe("level2")
	// 这里的 Handler 建议在初始化时就设置好，避免重复调用
	p.client.SetHandler("snapshot", p.handlingOrderBookSnapshot)
	p.client.SetHandler("l2update", p.handlingOrderBookDelta)

	// 启动定时器
	p.setSnapshotTimer(p.ctx, interval, callback)
}

// handlingOrderBookSnapshot 优化：修复 append bug，提取解析逻辑
func (p *Public) handlingOrderBookSnapshot(data []byte) error {
	var t1 payload.OrderBookSnapshot
	if err := json.Unmarshal(data, &t1); err != nil {
		return err
	}

	ob := p.bookManager.GetOrCreate(t1.ProductId)
	// 修正：容量设为总和，长度设为 0
	levels := make([]bookManager.Level, 0, len(t1.Bids)+len(t1.Asks))

	// 辅助解析闭包：减少重复代码
	parseAndAppend := func(items [][]string, isBid bool) {
		for _, item := range items {
			if len(item) < 2 {
				continue
			}
			px, err1 := bookManager.PriceTicks(item[0], 100)
			sz, err2 := bookManager.SizeFloat(item[1])
			if err1 != nil || err2 != nil {
				continue // 实际生产环境建议打一条采样日志
			}
			levels = append(levels, bookManager.Level{
				PriceTicks: px,
				Size:       sz,
				IsBids:     isBid,
			})
		}
	}

	parseAndAppend(t1.Bids, true)
	parseAndAppend(t1.Asks, false)

	return ob.ApplySnapshot(time.Now(), levels...)
}

// handlingOrderBookDelta 优化：Side 快速判断
func (p *Public) handlingOrderBookDelta(data []byte) error {
	var t1 payload.OrderBookUpdate
	if err := json.Unmarshal(data, &t1); err != nil {
		return err
	}

	ob := p.bookManager.GetOrCreate(t1.ProductId)
	levels := make([]bookManager.Level, 0, len(t1.Changes))

	for _, ch := range t1.Changes {
		if len(ch) < 3 {
			continue
		}

		px, err1 := bookManager.PriceTicks(ch[1], 100)
		sz, err2 := bookManager.SizeFloat(ch[2])
		if err1 != nil || err2 != nil {
			continue
		}

		// Coinbase 的 side 通常是 "buy" 或 "sell"
		isBid := ch[0] == "buy" || ch[0] == "BUY"
		levels = append(levels, bookManager.Level{
			PriceTicks: px,
			Size:       sz,
			IsBids:     isBid,
		})
	}
	return ob.ApplyL2Update(levels, time.Now())
}

// setSnapshotTimer
func (p *Public) setSnapshotTimer(ctx context.Context, interval time.Duration, callback func([]payload.OrderBook) error) {
	// 假设 StartSnapshotTimerAsync 内部已有 ticker 逻辑
	p.bookManager.StartSnapshotTimerAsync(ctx, interval, 20, func(snapshots []bookManager.TopNSnapshot) {
		if len(snapshots) == 0 {
			return
		}

		resp := make([]payload.OrderBook, len(snapshots))
		for i, snapshot := range snapshots {
			bids := make([]payload.Level, 0, len(snapshot.Bids))
			asks := make([]payload.Level, 0, len(snapshot.Asks))

			for _, b := range snapshot.Bids {
				bids = append(bids, payload.Level{
					Price: bookManager.PriceTo(b.PriceTicks),
					Size:  b.Size,
				})
			}
			for _, a := range snapshot.Asks {
				asks = append(asks, payload.Level{
					Price: bookManager.PriceTo(a.PriceTicks),
					Size:  a.Size,
				})
			}

			resp[i] = payload.OrderBook{
				ProductId: snapshot.ProductID,
				Bids:      bids,
				Asks:      asks,
				Time:      time.UnixMilli(snapshot.Ts),
			}
		}

		// 异步执行回调，防止阻塞管理器
		go func(d []payload.OrderBook) {
			if err := callback(d); err != nil && p.logger != nil {
				p.logger.Printf("[error] callback failed: %v", err)
			}
		}(resp)
	})
}
