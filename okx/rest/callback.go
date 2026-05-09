package rest

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/simonks2016/dex_plus/internal/httpClient"
	"github.com/simonks2016/dex_plus/okx/Response"
)

type asyncResult[T any] struct {
	data T
	err  error
}

func OKXCallback[T any](resultCh chan<- asyncResult[T]) httpClient.Callback {
	return func(resp *httpClient.Response, err error) {
		var zero T

		if err != nil {
			resultCh <- asyncResult[T]{zero, err}
			return
		}

		if resp == nil {
			resultCh <- asyncResult[T]{zero, fmt.Errorf("nil response")}
			return
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			resultCh <- asyncResult[T]{
				zero,
				fmt.Errorf("http status=%d, body=%s", resp.StatusCode, string(resp.Body)),
			}
			return
		}

		var out Response.BasicResponse[T]
		if err := json.Unmarshal(resp.Body, &out); err != nil {
			resultCh <- asyncResult[T]{zero, err}
			return
		}

		if out.Code != "0" {
			resultCh <- asyncResult[T]{
				zero,
				fmt.Errorf("okx code=%s, msg=%s", out.Code, out.Msg),
			}
			return
		}
		resultCh <- asyncResult[T]{out.Data, nil}
	}
}
