package gpio

import (
	"time"

	"github.com/brian-armstrong/gpio"
)

const (
	receiveTolerance int64 = 60
	separationLimit  int64 = 4600
	maxChanges       uint  = 67

	receiveResultChanLen = 32
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

// Watcher defines the interface for a gpio pin watcher
type Watcher interface {
	Watch() (uint, uint)
	AddPin(uint)
	Close()
}

// NativeReceiver type definition
type NativeReceiver struct {
	gpioPin     uint
	lastEvent   int64
	changeCount uint
	repeatCount uint
	timings     [maxChanges]int64

	watcher Watcher
	done    chan bool
	result  chan ReceiveResult
}

// NewNativeReceiver create a new receiver on the gpio pin using watcher
func NewNativeReceiver(gpioPin uint, watcher Watcher) *NativeReceiver {
	r := &NativeReceiver{
		gpioPin: gpioPin,
		watcher: watcher,
		done:    make(chan bool, 1),
		result:  make(chan ReceiveResult, receiveResultChanLen),
	}

	r.watcher.AddPin(r.gpioPin)

	go r.watch()

	return r
}

func (r *NativeReceiver) watch() {
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
func (r *NativeReceiver) Receive() <-chan ReceiveResult {
	return r.result
}

// Close stops the watcher and receiver goroutines and perform cleanup
func (r *NativeReceiver) Close() error {
	r.done <- true
	r.watcher.Close()

	return nil
}

func (r *NativeReceiver) handleEvent() {
	event := time.Now().UnixNano() / int64(time.Microsecond)
	duration := event - r.lastEvent

	if duration > separationLimit {
		if diff(duration, r.timings[0]) < 200 {
			r.repeatCount++

			if r.repeatCount == 2 {
				for i := 1; i <= len(Protocols); i++ {
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
func (r *NativeReceiver) receiveProtocol(protocol int) bool {
	p := Protocols[protocol-1]

	var code uint64
	var delay int64 = r.timings[0] / int64(p.Sync.Low)
	var delayTolerance int64 = delay * receiveTolerance / 100
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
			Protocol:    protocol,
		}

		select {
		case r.result <- result:
		default:
		}
	}

	return true
}

// NewReceiver create a new receiver on the gpio pin
func NewReceiver(gpioPin uint) CodeReceiver {
	w := gpio.NewWatcher()

	return NewNativeReceiver(gpioPin, w)
}

func diff(a, b int64) int64 {
	v := a - b

	if v < 0 {
		return -v
	}

	return v
}
