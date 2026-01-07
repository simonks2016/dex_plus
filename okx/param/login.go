package param

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

type LoginParameters struct {
	ApiKey     string `json:"apiKey"`
	Passphrase string `json:"passphrase"`
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
}

func NewLoginParameters(apiKey, passphrase, secretKey string) []byte {
	// 1) timestamp: Unix epoch seconds (string)
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	// 2) prehash = timestamp + method + requestPath
	method := "GET"
	requestPath := "/users/self/verify"
	prehash := timestamp + method + requestPath

	// 3) HMAC-SHA256(prehash, secretKey) then Base64
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, _ = mac.Write([]byte(prehash))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// 4) build WS login message
	return NewParameters[LoginParameters]("login", LoginParameters{
		ApiKey:     apiKey,
		Passphrase: passphrase,
		Timestamp:  timestamp,
		Sign:       sign,
	}).Encode()
}
