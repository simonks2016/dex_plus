package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/simonks2016/dex_plus/internal/client"
	"github.com/simonks2016/dex_plus/okx"
	"github.com/simonks2016/dex_plus/okx/param"
)

type OKXClient struct {
	client *client.WsClient
	auth   *Auth
	ctx    context.Context

	//
	logger *log.Logger
	pool   *ants.Pool

	handlerMap map[string][]okx.Caller

	authDone        atomic.Bool
	sendTimeOut     time.Duration
	subscribeParams [][]byte
	isNeedAuth      bool
	url             string
}

func NewOKXClient(ctx context.Context, auth *Auth, cfg *client.Config) *OKXClient {

	cli := &OKXClient{
		client:          client.NewWsClient(ctx, cfg),
		auth:            auth,
		ctx:             ctx,
		sendTimeOut:     cfg.SendTimeout,
		subscribeParams: make([][]byte, 0),
		handlerMap:      make(map[string][]okx.Caller),
		isNeedAuth:      cfg.IsNeedAuth,
		logger:          cfg.Logger,
		url:             cfg.URL,
	}
	cli.client.SetObserver(cli)

	return cli
}

func (o *OKXClient) OnError(err error) {
	//TODO implement me
	if o.logger != nil {
		o.logger.Println(err)
	}
	return
}

func (o *OKXClient) OnMessage(msg []byte) error {

	resp, err := okx.ConvertResponse(msg)
	if err != nil {
		o.logger.Printf("[error] failed to encode payload,%s", err.Error())
		return nil
	}

	// 将消息分流到各个处理单位上
	switch true {
	case resp.IsSubscribe():
		return o.onSubscribe(resp.GetChannel(), resp)
	case resp.IsEvent():
		return o.onEvent(resp.Event, resp)
	case resp.IsOperation():
		return o.onOpEvent(*resp.Op, resp)
	default:
		return nil
	}
}

func (o *OKXClient) onSubscribe(channel string, payload *okx.Payload) error {
	callers, exist := o.handlerMap[channel]
	if !exist || len(callers) == 0 {
		return nil
	}

	// 如果 handlerMap 运行时可能会增删（并发写），建议 copy 一份 slice 再遍历
	local := make([]okx.Caller, len(callers))
	copy(local, callers)

	for _, c := range local {
		caller := c // 避免闭包捕获循环变量

		// 关键：Submit 失败才返回 error；任务执行错误内部记录
		if err := o.pool.Submit(func() {
			// 防止某个 handler panic 把整个 worker 干崩（ants worker 会退出）
			defer func() {
				if r := recover(); r != nil {
					o.OnError(fmt.Errorf("handler panic, channel=%s: %v", channel, r))
				}
			}()
			if err := caller(payload); err != nil {
				o.OnError(err)
			}
		}); err != nil {
			// 只返回“提交失败”的错误（池满 / Nonblocking / 已关闭等）
			return err
		}
	}
	return nil
}

func (o *OKXClient) onEvent(event string, payload *okx.Payload) error {

	switch strings.ToLower(event) {
	case "login":
		if !o.authDone.Load() {
			// 将自动验证设置为真
			o.authDone.Swap(true)
			// 发送订阅信息到OKX
			o.sendSubscribeChannelMessage()
		}
	case "error":
		if o.logger != nil {
			o.logger.Printf("[error] %s", payload.Code)
		}
	case "notice":
		o.client.Reconnect("the okx command we ar reconnect")
	case "subscribe":
		if o.logger != nil {
			o.logger.Printf("[info]Successfully subscribed to “%s”",
				func() string {
					if payload.Arg != nil {
						return payload.Arg.Channel
					}
					return ""
				}())
		}

	}
	return nil
}

func (o *OKXClient) onOpEvent(event string, payload *okx.Payload) error {

	switch strings.ToLower(event) {
	case "order":
		if payload.Code == "0" {
			if o.logger != nil {
				o.logger.Printf("[info] success to submit order \n")
			}
		} else {
			if o.logger != nil {

				var errMsg []string

				d, err := okx.ParseDataToMap(payload.Data)
				if err != nil {
					return err
				}
				for _, m := range d {
					errMsg = append(errMsg, m["sMsg"].(string))
				}
				// 打印错误信息
				o.logger.Printf("[info] failed to submit order,%s \n", strings.Join(errMsg, ","))
			}
		}
		return nil
	default:
		if o.logger != nil {
			o.logger.Printf("[info.%s] recive op message from okx,%s \n", event, string(payload.Data))
		}
		return nil
	}
}

// buildSubParams 复用构建参数：减少重复 lambda + 更清晰
func buildSubParams(channel, instId, instType string) param.SubscribeChannelParams {
	p := param.SubscribeChannelParams{
		Channel: channel,
	}
	if instId != "" {
		p.InstId = &instId
	}
	if instType != "" {
		p.InstType = &instType
	}
	return p
}

func (o *OKXClient) sendWithTimeout(data []byte) error {

	if o.isNeedAuth && !o.authDone.Load() {
		return errors.New("authentication required. Identity verification is enabled for this API, but no valid authentication was found")
	}

	ctx, cancel := context.WithTimeout(o.ctx, o.sendTimeOut)
	defer cancel()
	//
	return o.client.Send(ctx, data)
}

// SubscribeChannel 订阅频道
// Parameters:
// @param string 可选
// @caller []okx.Okx 回调函数
func (o *OKXClient) SubscribeChannel(param []byte, channel string, caller ...okx.Caller) error {

	// 将订阅参数放入到数组里面
	o.subscribeParams = append(o.subscribeParams, param)
	// 将处理函数放入map当中
	o.handlerMap[channel] = append(o.handlerMap[channel], caller...)
	return nil
}

// UnsubscribeAll 取消订阅各个频道
// Parameters:
// @instId string 可选
// @instType string 可选
func (o *OKXClient) UnsubscribeAll() error {

	if len(o.handlerMap) == 0 {
		return nil
	}
	channels := make([]string, 0, len(o.handlerMap))
	for ch := range o.handlerMap {
		channels = append(channels, ch)
	}

	args := make([]param.SubscribeChannelParams, 0, len(channels))
	for _, ch := range channels {
		args = append(args, buildSubParams(ch, "", "ANY"))
	}

	data := param.NewUnsubscribeParameters(args...).Encode()

	if err := o.sendWithTimeout(data); err != nil {
		return err
	}

	for _, ch := range channels {
		delete(o.handlerMap, ch)
	}
	return nil
}

func (o *OKXClient) Send(msg []byte) error {
	return o.sendWithTimeout(msg)
}
func (o *OKXClient) Connect() {

	if o.pool == nil {
		panic("the thread pool is nil")
	}

	o.client.Start()
}
func (o *OKXClient) Close() {
	o.client.Close()
}
func (o *OKXClient) Reconnect(reason string) {
	o.client.Reconnect(reason)
}

func (o *OKXClient) sendSubscribeChannelMessage() {
	// 遍历订阅
	for _, subscribeParam := range o.subscribeParams {
		// 发送订阅信息
		if err := o.sendWithTimeout(subscribeParam); err != nil {
			if o.logger != nil {
				o.logger.Printf("[error] failed to subscribe channel: %v", err.Error())
			}
			return
		}
	}
	return
}

func (o *OKXClient) SetThreadPool(pool *ants.Pool) *OKXClient {
	o.pool = pool
	return o
}
