package rank

import "slices"

type IRankNode[T any, K comparable] interface {
	Id() K
	Less(T) bool
}

type RankNode struct {
	id    string
	score int
}

func (r *RankNode) Less(other *RankNode) bool {
	return r.score < other.score
}

func (r *RankNode) Id() string {
	return r.id
}

var _ IRankNode[*RankNode, string] = (*RankNode)(nil)

type NormalRank[K comparable, T IRankNode[T, K]] struct {
	m     map[K]T
	nodes []T
}

func (n *NormalRank[K, T]) Sort() {
	slices.SortFunc(n.nodes, func(i, j T) int {
		if i.Id() == j.Id() {
			return 0
		}
		if i.Less(j) {
			return -1
		}
		return 1
	})
}

func (n *NormalRank[K, T]) Len() int {
	return len(n.nodes)
}

func (n *NormalRank[K, T]) Less(i, j int) bool {
	return n.nodes[i].Less(n.nodes[j])
}

func (n *NormalRank[K, T]) Swap(i, j int) {
	n.nodes[i], n.nodes[j] = n.nodes[j], n.nodes[i]
}

func (n *NormalRank[K, T]) Get(id K) (T, bool) {
	node, ok := n.m[id]
	return node, ok
}

// Rank returns the rank of the given node.
func (n *NormalRank[K, T]) Rank(node T) int {
	idx, ok := n.Search(node)
	if !ok {
		return -1
	}
	return idx + 1
}

// Search returns the index of the first node in the rank whose score is greater than or equal to the given score.
func (n *NormalRank[K, T]) Search(node T) (int, bool) {
	return slices.BinarySearchFunc(n.nodes, node, func(a, b T) int {
		if a.Id() == b.Id() {
			return 0
		}
		if a.Less(b) {
			return -1
		}
		return 1
	})
}

// Remove removes a node from the rank.
func (n *NormalRank[K, T]) Remove(id K) {
	old, ok := n.Get(id)
	if !ok {
		return
	}
	delete(n.m, id)
	idx, ok := n.Search(old)
	if !ok {
		return
	}
	n.nodes = append(n.nodes[:idx], n.nodes[idx+1:]...)
}

// Add adds a node to the rank.
func (n *NormalRank[K, T]) Add(node T) {
	if old, ok := n.m[node.Id()]; ok {
		n.Remove(old.Id())
	}
	n.m[node.Id()] = node
	insertIdx, _ := n.Search(node)
	n.nodes = append(n.nodes, node)
	copy(n.nodes[insertIdx+1:], n.nodes[insertIdx:])
	n.nodes[insertIdx] = node
}

var _ = NormalRank[string, *RankNode]{}

func NewNormalRank[K comparable, T IRankNode[T, K]]() *NormalRank[K, T] {
	return &NormalRank[K, T]{
		m: make(map[K]T),
	}
}
