package client

type ConnectionObserver interface {
	OnConnecting(reason string)
	OnConnected()     // 表示已连接上
	OnDisconnecting() //表示断开连接
	OnDisconnected()
	OnMessage(data []byte) error // 表示已接收到信息
	OnError(err error)
}
