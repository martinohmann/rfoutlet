package controller

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/stretchr/testify/assert"
)

type fakeSwitcher struct{}

func (fakeSwitcher) Switch(*outlet.Outlet, outlet.State) error { return nil }

type testBroadcaster chan []byte

func (b testBroadcaster) Broadcast(msg []byte) {
	b <- msg
}

type testCommand struct {
	context   command.Context
	doneCh    chan struct{}
	broadcast bool
	err       error
}

func (c *testCommand) Execute(context command.Context) (bool, error) {
	c.context = context
	close(c.doneCh)
	return c.broadcast, c.err
}

func TestController(t *testing.T) {
	tests := []struct {
		name       string
		broadcast  bool
		commandErr error
	}{
		{
			name: "no broadcast",
		},
		{
			name:      "broadcast",
			broadcast: true,
		},
		{
			name:       "broadcast, but error -> do not broadcast",
			broadcast:  true,
			commandErr: errors.New("whoops"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			queue := make(chan command.Command)
			b := make(testBroadcaster)
			s := &fakeSwitcher{}

			stopCh := make(chan struct{})
			defer close(stopCh)

			c := &Controller{
				Registry:     outlet.NewRegistry(),
				Switcher:     s,
				Broadcaster:  b,
				CommandQueue: queue,
			}

			go c.Run(stopCh)

			ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			defer cancel()

			doneCh := make(chan struct{})
			cmd := &testCommand{
				doneCh:    doneCh,
				broadcast: test.broadcast,
				err:       test.commandErr,
			}

			go func() { queue <- cmd }()

			select {
			case <-ctx.Done():
				t.Fatal("timeout exceeded")
			case <-doneCh:
				expectedCtx := command.Context{
					Registry: c.Registry,
					Switcher: c.Switcher,
				}
				assert.Equal(t, expectedCtx, cmd.context)

				if test.broadcast && test.commandErr == nil {
					select {
					case <-ctx.Done():
						t.Fatal("timeout exceeded waiting for broadcast")
					case <-b:
					}
				}
			}
		})
	}
}
