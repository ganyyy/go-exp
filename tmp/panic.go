//go:build ignore

package main

import (
	"fmt"
	"runtime"
)

func main() {
	var s []int
	var ps *struct{ n int }
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover error is", err)

			fmt.Println("call stack is")
			fmt.Println()
			fmt.Println()

			var frame [32]uintptr
			n := runtime.Callers(0, frame[:])
			frames := runtime.CallersFrames(frame[:n])
			for {
				frame, more := frames.Next()
				fmt.Printf("%s:%d %s\n", frame.File, frame.Line, frame.Function)
				if !more {
					break
				}
			}
		}
	}()

	defer func() { println("index out of range"); _ = s[2] }()

	defer func() { println("nil pointer dereference"); _ = ps.n }()

	panic(3)
}
