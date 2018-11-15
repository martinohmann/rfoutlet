package context_test

import (
	"testing"

	"github.com/martinohmann/rfoutlet/internal/config"
	"github.com/martinohmann/rfoutlet/internal/context"
	"github.com/martinohmann/rfoutlet/internal/state"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		config  *config.Config
		state   *state.State
		wantErr bool
		errMsg  string
		assert  func(t *testing.T, context *context.Context)
	}{
		{
			config: &config.Config{
				GroupOrder: []string{"foo"},
				Groups: map[string]*config.Group{
					"foo": &config.Group{
						Name:    "bar",
						Outlets: []string{"baz"},
					},
				},
				Outlets: map[string]*config.Outlet{
					"baz": &config.Outlet{
						Name: "qux",
					},
				},
			},
			state: &state.State{},
			assert: func(t *testing.T, context *context.Context) {
				if assert.Len(t, context.Groups, 1) {
					assert.Equal(t, "bar", context.Groups[0].Name)

					if assert.Len(t, context.Groups[0].Outlets, 1) {
						assert.Equal(t, "qux", context.Groups[0].Outlets[0].Name)
					}
				}
			},
		},
		{
			config: &config.Config{
				GroupOrder: []string{"bar", "foo"},
				Groups: map[string]*config.Group{
					"foo": &config.Group{
						Name:    "bar",
						Outlets: []string{},
					},
					"bar": &config.Group{
						Name:    "baz",
						Outlets: []string{},
					},
				},
			},
			state: &state.State{},
			assert: func(t *testing.T, context *context.Context) {
				if assert.Len(t, context.Groups, 2) {
					assert.Equal(t, "baz", context.Groups[0].Name)
				}
			},
		},
		{
			config: &config.Config{
				GroupOrder: []string{"baz"},
				Groups: map[string]*config.Group{
					"foo": &config.Group{
						Name:    "bar",
						Outlets: []string{},
					},
				},
			},
			state:   &state.State{},
			wantErr: true,
			errMsg:  `invalid group identifier "baz"`,
		},
		{
			config: &config.Config{
				GroupOrder: []string{"foo"},
				Groups: map[string]*config.Group{
					"foo": &config.Group{
						Name:    "bar",
						Outlets: []string{"qux"},
					},
				},
				Outlets: map[string]*config.Outlet{
					"baz": &config.Outlet{
						Name: "qux",
					},
				},
			},
			state:   &state.State{},
			wantErr: true,
			errMsg:  `invalid outlet identifier "qux"`,
		},
	}

	for _, tt := range tests {
		context, err := context.New(tt.config, tt.state)

		if tt.wantErr {
			if assert.Error(t, err) {
				assert.Equal(t, tt.errMsg, err.Error())
			}
		} else if assert.NoError(t, err) {
			tt.assert(t, context)
		}
	}
}

func TestGet(t *testing.T) {
	c := &config.Config{
		GroupOrder: []string{"foo"},
		Groups: map[string]*config.Group{
			"foo": &config.Group{
				Name:    "bar",
				Outlets: []string{"baz"},
			},
		},
		Outlets: map[string]*config.Outlet{
			"baz": &config.Outlet{
				Name: "qux",
			},
		},
	}

	s := &state.State{}

	context, err := context.New(c, s)

	if assert.NoError(t, err) {
		g, err := context.GetGroup("foo")
		if assert.NoError(t, err) {
			assert.Equal(t, "bar", g.Name)
		}

		o, err := context.GetOutlet("baz")
		if assert.NoError(t, err) {
			assert.Equal(t, "qux", o.Name)
		}

		_, err = context.GetGroup("f")
		assert.Error(t, err)

		_, err = context.GetOutlet("b")
		assert.Error(t, err)
	}
}
