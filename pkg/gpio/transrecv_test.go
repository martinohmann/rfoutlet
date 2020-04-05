// +build flaky

package gpio_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/stretchr/testify/assert"
	"github.com/warthog618/gpiod"
)

type testWatcherPin struct {
	offset int
	w      *testWatcher
}

func newTestWatcherPin(offset int, watcher *testWatcher) *testWatcherPin {
	return &testWatcherPin{offset, watcher}
}

func (p *testWatcherPin) SetValue(value int) error {
	var eventType gpiod.LineEventType
	switch value {
	case 0:
		eventType = gpiod.LineEventFallingEdge
	case 1:
		eventType = gpiod.LineEventRisingEdge
	default:
		panic(fmt.Sprintf("invalid value: %d", value))
	}

	p.w.events <- gpiod.LineEvent{
		Offset: p.offset,
		Type:   eventType,
	}

	return nil
}

func (p *testWatcherPin) Close() error {
	return nil
}

func TestTransmitReceive(t *testing.T) {
	watcher := newTestWatcher()
	pin := newTestWatcherPin(10, watcher)

	receiver := gpio.NewWatcherReceiver(watcher)
	defer receiver.Close()

	transmitter := gpio.NewPinTransmitter(pin, gpio.TransmissionRetries(15))
	defer transmitter.Close()

	tests := []struct {
		code        uint64
		protocol    int
		pulseLength uint
	}{
		{5510451, 1, 184},
		{83281, 1, 305},
		{86356, 1, 305},
		{5510604, 1, 184},
		{5591317, 1, 330},
	}

	for _, tc := range tests {
		transmitter.Transmit(tc.code, tc.protocol, tc.pulseLength)
		transmitter.Wait()
	}

	var i int
	var lastCode uint64

	for i < len(tests) {
		select {
		case result := <-receiver.Receive():
			if result.Code == lastCode {
				continue
			}

			tc := tests[i]

			assert.Equalf(t, tc.code, result.Code, "received code %d != expected %d", result.Code, tc.code)

			lastCode = result.Code
			i++
		case <-time.After(time.Second):
			assert.FailNow(t, "receive timed out")
			break
		}
	}
}
