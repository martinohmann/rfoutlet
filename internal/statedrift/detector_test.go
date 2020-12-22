package statedrift

import (
	"testing"

	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/stretchr/testify/assert"
)

type fakeReceiver struct {
	results chan gpio.ReceiveResult
}

func (f *fakeReceiver) Receive() <-chan gpio.ReceiveResult {
	return f.results
}

func (f *fakeReceiver) Close() error { return nil }

func TestDetector(t *testing.T) {
	o1 := &outlet.Outlet{ID: "foo", CodeOn: 123, CodeOff: 456, State: outlet.StateOff}
	o2 := &outlet.Outlet{ID: "bar", CodeOn: 789, CodeOff: 234, State: outlet.StateOn}

	reg := outlet.NewRegistry()
	reg.RegisterOutlets(o1, o2)

	recv := &fakeReceiver{make(chan gpio.ReceiveResult)}

	queue := make(chan command.Command)
	stopCh := make(chan struct{})
	defer close(stopCh)

	d := NewDetector(reg, recv, queue)
	go func() {
		d.Run(stopCh)
		close(queue)
	}()

	go func() {
		defer close(recv.results)
		recv.results <- gpio.ReceiveResult{Code: 456}
		recv.results <- gpio.ReceiveResult{Code: 123}
		recv.results <- gpio.ReceiveResult{Code: 42}
		recv.results <- gpio.ReceiveResult{Code: 789}
		recv.results <- gpio.ReceiveResult{Code: 234}
	}()

	expected := []command.Command{
		command.StateCorrectionCommand{Outlet: o1, DesiredState: outlet.StateOn},
		command.StateCorrectionCommand{Outlet: o2, DesiredState: outlet.StateOff},
	}

	received := make([]command.Command, 0)

	for cmd := range queue {
		received = append(received, cmd)
	}

	assert.Equal(t, expected, received)
}
