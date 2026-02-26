package param

import (
	"github.com/goccy/go-json"
)

type SubscribeChannelParams struct {
	Channel     string  `json:"channel"`
	InstId      *string `json:"instId,omitempty"`
	InstFamily  *string `json:"instFamily,omitempty"`
	InstType    *string `json:"instType,omitempty"`
	ExtraParams *string `json:"extra_params,omitempty"`
}

func NewSubscribeParameters(args ...SubscribeChannelParams) *Parameters[SubscribeChannelParams] {

	return &Parameters[SubscribeChannelParams]{
		OP:   "subscribe",
		Args: args,
	}
}

func NewUnsubscribeParameters(args ...SubscribeChannelParams) *Parameters[SubscribeChannelParams] {
	return &Parameters[SubscribeChannelParams]{
		OP:   "unsubscribe",
		Args: args,
	}
}

func NewInstIdArg(instId string, channel string) SubscribeChannelParams {
	return SubscribeChannelParams{
		Channel: channel,
		InstId:  &instId,
	}
}

func NewInstFamilyArg(instFamily string, channel string) SubscribeChannelParams {
	return SubscribeChannelParams{
		Channel:    channel,
		InstFamily: &instFamily,
	}
}
func NewInstTypeArg(instType, channel string) SubscribeChannelParams {
	return SubscribeChannelParams{
		Channel:  channel,
		InstType: &instType,
	}
}

func NewExtraParam(data map[string]any) *string {

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	d := string(marshal)
	return &d
}
