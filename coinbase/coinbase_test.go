package coinbase

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/simonks2016/dex_plus/coinbase/payload"
	"github.com/simonks2016/dex_plus/common"
)

func TestCoinbase(t *testing.T) {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	// 释放资源
	defer stop()

	cli := NewPublic(ctx, common.CoinBaseSymbol(common.BTC))

	cli.Connect()

	cli.SubscribeTrade(func(trades payload.MatchedTrade) error {
		fmt.Println(trades)
		return nil
	})

	cli.SubscribeOrderBook(func(snapshot payload.OrderBook) error {
		fmt.Println(snapshot)
		return nil

	}, func(update payload.OrderBookUpdate) error {
		fmt.Println(update)
		return nil
	})

	for {
		select {
		case <-ctx.Done():
			cli.Close()
			return
		}
	}

}
