package schedule_test

import (
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
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Monday},
				From:     schedule.NewDayTime(0, 0),
				To:       schedule.NewDayTime(12, 35),
			},
			t:        time.Date(2018, 11, 12, 12, 34, 56, 0, time.UTC),
			expected: true,
		},
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Tuesday},
				From:     schedule.NewDayTime(0, 0),
				To:       schedule.NewDayTime(12, 35),
			},
			t:        time.Date(2018, 11, 12, 12, 34, 56, 0, time.UTC),
			expected: false,
		},
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Monday},
				From:     schedule.NewDayTime(0, 0),
				To:       schedule.NewDayTime(12, 35),
			},
			t:        time.Date(2018, 11, 12, 12, 36, 56, 0, time.UTC),
			expected: false,
		},
		{
			i: schedule.Interval{
				Enabled:  true,
				Weekdays: []time.Weekday{time.Monday},
				From:     schedule.NewDayTime(12, 35),
				To:       schedule.NewDayTime(0, 0),
			},
			t:        time.Date(2018, 11, 12, 12, 36, 56, 0, time.UTC),
			expected: true,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.i.Contains(tt.t))
	}
}
