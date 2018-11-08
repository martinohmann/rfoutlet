package outlet_test

import (
	"testing"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/stretchr/testify/assert"
)

var control = outlet.NewControl(&outlet.Config{}, outlet.NewNullStateManager(), transmitter)

func TestControlSwitchOn(t *testing.T) {
	s := &outlet.Outlet{Protocol: 1, State: outlet.StateOff}

	assert.Nil(t, control.SwitchOn(s))
	assert.Equal(t, outlet.StateOn, s.State)
}

func TestControlSwitchOff(t *testing.T) {
	s := &outlet.Outlet{Protocol: 1, State: outlet.StateOn}

	assert.Nil(t, control.SwitchOff(s))
	assert.Equal(t, outlet.StateOff, s.State)
}

func TestControlToggleState(t *testing.T) {
	s := &outlet.Outlet{Protocol: 1, State: outlet.StateOn}

	assert.Nil(t, control.ToggleState(s))
	assert.Equal(t, outlet.StateOff, s.State)
}
