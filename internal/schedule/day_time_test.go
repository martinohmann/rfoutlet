package schedule_test

import (
	"testing"

	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/stretchr/testify/assert"
)

func TestBefore(t *testing.T) {
	tests := []struct {
		a, b     schedule.DayTime
		expected bool
	}{
		{
			a:        schedule.NewDayTime(0, 0),
			b:        schedule.NewDayTime(0, 0),
			expected: false,
		},
		{
			a:        schedule.NewDayTime(0, 0),
			b:        schedule.NewDayTime(0, 1),
			expected: true,
		},
		{
			a:        schedule.NewDayTime(0, 1),
			b:        schedule.NewDayTime(0, 0),
			expected: false,
		},
		{
			a:        schedule.NewDayTime(1, 0),
			b:        schedule.NewDayTime(0, 0),
			expected: false,
		},
		{
			a:        schedule.NewDayTime(0, 59),
			b:        schedule.NewDayTime(1, 0),
			expected: true,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.a.Before(tt.b))
	}
}

func TestAfter(t *testing.T) {
	tests := []struct {
		a, b     schedule.DayTime
		expected bool
	}{
		{
			a:        schedule.NewDayTime(0, 0),
			b:        schedule.NewDayTime(0, 0),
			expected: false,
		},
		{
			a:        schedule.NewDayTime(0, 0),
			b:        schedule.NewDayTime(0, 1),
			expected: false,
		},
		{
			a:        schedule.NewDayTime(0, 1),
			b:        schedule.NewDayTime(0, 0),
			expected: true,
		},
		{
			a:        schedule.NewDayTime(1, 0),
			b:        schedule.NewDayTime(0, 0),
			expected: true,
		},
		{
			a:        schedule.NewDayTime(0, 59),
			b:        schedule.NewDayTime(1, 0),
			expected: false,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.a.After(tt.b))
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		a, b     schedule.DayTime
		expected bool
	}{
		{
			a:        schedule.NewDayTime(0, 0),
			b:        schedule.NewDayTime(0, 0),
			expected: true,
		},
		{
			a:        schedule.NewDayTime(0, 0),
			b:        schedule.NewDayTime(0, 1),
			expected: false,
		},
		{
			a:        schedule.NewDayTime(0, 1),
			b:        schedule.NewDayTime(0, 0),
			expected: false,
		},
		{
			a:        schedule.NewDayTime(1, 0),
			b:        schedule.NewDayTime(0, 0),
			expected: false,
		},
		{
			a:        schedule.NewDayTime(0, 59),
			b:        schedule.NewDayTime(1, 0),
			expected: false,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.a.Equal(tt.b))
	}
}
