package main

import (
	"flag"
	"fmt"
	"plugin"
)

func main() {
	var pluginPath string
	flag.StringVar(&pluginPath, "plugin", "", "path to the plugin")
	flag.Parse()

	if pluginPath == "" {
		flag.Usage()
		return
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		panic(err)
	}
	sum, err := p.Lookup("Sum3")
	if err != nil {
		panic("Sum3 not found")
	}
	sum3, ok := sum.(func([]int) int)
	if !ok {
		panic("Sum3 has incorrect signature")
	}

	fmt.Println("sum3([]int{1, 2, 3}):", sum3([]int{1, 2, 3})) // Call the plugin function
}
