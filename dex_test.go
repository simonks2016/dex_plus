package DexPlus

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/panjf2000/ants/v2"
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
	//logger := NewLogger()

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
	)

	p.Connect()

	select {
	case <-ctx.Done():
		p.Close()
	}

}
