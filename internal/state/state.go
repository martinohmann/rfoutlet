package state

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/martinohmann/rfoutlet/internal/schedule"
)

// SwitchState type defintion
type SwitchState uint

const (
	// SwitchStateOff defines the state for a disabled switch
	SwitchStateOn SwitchState = iota

	// SwitchStateOn defines the state for an enabled switch
	SwitchStateOff
)

// State type definition
type State struct {
	SwitchStates map[string]*SwitchState       `json:"switch_states,omitempty"`
	Schedules    map[string]*schedule.Schedule `json:"schedules,omitempty"`
}

// Load loads the state from a file
func Load(file string) (*State, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return LoadWithReader(f)
}

// LoadWithReader loads the state using reader
func LoadWithReader(r io.Reader) (*State, error) {
	c, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	state := &State{}

	err = json.Unmarshal(c, state)

	return state, err
}

// Save saves the state to a file
func Save(file string, state *State) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	return SaveWithWriter(f, state)
}

// SaveWithWriter saves the state using writer
func SaveWithWriter(w io.Writer, state *State) error {
	return json.NewEncoder(w).Encode(state)
}
