package gpio

import (
	"fmt"
	"os"
	"time"

	"github.com/brian-armstrong/gpio"
)

const (
	// DefaultTransmitPin defines the default gpio pin for transmitting rf codes
	DefaultTransmitPin uint = 17

	// DefaultReceivePin defines the default gpio pin for receiving rf codes
	DefaultReceivePin uint = 27

	// DefaultProtocol defines the default rf protocol
	DefaultProtocol int = 1

	// DefaultPulseLength defines the default pulse length
	DefaultPulseLength uint = 189

	transmissionChanLen = 32
	bitLength           = 24
)

// TransmitRetries defines how many times a code should be transmitted in a
// row. The higher the value, the more likely it is that an outlet actually
// received the code.
var TransmitRetries int = 10

type transmission struct {
	code        uint64
	protocol    Protocol
	pulseLength uint
}

// CodeTransmitter defines the interface for a rf code transmitter.
type CodeTransmitter interface {
	Transmit(uint64, int, uint) error
	Wait()
	Close() error
}

// Pin defines an interface for a gpio pin
type Pin interface {
	High() error
	Low() error
	Close()
}

// NativeTransmitter type definition.
type NativeTransmitter struct {
	gpioPin      Pin
	transmission chan transmission
	transmitted  chan bool
	done         chan bool
}

// NewNativeTransmitter create a native transmitter on the gpio pin.
func NewNativeTransmitter(gpioPin Pin) *NativeTransmitter {
	t := &NativeTransmitter{
		gpioPin:      gpioPin,
		transmission: make(chan transmission, transmissionChanLen),
		transmitted:  make(chan bool, transmissionChanLen),
		done:         make(chan bool, 1),
	}

	go t.watch()

	return t
}

// Transmit transmits a code using given protocol and pulse length. It will
// return an error if the provided protocol is does not exist.
//
// This method returns immediately. The code is transmitted in the background.
// If you need to ensure that a code has been fully transmitted, call Wait()
// after calling Transmit().
func (t *NativeTransmitter) Transmit(code uint64, protocol int, pulseLength uint) error {
	if protocol < 1 || protocol > len(Protocols) {
		return fmt.Errorf("Protocol %d does not exist", protocol)
	}

	trans := transmission{
		code:        code,
		protocol:    Protocols[protocol-1],
		pulseLength: pulseLength,
	}

	select {
	case t.transmission <- trans:
	default:
	}

	return nil
}

// transmit performs the acutal transmission of the remote control code.
func (t *NativeTransmitter) transmit(trans transmission) {
	for retry := 0; retry < TransmitRetries; retry++ {
		for j := bitLength - 1; j >= 0; j-- {
			if trans.code&(1<<uint64(j)) > 0 {
				t.send(trans.protocol.One, trans.pulseLength)
			} else {
				t.send(trans.protocol.Zero, trans.pulseLength)
			}
		}
		t.send(trans.protocol.Sync, trans.pulseLength)
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

// Close stops started goroutines and closes the gpio pin
func (t *NativeTransmitter) Close() error {
	t.done <- true
	t.gpioPin.Close()

	return nil
}

// Wait blocks until a code is fully transmitted.
func (t *NativeTransmitter) Wait() {
	<-t.transmitted
}

// watch listens on a channel and processes incoming transmissions.
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

// send sends a sequence of high and low pulses on the gpio pin.
func (t *NativeTransmitter) send(pulses HighLow, pulseLength uint) {
	t.gpioPin.High()
	time.Sleep(time.Microsecond * time.Duration(pulseLength*pulses.High))
	t.gpioPin.Low()
	time.Sleep(time.Microsecond * time.Duration(pulseLength*pulses.Low))
}

// NullTransmitter type definition.
type NullTransmitter struct{}

// NewNullTransmitter create a transmitter that does nothing except logging the
// transmissions. This is mainly useful for development on systems where
// /dev/gpiomem is not available.
func NewNullTransmitter() *NullTransmitter {
	return &NullTransmitter{}
}

// Transmit does nothing.
func (t *NullTransmitter) Transmit(code uint64, protocol int, pulseLength uint) error {
	if protocol < 1 || protocol > len(Protocols) {
		return fmt.Errorf("Protocol %d does not exist", protocol)
	}

	return nil
}

// Close does nothing.
func (t *NullTransmitter) Close() error {
	return nil
}

// Wait does nothing.
func (t *NullTransmitter) Wait() {}

// NewTransmitter creates a NativeTransmitter when /dev/gpiomem is available,
// NullTransmitter otherwise.
func NewTransmitter(gpioPin uint) CodeTransmitter {
	if _, err := os.Stat("/dev/gpiomem"); os.IsNotExist(err) {
		return NewNullTransmitter()
	}

	pin := gpio.NewOutput(gpioPin, false)

	return NewNativeTransmitter(pin)
}
