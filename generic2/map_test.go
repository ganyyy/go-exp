package generic2

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnyMap(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		const (
			NUM   = 10
			MULTI = 100
		)
		var wait sync.WaitGroup
		wait.Add(NUM)
		var m = NewAnyMap[int, int](MULTI * NUM)
		for i := 1; i <= NUM; i++ {
			var idx = i
			go func() {
				defer wait.Done()
				for j := (idx - 1) * MULTI; j < idx*MULTI; j++ {
					m.Add(j, j)
				}
			}()
		}
		wait.Wait()
		assert.Equal(t, m.Count(), NUM*MULTI)

		m.Range(func(k int, v int) {
			t.Logf("[%v] %v", k, v)
		})
	})
}
