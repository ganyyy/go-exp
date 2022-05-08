//go:build ignore

package main

import (
	"fmt"
	"patch"
	_ "plugin"
)

//go:noinline
func MainMul(a, b int) int {
	return a * b
}

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	const path = "./sum.so"

	var mainPlugin, sumPlugin *patch.Plugin
	var err error
	mainPlugin, err = patch.NewPlugin("")
	panicIfError(err)
	defer mainPlugin.Close()

	sumPlugin, err = patch.NewPlugin(path)
	panicIfError(err)
	defer sumPlugin.Close()

	var mulSymbol *patch.Symbol
	mulSymbol, err = mainPlugin.Symbol("main.MainMul")
	panicIfError(err)

	fmt.Printf("before patch calc result %v", MainMul(10, 20))

	var exchange = func(name string) {
		var sumSymbol, err = sumPlugin.Symbol(name)
		panicIfError(err)
		mulSymbol.Patch(sumSymbol)
		fmt.Printf("patch %v calc result %v", name, MainMul(10, 20))
		mulSymbol.Restore()
		fmt.Printf("un patch %v calc result %v", name, MainMul(10, 20))
	}

	exchange("plugin/unnamed-e25187071282631620d549f8ddd20d3a6e6dec10.Add")
	exchange("plugin/unnamed-e25187071282631620d549f8ddd20d3a6e6dec10.Sub")
}
