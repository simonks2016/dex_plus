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

	"github.com/panjf2000/ants/v2"
	"github.com/simonks2016/dex_plus/okx"
	"github.com/simonks2016/dex_plus/okx/param"
	"github.com/simonks2016/dex_plus/okx/private"
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

	p1 := private.NewPrivate(
		"", "", "",
		ctx,
		pool,
		okx.WithForbidIpV6(),
		okx.WithLogger(NewLogger()),
		okx.WithSandboxEnv())

	p1.Connect()

	p1.SubscribePositionAndBalance(func(posAndBala ...okx.PositionAndBalance) error {

		return nil

	})

	time.Sleep(2 * time.Second)

	err := p1.PlaceOrder(param.PlaceOrderParams{
		InstIdCode: NewInt(2021032601102993),
		TdMode:     "cross",
		Ccy:        NewString("USDT"),
		ClOrdId:    NewString("Ks3987haa"),
		Tag:        nil,
		Side:       "sell",
		PosSide:    NewString("short"),
		OrdType:    "limit",
		SZ:         "0.02",
		Px:         NewString("72165"),
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	select {
	case <-ctx.Done():

		//
		pool.Free()
		p1.Close()
	}

}

func NewInt(i int) *int {
	return &i
}

func NewString(s string) *string {
	return &s
}

// 存在问题：
// 第二: 验证失败会出现批量发送验证信息，不停的重新启动
