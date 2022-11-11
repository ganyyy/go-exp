package main

import (
	"strings"
	"testing"

	"github.com/nats-io/nats.go"
)

func TestGetConn(t *testing.T) {

	var urls = strings.Join([]string{
		"localhost:4225",
		//"localhost:4223",
		//"localhost:4224",
	}, ",")
	pubNc, _ := nats.Connect(urls)
	conn := getConn(pubNc)
	n, _ := conn.Write([]byte("hello world"))
	t.Logf("n:%v", n)

	e := conn.Close()
	t.Logf("err:%v", e)
}
