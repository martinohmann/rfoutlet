package state

import "github.com/martinohmann/rfoutlet/internal/outlet"

// Handler type definition
type Handler struct {
	f string
}

// NewHandler create a new state handler
func NewHandler(f string) *Handler {
	return &Handler{f}
}

// LoadState implements the outlet.StateHandler interface
func (h *Handler) LoadState(outlets []*outlet.Outlet) error {
	s, err := Load(h.f)
	s.Apply(outlets)

	return err
}

// SaveState implements the outlet.StateHandler interface
func (h *Handler) SaveState(outlets []*outlet.Outlet) error {
	if h.f == "" {
		return nil
	}

	s := Collect(outlets)

	return Save(h.f, s)
}
