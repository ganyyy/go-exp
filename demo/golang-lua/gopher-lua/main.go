package main

import (
	"fmt"

	"ganyyy.com/go-exp/demo/golang-lua/scripts"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

type User = scripts.User

func main() {
	var l = lua.NewState()
	defer l.Close()

	var u = &User{
		Name: "100",
		Age:  200,
	}

	l.SetGlobal("user", luar.New(l, u))

	if err := l.DoFile(scripts.CounterPath); err != nil {
		panic(err)
	}

	v := l.GetGlobal("NXT")
	l.Push(v)
	l.Call(0, 1)
	ret := l.Get(-1)
	fmt.Println(ret)

	fmt.Println("go", u)
}
