package generic2

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvert(t *testing.T) {
	var a = 10
	var b = 10.0

	assert.Equal(t, a, Convert[float64, int](b))
	assert.Equal(t, b, Convert[int, float64](a))
}

func TestAddNum(t *testing.T) {
	t.Logf("%v", AddNum[int8](10))
}

func TestAddString(t *testing.T) {
	assert.Equal(t, AddString("a"), "aa")
	//assert.Equal(t, AddString[int](10), 20) // 这个编译不通过
}

//func TestStringFish(t *testing.T) {
//
//	type Stringer interface {
//		String() string
//	}
//
//	RunStringFish("")
//	RunStringFish(MyString(""))
//	RunStringFish(Stringer(nil))
//	RunStringFish(10) // 编译不通过
//}
