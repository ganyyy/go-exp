package main

import (
	"fmt"

	"ganyyy.com/go-exp/demo/golang-lua/scripts"
	"github.com/Shopify/go-lua"
)

func main() {
	var L = lua.NewState()
	lua.OpenLibraries(L)

	if err := lua.DoFile(L, scripts.CounterPath); err != nil {
		fmt.Println(err)
	}

	var user = &scripts.User{
		Name: "123",
		Age:  100,
	}
	L.PushUserData(user)
	L.SetGlobal("user")
	L.Global("NXT")
	L.Call(0, 1)
	var ret, _ = L.ToNumber(-1)
	fmt.Println(ret)
}
