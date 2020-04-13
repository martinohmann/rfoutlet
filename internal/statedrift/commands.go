package statedrift

import (
	"github.com/martinohmann/rfoutlet/internal/command"
	"github.com/martinohmann/rfoutlet/internal/outlet"
)

// StateCorrectionCommand is sent out whenever an outlet should change its state
// based on detected rf codes.
type StateCorrectionCommand struct {
	// Outlet is the outlet that should be brought into the desired state.
	Outlet *outlet.Outlet
	// DesiredState is the state that the outlet should be in.
	DesiredState outlet.State
}

// Execute implements command.Command.
//
// It switch an outlet to the detected state.
func (c StateCorrectionCommand) Execute(context command.Context) (bool, error) {
	// If the outlet was already switched to the desired state after we
	// submitted the command, we can bail out early.
	if c.Outlet.GetState() == c.DesiredState {
		return false, nil
	}

	err := context.Switch(c.Outlet, c.DesiredState)
	if err != nil {
		return false, err
	}

	return true, nil
}
