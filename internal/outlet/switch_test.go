package outlet

import (
	"errors"
	"testing"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/stretchr/testify/assert"
)

func TestSwitch(t *testing.T) {
	s := NewSwitch(gpio.NewDiscardingTransmitter())
	o := &Outlet{State: StateOn, Protocol: 1}

	assert.NoError(t, s.Switch(o, StateOff))
	assert.Equal(t, StateOff, o.GetState())

	assert.NoError(t, s.Switch(o, StateOn))
	assert.Equal(t, StateOn, o.GetState())
}

func TestFakeSwitch(t *testing.T) {
	s := &FakeSwitch{}
	o := &Outlet{State: StateOn}

	assert.NoError(t, s.Switch(o, StateOff))
	assert.Equal(t, StateOff, o.GetState())

	s = &FakeSwitch{Err: errors.New("whoops")}

	assert.Error(t, s.Switch(o, StateOn))
	assert.Equal(t, StateOff, o.GetState())
}
