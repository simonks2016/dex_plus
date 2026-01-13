package private

import (
	"log"

	"github.com/simonks2016/dex_plus/okx"
	"github.com/simonks2016/dex_plus/okx/param"
)

// SubscribePositionAndBalance 	订阅持仓和余额
// parameters:
// @handler func(posAndBala []okx.PositionAndBalance) error
func (p *Private) SubscribePositionAndBalance(handler func(posAndBala ...okx.PositionAndBalance) error) {

	instType := "ANY"
	channel := "balance_and_position"

	p1 := param.NewSubscribeParameters(
		param.SubscribeChannelParams{
			Channel:  channel,
			InstType: &instType,
		},
	).Encode()

	if err := p.client.SubscribeChannel(p1, channel, func(payload *okx.Payload) error {
		data, err := okx.ParseData[okx.PositionAndBalance](payload)
		if err != nil {
			return err
		}
		return handler(data...)
	}); err != nil {

		return
	}
}

// SubscribeTrade 订阅交易信息
// parameters:
// @handler func(trade okx.Trade) error
func (p *Private) SubscribeTrade(handler func(trade ...okx.TradeFill) error) {

	instType := "ANY"
	channel := "fills"

	p1 := param.NewSubscribeParameters(
		param.SubscribeChannelParams{
			Channel:  channel,
			InstType: &instType,
		},
	).Encode()

	//fills
	if err := p.client.SubscribeChannel(p1, channel, func(payload *okx.Payload) error {
		data, err := okx.ParseData[okx.TradeFill](payload)
		if err != nil {
			return err
		}
		return handler(data...)
	}); err != nil {
		return
	}
}

func (p *Private) SubscribeOrderFilled(handler func(orders ...okx.OrderState) error) {

	instType := "ANY"
	channel := "orders"

	p1 := param.NewSubscribeParameters(
		param.SubscribeChannelParams{
			Channel:  channel,
			InstType: &instType,
		},
	).Encode()

	//order
	if err := p.client.SubscribeChannel(p1, channel, func(payload *okx.Payload) error {
		data, err := okx.ParseData[okx.OrderState](payload)
		if err != nil {
			return err
		}
		return handler(data...)
	}); err != nil {
		return
	}
}

// SetLogger 设置日志记录器
// parameters:
// @logger *log.Logger
func (p *Private) SetLogger(logger *log.Logger) OKXPrivate {
	p.logger = logger
	return p
}

// Connect 连接OKX交易所
func (p *Private) Connect() {
	p.client.Connect()
}

// Close 关闭连接
func (p *Private) Close() {
	p.client.Close()
}

// Reconnect 重新连接
func (p *Private) Reconnect() {
	// 升级
	p.client.Reconnect("")
}

// PlaceOrder 下单接口
func (p *Private) PlaceOrder(placeOrderParam ...param.PlaceOrderParams) error {

	p1 := param.NewParameters[param.PlaceOrderParams](
		"order",
		placeOrderParam...,
	).Encode()

	return p.client.Send(p1)
}

// CancelOrder 撤单接口
func (p *Private) CancelOrder(cancelOrderParam ...param.CancelOrder) error {
	p1 := param.NewParameters[param.CancelOrder](
		"cancel-order",
		cancelOrderParam...,
	).Encode()
	// 发送信息
	return p.client.Send(p1)
}

// AmendOrder 改单
func (p *Private) AmendOrder(params ...param.AmendOrder) error {
	p1 := param.NewParameters[param.AmendOrder](
		"amend-order",
		params...,
	).Encode()
	// 发送信息
	return p.client.Send(p1)
}
