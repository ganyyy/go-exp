package main

import (
	"log"
	"plugin"
)

func init() {
	log.Println("main init start")
}

func main() {
	log.Println("main func exec")

	var p *plugin.Plugin
	var s plugin.Symbol
	var err error
	p, err = plugin.Open("./plugin.so")
	if err != nil {
		log.Panicf("open plugin error:%v", err)
	}
	s, err = p.Lookup("Add")
	if err != nil {
		log.Panicf("lookup symbol error:%v", err)
	}
	var ret = s.(func(int, int) int)(10, 20)
	log.Println("ret is:", ret)
}
