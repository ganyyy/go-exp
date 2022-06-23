package enum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetData(t *testing.T) {
	data, ok := DataFromString(D1.String())
	assert.Equal(t, data, D1)
	assert.True(t, ok)

	data2, ok := Data2FromString(D1.String())
	assert.False(t, ok)
	assert.Equal(t, data2, DD1) // 这是个默认值?
}
