package rank

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestRank(t *testing.T) {
	r := NewNormalRank[string, *RankNode]()

	for i := 0; i < 100; i++ {
		r.Add(&RankNode{
			id:    strconv.Itoa(rand.Intn(100)),
			score: rand.Intn(100),
		})
	}

	t.Logf("rank length %v", r.Len())

	for _, node := range r.nodes {
		t.Logf("node %v", node)
	}
}
