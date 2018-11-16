package control_test

import (
	"testing"

	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/martinohmann/rfoutlet/internal/context"
	"github.com/martinohmann/rfoutlet/internal/control"
	"github.com/martinohmann/rfoutlet/internal/state"
	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/stretchr/testify/assert"
)

func createContext() *context.Context {
	c := &config.Config{
		GroupOrder: []string{"foo"},
		Groups: map[string]*config.Group{
			"foo": {
				Name:    "Foo",
				Outlets: []string{"bar", "baz", "qux"},
			},
		},
		Outlets: map[string]*config.Outlet{
			"bar": {Name: "Bar", Protocol: 1},
			"baz": {Name: "Baz", Protocol: 1},
			"qux": {Name: "Qux", Protocol: 1},
		},
	}

	s := state.New()
	s.SwitchStates["qux"] = state.SwitchStateOn

	ctx, _ := context.New(c, s)

	return ctx
}

func createControl(ctx *context.Context) *control.Control {
	t, _ := gpio.NewNullTransmitter()

	return control.New(ctx, t)
}

func TestSwitchOn(t *testing.T) {
	ctx := createContext()
	c := createControl(ctx)
	o := ctx.Groups[0].Outlets[0]

	err := c.SwitchOn(o)

	if assert.NoError(t, err) {
		assert.Equal(t, state.SwitchStateOn, o.State)
	}
}

func TestSwitchOff(t *testing.T) {
	ctx := createContext()
	c := createControl(ctx)
	o := ctx.Groups[0].Outlets[2]

	err := c.SwitchOff(o)

	if assert.NoError(t, err) {
		assert.Equal(t, state.SwitchStateOff, o.State)
	}
}

func TestToggle(t *testing.T) {
	ctx := createContext()
	c := createControl(ctx)
	o := ctx.Groups[0].Outlets[2]

	err := c.Toggle(o)

	if assert.NoError(t, err) {
		assert.Equal(t, state.SwitchStateOff, o.State)
	}
}
