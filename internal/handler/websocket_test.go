package handler_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/handler"
	"github.com/martinohmann/rfoutlet/internal/websocket"
	"github.com/posener/wstest"
	"github.com/stretchr/testify/assert"
)

func TestWebsocket(t *testing.T) {
	queue := make(chan command.Command)

	r := gin.New()
	r.GET("/ws", handler.Websocket(websocket.NewHub(), queue))
	d := wstest.NewDialer(r)

	_, rr, _ := d.Dial("ws://localhost/ws", nil)

	assert.Equal(t, http.StatusSwitchingProtocols, rr.StatusCode)
}
