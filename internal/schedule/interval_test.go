package schedule_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/stretchr/testify/assert"
)

func TestIntervalContains(t *testing.T) {
	tests := []struct {
		i        schedule.Interval
		t        time.Time
		expected bool
	}{
		// From > To
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Monday},
				From:     schedule.NewDayTime(0, 0),
				To:       schedule.NewDayTime(3, 0),
			},
			t:        time.Date(2018, 11, 5, 0, 0, 0, 0, time.UTC), // it's a monday
			expected: true,
		},
		{
			i: schedule.Interval{
				Enabled:  false,
				Weekdays: []time.Weekday{time.Monday},
				From:     schedule.NewDayTime(0, 0),
				To:       schedule.NewDayTime(3, 0),
			},
			t:        time.Date(2018, 11, 5, 0, 0, 0, 0, time.UTC), // it's a monday
			expected: false,
		},
		{
			i: schedule.Interval{
				Enabled:  false,
				Weekdays: []time.Weekday{time.Monday},
				From:     schedule.NewDayTime(0, 0),
				To:       schedule.NewDayTime(3, 0),
			},
			t:        time.Date(2018, 11, 5, 23, 59, 0, 0, time.UTC), // it's a monday
			expected: false,
		},
		{
			i: schedule.Interval{
				Enabled:  false,
				Weekdays: []time.Weekday{time.Monday},
				From:     schedule.NewDayTime(0, 0),
				To:       schedule.NewDayTime(3, 0),
			},
			t:        time.Date(2018, 11, 5, 1, 30, 0, 0, time.UTC), // it's a monday
			expected: false,
		},
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Monday},
				From:     schedule.NewDayTime(0, 0),
				To:       schedule.NewDayTime(3, 0),
			},
			t:        time.Date(2018, 11, 6, 0, 0, 0, 0, time.UTC), // it's a tuesday
			expected: false,
		},
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Monday, time.Tuesday},
				From:     schedule.NewDayTime(0, 0),
				To:       schedule.NewDayTime(3, 0),
			},
			t:        time.Date(2018, 11, 6, 0, 0, 0, 0, time.UTC), // it's a tuesday
			expected: true,
		},
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Monday, time.Tuesday},
				From:     schedule.NewDayTime(0, 0),
				To:       schedule.NewDayTime(3, 0),
			},
			t:        time.Date(2018, 11, 6, 3, 0, 0, 0, time.UTC), // it's a tuesday
			expected: false,
		},
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Monday, time.Tuesday},
				From:     schedule.NewDayTime(0, 0),
				To:       schedule.NewDayTime(3, 0),
			},
			t:        time.Date(2018, 11, 6, 2, 59, 0, 0, time.UTC), // it's a tuesday
			expected: true,
		},
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Monday, time.Tuesday},
				From:     schedule.NewDayTime(0, 0),
				To:       schedule.NewDayTime(3, 0),
			},
			t:        time.Date(2018, 11, 6, 3, 1, 0, 0, time.UTC), // it's a tuesday
			expected: false,
		},
		// From < To
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Tuesday},
				From:     schedule.NewDayTime(3, 0),
				To:       schedule.NewDayTime(0, 0),
			},
			t:        time.Date(2018, 11, 6, 6, 35, 0, 0, time.UTC), // it's a tuesday
			expected: true,
		},
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Tuesday},
				From:     schedule.NewDayTime(3, 0),
				To:       schedule.NewDayTime(0, 0),
			},
			t:        time.Date(2018, 11, 6, 3, 0, 0, 0, time.UTC), // it's a tuesday
			expected: true,
		},
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Tuesday},
				From:     schedule.NewDayTime(3, 0),
				To:       schedule.NewDayTime(0, 0),
			},
			t:        time.Date(2018, 11, 6, 0, 0, 0, 0, time.UTC), // it's a tuesday
			expected: false,
		},
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Tuesday},
				From:     schedule.NewDayTime(3, 0),
				To:       schedule.NewDayTime(0, 0),
			},
			t:        time.Date(2018, 11, 6, 2, 59, 0, 0, time.UTC), // it's a tuesday
			expected: false,
		},
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Tuesday},
				From:     schedule.NewDayTime(3, 0),
				To:       schedule.NewDayTime(0, 0),
			},
			t:        time.Date(2018, 11, 6, 0, 1, 0, 0, time.UTC), // it's a tuesday
			expected: false,
		},
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Tuesday},
				From:     schedule.NewDayTime(3, 0),
				To:       schedule.NewDayTime(0, 0),
			},
			t:        time.Date(2018, 11, 6, 23, 59, 0, 0, time.UTC), // it's a tuesday
			expected: true,
		},
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Tuesday},
				From:     schedule.NewDayTime(3, 0),
				To:       schedule.NewDayTime(0, 0),
			},
			t:        time.Date(2018, 11, 6, 3, 1, 0, 0, time.UTC), // it's a tuesday
			expected: true,
		},
		// From == To
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Tuesday},
				From:     schedule.NewDayTime(0, 0),
				To:       schedule.NewDayTime(0, 0),
			},
			t:        time.Date(2018, 11, 6, 0, 0, 0, 0, time.UTC), // it's a tuesday
			expected: true,
		},
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Tuesday},
				From:     schedule.NewDayTime(0, 0),
				To:       schedule.NewDayTime(0, 0),
			},
			t:        time.Date(2018, 11, 6, 0, 1, 0, 0, time.UTC), // it's a tuesday
			expected: false,
		},
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Tuesday},
				From:     schedule.NewDayTime(0, 0),
				To:       schedule.NewDayTime(0, 0),
			},
			t:        time.Date(2018, 11, 6, 23, 59, 0, 0, time.UTC), // it's a tuesday
			expected: false,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.i.Contains(tt.t),
			fmt.Sprintf("i=%v, t=%v", tt.i, tt.t))
	}
}
