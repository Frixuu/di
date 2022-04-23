package di

import "sync"

// typedSyncMap is a typed wrapper over sync.Map.
type typedSyncMap[K any, V any] struct {
	m sync.Map
}

// Get tries to return the value stored in the map for the key.
func (m *typedSyncMap[K, V]) Get(key K) (value V, ok bool) {
	v, ok := m.m.Load(key)
	if !ok {
		return
	}

	value, ok = v.(V)
	return
}

// Put sets the value for the key.
func (m *typedSyncMap[K, V]) Put(key K, value V) {
	m.m.Store(key, value)
}

// ComputeIfAbsent puts a value in a map if it is absent there.
func (m *typedSyncMap[K, V]) ComputeIfAbsent(key K, factory func() V) (value V) {
	value, ok := m.Get(key)
	if ok {
		return
	}

	v, _ := m.m.LoadOrStore(key, factory())
	value, _ = v.(V)
	return
}
