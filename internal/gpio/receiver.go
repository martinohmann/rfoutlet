package gpio

// Most of the receiver code is ported from the rc-switch c++ implementation to
// go. Check out the rc-switch repository at https://github.com/sui77/rc-switch
// for the original implementation.

import (
	"sync"
	"time"

	"github.com/brian-armstrong/gpio"
)

const (
	receiveTolerance int64 = 60
	separationLimit  int64 = 4600
	maxChanges       uint  = 67
)

// ReceiveFunc type definition
type ReceiveFunc func(uint64, int64, uint, int)

// CodeReceiver defines the interface for a rf code receiver.
type CodeReceiver interface {
	Receive(ReceiveFunc)
	Wait()
	Close() error
}

// Receiver type definition
type Receiver struct {
	gpioPin     uint
	lastTime    int64
	changeCount uint
	repeatCount uint

	ReceivedCode        uint64
	ReceivedBitLength   uint
	ReceivedPulseLength int64
	ReceivedProtocol    int

	timings [maxChanges]int64

	watcher       *gpio.Watcher
	stopWatching  chan bool
	stopReceiving chan bool
	wg            sync.WaitGroup
}

// NewReceiver create a new receiver on the gpio pin
func NewReceiver(gpioPin uint) *Receiver {
	watcher := gpio.NewWatcher()

	r := &Receiver{
		gpioPin:       gpioPin,
		watcher:       watcher,
		stopWatching:  make(chan bool),
		stopReceiving: make(chan bool),
	}

	return r
}

// Receive starts the goroutines for the watcher and receiver. The receive
// function will be called whenever a code has been received.
func (r *Receiver) Receive(f ReceiveFunc) {
	r.watcher.AddPin(r.gpioPin)
	r.wg.Add(2)

	go r.watch(r.stopWatching)
	go r.receive(f, r.stopReceiving)
}

func (r *Receiver) watch(done chan bool) {
	defer r.wg.Done()

	var lastValue uint

	for {
		select {
		case <-done:
			return
		default:
			pin, value := r.watcher.Watch()

			if pin == r.gpioPin && value != lastValue {
				r.handleInterrupt()
			}

			lastValue = value
		}
	}
}

func (r *Receiver) receive(f ReceiveFunc, done chan bool) {
	defer r.wg.Done()

	for {
		select {
		case <-done:
			return
		default:
			if r.hasReceivedCode() {
				f(r.ReceivedCode, r.ReceivedPulseLength, r.ReceivedBitLength, r.ReceivedProtocol)

				r.reset()
			}

			time.Sleep(time.Microsecond * 100)
		}
	}
}

// Wait blocks until Close is called
func (r *Receiver) Wait() {
	r.wg.Wait()
}

// Close stops the watcher and receiver goroutines and perform cleanup
func (r *Receiver) Close() error {
	r.stopReceiving <- true
	r.stopWatching <- true
	r.watcher.Close()

	return nil
}

func (r *Receiver) hasReceivedCode() bool {
	return r.ReceivedCode != 0
}

func (r *Receiver) reset() {
	r.ReceivedCode = 0
}

func (r *Receiver) handleInterrupt() {
	t := time.Now().UnixNano() / int64(time.Microsecond)
	duration := t - r.lastTime

	if duration > separationLimit {
		if diff(duration, r.timings[0]) < 200 {
			r.repeatCount++

			if r.repeatCount == 2 {
				for i := 1; i <= len(protocols); i++ {
					if r.receiveProtocol(i, r.changeCount) {
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
	r.lastTime = t
}

// receiveProtocol tries to receive a code using the provided protocol
func (r *Receiver) receiveProtocol(protocol int, changeCount uint) bool {
	p := protocols[protocol-1]

	var code uint64
	var delay int64 = r.timings[0] / int64(p.sync.low)
	var delayTolerance int64 = delay * receiveTolerance / 100
	var i uint = 1

	for ; i < changeCount-1; i += 2 {
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

	if changeCount > 7 {
		r.ReceivedCode = code
		r.ReceivedBitLength = (changeCount - 1) / 2
		r.ReceivedPulseLength = delay
		r.ReceivedProtocol = protocol
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
