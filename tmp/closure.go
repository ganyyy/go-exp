package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

//go:noinline
func Closure() func() int {
	var v int
	return func() int {
		v++
		return v
	}
}

func GenClosure() {
	var g1 = Closure()
	g1()
}

func HTTP2() {
	var bs = bytes.NewBuffer(nil)
	var req, _ = http.NewRequest(http.MethodGet, "https://www.baidu.com/", bs)
	var client http.Client

	var resp, err = client.Do(req)
	if err != nil {
		log.Fatalf("GET error:%v", err)
		return
	}
	defer resp.Body.Close()

	var ret, _ = ioutil.ReadAll(resp.Body)
	log.Printf("%s", string(ret))
}
