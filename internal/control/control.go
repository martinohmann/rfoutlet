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

// SwitchOn switches an outlet on
func (c *Control) SwitchOn(outlet *context.Outlet) error {
	return c.switchState(state.SwitchStateOn, outlet)
}

// SwitchOff switches an outlet off
func (c *Control) SwitchOff(outlet *context.Outlet) error {
	return c.switchState(state.SwitchStateOff, outlet)
}

// Toggle switches outlet on if it is off, otherwise switches it on
func (c *Control) Toggle(outlet *context.Outlet) error {
	if outlet.State == state.SwitchStateOn {
		return c.SwitchOff(outlet)
	}

	return c.SwitchOn(outlet)
}

func (c *Control) switchState(newState state.SwitchState, o *context.Outlet) error {
	code := o.CodeOn
	if newState == state.SwitchStateOff {
		code = o.CodeOff
	}

	if err := c.transmitter.Transmit(code, o.Protocol, o.PulseLength); err != nil {
		return err
	}

	o.State = newState
	c.ctx.State.SwitchStates[o.ID] = newState

	return c.saveState()
}

func (c *Control) saveState() error {
	if c.ctx.Config.StateFile == "" {
		return nil
	}

	return state.Save(c.ctx.Config.StateFile, c.ctx.State)
}
