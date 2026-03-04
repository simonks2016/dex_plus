package kraken

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/simonks2016/dex_plus/internal/bookManager"
	"github.com/simonks2016/dex_plus/internal/client"
	"github.com/simonks2016/dex_plus/kraken/internal"
	"github.com/simonks2016/dex_plus/kraken/payload"
)

type Public struct {
	client      *internal.KrakenClient
	ctx         context.Context
	symbols     []string
	logger      *log.Logger
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

	cli := internal.NewKrakenClient(ctx, cfg)

	p1 := Public{
		client:      cli,
		symbols:     symbols,
		logger:      cfg.Logger,
		ctx:         ctx,
		bookManager: bookManager.NewBookManager(),
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

func (p *Public) SubscribeOrderBook(interval time.Duration, callback func(ob []payload.OrderBook) error) {
	// 订阅盘口数据
	p.client.Subscribe(internal.SubscribeChannel{
		Channel: "book",
		Symbols: p.symbols,
		Caller: []internal.Caller{
			func(envelope *internal.KrakenEnvelope) error {
				return p.handlerOrderBook(envelope)
			},
		},
	})
	// 异步定时发送盘口快照
	p.setSnapshotTimer(p.ctx, interval, 20, callback)
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

func (p *Public) handlerOrderBook(env *internal.KrakenEnvelope) error {
	// 1. 安全检查 Type
	if env.Type == nil {
		return nil
	}
	msgType := strings.ToLower(*env.Type)

	// 2. 解析数据
	data, err := payload.ParseData[payload.OrderBook](env)
	if err != nil || len(data) == 0 {
		return err
	}

	// 3. 直接遍历数据进行处理，避免创建中间 map (o2)
	for _, datum := range data {
		// 预分配 level 切片容量
		levels := make([]bookManager.Level, 0, len(datum.Bids)+len(datum.Asks))

		// 填充 Bids
		for _, bid := range datum.Bids {
			levels = append(levels, bookManager.Level{
				PriceTicks: bookManager.NewPrice(bid.Price),
				Size:       bid.Qty,
				IsBids:     true,
			})
		}
		// 填充 Asks
		for _, ask := range datum.Asks {
			levels = append(levels, bookManager.Level{
				PriceTicks: bookManager.NewPrice(ask.Price),
				Size:       ask.Qty,
				IsBids:     false,
			})
		}

		// 4. 根据类型路由到不同的处理逻辑
		book := p.bookManager.GetOrCreate(datum.Symbol)

		switch msgType {
		case "snapshot":
			if err := book.ApplySnapshot(datum.Timestamp, levels...); err != nil {
				return err
			}
		case "update":
			// 假设 ApplyL2Update 接收 []Level 和 Time
			if err := book.ApplyL2Update(levels, datum.Timestamp); err != nil {
				return err
			}
		default:
			// 忽略未定义类型
		}
	}
	return nil
}

// setSnapshotTimer Set the timer of book snapshot
func (p *Public) setSnapshotTimer(ctx context.Context, interval time.Duration, n int, callback func([]payload.OrderBook) error) {

	// 异步执行Book Manager 定时发送快照
	p.bookManager.StartSnapshotTimerAsync(ctx, interval, n, func(snapshots []bookManager.TopNSnapshot) {

		response := make([]payload.OrderBook, len(snapshots))

		for index, snapshot := range snapshots {
			// 分别初始化 Bids 和 Asks
			bids := make([]payload.OrderBookItem, len(snapshot.Bids))
			asks := make([]payload.OrderBookItem, len(snapshot.Asks))

			// 填充 Bids
			for i, bid := range snapshot.Bids {
				bids[i] = payload.OrderBookItem{
					Price: bookManager.PriceTo(bid.PriceTicks),
					Qty:   bid.Size,
				}
			}

			// 填充 Asks (修正：写入 asks 数组)
			for i, ask := range snapshot.Asks {
				asks[i] = payload.OrderBookItem{
					Price: bookManager.PriceTo(ask.PriceTicks),
					Qty:   ask.Size,
				}
			}

			response[index] = payload.OrderBook{
				Symbol:    snapshot.ProductID,
				Bids:      bids,
				Asks:      asks,
				Checksum:  0, // 如果需要校验和，需在此计算
				Timestamp: time.UnixMilli(snapshot.Ts),
			}
		}
		// 执行回调（建议考虑是否需要 go callback(response) 异步处理）
		go func() {
			if err := callback(response); err != nil {
				if p.logger != nil {
					p.logger.Printf("snapshot error: %v", err)
				}
			}
		}()
	})
}
