package DexPlus

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/panjf2000/ants/v2"
	"github.com/simonks2016/dex_plus/okx"
	"github.com/simonks2016/dex_plus/okx/public"
	"github.com/simonks2016/dex_plus/option"
)

func NewLogger() *log.Logger {
	return log.New(
		os.Stdout, // 也可以换成你自己的 Writer（SLS、文件等）
		"[OKX] ",  // 前缀
		log.LstdFlags|log.Lmicroseconds|log.Lshortfile,
	)
}

func TestNew(t *testing.T) {

	ctx := context.Background()

	pool, _ := ants.NewPool(100, ants.WithNonblocking(true))
	logger := NewLogger()

	/*
		p := private.NewPrivate(
			"f668642d-ad76-4e72-8c85-0075acd73a5a",
			"E698B82FDF9CB97C333C306669D6AD4D",
			"testBbq20251002@",
			ctx,
			option.WithWriteBufferSize(4000),
			option.WithReadBufferSize(4000),
			option.WithThreadPool(pool),
			option.WithURL("wss://wspap.okx.com:8443/ws/v5/private"),
			option.WithForbidIpv6(true),
			option.WithLogger(logger),
		)*/

	p := public.NewPublic(
		ctx,
		option.WithWriteBufferSize(4000),
		option.WithReadBufferSize(4000),
		option.WithForbidIpv6(true),
		option.WithLogger(logger),
		option.WithThreadPool(pool),
	)

	p.SetLogger(logger)
	p.SetInstId("BTC-USDT")

	p.Connect()

	p.SubscribeTrade(func(trades []okx.Trade) error {
		for _, trade := range trades {
			fmt.Printf("%s:%s %s on %s\n", trade.InstId, trade.Side, trade.Sz, trade.Px)
		}
		return nil
	})

	// 创建频道
	sig := make(chan os.Signal, 1)
	// 假如系统关闭
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	select {

	case <-sig:
		p.Close()
	}

}
