package outlet

import (
	"encoding/json"
	"io/ioutil"

	"github.com/spf13/afero"
)

// State type definition
type State int

const (
	// StateUnknown defines an unknown outlet state
	StateUnknown State = iota

	// StateOn defines an outlet that is on
	StateOn

	// StateOff defines an outlet that is off
	StateOff
)

type stateInfo struct {
	Outlet      int   `json:"outlet"`
	Group       int   `json:"group"`
	SwitchState State `json:"state"`
}

type StateManager interface {
	SaveState(*Control) error
	RestoreState(*Control) error
}

type FileStateManager struct {
	f afero.File
}

func NewStateManager(stateFile afero.File) *FileStateManager {
	return &FileStateManager{f: stateFile}
}

func (m *FileStateManager) RestoreState(control *Control) error {
	m.f.Seek(0, 0)

	b, err := ioutil.ReadAll(m.f)
	if err != nil {
		return err
	}

	if len(b) == 0 {
		return nil
	}

	stateInfos := make([]stateInfo, 0)

	if err = json.Unmarshal(b, &stateInfos); err != nil {
		return err
	}

	return m.restoreState(control, stateInfos)
}

func (m *FileStateManager) SaveState(control *Control) error {
	stateInfos := make([]stateInfo, 0)

	for i, og := range control.OutletGroups() {
		for j, o := range og.Outlets {
			stateInfo := stateInfo{
				Group:       i,
				Outlet:      j,
				SwitchState: o.State,
			}

			stateInfos = append(stateInfos, stateInfo)
		}
	}

	return m.saveState(stateInfos)
}

func (m *FileStateManager) restoreState(control *Control, stateInfos []stateInfo) error {
	for _, s := range stateInfos {
		o, err := control.Outlet(s.Group, s.Outlet)
		if err != nil {
			return err
		}

		o.State = s.SwitchState
	}

	return nil
}

func (m *FileStateManager) saveState(stateInfos []stateInfo) error {
	m.f.Truncate(0)
	m.f.Seek(0, 0)

	if err := json.NewEncoder(m.f).Encode(stateInfos); err != nil {
		return err
	}

	return m.f.Sync()
}

type NullStateManager struct{}

func NewNullStateManager() *NullStateManager                    { return &NullStateManager{} }
func (m *NullStateManager) RestoreState(control *Control) error { return nil }
func (m *NullStateManager) SaveState(control *Control) error    { return nil }
