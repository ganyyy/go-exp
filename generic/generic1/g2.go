package main

import (
	"fmt"
	"strconv"
)

func pack[T any](params ...T) []T {
	return params
}

type Student struct {
}

type Addable interface {
	Student
	ToString() string
}

func GetOrDefault[T any](src interface{}, def T) T {
	if t, ok := src.(T); ok {
		return t
	}
	return def
}

type set[T comparable] map[T]struct{}

func (s set[T]) pack(params ...T) {
	for _, p := range params {
		s[p] = struct{}{}
	}
}

func (s set[T]) unpack() []T {
	var ret = make([]T, 0, len(map[T]struct{}(s)))
	for k := range s {
		ret = append(ret, k)
	}
	return ret
}

func packSet[T comparable](params ...T) set[T] {
	var m set[T] = make(map[T]struct{}, len(params))
	for _, p := range params {
		m[p] = struct{}{}
	}
	return m
}

type myInt int

type Addable2 interface {
	~int | ~string
	String() string
}

func (a myInt) String() string {
	return strconv.Itoa(int(a))
}

type myString string

func (s myString) String() string {
	return string(s)
}

func Map[T any](src []T, opt func(T) T) []T {
	var ret = make([]T, len(src))
	for i, v := range src {
		ret[i] = opt(v)
	}
	return ret
}

func Add[T, T2 Addable2](a T2, b T) string {
	return a.String() + b.String()
}

func Add2(a, b fmt.Stringer) string {
	return a.String() + b.String()
}

type Comparable[T comparable] interface {
	ComparableTo(T) int
}

type MyInt int

func (i MyInt) ComparableTo(i2 MyInt) int {
	return int(i - i2)
}

type ID interface {
	ID() string
}

type AID struct {
}

func (a *AID) ID() string {
	return ""
}

type IDPointer[T any] interface {
	ID
	*T
}

func GetID[T any, PT IDPointer[T]](t *T) string {
	return PT(t).ID()
}

func GetIDInterface[T ID](t T) string {
	return t.ID()
}

//go:noinline
func GetIDDirect(t *AID) string {
	return t.ID()
}
