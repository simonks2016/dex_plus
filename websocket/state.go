package websocket

type ConnectionObserver interface {
	OnConnecting(reason string)
	OnMessage(data []byte) error // 表示已接收到信息
	OnConnected()                // 表示已连接上
	OnDisconnecting()            //表示断开连接
	OnDisconnected()
	OnError(err error)
	OnAuth(sender func([]byte) error, done func(error))
}
