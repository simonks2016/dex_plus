package internal

import (
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/simonks2016/dex_plus/binance/payload"
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
	b.isConnected.Store(false)
	if b.logger != nil {
		b.logger.Printf("Binance client OnDisconnected")
	}
}

func (b *BinanceClient) OnMessage(data []byte) error {
	//TODO implement me

	var streams payload.Stream

	if err := json.Unmarshal(data, &streams); err != nil {
		if b.logger != nil {
			b.logger.Printf("[error] unmarshal json:%v", err)
		}
		return err
	}

	if streams.Id != nil {
		if b.logger != nil {
			b.logger.Printf("[success] Successfuly Subscribe Channel")
		}
		return nil
	}

	if streams.Stream != nil {

		streamName := *streams.Stream

		// 假如等于服务器即将关闭
		if strings.ToLower(streamName) == "serverShutdown" {
			if b.logger != nil {
				b.logger.Printf("[error] Binance server shutdown,date_time=%s", time.Now().Format("2006-01-02 15:04:05"))
			}
			b.client.Reconnect("The server is shutting down")
			return nil
		}
		// 分析Stream
		s := ParseStreamName(streamName)
		// 分析出ChannelName
		channelName := s[1]
		symbol := s[0]

		// 查看一下处理函数
		if callers, ex := b.handlerMap[channelName]; ex {
			for _, callback := range callers {
				if err := b.pool.Submit(func() {
					if err := callback(symbol, streams.Data); err != nil {
						if b.logger != nil {
							b.logger.Printf("[error]Failed to read message:%s,%s,%s", err.Error(), channelName, *streams.Stream)
						}
						return
					}
				}); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (b *BinanceClient) OnError(err error) {
	//TODO implement me
	if b.logger != nil {
		b.logger.Printf("Binance client OnError:%v", err.Error())
	}
}
