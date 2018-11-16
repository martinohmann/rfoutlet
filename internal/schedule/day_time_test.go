package schedule_test

import (
	"fmt"
	"testing"

	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/stretchr/testify/assert"
)

func TestDayTimeBefore(t *testing.T) {
	tests := []struct {
		a        schedule.DayTime
		b        schedule.DayTime
		expected bool
	}{
		{
			a:        schedule.NewDayTime(10, 0),
			b:        schedule.NewDayTime(12, 0),
			expected: true,
		},
		{
			a:        schedule.NewDayTime(10, 0),
			b:        schedule.NewDayTime(10, 1),
			expected: true,
		},
		{
			a:        schedule.NewDayTime(10, 0),
			b:        schedule.NewDayTime(10, 0),
			expected: false,
		},
		{
			a:        schedule.NewDayTime(10, 0),
			b:        schedule.NewDayTime(9, 59),
			expected: false,
		},
		{
			a:        schedule.NewDayTime(0, 0),
			b:        schedule.NewDayTime(23, 59),
			expected: true,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.a.Before(tt.b),
			fmt.Sprintf("a=%v, b=%v", tt.a, tt.b))
	}
}

func TestDayTimeAfter(t *testing.T) {
	tests := []struct {
		a        schedule.DayTime
		b        schedule.DayTime
		expected bool
	}{
		{
			a:        schedule.NewDayTime(10, 0),
			b:        schedule.NewDayTime(12, 0),
			expected: false,
		},
		{
			a:        schedule.NewDayTime(10, 0),
			b:        schedule.NewDayTime(10, 1),
			expected: false,
		},
		{
			a:        schedule.NewDayTime(10, 0),
			b:        schedule.NewDayTime(10, 0),
			expected: false,
		},
		{
			a:        schedule.NewDayTime(10, 0),
			b:        schedule.NewDayTime(9, 59),
			expected: true,
		},
		{
			a:        schedule.NewDayTime(23, 59),
			b:        schedule.NewDayTime(0, 0),
			expected: true,
		},
		{
			a:        schedule.NewDayTime(0, 0),
			b:        schedule.NewDayTime(23, 59),
			expected: false,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.a.After(tt.b),
			fmt.Sprintf("a=%v, b=%v", tt.a, tt.b))
	}
}
