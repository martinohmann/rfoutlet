package timeswitch

import (
	"errors"
	"testing"

	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/stretchr/testify/assert"
)

func TestTimeSwitchCommand(t *testing.T) {
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
			cmd := TimeSwitchCommand{
				Outlet:       o,
				DesiredState: test.desiredState,
			}

			ctx, _, s := command.NewTestContext()
			s.Err = test.switchErr

			broadcast, err := cmd.Execute(ctx)

			assert.Equal(t, test.expectedErr, err)
			assert.Equal(t, test.expectedBroadcast, broadcast)
			assert.Equal(t, test.expectedState, o.GetState())
		})
	}
}
