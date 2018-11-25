package scheduler

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/stretchr/testify/assert"
)

type testSwitcher struct {
	count int64
}

func (ts *testSwitcher) Switch(o *outlet.Outlet, s outlet.State) error {
	atomic.AddInt64(&ts.count, 1)
	o.SetState(s)

	return nil
}

func TestScheduler(t *testing.T) {
	ts := &testSwitcher{}

	s := NewWithInterval(ts, 5*time.Millisecond)

	now := time.Now()
	plus1 := now.Add(time.Hour)

	intervals := []schedule.Interval{
		{
			Enabled:  true,
			Weekdays: []time.Weekday{now.Weekday()},
			From:     schedule.NewDayTime(now.Hour(), now.Minute()),
			To:       schedule.NewDayTime(plus1.Hour(), plus1.Minute()),
		},
	}

	o := &outlet.Outlet{Schedule: schedule.NewWithIntervals(intervals)}

	s.Register(o)

	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, outlet.StateOn, o.GetState())
	assert.Equal(t, int64(1), atomic.LoadInt64(&ts.count))
}
