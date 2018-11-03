package gpio

import (
	"time"

	"github.com/brian-armstrong/gpio"
)

const (
	receiveTolerance int64 = 60
	separationLimit  int64 = 4600
)

type ReceiveFunc func(uint64, int64, uint, int)

type Receiver interface {
	Receive(ReceiveFunc)
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

	timings [67]int64

	watcher *gpio.Watcher
}

func NewNativeReceiver(gpioPin uint) *NativeReceiver {
	watcher := gpio.NewWatcher()

	r := &NativeReceiver{
		gpioPin: gpioPin,
		watcher: watcher,
	}

	return r
}

func (r *NativeReceiver) Receive(f ReceiveFunc) {
	r.watcher.AddPin(r.gpioPin)

	go func() {
		var lastValue uint

		for {
			_, value := r.watcher.Watch()

			if value != lastValue {
				r.handleInterrupt()
			}

			lastValue = value
		}
	}()

	for {
		if r.hasReceivedCode() {
			f(r.ReceivedCode, r.ReceivedPulseLength, r.ReceivedBitLength, r.ReceivedProtocol)

			r.reset()
		}

		time.Sleep(time.Microsecond * 100)
	}
}

func (r *NativeReceiver) Close() error {
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
				r.receiveProtocol(1, r.changeCount)
				r.repeatCount = 0
			}
		}

		r.changeCount = 0
	}

	if r.changeCount >= 67 {
		r.changeCount = 0
		r.repeatCount = 0
	}

	r.timings[r.changeCount] = duration
	r.changeCount++
	r.lastTime = t
}

func (r *NativeReceiver) receiveProtocol(p int, changeCount uint) bool {
	var code uint64
	var delay int64 = r.timings[0] / 31
	var delayTolerance int64 = delay * receiveTolerance / 100

	var i uint

	for i = 1; i < changeCount-1; i += 2 {
		code <<= 1

		if diff(r.timings[i], delay*1) < delayTolerance && diff(r.timings[i+1], delay*3) < delayTolerance {
			// zero
		} else if diff(r.timings[i], delay*3) < delayTolerance && diff(r.timings[i+1], delay*1) < delayTolerance {
			code |= 1
		} else {
			return false
		}
	}

	if changeCount > 7 {
		r.ReceivedCode = code
		r.ReceivedBitLength = (r.changeCount - 1) / 2
		r.ReceivedPulseLength = delay
		r.ReceivedProtocol = p
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
