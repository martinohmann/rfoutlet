package control

import (
	"fmt"

	"github.com/martinohmann/rfoutlet/internal/context"
	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/martinohmann/rfoutlet/internal/state"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
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
	ctx         *context.Context
	hub         *Hub
	transmitter gpio.CodeTransmitter
}

// New create a new control
func New(ctx *context.Context, transmitter gpio.CodeTransmitter) *Control {
	c := &Control{
		ctx:         ctx,
		hub:         NewHub(),
		transmitter: transmitter,
	}

	return c
}

// AddInterval adds a new schedule interval for an outlet
func (c *Control) AddInterval(o *context.Outlet, interval schedule.Interval) error {
	if err := o.AddInterval(interval); err != nil {
		return err
	}

	return c.SaveState()
}

// UpdateInterval updates an existing schedule interval for an outlet
func (c *Control) UpdateInterval(o *context.Outlet, interval schedule.Interval) error {
	if err := o.UpdateInterval(interval); err != nil {
		return err
	}

	return c.SaveState()
}

// DeleteInterval deletes a schedule interval for an outlet
func (c *Control) DeleteInterval(o *context.Outlet, interval schedule.Interval) error {
	if err := o.DeleteInterval(interval); err != nil {
		return err
	}

	return c.SaveState()
}

// SwitchState switches the outlet state
func (c *Control) SwitchState(o *context.Outlet, newState state.SwitchState) error {
	code := o.CodeOn
	if newState == state.SwitchStateOff {
		code = o.CodeOff
	}

	if err := c.transmitter.Transmit(code, o.Protocol, o.PulseLength); err != nil {
		return err
	}

	o.SetSwitchState(newState)

	return c.SaveState()
}

// Groups returns the configured outlet groups
func (c *Control) Groups() []*context.Group {
	return c.ctx.Groups
}

// Toggle switches outlet on if it is off, otherwise switches it on
func (c *Control) Toggle(o *context.Outlet) error {
	if o.State == state.SwitchStateOn {
		return c.SwitchState(o, state.SwitchStateOff)
	}

	return c.SwitchState(o, state.SwitchStateOn)
}

// SaveState saves the current state of all outlets
func (c *Control) SaveState() error {
	if c.ctx.Config.StateFile == "" {
		return nil
	}

	return state.Save(c.ctx.Config.StateFile, c.ctx.CollectState())
}

func (c *Control) handleOutletAction(o *context.Outlet, action string) error {
	if o.Schedule != nil && o.Schedule.Enabled() {
		return nil
	}

	switch action {
	case actionOn:
		return c.SwitchState(o, state.SwitchStateOn)
	case actionOff:
		return c.SwitchState(o, state.SwitchStateOff)
	case actionToggle:
		return c.Toggle(o)
	}

	return fmt.Errorf("invalid outlet action %q", action)
}

func (c *Control) handleIntervalAction(o *context.Outlet, i schedule.Interval, action string) error {
	switch action {
	case actionCreate:
		return c.AddInterval(o, i)
	case actionUpdate:
		return c.UpdateInterval(o, i)
	case actionDelete:
		return c.DeleteInterval(o, i)
	}

	return fmt.Errorf("invalid interval action %q", action)
}

// HandleMessage decode the message in the passed envelope and executes the
// command
func (c *Control) HandleMessage(message messageEnvelope) error {
	msg, err := decodeMessage(message)
	if err != nil {
		return err
	}

	switch msg.(type) {
	case outletMessage:
		data := msg.(outletMessage)

		o, err := c.ctx.GetOutlet(data.ID)
		if err != nil {
			return err
		}

		if err = c.handleOutletAction(o, data.Action); err != nil {
			return err
		}
	case groupMessage:
		data := msg.(groupMessage)

		og, err := c.ctx.GetGroup(data.ID)
		if err != nil {
			return err
		}

		for _, o := range og.Outlets {
			if err := c.handleOutletAction(o, data.Action); err != nil {
				return err
			}
		}
	case intervalMessage:
		data := msg.(intervalMessage)

		o, err := c.ctx.GetOutlet(data.ID)
		if err != nil {
			return err
		}

		if err := c.handleIntervalAction(o, data.Interval, data.Action); err != nil {
			return err
		}
	}

	b, err := encodeMessage(c.ctx.Groups)
	if err != nil {
		return err
	}

	c.Broadcast(b)

	return nil
}

func (c *Control) Broadcast(msg []byte) {
	c.hub.broadcast <- msg
}
