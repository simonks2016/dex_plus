package internal

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/simonks2016/dex_plus/okx"
	"github.com/simonks2016/dex_plus/okx/param"
	"github.com/simonks2016/dex_plus/option"
	"github.com/simonks2016/dex_plus/websocket"
)

type OKXClient struct {
	client     *websocket.WsClient
	auth       *Auth
	ctx        context.Context
	logger     *log.Logger
	handlerMap map[string][]okx.Caller
	pool       *ants.Pool
	authDone   atomic.Pointer[func(error)]
}

func NewOKXClient(ctx context.Context, auth *Auth, cfg *websocket.Config, opts ...option.Option) *OKXClient {

	var pool *ants.Pool

	if val := option.GetOption("read_buffer_size", opts...); val != nil {
		if v, ok := val.(int64); ok {
			cfg.ReadBufferSize = int(v)
		}
	}
	if val := option.GetOption("write_buffer_size", opts...); val != nil {
		if v, ok := val.(int64); ok {
			cfg.WriteBufferSize = int(v)
		}
	}
	// 获取链接URL
	if val := option.GetOption("url", opts...); val != nil {
		cfg.URL = val.(string)
	}
	// 获取是否禁止IPv6
	if val := option.GetOption("forbid_ipv6", opts...); val != nil {
		if v, ok := val.(bool); ok {
			cfg.IsForbidIPV6 = v
		}
	} // 获取新的线程池
	if p1 := option.GetOption("thread_pool", opts...); p1 != nil {
		if p, ok := p1.(*ants.Pool); !ok {
			panic("Need Set the thread pool")
		} else {
			pool = p
		}
	}
	// 获取日志记录器
	if val := option.GetOption("logger", opts...); val != nil {
		if v, ok := val.(*log.Logger); ok {
			cfg.Logger = v
		}
	}

	cli := &OKXClient{
		client: websocket.NewWsClient(ctx, cfg),
		auth:   auth,
		ctx:    ctx,
		pool:   pool,
	}
	cli.client.SetObserver(cli)
	cli.logger = cfg.Logger

	return cli
}
func (o *OKXClient) OnDisconnecting() {
	// 取消全部订阅
	if err := o.UnsubscribeAll(); err != nil {
		if o.logger != nil {
			o.logger.Printf("[error] error on disconnecting client: %v", err)
		}
		return
	}
	if o.logger != nil {
		o.logger.Println("[disconnecting]")
	}
}
func (o *OKXClient) OnDisconnected()            {}
func (o *OKXClient) OnConnected()               {}
func (o *OKXClient) OnConnecting(reason string) {}
func (o *OKXClient) OnAuth(sender func([]byte) error, done func(error)) {

	if o.auth == nil {
		_ = sender(nil)
		done(nil)
		return
	} else {
		// 生成信息
		data := param.NewLoginParameters(o.auth.ApiKey, o.auth.Passphrase, o.auth.SecretKey)
		if err := sender(data); err != nil {
			o.logger.Println(err)
			return
		}
		// 将验证之后回调函数保存下来
		o.setAuthDone(done)
	}
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
		if v := o.authDone.Swap(nil); v != nil {
			callback := *v
			//
			if payload.Code == "0" {
				callback(nil)
				if o.logger != nil {
					o.logger.Printf("[okx] successfully logged in, code=%s", payload.Code)
				}
			}
		}
	case "error":
		if o.logger != nil {
			o.logger.Printf("[error] %s", payload.Code)
		}
	case "notice":
		o.client.Reconnect("the okx command we ar reconnect")
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

func (o *OKXClient) setAuthDone(done func(error)) {
	o.authDone.Store(&done)
}

// 复用构建参数：减少重复 lambda + 更清晰
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
	ctx, cancel := context.WithTimeout(o.ctx, time.Second*5)
	defer cancel()
	return o.client.Send(ctx, data)
}

// SubscribeChannel：
// Parameters:
// @param string 可选
// @caller []okx.Okx 回调函数
func (o *OKXClient) SubscribeChannel(param []byte, channel string, caller ...okx.Caller) error {

	if err := o.sendWithTimeout(param); err != nil {
		return err
	}
	if len(caller) == 0 {
		return nil
	}

	if o.handlerMap == nil {
		o.handlerMap = make(map[string][]okx.Caller, 8)
	}
	o.handlerMap[channel] = append(o.handlerMap[channel], caller...)
	return nil
}

// UnsubscribeAll：
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
	o.client.Start()
}
func (o *OKXClient) Close() {
	o.client.Close()
}
func (o *OKXClient) Reconnect(reason string) {
	o.client.Reconnect(reason)
}

func (o *OKXClient) WithAntsPool(pool *ants.Pool) *OKXClient {
	o.pool = pool
	return o
}
