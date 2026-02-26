package internal

import (
	"github.com/goccy/go-json"
)

type BinanceParams struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     string   `json:"id"`
}

func NewBinanceParams(Method string, params ...string) *BinanceParams {

	return &BinanceParams{
		Method: Method,
		Params: params,
		Id:     "1",
	}
}

const (
	SubscribeMethod   = "SUBSCRIBE"
	UnsubscribeMethod = "UNSUBSCRIBE"
)

func (p *BinanceParams) Json() []byte {

	marshal, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return marshal
}

func (p *BinanceParams) Add(symbol, channel string, Is100Ms bool) *BinanceParams {

	p1 := symbol + "@" + channel
	if !Is100Ms {
		p1 = p1 + "@100ms"
	}
	p.Params = append(p.Params, p1)
	return p
}

func (p *BinanceParams) CopyNew(method string) *BinanceParams {
	return &BinanceParams{
		Method: method,
		Params: p.Params,
		Id:     "2",
	}
}
