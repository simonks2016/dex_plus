package params

import "github.com/goccy/go-json"

const (
	Subscribe   = "subscribe"
	Unsubscribe = "unsubscribe"
)

type SubscribeParams struct {
	Type       string   `json:"type"`
	Channels   []string `json:"channels"`
	ProductIDs []string `json:"product_ids"`
}

func NewSubscribeParams(method string, productIDs ...string) SubscribeParams {
	return SubscribeParams{
		Type:       method,
		Channels:   []string{},
		ProductIDs: productIDs,
	}
}

func (s *SubscribeParams) AddChannel(channels ...string) *SubscribeParams {
	s.Channels = append(s.Channels, channels...)
	return s
}

func (s *SubscribeParams) Json() []byte {

	marshal, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return marshal
}

func (s *SubscribeParams) IsEmpty() bool {

	return !(len(s.Channels) > 0 && len(s.ProductIDs) > 0)

}
