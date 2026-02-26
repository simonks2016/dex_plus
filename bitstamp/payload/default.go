package payload

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/simonks2016/dex_plus/bitstamp/internal"
)

type Trade struct {
	Id             int     `json:"id"`
	IdStr          string  `json:"id_str"`
	Amount         float64 `json:"amount"`
	AmountStr      string  `json:"amount_str"`
	Price          float64 `json:"price"`
	PriceStr       string  `json:"price_str"`
	Type           int     `json:"type"`
	Timestamp      string  `json:"timestamp"`
	MicroTimestamp string  `json:"microtimestamp"`
	BuyOrderId     int     `json:"buy_order_id"`
	SellOrderId    int     `json:"sell_order_id"`
}

type OrderBook struct {
	Timestamp      string     `json:"timestamp"`
	MicroTimestamp string     `json:"microtimestamp"`
	Bids           [][]string `json:"bids"`
	Asks           [][]string `json:"asks"`
}

type BitstampTypeInterface interface {
	Trade | OrderBook
}

func ParseData[T BitstampTypeInterface](env *internal.Envelope) (T, error) {

	var result T

	// 1. 基础校验：检查数据是否为空
	if len(env.Data) == 0 {
		return result, nil
	}

	// 2. 预处理：修剪空格并排除 "null" 字符串
	trimmed := bytes.TrimSpace(env.Data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return result, nil
	}
	// 3. 核心逻辑：利用泛型直接解析
	// Go 的 json 库可以直接根据传入指针的基类型（这里是 []T）进行反射解析
	if err := json.Unmarshal(trimmed, &result); err != nil {
		// 使用 %w 包装原始错误，方便外部通过 errors.Is 判断
		return result, fmt.Errorf("failed to parse kraken payload: %w", err)
	}
	return result, nil
}
