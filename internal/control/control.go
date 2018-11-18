package control

import (
	"fmt"

	"github.com/martinohmann/rfoutlet/internal/context"
	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/martinohmann/rfoutlet/internal/state"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	uuid "github.com/satori/go.uuid"
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
	if interval.ID == "" {
		interval.ID = uuid.NewV4().String()
	}

	o.Schedule = append(o.Schedule, interval)

	c.ctx.State.Schedules[o.ID] = o.Schedule

	return c.SaveState()
}

// UpdateInterval updates an existing schedule interval for an outlet
func (c *Control) UpdateInterval(o *context.Outlet, interval schedule.Interval) error {
	for i, intv := range o.Schedule {
		if intv.ID != interval.ID {
			continue
		}

		o.Schedule[i] = interval

		c.ctx.State.Schedules[o.ID] = o.Schedule

		return c.SaveState()
	}

	return fmt.Errorf("interval with identifier %q does not exist", interval.ID)
}

// DeleteInterval deletes a schedule interval for an outlet
func (c *Control) DeleteInterval(o *context.Outlet, intervalID string) error {
	for i, interval := range o.Schedule {
		if interval.ID == intervalID {
			o.Schedule = append(o.Schedule[:i], o.Schedule[i+1:]...)
			break
		}
	}

	c.ctx.State.Schedules[o.ID] = o.Schedule

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

	o.State = newState
	c.ctx.State.SwitchStates[o.ID] = newState

	return c.SaveState()
}

// Toggle switches outlet on if it is off, otherwise switches it on
func (c *Control) Toggle(o *context.Outlet) error {
	if o.State == state.SwitchStateOn {
		return c.SwitchState(o, state.SwitchStateOff)
	}

	return c.SwitchState(o, state.SwitchStateOn)
}

func (c *Control) SaveState() error {
	if c.ctx.Config.StateFile == "" {
		return nil
	}

	return state.Save(c.ctx.Config.StateFile, c.ctx.State)
}
