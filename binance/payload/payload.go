package payload

type AggTrade struct {
	EventType    string `json:"event_type" binance:"e"`
	EventTime    int64  `json:"event_time" binance:"E"`
	Symbol       string `json:"symbol" binance:"s"`
	TradeId      int    `json:"trade_id" binance:"a"`
	Price        string `json:"price" binance:"p"`
	Quantity     string `json:"quantity" binance:"q"`
	FirstTradeId int    `json:"first_trade_id" binance:"f"`
	LastTradeId  int    `json:"last_trade_id" binance:"l"`
	TradeTime    int64  `json:"trade_time" binance:"T"`
	IsMarket     bool   `json:"is_market" binance:"m"`
}

type Trade struct {
	EventType string `json:"event_type" binance:"e"`
	EventTime int64  `json:"event_time" binance:"E"`
	Symbol    string `json:"symbol" binance:"s"`
	TradeId   int    `json:"trade_id" binance:"t"`
	Price     string `json:"price" binance:"p"`
	Quantity  string `json:"quantity" binance:"q"`
	TradeTime int64  `json:"trade_time" binance:"T"`
	IsMarket  bool   `json:"is_market" binance:"m"`
}

type OrderBookSnapshot struct {
	LastUpdateID int64   `json:"last_update_id" binance:"lastUpdateId"`
	Bids         [][]any `json:"bids" binance:"bids"`
	Asks         [][]any `json:"asks" binance:"asks"`
}

type OrderBookDelta struct {
	EventType    string     `json:"event_type" binance:"e"`
	EventTime    int64      `json:"event_time" binance:"E"`
	Symbol       string     `json:"symbol" binance:"s"`
	UpdateId     int        `json:"update_id" binance:"U"`
	LastUpdateId int        `json:"last_update_id" binance:"u"`
	Bids         [][]string `json:"bids" binance:"b"`
	Asks         [][]string `json:"asks" binance:"a"`
}
