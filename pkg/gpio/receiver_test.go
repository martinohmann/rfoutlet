package gpio_test

import (
	"testing"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/stretchr/testify/assert"
)

type testNotification struct {
	pin   uint
	value uint
}

type testWatcher struct {
	pin          uint
	closed       bool
	notification chan testNotification
}

func newTestWatcher() *testWatcher {
	return &testWatcher{notification: make(chan testNotification, 32)}
}

func (w *testWatcher) Watch() (uint, uint) {
	notification := <-w.notification
	return notification.pin, notification.value
}

func (w *testWatcher) AddPin(pin uint) {
	w.pin = pin
}

func (w *testWatcher) Close() {
	w.closed = true
}

func TestReceiverClose(t *testing.T) {
	watcher := newTestWatcher()

	receiver := gpio.NewNativeReceiver(1, watcher)
	assert.Nil(t, receiver.Close())

	assert.True(t, watcher.closed)
}

func TestReceiverWatcherAddPin(t *testing.T) {
	watcher := newTestWatcher()

	receiver := gpio.NewNativeReceiver(17, watcher)
	defer receiver.Close()

	assert.Equal(t, uint(17), watcher.pin)
}
