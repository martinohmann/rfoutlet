package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/websocket"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("component", "handler")

// Websocket accepts websocket connections and creates a clients to handle them
func Websocket(hub *websocket.Hub, queue chan<- command.Command) gin.HandlerFunc {
	upgrader := ws.Upgrader{
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

		client := websocket.NewClient(hub, conn, queue)
		client.Listen()
	}
}
