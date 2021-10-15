package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
	"unsafe"

	"go-exp/export/package_b"
)

type Outer struct {
	Name string
	Age int
	Tmp [1<<20]int
}


func TestExchangeGC() *Outer {
	var inner = package_b.GetInner("1234", 1234)
	var outer = *(**Outer)(unsafe.Pointer(&inner))
	return outer
}

func main() {

	go func() {
		var out = TestExchangeGC()
		time.Sleep(time.Minute * 2)
		out.Tmp[100] = 1000

		time.Sleep(time.Minute * 2)

		out.Tmp[200] = 1000

		log.Println(out.Tmp[200])
	}()

	_ = http.ListenAndServe(":9999", nil)
}
