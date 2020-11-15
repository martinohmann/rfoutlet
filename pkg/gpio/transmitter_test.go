package gpio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransmitterTransmit(t *testing.T) {
	pin := NewFakeOutputPin()

	tx := NewPinTransmitter(pin, TransmissionCount(1))
	defer tx.Close()

	<-tx.Transmit(0x1, DefaultProtocols[0], 190)

	assert.Equal(
		t,
		[]int{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0},
		pin.Values,
	)
}

func TestTransmitterClose(t *testing.T) {
	pin := NewFakeOutputPin()

	tx := NewPinTransmitter(pin)
	assert.Nil(t, tx.Close())

	assert.True(t, pin.Closed)
}
