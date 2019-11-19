package handler_test

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/martinohmann/rfoutlet/internal/control"
	"github.com/martinohmann/rfoutlet/internal/handler"
	"github.com/martinohmann/rfoutlet/internal/message"
	"github.com/posener/wstest"
	"github.com/stretchr/testify/assert"
)

var nopDispatcher = new(testDispatcher)

type testDispatcher struct{}

func (testDispatcher) Dispatch(msg message.Envelope) error {
	return nil
}

func TestWebsocket(t *testing.T) {
	r := gin.New()
	r.GET("/ws", handler.Websocket(control.NewHub(), nopDispatcher))
	d := wstest.NewDialer(r)

	_, rr, _ := d.Dial("ws://localhost/ws", nil)

	assert.Equal(t, http.StatusSwitchingProtocols, rr.StatusCode)
}
