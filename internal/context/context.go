package context

import (
	ctx "context"
	"fmt"
	"sync"

	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/martinohmann/rfoutlet/internal/state"
	uuid "github.com/satori/go.uuid"
)

// Context type definition
type Context struct {
	ctx.Context

	state     *state.State
	groupMap  map[string]*Group
	outletMap map[string]*Outlet

	Config *config.Config
	Groups []*Group `json:"groups"`
}

// Group type definition
type Group struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Outlets []*Outlet `json:"outlets"`
}

// Outlet type definition
type Outlet struct {
	sync.RWMutex
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	CodeOn      uint64            `json:"-"`
	CodeOff     uint64            `json:"-"`
	Protocol    int               `json:"-"`
	PulseLength uint              `json:"-"`
	Schedule    schedule.Schedule `json:"schedule"`
	State       state.SwitchState `json:"state"`
}

// New create a new context for config and state
func New(c *config.Config, s *state.State) (*Context, error) {
	return Wrap(ctx.Background(), c, s)
}

// Wrap wraps an existing context
func Wrap(ctx ctx.Context, c *config.Config, s *state.State) (*Context, error) {
	context := &Context{
		Context:   ctx,
		Config:    c,
		state:     s,
		groupMap:  make(map[string]*Group),
		outletMap: make(map[string]*Outlet),
		Groups:    make([]*Group, 0, len(c.Groups)),
	}

	if err := context.buildGroups(); err != nil {
		return nil, err
	}

	return context, nil
}

func (c *Context) buildGroups() error {
	for _, id := range c.Config.GroupOrder {
		g, ok := c.Config.Groups[id]
		if !ok {
			return fmt.Errorf("invalid group identifier %q", id)
		}

		group := &Group{
			ID:      id,
			Name:    g.Name,
			Outlets: make([]*Outlet, 0, len(g.Outlets)),
		}

		if err := c.buildOutlets(g, group); err != nil {
			return err
		}

		c.groupMap[id] = group
		c.Groups = append(c.Groups, group)
	}

	return nil
}

func (c *Context) buildOutlets(g *config.Group, group *Group) error {
	for _, id := range g.Outlets {
		o, ok := c.Config.Outlets[id]
		if !ok {
			return fmt.Errorf("invalid outlet identifier %q", id)
		}

		outlet := &Outlet{
			ID:          id,
			Name:        o.Name,
			CodeOn:      o.CodeOn,
			CodeOff:     o.CodeOff,
			Protocol:    o.Protocol,
			PulseLength: o.PulseLength,
			State:       c.state.SwitchStates[id],
			Schedule:    c.state.Schedules[id],
		}

		c.outletMap[id] = outlet
		group.Outlets = append(group.Outlets, outlet)
	}

	return nil
}

// GetOutlet retrieves the outlet with given id
func (c *Context) GetOutlet(id string) (*Outlet, error) {
	outlet, ok := c.outletMap[id]
	if !ok {
		return nil, fmt.Errorf("outlet with identifier %q does not exist", id)
	}

	return outlet, nil
}

// GetGroup retrieves the group with given id
func (c *Context) GetGroup(id string) (*Group, error) {
	group, ok := c.groupMap[id]
	if !ok {
		return nil, fmt.Errorf("group with identifier %q does not exist", id)
	}

	return group, nil
}

// CollectState collects the state of all outlets
func (c *Context) CollectState() *state.State {
	state := state.New()

	for _, o := range c.outletMap {
		state.Schedules[o.ID] = o.GetSchedule()
		state.SwitchStates[o.ID] = o.GetSwitchState()
	}

	return state
}

// GetSwitchState returns the switch state of an outlet (thread-safe)
func (o *Outlet) GetSwitchState() state.SwitchState {
	o.RLock()
	defer o.RUnlock()

	return o.State
}

// SetSwitchState sets the switch state of an outlet (thread-safe)
func (o *Outlet) SetSwitchState(state state.SwitchState) {
	o.Lock()
	defer o.Unlock()

	o.State = state
}

// GetSchedule returns the schedule of an outlet (thread-safe)
func (o *Outlet) GetSchedule() schedule.Schedule {
	o.RLock()
	defer o.RUnlock()

	return o.Schedule
}

// AddInterval adds an interval to the schedule of an outlet (thread-safe)
func (o *Outlet) AddInterval(interval schedule.Interval) error {
	if interval.ID == "" {
		interval.ID = uuid.NewV4().String()
	}

	o.Lock()
	defer o.Unlock()

	for _, i := range o.Schedule {
		if i.ID == interval.ID {
			return fmt.Errorf("interval with identifier %q already exists", interval.ID)
		}
	}

	o.Schedule = append(o.Schedule, interval)

	return nil
}

// UpdateInterval updates an interval of the schedule of an outlet
// (thread-safe). Will return an error if the interval does not exist.
func (o *Outlet) UpdateInterval(interval schedule.Interval) error {
	o.Lock()
	defer o.Unlock()

	for i, intv := range o.Schedule {
		if intv.ID != interval.ID {
			continue
		}

		o.Schedule[i] = interval

		return nil
	}

	return fmt.Errorf("interval with identifier %q does not exist", interval.ID)
}

// DeleteInterval deletes an interval of the schedule of an outlet
// (thread-safe). Will return an error if the interval does not exist.
func (o *Outlet) DeleteInterval(interval schedule.Interval) error {
	o.Lock()
	defer o.Unlock()

	for i, intv := range o.Schedule {
		if intv.ID == interval.ID {
			o.Schedule = append(o.Schedule[:i], o.Schedule[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("interval with identifier %q does not exist", interval.ID)
}
