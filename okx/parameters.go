package okx

import "encoding/json"

type Parameters struct {
	Id   *string `json:"id,omitempty"`
	OP   string  `json:"op"`
	Args []*Arg  `json:"args"`
}

func NewSubscribeParameters(args ...*Arg) *Parameters {

	return &Parameters{
		OP:   "subscribe",
		Args: args,
	}
}

func NewUnsubscribeParameters(args ...*Arg) *Parameters {
	return &Parameters{
		OP:   "unsubscribe",
		Args: args,
	}
}

func NewInstIdArg(instId string, channel string) *Arg {
	return &Arg{
		Channel: channel,
		InstId:  &instId,
	}
}

func NewInstFamilyArg(instFamily string, channel string) *Arg {
	return &Arg{
		Channel:    channel,
		InstFamily: &instFamily,
	}
}

func (p *Parameters) Encode() []byte {
	d, _ := json.Marshal(p)
	return d
}
