package api

import (
	"iter"
	"strings"
	"testing"
)

func ReverseIter(n int) func(func(int) bool) {
	return func(f func(int) bool) {
		for i := n; i >= 0; i-- {
			if !f(i) {
				break
			}
		}
	}
}

func Reverse2Iter[T any](iter iter.Seq[T]) iter.Seq2[int, T] {
	return func(f func(int, T) bool) {
		var i int
		for v := range iter {
			if !f(i, v) {
				break
			}
			i++
		}
	}
}

func TestIter(t *testing.T) {
	var seq iter.Seq[int] = ReverseIter(10)
	for i := range seq {
		t.Logf("i: %d", i)
	}

	var seq2 iter.Seq2[int, int] = Reverse2Iter(ReverseIter(20))
	for i, v := range seq2 {
		t.Logf("i: %d, v: %d", i, v)
	}
}

func Split(s string, sep string) iter.Seq[string] {
	return func(f func(string) bool) {
		if len(sep) == 0 {
			for i := 0; i < len(s); i++ {
				if !f(s[i : i+1]) {
					break
				}
			}
			return
		}
		for len(s) != 0 {
			i := strings.Index(s, sep)
			if i < 0 {
				if !f(s) {
					break
				}
				break
			}
			if !f(s[:i]) {
				break
			}
			s = s[i+len(sep):]
		}
	}
}

func TestSplit(t *testing.T) {
	var testCase = []struct {
		s   string
		sep string
	}{
		{"a,b,c", ","},
		{"a,b,c,", ","},
		{"a,b,c", ""},
		{",a,b,c", ","},
		{"a,b,c", "b"},
	}

	for _, tc := range testCase {
		t.Logf("s: %s, sep: %s", tc.s, tc.sep)
		var seq iter.Seq[string] = Split(tc.s, tc.sep)
		for i := range seq {
			t.Logf("i: %s", i)
		}
	}
}

func BenchmarkIter(b *testing.B) {
	b.Run("Normal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var seq iter.Seq[int] = ReverseIter(100)
			for range seq {
			}
		}
		b.ReportAllocs()
	})

	b.Run("Seq2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var seq2 iter.Seq2[int, int] = Reverse2Iter(ReverseIter(100))
			for range seq2 {
			}
		}
		b.ReportAllocs()
	})

	b.Run("Simple", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 100; j >= 0; j-- {
			}
		}
		b.ReportAllocs()
	})
}

func BenchmarkMapIter(b *testing.B) {
	var m = Map[int, int]{
		m: map[int]int{},
	}
	for i := range 100 {
		m.m[i] = i
	}

	b.Run("Keys", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for range m.Keys() {
			}
		}
		b.ReportAllocs()
	})

	b.Run("Values", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for range m.Values() {
			}
		}
		b.ReportAllocs()
	})

	b.Run("NormalKeys", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for range m.m {
			}
		}
		b.ReportAllocs()
	})
}

type Map[K comparable, V any] struct {
	m map[K]V
}

// Keys
func (m Map[K, V]) Keys() iter.Seq[K] {
	return func(f func(K) bool) {
		for k := range m.m {
			if !f(k) {
				break
			}
		}
	}
}

// Values
func (m Map[K, V]) Values() iter.Seq[V] {
	return func(f func(V) bool) {
		for _, v := range m.m {
			if !f(v) {
				break
			}
		}
	}
}

// Iterator for map
func (m Map[K, V]) Iterator() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m.m {
			if !yield(k, v) {
				break
			}
		}
	}
}

func TestPull(t *testing.T) {
	var m = Map[int, int]{
		m: map[int]int{},
	}

	for i := range 100 {
		m.m[i] = i
	}

	next, stop := iter.Pull2(m.Iterator())
	defer stop()
	for {
		k, v, ok := next()
		if !ok {
			break
		}
		t.Logf("k: %d, v: %d", k, v)
	}
}
