package tool

import (
	"testing"
)

func TestWaitGroup(t *testing.T) {
	var wait WaitGroup

	wait.Do(func() {})

	wait.Wait()
}
