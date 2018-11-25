package schedule

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduleEnabled(t *testing.T) {
	tests := []struct {
		is       []Interval
		expected bool
	}{
		{
			is: []Interval{
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
			is: []Interval{
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
			is: []Interval{
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
			is: []Interval{
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
		s := NewWithIntervals(tt.is)
		assert.Equal(t, tt.expected, s.Enabled(), fmt.Sprintf("is=%v", tt.is))
	}
}

func TestScheduleContains(t *testing.T) {
	tests := []struct {
		is       []Interval
		t        time.Time
		expected bool
	}{
		{
			is: []Interval{
				{
					Enabled:  true,
					Weekdays: []time.Weekday{time.Monday},
					From:     NewDayTime(0, 0),
					To:       NewDayTime(3, 0),
				},
			},
			t:        time.Date(2018, 11, 5, 0, 0, 0, 0, time.UTC), // it's a monday
			expected: true,
		},
		{
			is: []Interval{
				{
					Enabled:  true,
					Weekdays: []time.Weekday{time.Tuesday},
					From:     NewDayTime(0, 0),
					To:       NewDayTime(3, 0),
				},
			},
			t:        time.Date(2018, 11, 5, 0, 0, 0, 0, time.UTC), // it's a monday
			expected: false,
		},
		{
			is: []Interval{
				{
					Enabled:  false,
					Weekdays: []time.Weekday{time.Monday},
					From:     NewDayTime(0, 0),
					To:       NewDayTime(3, 0),
				},
			},
			t:        time.Date(2018, 11, 5, 0, 0, 0, 0, time.UTC), // it's a monday
			expected: false,
		},
	}

	for _, tt := range tests {
		s := NewWithIntervals(tt.is)
		assert.Equal(t, tt.expected, s.Contains(tt.t), fmt.Sprintf("is=%v, t=%v", tt.is, tt.t))
	}
}

func TestAddInterval(t *testing.T) {
	s := New()
	i := Interval{}

	err := s.AddInterval(i)

	assert.NoError(t, err)
	assert.Len(t, s.intervals, 1)
}

func TestDeleteInterval(t *testing.T) {
	i := Interval{ID: "foo"}
	s := NewWithIntervals([]Interval{i})

	err := s.DeleteInterval(i)

	assert.NoError(t, err)
	assert.Len(t, s.intervals, 0)
}

func TestUpdateInterval(t *testing.T) {
	i := Interval{ID: "foo", Enabled: false}
	s := NewWithIntervals([]Interval{i})

	i2 := Interval{ID: "foo", Enabled: true}

	err := s.UpdateInterval(i2)

	assert.NoError(t, err)
	assert.Equal(t, true, s.intervals[0].Enabled)
}
