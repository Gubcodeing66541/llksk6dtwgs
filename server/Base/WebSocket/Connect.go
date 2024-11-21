package WebSocket

import (
	"github.com/gorilla/websocket"
)

type Connect struct {
	ConnId string
	UserId string
	Conn   *websocket.Conn
}
