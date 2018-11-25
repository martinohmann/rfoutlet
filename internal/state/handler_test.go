package state

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/stretchr/testify/assert"
)

func TestHandlerSaveState(t *testing.T) {
	defer os.Remove("testdata/state.json")

	h := NewHandler("testdata/state.json")

	outlets := []*outlet.Outlet{{ID: "foo", State: outlet.StateOn}}

	err := h.SaveState(outlets)

	assert.NoError(t, err)

	data, err := ioutil.ReadFile("testdata/state.json")

	assert.NoError(t, err)
	assert.Equal(t, "{\"switch_states\":{\"foo\":1},\"schedules\":{}}\n", string(data))
}

func TestHandlerLoad(t *testing.T) {
	h := NewHandler("testdata/valid.json")

	outlets := []*outlet.Outlet{{ID: "foo", State: outlet.StateOff}}

	err := h.LoadState(outlets)
	assert.NoError(t, err)

	assert.Equal(t, outlet.StateOn, outlets[0].GetState())
	assert.True(t, outlets[0].Schedule.Enabled())
}
