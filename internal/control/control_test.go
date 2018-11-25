package control

import (
	"encoding/json"
	"testing"

	"github.com/martinohmann/rfoutlet/internal/message"
	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/stretchr/testify/assert"
)

var (
	testStateHandler = new(nopStateHandler)
	testSwitcher     = new(nopSwitcher)
)

type nopStateHandler int

func (nopStateHandler) LoadState(o []*outlet.Outlet) error {
	return nil
}

func (nopStateHandler) SaveState(o []*outlet.Outlet) error {
	return nil
}

type nopSwitcher int

func (nopSwitcher) Switch(o *outlet.Outlet, s outlet.State) error {
	return nil
}

func createControl(m *outlet.Manager) *Control {
	return New(m, testSwitcher, NewHub())
}

func TestSwitch(t *testing.T) {
	m := outlet.NewManager(testStateHandler)
	c := createControl(m)
	o := &outlet.Outlet{ID: "foo", Protocol: 1}

	assert.NoError(t, c.Switch(o, outlet.StateOn))
}

func TestDispatch(t *testing.T) {
	m := outlet.NewManager(testStateHandler)
	c := createControl(m)
	o := &outlet.Outlet{ID: "foo", Protocol: 1}
	m.Register("foo", o)

	data := json.RawMessage([]byte(`{"id":"foo","action":"on"}`))
	env := message.Envelope{
		Type: message.OutletActionType,
		Data: &data,
	}

	assert.NoError(t, c.Dispatch(env))
}

func TestDispatchError(t *testing.T) {
	m := outlet.NewManager(testStateHandler)
	c := createControl(m)

	env := message.Envelope{
		Type: message.OutletActionType,
		Data: &json.RawMessage{},
	}

	assert.Error(t, c.Dispatch(env))
}
