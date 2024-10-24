package readonly

import (
	"iter"
	"sync"
)

func OnceList[T any](src *[]T) func() List[T] {
	return sync.OnceValue(func() List[T] {
		return NewList(src)
	})
}

func NewList[T any, S ~[]T](p *S) List[T] {
	if p == nil || len(*p) == 0 {
		return List[T]{inner: nil}
	}
	inner := *p
	*p = nil
	return List[T]{inner: inner}
}

func OnceListFrom[T any, RT any](src *[]T, to func(*T) RT) func() List[RT] {
	return sync.OnceValue(func() List[RT] {
		return NewListFrom(src, to)
	})
}

func NewListFrom[T any, RT any, S ~[]T](p *S, to func(*T) RT) List[RT] {
	if p == nil || len(*p) == 0 {
		return NewList[RT, []RT](nil)
	}
	inner := *p
	*p = nil
	var result = make([]RT, 0, len(inner))
	for _, v := range inner {
		result = append(result, to(&v))
	}
	return NewList(&result)
}

type List[T any] struct {
	inner []T
}

func (x *List[T]) Get(i int) (T, bool) {
	if i < 0 || i >= len(x.inner) {
		var empty T
		return empty, false
	}
	return x.inner[i], true
}

func (x *List[T]) Len() int {
	return len(x.inner)
}

func (x *List[T]) Range() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range x.inner {
			if !yield(v) {
				break
			}
		}
	}
}
