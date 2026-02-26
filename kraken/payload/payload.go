package payload

import "time"

type Trade struct {
	Symbol    string    `json:"symbol"`
	Side      string    `json:"side"`
	Qty       float64   `json:"qty"`
	Price     float64   `json:"price"`
	OrdType   string    `json:"ord_type"`
	TradeId   int       `json:"trade_id"`
	Timestamp time.Time `json:"timestamp"`
}

type OrderBookItem struct {
	Price float64 `json:"price"`
	Qty   float64 `json:"qty"`
}

type OrderBook struct {
	Symbol    string          `json:"symbol"`
	Bids      []OrderBookItem `json:"bids"`
	Asks      []OrderBookItem `json:"asks"`
	Checksum  int64           `json:"checksum"`
	Timestamp time.Time       `json:"timestamp"`
}

type L3OrderEvent struct {
	Event      string    `json:"event,omitempty"`
	OrderId    string    `json:"order_id"`
	LimitPrice float64   `json:"limit_price"`
	OrderQty   float64   `json:"order_qty"`
	Timestamp  time.Time `json:"timestamp"`
}

type L3BookUpdate struct {
	Checksum int64          `json:"checksum"`
	Symbol   string         `json:"symbol"`
	Bids     []L3OrderEvent `json:"bids"`
	Asks     []L3OrderEvent `json:"asks"`
}

type Ticker struct {
	Symbol    string    `json:"symbol"`
	Bid       float64   `json:"bid"`
	BidQty    float64   `json:"bid_qty"`
	Ask       float64   `json:"ask"`
	AskQty    float64   `json:"ask_qty"`
	Last      float64   `json:"last"`
	Volume    float64   `json:"volume"`
	VWAP      float64   `json:"vwap"`
	Low       float64   `json:"low"`
	High      float64   `json:"high"`
	Change    float64   `json:"change"`
	ChangePct float64   `json:"change_pct"`
	Timestamp time.Time `json:"timestamp"`
}
