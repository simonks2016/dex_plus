package rest

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/simonks2016/dex_plus/internal/httpClient"
	"github.com/simonks2016/dex_plus/okx/response"
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

		var out response.BasicResponse[T]
		if err := json.Unmarshal(resp.Body, &out); err != nil {
			resultCh <- asyncResult[T]{zero, err}
			return
		}

		if out.Code != "0" {
			resultCh <- asyncResult[T]{
				out.Data,
				fmt.Errorf("okx code=%s%s", out.Code, func() string {
					if len(out.Msg) >= 0 {
						return ""
					}
					return fmt.Sprintf(" , msg=%s", out.Msg)
				}()),
			}
			return
		}
		resultCh <- asyncResult[T]{out.Data, nil}
	}
}
