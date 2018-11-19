package context_test

import (
	"testing"

	"github.com/martinohmann/rfoutlet/internal/context"
	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/stretchr/testify/assert"
)

func TestAddInterval(t *testing.T) {
	o := &context.Outlet{Schedule: schedule.Schedule{}}

	i := schedule.Interval{ID: "foo"}

	assert.NoError(t, o.AddInterval(i))
	assert.Error(t, o.AddInterval(i))
}

func TestDeleteInterval(t *testing.T) {

	o := &context.Outlet{Schedule: schedule.Schedule{{ID: "foo"}}}

	i := schedule.Interval{ID: "foo"}
	i2 := schedule.Interval{ID: "bar"}

	assert.NoError(t, o.DeleteInterval(i))
	assert.Error(t, o.DeleteInterval(i2))

	assert.Len(t, o.Schedule, 0)
}

func TestUpdateInterval(t *testing.T) {
	o := &context.Outlet{Schedule: schedule.Schedule{{ID: "foo"}}}

	i := schedule.Interval{ID: "foo", Enabled: true}
	i2 := schedule.Interval{ID: "bar"}

	assert.NoError(t, o.UpdateInterval(i))
	assert.Error(t, o.UpdateInterval(i2))
	assert.True(t, o.Schedule[0].Enabled)
}
