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

	transmissionChanLen = 32
)

type transmission struct {
	code        uint64
	protocol    protocol
	pulseLength uint
}

// CodeTransmitter defines the interface for a rf code transmitter.
type CodeTransmitter interface {
	Transmit(uint64, int, uint) error
	Wait()
	Close() error
}

// NativeTransmitter type definition
type NativeTransmitter struct {
	gpioPin      gpio.Pin
	transmission chan transmission
	transmitted  chan bool
	done         chan bool
}

// NewNativeTransmitter create a native transmitter on the gpio pin
func NewNativeTransmitter(gpioPin uint) (*NativeTransmitter, error) {
	t := &NativeTransmitter{
		gpioPin:      gpio.NewOutput(gpioPin, false),
		transmission: make(chan transmission, transmissionChanLen),
		transmitted:  make(chan bool, transmissionChanLen),
		done:         make(chan bool, 1),
	}

	go t.watch()

	return t, nil
}

// Transmit transmits a code using given protocol and pulse length
func (t *NativeTransmitter) Transmit(code uint64, protocol int, pulseLength uint) error {
	if protocol < 1 || protocol > len(protocols) {
		return fmt.Errorf("Protocol %d does not exist", protocol)
	}

	trans := transmission{
		code:        code,
		protocol:    protocols[protocol-1],
		pulseLength: pulseLength,
	}

	select {
	case t.transmission <- trans:
	default:
	}

	return nil
}

// Transmit transmits a code using given protocol and pulse length
func (t *NativeTransmitter) transmit(trans transmission) {
	for retry := 0; retry < numRetries; retry++ {
		for j := bitLength - 1; j >= 0; j-- {
			if trans.code&(1<<uint64(j)) > 0 {
				t.send(trans.protocol.one, trans.pulseLength)
			} else {
				t.send(trans.protocol.zero, trans.pulseLength)
			}
		}
		t.send(trans.protocol.sync, trans.pulseLength)
	}

	select {
	case t.transmitted <- true:
	default:
	}

	// if we send out codes too quickly in a row it will confuse outlets and
	// they wont react on it. this is especially the case when sending out
	// codes to multiple different outlets in a loop. we sleep a little bit
	// after each transmission to better separate signals flying around.
	time.Sleep(time.Millisecond * 200)
}

// Close triggers rpio cleanup
func (t *NativeTransmitter) Close() error {
	t.done <- true
	t.gpioPin.Close()

	return nil
}

// Transmitted blocks until code is transmitted
func (t *NativeTransmitter) Wait() {
	<-t.transmitted
}

func (t *NativeTransmitter) watch() {
	for {
		select {
		case <-t.done:
			close(t.transmitted)
			return
		case trans := <-t.transmission:
			t.transmit(trans)
		}
	}
}

// transmit sends a sequence of high and low pulses on the gpio pin
func (t *NativeTransmitter) send(pulses highLow, pulseLength uint) {
	t.gpioPin.High()
	time.Sleep(time.Microsecond * time.Duration(pulseLength*pulses.high))
	t.gpioPin.Low()
	time.Sleep(time.Microsecond * time.Duration(pulseLength*pulses.low))
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

// Transmitted blocks until code is transmitted
func (t *NullTransmitter) Wait() {}

// NewTransmitter creates a NativeTransmitter when /dev/gpiomem is available,
// NullTransmitter otherwise.
func NewTransmitter(gpioPin uint) (CodeTransmitter, error) {
	if _, err := os.Stat("/dev/gpiomem"); os.IsNotExist(err) {
		return NewNullTransmitter()
	}

	return NewNativeTransmitter(gpioPin)
}
