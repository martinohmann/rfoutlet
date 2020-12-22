package command

import (
	"encoding/json"
	"fmt"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/internal/schedule"
)

const (
	outletOnAction     = "on"
	outletOffAction    = "off"
	outletToggleAction = "toggle"

	intervalCreateAction = "create"
	intervalUpdateAction = "update"
	intervalDeleteAction = "delete"
)

// StatusCommand is sent by a connected client to retrieve the list of current
// outlet groups. This usually happens when the client first connects.
type StatusCommand struct {
	sender Sender
}

// Execute implements Command.
//
// It sends the registered outlet groups back to the sender.
func (c StatusCommand) Execute(context Context) (bool, error) {
	groups := context.GetGroups()

	msg, err := json.Marshal(groups)
	if err != nil {
		return false, err
	}

	c.sender.Send(msg)

	return false, nil
}

// SetSender implements SenderAwareCommand.
func (c *StatusCommand) SetSender(sender Sender) {
	c.sender = sender
}

// OutletCommand switches a specific outlet based on the action.
type OutletCommand struct {
	// OutletID is the ID of the outlet that the action should be performed on.
	OutletID string `json:"outletID"`
	// Action defines the action type that should be performed on the outlet.
	Action string `json:"action"`
}

// Execute implements Command.
//
// It switches an outlet based on the transmitted action.
func (c OutletCommand) Execute(context Context) (bool, error) {
	outlet, ok := context.GetOutlet(c.OutletID)
	if !ok {
		return false, fmt.Errorf("outlet %q does not exist", c.OutletID)
	}

	if outlet.Schedule.Enabled() {
		return false, nil
	}

	targetState, err := getTargetState(outlet, c.Action)
	if err != nil {
		return false, err
	}

	err = context.Switch(outlet, targetState)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GroupCommand switches a group of outlets based on the action.
type GroupCommand struct {
	// GroupID is the ID of the outlet group that the action should be
	// performed on.
	GroupID string `json:"groupID"`
	// Action defines the action type that should be performed on the outlet
	// group.
	Action string `json:"action"`
}

// Execute implements Command.
//
// It switches a group of outlets based on the transmitted action.
func (c GroupCommand) Execute(context Context) (bool, error) {
	group, ok := context.GetGroup(c.GroupID)
	if !ok {
		return false, fmt.Errorf("outlet group %q does not exist", c.GroupID)
	}

	var modified bool

	for _, outlet := range group.Outlets {
		if outlet.Schedule.Enabled() {
			continue
		}

		targetState, err := getTargetState(outlet, c.Action)
		if err != nil {
			return modified, err
		}

		err = context.Switch(outlet, targetState)
		if err != nil {
			return modified, err
		}

		modified = true
	}

	return modified, nil
}

// IntervalCommand changes the intervals of an outlet based on the action.
type IntervalCommand struct {
	// OutletID is the ID of the outlet where the intervals of the schedule
	// should be changed.
	OutletID string `json:"outletID"`
	// Action defines the action type that should be performed on the interval.
	Action string `json:"action"`
	// Interval is the configuration of the interval.
	Interval schedule.Interval `json:"interval"`
}

// Execute implements Command.
//
// It add, updates or deletes an interval from the outlet's schedule based on
// the transmitted action.
func (c IntervalCommand) Execute(context Context) (bool, error) {
	outlet, ok := context.GetOutlet(c.OutletID)
	if !ok {
		return false, fmt.Errorf("outlet %q does not exist", c.OutletID)
	}

	err := c.handle(outlet)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c IntervalCommand) handle(outlet *outlet.Outlet) error {
	switch c.Action {
	case intervalCreateAction:
		return outlet.Schedule.AddInterval(c.Interval)
	case intervalUpdateAction:
		return outlet.Schedule.UpdateInterval(c.Interval)
	case intervalDeleteAction:
		return outlet.Schedule.DeleteInterval(c.Interval)
	default:
		return fmt.Errorf("invalid interval action %q", c.Action)
	}
}

func getTargetState(o *outlet.Outlet, action string) (outlet.State, error) {
	switch action {
	case outletOnAction:
		return outlet.StateOn, nil
	case outletOffAction:
		return outlet.StateOff, nil
	case outletToggleAction:
		if o.GetState() == outlet.StateOn {
			return outlet.StateOff, nil
		}

		return outlet.StateOn, nil
	default:
		return 0, fmt.Errorf("invalid outlet action %q", action)
	}
}

// StateCorrectionCommand is sent out whenever an outlet should change its state
// based on detected rf codes.
type StateCorrectionCommand struct {
	// Outlet is the outlet that should be brought into the desired state.
	Outlet *outlet.Outlet
	// DesiredState is the state that the outlet should be in.
	DesiredState outlet.State
}

// Execute implements Command.
//
// It switch an outlet to the detected state.
func (c StateCorrectionCommand) Execute(context Context) (bool, error) {
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
