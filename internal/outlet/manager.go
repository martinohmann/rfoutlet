package outlet

import (
	"sync"
)

// Manager type definition
type Manager struct {
	sync.Mutex
	sh         StateHandler
	outlets    map[string]*Outlet
	groups     map[string]*Group
	groupOrder []string
}

// NewManager create a new Manager
func NewManager(sh StateHandler) *Manager {
	return &Manager{
		sh:         sh,
		outlets:    make(map[string]*Outlet),
		groups:     make(map[string]*Group),
		groupOrder: make([]string, 0),
	}
}

// SaveState saves the outlet state using the state handler
func (m *Manager) SaveState() error {
	return m.sh.SaveState(m.Outlets())
}

// LoadState loads the outlet state using the state handler
func (m *Manager) LoadState() error {
	return m.sh.LoadState(m.Outlets())
}
