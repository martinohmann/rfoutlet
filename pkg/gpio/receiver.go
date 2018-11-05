package gpio

// Most of the receiver code is ported from the rc-switch c++ implementation to
// go. Check out the rc-switch repository at https://github.com/sui77/rc-switch
// for the original implementation.

import (
	"errors"
	"time"

	"github.com/brian-armstrong/gpio"
)

const (
	receiveTolerance int64 = 60
	separationLimit  int64 = 4600
	maxChanges       uint  = 67
)

// ReceiveResult type definition
type ReceiveResult struct {
	Code        uint64
	BitLength   uint
	PulseLength int64
	Protocol    int
}

// CodeReceiver defines the interface for a rf code receiver.
type CodeReceiver interface {
	Receive() <-chan ReceiveResult
	Close() error
}

// Receiver type definition
type Receiver struct {
	gpioPin     uint
	lastEvent   int64
	changeCount uint
	repeatCount uint
	timings     [maxChanges]int64

	watcher *gpio.Watcher
	closed  bool
	done    chan bool
	result  chan ReceiveResult
}

// NewReceiver create a new receiver on the gpio pin
func NewReceiver(gpioPin uint) *Receiver {
	watcher := gpio.NewWatcher()

	r := &Receiver{
		gpioPin: gpioPin,
		watcher: watcher,
		done:    make(chan bool, 1),
		result:  make(chan ReceiveResult, 1),
	}

	r.watcher.AddPin(r.gpioPin)

	go r.watch()

	return r
}

func (r *Receiver) watch() {
	var lastVal uint

	for {
		select {
		case <-r.done:
			close(r.result)
			return
		default:
			pin, val := r.watcher.Watch()

			if pin == r.gpioPin && val != lastVal {
				r.handleEvent()
			}

			lastVal = val
		}
	}
}

// Receive blocks until there is a result on the receive channel
func (r *Receiver) Receive() <-chan ReceiveResult {
	return r.result
}

// Close stops the watcher and receiver goroutines and perform cleanup
func (r *Receiver) Close() error {
	if r.closed {
		return errors.New("receiver already closed")
	}

	r.closed = true
	r.done <- true
	r.watcher.Close()

	return nil
}

func (r *Receiver) handleEvent() {
	event := time.Now().UnixNano() / int64(time.Microsecond)
	duration := event - r.lastEvent

	if duration > separationLimit {
		if diff(duration, r.timings[0]) < 200 {
			r.repeatCount++

			if r.repeatCount == 2 {
				for i := 1; i <= len(protocols); i++ {
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
	p := protocols[protocol-1]

	var code uint64
	var delay int64 = r.timings[0] / int64(p.sync.low)
	var delayTolerance int64 = delay * receiveTolerance / 100
	var i uint = 1

	for ; i < r.changeCount-1; i += 2 {
		code <<= 1

		if diff(r.timings[i], delay*int64(p.zero.high)) < delayTolerance &&
			diff(r.timings[i+1], delay*int64(p.zero.low)) < delayTolerance {
			// zero
		} else if diff(r.timings[i], delay*int64(p.one.high)) < delayTolerance &&
			diff(r.timings[i+1], delay*int64(p.one.low)) < delayTolerance {
			code |= 1
		} else {
			return false
		}
	}

	if r.changeCount > 7 {
		r.result <- ReceiveResult{
			Code:        code,
			BitLength:   (r.changeCount - 1) / 2,
			PulseLength: delay,
			Protocol:    protocol,
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
