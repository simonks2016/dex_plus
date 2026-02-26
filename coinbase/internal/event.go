package internal

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/simonks2016/dex_plus/coinbase/params"
)

func (c *CoinbaseClient) OnConnecting(reason string) {
	//TODO implement me
	fmt.Println(reason)
}

func (c *CoinbaseClient) OnConnected() {
	//TODO implement me
	fmt.Println("On Connected")

	c.isConnected.Store(true)
	if !c.isRequireAuth {
		c.authDone.Store(true)
		// 构建订阅参数
		p := params.NewSubscribeParams(params.Subscribe, c.symbols...)
		p.AddChannel(c.channels...)
		if !p.IsEmpty() {
			// 发送订阅信息
			if err := c.Send(p.Json()); err != nil {
				if c.logger != nil {
					c.logger.Printf("failed to subscribe channel,%s", err.Error())
				}
			}
		}

	}

}

func (c *CoinbaseClient) OnDisconnecting() {
	//TODO implement me
	p := params.NewSubscribeParams(params.Unsubscribe, c.symbols...)
	p.AddChannel(c.channels...)
	if !p.IsEmpty() {
		// 发送取消订阅信息
		if err := c.Send(p.Json()); err != nil {
			if c.logger != nil {
				c.logger.Printf("send subscribe error,%s", err.Error())
			}
			return
		}
	}

}

func (c *CoinbaseClient) OnDisconnected() {
	//TODO implement me
	fmt.Println("On Disconnected")
}

func (c *CoinbaseClient) OnMessage(data []byte) error {
	//TODO implement me

	var e Envelope

	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}

	channelName := e.Type

	if callers, ex := c.handler[channelName]; ex {
		for _, caller := range callers {
			if err := c.pool.Submit(func() {
				if err := caller(data); err != nil {
					if c.logger != nil {
						c.logger.Printf("[error]failed to handler channel(%s),%s", channelName, err.Error())
					}
					return
				}
			}); err != nil {
				return fmt.Errorf("failed to submit task of pool,%s", err.Error())
			}
		}
	}
	return nil
}

func (c *CoinbaseClient) OnError(err error) {
	//TODO implement me
	fmt.Println("On Error", err)
}
