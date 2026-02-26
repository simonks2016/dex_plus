package internal

import (
	"time"

	"github.com/goccy/go-json"
)

type KrakenEnvelope struct {
	Channel *string         `json:"channel,omitempty"`
	Type    *string         `json:"type,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	Method  *string         `json:"method,omitempty"`
	Result  map[string]any  `json:"result,omitempty"`
	Success *bool           `json:"success,omitempty"`
	Error   *string         `json:"error,omitempty"`
	TimeIn  time.Time       `json:"time_in,omitempty"`
	TimeOut time.Time       `json:"time_out,omitempty"`
}

func (e *KrakenEnvelope) IsAck() bool { return e.Channel == nil && e.Method != nil }
func (e *KrakenEnvelope) IsSuccess() bool {
	return e.Success != nil && *e.Success
}
func (e *KrakenEnvelope) IsSubscription() bool { return e.Channel != nil && e.Type != nil }

func (e *KrakenEnvelope) GetChannel() string {
	if e.Channel == nil {
		return ""
	}
	return *e.Channel
}
