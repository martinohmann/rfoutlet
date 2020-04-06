package gpio

import (
	"sync/atomic"
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

type transmission struct {
	code        uint64
	protocol    Protocol
	pulseLength uint
	done        chan struct{}
}

// Transmitter type definition.
type Transmitter struct {
	pin          OutputPin
	transmission chan transmission
	closed       int32
	retries      int
}

// NewTransmitter creates a Transmitter which attaches to the chip's pin at
// offset.
func NewTransmitter(chip *gpiod.Chip, offset int, options ...TransmitterOption) (*Transmitter, error) {
	pin, err := chip.RequestLine(offset, gpiod.AsOutput(0))
	if err != nil {
		return nil, err
	}

	return NewPinTransmitter(pin, options...), nil
}

// NewTransmitter creates a *Transmitter that sends on pin.
func NewPinTransmitter(pin OutputPin, options ...TransmitterOption) *Transmitter {
	t := &Transmitter{
		pin:          pin,
		transmission: make(chan transmission, transmissionChanLen),
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

// Transmit transmits a code using given protocol and pulse length.
//
// This method returns immediately. The code is transmitted in the background.
// If you need to ensure that a code has been fully transmitted, wait for the
// returned channel to be closed.
func (t *Transmitter) Transmit(code uint64, protocol Protocol, pulseLength uint) <-chan struct{} {
	done := make(chan struct{})

	if atomic.LoadInt32(&t.closed) == 1 {
		close(done)
		return done
	}

	t.transmission <- transmission{
		code:        code,
		protocol:    protocol,
		pulseLength: pulseLength,
		done:        done,
	}

	return done
}

// transmit performs the acutal transmission of the remote control code.
func (t *Transmitter) transmit(trans transmission) {
	defer close(trans.done)

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
}

// Close stops started goroutines and closes the gpio pin.
func (t *Transmitter) Close() error {
	atomic.StoreInt32(&t.closed, 0)
	close(t.transmission)
	return t.pin.Close()
}

// watch listens on a channel and processes incoming transmissions.
func (t *Transmitter) watch() {
	for {
		select {
		case trans, ok := <-t.transmission:
			if !ok {
				return
			}

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
