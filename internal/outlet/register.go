package outlet

import (
	"fmt"

	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/martinohmann/rfoutlet/internal/schedule"
)

// RegisterFromConfig registers configured outlets and groups to the manager
func RegisterFromConfig(m *Manager, config *config.Config) error {
	if err := registerGroups(m, config); err != nil {
		return err
	}

	return nil
}

// registerGroups populates the groups from config
func registerGroups(m *Manager, config *config.Config) error {
	groupOrder := config.GroupOrder

	m.SetGroupOrder(groupOrder)

	for _, id := range groupOrder {
		g, ok := config.Groups[id]
		if !ok {
			return fmt.Errorf("unknown group %q", id)
		}

		group := &Group{
			ID:      id,
			Name:    g.Name,
			Outlets: make([]*Outlet, 0, len(g.Outlets)),
		}

		if err := registerOutlets(m, config, group, g.Outlets); err != nil {
			return err
		}

		m.RegisterGroup(id, group)
	}

	return nil
}

// registerOutlets populates the outlets and assigns them to groups as defined in
// the config
func registerOutlets(m *Manager, config *config.Config, group *Group, outletOrder []string) error {
	for _, id := range outletOrder {
		o, ok := config.Outlets[id]
		if !ok {
			return fmt.Errorf("unknown outlet %q", id)
		}

		outlet := &Outlet{
			ID:          id,
			Name:        o.Name,
			CodeOn:      o.CodeOn,
			CodeOff:     o.CodeOff,
			Protocol:    o.Protocol,
			PulseLength: o.PulseLength,
			Schedule:    schedule.New(),
		}

		group.Outlets = append(group.Outlets, outlet)

		m.Register(id, outlet)
	}

	return nil
}
