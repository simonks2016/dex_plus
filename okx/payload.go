package okx

import (
	"encoding/json"
	"strings"
)

type Payload struct {
	Event   string          `json:"event"`
	Id      string          `json:"id"`
	Arg     *Arg            `json:"arg,omitempty"`
	Code    string          `json:"code"`
	Msg     string          `json:"msg"`
	ConnId  string          `json:"connId"`
	Data    json.RawMessage `json:"data"`
	Op      *string         `json:"op,omitempty"`
	InTime  *string         `json:"inTime,omitempty"`
	OutTime *string         `json:"outTime,omitempty"`
}

type Arg struct {
	Channel    string  `json:"channel"`
	InstId     *string `json:"instId,omitempty"`
	InstFamily *string `json:"instFamily,omitempty"`
}

func (o *Payload) IsSubscribe() bool {

	if o.Arg == nil {
		return false
	}
	return len(o.Arg.Channel) > 0
}

func (o *Payload) IsEvent() bool {
	return len(o.Event) > 0
}

func (o *Payload) IsOperation() bool {
	return o.Op != nil && len(*o.Op) > 0
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
	InstId    string     `json:"instId"`
}

type BookLevel struct {
	Price     string `json:"price"`
	Size      string `json:"size"`
	OrderSize string `json:"orderSize"`
}

func (o *OrderBook) GetAsks() []BookLevel {

	var response = make([]BookLevel, 0, len(o.Asks))

	for _, ask := range o.Asks {
		if len(ask) >= 4 {
			response = append(response, BookLevel{
				Price:     ask[0],
				Size:      ask[1],
				OrderSize: ask[3],
			})
		}
	}
	return response
}
func (o *OrderBook) GetBids() []BookLevel {

	var response []BookLevel

	for _, bid := range o.Bids {
		if len(bid) >= 4 {
			response = append(response, BookLevel{
				Price:     bid[0],
				Size:      bid[1],
				OrderSize: bid[3],
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

func ParseDataToMap(data json.RawMessage) ([]map[string]any, error) {

	var d1 []map[string]any

	if err := json.Unmarshal(data, &d1); err != nil {
		return nil, err
	} else {
		return d1, nil
	}
}

type PositionAndBalance struct {
	PTime     string    `json:"pTime"`
	EventType string    `json:"eventType"`
	BalData   []BalData `json:"balData"`
	PosData   []PosData `json:"posData"`
	Trades    []struct {
		InstId  string `json:"instId"`
		TradeId string `json:"tradeId"`
	} `json:"trades"`
}

type BalData struct {
	Ccy     string `json:"ccy"`
	CashBal string `json:"cashBal"`
	UTime   string `json:"uTime"`
}

type PosData struct {
	PosId          string `json:"posId"`
	TradeId        string `json:"tradeId"`
	InstId         string `json:"instId"`
	InstType       string `json:"instType"`
	MgnMode        string `json:"mgnMode"`
	AvgPx          string `json:"avgPx"`
	Ccy            string `json:"ccy"`
	PosSide        string `json:"posSide"`
	Pos            string `json:"pos"`
	PosCcy         string `json:"posCcy"`
	NonSettleAvgPx string `json:"nonSettleAvgPx"`
	SettledPnl     string `json:"settledPnl"`
	UTime          string `json:"uTime"`
}

type Position struct {
	InstId                 string `json:"instId"`
	InstType               string `json:"instType"`
	Adl                    string `json:"adl"`
	AvailPos               string `json:"availPos"`
	AvgPx                  string `json:"avgPx"`
	CTime                  string `json:"cTime"`
	Ccy                    string `json:"ccy"`
	DeltaBS                string `json:"deltaBS"`
	DeltaPA                string `json:"deltaPA"`
	GammaBS                string `json:"gammaBS"`
	GammaPA                string `json:"gammaPA"`
	HedgedPos              string `json:"hedgedPos"`
	Imr                    string `json:"imr"`
	Interest               string `json:"interest"`
	IdxPx                  string `json:"idxPx"`
	Last                   string `json:"last"`
	Lever                  string `json:"lever"`
	Liab                   string `json:"liab"`
	LiabCcy                string `json:"liabCcy"`
	LiqPx                  string `json:"liqPx"`
	MarkPx                 string `json:"markPx"`
	Margin                 string `json:"margin"`
	MgnMode                string `json:"mgnMode"`
	MgnRatio               string `json:"mgnRatio"`
	Mmr                    string `json:"mmr"`
	NotionalUsd            string `json:"notionalUsd"`
	OptVal                 string `json:"optVal"`
	PTime                  string `json:"pTime"`
	PendingCloseOrdLiabVal string `json:"pendingCloseOrdLiabVal"`
	Pos                    string `json:"pos"`
	BaseBorrowed           string `json:"baseBorrowed"`
	BaseInterest           string `json:"baseInterest"`
	QuoteBorrowed          string `json:"quoteBorrowed"`
	QuoteInterest          string `json:"quoteInterest"`
	PosCcy                 string `json:"posCcy"`
	PosId                  string `json:"posId"`
	PosSide                string `json:"posSide"`
	SpotInUseAmt           string `json:"spotInUseAmt"`
	ClSpotInUseAmt         string `json:"clSpotInUseAmt"`
	MaxSpotInUseAmt        string `json:"maxSpotInUseAmt"`
	SpotInUseCcy           string `json:"spotInUseCcy"`
	BizRefId               string `json:"bizRefId"`
	BizRefType             string `json:"bizRefType"`
	ThetaBS                string `json:"thetaBS"`
	ThetaPA                string `json:"thetaPA"`
	TradeId                string `json:"tradeId"`
	UTime                  string `json:"uTime"`
	Upl                    string `json:"upl"`
	UplLastPx              string `json:"uplLastPx"`
	UplRatio               string `json:"uplRatio"`
	UplRatioLastPx         string `json:"uplRatioLastPx"`
	VegaBS                 string `json:"vegaBS"`
	VegaPA                 string `json:"vegaPA"`
	RealizedPnl            string `json:"realizedPnl"`
	Pnl                    string `json:"pnl"`
	Fee                    string `json:"fee"`
	FundingFee             string `json:"fundingFee"`
	LiqPenalty             string `json:"liqPenalty"`
	NonSettleAvgPx         string `json:"nonSettleAvgPx"`
	SettledPnl             string `json:"settledPnl"`

	CloseOrderAlgo []CloseOrderAlgo `json:"closeOrderAlgo"`
}

type CloseOrderAlgo struct {
	AlgoId          string `json:"algoId"`
	SlTriggerPx     string `json:"slTriggerPx"`
	SlTriggerPxType string `json:"slTriggerPxType"`
	TpTriggerPx     string `json:"tpTriggerPx"`
	TpTriggerPxType string `json:"tpTriggerPxType"`
	CloseFraction   string `json:"closeFraction"`
}

type TradeFill struct {
	InstId   string `json:"instId"`   // 产品ID，如 BTC-USDT-SWAP
	FillSz   string `json:"fillSz"`   // 成交数量
	FillPx   string `json:"fillPx"`   // 成交价格
	Side     string `json:"side"`     // buy / sell
	Ts       string `json:"ts"`       // 成交时间戳（毫秒）
	OrdId    string `json:"ordId"`    // 订单ID
	ClOrdId  string `json:"clOrdId"`  // 客户自定义订单ID
	TradeId  string `json:"tradeId"`  // 成交ID
	ExecType string `json:"execType"` // 执行类型（T = Trade）
	Count    string `json:"count"`    // 成交笔数
}

type LinkedAlgoOrd struct {
	AlgoId string `json:"algoId"`
}

type OrderState struct {
	AccFillSz         string        `json:"accFillSz"`
	AmendResult       string        `json:"amendResult"`
	AvgPx             string        `json:"avgPx"`
	CTime             string        `json:"cTime"`
	Category          string        `json:"category"`
	Ccy               string        `json:"ccy"`
	ClOrdId           string        `json:"clOrdId"`
	Code              string        `json:"code"`
	ExecType          string        `json:"execType"`
	Fee               string        `json:"fee"`
	FeeCcy            string        `json:"feeCcy"`
	FillFee           string        `json:"fillFee"`
	FillFeeCcy        string        `json:"fillFeeCcy"`
	FillNotionalUsd   string        `json:"fillNotionalUsd"`
	FillPx            string        `json:"fillPx"`
	FillSz            string        `json:"fillSz"`
	FillPnl           string        `json:"fillPnl"`
	FillTime          string        `json:"fillTime"`
	FillPxVol         string        `json:"fillPxVol"`
	FillPxUsd         string        `json:"fillPxUsd"`
	FillMarkVol       string        `json:"fillMarkVol"`
	FillFwdPx         string        `json:"fillFwdPx"`
	FillMarkPx        string        `json:"fillMarkPx"`
	FillIdxPx         string        `json:"fillIdxPx"`
	InstId            string        `json:"instId"`
	InstType          string        `json:"instType"`
	Lever             string        `json:"lever"`
	Msg               string        `json:"msg"`
	NotionalUsd       string        `json:"notionalUsd"`
	OrdId             string        `json:"ordId"`
	OrdType           string        `json:"ordType"`
	Pnl               string        `json:"pnl"`
	PosSide           string        `json:"posSide"`
	Px                string        `json:"px"`
	PxUsd             string        `json:"pxUsd"`
	PxVol             string        `json:"pxVol"`
	PxType            string        `json:"pxType"`
	Rebate            string        `json:"rebate"`
	RebateCcy         string        `json:"rebateCcy"`
	ReduceOnly        string        `json:"reduceOnly"`
	ReqId             string        `json:"reqId"`
	Side              string        `json:"side"`
	AttachAlgoClOrdId string        `json:"attachAlgoClOrdId"`
	SlOrdPx           string        `json:"slOrdPx"`
	SlTriggerPx       string        `json:"slTriggerPx"`
	SlTriggerPxType   string        `json:"slTriggerPxType"`
	Source            string        `json:"source"`
	State             string        `json:"state"`
	StpId             string        `json:"stpId"`
	StpMode           string        `json:"stpMode"`
	Sz                string        `json:"sz"`
	Tag               string        `json:"tag"`
	TdMode            string        `json:"tdMode"`
	TgtCcy            string        `json:"tgtCcy"`
	TpOrdPx           string        `json:"tpOrdPx"`
	TpTriggerPx       string        `json:"tpTriggerPx"`
	TpTriggerPxType   string        `json:"tpTriggerPxType"`
	TradeId           string        `json:"tradeId"`
	TradeQuoteCcy     string        `json:"tradeQuoteCcy"`
	LastPx            string        `json:"lastPx"`
	QuickMgnType      string        `json:"quickMgnType"`
	AlgoClOrdId       string        `json:"algoClOrdId"`
	AttachAlgoOrds    []any         `json:"attachAlgoOrds"` // 这里是空数组，如果有具体结构可替换
	AlgoId            string        `json:"algoId"`
	AmendSource       string        `json:"amendSource"`
	CancelSource      string        `json:"cancelSource"`
	IsTpLimit         string        `json:"isTpLimit"`
	UTime             string        `json:"uTime"`
	LinkedAlgoOrd     LinkedAlgoOrd `json:"linkedAlgoOrd"`
}
