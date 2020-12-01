package gpio

import (
	"time"
	_ "unsafe"
)

// nanotime returns the current time in nanoseconds from a monotonic clock.
//go:linkname nanotime runtime.nanotime
func nanotime() int64

// delay waits for given duration using busy waiting. The godoc for
// time.Sleep states:
//
//   Sleep pauses the current goroutine for *at least* the duration d
//
// This means that for sub-millisecond sleep durations it will pause the
// current goroutine for longer than we can afford as for us the sleep duration
// needs to be as precise as possible to send out the correct codes to the
// outlets.
// time.Sleep causes sleep pauses to be off by 100+ microseconds on average
// whereas, depending on the platform, we can bring this down to ~1
// microseconds at worst using busy waiting in combination with
// runtime.nanotime.
func delay(duration time.Duration) {
	end := nanotime() + int64(duration)

	for nanotime() < end {
	}
}
