package okx

import (
	"crypto/md5"
	"fmt"

	"github.com/google/uuid"
	"github.com/simonks2016/dex_plus/okx/param"
)

type OrderBuilder struct {
	Side          string  `json:"side"`
	Size          string  `json:"size"`
	Price         string  `json:"price"`
	InstId        *string `json:"instId"`
	InstCode      *int    `json:"instCode"`
	OrderType     string  `json:"orderType"`
	TdMode        string  `json:"tdMode"`
	Ccy           *string `json:"ccy"`
	ClientOrderId string  `json:"clientOrderId"`
	PosSide       string  `json:"posSide"`
}

func NewOrderBuilder() *OrderBuilder {
	return &OrderBuilder{
		PosSide: "net",
	}
}

func (o *OrderBuilder) Buy() *OrderBuilder {
	o.Side = "buy"
	return o
}
func (o *OrderBuilder) Sell() *OrderBuilder {
	o.Side = "sell"
	return o
}
func (o *OrderBuilder) OnMarketOrder() *OrderBuilder {
	o.Price = "-1"
	o.OrderType = "market"
	return o
}
func (o *OrderBuilder) OnPrice(price float64) *OrderBuilder {
	o.Price = fmt.Sprintf("%.2f", price)
	o.OrderType = "limit"
	return o
}

func (o *OrderBuilder) OnSize(size float64) *OrderBuilder {
	if size < 0.01 {
		panic("the size must be greater than 0.01")
	}
	o.Size = fmt.Sprintf("%.2f", size)
	return o
}
func (o *OrderBuilder) OnInstId(instId string) *OrderBuilder {
	o.InstId = &instId
	return o
}
func (o *OrderBuilder) OnInstCode(code int) *OrderBuilder {
	o.InstCode = &code
	return o
}

func (o *OrderBuilder) OnCross(ccy string) *OrderBuilder {
	o.TdMode = "cross"
	o.Ccy = &ccy
	return o
}
func (o *OrderBuilder) OnCash() *OrderBuilder {
	o.TdMode = "cash"
	o.Ccy = nil
	return o
}
func (o *OrderBuilder) OnIsolated() *OrderBuilder {
	o.TdMode = "isolated"
	return o
}
func (o *OrderBuilder) SetOrderId() *OrderBuilder {

	u := uuid.New().String()

	m5 := md5.New()
	m5.Write([]byte(u))

	o.ClientOrderId = fmt.Sprintf("%x", m5.Sum(nil))
	return o
}

func (o *OrderBuilder) Build() param.PlaceOrderParams {

	if o.InstCode == nil && o.InstId == nil {
		panic("InstCode or InstId is required")
	}

	return param.PlaceOrderParams{
		InstIdCode: o.InstCode,
		InstId:     o.InstId,
		TdMode:     o.TdMode,
		Ccy:        o.Ccy,
		ClOrdId:    &o.ClientOrderId,
		Tag:        nil,
		Side:       o.Side,
		PosSide:    &o.PosSide,
		OrdType:    o.OrderType,
		SZ:         o.Size,
		Px: func() *string {
			if o.Price == "-1" {
				return nil
			}
			return &o.Price
		}(),
		PxUSD:       nil,
		PxVol:       nil,
		ReduceOnly:  nil,
		TgtCcy:      nil,
		BanAmend:    nil,
		PxAmendType: nil,
		StpMode:     nil,
	}

}
