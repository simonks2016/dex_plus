package internal

import (
	"context"

	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/goccy/go-json"
	"github.com/simonks2016/dex_plus/internal/client"
)

type Ob struct {
	client *client.WsClient
}

func (ob *Ob) OnConnecting(reason string) {
	//TODO implement me
	fmt.Println("OnConnecting", reason)
}

func (ob *Ob) OnConnected() {
	//TODO implement me

	data := map[string]interface{}{

		"op": "subscribe",
		"args": []map[string]interface{}{
			{
				"channel": "trades",
				"instId":  "BTC-USDT",
			},
		},
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := ob.client.Send(context.Background(), marshal); err != nil {
		fmt.Println(err)
		return
	}

}

func (ob *Ob) OnDisconnecting() {
	//TODO implement me
	fmt.Println("OnDisconnecting")
}

func (ob *Ob) OnDisconnected() {
	//TODO implement me
	fmt.Println("OnDisconnected")
}

func (ob *Ob) OnMessage(data []byte) error {
	//TODO implement me
	fmt.Println("OnMessage", string(data))
	return nil
}

func (ob *Ob) OnError(err error) {
	//TODO implement me
	fmt.Println("OnError", err.Error())
}

func TestClient(t *testing.T) {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	// 释放资源
	defer stop()

	cfg := client.NewConfig().WithURL("wss://ws.okx.com:8443/ws/v5/public")
	cfg.IsForbidIPV6 = true

	cli := client.NewWsClient(
		ctx,
		cfg,
	)
	var ob1 = Ob{client: cli}

	cli.SetObserver(&ob1)
	cli.Start()

	for {
		select {
		case <-ctx.Done():
			return
		}
	}

}
