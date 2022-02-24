package generic2

import "sync"

type AnyMap[K comparable, V any] struct {
	mutex sync.RWMutex
	m     map[K]V
}

func NewAnyMap[K comparable, V any](capacity int) *AnyMap[K, V] {
	return &AnyMap[K, V]{
		m: make(map[K]V, capacity),
	}
}

func (m *AnyMap[K, V]) Add(k K, v V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.m == nil {
		m.m = map[K]V{}
	}
	m.m[k] = v
}

func (m *AnyMap[K, V]) Get(k K) (V, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	var v, ok = m.m[k]
	return v, ok
}

func (m *AnyMap[K, V]) Del(k K) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.m, k)
}

func (m *AnyMap[K, V]) Range(f func(K, V)) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for k, v := range m.m {
		f(k, v)
	}
}

func (m *AnyMap[K, V]) Count() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.m)
}
