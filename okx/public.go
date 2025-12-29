package okx

import (
	"DexPlus/websocket"
	"context"
	"fmt"
	"log"

	"github.com/panjf2000/ants/v2"
)

type Public struct {
	ws *websocket.WsClient

	callbackMap map[string][]Caller
	pool        *ants.Pool
	logger      *log.Logger
	instId      []string
	instFamily  []string

	msgChan chan []byte
}

func NewPublic(bg context.Context, pool *ants.Pool) *Public {
	return &Public{
		callbackMap: make(map[string][]Caller),
		ws: websocket.NewWsClient(
			bg,
			websocket.NewConfig().WithURL("wss://ws.okx.com:8443/ws/v5/public"),
		),
		pool:    pool,
		msgChan: make(chan []byte, 4000),
	}
}

func (ws *Public) SetLogger(logger *log.Logger) *Public {
	ws.logger = logger
	return ws
}

func (ws *Public) SetInstId(id ...string) *Public {
	ws.instId = append(ws.instId, id...)
	return ws
}
func (ws *Public) SetInstFamily(f ...string) *Public {
	ws.instFamily = append(ws.instFamily, f...)
	return ws
}

func (p *Public) Connect() error {

	p.ws.SetHandler(p.handleMessage)

	// 异步启动线程池，分批写入
	go func() {
		for {
			select {
			case msg := <-p.msgChan:
				resp, err := ConvertResponse(msg)
				if err != nil {
					p.handlerError(err)
					return
				}

				if callers, exist := p.callbackMap[resp.GetChannel()]; exist {
					for _, caller := range callers {
						if err = p.pool.Submit(func() {
							if err = caller(resp); err != nil {
								p.handlerError(err)
								return
							}
						}); err != nil {
							p.handlerError(err)
							return
						}
					}
				}
			}
		}
	}()

	// 设置重新连接提醒和重连之前工作
	p.ws.SetReconnectNotify(func(s string) {
		var channels []string
		// 将处理函数中全部订阅
		for k, _ := range p.callbackMap {
			channels = append(channels, k)
		}
		p.subscribeChannel(channels...)
	}, func() {
		p.unsubscribe()
	})
	// 连接WS
	p.ws.Connect()
	return nil
}

func (p *Public) handleMessage(message []byte) error {
	p.msgChan <- message
	return nil
}

func (p *Public) subscribeChannel(channel ...string) {

	var args []*Arg
	for _, s := range channel {
		if len(p.instId) > 0 {
			for _, i := range p.instId {
				args = append(args, NewInstIdArg(i, s))
			}
		} else if len(p.instFamily) > 0 {
			for _, i := range p.instFamily {
				args = append(args, NewInstIdArg(i, s))
			}
		} else {
			panic("Not Set the inst id or family")
		}
	}

	param := NewSubscribeParameters(args...)

	if !p.ws.IsConnecting() {
		p.ws.Connect()
	}
	// 发送信息
	p.ws.WriteMessage(param.Encode())
	return
}

func (p *Public) SubscribeKline(channel string, callback func(kline []Kline) error) {

	var caller = func(resp *Payload) error {
		data, err := ParseData[Kline](resp)
		if err != nil {
			return err
		}
		return callback(data)
	}

	if callbacks, ok := p.callbackMap[channel]; !ok {
		p.callbackMap[channel] = []Caller{caller}
	} else {
		callbacks = append(callbacks, caller)
	}
	p.subscribeChannel(channel)
}
func (p *Public) SubscribeTrade(callback func(trade []Trade) error) {

	var caller = func(resp *Payload) error {
		data, err := ParseData[Trade](resp)
		if err != nil {
			return err
		}
		return callback(data)
	}

	if callbacks, ok := p.callbackMap[TradesChannel]; !ok {
		p.callbackMap[TradesChannel] = []Caller{caller}
	} else {
		callbacks = append(callbacks, caller)
	}
	p.subscribeChannel(TradesChannel)
}
func (p *Public) SubscribeBook(channel string, callback func(books []OrderBook) error) {

	var caller = func(resp *Payload) error {
		data, err := ParseData[OrderBook](resp)
		if err != nil {
			return err
		}
		return callback(data)
	}

	if callbacks, ok := p.callbackMap[channel]; !ok {
		p.callbackMap[channel] = []Caller{caller}
	} else {
		callbacks = append(callbacks, caller)
	}
	p.subscribeChannel(channel)
}
func (p *Public) SubscribeTicker(callback func(tickers []Ticker) error) {
	var caller = func(resp *Payload) error {
		data, err := ParseData[Ticker](resp)
		if err != nil {
			return err
		}
		return callback(data)
	}

	if callbacks, ok := p.callbackMap[TickersChannel]; !ok {
		p.callbackMap[TickersChannel] = []Caller{caller}
	} else {
		callbacks = append(callbacks, caller)
	}
	p.subscribeChannel(TickersChannel)
}

func (p *Public) handlerError(err error) {
	log.Println(err)
}

func (p *Public) unsubscribe() {

	var channels []string

	for k, _ := range p.callbackMap {
		channels = append(channels, k)
	}

	var args []*Arg
	for _, s := range channels {
		if len(p.instId) > 0 {
			for _, i := range p.instId {
				args = append(args, NewInstIdArg(i, s))
			}
		} else if len(p.instFamily) > 0 {
			for _, i := range p.instFamily {
				args = append(args, NewInstIdArg(i, s))
			}
		} else {
			panic("Not Set the inst id or family")
		}
	}
	// 创建新的取消订阅参数
	param := NewUnsubscribeParameters(args...)

	// 发送信息
	p.ws.WriteMessage(param.Encode())
	return
}

func (p *Public) Close() {
	// 取消订阅
	p.unsubscribe()
	//
	fmt.Println("Success To Unsubscribe")
	// 关闭连接
	p.ws.Close()

}
