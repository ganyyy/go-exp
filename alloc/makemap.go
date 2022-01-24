package main

import (
	"log"
)

func IteratorMap(m map[string]struct{}) {
	for k := range m {
		log.Printf("%v", k)
	}
}

func MakeMap1() {
	var m = make(map[string]struct{})
	m["123"] = struct{}{}

	IteratorMap(m)
}

func TempMap() {
	IteratorMap(map[string]struct{}{
		"123": {},
	})
}
