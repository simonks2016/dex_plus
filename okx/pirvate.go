package okx

import "DexPlus/websocket"

type Private struct {
	ws      *websocket.WsClient
	channel []string
	auth    string
}
