package testutil

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func Retry(t *testing.T, attempts int, sleep time.Duration, fn func(*R)) bool {
	for attempt := 1; attempt <= attempts; attempt++ {
		r := &R{log: &bytes.Buffer{}}

		fn(r)

		if !r.failed {
			return true
		}

		if attempt == attempts {
			t.Errorf("Failed after %d attempts: %s", attempt, r.log.String())
			break
		}

		time.Sleep(sleep)
	}

	return false
}

type R struct {
	failed bool
	log    *bytes.Buffer
}

func (r *R) Errorf(format string, args ...interface{}) {
	fmt.Fprintln(r.log)
	fmt.Fprintf(r.log, format, args...)
	r.failed = true
}
