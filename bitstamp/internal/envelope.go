package internal

import (
	"strings"

	"github.com/goccy/go-json"
)

type Envelope struct {
	Event   string          `json:"event"`
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data"`
}

func (e *Envelope) GetSymbol() string {
	if e.Channel == "" {
		return ""
	}
	idx := strings.LastIndex(e.Channel, "_")
	if idx == -1 || idx == len(e.Channel)-1 {
		return ""
	}
	return e.Channel[idx+1:]
}
