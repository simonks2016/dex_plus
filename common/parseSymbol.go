package common

import "strings"

// ParseOKXSymbol 处理 "BTC-USDT"
func ParseOKXSymbol(symbol string) SymbolInfo {
	parts := strings.Split(symbol, "-")
	if len(parts) != 2 {
		return SymbolInfo{Base: symbol} // 异常处理逻辑
	}
	return SymbolInfo{Base: parts[0], Quote: parts[1]}.Standardize()
}

// ParseKrakenSymbol 处理 "BTC/USDT"
func ParseKrakenSymbol(symbol string) SymbolInfo {
	parts := strings.Split(symbol, "/")
	if len(parts) != 2 {
		return SymbolInfo{Base: symbol}
	}
	return SymbolInfo{Base: parts[0], Quote: parts[1]}.Standardize()
}

// ParseCoinBaseSymbol 处理 "BTC-USD"
func ParseCoinBaseSymbol(symbol string) SymbolInfo {
	return ParseOKXSymbol(symbol) // 逻辑与 OKX 一致
}

// ParseNoDelimiterSymbol 处理 "btcusdt" 或 "btcusd" (Binance/Bitstamp)
// 注意：由于没有分隔符，需要匹配常见的计价货币(Quote)
func ParseNoDelimiterSymbol(symbol string) SymbolInfo {
	symbol = strings.ToUpper(symbol)
	// 常见的计价币种（按长度倒序排列，防止先匹配到 USD 而错过 USDT）
	quotes := []string{USD, USDT, USDC, "FDUSD", "TUSD", SOL, ETH, BTC, BNB, "OKB", "HKD", "EUR"}

	for _, q := range quotes {
		if strings.HasSuffix(symbol, q) {
			base := strings.TrimSuffix(symbol, q)
			return SymbolInfo{Base: base, Quote: q}
		}
	}
	return SymbolInfo{Base: symbol} // 未匹配到则返回原样
}

type SymbolInfo struct {
	Base  string
	Quote string
}

// Standardize 将结果统一转为大写，方便后续业务逻辑判断
func (s SymbolInfo) Standardize() SymbolInfo {
	return SymbolInfo{
		Base:  strings.ToUpper(s.Base),
		Quote: strings.ToUpper(s.Quote),
	}
}
func (s SymbolInfo) StandardizeString(delimiter string) string {
	return s.Base + "-" + s.Quote
}

func ParseSymbol(exchange string, symbol string) SymbolInfo {
	exchange = strings.ToLower(exchange)

	switch exchange {
	case "okx", "coinbase":
		return ParseOKXSymbol(symbol)
	case "kraken":
		return ParseKrakenSymbol(symbol)
	case "binance", "bitstamp":
		return ParseNoDelimiterSymbol(symbol)
	default:
		// 如果无法识别交易所，尝试常见的几种分隔符
		if strings.Contains(symbol, "-") {
			return ParseOKXSymbol(symbol)
		} else if strings.Contains(symbol, "/") {
			return ParseKrakenSymbol(symbol)
		}
		return SymbolInfo{Base: symbol}
	}
}
