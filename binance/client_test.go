package binance

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/simonks2016/dex_plus/binance/payload"
	"github.com/simonks2016/dex_plus/common"
)

func TestClient(t *testing.T) {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	// 释放资源
	defer stop()

	cli := NewPublic(ctx,
		common.BinanceSymbol(common.BTC, common.USDT),
	)

	cli.Connect()

	/*
		cli.SubscribeOrderBook(func(symbol string, snapshot payload.OrderBookSnapshot) error {
			fmt.Println(symbol, snapshot)
			return nil
		})*/

	cli.SubscribeAggTrade(func(s string, trade payload.AggTrade) error {
		fmt.Println(s, trade)
		return nil
	})

	select {
	case <-ctx.Done():
		cli.Close()
	}

}
