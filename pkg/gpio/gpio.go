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
	// Transmit transmits a code using given protocol and pulse length. It will
	// return an error if the provided protocol is does not exist.
	//
	// This method returns immediately. The code is transmitted in the background.
	// If you need to ensure that a code has been fully transmitted, call Wait()
	// after calling Transmit().
	Transmit(uint64, int, uint) error

	// Wait blocks until a code is fully transmitted.
	Wait()
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

	// Protocol is the detected protocol.
	Protocol int
}

// HighLow type definition
type HighLow struct {
	High, Low uint
}

// Protocol type definition
type Protocol struct {
	PulseLength     uint
	Sync, Zero, One HighLow
}

// Protocols defines known remote control protocols. These are exported to give
// users the ability to add more protocols if needed.
var Protocols = []Protocol{
	{350, HighLow{1, 31}, HighLow{1, 3}, HighLow{3, 1}},
	{650, HighLow{1, 10}, HighLow{1, 2}, HighLow{2, 1}},
	{100, HighLow{30, 71}, HighLow{4, 11}, HighLow{9, 6}},
	{380, HighLow{1, 6}, HighLow{1, 3}, HighLow{3, 1}},
	{500, HighLow{6, 14}, HighLow{1, 2}, HighLow{2, 1}},
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

// Watch implements Watcher.
func (w *FakeWatcher) Watch() <-chan gpiod.LineEvent {
	return w.Events
}

// Close implements Closer.
func (w *FakeWatcher) Close() error {
	defer close(w.Events)
	w.Closed = true
	return w.Err
}

// FakeOutputPin can be used in tests as an OutputPin.
type FakeOutputPin struct {
	// Value holds the value the was set via SetValue.
	Value int

	// Err controls the error returned by Close.
	Err error

	// Closed indicates whether Close was called or not.
	Closed bool
}

// SetValue implements OutputPin.
func (p *FakeOutputPin) SetValue(value int) error {
	p.Value = value
	return p.Err
}

// Close implements Closer.
func (p *FakeOutputPin) Close() error {
	p.Closed = true
	return p.Err
}
