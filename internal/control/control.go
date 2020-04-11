package control

import (
	"encoding/json"
	"fmt"

	"github.com/martinohmann/rfoutlet/internal/message"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/internal/schedule"
)

const (
	actionOn     = "on"
	actionOff    = "off"
	actionToggle = "toggle"
	actionCreate = "create"
	actionUpdate = "update"
	actionDelete = "delete"
)

// Control type definition
type Control struct {
	registry *outlet.Registry
	switcher outlet.Switcher
	hub      *Hub
}

// New creates a new controller
func New(registry *outlet.Registry, switcher outlet.Switcher, hub *Hub) *Control {
	return &Control{registry, switcher, hub}
}

// Switch implements the outlet.Switcher interface
func (c *Control) Switch(o *outlet.Outlet, newState outlet.State) error {
	err := c.switcher.Switch(o, newState)
	if err != nil {
		return err
	}

	return c.broadcastState()
}

// Dispatch implements the messsage.Dispatcher interface
func (c *Control) Dispatch(envelope message.Envelope) error {
	msg, err := message.Decode(envelope)
	if err != nil {
		return err
	}

	switch data := msg.(type) {
	case *message.StatusMessage:
		return c.broadcastState()
	case *message.OutletMessage:
		o, ok := c.registry.GetOutlet(data.ID)
		if !ok {
			return fmt.Errorf("outlet %q does not exist", data.ID)
		}

		return c.outletAction(o, data.Action)
	case *message.GroupMessage:
		g, ok := c.registry.GetGroup(data.ID)
		if !ok {
			return fmt.Errorf("outlet group %q does not exist", data.ID)
		}

		for _, o := range g.Outlets {
			if err := c.outletAction(o, data.Action); err != nil {
				return err
			}
		}
	case *message.IntervalMessage:
		o, ok := c.registry.GetOutlet(data.ID)
		if !ok {
			return fmt.Errorf("outlet %q does not exist", data.ID)
		}

		return c.intervalAction(o, data.Interval, data.Action)
	}

	return nil
}

// broadcastState broadcasts the state of all outlet groups to all connected
// clients. the is called whenever switch states or outlet schedules are
// changed.
func (c *Control) broadcastState() error {
	b, err := json.Marshal(c.registry.GetGroups())
	if err != nil {
		return err
	}

	select {
	case c.hub.broadcast <- b:
	default:
	}

	return nil
}

// outletAction switches an outlet on or off, depending on the action provided.
// Outlets with enabled schedules will not be switched as they are managed by
// the scheduler.
func (c *Control) outletAction(o *outlet.Outlet, action string) error {
	if o.Schedule != nil && o.Schedule.Enabled() {
		return nil
	}

	switch action {
	case actionOn:
		return c.Switch(o, outlet.StateOn)
	case actionOff:
		return c.Switch(o, outlet.StateOff)
	case actionToggle:
		if o.State == outlet.StateOn {
			return c.Switch(o, outlet.StateOff)
		}

		return c.Switch(o, outlet.StateOn)
	}

	return fmt.Errorf("invalid outlet action %q", action)
}

// intervalAction
func (c *Control) intervalAction(o *outlet.Outlet, i schedule.Interval, action string) (err error) {
	switch action {
	case actionCreate:
		err = o.Schedule.AddInterval(i)
	case actionUpdate:
		err = o.Schedule.UpdateInterval(i)
	case actionDelete:
		err = o.Schedule.DeleteInterval(i)
	default:
		err = fmt.Errorf("invalid interval action %q", action)
	}

	if err != nil {
		return err
	}

	return c.broadcastState()
}
