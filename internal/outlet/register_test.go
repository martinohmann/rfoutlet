package outlet

import (
	"testing"

	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/stretchr/testify/assert"
)

var testStateHandler = &nopStateHandler{}

type nopStateHandler struct{}

func (nopStateHandler) LoadState(o []*Outlet) error {
	return nil
}

func (nopStateHandler) SaveState(o []*Outlet) error {
	return nil
}

func TestRegisterFromConfig(t *testing.T) {

	tests := []struct {
		config  *config.Config
		wantErr bool
		errMsg  string
		assert  func(t *testing.T, m *Manager)
	}{
		{
			config: &config.Config{
				GroupOrder: []string{"foo"},
				Groups: map[string]*config.Group{
					"foo": {Name: "bar", Outlets: []string{"baz"}},
				},
				Outlets: map[string]*config.Outlet{
					"baz": {Name: "qux"},
				},
			},
			assert: func(t *testing.T, m *Manager) {
				assert.Len(t, m.Outlets(), 1)

				if assert.Len(t, m.Groups(), 1) {
					g, err := m.GetGroup("foo")
					assert.NoError(t, err)
					assert.Equal(t, "bar", g.Name)
					assert.Len(t, g.Outlets, 1)
				}
			},
		},
		{
			config: &config.Config{
				GroupOrder: []string{"bar", "foo", "baz"},
				Groups: map[string]*config.Group{
					"foo": {Name: "Foo"},
					"bar": {Name: "Bar"},
					"baz": {Name: "Baz"},
				},
			},
			assert: func(t *testing.T, m *Manager) {
				assert.Len(t, m.groupOrder, 3)
				assert.Equal(t, []string{"bar", "foo", "baz"}, m.groupOrder)
			},
		},
		{
			config: &config.Config{
				GroupOrder: []string{"baz"},
				Groups: map[string]*config.Group{
					"foo": {
						Name:    "bar",
						Outlets: []string{"baz"},
					},
				},
			},
			wantErr: true,
			errMsg:  `unknown group "baz"`,
		},
		{
			config: &config.Config{
				GroupOrder: []string{"foo"},
				Groups: map[string]*config.Group{
					"foo": {Name: "bar", Outlets: []string{"qux"}},
				},
			},
			wantErr: true,
			errMsg:  `unknown outlet "qux"`,
		},
	}

	for _, tt := range tests {
		m := NewManager(testStateHandler)
		err := RegisterFromConfig(m, tt.config)

		if tt.wantErr {
			if assert.Error(t, err) && tt.errMsg != "" {
				assert.Equal(t, tt.errMsg, err.Error())
			}
		} else {
			if assert.NoError(t, err) {
				tt.assert(t, m)
			}
		}
	}
}
