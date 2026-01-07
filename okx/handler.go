package okx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type MarketEvent interface {
	Trade | Ticker | OrderBook | Kline | CallAuctionDetails | TradeFill | PositionAndBalance | Position
}

type Caller func(payload *Payload) error

func ParseData[T MarketEvent](resp *Payload) ([]T, error) {
	ch := resp.GetChannel()

	if len(resp.Data) == 0 {
		return nil, nil
	}

	trimmed := bytes.TrimSpace(resp.Data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}

	// 先用 T 的“零值类型”判断你想要什么
	var zero T
	switch any(zero).(type) {

	case Kline:
		if !strings.HasPrefix(ch, "candle") {
			return nil, fmt.Errorf("type/channel mismatch: want kline(T=OKXKline) but channel=%s", ch)
		}
		// OKX kline data: [][]string
		var raw [][]string
		if err := json.Unmarshal(resp.Data, &raw); err != nil {
			return nil, fmt.Errorf("unmarshal kline [][]string failed: %w", err)
		}

		if len(raw) > 0 {
			fmt.Println(string(resp.Data))
		}

		kl, err := DecodeOKXLine(raw...) // 你的解码函数：([]OKXKline, error)
		if err != nil {
			return nil, err
		}
		// []OKXKline -> []T
		ret := make([]T, len(kl))
		for i := range kl {
			ret[i] = any(kl[i]).(T)
		}
		return ret, nil

	case Ticker:
		if ch != "tickers" {
			return nil, fmt.Errorf("type/channel mismatch: want ticker(T=OKXTicker) but channel=%s", ch)
		}
		var v []Ticker
		if err := json.Unmarshal(resp.Data, &v); err != nil {
			return nil, fmt.Errorf("unmarshal ticker failed: %w", err)
		}
		ret := make([]T, len(v))
		for i := range v {
			ret[i] = any(v[i]).(T)
		}
		return ret, nil

	case OrderBook:
		// 你可能订阅的是 books / books5 / books50
		if !strings.HasPrefix(ch, "books") {
			return nil, fmt.Errorf("type/channel mismatch: want orderbook(T=OKXOrderBook) but channel=%s", ch)
		}

		var v []OrderBook
		if err := json.Unmarshal(resp.Data, &v); err != nil {
			return nil, fmt.Errorf("unmarshal orderbook failed: %w", err)
		}
		ret := make([]T, len(v))
		for i := range v {
			ret[i] = any(v[i]).(T)
		}
		return ret, nil

	case Trade:
		if ch != "trades" {
			return nil, fmt.Errorf("type/channel mismatch: want trade(T=OKXTrade) but channel=%s", ch)
		}
		var v []Trade
		if err := json.Unmarshal(resp.Data, &v); err != nil {
			return nil, fmt.Errorf("unmarshal trade failed: %w", err)
		}
		ret := make([]T, len(v))
		for i := range v {
			ret[i] = any(v[i]).(T)
		}
		return ret, nil
	case CallAuctionDetails:
		if ch != "call-auction-details" {
			return nil, fmt.Errorf("type/channel mismatch: want call-auction-details but channel=%s", ch)
		}
		var v []CallAuctionDetails
		if err := json.Unmarshal(resp.Data, &v); err != nil {
			return nil, fmt.Errorf("unmarshal trade failed: %w", err)
		}
		ret := make([]T, len(v))
		for i := range v {
			ret[i] = any(v[i]).(T)
		}
		return ret, nil
	case TradeFill:
		if ch != "fills" {
			return nil, fmt.Errorf("type/channel mismatch: want trade fills but channel=%s", ch)
		}
		var v []TradeFill
		if err := json.Unmarshal(resp.Data, &v); err != nil {
			return nil, fmt.Errorf("unmarshal trade failed: %w", err)
		}
		ret := make([]T, len(v))
		for i := range v {
			ret[i] = any(v[i]).(T)
		}
		return ret, nil
	case PositionAndBalance:
		if ch != "balance_and_position" {
			return nil, fmt.Errorf("type/channel mismatch: want balance_and_position but channel=%s", ch)
		}
		var v []PositionAndBalance
		if err := json.Unmarshal(resp.Data, &v); err != nil {
			return nil, fmt.Errorf("unmarshal position_and_balance: %w", err)
		}
		ret := make([]T, len(v))
		for i := range v {
			ret[i] = any(v[i]).(T)
		}
		return ret, nil
	case Position:
		if ch != "positions" {
			return nil, fmt.Errorf("type/channel mismatch: want positions but channel=%s", ch)
		}
		var v []Position
		if err := json.Unmarshal(resp.Data, &v); err != nil {
			return nil, fmt.Errorf("unmarshal position: %w", err)
		}
		ret := make([]T, len(v))
		for i := range v {
			ret[i] = any(v[i]).(T)
		}
		return ret, nil
	default:
		return nil, fmt.Errorf("unsupported generic type")
	}
}
