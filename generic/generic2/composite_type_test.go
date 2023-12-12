package generic2

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	var src = []string{"1", "2", "3"}

	t.Logf("%+v", Join(src, ","))

	var src2 = [][]byte{{'1'}, {'2'}, {'3'}}
	t.Logf("%s", Join(src2, []byte(",")))
}

func TestEntry(t *testing.T) {
	type E = []int
	var a E = []int{1, 2, 3, 4}
	var b = map[int]int{
		1: 1,
		2: 2,
		3: 3,
	}

	assert.Equal(t, Entry(a, 1), Entry(b, 2))
}

func TestMapSlice(t *testing.T) {
	type MyInt int
	type MyString string

	var src = make([]MyInt, 5)
	var src2 = make([]MyString, 5)
	for i := range src {
		src[i] = MyInt(i)
		src2[i] = MyString(strconv.Itoa(i))
	}

	var ret = MapSlice(src, func(e MyInt) MyInt {
		return MyInt(int(e) + 10)
	})
	t.Logf("%#v", ret)

	var ret2 = MapSlice(src2, func(e MyString) MyString {
		return MyString(string(e) + "aab")
	})
	t.Logf("%#v", ret2)
}

func TestIncrementX(t *testing.T) {
	var a struct {
		a int
		x int
	}

	var b = struct {
		b int
		x float64
	}{b: 10, x: 200}

	var c struct {
		c int
		x bool
	}

	var d struct {
		x int
	}

	IncrementX(&a)
	IncrementX(&b)
	IncrementX(&c)
	//IncrementX(&d) // 这个不满足约束, 所以无法使用

	t.Logf("a:%+v, b:%+v, c:%+v, d:%+v", a, b, c, d)
}
