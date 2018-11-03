package gpio

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

type ReceiveFunc func(uint64, int64, uint, int)

type Receiver interface {
	Receive(ReceiveFunc)
	Wait()
	Close() error
}

type NativeReceiver struct {
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

func NewNativeReceiver(gpioPin uint) *NativeReceiver {
	watcher := gpio.NewWatcher()

	r := &NativeReceiver{
		gpioPin:       gpioPin,
		watcher:       watcher,
		stopWatching:  make(chan bool),
		stopReceiving: make(chan bool),
	}

	return r
}

func (r *NativeReceiver) Receive(f ReceiveFunc) {
	r.watcher.AddPin(r.gpioPin)

	go r.watch(r.stopWatching)
	go r.receive(f, r.stopReceiving)
}

func (r *NativeReceiver) watch(done chan bool) {
	r.wg.Add(1)

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

func (r *NativeReceiver) receive(f ReceiveFunc, done chan bool) {
	r.wg.Add(1)

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

func (r *NativeReceiver) Wait() {
	r.wg.Wait()
}

func (r *NativeReceiver) Close() error {
	r.stopReceiving <- true
	r.stopWatching <- true
	r.watcher.Close()

	return nil
}

func (r *NativeReceiver) hasReceivedCode() bool {
	return r.ReceivedCode != 0
}

func (r *NativeReceiver) reset() {
	r.ReceivedCode = 0
}

func (r *NativeReceiver) handleInterrupt() {
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

func (r *NativeReceiver) receiveProtocol(protocol int, changeCount uint) bool {
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
		r.ReceivedBitLength = (r.changeCount - 1) / 2
		r.ReceivedPulseLength = delay
		r.ReceivedProtocol = protocol
	}

	return true
}

func NewReceiver(gpioPin uint) Receiver {
	return NewNativeReceiver(gpioPin)
}

func diff(a, b int64) int64 {
	v := a - b

	if v < 0 {
		return -v
	}

	return v
}
