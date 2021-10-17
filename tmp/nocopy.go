//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"sync"
)

func main() {

	var wait sync.WaitGroup
	wait.Add(1)
	go func(wg sync.WaitGroup) {
		defer wg.Done()
		fmt.Println("hello")
	}(wait)
	wait.Wait()
}
