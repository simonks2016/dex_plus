package payload

import "github.com/goccy/go-json"

type AggTrade struct {
	EventType    string `json:"e" binance:"e"`
	EventTime    int64  `json:"E" binance:"E"`
	Symbol       string `json:"s" binance:"s"`
	TradeId      int    `json:"a" binance:"a"`
	Price        string `json:"p" binance:"p"`
	Quantity     string `json:"q" binance:"q"`
	FirstTradeId int    `json:"f" binance:"f"`
	LastTradeId  int    `json:"l" binance:"l"`
	TradeTime    int64  `json:"T" binance:"T"`
	IsMarket     bool   `json:"m" binance:"m"`
}

type Trade struct {
	EventType string `json:"e" binance:"e"`
	EventTime int64  `json:"E" binance:"E"`
	Symbol    string `json:"s" binance:"s"`
	TradeId   int    `json:"t" binance:"t"`
	Price     string `json:"p" binance:"p"`
	Quantity  string `json:"q" binance:"q"`
	TradeTime int64  `json:"T" binance:"T"`
	IsMarket  bool   `json:"m" binance:"m"`
}

type OrderBookSnapshot struct {
	LastUpdateID int64   `json:"last_update_id" binance:"lastUpdateId"`
	Bids         [][]any `json:"bids" binance:"bids"`
	Asks         [][]any `json:"asks" binance:"asks"`
}

type OrderBookDelta struct {
	EventType    string     `json:"e" binance:"e"`
	EventTime    int64      `json:"E" binance:"E"`
	Symbol       string     `json:"s" binance:"s"`
	UpdateId     int        `json:"U" binance:"U"`
	LastUpdateId int        `json:"u" binance:"u"`
	Bids         [][]string `json:"b" binance:"b"`
	Asks         [][]string `json:"a" binance:"a"`
}

type Stream struct {
	Stream *string         `json:"stream"`
	Data   json.RawMessage `json:"data"`
	Id     *string         `json:"id"`
	Result json.RawMessage `json:"result"`
}
