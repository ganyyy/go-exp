package main

import "log"

func init() {
	log.Println("plugin init")
}

func Add(a, b int) int {
	return a + b
}
