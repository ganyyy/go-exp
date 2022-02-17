package generic2

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDouble(t *testing.T) {

	t.Run("Double", func(t *testing.T) {
		var src = []int{1, 2, 3, 4}
		t.Logf("%#v", Double(src))
	})

	t.Run("DoubleDefined", func(t *testing.T) {
		type MySlice []int

		var src MySlice = []int{1, 2, 3, 4}

		t.Logf("%#v", Double(src))
		// 不同的约束会推导出不同的结果
		t.Logf("%#v", DoubleDefined(src)) // IDE的Bug, 还是推断有问题
		t.Logf("%#v", DoubleDefined2(src))
	})

}

func TestFromStrings(t *testing.T) {
	t.Run("Panic", func(t *testing.T) {
		defer func() {
			var err = recover()
			t.Logf("panic:%v", err)
			if err != nil {
				t.Logf("%s", string(debug.Stack()))
			}
		}()

		var nums = FromStrings[*Settable]([]string{"1", "2"})
		t.Logf("%+v", nums)
	})

	t.Run("Good", func(t *testing.T) {
		// 约束推导可以让我们没必要写第二个参数
		// 第二个约束类型 Setter2[T] 限定了 *T 必须要实现 Set(string).
		// 不满足条件的话无法编译
		// T -> Setter
		// PT -> *T
		var nums = FromStrings2[Settable]([]string{"1", "2"})
		var nums2 = FromStrings2[Settable]([]string{"1", "2"})
		t.Logf("%#v, %#v", nums, nums2)
	})
}

func TestIndexEqual(t *testing.T) {
	var src = []EqualNum{
		1, 2, 3, 4,
	}

	assert.Equal(t, IndexEqual(src, EqualNum(1)), 0)
	assert.Equal(t, IndexEqual2(src, EqualNum(1)), 0)
	assert.Equal(t, IndexEqual(src, EqualNum(10)), -1)
	assert.Equal(t, IndexEqual2(src, EqualNum(10)), -1)
}
