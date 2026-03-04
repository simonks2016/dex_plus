package payload

import "time"

type OrderBook struct {
	ProductId string    `json:"product_id"`
	Bids      []Level   `json:"bids"`
	Asks      []Level   `json:"asks"`
	Time      time.Time `json:"time"`
}

type OrderBookUpdate struct {
	Type      string     `json:"type"`
	ProductId string     `json:"product_id"`
	Changes   [][]string `json:"changes"`
	Time      time.Time  `json:"time"`
}

type OrderBookSnapshot struct {
	Type      string     `json:"type"`
	ProductId string     `json:"product_id"`
	Bids      [][]string `json:"bids"`
	Asks      [][]string `json:"asks"`
}

type Level struct {
	Price float64 `json:"price"`
	Size  float64 `json:"size"`
}
