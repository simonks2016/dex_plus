package rest

import (
	"net/url"
	"strings"

	"github.com/simonks2016/dex_plus/okx/internal"
)

type Option func(*Client)

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.BaseUrl = baseURL
		return
	}
}

func WithAuth(apiKey, secretKey, passphrase string) Option {
	return func(c *Client) {
		c.auth = internal.NewAuth(apiKey, passphrase, secretKey)
	}
}

func WithSandbox() Option {
	return func(c *Client) {
		c.isSandBox = true
	}
}

type QueryParam func(map[string]string)

func WithQueryParam(name string, value string) QueryParam {
	return func(m map[string]string) {
		m[name] = value
	}
}

func buildPath(path string, queryParams ...QueryParam) string {
	query := make(map[string]string)

	for _, setParam := range queryParams {
		if setParam != nil {
			setParam(query)
		}
	}

	values := url.Values{}
	for k, v := range query {
		if v == "" {
			continue
		}
		values.Set(k, v)
	}

	if encoded := values.Encode(); encoded != "" {
		return path + "?" + encoded
	}

	return path
}

func joinValues(values ...string) string {
	arr := make([]string, 0, len(values))

	for _, v := range values {
		if v == "" {
			continue
		}
		arr = append(arr, v)
	}
	return strings.Join(arr, ",")
}
func WithInstType(instType string) QueryParam {
	return WithQueryParam("instType", instType)
}

func WithInstId(instId string) QueryParam {
	return WithQueryParam("instId", instId)
}

func WithInstIds(instIds ...string) QueryParam {
	return WithQueryParam("instId", joinValues(instIds...))
}

func WithPosId(posId string) QueryParam {
	return WithQueryParam("posId", posId)
}

func WithPosIds(posIds ...string) QueryParam {
	return WithQueryParam("posId", joinValues(posIds...))
}

func WithCcy(ccy string) QueryParam {
	return WithQueryParam("ccy", ccy)
}

func WithCcies(ccies ...string) QueryParam {
	return WithQueryParam("ccy", joinValues(ccies...))
}

func WithOrdId(ordId string) QueryParam {
	return WithQueryParam("ordId", ordId)
}

func WithClOrdId(clOrdId string) QueryParam {
	return WithQueryParam("clOrdId", clOrdId)
}

func WithInstFamily(instFamily string) QueryParam {
	return WithQueryParam("instFamily", instFamily)
}
