package tests

import (
	"os"
	"testing"

	"ganyyy.com/go-exp/demo/golang-lua/scripts"
	go_lua "github.com/Shopify/go-lua"
	clua "github.com/aarzilli/golua/lua"
	"github.com/arnodel/golua/runtime"
	"github.com/folays/luajit-go"
	lua "github.com/yuin/gopher-lua"
)

const (
	V = 20
)

func loadGo_Lua() *go_lua.State {
	var l = go_lua.NewState()
	go_lua.DoFile(l, scripts.FibPath)
	return l
}

func testGo_Lua(l *go_lua.State) int {
	l.Global("Fib")
	l.PushInteger(V)
	l.Call(1, 1)
	ret, _ := l.ToInteger(-1)
	l.Pop(1)
	return ret
}

func loadGoLua() *runtime.Runtime {
	var run = runtime.New(os.Stdout)
	var source, _ = os.ReadFile(scripts.FibPath)
	chunk, _ := run.CompileAndLoadLuaChunk("fib", source, runtime.TableValue(run.GlobalEnv()))
	term := runtime.NewTerminationWith(run.MainThread().CurrentCont(), 0, false)
	// 执行脚本
	runtime.Call(run.MainThread(), runtime.FunctionValue(chunk), nil, term)
	return run
}

func testGoLua(run *runtime.Runtime, v runtime.Value) int64 {
	ret, _ := runtime.Call1(run.MainThread(), v, runtime.IntValue(V))
	val := ret.AsInt()
	return val
}

func loadGopherLua() *lua.LState {
	var l = lua.NewState()
	l.DoFile(scripts.FibPath)
	return l
}

func testGopherLua(state *lua.LState, v lua.LValue) int {
	state.Push(v)
	state.Push(lua.LNumber(V))
	state.Call(1, 1)
	ret := int(lua.LVAsNumber(state.Get(-1)))
	state.Pop(1)
	return ret
}

func loadCLua() *clua.State {
	l := clua.NewState()

	l.DoFile(scripts.FibPath)
	return l
}

func testClua(l *clua.State) int {
	l.GetGlobal("Fib")
	l.PushInteger(V)
	l.Call(1, 1)
	val := l.ToInteger(-1)
	l.Pop(1)
	return val
}

func loadGoLuaJit() (*luajit.State, luajit.Ref) {
	L := luajit.NewState()
	L.ReasonableDefaults()
	err := L.RunFile(scripts.FibPath)
	if err != nil {
		println(err)
	}
	L.GetGlobal("Fib")
	return L, L.RegistryRef()
}

// func testGoLuaJit(L *luajit.State, ref luajit.Ref) int {
// 	L.RunRefWithResults(1, ref, V)
// 	// L.RunFuncWithResults(1, "Fib", V)
// 	val := L.ToInt(L.GetTop())
// 	L.Pop(1)
// 	return val
// }

func testGoLuaJitFunc(L *luajit.State) int {
	L.RunFuncWithResults(1, "Fib", V)
	val := L.ToInt(L.GetTop())
	L.Pop(1)
	return val
}

func TestRunCheck(t *testing.T) {

	t.Run("go-lua", func(t *testing.T) {
		state := loadGo_Lua()
		ret := testGo_Lua(state)
		t.Log(ret)
	})

	t.Run("golua", func(t *testing.T) {
		run := loadGoLua()
		v := run.MainThread().GlobalEnv().Get(runtime.StringValue("Fib"))
		ret := testGoLua(run, v)
		t.Log(ret)
	})

	t.Run("gopher-lua", func(t *testing.T) {
		st := loadGopherLua()
		v := st.GetGlobal("Fib")
		ret := testGopherLua(st, v)
		t.Log(ret)
	})

	t.Run("clua", func(t *testing.T) {
		cl := loadCLua()
		v := testClua(cl)
		t.Log(v)
	})

	t.Run("go", func(t *testing.T) {
		ret := fib(V)
		t.Log(ret)
	})

	// t.Run("go-luajit", func(t *testing.T) {
	// 	L, ref := loadGoLuaJit()
	// 	ret := testGoLuaJit(L, ref)
	// 	t.Log(ret)
	// })

	t.Run("go-luajit-func", func(t *testing.T) {
		L, _ := loadGoLuaJit()
		ret := testGoLuaJitFunc(L)
		t.Log(ret)
	})
}

//go:noinline
func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func BenchmarkLuaFib(b *testing.B) {
	b.Run("go-lua", func(b *testing.B) {
		var l = loadGo_Lua()
		b.ResetTimer()
		for b.Loop() {
			testGo_Lua(l)
		}
	})

	b.Run("golua", func(b *testing.B) {
		run := loadGoLua()
		v := run.MainThread().GlobalEnv().Get(runtime.StringValue("Fib"))
		b.ResetTimer()
		for b.Loop() {
			testGoLua(run, v)
		}
	})

	b.Run("gopher-lua", func(b *testing.B) {
		st := loadGopherLua()
		v := st.GetGlobal("Fib")
		b.ResetTimer()
		for b.Loop() {
			testGopherLua(st, v)
		}
	})

	b.Run("clua", func(b *testing.B) {
		ct := loadCLua()
		b.ResetTimer()
		for b.Loop() {
			testClua(ct)
		}
	})

	// b.Run("go-luajit", func(b *testing.B) {
	// 	L, ref := loadGoLuaJit()
	// 	b.ResetTimer()
	// 	for b.Loop() {
	// 		testGoLuaJit(L, ref)
	// 	}
	// })

	b.Run("go-luajit-func", func(b *testing.B) {
		L, _ := loadGoLuaJit()
		b.ResetTimer()
		for b.Loop() {
			testGoLuaJitFunc(L)
		}
	})

	b.Run("go", func(b *testing.B) {
		for b.Loop() {
			fib(V)
		}
	})
}
