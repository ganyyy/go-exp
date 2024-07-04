package meta

type innerList[V, T any] struct {
	dyeingMark
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

// Set sets the value of the index.
func (m *innerList[V, T]) Set(index int, value V) {
	m.dyed() // mark self as dirty
	if index < 0 || index >= len(m.values) {
		return
	}
	m.setHook(value)
	m.values[index] = value
}

// Add adds the value to the list.
func (m *innerList[V, T]) Add(value ...V) {
	m.dyed() // mark self as dirty
	for _, v := range value {
		m.setHook(v)
		m.values = append(m.values, v)
	}
}

// Remove removes the value from the list.
func (m *innerList[V, T]) Remove(index int) V {
	if index < 0 || index >= len(m.values) {
		var empty V
		return empty
	}
	m.dyed() // mark self as dirty
	v := m.values[index]
	m.delHook(v)
	m.values = append(m.values[:index], m.values[index+1:]...)
	return v
}

// FromProto sets the value from the target.
func (m *innerList[V, T]) FromProto(target []T) {
	m.dyed() // mark self as dirty
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
func (m *innerList[V, T]) Range(f func(int, V) bool) {
	for i, v := range m.values {
		if !f(i, v) {
			break
		}
	}
}

// DirtyCollect collects the list from the target.
func (m *innerList[V, T]) DirtyCollect(target []T) []T {
	for _, v := range m.values {
		target = append(target, m.v2t(v))
	}
	return target
}

// ============================

type ValueList[V any] struct {
	innerList[V, V]
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
