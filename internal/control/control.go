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
	manager  *outlet.Manager
	switcher outlet.Switcher
	hub      *Hub
}

// New creates a new controller
func New(manager *outlet.Manager, switcher outlet.Switcher, hub *Hub) *Control {
	return &Control{manager, switcher, hub}
}

// Switch implements the outlet.Switcher interface
func (c *Control) Switch(o *outlet.Outlet, newState outlet.State) error {
	err := c.switcher.Switch(o, newState)
	if err != nil {
		return err
	}

	if err = c.broadcastState(); err != nil {
		return err
	}

	return c.manager.SaveState()
}

// Dispatch implements the messsage.Dispatcher interface
func (c *Control) Dispatch(env message.Envelope) error {
	msg, err := message.Decode(env)
	if err != nil {
		return err
	}

	switch msg.(type) {
	case message.Unknown:
		return c.broadcastState()
	case message.OutletAction:
		data := msg.(message.OutletAction)

		o, err := c.manager.Get(data.ID)
		if err != nil {
			return err
		}

		if err = handleOutletAction(o, c, data.Action); err != nil {
			return err
		}
	case message.GroupAction:
		data := msg.(message.GroupAction)

		og, err := c.manager.GetGroup(data.ID)
		if err != nil {
			return err
		}

		for _, o := range og.Outlets {
			if err := handleOutletAction(o, c, data.Action); err != nil {
				return err
			}
		}
	case message.IntervalAction:
		data := msg.(message.IntervalAction)

		o, err := c.manager.Get(data.ID)
		if err != nil {
			return err
		}

		if err = handleIntervalAction(o, data.Interval, data.Action); err != nil {
			return err
		}

		return c.broadcastState()
	}

	return nil
}

func (c *Control) broadcastState() error {
	if b, err := json.Marshal(c.manager.Groups()); err != nil {
		return err
	} else {
		c.hub.broadcast <- b
	}

	return nil
}

func handleOutletAction(o *outlet.Outlet, s outlet.Switcher, action string) error {
	if o.Schedule.Enabled() {
		return nil
	}

	switch action {
	case actionOn:
		return s.Switch(o, outlet.StateOn)
	case actionOff:
		return s.Switch(o, outlet.StateOff)
	case actionToggle:
		if o.State == outlet.StateOn {
			return s.Switch(o, outlet.StateOff)
		}

		return s.Switch(o, outlet.StateOn)
	}

	return fmt.Errorf("invalid outlet action %q", action)
}

func handleIntervalAction(o *outlet.Outlet, i schedule.Interval, action string) (err error) {
	switch action {
	case actionCreate:
		return o.Schedule.AddInterval(i)
	case actionUpdate:
		return o.Schedule.UpdateInterval(i)
	case actionDelete:
		return o.Schedule.DeleteInterval(i)
	}

	return fmt.Errorf("invalid interval action %q", action)
}
