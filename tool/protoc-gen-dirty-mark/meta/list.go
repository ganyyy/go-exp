package meta

import (
	"iter"
	"slices"
)

type innerList[V, T any] struct {
	dirtyMark
	values []V
	transfer[V, T]
}

// Count returns the count of the list.
func (m *innerList[V, T]) Count() int {
	return len(m.values)
}

// Get returns the value of the index.
func (m *innerList[V, T]) Get(index int) V {
	v, _ := m.GetExisted(index)
	return v
}

// GetExisted returns the value of the index if it exists.
func (m *innerList[V, T]) GetExisted(index int) (V, bool) {
	if index < 0 || index >= len(m.values) {
		var empty V
		return empty, false
	}
	v := m.values[index]
	m.setHook(v)
	return v, true
}

// Insert inserts the value to the list.
func (m *innerList[V, T]) Insert(index int, value V) {
	m.dirty()
	if index < 0 || index > len(m.values) {
		return
	}
	m.setHook(value)
	m.values = append(m.values[:index], append([]V{value}, m.values[index:]...)...)
}

// Add adds the value to the list.
func (m *innerList[V, T]) Add(value ...V) {
	m.dirty()
	for _, v := range value {
		m.setHook(v)
		m.values = append(m.values, v)
	}
}

// Remove removes the value from the list.
func (m *innerList[V, T]) Remove(index int) (V, bool) {
	if index < 0 || index >= len(m.values) {
		var empty V
		return empty, false
	}
	m.dirty()
	v := m.values[index]
	m.delHook(v)
	m.values = append(m.values[:index], m.values[index+1:]...)
	return v, true
}

// FromProto sets the value from the target.
func (m *innerList[V, T]) FromProto(target []T) {
	m.dirty()
	m.values = make([]V, len(target))
	for i, t := range target {
		m.values[i] = m.t2v(t)
	}
}

// ToProto gets the target from the value.
func (m *innerList[V, T]) ToProto() []T {
	ret := make([]T, len(m.values))
	for i, v := range m.values {
		ret[i] = m.v2t(v)
	}
	return ret
}

// Range iterates the list.
func (m *innerList[V, T]) Range() iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		for i, v := range m.values {
			if !yield(i, v) {
				break
			}
		}
	}
}

// DirtyCollect collects the list from the target.
func (m *innerList[V, T]) DirtyCollect(target []T) []T {
	target = slices.Grow(target, m.Count())
	for _, v := range m.values {
		target = append(target, m.v2t(v))
	}
	if len(target) == 0 {
		return make([]T, 0)
	}
	return target
}

// ============================

type ValueList[V any] struct {
	innerList[V, V]
}

// FromProto sets the value from the target.
func (m *ValueList[V]) FromProto(target []V) {
	m.dirty()
	m.values = make([]V, len(target))
	copy(m.values, target)
}

// ToProto gets the target from the value.
func (m *ValueList[V]) ToProto() []V {
	ret := make([]V, len(m.values))
	copy(ret, m.values)
	return ret
}

func NewValueList[V any]() *ValueList[V] {
	var m = &ValueList[V]{}
	m.transfer = ValueTransfer[V]()
	return m
}

type ReferenceList[V IValue[T], T any] struct {
	innerList[V, T]
}

func NewReferenceList[V IValue[T], T any]() *ReferenceList[V, T] {
	var m = &ReferenceList[V, T]{}
	m.transfer = ReferenceTransfer[V](m)
	return m
}
