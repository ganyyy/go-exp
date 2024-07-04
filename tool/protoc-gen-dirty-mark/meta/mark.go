package meta

type Mark[T any] struct {
	dyeingMark                    // dyeing is a function that will be called when any field is set
	dirty      map[uint16]func(T) // dirty is a map that stores the dirty mark of each field
}

// Dirty sets the dirty mark of the field.
func (m *Mark[T]) Dirty(idx uint16, set func(T)) {
	if m.dirty == nil {
		m.dirty = make(map[uint16]func(T))
	}
	if _, ok := m.dirty[idx]; !ok {
		m.dirty[idx] = set
	}
	m.dyed()
}

// DirtyCollect applies the dirty mark to the target.
func (m *Mark[T]) DirtyCollect(target T) T {
	for _, set := range m.dirty {
		set(target)
	}
	return target
}

// ResetDirty resets the dirty mark.
func (m *Mark[T]) ResetDirty() {
	m.dirty = nil
}

type dyeingMark struct {
	dyeing func() // dyeing is a function that will be called when any field is set
}

// Dyeing sets the dyeing function.
func (m *dyeingMark) Dyeing(dyeing func()) {
	m.dyeing = dyeing
}

// dyed calls the dyeing function.
func (m *dyeingMark) dyed() {
	if m.dyeing != nil {
		m.dyeing()
	}
}

// GetDyeing gets the dyeing function.
func (m *dyeingMark) GetDyeing() func() {
	return m.dyeing
}
