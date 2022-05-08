package main

import (
	"log"
	"patch"
	"plugin"
	"reflect"
	"unsafe"
)

//go:noinline
func Add2(a, b int) int {
	return a * b
}

type funcVal struct {
	_   uintptr
	ptr unsafe.Pointer
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
		var handle, _ = plugin.Open("./sum.so")
		var s, err = handle.Lookup("Sub")
		log.Printf("symbol %v, error %v", s, err)

		var sub = s.(func(int, int) int)
		var value = reflect.ValueOf(s)
		var toAddr = (*funcVal)(unsafe.Pointer(&value)).ptr

		log.Printf("symbol %v, error %v, to addr:%v", s, err, toAddr)
		log.Println(sub(10, 20))
		log.Println(Add2(10, 20))
		patch.Patch(symbol, toAddr)
		log.Println(Add2(10, 20))
	}

	// log.Printf("load sum.so")

	// 这个是通过 go tool nm xxx.so 获取到的符号表
	// 和直接使用plugin还是有点区别的
	{
		handler, err := patch.PluginOpen(path)
		log.Printf("handler %v, error %v", handler, err)
		if err != nil {
			return
		}
		symbol, err := patch.LookupSymbol(handler, "plugin/unnamed-e25187071282631620d549f8ddd20d3a6e6dec10.Sub")
		// 尽管地址相同, 但是无法转换到具体的.text地址, 所以直接look是不可行滴
		var toAddr = *(*unsafe.Pointer)(symbol)
		log.Printf("symbol %v, error %v, to addr %v", symbol, err, toAddr)
		err = patch.PluginClose(handler)
		log.Printf("close handler error:%v", err)
	}

}
