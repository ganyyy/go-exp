package meta

import (
	"iter"
	"unsafe"
)

type Mark[T any] struct {
	dirtyMark                    // dyeing is a function that will be called when any field is set
	dirty     map[uint16]func(T) // dirty is a map that stores the dirty mark of each field
}

// Dirty sets the dirty mark of the field.
func (m *Mark[T]) Dirty(idx uint16, set func(T)) {
	if m.dirty == nil {
		m.dirty = make(map[uint16]func(T))
	}
	if _, ok := m.dirty[idx]; !ok {
		m.dirty[idx] = set
	}
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

type IMark interface {
	mark(uint32)                 // mark self as dirty
	setMark(IMark, uint32, bool) // set root and idx
}

type dirtyMark struct {
	idx     uint32 // idx is the index of the this field
	isDirty bool   // isDirty is the flag of the dirty mark
	root    IMark  // dyeing is a function that will be called when any field is set
}

// dirty marks the field as dirty.
func (m *dirtyMark) dirty() {
	m.isDirty = true
	if m.root != nil {
		// Mark the parent field as dirty.
		m.root.mark(m.idx)
	}
}

// reset resets the dirty mark.
func (m *dirtyMark) reset() {
	m.isDirty = false
}

// setMark sets the mark.
func (m *dirtyMark) setMark(mark IMark, idx uint32, needNew bool) {
	if needNew && m.root != nil {
		panic("mark reset")
	}
	if mark == nil || m.root == nil {
		m.isDirty = false
		m.root = mark
		m.idx = idx
		return
	}
	if !Equal[IMark](m.root, mark) || m.idx != idx {
		panic("mark should not be changed")
	}
}

// mark marks the field as dirty. default implementation. direct call Dirty.
func (m *dirtyMark) mark(uint32) {
	m.dirty()
}

type BitType uint8

const (
	bitSize = uint32(unsafe.Sizeof(BitType(0)) * 8)
)

type BitsetMark struct {
	dirtyMark

	maxIdx uint32
	bits   []BitType
}

func NewBitsetMark(size uint32) *BitsetMark {
	return &BitsetMark{bits: make([]BitType, (size+bitSize-1)/bitSize), maxIdx: size}
}

func (b *BitsetMark) mark(index uint32) {
	if index >= b.maxIdx {
		return
	}
	if !b.isSet(index) {
		b.dirty()
		b.bits[index/bitSize] |= 1 << (index % bitSize)
	}
}

// isSet checks if the bit at the index is set.
func (b *BitsetMark) isSet(index uint32) bool {
	return b.bits[index/bitSize]&(1<<(index%bitSize)) != 0
}

// reset resets the bitset.
func (b *BitsetMark) reset() {
	clear(b.bits)
	b.dirtyMark.reset()
}

func (b *BitsetMark) dirtyBits() iter.Seq[uint32] {
	return func(f func(uint32) bool) {
		if !b.isDirty {
			return
		}
		for i, bits := range b.bits {
			for j := uint32(0); j < bitSize; j++ {
				if bits&(1<<j) != 0 {
					idx := uint32(i)*bitSize + j
					if idx >= b.maxIdx {
						return
					}
					if !f(idx) {
						return
					}
				}
			}
		}
	}
}

func SetMarkHelper(m IMark, root IMark, idx uint32) {
	m.setMark(root, idx, false)
}

func MarkHelper(m IMark, idx uint32) {
	m.mark(idx)
}

type IReset interface{ reset() }
type IDirtyBits interface{ dirtyBits() iter.Seq[uint32] }

func ResetHelper(m IReset) {
	m.reset()
}

func DirtyBitsHelper(m IDirtyBits) iter.Seq[uint32] {
	return m.dirtyBits()
}
