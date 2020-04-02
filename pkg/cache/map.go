package cache

import (
	"reflect"
	"sync"

	"bingo/pkg/utils"
)

const (
	existed = iota
	created
	changed
	deleted
)

type Map struct {
	// length  int
	change  bool
	details sync.Map // map[key] int:existed/created/changed/deleted
	utils.Map
}

// NewMap creates a new Map instance
// func NewMap() Map { return Map{change: false, Map: utils.NewMap()} }

// // Load returns the value stored in the map for a key, or nil if no
// // value is present.
// // The ok result indicates whether value was found in the map.
// func (m *Map) Load(key interface{}) (value interface{}, ok bool) { return m.Load(key) }

// // Range calls f sequentially for each key and value present in the map.
// // If f returns false, range stops the iteration.
// //
// // Range does not necessarily correspond to any consistent snapshot of the Map's
// // contents: no key will be visited more than once, but if the value for any key
// // is stored or deleted concurrently, Range may reflect any mapping for that key
// // from any point during the Range call.
// //
// // Range may be O(N) with the number of elements in the map even if f returns
// // false after a constant number of calls.
// func (m *Map) Range(f func(key, value interface{}) bool) { m.Range(f) }

func (m *Map) reset() { m.change = false }

// Store sets the value for a key.
func (m *Map) Store(key, value interface{}) {
	if v, ok := m.Map.Load(key); !ok || !reflect.DeepEqual(v, value) {
		m.change = true
		m.Map.Store(key, value)

		if !ok {
			// TODO: new content, should update index
			// c.
		}
	}
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *Map) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	actual, loaded = m.Map.LoadOrStore(key, value)
	if !loaded {
		m.change = true
		// TODO: new content, should update index
		// c.
	}
	return
}

// Delete deletes the value for a key.
func (m *Map) Delete(key interface{}) {
	if _, ok := m.Map.Load(key); ok {
		m.change = true
		m.Map.Delete(key)

		if !ok {
			// TODO: new content, should update index
			// c.
		}
	}
}
