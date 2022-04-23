package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapLoadAndStore(t *testing.T) {
	m := typedSyncMap[int, string]{}
	m.Put(1, "one")
	m.Put(2, "two")
	m.Put(3, "three")

	v, ok := m.Get(2)
	assert.True(t, ok)
	assert.Equal(t, "two", v)

	_, ok = m.Get(4)
	assert.False(t, ok)
}
