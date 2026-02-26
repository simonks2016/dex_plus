package payload

import "time"

type MatchedTrade struct {
	Type         string    `json:"type"`
	TradeId      int       `json:"trade_id"`
	Sequence     int       `json:"sequence"`
	MakerOrderId string    `json:"maker_order_id"`
	TakerOrderId string    `json:"taker_order_id"`
	Time         time.Time `json:"time"`
	ProductId    string    `json:"product_id"`
	Size         string    `json:"size"`
	Price        string    `json:"price"`
	Side         string    `json:"side"`
}
