package gpio

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/warthog618/gpiod"
)

type pinWatcherPipe struct {
	*FakeWatcher
	clock clockwork.Clock
	once  sync.Once
}

func newPinWatcherPipe(clock clockwork.Clock) *pinWatcherPipe {
	return &pinWatcherPipe{
		FakeWatcher: NewFakeWatcher(),
		clock:       clock,
	}
}

func (p *pinWatcherPipe) SetValue(value int) error {
	var eventType gpiod.LineEventType
	switch value {
	case 0:
		eventType = gpiod.LineEventFallingEdge
	case 1:
		eventType = gpiod.LineEventRisingEdge
	default:
		panic(fmt.Sprintf("invalid value: %d", value))
	}

	p.Events <- gpiod.LineEvent{
		Type:      eventType,
		Timestamp: time.Duration(p.clock.Now().UnixNano()),
	}

	return nil
}

func (p *pinWatcherPipe) Close() (err error) {
	p.once.Do(func() {
		err = p.FakeWatcher.Close()
	})
	return
}

func TestReceiverReceive(t *testing.T) {
	fakeClock := clockwork.NewFakeClockAt(time.Now())

	pipe := newPinWatcherPipe(fakeClock)

	tx := NewPinTransmitter(pipe, TransmissionCount(10))
	tx.delay = fakeClock.Advance
	defer tx.Close()

	rx := NewWatcherReceiver(pipe, ReceiverProtocols(DefaultProtocols))
	defer rx.Close()

	transmissions := []struct {
		code        uint64
		protocol    int
		pulseLength uint
	}{
		{5510451, 1, 184},
		{83281, 1, 305},
		{86356, 1, 305},
		{5510604, 1, 184},
		{5591317, 1, 330},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go func() {
		defer cancel()

		var i int
		var lastCode uint64

		for i < len(transmissions) {
			result, ok := <-rx.Receive()
			if !ok {
				return
			}
			if result.Code == lastCode {
				continue
			}

			tm := transmissions[i]

			assert.Equalf(t, tm.code, result.Code,
				"received code %d != expected %d", result.Code, tm.code)

			lastCode = result.Code
			i++
		}
	}()

	for _, tm := range transmissions {
		<-tx.Transmit(tm.code, DefaultProtocols[tm.protocol-1], tm.pulseLength)
	}

	<-ctx.Done()

	if err := ctx.Err(); err == context.DeadlineExceeded {
		t.Fatal(err)
	}
}

func TestReceiverClose(t *testing.T) {
	w := NewFakeWatcher()

	receiver := NewWatcherReceiver(w)
	assert.Nil(t, receiver.Close())

	assert.True(t, w.Closed)
}
