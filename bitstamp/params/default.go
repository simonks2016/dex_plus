package params

import "github.com/goccy/go-json"

type SubscribeParams struct {
	Event string                `json:"event"`
	Data  SubscribeChannelParam `json:"data"`
}

type SubscribeChannelParam struct {
	Channel string `json:"channel"`
}

func NewSubscribeParams(channel string) SubscribeParams {
	return SubscribeParams{
		Event: "bts:subscribe",
		Data: SubscribeChannelParam{
			Channel: channel,
		},
	}
}

func NewUnsubscribeParams(channel string) SubscribeParams {
	return SubscribeParams{
		Event: "bts:unsubscribe",
		Data: SubscribeChannelParam{
			Channel: channel,
		},
	}
}

func (p *SubscribeParams) Json() []byte {

	marshal, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return marshal
}
