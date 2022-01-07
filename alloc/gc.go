package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"
)

type SimpleStruct struct {
	_ [100]int
}

type Struct2 struct {
	_   int
	Tmp []*SimpleStruct
}

var src []Struct2

func sliceStructGC() {
	for i := 0; i < 100; i++ {
		var tmp Struct2
		for j := 0; j < 100; j++ {
			tmp.Tmp = append(tmp.Tmp, &SimpleStruct{})
		}
		if rand.Intn(2) == 1 {
			src = src[len(src):]
		}
		src = append(src, tmp)
	}

	time.Sleep(time.Second)
	runtime.GC()

	log.Printf("len:%v, cap:%v", len(src), cap(src))

	time.Sleep(time.Second)
	runtime.GC()

	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	log.Printf("%+v", memStat)

	src = make([]Struct2, 0, 5)
	time.Sleep(time.Second)
	runtime.GC()

	runtime.ReadMemStats(&memStat)
	log.Printf("%+v", memStat)
}
