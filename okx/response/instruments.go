package response

type Instruments struct {
	Alias             string   `json:"alias"`
	AuctionEndTime    string   `json:"auctionEndTime"`
	BaseCcy           string   `json:"baseCcy"`
	Category          string   `json:"category"`
	CtMult            string   `json:"ctMult"`
	CtType            string   `json:"ctType"`
	CtVal             string   `json:"ctVal"`
	CtValCcy          string   `json:"ctValCcy"`
	ContTdSwTime      string   `json:"contTdSwTime"`
	ExpTime           string   `json:"expTime"`
	FutureSettlement  bool     `json:"futureSettlement"`
	GroupId           string   `json:"groupId"`
	InstFamily        string   `json:"instFamily"`
	InstId            string   `json:"instId"`
	InstType          string   `json:"instType"`
	Lever             string   `json:"lever"`
	ListTime          string   `json:"listTime"`
	LotSz             string   `json:"lotSz"`
	MaxIcebergSz      string   `json:"maxIcebergSz"`
	MaxLmtAmt         string   `json:"maxLmtAmt"`
	MaxLmtSz          string   `json:"maxLmtSz"`
	MaxMktAmt         string   `json:"maxMktAmt"`
	MaxMktSz          string   `json:"maxMktSz"`
	MaxStopSz         string   `json:"maxStopSz"`
	MaxTriggerSz      string   `json:"maxTriggerSz"`
	MaxTwapSz         string   `json:"maxTwapSz"`
	MinSz             string   `json:"minSz"`
	OptType           string   `json:"optType"`
	OpenType          string   `json:"openType"`
	PreMktSwTime      string   `json:"preMktSwTime"`
	QuoteCcy          string   `json:"quoteCcy"`
	TradeQuoteCcyList []string `json:"tradeQuoteCcyList"`
	SettleCcy         string   `json:"settleCcy"`
	State             string   `json:"state"`
	RuleType          string   `json:"ruleType"`
	Stk               string   `json:"stk"`
	TickSz            string   `json:"tickSz"`
	Uly               string   `json:"uly"`
	InstIdCode        int      `json:"instIdCode"`
	InstCategory      string   `json:"instCategory"`
	UpcChg            []UpcChg `json:"upcChg"`
}

type UpcChg struct {
	Param    string `json:"param"`
	NewValue string `json:"newValue"`
	EffTime  string `json:"effTime"`
}
