package internal

import (
	"context"

	"github.com/simonks2016/dex_plus/okx/param"
)

func (o *OKXClient) OnDisconnecting() {
	// 取消全部订阅
	if err := o.UnsubscribeAll(); err != nil {
		if o.logger != nil {
			o.logger.Printf("[error]Failed to unsubscribe channel:%v",
				err)
		}
		return
	}
	if o.logger != nil {
		o.logger.Printf("[info]Disconnecting from OKX(%s)...",
			o.url)
	}
}
func (o *OKXClient) OnDisconnected() {}
func (o *OKXClient) OnConnected() {

	if o.logger != nil {
		o.logger.Printf("[info]Successfully connected to OKX(%s)",
			o.url,
		)
	}

	if !o.isNeedAuth {
		// 设置已经完成验证
		o.authDone.Swap(true)
		// 发送订阅信息
		o.sendSubscribeChannelMessage()
		return
	} else {
		if o.auth == nil {
			o.client.Close()
			panic("Authentication is required but credentials are not configured.")
		} else {
			ctx, cancel := context.WithTimeout(o.ctx, o.sendTimeOut)
			defer cancel()
			// 生成信息
			data := param.NewLoginParameters(o.auth.ApiKey, o.auth.Passphrase, o.auth.SecretKey)
			// 发送消息
			if err := o.client.Send(ctx, data); err != nil {
				if o.logger != nil {
					o.logger.Printf("[error] Failed to connect to the server:%s",
						err.Error())
				}
				return
			}
		}
	}
}
func (o *OKXClient) OnConnecting(reason string) {}
