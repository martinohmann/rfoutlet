package websocket

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHub(t *testing.T) {
	h := NewHub()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go h.Run(ctx.Done())

	c1 := newClient(h, nil, nil)
	c2 := newClient(h, nil, nil)

	go func() {
		defer cancel()

		h.register <- c1
		h.register <- c2

		h.Send(c1, []byte{0x1})
		assert.Equal(t, []byte{0x1}, <-c1.send)

		h.Send(c2, []byte{0x2})
		assert.Equal(t, []byte{0x2}, <-c2.send)

		h.Broadcast([]byte{0x3})
		assert.Equal(t, []byte{0x3}, <-c1.send)
		assert.Equal(t, []byte{0x3}, <-c2.send)

		h.unregister <- c1

		h.Broadcast([]byte{0x4})
		assert.Equal(t, []byte{0x4}, <-c2.send)
	}()

	<-ctx.Done()

	if err := ctx.Err(); err == context.DeadlineExceeded {
		t.Fatal(err)
	}
}
