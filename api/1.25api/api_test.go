package api

import (
	"encoding/json/v2"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncTest(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var counter atomic.Uint64
		for range 100 {
			go func() {
				counter.Add(1)
			}()
		}
		synctest.Wait()
		assert.Equal(t, uint64(100), counter.Load())
	})
}

func TestJsonV2(t *testing.T) {
	type Data struct {
		Name     string            `json:"name"`
		Age      int               `json:"age"`
		Address  []string          `json:"address"`
		Address2 []string          `json:"address2,format:emitempty"`
		Tags     map[string]string `json:"tags"`
	}

	var data Data
	data.Name = "Alice"
	data.Age = 30

	bs, err := json.Marshal(data,
		json.FormatNilSliceAsNull(true),
		json.OmitZeroStructFields(true),
	)
	require.NoError(t, err)
	t.Logf("JSON: %s", bs)
}
