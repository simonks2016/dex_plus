package okx

const (
	TickersChannel            = "tickers"
	CallAuctionDetailsChannel = "call-auction-details"
	TradesChannel             = "trades"
	Books5Channel             = "books5"
	KLine1SChannel            = "candle1s"
	KLine1DChannel            = "candle1d"
	SandboxPublicURL          = "wss://wspap.okx.com:8443/ws/v5/public"
	SandboxPrivateURL         = "wss://wspap.okx.com:8443/ws/v5/private"
	SandBoxBusinessURL        = "wss://wspap.okx.com:8443/ws/v5/business"
	ProductionPublicURL       = "wss://ws.okx.com:8443/ws/v5/public"
	ProductionPrivateURL      = "wss://ws.okx.com:8443/ws/v5/private"
	ProductionBusinessURL     = "wss://ws.okx.com:8443/ws/v5/business"
)

func PublicURL(isProduction bool) string {

	if isProduction {
		return ProductionPublicURL
	}
	return SandboxPublicURL
}

func PrivateURL(isProduction bool) string {
	if isProduction {
		return ProductionPrivateURL
	}
	return SandboxPrivateURL
}
func BusinessURL(isProduction bool) string {
	if isProduction {
		return ProductionBusinessURL
	}
	return SandboxPublicURL
}
