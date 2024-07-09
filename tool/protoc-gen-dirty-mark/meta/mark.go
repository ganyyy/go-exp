package meta

import "unsafe"

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

type bitType uint8

const (
	bitSize = uint32(unsafe.Sizeof(bitType(0)) * 8)
)

type BitsetMark struct {
	dyeingMark
	bits     []bitType
	maxIndex uint32
	isDirty  bool
}

func NewBitsetMark(size uint32) *BitsetMark {
	return &BitsetMark{bits: make([]bitType, (size+bitSize-1)/bitSize), maxIndex: size}
}

// Dirty sets the bit at the index.
func (b *BitsetMark) Dirty(index uint32) {
	if index >= b.maxIndex {
		return
	}
	if !b.isSet(index) {
		b.dyed()
		b.isDirty = true
		b.bits[index/bitSize] |= 1 << (index % bitSize)
	}
}

// isSet checks if the bit at the index is set.
func (b *BitsetMark) isSet(index uint32) bool {
	return b.bits[index/bitSize]&(1<<(index%bitSize)) != 0
}

// Reset resets the bitset.
func (b *BitsetMark) Reset() {
	clear(b.bits)
	b.isDirty = false
}

// AllBits iterates over the bitset.
func (b *BitsetMark) AllBits() func(func(uint32) bool) {
	return func(f func(uint32) bool) {
		if !b.isDirty {
			return
		}
		for i, bits := range b.bits {
			for j := uint32(0); j < bitSize; j++ {
				if bits&(1<<j) != 0 {
					if !f(uint32(i)*bitSize + j) {
						return
					}
				}
			}
		}
	}
}
