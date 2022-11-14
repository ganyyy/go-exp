package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/nats-io/nats.go"
)

const subject = "test.cluster.subj"
const topicNum = 5
const connField = "conn"

func init() {
	rand.Seed(time.Now().UnixNano())
}

var field = func() reflect.StructField {
	var conn nats.Conn
	var rt = reflect.TypeOf(&conn).Elem()
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if f.Name == connField {
			return f
		}
	}
	panic("cannot find conn field")
}()

func getConn(conn *nats.Conn) net.Conn {
	return *(*net.Conn)(unsafe.Pointer(uintptr(unsafe.Pointer(conn)) + field.Offset))
}

func main() {
	var urls = strings.Join([]string{
		"localhost:4225",
		"localhost:4223",
		"localhost:4224",
	}, ",")

	genSubject := func(i int) string {
		return subject + "-" + strconv.Itoa(i)
	}
	var pubNc *nats.Conn
	var e error
	go func() {
		pubNc, e = nats.Connect(urls, nats.DisconnectErrHandler(func(conn *nats.Conn, err error) {
			log.Printf("error:%v", err)
		}))
		log.Println("sub connect:", e)
		for i := 0; i < topicNum; i++ {
			topic := genSubject(i)
			_, e = pubNc.Subscribe(topic, func(msg *nats.Msg) {
				log.Printf("topic %v publish data:%s", topic, msg.Data)
			})
			if e != nil {
				log.Printf("sub %v error :%v", topic, e)
			}
		}
	}()
	nc, e := nats.Connect(urls)
	log.Println("pub connect error:", e)
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
		case "error":
			conn := getConn(pubNc)
			_, err = conn.Write([]byte("hello world"))
			log.Printf("writer error:%v", err)
		default:
			subj := genSubject(rand.Intn(topicNum))
			err = nc.Publish(subj, []byte(time.Now().String()))
			log.Printf("send to %v", subj)
			if err != nil {
				panic(err)
			}
		}
	}
}
