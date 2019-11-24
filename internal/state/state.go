package state

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/martinohmann/rfoutlet/internal/outlet"
	"github.com/martinohmann/rfoutlet/internal/schedule"
)

// State type definition
type State struct {
	SwitchStates map[string]outlet.State       `json:"switch_states"`
	Schedules    map[string]*schedule.Schedule `json:"schedules"`
}

// New create a new empty state
func New() *State {
	return &State{
		SwitchStates: make(map[string]outlet.State),
		Schedules:    make(map[string]*schedule.Schedule),
	}
}

// Load loads the state from a file
func Load(file string) (*State, error) {
	f, err := os.Open(file)
	if err != nil {
		return New(), err
	}

	return LoadWithReader(f)
}

// LoadWithReader loads the state using reader
func LoadWithReader(r io.Reader) (*State, error) {
	c, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	state := New()

	err = json.Unmarshal(c, state)

	return state, err
}

// Save saves the state to a file
func Save(file string, state *State) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return err
	}

	f.Truncate(0)

	return SaveWithWriter(f, state)
}

// SaveWithWriter saves the state using writer
func SaveWithWriter(w io.Writer, state *State) error {
	return json.NewEncoder(w).Encode(state)
}

// Apply applies the state to outlets
func (s *State) Apply(outlets []*outlet.Outlet) {
	for _, o := range outlets {
		if state, ok := s.SwitchStates[o.ID]; ok {
			o.SetState(state)
		}

		if schedule := s.Schedules[o.ID]; schedule != nil {
			o.Schedule = schedule
		}
	}
}

// Collect collects the states of passed outlets
func Collect(outlets []*outlet.Outlet) *State {
	s := New()

	for _, o := range outlets {
		if o.Schedule != nil {
			s.Schedules[o.ID] = o.Schedule
		}

		s.SwitchStates[o.ID] = o.GetState()
	}

	return s
}
