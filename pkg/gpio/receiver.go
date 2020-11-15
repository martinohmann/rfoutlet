package gpio

import (
	"time"

	"github.com/warthog618/gpiod"
)

const (
	receiveTolerance int64 = 60
	separationLimit  int64 = 4600
	maxChanges       uint  = 67

	receiveResultChanLen = 32
)

// Receiver can detect and unserialize rf codes received on a gpio pin.
type Receiver struct {
	lastEvent   int64
	changeCount uint
	repeatCount uint
	timings     [maxChanges]int64

	watcher   Watcher
	protocols []Protocol
	result    chan ReceiveResult
}

// NewReceiver creates a *Receiver which listens on the chip's pin at offset
// for rf codes.
func NewReceiver(chip *gpiod.Chip, offset int, options ...ReceiverOption) (*Receiver, error) {
	watcher, err := NewWatcher(chip, offset)
	if err != nil {
		return nil, err
	}

	return NewWatcherReceiver(watcher, options...), nil
}

// NewWatcherReceiver create a new receiver which uses given Watcher to detect
// sent rf codes.
func NewWatcherReceiver(watcher Watcher, options ...ReceiverOption) *Receiver {
	r := &Receiver{
		watcher:   watcher,
		result:    make(chan ReceiveResult, receiveResultChanLen),
		protocols: DefaultProtocols,
	}

	for _, option := range options {
		option(r)
	}

	go r.watch()

	return r
}

func (r *Receiver) watch() {
	defer close(r.result)

	var lastEventType gpiod.LineEventType

	for evt := range r.watcher.Watch() {
		if lastEventType != evt.Type {
			r.handleEvent(evt)
			lastEventType = evt.Type
		}
	}
}

// Receive blocks until there is a result on the receive channel
func (r *Receiver) Receive() <-chan ReceiveResult {
	return r.result
}

// Close stops the watcher and receiver goroutines and perform cleanup.
func (r *Receiver) Close() error {
	return r.watcher.Close()
}

func (r *Receiver) handleEvent(evt gpiod.LineEvent) {
	event := int64(evt.Timestamp) / int64(time.Microsecond)
	duration := event - r.lastEvent

	if duration > separationLimit {
		if diff(duration, r.timings[0]) < 200 {
			r.repeatCount++

			if r.repeatCount == 2 {
				for i := 0; i < len(r.protocols); i++ {
					if r.receiveProtocol(i) {
						break
					}
				}

				r.repeatCount = 0
			}
		}

		r.changeCount = 0
	}

	if r.changeCount >= maxChanges {
		r.changeCount = 0
		r.repeatCount = 0
	}

	r.timings[r.changeCount] = duration
	r.changeCount++
	r.lastEvent = event
}

// receiveProtocol tries to receive a code using the provided protocol
func (r *Receiver) receiveProtocol(protocol int) bool {
	p := r.protocols[protocol]

	delay := r.timings[0] / int64(p.Sync.Low)
	delayTolerance := delay * receiveTolerance / 100

	var code uint64
	var i uint = 1

	for ; i < r.changeCount-1; i += 2 {
		code <<= 1

		if diff(r.timings[i], delay*int64(p.Zero.High)) < delayTolerance &&
			diff(r.timings[i+1], delay*int64(p.Zero.Low)) < delayTolerance {
			// zero
		} else if diff(r.timings[i], delay*int64(p.One.High)) < delayTolerance &&
			diff(r.timings[i+1], delay*int64(p.One.Low)) < delayTolerance {
			code |= 1
		} else {
			return false
		}
	}

	if r.changeCount > 7 {
		result := ReceiveResult{
			Code:        code,
			BitLength:   (r.changeCount - 1) / 2,
			PulseLength: delay,
			Protocol:    protocol + 1,
		}

		select {
		case r.result <- result:
		default:
		}
	}

	return true
}

func diff(a, b int64) int64 {
	v := a - b

	if v < 0 {
		return -v
	}

	return v
}
