package gpio_test

import (
	"testing"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/stretchr/testify/assert"
	"github.com/warthog618/gpiod"
)

type testWatcher struct {
	closed bool
	events chan gpiod.LineEvent
}

func newTestWatcher() *testWatcher {
	return &testWatcher{events: make(chan gpiod.LineEvent)}
}

func (w *testWatcher) Watch() <-chan gpiod.LineEvent {
	return w.events
}

func (w *testWatcher) Close() error {
	defer close(w.events)
	w.closed = true
	return nil
}

func TestReceiverClose(t *testing.T) {
	watcher := newTestWatcher()

	receiver := gpio.NewWatcherReceiver(watcher)
	assert.Nil(t, receiver.Close())

	assert.True(t, watcher.closed)
}
