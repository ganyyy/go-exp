package main

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func add1(src []int, add int) {
	for i := 1; i < len(src); i++ {
		src[i] = src[i-1] + add
	}
}

func add2(src []int, add int) {
	var tmp, tmp2 = src[0], 0
	for i := 1; i < len(src); i++ {
		tmp2 = tmp + add
		src[i] = tmp2
		tmp = tmp2
	}
}

func TestSliceAdd(t *testing.T) {
	var src1, src2 = make([]int, 10), make([]int, 10)
	for i := range src1 {
		src1[i] = i
		src2[i] = i
	}

	t.Logf("%+v, %+v", src1, src2)

	add1(src1, 10)
	add2(src2, 10)

	t.Logf("%+v, %+v", src1, src2)

	assert.True(t, reflect.DeepEqual(src1, src2))
}

func BenchmarkSliceAdd(b *testing.B) {
	b.Run("Add1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			add1(make([]int, 10), 10)
		}
	})

	b.Run("Add2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			add2(make([]int, 10), 10)
		}
	})
}

func TestContextWithValue(t *testing.T) {

	var key struct{}
	var valCtx = context.WithValue(context.TODO(), key, "b")
	var bs, _ = json.Marshal(valCtx)

	t.Logf("%v", valCtx.Value(struct{}{}))
	t.Logf("%v", valCtx.Value(key))

	t.Logf("%v", string(bs))

	bs, _ = json.Marshal(context.TODO())
	t.Logf("%v", string(bs))

}

type NullStruct struct {
	a int
	b int
	c int
}

func (m *NullStruct) method(t *testing.T) {
	t.Logf("null struct function %+v", m)
}

func TestNullMethod(t *testing.T) {
	(*NullStruct)(unsafe.Pointer((*bool)(nil))).method(t)
}
