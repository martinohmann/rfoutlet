package gpio

import (
	"fmt"
	"time"

	"github.com/warthog618/gpiod"
)

const (
	// DefaultTransmissionRetries defines how many times a code should be
	// transmitted in a row by default.
	DefaultTransmissionRetries = 10

	transmissionChanLen = 32
	bitLength           = 24
)

// TransmitterOption is the signature of funcs that are used to configure a
// *Transmitter.
type TransmitterOption func(*Transmitter)

// TransmissionRetries configures how many times a code should be transmitted
// in a row. The higher the value, the more likely it is that an outlet
// actually received the code.
func TransmissionRetries(retries int) func(*Transmitter) {
	return func(t *Transmitter) {
		t.retries = retries
	}
}

type transmission struct {
	code        uint64
	protocol    Protocol
	pulseLength uint
}

// Transmitter type definition.
type Transmitter struct {
	pin          OutputPin
	transmission chan transmission
	transmitted  chan bool
	done         chan bool
	retries      int
}

// NewTransmitter creates a Transmitter which attaches to the chip's pin at
// offset.
func NewTransmitter(chip *gpiod.Chip, offset int, options ...TransmitterOption) (*Transmitter, error) {
	line, err := chip.RequestLine(offset, gpiod.AsOutput(0))
	if err != nil {
		return nil, err
	}

	return NewPinTransmitter(line, options...), nil
}

// NewTransmitter creates a *Transmitter that sends on pin.
func NewPinTransmitter(pin OutputPin, options ...TransmitterOption) *Transmitter {
	t := &Transmitter{
		pin:          pin,
		transmission: make(chan transmission, transmissionChanLen),
		transmitted:  make(chan bool, transmissionChanLen),
		done:         make(chan bool, 1),
		retries:      DefaultTransmissionRetries,
	}

	for _, option := range options {
		option(t)
	}

	if t.retries <= 0 {
		t.retries = 1
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
func (t *Transmitter) Transmit(code uint64, protocol int, pulseLength uint) error {
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
func (t *Transmitter) transmit(trans transmission) {
	for retry := 0; retry < t.retries; retry++ {
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
}

// Close stops started goroutines and closes the gpio pin
func (t *Transmitter) Close() error {
	t.done <- true
	t.pin.Close()

	return nil
}

// Wait blocks until a code is fully transmitted.
func (t *Transmitter) Wait() {
	<-t.transmitted
}

// watch listens on a channel and processes incoming transmissions.
func (t *Transmitter) watch() {
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
func (t *Transmitter) send(pulses HighLow, pulseLength uint) {
	t.pin.SetValue(1)
	sleepFor(time.Microsecond * time.Duration(pulseLength*pulses.High))
	t.pin.SetValue(0)
	sleepFor(time.Microsecond * time.Duration(pulseLength*pulses.Low))
}

// NewDiscardingTransmitter creates a *Transmitter that does not send anything.
func NewDiscardingTransmitter() *Transmitter {
	return NewPinTransmitter(&FakeOutputPin{})
}

// sleepFor sleeps for given duration using busy waiting. The godoc for
// time.Sleep states:
//
//   Sleep pauses the current goroutine for *at least* the duration d
//
// This means that for sub-millisecond sleep durations it will pause the
// current goroutine for longer than we can afford as for us the sleep duration
// needs to be as precise as possible to send out the correct codes to the
// outlets. time.Sleep causes sleep pauses to be off by 100+ microseconds on
// average whereas we can bring this down to < 5 microseconds using busy
// waiting.
func sleepFor(duration time.Duration) {
	now := time.Now()

	for {
		if time.Since(now) >= duration {
			break
		}
	}
}
