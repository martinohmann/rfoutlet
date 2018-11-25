package outlet

import (
	"testing"

	"github.com/martinohmann/rfoutlet/pkg/gpio"
	"github.com/stretchr/testify/assert"
)

func TestSwitch(t *testing.T) {
	trans, _ := gpio.NewNullTransmitter()

	s := NewSwitch(trans)
	o := &Outlet{State: StateOn, Protocol: 1}

	assert.NoError(t, s.Switch(o, StateOff))
	assert.Equal(t, StateOff, o.GetState())

	assert.NoError(t, s.Switch(o, StateOn))
	assert.Equal(t, StateOn, o.GetState())
}
