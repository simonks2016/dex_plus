package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"time"
)

type Auth struct {
	ApiKey     string
	Passphrase string
	SecretKey  string
}

func NewAuth(apiKey, passphrase, secretKey string) *Auth {
	return &Auth{
		ApiKey:     apiKey,
		Passphrase: passphrase,
		SecretKey:  secretKey,
	}
}

// Timestamp 返回 OKX 要求的 UTC ISO8601 毫秒时间
func (auth *Auth) Timestamp() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
}

// Signature 生成签名
func (auth *Auth) Signature(timestamp, method, requestPath, body string) string {
	method = strings.ToUpper(method)

	preHash := timestamp + method + requestPath + body

	mac := hmac.New(sha256.New, []byte(auth.SecretKey))
	mac.Write([]byte(preHash))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// Headers 生成请求头
func (auth *Auth) Headers(method, requestPath, body string, extraHeaders ...ExtraHeader) map[string]string {
	timestamp := auth.Timestamp()
	sign := auth.Signature(timestamp, method, requestPath, body)

	d := map[string]string{
		"OK-ACCESS-KEY":        auth.ApiKey,
		"OK-ACCESS-SIGN":       sign,
		"OK-ACCESS-TIMESTAMP":  timestamp,
		"OK-ACCESS-PASSPHRASE": auth.Passphrase,
		"Content-Type":         "application/json",
	}

	for _, setHeader := range extraHeaders {
		setHeader(d)
	}
	return d
}

type ExtraHeader func(map[string]string)

func AddHeaders(headerName string, value string) ExtraHeader {
	return func(m map[string]string) {
		m[headerName] = value
	}
}
