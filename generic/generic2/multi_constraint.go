package generic2

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

//AddInterface 接口也允许添加泛型约束
type AddInterface[T constraints.Integer] interface {
	constraints.Integer
	Add(T) T
}

type MyAddNumber int

func (m MyAddNumber) Add(t MyAddNumber) MyAddNumber {
	return MyAddNumber(int(m) + int(t))
}

//AddInterfaceFunc 神奇的用法
//关于接口中的类型约束, 必须要在使用时声明
func AddInterfaceFunc[ET constraints.Integer, AT AddInterface[ET]](a AT, b ET) ET {
	return a.Add(b)
}

/**
=========== 复杂的多引用 ===========
*/

type NodeConstraint[Edge any] interface {
	Edges() []Edge
}

type EdgeConstraint[Node any] interface {
	Nodes() (from, to Node)
}

/**
结构体实例化约束
*/

type MyNode struct{}

func (n MyNode) Edges() []MyEdge { return nil }

type MyEdge struct{}

func (e MyEdge) Nodes() (from, to MyNode) { return MyNode{}, MyNode{} }

/**
接口实例化泛型约束
*/

type NodeInterface interface {
	Edges() []EdgeInterface
}

type EdgeInterface interface {
	Nodes() (NodeInterface, NodeInterface)
}

type NodeII struct{}

func (n NodeII) Edges() []EdgeInterface { return nil }

type EdgeII struct{}

func (e EdgeII) Nodes() (NodeInterface, NodeInterface) { return NodeII{}, NodeII{} }

/**
图的定义
*/

type Graph[Node NodeConstraint[Edge], Edge EdgeConstraint[Node]] struct {
}

func NewGraph[Node NodeConstraint[Edge], Edge EdgeConstraint[Node]](_ []Node) *Graph[Node, Edge] {
	var graph Graph[Node, Edge]
	return &graph
}

type AnyConvert[From any, To any] interface {
	Convert(From) To
}

type Convert2Int struct{}

func (c Convert2Int) Convert(s string) int {
	var v, _ = strconv.Atoi(s)
	return v
}

type Convert2String struct{}

func (c Convert2String) Convert(i int) string {
	return strconv.Itoa(i)
}

func ConvertTo[CT AnyConvert[From, To], From any, To any](c CT, v From) To {
	return c.Convert(v)
}

func ConvertSlice[CT AnyConvert[From, To], From any, To any](c CT, src []From) []To {
	var ret = make([]To, len(src))
	for i, v := range src {
		ret[i] = c.Convert(v)
	}
	return ret
}

type ConvertToOther[To any] interface {
	Convert() To
}

func (m MyString) Convert() int {
	var v, _ = strconv.Atoi(string(m))
	return v
}

func Convert3[T ConvertToOther[To], To any](v T) To {
	return v.Convert()
}

type ConvertFunc[From, To any] func(From) To

func AnyConvertFunc[From any, To any, F ConvertFunc[From, To]](f F, ele ...From) []To {
	var ret = make([]To, 0, len(ele))
	for _, v := range ele {
		ret = append(ret, f(v))
	}
	return ret
}
