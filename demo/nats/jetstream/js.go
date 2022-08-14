package main

import (
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	nc, _ = nats.Connect(nats.DefaultURL)
)

const (
	ORDER   = "ORDER"
	MONITOR = "MONITOR"
	Sub     = "scratch"
)

func Subject(src ...string) string {
	return strings.Join(src, ".")
}

var sbj = Subject(ORDER, Sub)

func CreateStream() {
	js, _ := nc.JetStream(nats.PublishAsyncMaxPending(256))

	// 删除流/消费者
	js.DeleteConsumer(ORDER, MONITOR)
	js.DeleteStream(ORDER)

	js.AddStream(&nats.StreamConfig{
		Name:      ORDER,
		Subjects:  []string{sbj},
		Retention: nats.WorkQueuePolicy,
	})
	js.UpdateStream(&nats.StreamConfig{
		MaxBytes: 8,
	})
	js.AddConsumer(ORDER, &nats.ConsumerConfig{
		Durable: MONITOR,
	})

	js.Publish(sbj, []byte("hello"))

	for i := 0; i < 5000; i++ {
		js.PublishAsync(sbj, []byte("hello "+strconv.Itoa(i)))
		time.Sleep(time.Second)
	}

}

func SubcribeStream(id int) {
	js, _ := nc.JetStream()
	js.AddStream(&nats.StreamConfig{
		Name:      ORDER,
		Subjects:  []string{sbj},
		Retention: nats.WorkQueuePolicy,
	})
	js.UpdateStream(&nats.StreamConfig{
		Name:     ORDER,
		MaxBytes: 8,
	})

	js.UpdateConsumer(ORDER, &nats.ConsumerConfig{
		Durable: MONITOR,
	})

	js.Subscribe(sbj, func(msg *nats.Msg) {
		log.Printf("[%v] Received a message: %s", id, string(msg.Data))
	})

}

func main() {
	var wait sync.WaitGroup
	wait.Add(1)

	go func() {
		defer wait.Done()
		CreateStream()
	}()

	SubcribeStream(1)
	wait.Wait()
}
