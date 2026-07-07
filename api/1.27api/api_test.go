package api

import (
	"iter"
	"slices"
	"strconv"
	"testing"
	"uuid"

	"github.com/stretchr/testify/require"
)

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person  // embedded field
	Company string
}

func TestEmbeddedStruct(t *testing.T) {
	var _ = Employee{
		Name:    "John Doe",
		Age:     30,
		Company: "Acme Corp",
	}
}

type Array[T any] []T

func (a Array[T]) Map[U any](fn func(T) U) Array[U] {
	result := make(Array[U], len(a))
	for i, v := range a {
		result[i] = fn(v)
	}
	return result
}

func (a Array[T]) MapIter[U any](fn func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for _, v := range a {
			if !yield(fn(v)) {
				break
			}
		}
	}
}

func (a Array[T]) Filter(fn func(T) bool) Array[T] {
	result := make(Array[T], 0)
	for _, v := range a {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func (a Array[T]) Reduce(fn func(T, T) T) T {
	if len(a) == 0 {
		var zero T
		return zero
	}
	result := a[0]
	for i := 1; i < len(a); i++ {
		result = fn(result, a[i])
	}
	return result
}

func TestArray(t *testing.T) {
	var src Array[int] = Array[int]{1, 2, 3, 4, 5}
	var dst Array[string] = src.Map(strconv.Itoa)
	var dst2 Array[int] = dst.Map(func(s string) int {
		val, _ := strconv.Atoi(s)
		return val
	})

	require.Equal(t, src, dst2)

	var ret = src.Reduce(func(i1, i2 int) int { return i1 + i2 })
	require.Equal(t, 15, ret)

	var dst3 = slices.Collect(src.MapIter(strconv.Itoa))
	var dst4 = slices.Collect(dst.MapIter(func(s string) string { return s }))
	require.Equal(t, dst3, dst4)
}

func TestUUID(t *testing.T) {
	t.Logf("v4: %s", uuid.NewV4().String())
	t.Logf("v7: %s", uuid.NewV7().String())
}
