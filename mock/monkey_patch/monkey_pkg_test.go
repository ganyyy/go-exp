package monkey_patch_test

import (
	"testing"

	"github.com/go-kiss/monkey"
	"github.com/stretchr/testify/assert"

	"ganyyy.com/go-exp/mock/monkey_patch"
)

//go:noinline
func MonkeyAdd(a, b int) int {
	return a - b
}

func TestMonkeyAdd(t *testing.T) {
	var a, b = 100, 200

	assert.Equal(t, a+b, monkey_patch.Add(a, b))

	{
		var patch = monkey.Patch(monkey_patch.Add, MonkeyAdd)
		assert.Equal(t, monkey_patch.Add(a, b), MonkeyAdd(a, b))
		patch.Unpatch()
	}

	assert.Equal(t, a+b, monkey_patch.Add(a, b))
}
