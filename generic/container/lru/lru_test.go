package lru

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRU(t *testing.T) {
	lru := NewLRU[int, string](3)
	lru.Push(1, "one")
	lru.Push(2, "two")
	lru.Push(3, "three")
	lru.Push(4, "four")
	if _, ok := lru.Get(1); ok {
		t.Error("key 1 should be evicted")
	}
	if _, ok := lru.Get(2); !ok {
		t.Error("key 2 should be in the cache")
	}
	if _, ok := lru.Get(3); !ok {
		t.Error("key 3 should be in the cache")
	}
	if _, ok := lru.Get(4); !ok {
		t.Error("key 4 should be in the cache")
	}

	items := lru.Items()
	assert.Equal(t, []string{"four", "three", "two"}, items)

	lru.Push(2, "two2")
	items = lru.Items()
	assert.Equal(t, []string{"two2", "four", "three"}, items)
}
