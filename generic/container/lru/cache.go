package lru

import (
	"container/list"
)

type LRUCache[V any] struct {
	list list.List
	keys map[interface{}]*list.Element
}

func New[V any](capicity int) *LRUCache[V] {
	return &LRUCache[V]{
		keys: make(map[interface{}]*list.Element, capicity),
	}
}

func (l *LRUCache[V]) Push(v V) {
	l.keys[interface{}(v)] = l.list.PushFront(v)
}
