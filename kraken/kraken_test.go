package kraken

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/simonks2016/dex_plus/common"
	"github.com/simonks2016/dex_plus/kraken/payload"
)

func TestClient(t *testing.T) {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	// 释放资源
	defer stop()

	p1 := NewPublic(ctx, common.KrakenSymbol(common.BTC))

	p1.SubscribeTrade(func(trades []payload.Trade) error {
		fmt.Println(trades)
		return nil
	})

	p1.SubscribeOrderBook(func(ob []payload.OrderBook) error {
		fmt.Println(ob)
		return nil
	})

	p1.Connect()
	for {
		select {
		case <-ctx.Done():
			p1.Close()
		}
	}

}
