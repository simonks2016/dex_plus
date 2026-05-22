package internal

import (
	"fmt"
	"strings"

	"github.com/goccy/go-json"
	"github.com/simonks2016/dex_plus/kraken/params"
	"github.com/simonks2016/dex_plus/kraken/payload"
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
		for channel, s := range k.subscribeRequest {
			if len(s) > 0 {
				// 构建订阅参数
				p := params.NewKrakenParams(params.Subscribe, channel, s...)
				// 发送订阅参数
				if err := k.Send(p.Json()); err != nil {
					if k.logger != nil {
						k.logger.Printf("send subscribe error,%s", err.Error())
					}
					return
				}
			}
		}

		// 订阅instrument频道
		p := params.NewKrakenParams(params.Subscribe, "instrument")
		if err := k.Send(p.Json()); err != nil {
			if k.logger != nil {
				k.logger.Printf("failed to subscribe instrument channel,%s", err.Error())
			}
			return
		}

	}
}

func (k *KrakenClient) OnDisconnecting() {
	//TODO implement me
	if k.logger != nil {
		k.logger.Println("on disconnecting")
	}
	// 全部取消订阅
	for channel, strs := range k.subscribeRequest {
		// 构建订阅参数
		p := params.NewKrakenParams(params.Unsubscribe, channel, strs...)
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
	var e payload.KrakenEnvelope
	// 转码JSON对象
	if err := json.Unmarshal(data, &e); err != nil {
		return fmt.Errorf("failed to unmarshal message,%s", err.Error())
	}

	// 假如是ACK 消息
	if e.IsAck() {
		if e.Method != nil {
			switch strings.ToLower(*e.Method) {
			case "subscribe":
				return k.onSubscribeAck(&e)
			case "unsubscribe":
				return k.onUnsubscribeAck(&e)
			default:
				return fmt.Errorf("invalid method %s", *e.Method)
			}
		}
	}

	if e.IsSubscription() {
		channel := e.GetChannel()

		if strings.EqualFold(channel, "status") {
			// 提交任务处理status
			if err := k.pool.Submit(func() {
				if err := k.onStatus(&e); err != nil {
					if k.logger != nil {
						k.logger.Printf("[error]failed to handler message,channel=%s,error=%s", channel, err.Error())
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

func (k *KrakenClient) onStatus(data *payload.KrakenEnvelope) error {

	var s2 []payload.Status

	if err := json.Unmarshal(data.Data, &s2); err != nil {

		return fmt.Errorf("failed to unmarshal status,error=%s,data=%s", err.Error(), string(data.Data))
	} else {

		for _, status := range s2 {
			if strings.EqualFold(status.System, "online") {
				if k.logger != nil {
					k.logger.Printf("[status]Successfully connected to Kraken. Status: Online. Version: %s. API Version: %s.",
						status.Version, status.ApiVersion)
				}
			} else {
				if k.logger != nil {
					k.logger.Printf("[status]recive status from kraken,%s", string(data.Data))
				}
			}
		}
	}
	return nil
}

func (k *KrakenClient) onInstrument(envelope *payload.KrakenEnvelope) error {

	data, err := payload.ParseSignalData[payload.Instrument](envelope)
	if err != nil {
		return err
	}

	// 添加到类中
	k.instrumentService.AddTradingPairs(data.Pairs...)
	k.instrumentService.AddAsset(data.Assets...)
	return nil
}

func (k *KrakenClient) onSubscribeAck(e *payload.KrakenEnvelope) error {

	channel, _ := e.Result["channel"].(string)
	symbol, _ := e.Result["symbol"].(string)

	if e.Success != nil && *e.Success {
		if symbol != "" {
			k.channelState.Switch(channel, symbol, Subscribed)
		}
		if k.logger != nil {
			k.logger.Printf("[success] Successfully subscribed to channel=%s symbol=%s", channel, symbol)
		}
	} else {
		if symbol != "" {
			k.channelState.Switch(channel, symbol, SubscribeFailed)
		}
		if k.logger != nil {
			errMsg := ""
			if e.Error != nil {
				errMsg = *e.Error
			}
			k.logger.Printf("[error] failed to subscribe channel=%s symbol=%s because %s", channel, symbol, errMsg)
		}
	}
	return nil
}

func (k *KrakenClient) onUnsubscribeAck(e *payload.KrakenEnvelope) error {

	channel, _ := e.Result["channel"].(string)
	symbol, _ := e.Result["symbol"].(string)

	if e.Success != nil && *e.Success {
		if symbol != "" {
			k.channelState.Switch(channel, symbol, Unsubscribed)
		}
		if k.logger != nil {
			k.logger.Printf("[success] Successfully unsubscribed to channel=%s symbol=%s", channel, symbol)
		}
	} else {
		if symbol != "" {
			k.channelState.Switch(channel, symbol, SubscribeFailed)
		}
		if k.logger != nil {
			errMsg := ""
			if e.Error != nil {
				errMsg = *e.Error
			}
			k.logger.Printf("[error] failed to unsubscribe channel=%s symbol=%s because %s", channel, symbol, errMsg)
		}
	}
	return nil
}
