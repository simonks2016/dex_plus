package bitstamp

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/simonks2016/dex_plus/bitstamp/payload"
	"github.com/simonks2016/dex_plus/common"
)

func TestBitstampClient(t *testing.T) {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	// 释放资源
	defer stop()

	p := NewPublic(ctx,
		common.BitstampSymbol(common.BTC),
		common.BinanceSymbol(common.SOL),
	)

	p.Connect()

	p.SubscribeOrderBook(func(s string, book payload.OrderBook) error {

		fmt.Println(s, book)
		return nil

	})

	for {
		select {
		case <-ctx.Done():
			return
		}
	}

}
