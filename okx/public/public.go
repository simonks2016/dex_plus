package public

import (
	"log"

	"github.com/simonks2016/dex_plus/okx"
	"github.com/simonks2016/dex_plus/okx/param"
)

// SetLogger 设置日志记录器
// parameters:
// @logger *log.logger
func (ws *Public) SetLogger(logger *log.Logger) OKXPublic {
	ws.logger = logger
	return ws
}

// SetInstId 设置订阅品种名字
// parameters:
// @id []string 品种ID
func (ws *Public) SetInstId(id ...string) OKXPublic {
	ws.instId = append(ws.instId, id...)
	return ws
}

// SetInstFamily 设置订阅品种类型
// parameters:
// @f []string 品种类型名字
func (ws *Public) SetInstFamily(f ...string) OKXPublic {
	ws.instFamily = append(ws.instFamily, f...)
	return ws
}

// Connect 连接OKX交易所
func (p *Public) Connect() {
	p.client.Connect()
}

// Close 关闭并且取消订阅
func (p *Public) Close() {
	p.client.Close()
}

// Reconnect 重新连接
func (p *Public) Reconnect() {
	// 升级
	p.client.Reconnect("")
}

func (p *Public) buildSubscribeArgs(channel string) []param.SubscribeChannelParams {
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

// subscribe: 通用订阅逻辑
func subscribe[T okx.MarketEvent](channel string, callback func([]T) error, p *Public) {
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

// SubscribeKline 订阅k线频道
func (p *Public) SubscribeKline(channel string, callback func([]okx.Kline) error) {
	subscribe[okx.Kline](channel, callback, p)
}

// SubscribeTrade 订阅公共聚合交易数据
func (p *Public) SubscribeTrade(callback func(trade []okx.AggregatedTrades) error) {
	//TODO implement me
	subscribe[okx.AggregatedTrades](okx.TradesChannel, callback, p)
}

// SubscribeTrade 订阅公共全部交易数据
func (p *Public) SubscribeTradeAll(callback func(trade []okx.RawTrades) error) {
	//TODO implement me
	subscribe[okx.RawTrades](okx.TradesChannel, callback, p)
}

// SubscribeBook 订阅实时盘口数据
func (p *Public) SubscribeBook(channel string, callback func([]okx.OrderBook) error) {
	subscribe[okx.OrderBook](channel, callback, p)
}

// SubscribeTicker 订阅Tick数据行情
func (p *Public) SubscribeTicker(callback func([]okx.Ticker) error) {
	subscribe[okx.Ticker](okx.TickersChannel, callback, p)
}

func (p *Public) ExchangeName() string {
	return "okx"
}
