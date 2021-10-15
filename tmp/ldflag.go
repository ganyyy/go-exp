//go:build ignore
package main

import (
	"fmt"
)

//go:generate go build -tags release -ldflags "-X 'main.Version=$(go version)' -X 'main.Val=`date +%s`'"  -o version ldflag.go

var Version string
var Val string

func main() {
	fmt.Println(Version, Val, Ver)
}
