package Response

type Position struct {
	InstType               string           `json:"instType"`
	MgnMode                string           `json:"mgnMode"`
	PosId                  string           `json:"posId"`
	PosSide                string           `json:"posSide"`
	Pos                    string           `json:"pos"`
	HedgedPos              string           `json:"hedgedPos"`
	BaseBal                string           `json:"baseBal"`
	QuoteBal               string           `json:"quoteBal"`
	BaseBorrowed           string           `json:"baseBorrowed"`
	BaseInterest           string           `json:"baseInterest"`
	QuoteBorrowed          string           `json:"quoteBorrowed"`
	QuoteInterest          string           `json:"quoteInterest"`
	PosCcy                 string           `json:"posCcy"`
	AvailPos               string           `json:"availPos"`
	AvgPx                  string           `json:"avgPx"`
	NonSettleAvgPx         string           `json:"nonSettleAvgPx"`
	Upl                    string           `json:"upl"`
	UplRatio               string           `json:"uplRatio"`
	UplLastPx              string           `json:"uplLastPx"`
	UplRatioLastPx         string           `json:"uplRatioLastPx"`
	InstId                 string           `json:"instId"`
	Lever                  string           `json:"lever"`
	LiqPx                  string           `json:"liqPx"`
	MarkPx                 string           `json:"markPx"`
	Imr                    string           `json:"imr"`
	Margin                 string           `json:"margin"`
	MgnRatio               string           `json:"mgnRatio"`
	Mmr                    string           `json:"mmr"`
	Liab                   string           `json:"liab"`
	LiabCcy                string           `json:"liabCcy"`
	Interest               string           `json:"interest"`
	TradeId                string           `json:"tradeId"`
	OptVal                 string           `json:"optVal"`
	PendingCloseOrdLiabVal string           `json:"pendingCloseOrdLiabVal"`
	NotionalUsd            string           `json:"notionalUsd"`
	Adl                    string           `json:"adl"`
	Ccy                    string           `json:"ccy"`
	Last                   string           `json:"last"`
	IdxPx                  string           `json:"idxPx"`
	UsdPx                  string           `json:"usdPx"`
	BePx                   string           `json:"bePx"`
	DeltaBS                string           `json:"deltaBS"`
	DeltaPA                string           `json:"deltaPA"`
	GammaBS                string           `json:"gammaBS"`
	GammaPA                string           `json:"gammaPA"`
	ThetaBS                string           `json:"thetaBS"`
	ThetaPA                string           `json:"thetaPA"`
	VegaBS                 string           `json:"vegaBS"`
	VegaPA                 string           `json:"vegaPA"`
	SpotInUseAmt           string           `json:"spotInUseAmt"`
	SpotInUseCcy           string           `json:"spotInUseCcy"`
	ClSpotInUseAmt         string           `json:"clSpotInUseAmt"`
	MaxSpotInUseAmt        string           `json:"maxSpotInUseAmt"`
	RealizedPnl            string           `json:"realizedPnl"`
	SettledPnl             string           `json:"settledPnl"`
	Pnl                    string           `json:"pnl"`
	Fee                    string           `json:"fee"`
	FundingFee             string           `json:"fundingFee"`
	LiqPenalty             string           `json:"liqPenalty"`
	CloseOrderAlgo         []CloseOrderAlgo `json:"closeOrderAlgo"`
	CTime                  string           `json:"cTime"`
	UTime                  string           `json:"uTime"`
	BizRefId               string           `json:"bizRefId"`
	BizRefType             string           `json:"bizRefType"`
}

type CloseOrderAlgo struct {
	AlgoId          string `json:"algoId"`
	SlTriggerPx     string `json:"slTriggerPx"`
	SlTriggerPxType string `json:"slTriggerPxType"`
	TpTriggerPx     string `json:"tpTriggerPx"`
	TpTriggerPxType string `json:"tpTriggerPxType"`
	CloseFraction   string `json:"closeFraction"`
}
