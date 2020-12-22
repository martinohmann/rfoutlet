package command

import (
	"bytes"
	"errors"
	"testing"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeSender struct {
	buf bytes.Buffer
}

func (s *fakeSender) Send(msg []byte) {
	s.buf.Write(msg)
}

func TestStatusCommand(t *testing.T) {
	ctx, r, _ := NewTestContext()

	s := &fakeSender{}

	r.RegisterGroups(&outlet.Group{ID: "foo"})

	cmd := StatusCommand{}
	cmd.SetSender(s)

	broadcast, err := cmd.Execute(ctx)

	require.NoError(t, err)
	assert.False(t, broadcast)

	assert.Equal(t, `[{"id":"foo","displayName":"","outlets":null}]`, s.buf.String())
}

func TestOutletCommand(t *testing.T) {
	ctx, r, _ := NewTestContext()

	o := &outlet.Outlet{ID: "foo"}

	r.RegisterOutlets(o)

	cmd := OutletCommand{"foo", "on"}

	broadcast, err := cmd.Execute(ctx)

	require.NoError(t, err)
	assert.True(t, broadcast)
	assert.Equal(t, outlet.StateOn, o.GetState())
}

func TestGroupCommand(t *testing.T) {
	ctx, r, _ := NewTestContext()

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
	ctx, r, _ := NewTestContext()

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

func TestStateCorrectionCommand(t *testing.T) {
	tests := []struct {
		name              string
		outletState       outlet.State
		desiredState      outlet.State
		switchErr         error
		expectedState     outlet.State
		expectedBroadcast bool
		expectedErr       error
	}{
		{
			name:          "outlet is already in desired state",
			outletState:   outlet.StateOn,
			desiredState:  outlet.StateOn,
			expectedState: outlet.StateOn,
		},
		{
			name:              "outlet is not in desired state",
			outletState:       outlet.StateOn,
			desiredState:      outlet.StateOff,
			expectedState:     outlet.StateOff,
			expectedBroadcast: true,
		},
		{
			name:              "outlet is not in desired state, switch error",
			outletState:       outlet.StateOn,
			desiredState:      outlet.StateOff,
			switchErr:         errors.New("whoops"),
			expectedState:     outlet.StateOn,
			expectedBroadcast: false,
			expectedErr:       errors.New("whoops"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o := &outlet.Outlet{State: test.outletState}
			cmd := StateCorrectionCommand{
				Outlet:       o,
				DesiredState: test.desiredState,
			}

			ctx, _, s := NewTestContext()
			s.Err = test.switchErr

			broadcast, err := cmd.Execute(ctx)

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedBroadcast, broadcast)
			assert.Equal(t, test.expectedState, o.GetState())
		})
	}
}
