package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

const subject = "test.cluster.subj"

func main() {
	var urls = strings.Join([]string{
		"localhost:4225",
		"localhost:4223",
		"localhost:4224",
	}, ",")

	go func() {
		nc, e := nats.Connect(urls)
		log.Println("sub connect:", e)
		nc.Subscribe(subject, func(msg *nats.Msg) {
			log.Printf("msg data:%s", msg.Data)
		})
	}()
	nc, e := nats.Connect(urls)
	log.Println("pub connect:", e)
	for {
		var b string
		n, err := fmt.Scanf("%v", &b)
		if err != nil {
			panic(err)
		}
		if n != 1 {
			return
		}
		switch b {
		case "quit":
			return
		default:
			err = nc.Publish(subject, []byte(time.Now().String()))
			if err != nil {
				panic(err)
			}
		}
	}
}
