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
			require.NoError(t, err)
			defer c.Close()

			go func(data interface{}) {
				require.NoError(t, c.WriteJSON(data))
			}(test.data)

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
	stopCh := make(chan struct{})
	defer close(stopCh)

	hub := NewHub()
	go hub.Run(stopCh)

	queue := make(chan command.Command)

	r := gin.New()
	r.GET("/ws", Handler(hub, queue))

	c, rr, err := wstest.NewDialer(r).Dial("ws://localhost/ws", nil)
	require.NoError(t, err)
	defer c.Close()

	assert.Equal(t, http.StatusSwitchingProtocols, rr.StatusCode)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	go func() {
		// give the client some time to register to hub before broadcasting
		<-time.After(20 * time.Millisecond)
		hub.Broadcast([]byte(`{"name":"bar"}`))

		<-ctx.Done()
		c.Close()
	}()

	type foo struct {
		Name string
	}

	val := foo{}

	require.NoError(t, c.ReadJSON(&val))
	assert.Equal(t, foo{Name: "bar"}, val)
}
