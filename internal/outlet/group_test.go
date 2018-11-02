package outlet_test

import (
	"testing"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/stretchr/testify/assert"
)

func TestNewOutletGroup(t *testing.T) {
	og := outlet.NewOutletGroup("foo")

	assert.Equal(t, "foo", og.Identifier)
}

func TestAddOutlet(t *testing.T) {
	og := outlet.NewOutletGroup("foo")
	o := &outlet.Outlet{}

	og.AddOutlet(o)

	assert.Len(t, og.Outlets, 1)
	assert.Equal(t, o, og.Outlets[0])
}

func TestOutlet(t *testing.T) {
	o := &outlet.Outlet{}

	og := &outlet.OutletGroup{
		Outlets: []*outlet.Outlet{o},
	}

	res, err := og.Outlet(0)

	assert.Nil(t, err)
	assert.Equal(t, o, res)

	res, err = og.Outlet(1)

	assert.Nil(t, res)
	assert.EqualError(t, err, "invalid offset 1")
}

func TestOutputGroupSwitchOn(t *testing.T) {
	o := &outlet.Outlet{State: outlet.StateOff}

	og := &outlet.OutletGroup{
		Outlets: []*outlet.Outlet{o},
	}

	err := og.SwitchOn(transmitter)

	assert.Nil(t, err)
	assert.Equal(t, outlet.StateOn, o.State)
}

func TestOutputGroupSwitchOff(t *testing.T) {
	o := &outlet.Outlet{State: outlet.StateOn}

	og := &outlet.OutletGroup{
		Outlets: []*outlet.Outlet{o},
	}

	err := og.SwitchOff(transmitter)

	assert.Nil(t, err)
	assert.Equal(t, outlet.StateOff, o.State)
}

func TestOutputGroupToggleState(t *testing.T) {
	o := &outlet.Outlet{State: outlet.StateOff}

	og := &outlet.OutletGroup{
		Outlets: []*outlet.Outlet{o},
	}

	err := og.ToggleState(transmitter)

	assert.Nil(t, err)
	assert.Equal(t, outlet.StateOn, o.State)

	err = og.ToggleState(transmitter)

	assert.Nil(t, err)
	assert.Equal(t, outlet.StateOff, o.State)
}
