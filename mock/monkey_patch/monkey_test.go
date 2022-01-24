package monkey_patch

import (
	"testing"

	"github.com/go-kiss/monkey"
	"github.com/stretchr/testify/assert"
)

//go:noinline
func MonkeyAdd(a, b int) int {
	return a * b
}

func TestMonkeyAdd(t *testing.T) {
	var a, b = 100, 200

	assert.Equal(t, a+b, Add(a, b))

	var patch = monkey.Patch(Add, MonkeyAdd)
	defer patch.Unpatch()

	assert.Equal(t, MonkeyAdd(a, b), Add(a, b))
}
