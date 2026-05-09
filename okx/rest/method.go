package rest

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/simonks2016/dex_plus/internal/httpClient"
	"github.com/simonks2016/dex_plus/okx/internal"
)

// GET方法
func doGET[T any](host, path string, client *Client) (T, error) {
	var zero T

	resultCh := make(chan asyncResult[T], 1)

	req := httpClient.Request{
		RequestId: uuid.New().String(),
		Method:    httpClient.GET,
		URL:       host + path,
		Header: client.auth.Headers(
			"GET",
			path,
			"",
			internal.AddHeaders(
				"Content-Type", "application/json",
			),
			internal.AddHeaders(
				"x-simulated-trading", func() string {
					if client.isSandBox {
						return "1"
					}
					return "0"
				}()),
		),
		Timeout:   5 * time.Second,
		Retry:     2,
		CreatedAt: time.Now(),
		Callback:  OKXCallback[T](resultCh),
	}

	if err := client.client.DoAsync(req); err != nil {
		return zero, err
	}

	select {
	case result := <-resultCh:
		return result.data, result.err
	case <-time.After(30 * time.Minute):
		return zero, fmt.Errorf("request timeout: %s", path)
	}
}
