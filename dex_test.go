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
	"github.com/simonks2016/dex_plus/okx/private"
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

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	// 释放资源
	defer stop()

	pool, _ := ants.NewPool(100, ants.WithNonblocking(true))
	logger := NewLogger()

	p := private.NewPrivate(
		"f668642d-ad76-4e72-8c85-0075acd73a5a",
		"E698B82FDF9CB97C333C306669D6AD4D",
		"testBbq20251002@",
		ctx,
		option.WithWriteBufferSize(4000),
		option.WithReadBufferSize(4000),
		option.WithThreadPool(pool),
		option.WithSandBoxEnvironment(),
		option.WithForbidIpv6(true),
		option.WithLogger(logger),
	)

	p.Connect()

	time.Sleep(time.Second * time.Duration(5))

	p.SubscribePosition(func(pos ...okx.Position) error {
		for _, p1 := range pos {
			side := "long"
			if p1.PosCcy == p1.Ccy {
				side = "short"
			}

			fmt.Println(side, p1.Pos, p1.Liab)
		}
		return nil
	}, nil)

	select {
	case <-ctx.Done():
		p.Close()
	}

}

// 存在问题：
// 第一：假如网络抖动会导致订阅消息，因没有得到验证而吞了。
// 第二: 验证失败会出现批量发送验证信息，不停的重新启动
