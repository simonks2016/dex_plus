package DexPlus

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
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
	client.SubscribeOrderBook(time.Duration(100)*time.Millisecond, func(ob []payload.OrderBook) error {
		fmt.Println(ob)
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
