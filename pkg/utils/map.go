package utils

import (
	"container/list"
	"reflect"
	"sync"
)

// Map is a safe map type when parally read & write
type Map struct {
	length int
	// change  bool
	// details sync.Map // map[key] int:existed/created/changed/deleted
	sync.Map
}

// NewMap creates a new Map instance
func NewMap() Map { return Map{length: 0, Map: sync.Map{}} }

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

// Count returns
func (m *Map) Count() int {
	m.length = 0
	m.Map.Range(func(k, v interface{}) bool {
		m.length++
		return true
	})
	return m.length
}

// Len returns
func (m *Map) Len() int {
	if m == nil {
		return -1
	}
	return m.length
}

// Store sets the value for a key.
func (m *Map) Store(key, value interface{}) {
	defer Recover()

	if v, ok := m.Map.Load(key); !ok || !reflect.DeepEqual(v, value) {
		m.length++
		m.Map.Store(key, value)
	}
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (m *Map) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	actual, loaded = m.Map.LoadOrStore(key, value)
	if !loaded {
		m.length++
	}
	return
}

// Delete deletes the value for a key.
func (m *Map) Delete(key interface{}) {
	if _, ok := m.Map.Load(key); ok {
		m.length--
		m.Map.Delete(key)
	}
}

type MapList struct {
	dataMap  map[string]*list.Element
	dataList *list.List
}

func NewMapList() *MapList {
	return &MapList{
		dataMap:  make(map[string]*list.Element),
		dataList: list.New(),
	}
}

func (mapList *MapList) Exists(data Keyer) bool {
	_, exists := mapList.dataMap[string(data.GetKey())]
	return exists
}

func (mapList *MapList) Push(data Keyer) bool {
	if mapList.Exists(data) {
		return false
	}
	elem := mapList.dataList.PushBack(data)
	mapList.dataMap[data.GetKey()] = elem
	return true
}

func (mapList *MapList) Remove(data Keyer) {
	if !mapList.Exists(data) {
		return
	}
	mapList.dataList.Remove(mapList.dataMap[data.GetKey()])
	delete(mapList.dataMap, data.GetKey())
}

func (mapList *MapList) Len() int { return mapList.dataList.Len() }

func (mapList *MapList) Walk(cb func(data Keyer)) {
	for elem := mapList.dataList.Front(); elem != nil; elem = elem.Next() {
		cb(elem.Value.(Keyer))
	}
}

type Keyer interface {
	GetKey() string
}
type Elements struct {
	value string
}

func (e Elements) GetKey() string {
	return e.value
}

// func main() {
// 	fmt.Println("Starting test...")
// 	ml := NewMapList()
// 	var a, b, c Keyer
// 	a = &Elements{"Alice"}
// 	b = &Elements{"Bob"}
// 	c = &Elements{"Conrad"}
// 	ml.Push(a)
// 	ml.Push(b)
// 	ml.Push(c)
// 	cb := func(data Keyer) {
// 		fmt.Println(ml.dataMap[data.GetKey()].Value.(*Elements).value)
// 	}
// 	fmt.Println("Print elements in the order of pushing:")
// 	ml.Walk(cb)
// 	fmt.Printf("Size of MapList: %d \n", ml.Size())
// 	ml.Remove(b)
// 	fmt.Println("After removing b:")
// 	ml.Walk(cb)
// 	fmt.Printf("Size of MapList: %d \n", ml.Size())
// }
