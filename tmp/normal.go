package main

import (
	"log"
	"strconv"
)

func LoopTmp() {
	type LoopValue struct {
		_   [10]int
		AAA string
	}

	var testFunc = func(val interface{}) {
		loop, ok := val.(LoopValue)
		if !ok {
			return
		}
		log.Printf("%+v", loop)
	}

	for i := 0; i < 10; i++ {
		var t LoopValue
		t.AAA = strconv.Itoa(i)
		testFunc(t)
	}
}
