package schedule_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/stretchr/testify/assert"
)

func TestScheduleEnabled(t *testing.T) {
	tests := []struct {
		s        schedule.Schedule
		expected bool
	}{
		{
			s: schedule.Schedule{
				{
					Enabled: true,
				},
				{
					Enabled: true,
				},
			},
			expected: true,
		},
		{
			s: schedule.Schedule{
				{
					Enabled: true,
				},
				{
					Enabled: false,
				},
			},
			expected: true,
		},
		{
			s: schedule.Schedule{
				{
					Enabled: false,
				},
				{
					Enabled: true,
				},
			},
			expected: true,
		},
		{
			s: schedule.Schedule{
				{
					Enabled: false,
				},
				{
					Enabled: false,
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.s.Enabled(),
			fmt.Sprintf("i=%v", tt.s))
	}
}

func TestScheduleContains(t *testing.T) {
	tests := []struct {
		s        schedule.Schedule
		t        time.Time
		expected bool
	}{
		{
			s: schedule.Schedule{
				{
					Enabled:  true,
					Weekdays: []time.Weekday{time.Monday},
					From:     schedule.NewDayTime(0, 0),
					To:       schedule.NewDayTime(3, 0),
				},
			},
			t:        time.Date(2018, 11, 5, 0, 0, 0, 0, time.UTC), // it's a monday
			expected: true,
		},
		{
			s: schedule.Schedule{
				{
					Enabled:  true,
					Weekdays: []time.Weekday{time.Tuesday},
					From:     schedule.NewDayTime(0, 0),
					To:       schedule.NewDayTime(3, 0),
				},
			},
			t:        time.Date(2018, 11, 5, 0, 0, 0, 0, time.UTC), // it's a monday
			expected: false,
		},
		{
			s: schedule.Schedule{
				{
					Enabled:  false,
					Weekdays: []time.Weekday{time.Monday},
					From:     schedule.NewDayTime(0, 0),
					To:       schedule.NewDayTime(3, 0),
				},
			},
			t:        time.Date(2018, 11, 5, 0, 0, 0, 0, time.UTC), // it's a monday
			expected: false,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.s.Contains(tt.t),
			fmt.Sprintf("i=%v, t=%v", tt.s, tt.t))
	}
}
