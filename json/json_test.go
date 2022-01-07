package json

import (
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

type T struct {
	M map[int]int
}

func TestUnMarshal(t *testing.T) {
	var tt T
	tt.M = map[int]int{}
	var bs, _ = json.Marshal(tt)

	t.Logf("%v", string(bs))
	tt.M = nil

	err := json.Unmarshal(bs, &tt)

	t.Logf("%v, %+v", err, tt)
	assert.NotNil(t, tt.M)
}

func TestMapSlice(t *testing.T) {
	var src = map[int][]int{
		1: {1, 2, 3, 4, 5},
		2: {4, 5},
		3: {1, 2},
	}
	var dst = make(map[int][]int, len(src))
	for k, v := range src {
		dst[k] = v
	}

	var addDst = make(map[int]int, 3)
	for i := 0; i <= 100; i++ {
		addDst[rand.Intn(10)] += 10
	}

	t.Logf("%+v, %+v", src, dst)
	t.Logf("%+v", addDst)
}
