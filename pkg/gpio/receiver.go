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
type ReceiveFunc func(ReceiveResult)

// ReceiveResult type definition
type ReceiveResult struct {
	Code        uint64
	BitLength   uint
	PulseLength int64
	Protocol    int
}

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
	timings     [maxChanges]int64

	watcher       *gpio.Watcher
	wg            sync.WaitGroup
	done          chan bool
	receiveResult chan ReceiveResult
}

// NewReceiver create a new receiver on the gpio pin
func NewReceiver(gpioPin uint) *Receiver {
	watcher := gpio.NewWatcher()

	r := &Receiver{
		gpioPin:       gpioPin,
		watcher:       watcher,
		done:          make(chan bool),
		receiveResult: make(chan ReceiveResult),
	}

	return r
}

// Receive starts the goroutines for the watcher and receiver. The receive
// function will be called whenever a code has been received.
func (r *Receiver) Receive(f ReceiveFunc) {
	r.watcher.AddPin(r.gpioPin)
	r.wg.Add(1)

	go r.receive(f)
}

func (r *Receiver) receive(f ReceiveFunc) {
	defer r.wg.Done()

	var lastValue uint

	for {
		select {
		case <-r.done:
			return
		case res := <-r.receiveResult:
			f(res)
		default:
			pin, value := r.watcher.Watch()

			if pin == r.gpioPin && value != lastValue {
				r.handleInterrupt()
			}

			lastValue = value
		}
	}
}

// Wait blocks until Close is called
func (r *Receiver) Wait() {
	r.wg.Wait()
}

// Close stops the watcher and receiver goroutines and perform cleanup
func (r *Receiver) Close() error {
	r.done <- true
	r.watcher.Close()

	return nil
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
		r.receiveResult <- ReceiveResult{
			Code:        code,
			BitLength:   (changeCount - 1) / 2,
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
