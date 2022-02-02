package generic2

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
