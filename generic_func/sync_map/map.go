package sync_map

import (
	"sync"
)

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{}
}

type SyncMap[K comparable, V any] struct {
	m sync.Map
}

func (s *SyncMap[K, V]) Load(key K) (V, bool) {
	var val, ok = s.m.Load(key)
	if !ok {
		var v V
		return v, false
	}
	return val.(V), true
}

func (s *SyncMap[K, V]) Store(key K, value V) {
	s.m.Store(key, value)
}

func (s *SyncMap[K, V]) LoadOrStore(key K, value V) (V, bool) {
	var actual, ok = s.m.LoadOrStore(key, value)
	return actual.(V), ok
}

func (s *SyncMap[K, V]) LoadAndDelete(key K) (V, bool) {
	var v, loaded = s.m.LoadAndDelete(key)
	if !loaded {
		var v V
		return v, loaded
	}
	return v.(V), loaded
}

func (s *SyncMap[K, V]) Delete(key K) {
	s.m.Delete(key)
}

func (s *SyncMap[K, V]) Range(f func(K, V) bool) {
	var fun = func(key, value any) bool {
		return f(key.(K), value.(V))
	}
	s.m.Range(fun)
}
