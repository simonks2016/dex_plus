package httpClient

import (
	"context"
	"time"
)

type Callback func(resp *Response, err error)

type Method int

const (
	GET Method = iota
	POST
	PATCH
	DELETE
	PUT
)

type Request struct {
	RequestId     string            `json:"requestId"`
	Method        Method            `json:"method"`
	URL           string            `json:"url"`
	Query         map[string]string `json:"query"`
	Header        map[string]string `json:"header"`
	Body          []byte            `json:"body"`
	Ctx           context.Context   `json:"-"`
	Timeout       time.Duration     `json:"timeout"`
	Retry         int               `json:"retry"`
	RetryInterval time.Duration     `json:"retryInterval"`
	EnableBackoff bool              `json:"enable_backoff"`
	Priority      int               `json:"priority"`
	Tags          []string          `json:"tags"`
	Meta          map[string]any    `json:"meta"`
	CreatedAt     time.Time         `json:"created_at"`
	Callback      Callback          `json:"-"`
}
