package gpio_test

import (
	"testing"
	"time"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/stretchr/testify/assert"
)

type testWatcherPin struct {
	p uint
	w *testWatcher
}

func newTestWatcherPin(pin uint, watcher *testWatcher) *testWatcherPin {
	return &testWatcherPin{pin, watcher}
}

func (p *testWatcherPin) High() error {
	p.w.notification <- testNotification{p.p, 1}
	return nil
}

func (p *testWatcherPin) Low() error {
	p.w.notification <- testNotification{p.p, 0}
	return nil
}

func (p *testWatcherPin) Close() {}

func TestTransmitReceive(t *testing.T) {
	var gpioPin uint = 17
	gpio.TransmitRetries = 15

	watcher := newTestWatcher()
	pin := newTestWatcherPin(gpioPin, watcher)

	receiver := gpio.NewNativeReceiver(gpioPin, watcher)
	defer receiver.Close()

	transmitter := gpio.NewNativeTransmitter(pin)
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
