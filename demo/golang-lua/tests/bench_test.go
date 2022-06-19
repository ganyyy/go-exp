package tests

import (
	"os"
	"testing"

	"ganyyy.com/go-exp/demo/golang-lua/scripts"
	go_lua "github.com/Shopify/go-lua"
	clua "github.com/aarzilli/golua/lua"
	"github.com/arnodel/golua/runtime"
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
	return ret.AsInt()
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
	return int(lua.LVAsNumber(state.Get(-1)))
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
	return l.ToInteger(-1)
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
}

func BenchmarkLuaFib(b *testing.B) {
	b.Run("go-lua", func(b *testing.B) {
		var l = loadGo_Lua()
		for i := 0; i < b.N; i++ {
			testGo_Lua(l)
		}
	})

	b.Run("golua", func(b *testing.B) {
		run := loadGoLua()
		v := run.MainThread().GlobalEnv().Get(runtime.StringValue("Fib"))
		for i := 0; i < b.N; i++ {
			testGoLua(run, v)
		}
	})

	b.Run("gopher-lua", func(b *testing.B) {
		st := loadGopherLua()
		v := st.GetGlobal("Fib")
		for i := 0; i < b.N; i++ {
			testGopherLua(st, v)
		}
	})

	b.Run("clua", func(b *testing.B) {
		ct := loadCLua()
		for i := 0; i < b.N; i++ {
			testClua(ct)
		}
	})
}
