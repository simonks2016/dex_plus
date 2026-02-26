package param

import (
	"crypto/md5"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type Parameters[T UnionArg] struct {
	Id            *string `json:"id,omitempty"`
	OP            string  `json:"op"`
	Args          []T     `json:"args"`
	ExpTime       *string `json:"exp_time,omitempty"`
	TradeQuoteCcy *string `json:"tradeQuoteCcy,omitempty"`
}

func NewParameters[T UnionArg](op string, args ...T) *Parameters[T] {

	m := md5.New()
	m.Write([]byte(uuid.New().String()))
	u := fmt.Sprintf("%x", m.Sum(nil))

	return &Parameters[T]{
		Id:   &u,
		OP:   op,
		Args: args,
	}
}

type UnionArg interface {
	SubscribeChannelParams | PlaceOrderParams | CancelOrder | AmendOrder | LoginParameters
}

func (p *Parameters[T]) Encode() []byte {
	d, _ := json.Marshal(p)
	return d
}
