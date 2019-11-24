package outlet

import "fmt"

// Group type definition
type Group struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Outlets []*Outlet `json:"outlets"`
}

// RegisterGroup registers a group to given name
func (m *Manager) RegisterGroup(name string, group *Group) {
	m.Lock()
	m.groups[name] = group
	m.Unlock()
}

// GetGroup retieves the group for given name
func (m *Manager) GetGroup(name string) (*Group, error) {
	m.Lock()
	defer m.Unlock()

	g, ok := m.groups[name]
	if !ok {
		return nil, fmt.Errorf("unknown group %q", name)
	}

	return g, nil
}

// Groups returns a slice with all registered groups. The slice is ordered by
// the cofigured group order.
func (m *Manager) Groups() []*Group {
	m.Lock()
	defer m.Unlock()

	s := make([]*Group, 0, len(m.groups))
	for _, name := range m.groupOrder {
		s = append(s, m.groups[name])
	}

	return s
}

// SetGroupOrder sets the display order for the groups
func (m *Manager) SetGroupOrder(order []string) {
	m.Lock()
	m.groupOrder = order
	m.Unlock()
}
