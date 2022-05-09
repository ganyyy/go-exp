package main

import (
	"log"
	"patch"
	"plugin"
)

//go:noinline
func Add2(a, b int) int {
	return a * b
}

func main() {
	const path = "./sum.so"

	mainHandler, err := patch.PluginOpen("")
	log.Printf("mainHandler %v, error %v", mainHandler, err)
	if err != nil {
		return
	}

	symbol, err := patch.LookupSymbol(mainHandler, "main.Add2")
	log.Printf("mainHandler symbol %v, error %v", symbol, err)
	// patch.PluginClose(mainHandler)
	{
		var handle, _ = plugin.Open(path)
		var s, err = handle.Lookup("Sub")

		var toAddr = patch.FuncAddr(s)

		log.Printf("symbol %v, error %v, to addr:%v", s, err, toAddr)
		log.Println(Add2(10, 20))
		var data = patch.Backup(symbol)
		patch.Patch(symbol, toAddr)
		log.Println(Add2(10, 20))
		patch.Restore(symbol, data)
		log.Println(Add2(10, 20))

		s, _ = handle.Lookup("Add")
		toAddr = patch.FuncAddr(s)
		patch.Patch(symbol, toAddr)
		log.Println(Add2(10, 20))
		patch.Restore(symbol, data)
		log.Println(Add2(10, 20))
	}

	// log.Printf("load sum.so")

	// 这个是通过 go tool nm xxx.so 获取到的符号表
	// 和直接使用plugin还是有点区别的
	// 还是要通过plugin加载的形式, 才能取到函数的地址
	// 放弃了~
	{
		// handler, err := patch.PluginOpen(path)
		// log.Printf("handler %v, error %v", handler, err)
		// if err != nil {
		// 	return
		// }
		// symbol, err := patch.LookupSymbol(handler, "plugin/unnamed-e25187071282631620d549f8ddd20d3a6e6dec10.Sub")
		// // 尽管地址相同, 但是无法转换到具体的.text地址, 所以直接look是不可行滴
		// var pair = *(*uintptrPair)(unsafe.Pointer(&symbol))
		// var toAddr = (unsafe.Pointer)(uintptr(unsafe.Pointer(&symbol)))

		// log.Printf("symbol %v, error %v, to addr %v, pair:%v", symbol, err, toAddr, pair)
		// err = patch.PluginClose(handler)
		// log.Printf("close handler error:%v", err)
	}

}
