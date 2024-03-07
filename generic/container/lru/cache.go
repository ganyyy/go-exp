package lru

import "container/list"

type Node[K comparable, T any] struct {
	key   K
	value T
}

// Key return the key of the node
func (n *Node[K, V]) Key() K { return n.key }

// Value return the value of the node
func (n *Node[K, V]) Value() V { return n.value }

// Set the value of the node
func (n *Node[K, V]) Set(value V) { n.value = value }

type LRU[K comparable, T any] struct {
	capacity int
	cache    map[K]*list.Element
	items    *list.List
}

// NewLRU create a new LRU cache
func NewLRU[K comparable, V any](capacity int) *LRU[K, V] {
	return &LRU[K, V]{
		capacity: capacity,
		cache:    make(map[K]*list.Element, capacity),
		items:    list.New(),
	}
}

// Get the value of the key
func (l *LRU[K, V]) Get(key K) (V, bool) {
	if ele, ok := l.cache[key]; ok {
		l.items.MoveToFront(ele)
		return l.node(ele.Value).Value(), true
	}
	var zero V
	return zero, false
}

// Push a new key-value pair into the cache
func (l *LRU[K, V]) Push(key K, value V) {
	if ele, ok := l.cache[key]; ok {
		l.items.MoveToFront(ele)
		l.node(ele.Value).Set(value)
		return
	}
	if l.items.Len() >= l.capacity {
		delete(l.cache, l.node(l.items.Remove(l.items.Back())).Key())
	}
	ele := l.items.PushFront(&Node[K, V]{key, value})
	l.cache[key] = ele
}

// Items return all the items in the cache order by the most recently used
func (l *LRU[K, V]) Items() []V {
	items := make([]V, 0, l.items.Len())
	for ele := l.items.Front(); ele != nil; ele = ele.Next() {
		items = append(items, l.node(ele.Value).Value())
	}
	return items
}

// node return the node of the element
func (l *LRU[K, V]) node(v any) *Node[K, V] {
	if val, ok := v.(*Node[K, V]); ok {
		return val
	}
	return &Node[K, V]{}
}
