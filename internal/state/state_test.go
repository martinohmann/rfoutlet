package state_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/martinohmann/rfoutlet/internal/state"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	s, err := state.Load("testdata/valid.json")
	assert.NoError(t, err)

	assert.Equal(t, state.SwitchStateOn, s.SwitchStates["foo"])
	assert.True(t, s.Schedules["foo"].Enabled())
}

func TestLoadInvalid(t *testing.T) {
	_, err := state.Load("testdata/invalid.json")
	assert.Error(t, err)
}

func TestLoadNonexistent(t *testing.T) {
	_, err := state.Load("testdata/idonotexist.json")
	assert.Error(t, err)
}

func TestLoadWithReader(t *testing.T) {
	r := strings.NewReader(`{"switch_states":{"foo":1}}`)
	s, err := state.LoadWithReader(r)
	assert.NoError(t, err)
	assert.Equal(t, state.SwitchStateOn, s.SwitchStates["foo"])
}

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 1, fmt.Errorf("error")
}

func TestLoadWithBadReader(t *testing.T) {
	_, err := config.LoadWithReader(errorReader{})
	assert.Error(t, err)
}

func TestSave(t *testing.T) {
	defer os.Remove("testdata/state.json")

	s := state.New()
	s.SwitchStates["foo"] = state.SwitchStateOff

	assert.NoError(t, state.Save("testdata/state.json", s))

	data, err := ioutil.ReadFile("testdata/state.json")

	assert.NoError(t, err)
	assert.Equal(t, "{\"switch_states\":{\"foo\":0},\"schedules\":{}}\n", string(data))
}

func TestSaveInNonexistentDir(t *testing.T) {
	s := state.New()
	assert.Error(t, state.Save("testdata/thisdoesnotexist/state.json", s))
}

type errorWriter struct{}

func (errorWriter) Write(p []byte) (n int, err error) {
	return 1, fmt.Errorf("error")
}

func TestSaveWithWriter(t *testing.T) {
	var buf []byte

	w := bytes.NewBuffer(buf)

	assert.NoError(t, state.SaveWithWriter(w, state.New()))
	assert.Equal(t, "{\"switch_states\":{},\"schedules\":{}}\n", w.String())
}

func TestSaveWithBadWriter(t *testing.T) {
	assert.Error(t, state.SaveWithWriter(errorWriter{}, state.New()))
}
