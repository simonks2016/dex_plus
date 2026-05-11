package response

type AccountBalance struct {
	UTime                 string                 `json:"uTime"`
	TotalEq               string                 `json:"totalEq"`
	IsoEq                 string                 `json:"isoEq"`
	AdjEq                 string                 `json:"adjEq"`
	AvailEq               string                 `json:"availEq"`
	OrdFroz               string                 `json:"ordFroz"`
	Imr                   string                 `json:"imr"`
	Mmr                   string                 `json:"mmr"`
	BorrowFroz            string                 `json:"borrowFroz"`
	MgnRatio              string                 `json:"mgnRatio"`
	NotionalUsd           string                 `json:"notionalUsd"`
	NotionalUsdForBorrow  string                 `json:"notionalUsdForBorrow"`
	NotionalUsdForSwap    string                 `json:"notionalUsdForSwap"`
	NotionalUsdForFutures string                 `json:"notionalUsdForFutures"`
	NotionalUsdForOption  string                 `json:"notionalUsdForOption"`
	Upl                   string                 `json:"upl"`
	Delta                 string                 `json:"delta"`
	DeltaLever            string                 `json:"deltaLever"`
	DeltaNeutralStatus    string                 `json:"deltaNeutralStatus"`
	Details               []AccountBalanceDetail `json:"details"`
}

type AccountBalanceDetail struct {
	Ccy                   string `json:"ccy"`
	Eq                    string `json:"eq"`
	CashBal               string `json:"cashBal"`
	UTime                 string `json:"uTime"`
	IsoEq                 string `json:"isoEq"`
	AvailEq               string `json:"availEq"`
	DisEq                 string `json:"disEq"`
	FixedBal              string `json:"fixedBal"`
	AvailBal              string `json:"availBal"`
	FrozenBal             string `json:"frozenBal"`
	OrdFrozen             string `json:"ordFrozen"`
	Liab                  string `json:"liab"`
	Upl                   string `json:"upl"`
	UplLiab               string `json:"uplLiab"`
	CrossLiab             string `json:"crossLiab"`
	IsoLiab               string `json:"isoLiab"`
	RewardBal             string `json:"rewardBal"`
	MgnRatio              string `json:"mgnRatio"`
	Imr                   string `json:"imr"`
	Mmr                   string `json:"mmr"`
	Interest              string `json:"interest"`
	Twap                  string `json:"twap"`
	FrpType               string `json:"frpType"`
	MaxLoan               string `json:"maxLoan"`
	EqUsd                 string `json:"eqUsd"`
	BorrowFroz            string `json:"borrowFroz"`
	NotionalLever         string `json:"notionalLever"`
	StgyEq                string `json:"stgyEq"`
	IsoUpl                string `json:"isoUpl"`
	SpotInUseAmt          string `json:"spotInUseAmt"`
	ClSpotInUseAmt        string `json:"clSpotInUseAmt"`
	MaxSpotInUse          string `json:"maxSpotInUse"`
	SpotIsoBal            string `json:"spotIsoBal"`
	SmtSyncEq             string `json:"smtSyncEq"`
	SpotCopyTradingEq     string `json:"spotCopyTradingEq"`
	SpotBal               string `json:"spotBal"`
	OpenAvgPx             string `json:"openAvgPx"`
	AccAvgPx              string `json:"accAvgPx"`
	SpotUpl               string `json:"spotUpl"`
	SpotUplRatio          string `json:"spotUplRatio"`
	TotalPnl              string `json:"totalPnl"`
	TotalPnlRatio         string `json:"totalPnlRatio"`
	ColRes                string `json:"colRes"`
	ColBorrAutoConversion string `json:"colBorrAutoConversion"`
	CollateralRestrict    bool   `json:"collateralRestrict"`
	CollateralEnabled     bool   `json:"collateralEnabled"`
	AutoLendStatus        string `json:"autoLendStatus"`
	AutoLendMtAmt         string `json:"autoLendMtAmt"`
}
