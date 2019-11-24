package state

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/stretchr/testify/assert"
)

func TestLoadInvalid(t *testing.T) {
	_, err := Load("testdata/invalid.json")
	assert.Error(t, err)
}

func TestLoadNonexistent(t *testing.T) {
	_, err := Load("testdata/idonotexist.json")
	assert.Error(t, err)
}

func TestLoadWithReader(t *testing.T) {
	r := strings.NewReader(`{"switch_states":{"foo":1}}`)
	s, err := LoadWithReader(r)
	assert.NoError(t, err)
	assert.Equal(t, outlet.StateOn, s.SwitchStates["foo"])
}

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 1, fmt.Errorf("error")
}

func TestLoadWithBadReader(t *testing.T) {
	_, err := LoadWithReader(errorReader{})
	assert.Error(t, err)
}

func TestSaveInNonexistentDir(t *testing.T) {
	s := New()
	assert.Error(t, Save("testdata/thisdoesnotexist/json", s))
}

type errorWriter struct{}

func (errorWriter) Write(p []byte) (n int, err error) {
	return 1, fmt.Errorf("error")
}

func TestSaveWithWriter(t *testing.T) {
	var buf []byte

	w := bytes.NewBuffer(buf)

	assert.NoError(t, SaveWithWriter(w, New()))
	assert.Equal(t, "{\"switch_states\":{},\"schedules\":{}}\n", w.String())
}

func TestSaveWithBadWriter(t *testing.T) {
	assert.Error(t, SaveWithWriter(errorWriter{}, New()))
}

func TestCollect(t *testing.T) {
	schedule := schedule.New()
	outlets := []*outlet.Outlet{{ID: "foo", State: outlet.StateOn, Schedule: schedule}}

	s := Collect(outlets)

	if assert.Len(t, s.Schedules, 1) {
		assert.Equal(t, s.Schedules["foo"], schedule)
	}

	if assert.Len(t, s.SwitchStates, 1) {
		assert.Equal(t, s.SwitchStates["foo"], outlet.StateOn)
	}
}

func TestApply(t *testing.T) {
	outlets := []*outlet.Outlet{{ID: "foo", State: outlet.StateOn}}
	sch := schedule.New()

	s := &State{
		SwitchStates: map[string]outlet.State{"foo": outlet.StateOff},
		Schedules:    map[string]*schedule.Schedule{"foo": sch},
	}

	s.Apply(outlets)

	assert.Equal(t, outlet.StateOff, outlets[0].State)
	assert.Equal(t, sch, outlets[0].Schedule)
}
