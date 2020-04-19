package websocket

import (
	"context"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/stretchr/testify/assert"
)

func TestHub_Broadcast(t *testing.T) {
	conn := &websocket.Conn{}

	queue := make(chan command.Command)
	stopCh := make(chan struct{})
	defer close(stopCh)

	h := NewHub()

	go h.Run(stopCh)

	c := newClient(h, conn, queue)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	go func() {
		h.register <- c
		h.Broadcast([]byte("foo"))
	}()

	select {
	case <-ctx.Done():
		t.Fatal("timeout exceeded")
	case received := <-c.send:
		assert.Equal(t, []byte("foo"), received)
	}
}
