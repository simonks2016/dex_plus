package okx

import (
	"encoding/json"
	"strings"
)

type Payload struct {
	Event  string          `json:"event"`
	Id     string          `json:"id"`
	Arg    Arg             `json:"arg"`
	Code   string          `json:"code"`
	Msg    string          `json:"msg"`
	ConnId string          `json:"connId"`
	Data   json.RawMessage `json:"data"`
}

type Arg struct {
	Channel    string  `json:"channel"`
	InstId     *string `json:"instId,omitempty"`
	InstFamily *string `json:"instFamily,omitempty"`
}

func (o *Payload) IsNotice() bool {
	return strings.ToLower(o.Event) == "notice" || o.Code == "64008"
}
func (o *Payload) GetChannel() string {
	return o.Arg.Channel
}
func (o *Payload) IsError() bool {
	return strings.ToLower(o.Event) == "error"
}

type CallAuctionDetails struct {
	InstId         string `json:"instId"`
	EqPx           string `json:"eqPx"`
	MatchedSz      string `json:"matchedSz"`
	UnmatchedSz    string `json:"unmatchedSz"`
	AuctionEndTime string `json:"auctionEndTime"`
	State          string `json:"state"`
	Ts             string `json:"ts"`
}

type OrderBook struct {
	Asks      [][]string `json:"asks"`
	Bids      [][]string `json:"bids"`
	Ts        string     `json:"ts"`
	Checksum  int64      `json:"checksum"`
	PrevSeqId int64      `json:"prevSeqId"`
	SeqId     int64      `json:"seqId"`
}

func (o *OrderBook) GetAsks() []map[string]string {

	var response []map[string]string

	for _, ask := range o.Asks {
		if len(ask) >= 4 {
			response = append(response, map[string]string{
				"price":      ask[0],
				"size":       ask[1],
				"order_size": ask[3],
			})
		}
	}
	return response
}
func (o *OrderBook) GetBids() []map[string]string {

	var response []map[string]string

	for _, bid := range o.Bids {
		if len(bid) >= 4 {
			response = append(response, map[string]string{
				"price":      bid[0],
				"size":       bid[1],
				"order_size": bid[3],
			})
		}
	}
	return response
}

type Trade struct {
	InstId  string `json:"instId"`
	TradeId string `json:"tradeId"`
	Ts      string `json:"ts"`
	Px      string `json:"px"`
	Sz      string `json:"sz"`
	Side    string `json:"side"`
	Count   string `json:"count"`
	Source  string `json:"source"`
	SeqId   int64  `json:"seqId"`
}

type Kline struct {
	Ts          string `json:"ts"`
	OpenPrice   string `json:"open_price"`
	HighPrice   string `json:"high_price"`
	LowPrice    string `json:"low_price"`
	ClosePrice  string `json:"close_price"`
	Volume      string `json:"volume"`
	VolCcy      string `json:"vol_ccy"`
	VolCcyQuote string `json:"vol_ccy_quote"`
	Confirm     string `json:"confirm"`
}

func DecodeOKXLine(raws ...[]string) ([]Kline, error) {

	var klines []Kline

	for _, raw := range raws {

		klines = append(klines, Kline{
			Ts:          raw[0],
			OpenPrice:   raw[1],
			HighPrice:   raw[2],
			LowPrice:    raw[3],
			ClosePrice:  raw[4],
			Volume:      raw[5],
			VolCcy:      raw[6],
			VolCcyQuote: raw[7],
			Confirm:     raw[8],
		})

	}
	return klines, nil
}

type Ticker struct {
	InstId    string `json:"instId"`
	InstType  string `json:"instType"`
	Last      string `json:"last"`
	LastSz    string `json:"lastSz"`
	AskPx     string `json:"askPx"`
	AskSz     string `json:"askSz"`
	BidPx     string `json:"bidPx"`
	BidSz     string `json:"bidSz"`
	Open24h   string `json:"open24h"`
	High24h   string `json:"high24h"`
	Low24h    string `json:"low24h"`
	VolCcy24h string `json:"volCcy24h"`
	Vol24h    string `json:"vol24h"`
	SodUtc0   string `json:"sodUtc0"`
	SodUtc8   string `json:"sodUtc8"`
	Ts        string `json:"ts"`
}

func ConvertResponse(msg []byte) (*Payload, error) {

	var response Payload

	err := json.Unmarshal(msg, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
