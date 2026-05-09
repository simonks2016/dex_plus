package httpClient

import (
	"net/http"
	"time"
)

type Response struct {
	RequestID  string        `json:"request_id"`
	StatusCode int           `json:"status_code"`
	Header     http.Header   `json:"header"`
	Body       []byte        `json:"body"`
	Latency    time.Duration `json:"latency"`
	RetryCount int           `json:"retry_count"`
	Err        error         `json:"err"`
}
