package main

import (
	"bufio"
	"log"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

var urls = strings.Join([]string{
	"localhost:4225",
	"localhost:4223",
	"localhost:4224",
}, ",")

func main() {

	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))

	forNatsCallback()

	// nc, e := nats.Connect(urls)
	// _ = nc
	// if e != nil {
	// 	slog.Error("connect nats error", slog.String("err", e.Error()))
	// 	return
	// }
	// defer nc.Close()

	// forMaxPayload(nc)
	// forMultiSubscribe(nc)
}

func forMultiSubscribe(nc *nats.Conn) {
	var subjects = []string{
		"test.cluster.subj",
		"test.cluster.subj2",
		"test.cluster.subj3",
		"test.cluster.subj4",
		"test.cluster.subj5",
	}

	for _, subject := range subjects {
		subject := subject
		nc.Subscribe(subject, func(msg *nats.Msg) {
			slog.Info("receive nats msg",
				slog.String("subject", subject),
				slog.String("msg", string(msg.Data)),
			)
		})
	}

	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	var ticker = time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			slog.Info("receive signal, exit")
			return
			// case <-ticker.C:
			// 	for _, subject := range subjects {
			// 		nc.Publish(subject, []byte("hello world"))
			// 	}
		}
	}
}

func forMaxPayload(nc *nats.Conn) {
	const output = "./msg.log"

	_ = syscall.Unlink(output)

	const subject = "test.cluster.subj"

	f, _ := os.Create(output)
	var buf = bufio.NewWriter(f)
	defer f.Close()
	defer buf.Flush()
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
					slog.Error("send nats msg error",
						slog.Int("idx", i),
						slog.String("err", sendErr.Error()),
					)
				} else {
					slog.Info("send nats msg success", slog.Int("idx", i), slog.Int("cnt", cnt))
				}
				time.Sleep(time.Second)
			}
		}(i)
	}

	time.Sleep(time.Minute * 20)
}
