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
		common.BinanceSymbol(common.ETH, common.USDT),
		common.BinanceSymbol(common.XRP, common.USDT),
		common.BinanceSymbol(common.SOL, common.USDT),
	)

	cli.Connect()
	cli.SubscribeOrderBookDelta(func(ob payload.OrderBookDelta) error {
		fmt.Println(ob)
		return nil
	})
	cli.SubscribeAggTrade(func(trade payload.AggTrade) error {
		fmt.Println(trade)
		return nil
	})

	select {
	case <-ctx.Done():
		cli.Close()
	}

}
