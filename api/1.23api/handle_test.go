package api

import (
	"testing"
	"unique"

	"github.com/stretchr/testify/require"
)

func TestHandle(t *testing.T) {
	// 这是一个很神奇的玩意
	// 本质上, 是将字面量相等的值使用同一块指针进行存储
	h := unique.Make(100)
	h2 := unique.Make(100)
	t.Logf("h: %v", h.Value())
	require.Equal(t, h, h2)
}
