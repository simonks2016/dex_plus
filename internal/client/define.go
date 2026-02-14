package client

import "context"

type Client interface {
	SetObserver(ob ConnectionObserver) Client
	Send(context.Context, []byte) error
	Close()
	Reconnect(reason string)
	Start()
}
