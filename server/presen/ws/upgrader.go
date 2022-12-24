package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// https://nananao-dev.hatenablog.com/entry/websocket-upgrader-cors
	// https://foresuke.com/post/go_react_websocket/
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
