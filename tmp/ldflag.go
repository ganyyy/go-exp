package main

import (
	"fmt"
)

//go:generate go build -ldflags "-X 'main.Version=$(go version)' -X 'main.Val=`date +%s`'"  -o version ldflag.go

var Version string
var Val string

func main() {
	fmt.Println(Version, Val)
}
