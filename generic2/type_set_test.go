package generic2

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSmall(t *testing.T) {
	var src = []int{1, 2, 3, 4, 5}

	t.Logf("%+v", Small(src))
}

func TestIndex2(t *testing.T) {
	var src = []int{1, 2, 3, 4, 5}
	assert.Equal(t, Index(src, 10), -1)
	assert.Equal(t, Index(src, 5), 4)
}
