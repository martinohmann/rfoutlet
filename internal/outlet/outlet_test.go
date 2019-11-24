package outlet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	m := NewManager(testStateHandler)
	_, err := m.Get("foo")

	assert.Error(t, err)

	g := &Outlet{}

	m.Register("foo", g)

	r, err := m.Get("foo")

	assert.NoError(t, err)
	assert.Equal(t, g, r)
}

func TestOutlets(t *testing.T) {
	m := NewManager(testStateHandler)
	names := []string{"foo", "baz", "bar"}

	for _, name := range names {
		m.Register(name, &Outlet{ID: name})
	}

	outlets := m.Outlets()

	if assert.Len(t, outlets, 3) {
		for _, name := range names {
			assert.True(t, hasOutletWithName(outlets, name))
		}
	}
}

func hasOutletWithName(outlets []*Outlet, name string) bool {
	for _, o := range outlets {
		if o.ID == name {
			return true
		}
	}

	return false
}
