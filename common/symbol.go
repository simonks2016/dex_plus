package common

import (
	"strings"
)

const (
	BTC  = "BTC"
	USDC = "USDC"
	USDT = "USDT"
	ETH  = "ETH"
	SOL  = "SOL"
	BNB  = "BNB"
	XRP  = "XRP"
	ADA  = "ADA"
	DOGE = "DOGE"
	AVAX = "AVAX"
	LINK = "LINK"
	TON  = "TON"
	USD  = "USD"
)

func buildSymbol(base string, quote string, relimiter string) string {
	return base + relimiter + quote
}

func OKXSymbol(currency ...string) string {
	// 逻辑：将 base 和 quote 拼接成 OKX 识别的格式，如 BTC-USDT

	if len(currency) <= 0 {
		panic("the base currency not set")
	}
	base := currency[0]
	quote := USD

	if len(currency) >= 2 {
		quote = currency[1]
	}

	return buildSymbol(base, quote, "-")
}

func BinanceSymbol(currency ...string) string {
	// 逻辑：将 base 和 quote 拼接成 Binance 识别的格式，如 btcusdt
	if len(currency) <= 0 {
		panic("the base currency not set")
	}
	base := currency[0]
	quote := USD

	if len(currency) >= 2 {
		quote = currency[1]
	}
	return buildSymbol(strings.ToLower(base), strings.ToLower(quote), "")
}

func KrakenSymbol(currency ...string) string {
	// 逻辑：将 base 和 quote 拼接成 Kraken 识别的格式，如 BTC/USDT
	if len(currency) <= 0 {
		panic("the base currency not set")
	}
	base := currency[0]
	quote := USD

	if len(currency) >= 2 {
		quote = currency[1]
	}
	base = strings.ToUpper(base)
	quote = strings.ToUpper(quote)

	return buildSymbol(base, quote, "/")
}

func CoinBaseSymbol(currency ...string) string {
	// 逻辑：将 base 和 quote 拼接成 Kraken 识别的格式，如 BTC/USDT
	if len(currency) <= 0 {
		panic("the base currency not set")
	}
	base := currency[0]
	quote := USD

	if len(currency) >= 2 {
		quote = currency[1]
	}
	base = strings.ToUpper(base)
	quote = strings.ToUpper(quote)

	return buildSymbol(base, quote, "-")
}

func BitstampSymbol(currency ...string) string {
	// 逻辑：将 base 和 quote 拼接成 Bitstamp 识别的格式，如 btcusd
	if len(currency) <= 0 {
		panic("the base currency not set")
	}
	base := currency[0]
	quote := USD

	if len(currency) >= 2 {
		quote = currency[1]
	}
	base = strings.ToLower(base)
	quote = strings.ToLower(quote)

	return buildSymbol(base, quote, "")
}
