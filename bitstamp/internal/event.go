package internal

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/simonks2016/dex_plus/bitstamp/params"
)

func (b *BitstampClient) OnConnecting(reason string) {
	//TODO implement me
	fmt.Println(reason)
}

func (b *BitstampClient) OnConnected() {
	//TODO implement me
	fmt.Println("On Connected")
	b.isConnected.Store(true)

	if !b.isRequireAuth {
		b.authDone.Store(true)
		// 循环订阅频道
		for channel, _ := range b.handler {
			p1 := params.NewSubscribeParams(channel)
			// 发送订阅信息
			if err := b.Send(p1.Json()); err != nil {
				if b.logger != nil {
					b.logger.Printf("[error]failed to subscribe channel,%s,error:%s", channel, err.Error())
				}
				return
			}
			time.Sleep(time.Millisecond * time.Duration(5))
		}
	}

}

func (b *BitstampClient) OnDisconnecting() {
	//TODO implement me
	fmt.Println("On Disconnecting")
	// 取消全部订阅
	b.Unsubscribe()
}

func (b *BitstampClient) OnDisconnected() {
	//TODO implement me
	fmt.Println("On Disconnected")
}

func (b *BitstampClient) OnMessage(data []byte) error {
	var result Envelope
	if err := json.Unmarshal(data, &result); err != nil {
		// 建议使用 logger 而不是 fmt.Println，保持日志一致性
		return fmt.Errorf("unmarshal envelope: %w", err)
	}

	// 1. 预处理 Event 字符串，避免多次调用 ToLower
	event := strings.ToLower(result.Event)
	channel := result.Channel

	// 2. 使用 switch 替代多个 if，逻辑更清晰且性能稍好
	switch event {
	case "bts:subscription_succeeded", "bts:unsubscription_succeeded":
		if b.logger != nil {
			b.logger.Printf("[success] %s channel: %s", event, channel)
		}
		return nil

	case "bts:error":
		return b.handleErrorMessage(result.Data)
	}

	// 3. 处理业务消息
	callers, exists := b.handler[channel]
	if !exists {
		return nil
	}

	for _, caller := range callers {
		// 闭包捕获变量优化：在循环内定义局部变量，避免并发竞争（Go 1.22 之前版本尤为重要）
		currentCaller := caller
		// 复制一份 result 的副本传递给异步任务，防止 data 竞争
		currentResult := result

		err := b.pool.Submit(func() {
			if err := currentCaller(&currentResult); err != nil {
				if b.logger != nil {
					b.logger.Printf("[error] failed to parse message of %s, error: %v", event, err)
				}
			}
		})
		if err != nil {
			return fmt.Errorf("submit to pool: %w", err)
		}
	}

	return nil
}

// 提取错误处理逻辑，增加类型断言保护
func (b *BitstampClient) handleErrorMessage(rawData json.RawMessage) error {
	var d1 map[string]any
	if err := json.Unmarshal(rawData, &d1); err != nil {
		return err
	}

	msg, _ := d1["message"].(string) // 安全断言，避免 panic
	if msg == "" {
		msg = "unknown bitstamp error"
	}

	b.OnError(errors.New(msg))
	return nil
}

func (b *BitstampClient) OnError(err error) {
	//TODO implement me
	fmt.Println(err.Error())
}
