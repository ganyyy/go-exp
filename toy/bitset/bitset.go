package bitset

import (
	"iter"
	"math/bits"
)

type layer []uint

func (l layer) len() uint {
	return uint(len(l))
}

func (l *layer) resize(n uint, v uint) {
	if n <= l.len() {
		return
	}
	tmp := make([]uint, n)
	// 复制
	copy(tmp, *l)
	// 填充
	if v != DefaultVal {
		for i := l.len(); i < n; i++ {
			tmp[i] = v
		}
	}
	*l = tmp
}

func (l *layer) fillUp(n uint, v uint) {
	n++
	if l.len() >= n {
		return
	}
	l.resize(n, v)
}

func (l layer) set(index, mask uint) {
	if index >= l.len() {
		return
	}
	l[index] |= mask
}

func (l layer) check(index, mask uint) bool {
	if index >= l.len() {
		return false
	}
	return l[index]&mask != 0
}

func (l layer) remove(index, mask uint) uint {
	if index >= l.len() {
		return 0
	}
	l[index] &^= mask
	return l[index]
}

type BitSet struct {
	/*
		逐级递减
	*/

	layer3 uint
	layer2 layer
	layer1 layer
	layer0 layer
}

func (b *BitSet) extend(idx Index) {
	validRange(idx)
	var p0, p1, p2 = Offsets(idx)
	b.layer2.fillUp(p2, DefaultVal)
	b.layer1.fillUp(p1, DefaultVal)
	b.layer0.fillUp(p0, DefaultVal)
}

func (b *BitSet) addSlow(id Index) {
	var _, p1, p2 = Offsets(id)
	b.layer1.set(p1, id.Mask(Shift1))
	b.layer2.set(p2, id.Mask(Shift2))
	b.layer3 |= id.Mask(Shift3)
}

// Add 添加Index到集合中
func (b *BitSet) Add(idx Index) bool {
	var p0, mask = idx.Offset(Shift1), idx.Mask(Shift0)
	if p0 >= b.layer0.len() {
		b.extend(idx)
	}
	old := b.layer0[p0]
	if old&mask != 0 {
		return true
	}
	b.layer0.set(p0, mask)
	if old == 0 {
		// 当前位置首次增加元素, 需要同步到上层
		b.addSlow(idx)
	}
	return false
}

// Remove 从集合中移除指定元素
func (b *BitSet) Remove(idx Index) (has bool) {
	var p0, p1, p2 = Offsets(idx)
	if p0 >= b.layer0.len() {
		return
	}
	if !b.layer0.check(p0, idx.Mask(Shift0)) {
		return
	}
	has = true
	nv := b.layer0.remove(p0, idx.Mask(Shift0))
	if nv == 0 {
		return
	}
	nv = b.layer1.remove(p1, idx.Mask(Shift1))
	if nv == 0 {
		return
	}
	nv = b.layer2.remove(p2, idx.Mask(Shift2))
	if nv == 0 {
		return
	}
	b.layer3 &^= idx.Mask(Shift3)
	return
}

func (b *BitSet) Contain(idx Index) bool {
	var p0, mask = idx.Offset(Shift1), idx.Mask(Shift0)
	if p0 >= b.layer0.len() {
		return false
	}
	return b.layer0.check(p0, mask)
}

func (b *BitSet) Clear() {
	b.layer3 = 0
	b.layer2 = nil
	b.layer1 = nil
	b.layer0 = nil
}

func (b *BitSet) Range() iter.Seq[Index] {
	return func(yield func(Index) bool) {

		// 每一个level对应的掩码
		var masks = [Layers]uint{3: b.layer3}
		// 1
		// 1--------1--------1
		// 1--1--1--1--1--1--1--1--1
		// 111111111111111111111111111
		var prefix [Layers - 1]uint32

		type State int32

		const (
			Empty State = iota
			Continue
			Value
		)

		getPrefix := func(level uint) uint32 {
			if level >= Layers-1 {
				return 0
			}
			return prefix[level]
		}

		var idx Index

		handleLevel := func(level uint) State {
			mask := masks[level]
			if mask == 0 {
				return Empty
			}
			// 首个1出现的位置
			firstBit := bits.TrailingZeros(mask)
			// 清理该位
			masks[level] &^= (1 << firstBit)
			// 获取该位对应的值
			idx = Index(getPrefix(level) | uint32(firstBit))
			if level == 0 {
				return Value
			}
			// 设置缓存
			masks[level-1] = b.getFromLevel(level-1, uint(idx))
			// 更新前缀. 类似于一种前缀树的概念
			prefix[level-1] = uint32(idx) << Bits
			return Continue
		}
	loop:
		for level := uint(0); level < Layers; level++ {
			state := handleLevel(level)
			if state == Empty {
				continue
			}
			if state == Continue {
				goto loop
			}
			// Value
			if yield(idx) {
				goto loop
			}
			return
		}
	}
}

func (b *BitSet) getFromLevel(level, idx uint) uint {
	switch level {
	case 0:
		return b.layer0[idx]
	case 1:
		return b.layer1[idx]
	case 2:
		return b.layer2[idx]
	case 3:
		return b.layer3
	}
	panic("invalid level")
}
