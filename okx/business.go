package okx

import "DexPlus/websocket"

type Business struct {
	ws      *websocket.WsClient
	channel []string
}
