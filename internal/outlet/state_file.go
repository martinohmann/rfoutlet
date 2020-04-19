package outlet

import (
	"encoding/json"
	"io/ioutil"

	"github.com/martinohmann/rfoutlet/internal/schedule"
)

// outletState represents the state of a single outlet.
type outletState struct {
	State    State              `json:"state,omitempty"`
	Schedule *schedule.Schedule `json:"schedule,omitempty"`
}

// StateFile holds the state and schedule of all configured outlets. This is
// used as persistence across rfoutlet restarts.
type StateFile struct {
	Filename string
}

// NewStateFile creates a new *StateFile with filename.
func NewStateFile(filename string) *StateFile {
	return &StateFile{
		Filename: filename,
	}
}

// ReadBack reads outlet state from the state file back into the passed in
// outlets. Returns an error if the state file cannot be read or if its
// contents are invalid.
func (f *StateFile) ReadBack(outlets []*Outlet) error {
	buf, err := ioutil.ReadFile(f.Filename)
	if err != nil {
		return err
	}

	var stateMap map[string]outletState

	err = json.Unmarshal(buf, &stateMap)
	if err != nil {
		return err
	}

	applyOutletStates(outlets, stateMap)

	return nil
}

// WriteOut writes out the outlet states to the state file. Returns an error if
// writing the state file fails.
func (f *StateFile) WriteOut(outlets []*Outlet) error {
	stateMap := collectOutletStates(outlets)

	buf, err := json.Marshal(stateMap)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(f.Filename, buf, 0664)
}

func applyOutletStates(outlets []*Outlet, stateMap map[string]outletState) {
	for _, o := range outlets {
		outletState, ok := stateMap[o.ID]
		if !ok {
			continue
		}

		o.SetState(outletState.State)
		o.Schedule = outletState.Schedule
		if o.Schedule == nil {
			o.Schedule = schedule.New()
		}
	}
}

func collectOutletStates(outlets []*Outlet) map[string]outletState {
	stateMap := make(map[string]outletState)

	for _, o := range outlets {
		stateMap[o.ID] = outletState{
			State:    o.GetState(),
			Schedule: o.Schedule,
		}
	}

	return stateMap
}
