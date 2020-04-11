package outlet

import (
	"fmt"
)

// Registry holds references to all outlets and outlet groups.
type Registry struct {
	outlets   []*Outlet
	outletMap map[string]*Outlet
	groups    []*Group
	groupMap  map[string]*Group
}

func NewRegistry() *Registry {
	return &Registry{
		outlets:   make([]*Outlet, 0),
		outletMap: make(map[string]*Outlet),
		groups:    make([]*Group, 0),
		groupMap:  make(map[string]*Group),
	}
}

func (r *Registry) RegisterGroups(groups ...*Group) error {
	for _, group := range groups {
		_, ok := r.groupMap[group.ID]
		if ok {
			return fmt.Errorf("duplicate group ID %q", group.ID)
		}

		err := r.RegisterOutlets(group.Outlets...)
		if err != nil {
			return err
		}

		r.groupMap[group.ID] = group
		r.groups = append(r.groups, group)
	}

	return nil
}

func (r *Registry) RegisterOutlets(outlets ...*Outlet) error {
	for _, outlet := range outlets {
		_, ok := r.outletMap[outlet.ID]
		if ok {
			return fmt.Errorf("duplicate outlet ID %q", outlet.ID)
		}

		r.outletMap[outlet.ID] = outlet
		r.outlets = append(r.outlets, outlet)
	}

	return nil
}

func (r *Registry) GetOutlet(id string) (*Outlet, bool) {
	outlet, ok := r.outletMap[id]
	return outlet, ok
}

func (r *Registry) GetOutlets() []*Outlet {
	return r.outlets
}

func (r *Registry) GetGroup(id string) (*Group, bool) {
	group, ok := r.groupMap[id]
	return group, ok
}

func (r *Registry) GetGroups() []*Group {
	return r.groups
}
