package meta

import "iter"

// innerMap is a map that stores the value of each field.
// K: the type of the key
// V: the type of the value
type innerMap[K comparable, V, T any] struct {
	dirtyMark
	transfer[V, T]
	kvs map[K]V
}

// Count returns the count of the map.
func (m *innerMap[K, V, T]) Count() int {
	return len(m.kvs)
}

// Get returns the value of the key.
func (m *innerMap[K, V, T]) Get(key K) V {
	v, _ := m.GetExisted(key)
	return v
}

// GetExisted returns the value of the key if it exists.
func (m *innerMap[K, V, T]) GetExisted(key K) (V, bool) {
	v, ok := m.kvs[key]
	if ok {
		m.setHook(v)
	}
	return v, ok
}

// Set sets the value of the key.
func (m *innerMap[K, V, T]) Set(key K, value V) {
	m.dirty()
	if m.kvs == nil {
		m.kvs = make(map[K]V)
	}
	m.setHook(value)
	m.kvs[key] = value
}

// Del deletes the value of the key.
func (m *innerMap[K, V, T]) Del(key K) V {
	if v, ok := m.kvs[key]; !ok {
		return v // return the default value
	} else {
		m.dirty()
		m.delHook(v)
		delete(m.kvs, key)
		return v
	}
}

// FromProto sets the value from the target.
func (m *innerMap[K, V, T]) FromProto(target map[K]T) {
	m.dirty()
	m.kvs = make(map[K]V, len(target))
	for k, t := range target {
		m.kvs[k] = m.t2v(t)
	}
}

// ToProto gets the target from the value.
func (m *innerMap[K, V, T]) ToProto() map[K]T {
	ret := make(map[K]T, len(m.kvs))
	for k, v := range m.kvs {
		ret[k] = m.v2t(v)
	}
	return ret
}

// Merge merges the value from the target.
func (m *innerMap[K, V, T]) Merge(target map[K]T) {
	if len(target) == 0 {
		return
	}
	m.dirty()
	for k, v := range target {
		m.kvs[k] = m.t2v(v)
	}
}

// Range applies the function to each key-value pair.
func (m *innerMap[K, V, T]) Range() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m.kvs {
			if !yield(k, v) {
				break
			}
		}
	}
}

// DirtyCollect overwrites Mark collect method. It collects the value of each field, not only the dirty mark.
func (m *innerMap[K, V, T]) DirtyCollect(target map[K]T) map[K]T {
	if target == nil {
		target = make(map[K]T)
	}
	for k, v := range m.kvs {
		target[k] = m.v2t(v)
	}
	return target
}

func NewValueMap[K comparable, V any]() *ValueMap[K, V] {
	var m = &ValueMap[K, V]{}
	m.transfer = ValueTransfer[V]()
	return m
}

type ValueMap[K comparable, V any] struct {
	innerMap[K, V, V]
}

func NewReferenceMap[K comparable, V IValue[T], T any]() *ReferenceMap[K, V, T] {
	var m = &ReferenceMap[K, V, T]{}
	m.transfer = ReferenceTransfer[V](m)
	return m
}

type ReferenceMap[K comparable, V IValue[T], T any] struct {
	innerMap[K, V, T]
}

// ============================
