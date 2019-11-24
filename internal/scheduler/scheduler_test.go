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
	now := time.Now()
	plus1 := now.Add(time.Hour)

	tests := []struct {
		outlet               *outlet.Outlet
		intervals            []schedule.Interval
		expectedState        outlet.State
		expectedStateChanges int64
	}{
		{
			outlet: &outlet.Outlet{},
			intervals: []schedule.Interval{
				{
					Enabled:  true,
					Weekdays: []time.Weekday{now.Weekday()},
					From:     schedule.NewDayTime(now.Hour(), now.Minute()),
					To:       schedule.NewDayTime(plus1.Hour(), plus1.Minute()),
				},
			},
			expectedState:        outlet.StateOn,
			expectedStateChanges: 1,
		},
		{
			outlet: &outlet.Outlet{State: outlet.StateOn},
			intervals: []schedule.Interval{
				{
					Enabled:  true,
					Weekdays: []time.Weekday{now.Weekday()},
					From:     schedule.NewDayTime(plus1.Hour(), plus1.Minute()),
					To:       schedule.NewDayTime(now.Hour(), now.Minute()),
				},
			},
			expectedState:        outlet.StateOff,
			expectedStateChanges: 1,
		},
		{
			outlet: &outlet.Outlet{State: outlet.StateOn},
			intervals: []schedule.Interval{
				{
					Enabled:  false,
					Weekdays: []time.Weekday{now.Weekday()},
					From:     schedule.NewDayTime(now.Hour(), now.Minute()),
					To:       schedule.NewDayTime(plus1.Hour(), plus1.Minute()),
				},
			},
			expectedState:        outlet.StateOn,
			expectedStateChanges: 0,
		},
	}

	for _, tt := range tests {
		tt.outlet.Schedule = schedule.NewWithIntervals(tt.intervals)
		testScheduler(t, tt.outlet, tt.expectedState, tt.expectedStateChanges)
	}
}

func testScheduler(t *testing.T, o *outlet.Outlet, expectedState outlet.State, expectedStateChanges int64) {
	ts := &testSwitcher{}

	s := NewWithInterval(ts, 5*time.Millisecond)
	defer s.ticker.Stop()

	s.Register(o)

	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, expectedState, o.GetState())
	assert.Equal(t, expectedStateChanges, atomic.LoadInt64(&ts.count))
}
