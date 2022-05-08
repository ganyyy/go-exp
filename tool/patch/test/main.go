//go:build ignore

package main

import (
	"log"
	"patch"
	_ "plugin" // 需要依赖于一些cgo的编译选项, 否则无法加载全部的符号
)

//go:noinline
func Add2() {

}

func main() {
	const path = "./sum.so"

	var load = func() {
		mainHandler, err := patch.PluginOpen("")
		log.Printf("mainHandler %v, error %v", mainHandler, err)
		if err != nil {
			return
		}

		symbol, err := patch.LookupSymbol(mainHandler, "main.Add2")
		log.Printf("mainHandler symbol %v, error %v", symbol, err)
		patch.PluginClose(mainHandler)
	}

	load()

	// log.Printf("load sum.so")

	handler, err := patch.PluginOpen(path)
	log.Printf("handler %v, error %v", handler, err)
	if err != nil {
		return
	}

	// 这个是通过 go tool nm xxx.so 获取到的符号表
	// 和直接使用plugin还是有点区别的
	symbol, err := patch.LookupSymbol(handler, "plugin/unnamed-e25187071282631620d549f8ddd20d3a6e6dec10.Sub")
	log.Printf("symbol %v, error %v", symbol, err)

	err = patch.PluginClose(handler)
	log.Printf("close handler error:%v", err)

	Add2()
}
