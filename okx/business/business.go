package business

import (
	"log"

	"github.com/simonks2016/dex_plus/okx"
	"github.com/simonks2016/dex_plus/okx/param"
)

func (O *Business) SubscribeTradeAll(callback func(trade []okx.RawTrades) error) {
	//TODO implement me
	subscribe[okx.RawTrades]("trades-all", callback, O)
}

func (O *Business) SetLogger(logger *log.Logger) OKXBusiness {
	//TODO implement me
	O.logger = logger
	return O
}

func (O *Business) SetInstId(id ...string) OKXBusiness {
	//TODO implement me
	O.instId = append(O.instId, id...)
	return O
}

func (O *Business) SetInstFamily(id ...string) OKXBusiness {
	//TODO implement me
	O.instFamily = append(O.instFamily, id...)
	return O
}

func (p *Business) Connect() {
	// Connect 连接OKX交易所
	p.client.Connect()
}

// Close 关闭并且取消订阅
func (p *Business) Close() { p.client.Close() }

// Reconnect 重新连接
func (p *Business) Reconnect()           { p.client.Reconnect("") }
func (O *Business) ExchangeName() string { return "okx" }

// subscribe: 通用订阅逻辑
func subscribe[T okx.MarketEvent](channel string, callback func([]T) error, p *Business) {
	caller := func(resp *okx.Payload) error {
		data, err := okx.ParseData[T](resp)
		if err != nil {
			return err
		}
		return callback(data)
	}

	args := p.buildSubscribeArgs(channel)
	payload := param.NewSubscribeParameters(args...).Encode()

	if err := p.client.SubscribeChannel(payload, channel, caller); err != nil {
		if p.logger != nil {
			p.logger.Printf("[ERROR] %s", err)
		}
		return
	}
}

func (p *Business) buildSubscribeArgs(channel string) []param.SubscribeChannelParams {
	var args []param.SubscribeChannelParams

	switch {
	case len(p.instId) > 0:
		// 注意：for range 变量地址问题，用索引或创建新变量
		for i := range p.instId {
			s := p.instId[i]
			args = append(args, param.SubscribeChannelParams{
				Channel: channel,
				InstId:  &s,
			})
		}
	case len(p.instFamily) > 0:
		for i := range p.instFamily {
			s := p.instFamily[i]
			args = append(args, param.SubscribeChannelParams{
				Channel:    channel,
				InstFamily: &s, // 如果你的结构体没有 InstFamily 字段，就删掉这行改成 InstId:&s
			})
		}
	default:
		panic("You need to specify either --instId or --instFamily")
	}
	return args
}
