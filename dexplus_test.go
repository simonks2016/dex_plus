package main

import (
	"DexPlus/okx"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/panjf2000/ants/v2"
)

func TestDe(t *testing.T) {

	pool, err := ants.NewPool(100, ants.WithNonblocking(true))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	o1 := okx.NewPublic(context.Background(), pool)
	//
	o1.SetInstId("BTC-USDT")
	// connect
	if err = o1.Connect(); err != nil {
		fmt.Println(err.Error())
		return
	}
	o1.SubscribeTrade(func(trade []okx.Trade) error {

		for _, o := range trade {

			fmt.Printf("%s:%s %s:%s\n", o.InstId, o.Side, o.Px, o.Sz)
		}
		return nil
	})

	o1.SubscribeTicker(func(ticker []okx.Ticker) error {
		for _, o := range ticker {
			fmt.Printf("%s +> %s --> %s\n", o.InstId, o.Last, o.LastSz)
		}
		return nil
	})

	o1.SubscribeBook(okx.Books5Channel, func(books []okx.OrderBook) error {
		for _, book := range books {
			fmt.Printf("卖一%s:买一%s\n", book.Asks[0][0], book.Bids[0][0])
		}
		return nil
	})

	// 创建频道
	sig := make(chan os.Signal, 1)
	// 假如系统关闭
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	select {
	case <-sig:
		o1.Close()
	}

}
