package okx

import (
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/simonks2016/dex_plus/internal/client"
)

func WithLogger(log *log.Logger) client.Option {

	return func(cfg *client.Config) {
		cfg.Logger = log
	}
}

func WithSendTimeout(timeout time.Duration) client.Option {
	return func(cfg *client.Config) {
		cfg.SendTimeout = timeout
	}
}

func WithSetAuth() client.Option {
	return func(cfg *client.Config) {
		cfg.IsNeedAuth = true
	}
}

func WithForbidIpV6() client.Option {
	return func(cfg *client.Config) {
		cfg.IsForbidIPV6 = true
	}
}
func WithNetDialer(dialer *websocket.Dialer) client.Option {
	return func(cfg *client.Config) {
		cfg.Dialer = dialer
	}
}

func WithSandboxEnv() client.Option {
	return func(cfg *client.Config) {

		u, err := url.Parse(cfg.URL)
		if err != nil {
			return
		}

		host := u.Hostname() // 不带端口
		port := u.Port()

		// 只把 host 的第一个 label 从 "ws" 换成 "wspap"
		labels := strings.Split(host, ".")
		if len(labels) > 0 && labels[0] == "ws" {
			labels[0] = "wspap"
			u.Host = strings.Join(labels, ".")
			if port != "" {
				u.Host += ":" + port
			}
		}
		cfg.URL = u.String()
		return
	}
}
