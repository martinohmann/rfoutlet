package control

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/martinohmann/rfoutlet/internal/message"
	"github.com/stretchr/testify/assert"
)

var testDispatcher = new(nopDispatcher)

type nopDispatcher int

func (nopDispatcher) Dispatch(e message.Envelope) error {
	return nil
}

func TestHubBroadcast(t *testing.T) {
	conn := &websocket.Conn{}

	h := NewHub()
	c := NewClient(h, testDispatcher, conn)

	h.register <- c

	sent := []byte("foo")

	h.broadcast <- sent

	received := <-c.send

	assert.Equal(t, sent, received)
}
