package command

import (
	"encoding/json"
	"fmt"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/internal/schedule"
)

// Type is the type of a command.
type Type string

// Command types.
const (
	GroupType    Type = "group"
	IntervalType Type = "interval"
	OutletType   Type = "outlet"
	StatusType   Type = "status"
)

// StatusCommand...
type StatusCommand struct {
	sender Sender
}

func (c StatusCommand) Execute(context Context) (bool, error) {
	groups := context.Registry.GetGroups()

	msg, err := json.Marshal(groups)
	if err != nil {
		return false, err
	}

	c.sender.Send(msg)

	return false, nil
}

func (c *StatusCommand) SetSender(sender Sender) {
	c.sender = sender
}

// OutletCommand...
type OutletCommand struct {
	OutletID string `json:"id"`
	Action   string `json:"action"`
}

func (c OutletCommand) Execute(context Context) (bool, error) {
	outlet, ok := context.Registry.GetOutlet(c.OutletID)
	if !ok {
		return false, fmt.Errorf("outlet %q does not exist", c.OutletID)
	}

	if outlet.Schedule.Enabled() {
		return false, nil
	}

	err := context.Switcher.Switch(outlet, getTargetState(outlet, c.Action))
	if err != nil {
		return false, err
	}

	return true, nil
}

// GroupCommand...
type GroupCommand struct {
	GroupID string `json:"id"`
	Action  string `json:"action"`
}

func (c GroupCommand) Execute(context Context) (bool, error) {
	group, ok := context.Registry.GetGroup(c.GroupID)
	if !ok {
		return false, fmt.Errorf("outlet group %q does not exist", c.GroupID)
	}

	var modified bool

	for _, outlet := range group.Outlets {
		if outlet.Schedule.Enabled() {
			continue
		}

		err := context.Switcher.Switch(outlet, getTargetState(outlet, c.Action))
		if err != nil {
			return modified, err
		}

		modified = true
	}

	return modified, nil
}

// IntervalCommand...
type IntervalCommand struct {
	OutletID string            `json:"id"`
	Action   string            `json:"action"`
	Interval schedule.Interval `json:"interval"`
}

func (c IntervalCommand) Execute(context Context) (bool, error) {
	outlet, ok := context.Registry.GetOutlet(c.OutletID)
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
	case "create":
		return outlet.Schedule.AddInterval(c.Interval)
	case "update":
		return outlet.Schedule.UpdateInterval(c.Interval)
	case "delete":
		return outlet.Schedule.DeleteInterval(c.Interval)
	default:
		return fmt.Errorf("invalid interval action %q", c.Action)
	}
}

type ScheduleCommand struct {
	Outlet       *outlet.Outlet
	DesiredState outlet.State
}

func (c ScheduleCommand) Execute(context Context) (bool, error) {
	if c.Outlet.GetState() == c.DesiredState {
		return false, nil
	}

	err := context.Switcher.Switch(c.Outlet, c.DesiredState)
	if err != nil {
		return false, err
	}

	return true, nil
}

func getTargetState(o *outlet.Outlet, action string) outlet.State {
	switch action {
	default:
		return outlet.StateOn
	case "off":
		return outlet.StateOff
	case "toggle":
		if o.GetState() == outlet.StateOn {
			return outlet.StateOff
		}

		return outlet.StateOn
	}
}
