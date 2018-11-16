package context

import (
	ctx "context"
	"fmt"

	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/martinohmann/rfoutlet/internal/schedule"
	"github.com/martinohmann/rfoutlet/internal/state"
)

// Context type definition
type Context struct {
	ctx.Context

	Config    *config.Config
	State     *state.State
	groupMap  map[string]*Group
	outletMap map[string]*Outlet

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
		State:     s,
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
			State:       c.State.SwitchStates[id],
			Schedule:    c.State.Schedules[id],
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
