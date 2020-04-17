package gpio_test

import (
	"testing"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/stretchr/testify/assert"
)

func TestTransmitterTransmit(t *testing.T) {
	pin := gpio.NewFakeOutputPin()

	transmitter := gpio.NewPinTransmitter(pin, gpio.TransmissionCount(1))
	defer transmitter.Close()

	<-transmitter.Transmit(0x1, gpio.DefaultProtocols[0], 190)

	assert.Equal(
		t,
		[]int{1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0},
		pin.Values,
	)
}

func TestTransmitterClose(t *testing.T) {
	pin := gpio.NewFakeOutputPin()

	transmitter := gpio.NewPinTransmitter(pin)
	assert.Nil(t, transmitter.Close())

	assert.True(t, pin.Closed)
}
