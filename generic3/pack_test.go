package generic3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPack(t *testing.T) {
	t.Log(Pack(1, 2, 3, 4))
	t.Log(Pack(1.2, 2.3, 3.4, 4.5))
	t.Log(Pack("1", "2", "3", "4"))

	assert.Equal(t, []int{1, 2, 3, 4}, UnpackSet(PackSet(1, 2, 3, 4)))
}
