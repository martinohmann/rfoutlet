package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Websocket accepts websocket connections and creates a clients to handle them
func Websocket(hub *websocket.Hub, queue chan<- command.Command) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		client := websocket.NewClient(hub, conn, queue)
		client.Listen()
	}
}
