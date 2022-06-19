//go:build linux

package main

import (
	"fmt"

	"ganyyy.com/go-exp/demo/golang-lua/scripts"
	"github.com/aarzilli/golua/lua"
)

//go build -tags=luaa
func main() {
	{
		l := lua.NewState()
		defer l.Close()
		l.DoFile(scripts.FibPath)

		l.GetGlobal("Fib")
		l.PushInteger(15)
		l.Call(1, 1)
		fmt.Println(l.ToInteger(-1))
	}

	{

		var user = &scripts.User{
			Name: "123",
			Age:  12321,
		}

		l := lua.NewState()

		defer l.Close()
		l.OpenLibs()
		l.PushGoStruct(user)
		l.SetGlobal("user")

		l.SetTop(0)

		l.DoFile(scripts.CounterPath)
		// l.GetGlobal("user")
		// var u = l.ToGoStruct(-1).(*scripts.User)
		// fmt.Println(user, u)
		l.GetGlobal("NXT")
		l.Call(0, 1)
		fmt.Println(l.ToInteger(-1))
	}

}
