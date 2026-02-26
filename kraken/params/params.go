package params

import (
	"github.com/goccy/go-json"
)

const (
	Subscribe   = "subscribe"
	Unsubscribe = "unsubscribe"
)

type KrakenParams struct {
	Method string `json:"method"`
	Params Param  `json:"params"`
}

type Param struct {
	Channel string   `json:"channel"`
	Symbol  []string `json:"symbol"`
}

func NewKrakenParams(method, channel string, symbols ...string) *KrakenParams {
	return &KrakenParams{
		Method: method,
		Params: Param{
			Channel: channel,
			Symbol:  symbols,
		},
	}
}

func (k *KrakenParams) Json() []byte {

	marshal, err := json.Marshal(k)
	if err != nil {
		panic(err)
	}
	return marshal
}
