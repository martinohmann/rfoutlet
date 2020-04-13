package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/martinohmann/rfoutlet/internal/command"
)

// Handler accepts websocket connections and creates a clients to handle them.
func Handler(hub *Hub, queue chan<- command.Command) gin.HandlerFunc {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Errorf("failed to upgrade websocket connection: %v", err)
			return
		}

		NewClient(hub, conn, queue).Listen()
	}
}
