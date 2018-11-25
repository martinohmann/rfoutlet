package outlet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterGroup(t *testing.T) {
	m := NewManager(testStateHandler)
	_, err := m.GetGroup("foo")

	assert.Error(t, err)

	g := &Group{}

	m.RegisterGroup("foo", g)

	r, err := m.GetGroup("foo")

	assert.NoError(t, err)
	assert.Equal(t, g, r)
}

func TestGroups(t *testing.T) {
	m := NewManager(testStateHandler)
	for _, name := range []string{"foo", "baz", "bar"} {
		m.RegisterGroup(name, &Group{ID: name})
	}

	m.SetGroupOrder([]string{"baz", "bar", "foo"})

	groups := m.Groups()

	if assert.Len(t, groups, 3) {
		assert.Equal(t, "baz", groups[0].ID)
	}
}
