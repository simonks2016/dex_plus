package internal

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/simonks2016/dex_plus/kraken/params"
)

func (k *KrakenClient) OnConnecting(reason string) {
	//TODO implement me
	if k.logger != nil {
		k.logger.Println(reason)
	}
}

func (k *KrakenClient) OnConnected() {
	//TODO implement me
	//存储已连接状态
	k.isConnected.Store(true)
	if !k.isRequireAuth {
		// 存储验证状态
		k.isAuthDone.Store(true)
		//
		for channel, strings := range k.subscribeRequest {
			// 构建订阅参数
			p := params.NewKrakenParams(params.Subscribe, channel, strings...)
			// 发送订阅参数
			if err := k.Send(p.Json()); err != nil {
				if k.logger != nil {
					k.logger.Printf("send subscribe error,%s", err.Error())
				}
				return
			}
		}
	}
}

func (k *KrakenClient) OnDisconnecting() {
	//TODO implement me
	if k.logger != nil {
		k.logger.Println("on disconnecting")
	}
	// 全部取消订阅
	for channel, strings := range k.subscribeRequest {
		// 构建订阅参数
		p := params.NewKrakenParams(params.Unsubscribe, channel, strings...)
		// 发送订阅参数
		if err := k.Send(p.Json()); err != nil {
			if k.logger != nil {
				k.logger.Printf("send subscribe error,%s", err.Error())
			}
			return
		}
	}
}

func (k *KrakenClient) OnDisconnected() {
	//TODO implement me
	if k.logger != nil {
		k.logger.Println("on disconnected")
	}
}

func (k *KrakenClient) OnMessage(data []byte) error {
	//TODO implement me
	var e KrakenEnvelope
	// 转码JSON对象
	if err := json.Unmarshal(data, &e); err != nil {
		return fmt.Errorf("failed to unmarshal message,%s", err.Error())
	}

	if e.IsSubscription() {
		channel := e.GetChannel()

		if callers, ex := k.handler[channel]; ex {
			// 遍历处理字典
			for _, caller := range callers {
				if err := k.pool.Submit(func() {
					if err := caller(&e); err != nil {
						if k.logger != nil {
							k.logger.Printf("[error]failed to handler message,%s,%s", channel, err.Error())
						}
						return
					}
				}); err != nil {
					if k.logger != nil {
						k.logger.Printf("[error]failed to submit task to ants pool,%s", err.Error())
					}
					return err
				}
			}
		}
	}
	return nil
}

func (k *KrakenClient) OnError(err error) {
	//TODO implement me
	if k.logger != nil {
		k.logger.Printf("[error]get error,%s", err.Error())
	}
}
