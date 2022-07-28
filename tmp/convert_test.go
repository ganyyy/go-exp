package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Interface1 interface {
	Add()
}

type Interface2 interface {
	Interface1
	Sub()
}

func convert() {
	var a any = 1
	var b any = 1.0
	var c any = true
	var d any = "123"

	var i Interface2
	var i1 any = i

	_ = a.(int)
	_ = b.(float64)
	_ = c.(bool)
	_ = d.(string)
	_ = i1.(Interface2)
	_ = i.(Interface1)
}

func TestAssertI2I(t *testing.T) {
	var v int
	var a interface{} = &v

	*(a.(*int)) = 100

	assert.Equal(t, v, 100)
}

func TestHTTP2(t *testing.T) {
	HTTP2()
}

func TestZigzag(t *testing.T) {
	var a int64 = -2
	var zigA = (a << 1) ^ (a >> 63)
	t.Log(a, zigA, (a << 1), (a >> 63))

	// 0000 0100 +
	// 1111 1100 = 1 0000 0000
	// 0000 0001 +
	// 1111 1111 = 1 0000 0000

	// 1111 1111 ^
	// 1111 1100 = 0000 0011
}
