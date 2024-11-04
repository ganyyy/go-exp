package readonly

import (
	"iter"
	"sync"
)

func OnceMap[K comparable, V any](src map[K]V) func() Map[K, V] {
	return sync.OnceValue(func() Map[K, V] {
		return NewMap(src)
	})
}

func NewMap[K comparable, V any](p map[K]V) Map[K, V] {
	if len(p) == 0 {
		return Map[K, V]{inner: nil}
	}
	m := p
	return Map[K, V]{inner: m}
}

func OnceMapFrom[K comparable, V any, RV any](src map[K]V, to func(V) RV) func() Map[K, RV] {
	return sync.OnceValue(func() Map[K, RV] {
		return NewMapFrom(src, to)
	})
}

func NewMapFrom[K comparable, V any, RV any](p map[K]V, to func(V) RV) Map[K, RV] {
	if len(p) == 0 {
		return NewMap[K, RV](nil)
	}
	m := p
	var result = make(map[K]RV, len(m))
	for k, v := range m {
		result[k] = to(v)
	}
	return NewMap(result)
}

type Map[K comparable, V any] struct {
	inner map[K]V
}

func (x *Map[K, V]) Get(k K) (V, bool) {
	v, ok := x.inner[k]
	return v, ok
}

func (x *Map[K, V]) Len() int {
	return len(x.inner)
}

func (x *Map[K, V]) Range() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range x.inner {
			if !yield(k, v) {
				break
			}
		}
	}
}
