package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

const output = "./msg.log"

const subject = "test.cluster.subj"

func main() {
	_ = syscall.Unlink(output)

	f, _ := os.Create(output)
	var buf = bufio.NewWriter(f)
	defer f.Close()
	defer buf.Flush()

	var urls = strings.Join([]string{
		"localhost:4225",
		"localhost:4223",
		"localhost:4224",
	}, ",")

	nc, e := nats.Connect(urls)
	log.Println(e)

	nc.Subscribe(subject, func(msg *nats.Msg) {
		buf.WriteString(string(msg.Data))
		buf.WriteByte('\n')
	})

	var sendBuf = make([]byte, 6291456+1024)
	for i := 0; i < 10; i++ {
		go func(i int) {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(30)) * 100)
			nc, e := nats.Connect(urls)
			log.Println(e)
			var cnt int
			for {
				cnt++
				sendErr := nc.Publish(subject, sendBuf[:])
				if sendErr != nil {
					log.Println(i, sendErr)
				} else {
					log.Println(i, "success")
				}
				time.Sleep(time.Second)
			}
		}(i)
	}

	time.Sleep(time.Minute * 20)
}
