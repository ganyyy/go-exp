package tests

import (
	"testing"

	"github.com/folays/luajit-go"
	"github.com/stretchr/testify/require"
)

func add(a, b int) int {
	return a + b
}

func sub(a, b int) int {
	return a - b
}

func TestLuaJit(t *testing.T) {

	t.Run("LuaStackTrace", func(t *testing.T) {
		L := luajit.NewState()
		L.ReasonableDefaults()
		require.NoError(t, L.RunString(`function buggy_function() error('Something went wrong') end`))
		err := L.RunFunc("buggy_function")
		t.Logf("Error: %v", err)
	})

	t.Run("BindFunc", func(t *testing.T) {
		L := luajit.NewState()
		L.ReasonableDefaults()
		L.SetGlobalAny("add", add)
		n, err := L.RunFuncWithResults(1, "add", 1)
		require.NoError(t, err)
		require.Equal(t, n, 1)
		t.Logf("val %v", L.ToInt(L.GetTop()))
		t.Logf("top: %v", L.GetTop())
		L.RunClearResults()
		t.Logf("top: %v", L.GetTop())
	})

	t.Run("GoTaleFunc", func(t *testing.T) {
		L := luajit.NewState()
		L.ReasonableDefaults()
		L.FuncAdd("go", "add", add)
		L.FuncAdd("go", "sub", sub)
		require.NoError(t, L.RunString(`print(go.add(1, 2))`))
		require.NoError(t, L.RunString(`print(go.sub(1, 2))`))

	})
}
