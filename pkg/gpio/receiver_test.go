package gpio_test

import (
	"testing"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/stretchr/testify/assert"
)

func TestReceiverClose(t *testing.T) {
	watcher := gpio.NewFakeWatcher()

	receiver := gpio.NewWatcherReceiver(watcher, gpio.ReceiverProtocols(gpio.DefaultProtocols))
	assert.Nil(t, receiver.Close())

	assert.True(t, watcher.Closed)
}
