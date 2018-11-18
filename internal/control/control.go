package control

import (
	"time"

	"github.com/martinohmann/rfoutlet/internal/context"
	"github.com/martinohmann/rfoutlet/internal/state"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
)

// Control type definition
type Control struct {
	ctx         *context.Context
	scheduler   *Scheduler
	transmitter gpio.CodeTransmitter
}

// New create a new control
func New(ctx *context.Context, transmitter gpio.CodeTransmitter) *Control {
	c := &Control{
		ctx:         ctx,
		scheduler:   NewScheduler(ctx, 10*time.Second),
		transmitter: transmitter,
	}

	c.scheduler.Start()

	return c
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
