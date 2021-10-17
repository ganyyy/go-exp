//go:build ignore
// +build ignore

package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

type Test struct {
}

func (t Test) Get() {

}

func main() {

	var t = &Test{}
	t.Get()

	go func() {
		var sigChan = make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGALRM, syscall.SIGPROF, syscall.SIGURG)
		for sig := range sigChan {
			println(sig)
		}
	}()

	http.ListenAndServe(":9999", nil)

}
