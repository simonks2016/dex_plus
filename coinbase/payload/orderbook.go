package payload

import "time"

type OrderBook struct {
	Type      string     `json:"type"`
	ProductId string     `json:"product_id"`
	Bids      [][]string `json:"bids"`
	Asks      [][]string `json:"asks"`
}

type OrderBookUpdate struct {
	Type      string     `json:"type"`
	ProductId string     `json:"product_id"`
	Changes   [][]string `json:"changes"`
	Time      time.Time  `json:"time"`
}
