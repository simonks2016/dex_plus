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
	"github.com/simonks2016/dex_plus/common"
	"github.com/simonks2016/dex_plus/okx"
	"github.com/simonks2016/dex_plus/okx/public"
)

func NewLogger() *log.Logger {
	return log.New(
		os.Stdout, // 也可以换成你自己的 Writer（SLS、文件等）
		"[OKX] ",  // 前缀
		log.LstdFlags|log.Lmicroseconds|log.Lshortfile,
	)
}

func TestNew(t *testing.T) {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	// 释放资源
	defer stop()

	pool, _ := ants.NewPool(ants.DefaultAntsPoolSize, ants.WithNonblocking(true))

	p1 := public.NewPublic(
		ctx,
		pool,
		okx.WithForbidIpV6(),
		okx.WithLogger(NewLogger()),
		okx.WithSandboxEnv(),
	)

	p1.SetInstId(common.OKXSymbol(common.BTC))
	p1.Connect()

	p1.SubscribeTrade(func(trades []okx.AggregatedTrades) error {

		for _, trade := range trades {

			fmt.Println(trade)
		}

		return nil
	})

	select {
	case <-ctx.Done():

		//
		pool.Free()
		p1.Close()
	}

}

// 存在问题：
// 第二: 验证失败会出现批量发送验证信息，不停的重新启动
