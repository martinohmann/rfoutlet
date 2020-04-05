package gpio_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/stretchr/testify/assert"
)

type testPin struct {
	sequence bytes.Buffer
	closed   bool
}

func newTestPin() *testPin {
	return &testPin{}
}

func (p *testPin) SetValue(value int) error {
	switch value {
	case 1:
		p.sequence.WriteRune('1')
	case 0:
		p.sequence.WriteRune('0')
	default:
		panic(fmt.Sprintf("unexpected value: %d", value))
	}
	return nil
}

func (p *testPin) Close() error {
	p.closed = true
	return nil
}

func TestTransmitterTransmit(t *testing.T) {
	pin := newTestPin()

	transmitter := gpio.NewPinTransmitter(pin, gpio.TransmissionRetries(1))
	defer transmitter.Close()

	err := transmitter.Transmit(0x1, 1, 190)

	assert.Nil(t, err)

	transmitter.Wait()

	assert.Equal(t, "10101010101010101010101010101010101010101010101010", pin.sequence.String())
}

func TestTransmitInvalidProtocol(t *testing.T) {
	pin := newTestPin()

	transmitter := gpio.NewPinTransmitter(pin)
	defer transmitter.Close()

	err := transmitter.Transmit(0x1, 999, 190)

	if assert.Error(t, err) {
		assert.Equal(t, errors.New("Protocol 999 does not exist"), err)
	}
}

func TestTransmitterClose(t *testing.T) {
	pin := newTestPin()

	transmitter := gpio.NewPinTransmitter(pin)
	assert.Nil(t, transmitter.Close())

	assert.True(t, pin.closed)
}
