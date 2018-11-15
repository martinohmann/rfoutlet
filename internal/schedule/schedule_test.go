package schedule_test

import (
	"testing"
	"time"

	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/stretchr/testify/assert"
)

func TestScheduleEnabled(t *testing.T) {
	s := schedule.Schedule{{Enabled: true}}

	assert.True(t, s.Enabled())
}

func TestScheduleContains(t *testing.T) {
	s := schedule.Schedule{{Enabled: true}}

	assert.False(t, s.Contains(time.Date(2018, 12, 12, 12, 34, 56, 0, time.UTC)))
}
