package param

type PlaceOrderParams struct {
	InstIdCode  *int    `json:"instIdCode,omitempty"`
	InstId      *string `json:"instId,omitempty"`
	TdMode      string  `json:"tdMode"`
	Ccy         *string `json:"ccy"`
	ClOrdId     *string `json:"clOrdId,omitempty"`
	Tag         *string `json:"tag,omitempty"`
	Side        string  `json:"side"`
	PosSide     *string `json:"posSide,omitempty"`
	OrdType     string  `json:"ordType"`
	SZ          string  `json:"sz"`
	Px          *string `json:"px,omitempty"`
	PxUSD       *string `json:"pxUSD,omitempty"`
	PxVol       *string `json:"pxVol,omitempty"`
	ReduceOnly  *bool   `json:"reduceOnly,omitempty"`
	TgtCcy      *string `json:"tgtCcy,omitempty"`
	BanAmend    *bool   `json:"banAmend,omitempty"`
	PxAmendType *string `json:"pxAmendType,omitempty"`
	StpMode     *string `json:"stpMode,omitempty"`
}

type CancelOrder struct {
	InstId     *string `json:"instId,omitempty"`
	InstIdCode *int    `json:"instIdCode,omitempty"`
	OrdId      *string `json:"ordId,omitempty"`
	ClOrdId    *string `json:"clOrdId,omitempty"`
}

type AmendOrder struct {
	InstIdCode  *int    `json:"instIdCode,omitempty"`
	CxlOnFail   *bool   `json:"cxlOnFail,omitempty"`
	OrdId       *string `json:"ordId,omitempty"`
	ClOrdId     *string `json:"clOrdId,omitempty"`
	ReqId       *string `json:"reqId"`
	NewSz       *string `json:"newSz,omitempty"`
	NewPx       *string `json:"newPx,omitempty"`
	NewPxUsd    *string `json:"newPxUsd,omitempty"`
	NewPxVol    *string `json:"newPxVol,omitempty"`
	PxAmendType *string `json:"pxAmendType,omitempty"`
}
