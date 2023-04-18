package json

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type BaseInner struct {
	Type string          `json:"type"` // 基础类型
	Data json.RawMessage `json:"data"` // 实际数据
}

type DataA struct {
	A int    `json:"a"`
	B string `json:"b"`
}

type DataB struct {
	Age  int    `json:"age"`
	Name string `json:"name"`
}

type DataC struct {
	Age2  int    `json:"age"`
	Name3 string `json:"name"`
}

// test BaseInner Unmarshal
func TestBaseInnerUnmarshal(t *testing.T) {
	var baseInner BaseInner
	var dataA DataA
	var dataB DataB
	var dataC DataC

	// test dataA
	baseInner.Type = "dataA"
	baseInner.Data = []byte(`{"a":1,"b":"test"}`)
	err := json.Unmarshal([]byte(`{"type":"dataA","data":{"a":1,"b":"test"}}`), &baseInner)
	assert.NoError(t, err)
	err = json.Unmarshal(baseInner.Data, &dataA)
	assert.NoError(t, err)
	assert.Equal(t, dataA.A, 1)
	assert.Equal(t, dataA.B, "test")

	// test dataB
	baseInner.Type = "dataB"
	baseInner.Data = []byte(`{"age":1,"name":"test"}`)
	err = json.Unmarshal([]byte(`{"type":"dataB","data":{"age":1,"name":"test"}}`), &baseInner)
	assert.NoError(t, err)
	err = json.Unmarshal(baseInner.Data, &dataB)
	assert.NoError(t, err)
	assert.Equal(t, dataB.Age, 1)
	assert.Equal(t, dataB.Name, "test")

	// test dataC
	baseInner.Type = "dataC"
	baseInner.Data = []byte(`{"age":1,"name":"test"}`)
	err = json.Unmarshal([]byte(`{"type":"dataC","data":{"age":1,"name":"test"}}`), &baseInner)
	assert.NoError(t, err)
	err = json.Unmarshal(baseInner.Data, &dataC)
	assert.NoError(t, err)
	assert.Equal(t, dataC.Age2, 1)
	assert.Equal(t, dataC.Name3, "test")
}
