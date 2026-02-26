package payload

import (
	"bytes"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/simonks2016/dex_plus/kraken/internal"
)

type KrakenPayloadType interface {
	Trade | Ticker | OrderBook | L3OrderEvent
}

func ParseData[T KrakenPayloadType](env *internal.KrakenEnvelope) ([]T, error) {
	// 1. 基础校验：检查数据是否为空
	if len(env.Data) == 0 {
		return nil, nil
	}

	// 2. 预处理：修剪空格并排除 "null" 字符串
	trimmed := bytes.TrimSpace(env.Data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil, nil
	}

	// 3. 核心逻辑：利用泛型直接解析
	// Go 的 json 库可以直接根据传入指针的基类型（这里是 []T）进行反射解析
	var result []T
	if err := json.Unmarshal(trimmed, &result); err != nil {
		// 使用 %w 包装原始错误，方便外部通过 errors.Is 判断
		return nil, fmt.Errorf("failed to parse kraken payload: %w", err)
	}

	return result, nil
}
