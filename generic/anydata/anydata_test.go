package anydata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnyData(t *testing.T) {
	var p = Player{
		Name:    "John",
		Address: "China",
		Age:     18,
	}
	var p2 = Player{}
	bs, err := p.Encode()
	assert.NoError(t, err)
	err = p2.Decode(bs)
	assert.NoError(t, err)
	if p2.Name != p.Name || p2.Address != p.Address || p2.Age != p.Age {
		t.Fatal("anydata test failed")
	}
}
