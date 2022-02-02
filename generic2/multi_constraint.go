package generic2

import "constraints"

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
