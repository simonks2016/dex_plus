package response

type PendingOrder = Order
type OrderStatus = Order

type Order struct {
	InstType           string          `json:"instType"`
	InstId             string          `json:"instId"`
	TgtCcy             string          `json:"tgtCcy"`
	Ccy                string          `json:"ccy"`
	OrdId              string          `json:"ordId"`
	ClOrdId            string          `json:"clOrdId"`
	Tag                string          `json:"tag"`
	Px                 string          `json:"px"`
	PxUsd              string          `json:"pxUsd"`
	PxVol              string          `json:"pxVol"`
	PxType             string          `json:"pxType"`
	Sz                 string          `json:"sz"`
	Pnl                string          `json:"pnl"`
	OrdType            string          `json:"ordType"`
	Side               string          `json:"side"`
	PosSide            string          `json:"posSide"`
	TdMode             string          `json:"tdMode"`
	AccFillSz          string          `json:"accFillSz"`
	FillPx             string          `json:"fillPx"`
	TradeId            string          `json:"tradeId"`
	FillSz             string          `json:"fillSz"`
	FillTime           string          `json:"fillTime"`
	AvgPx              string          `json:"avgPx"`
	State              string          `json:"state"`
	Lever              string          `json:"lever"`
	AttachAlgoClOrdId  string          `json:"attachAlgoClOrdId"`
	TpTriggerPx        string          `json:"tpTriggerPx"`
	TpTriggerPxType    string          `json:"tpTriggerPxType"`
	TpOrdPx            string          `json:"tpOrdPx"`
	SlTriggerPx        string          `json:"slTriggerPx"`
	SlTriggerPxType    string          `json:"slTriggerPxType"`
	SlOrdPx            string          `json:"slOrdPx"`
	AttachAlgoOrds     []AttachAlgoOrd `json:"attachAlgoOrds"`
	LinkedAlgoOrd      *LinkedAlgoOrd  `json:"linkedAlgoOrd,omitempty"`
	StpId              string          `json:"stpId"`
	StpMode            string          `json:"stpMode"`
	FeeCcy             string          `json:"feeCcy"`
	Fee                string          `json:"fee"`
	RebateCcy          string          `json:"rebateCcy"`
	Rebate             string          `json:"rebate"`
	Source             string          `json:"source"`
	Category           string          `json:"category"`
	ReduceOnly         string          `json:"reduceOnly"`
	CancelSource       string          `json:"cancelSource"`
	CancelSourceReason string          `json:"cancelSourceReason"`
	QuickMgnType       string          `json:"quickMgnType"`
	AlgoClOrdId        string          `json:"algoClOrdId"`
	AlgoId             string          `json:"algoId"`
	IsTpLimit          string          `json:"isTpLimit"`
	UTime              string          `json:"uTime"`
	CTime              string          `json:"cTime"`
	TradeQuoteCcy      string          `json:"tradeQuoteCcy"`
	Outcome            string          `json:"outcome"`
}

type AttachAlgoOrd struct {
	AttachAlgoId         string `json:"attachAlgoId"`
	AttachAlgoClOrdId    string `json:"attachAlgoClOrdId"`
	TpOrdKind            string `json:"tpOrdKind"`
	TpTriggerPx          string `json:"tpTriggerPx"`
	TpTriggerRatio       string `json:"tpTriggerRatio"`
	TpTriggerPxType      string `json:"tpTriggerPxType"`
	TpOrdPx              string `json:"tpOrdPx"`
	SlTriggerPx          string `json:"slTriggerPx"`
	SlTriggerRatio       string `json:"slTriggerRatio"`
	SlTriggerPxType      string `json:"slTriggerPxType"`
	SlOrdPx              string `json:"slOrdPx"`
	Sz                   string `json:"sz"`
	AmendPxOnTriggerType string `json:"amendPxOnTriggerType"`
	CallbackRatio        string `json:"callbackRatio"`
	CallbackSpread       string `json:"callbackSpread"`
	ActivePx             string `json:"activePx"`
	FailCode             string `json:"failCode"`
	FailReason           string `json:"failReason"`
}

type LinkedAlgoOrd struct {
	AlgoId string `json:"algoId"`
}

type ResultAsPlaceOrder struct {
	ClOrdId string `json:"clOrdId"`
	OrdId   string `json:"ordId"`
	Tag     string `json:"tag"`
	Ts      string `json:"ts"`
	SCode   string `json:"sCode"`
	SMsg    string `json:"sMsg"`
	SubCode string `json:"subCode"`
}

type ResultAsCancelOrder struct {
	ClOrdId string `json:"clOrdId"`
	OrdId   string `json:"ordId"`
	Ts      string `json:"ts"`
	SCode   string `json:"sCode"`
	SMsg    string `json:"sMsg"`
}
