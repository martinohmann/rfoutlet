package outlet

import "fmt"

// Registry holds references to all outlets and outlet groups.
type Registry struct {
	outlets   []*Outlet
	outletMap map[string]*Outlet
	groups    []*Group
	groupMap  map[string]*Group
}

// NewRegistry creates a new *Registry.
func NewRegistry() *Registry {
	return &Registry{
		outlets:   make([]*Outlet, 0),
		outletMap: make(map[string]*Outlet),
		groups:    make([]*Group, 0),
		groupMap:  make(map[string]*Group),
	}
}

// RegisterGroups registers groups and the outlets contained in those groups.
// Returns an error if groups or outlets with duplicate IDs are found.
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

		log.WithField("groupID", group.ID).Info("registered outlet group")
	}

	return nil
}

// RegisterOutlets registers outlets. Returns an error if outlets with duplicate
// IDs are found.
func (r *Registry) RegisterOutlets(outlets ...*Outlet) error {
	for _, outlet := range outlets {
		_, ok := r.outletMap[outlet.ID]
		if ok {
			return fmt.Errorf("duplicate outlet ID %q", outlet.ID)
		}

		r.outletMap[outlet.ID] = outlet
		r.outlets = append(r.outlets, outlet)

		log.WithField("outletID", outlet.ID).Info("registered outlet")
	}

	return nil
}

// GetOutlet fetches the an outlet from the registry by ID. The second return
// value is true if the outlet was found, false otherwise.
func (r *Registry) GetOutlet(id string) (*Outlet, bool) {
	outlet, ok := r.outletMap[id]
	return outlet, ok
}

// GetOutlets returns all registered outlets.
func (r *Registry) GetOutlets() []*Outlet {
	return r.outlets
}

// GetGroup fetches the an group from the registry by ID. The second return
// value is true if the group was found, false otherwise.
func (r *Registry) GetGroup(id string) (*Group, bool) {
	group, ok := r.groupMap[id]
	return group, ok
}

// GetGroups returns all registered groups.
func (r *Registry) GetGroups() []*Group {
	return r.groups
}
