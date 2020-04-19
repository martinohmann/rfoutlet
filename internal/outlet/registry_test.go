package outlet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	r := NewRegistry()

	require.NoError(t, r.RegisterGroups(&Group{ID: "foo"}))
	assert.Len(t, r.GetGroups(), 1)
	assert.Len(t, r.GetOutlets(), 0)

	require.Error(t, r.RegisterGroups(&Group{ID: "foo"}))
	require.NoError(t, r.RegisterGroups(&Group{ID: "bar", Outlets: []*Outlet{{ID: "baz"}}}))
	require.Error(t, r.RegisterGroups(&Group{ID: "baz", Outlets: []*Outlet{{ID: "baz"}}}))
	assert.Len(t, r.GetGroups(), 2)
	assert.Len(t, r.GetOutlets(), 1)

	group, ok := r.GetGroup("foo")
	assert.True(t, ok)
	assert.Equal(t, &Group{ID: "foo"}, group)

	_, ok = r.GetGroup("non-existent")
	assert.False(t, ok)

	outlet, ok := r.GetOutlet("baz")
	assert.True(t, ok)
	assert.Equal(t, &Outlet{ID: "baz"}, outlet)

	_, ok = r.GetOutlet("non-existent")
	assert.False(t, ok)
}
