package websocket

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/controller/commands"
	"github.com/posener/wstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_listenRead(t *testing.T) {
	tests := []struct {
		name            string
		data            interface{}
		expectedCmdType command.Command
	}{
		{
			name: "outlet action",
			data: map[string]interface{}{
				"type": "outlet",
				"data": map[string]string{
					"action": "on",
					"id":     "foo",
				},
			},
			expectedCmdType: &commands.OutletCommand{},
		},
		{
			name: "status action",
			data: map[string]interface{}{
				"type": "status",
			},
			expectedCmdType: &commands.StatusCommand{},
		},
		{
			name: "unknown command",
			data: map[string]interface{}{
				"type": "foo",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hub := NewHub()

			queue := make(chan command.Command)

			r := gin.New()
			r.GET("/ws", Handler(hub, queue))

			c, _, err := wstest.NewDialer(r).Dial("ws://localhost/ws", nil)
			defer c.Close()

			require.NoError(t, err)

			go func() {
				err := c.WriteJSON(test.data)
				require.NoError(t, err)
			}()

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			go hub.Run(ctx.Done())

			select {
			case <-ctx.Done():
				if test.expectedCmdType != nil {
					t.Fatal(ctx.Err())
				}
			case cmd := <-queue:
				if test.expectedCmdType == nil {
					t.Fatalf("did not expect command, but got %T", cmd)
				}

				assert.IsType(t, test.expectedCmdType, cmd)
			}
		})
	}
}

func TestClient_listenWrite(t *testing.T) {
	hub := NewHub()

	queue := make(chan command.Command)
	done := make(chan struct{})

	r := gin.New()
	r.GET("/ws", Handler(hub, queue))

	c, rr, err := wstest.NewDialer(r).Dial("ws://localhost/ws", nil)
	defer c.Close()

	require.NoError(t, err)
	assert.Equal(t, http.StatusSwitchingProtocols, rr.StatusCode)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go hub.Run(ctx.Done())

	type foo struct {
		Name string
	}

	go func() {
		defer close(done)
		val := foo{}

		err := c.ReadJSON(&val)
		require.NoError(t, err)

		assert.Equal(t, foo{Name: "bar"}, val)
	}()

	hub.Broadcast([]byte(`{"name":"bar"}`))

	select {
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	case <-done:
	}
}
