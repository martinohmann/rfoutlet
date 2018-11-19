package control

import (
	"github.com/martinohmann/rfoutlet/internal/context"
	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/martinohmann/rfoutlet/internal/state"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
)

// Control type definition
type Control struct {
	ctx         *context.Context
	transmitter gpio.CodeTransmitter
}

// New create a new control
func New(ctx *context.Context, transmitter gpio.CodeTransmitter) *Control {
	c := &Control{
		ctx:         ctx,
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
