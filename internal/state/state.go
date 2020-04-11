package state

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/internal/schedule"
)

// State holds the state and schedule of each configured outlet.
type State map[string]OutletState

// OutletState represents the state of a single outlet.
type OutletState struct {
	State    outlet.State       `json:"state"`
	Schedule *schedule.Schedule `json:"schedule,omitempty"`
}

// Load loads the state from a file
func Load(file string) (State, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return LoadWithReader(f)
}

// LoadWithReader loads the state using reader
func LoadWithReader(r io.Reader) (State, error) {
	c, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	state := State{}

	err = json.Unmarshal(c, &state)
	if err != nil {
		return nil, err
	}

	return state, nil
}

// Save saves the state to a file
func Save(file string, state State) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return err
	}

	f.Truncate(0)

	return SaveWithWriter(f, state)
}

// SaveWithWriter saves the state using writer
func SaveWithWriter(w io.Writer, state State) error {
	return json.NewEncoder(w).Encode(state)
}

// Apply applies the state to outlets.
func (s State) Apply(outlets []*outlet.Outlet) {
	for _, o := range outlets {
		outletState, ok := s[o.ID]
		if !ok {
			continue
		}

		o.SetState(outletState.State)
		o.Schedule = outletState.Schedule
	}
}

// Collect collects the states of passed outlets.
func Collect(outlets []*outlet.Outlet) State {
	s := State{}

	for _, o := range outlets {
		s[o.ID] = OutletState{
			State:    o.GetState(),
			Schedule: o.Schedule,
		}
	}

	return s
}
