package gpio

// Most of the transmitter code is ported from the rc-switch c++ implementation to
// go. Check out the rc-switch repository at https://github.com/sui77/rc-switch
// for the original implementation.

import (
	"fmt"
	"os"
	"time"

	"github.com/brian-armstrong/gpio"
)

const (
	DefaultTransmitPin uint = 17
	DefaultReceivePin  uint = 27
	DefaultProtocol    int  = 1
	DefaultPulseLength uint = 189

	numRetries int = 10
	bitLength  int = 24
)

// CodeTransmitter defines the interface for a rf code transmitter.
type CodeTransmitter interface {
	Transmit(uint64, int, uint) error
	Close() error
}

// NativeTransmitter type definition
type NativeTransmitter struct {
	gpioPin  gpio.Pin
	protocol protocol
}

func NewNativeTransmitter(gpioPin uint) (*NativeTransmitter, error) {
	t := &NativeTransmitter{
		gpioPin: gpio.NewOutput(gpioPin, false),
	}

	return t, nil
}

// Transmit transmits a code using given protocol and pulse length
func (t *NativeTransmitter) Transmit(code uint64, protocol int, pulseLength uint) error {
	if err := t.selectProtocol(protocol); err != nil {
		return err
	}

	t.setPulseLength(pulseLength)

	for retry := 0; retry < numRetries; retry++ {
		for j := bitLength - 1; j >= 0; j-- {
			if code&(1<<uint64(j)) > 0 {
				t.transmit(t.protocol.one)
			} else {
				t.transmit(t.protocol.zero)
			}
		}
		t.transmit(t.protocol.sync)
	}

	return nil
}

// Close triggers rpio cleanup
func (t *NativeTransmitter) Close() error {
	t.gpioPin.Close()

	return nil
}

func (t *NativeTransmitter) selectProtocol(protocol int) error {
	if protocol < 1 || protocol > len(protocols) {
		return fmt.Errorf("Protocol %d does not exist", protocol)
	}

	t.protocol = protocols[protocol-1]

	return nil
}

func (t *NativeTransmitter) setPulseLength(pulseLength uint) {
	t.protocol.pulseLength = pulseLength
}

// transmit sends a sequence of high and low pulses on the gpio pin
func (t *NativeTransmitter) transmit(pulses highLow) {
	t.gpioPin.High()
	time.Sleep(time.Microsecond * time.Duration(t.protocol.pulseLength*pulses.high))
	t.gpioPin.Low()
	time.Sleep(time.Microsecond * time.Duration(t.protocol.pulseLength*pulses.low))
}

// NullTransmitter type definition
type NullTransmitter struct{}

// NewNullTransmitter create a transmitter that does nothing except logging the
// transmissions. This is mainly useful for development on systems where
// /dev/gpiomem is not available.
func NewNullTransmitter() (*NullTransmitter, error) {
	t := &NullTransmitter{}

	return t, nil
}

// Transmit transmits the given code via the configured gpio pin
func (t *NullTransmitter) Transmit(code uint64, protocol int, pulseLength uint) error {
	return nil
}

// Close performs cleanup
func (t *NullTransmitter) Close() error {
	return nil
}

// NewTransmitter creates a NativeTransmitter when /dev/gpiomem is available,
// NullTransmitter otherwise.
func NewTransmitter(gpioPin uint) (CodeTransmitter, error) {
	if _, err := os.Stat("/dev/gpiomem"); os.IsNotExist(err) {
		return NewNullTransmitter()
	}

	return NewNativeTransmitter(gpioPin)
}
