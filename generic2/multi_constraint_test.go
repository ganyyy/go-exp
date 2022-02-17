package generic2

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddInterfaceFunc(t *testing.T) {
	var a = MyAddNumber(100)
	var b = MyAddNumber(200)
	t.Logf("%#v", AddInterfaceFunc(a, b))
}

func TestNewGraph(t *testing.T) {
	// 注意, 泛型本质上是针对的函数签名
	// 也就是说, 如果 MyNode满足 NodeConstraints, MyEdge满足 EdgeConstraints
	// 那么就可以实例化这个泛型函数函数
	// 核心: 约束和实例化
	var g = NewGraph[MyNode, MyEdge]([]MyNode(nil))

	// 接口也可以实列化
	// 核心: 约束和实例化
	var g2 = NewGraph[NodeInterface, EdgeInterface]([]NodeInterface(nil))

	assert.NotNil(t, g)
	assert.NotNil(t, g2)
}

func TestConvertTo(t *testing.T) {

	var convertStringToInt Convert2Int
	var convertIntToString Convert2String

	assert.Equal(t, ConvertTo[Convert2Int, string, int](convertStringToInt, "2"), 2)

	var src = []string{
		"1",
		"2",
		"3",
		"4",
	}

	var ret = ConvertSlice[Convert2Int, string, int](convertStringToInt, src)
	assert.Equal(t, ret, []int{1, 2, 3, 4})

	assert.Equal(t, src, ConvertSlice[Convert2String, int, string](convertIntToString, ret))
}

func TestConvert3(t *testing.T) {
	var c = MyString("12313")

	assert.Equal(t, Convert3[MyString, int](c), 12313)
}

func TestAnyConvertFunc(t *testing.T) {

	type StructA struct {
		Name string
	}

	type StructB struct {
		Age int
	}

	type StructC struct {
		Address string
	}

	var a2b ConvertFunc[*StructA, *StructB] = func(a *StructA) *StructB {
		return &StructB{Age: rand.Intn(10)}
	}

	var c2a ConvertFunc[*StructC, *StructA] = func(c *StructC) *StructA {
		return &StructA{Name: c.Address}
	}

	var b2c ConvertFunc[*StructB, *StructC] = func(b *StructB) *StructC {
		return &StructC{Address: strconv.Itoa(b.Age)}
	}

	var srcA = make([]*StructA, 10)

	for i := range srcA {
		srcA[i] = &StructA{}
		srcA[i].Name = strconv.Itoa(i)
	}
	t.Logf("%+v", srcA)
	var srcB = AnyConvertFunc[*StructA, *StructB, ConvertFunc[*StructA, *StructB]](a2b, srcA...)
	t.Logf("%+v", srcB)
	var srcC = AnyConvertFunc[*StructB, *StructC, ConvertFunc[*StructB, *StructC]](b2c, srcB...)
	t.Logf("%+v", srcC)
	srcA = AnyConvertFunc[*StructC, *StructA, ConvertFunc[*StructC, *StructA]](c2a, srcC...)
	t.Logf("%+v", srcA)
}
