package websocket

import (
	"testing"

	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutletCommand(t *testing.T) {
	ctx, r, _ := command.NewTestContext()

	o := &outlet.Outlet{ID: "foo"}

	r.RegisterOutlets(o)

	cmd := OutletCommand{"foo", "on"}

	broadcast, err := cmd.Execute(ctx)

	require.NoError(t, err)
	assert.True(t, broadcast)
	assert.Equal(t, outlet.StateOn, o.GetState())
}

func TestGroupCommand(t *testing.T) {
	ctx, r, _ := command.NewTestContext()

	o1 := &outlet.Outlet{ID: "foo", State: outlet.StateOn}
	o2 := &outlet.Outlet{ID: "baz"}
	o3 := &outlet.Outlet{
		ID: "qux",
		Schedule: schedule.NewWithIntervals([]schedule.Interval{
			{
				Enabled: true,
			},
		}),
	}

	r.RegisterGroups(&outlet.Group{
		ID:      "bar",
		Outlets: []*outlet.Outlet{o1, o2, o3},
	})

	cmd := GroupCommand{"bar", "toggle"}

	broadcast, err := cmd.Execute(ctx)

	require.NoError(t, err)
	assert.True(t, broadcast)
	assert.Equal(t, outlet.StateOff, o1.GetState())
	assert.Equal(t, outlet.StateOn, o2.GetState())
	assert.Equal(t, outlet.StateOff, o3.GetState())
}

func TestIntervalCommand(t *testing.T) {
	ctx, r, _ := command.NewTestContext()

	o := &outlet.Outlet{ID: "foo", Schedule: schedule.New()}

	r.RegisterOutlets(o)

	cmd := IntervalCommand{"foo", "create", schedule.Interval{ID: "bar"}}

	broadcast, err := cmd.Execute(ctx)

	require.NoError(t, err)
	assert.True(t, broadcast)

	cmd = IntervalCommand{"foo", "create", schedule.Interval{ID: "bar"}}

	broadcast, err = cmd.Execute(ctx)

	require.Error(t, err)
	assert.False(t, broadcast)

	cmd = IntervalCommand{"foo", "update", schedule.Interval{ID: "bar"}}

	broadcast, err = cmd.Execute(ctx)

	require.NoError(t, err)
	assert.True(t, broadcast)

	cmd = IntervalCommand{"foo", "delete", schedule.Interval{ID: "bar"}}

	broadcast, err = cmd.Execute(ctx)

	require.NoError(t, err)
	assert.True(t, broadcast)
}
