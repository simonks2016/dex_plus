package internal

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goccy/go-json"
)

func (b *BinanceClient) OnConnecting(reason string) {
	//TODO implement me
	if b.logger != nil {
		b.logger.Println("Binance client OnConnecting:", reason)
	}
}

func (b *BinanceClient) OnConnected() {
	//TODO implement me
	b.isConnected.Store(true)

	if !b.IsRequireAuth {

		// 设置已验证
		b.authDone.Store(true)

		if b.subscribedParams != nil {
			if err := b.Send(b.subscribedParams.Json()); err != nil {
				if b.logger != nil {
					b.logger.Println(err)
				}
			}
		}
		// 定时处理死信队列
		go b.clearDeadMessage()

	}

}

func (b *BinanceClient) OnDisconnecting() {
	//TODO implement me
	b.isConnected.Store(false)
	b.authDone.Store(false)

	if b.subscribedParams != nil {
		b.Unsubscribe()
	}
}

func (b *BinanceClient) OnDisconnected() {
	//TODO implement me
	fmt.Println("Binance client OnDisconnected")
}

func (b *BinanceClient) OnMessage(data []byte) error {
	//TODO implement me

	var result map[string]interface{}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	} else if result == nil {
		return errors.New("invalid result")
	}

	if id, ex := result["id"]; ex {
		if b.logger != nil {
			b.logger.Printf("[info]Successfully to subscribe channel,result id : %s", id.(string))
		}
		return nil
	}

	if eventType, ex := result["e"]; ex {
		// 转化成string
		event := eventType.(string)
		// 获取处理函数
		if callers, ex := b.handlerMap[event]; ex {
			for _, callback := range callers {
				if err := b.pool.Submit(func() {
					if err := callback(result); err != nil {
						if b.logger != nil {
							b.logger.Printf("[error]Failed to read message:%s,%s", err.Error(), event)
						}
						return
					}
				}); err != nil {
					return err
				}
			}
		}
	}

	if lastUpdateId, ex := result["lastUpdateId"]; ex {
		for k, callers := range b.handlerMap {

			if strings.HasPrefix(k, "depth") && k != "depth" {

				for _, callback := range callers {
					if err := b.pool.Submit(func() {
						if err := callback(result); err != nil {
							if b.logger != nil {
								b.logger.Printf("[error]Failed to read message:%s,%s,%s", err.Error(), k, lastUpdateId)
							}
							return
						}
					}); err != nil {
						return err
					}
				}
			}
		}

	}

	return nil
}

func (b *BinanceClient) OnError(err error) {
	//TODO implement me
	fmt.Println("Binance client OnError:", err.Error())
}
