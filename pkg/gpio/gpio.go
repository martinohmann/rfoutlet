package gpio

import "github.com/warthog618/gpiod"

// Closer is something that can be closed.
type Closer interface {
	// Close closes the thing.
	Close() error
}

// Watcher watches a pin.
type Watcher interface {
	Closer
	// Watch returns a channel of observed RisingEdge and FallingEdge pin
	// events.
	Watch() <-chan gpiod.LineEvent
}

// OutputPin is a pin that can be written to.
type OutputPin interface {
	Closer
	// SetValue sets the pin value.
	SetValue(value int) error
}

// CodeTransmitter defines the interface for a rf code transmitter.
type CodeTransmitter interface {
	Closer
	// Transmit transmits a code using given protocol and pulse length.
	//
	// This method returns immediately. The code is transmitted in the background.
	// If you need to ensure that a code has been fully transmitted, wait for the
	// returned channel to be closed.
	Transmit(code uint64, protocol Protocol, pulseLength uint) <-chan struct{}
}

// CodeReceiver defines the interface for a rf code receiver.
type CodeReceiver interface {
	Closer
	// Receive blocks until there is a result on the receive channel.
	Receive() <-chan ReceiveResult
}

// ReceiveResult contains information about a detected code sent by an rf code
// transmitter.
type ReceiveResult struct {
	// Code is the detected code.
	Code uint64

	// BitLength is the detected bit length.
	BitLength uint

	// PulseLength is the detected pulse length.
	PulseLength int64

	// Protocol is the detected protocol. The protocol is 1-indexed.
	Protocol int
}

// FakeWatcher can be used in tests as a Watcher.
type FakeWatcher struct {
	// Events can be used to make the FakeWatcher return arbitrary events.
	Events chan gpiod.LineEvent

	// Err controls the error returned by Close.
	Err error

	// Closed indicates whether Close was called or not.
	Closed bool
}

// NewFakeWatcher creates a new *FakeWatcher that can be used in tests as a
// Watcher.
func NewFakeWatcher() *FakeWatcher {
	return &FakeWatcher{
		Events: make(chan gpiod.LineEvent),
	}
}

// Watch implements Watcher.
func (w *FakeWatcher) Watch() <-chan gpiod.LineEvent {
	return w.Events
}

// Close implements Closer.
func (w *FakeWatcher) Close() error {
	close(w.Events)
	w.Closed = true
	return w.Err
}

// FakeOutputPin can be used in tests as an OutputPin.
type FakeOutputPin struct {
	// Values holds the sequence of values the were set via SetValue.
	Values []int

	// Err controls the error returned by Close.
	Err error

	// Closed indicates whether Close was called or not.
	Closed bool
}

// NewFakeOutputPin creates a new *FakeOutputPin that can be used in tests as
// an OutputPin.
func NewFakeOutputPin() *FakeOutputPin {
	return &FakeOutputPin{
		Values: make([]int, 0),
	}
}

// SetValue implements OutputPin.
func (p *FakeOutputPin) SetValue(value int) error {
	p.Values = append(p.Values, value)
	return p.Err
}

// Close implements Closer.
func (p *FakeOutputPin) Close() error {
	p.Closed = true
	return p.Err
}
