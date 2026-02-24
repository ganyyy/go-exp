package api

import "testing"

func TestAPI(t *testing.T) {
	var p = new(100)
	t.Logf("p: %v, %p", *p, p)
}

type Addable[T Addable[T]] interface {
	Add(other T) T
	Sub(other T) T
	Zero() T
}

type Addable2[T any] interface {
	Add(Addable2[T]) Addable2[T]
	Sub(Addable2[T]) Addable2[T]
	Zero() Addable2[T]
}

type MyInt int

func (a MyInt) Add(other MyInt) MyInt {
	return a + other
}

func (a MyInt) Sub(other MyInt) MyInt {
	return a - other
}

func (a MyInt) Zero() MyInt {
	return 0
}

func testMyInt[T Addable[T]](a, b T) T {
	return a.Zero().Add(b).Sub(a)
}

func testMyInt2[T Addable2[T]](a, b T) Addable2[T] {
	return a.Zero().Add(b).Sub(a)
}

func TestMyInt(t *testing.T) {
	var a, b MyInt = 5, 3
	result := testMyInt(a, b)
	// result2 := testMyInt2(a, b)
	t.Logf("result: %v", result)
}
