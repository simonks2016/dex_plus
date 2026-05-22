package DexPlus

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"testing"
	"time"

	"github.com/simonks2016/dex_plus/kraken"
	"github.com/simonks2016/dex_plus/kraken/payload"
)

func NewLogger() *log.Logger {
	return log.New(
		os.Stdout, // 也可以换成你自己的 Writer（SLS、文件等）
		"[OKX] ",  // 前缀
		log.LstdFlags|log.Lmicroseconds|log.Lshortfile,
	)
}

func TestNew(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := kraken.NewPublic(
		ctx,
		kraken.WithLogger(NewLogger()),
	)

	client.SetSymbols("BTC/USDT")
	client.SubscribeOrderBook(time.Duration(1)*time.Second, func(ob []payload.OrderBook) error {

		var allAsks []float64
		var allBids []float64

		for _, book := range ob {

			for _, ask := range book.Asks {
				allAsks = append(allAsks, ask.Price)
			}

			for _, bid := range book.Bids {
				allBids = append(allBids, bid.Price)
			}
		}

		// asks 正常应该：从低到高
		sort.Float64s(allAsks)

		// bids 正常应该：从高到低
		sort.Slice(allBids, func(i, j int) bool {
			return allBids[i] > allBids[j]
		})

		fmt.Println("========== ASKS ==========")
		for i, ask := range allAsks {
			fmt.Printf("[%d] %.2f\n", i, ask)
		}

		fmt.Println("========== BIDS ==========")
		for i, bid := range allBids {
			fmt.Printf("[%d] %.2f\n", i, bid)
		}

		// 检查是否 crossed
		if len(allAsks) > 0 && len(allBids) > 0 {

			bestAsk := allAsks[0]
			bestBid := allBids[0]

			crossed := bestBid - bestAsk

			fmt.Printf(
				"\nBestBid=%.2f BestAsk=%.2f Crossed=%.2f\n",
				bestBid,
				bestAsk,
				crossed,
			)
		}
		return nil
	})
	client.Connect()

	// 监听 Ctrl+C
	sigChan := make(chan os.Signal, 1)

	signal.Notify(
		sigChan,
		os.Interrupt,
		syscall.SIGTERM,
	)

	select {
	case <-sigChan:
		fmt.Println("received shutdown signal")

	case <-ctx.Done():
		fmt.Println("context canceled")
	}

	client.Close()

	fmt.Println("client closed")
}

func NewInt(i int) *int {
	return &i
}

func NewString(s string) *string {
	return &s
}

// 存在问题：
// 第二: 验证失败会出现批量发送验证信息，不停的重新启动
