package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

func createNatsConn(url string, log *slog.Logger) (*nats.Conn, error) {
	if log == nil {
		log = slog.Default()
	}
	return nats.Connect(url,
		nats.DisconnectHandler(func(c *nats.Conn) {
			log.Info("disconnect nats",
				slog.String("url", c.ConnectedAddr()))
		}),
		nats.DisconnectErrHandler(func(c *nats.Conn, _ error) {
			log.Info("disconnect nats with error",
				slog.String("url",
					c.ConnectedAddr()),
				slog.Any("err", c.LastError()))
		}),
		nats.ConnectHandler(func(c *nats.Conn) {
			log.Info("connect nats",
				slog.String("url", c.ConnectedAddr()))
		}),
		nats.ReconnectHandler(func(c *nats.Conn) {
			log.Info("reconnect nats",
				slog.String("url", c.ConnectedAddr()))
		}),
		nats.ClosedHandler(func(c *nats.Conn) {
			log.Info("close nats",
				slog.String("url", c.ConnectedAddr()))
		}),
		nats.ErrorHandler(func(c *nats.Conn, s *nats.Subscription, err error) {
			log.Info("nats error",
				slog.String("url", c.ConnectedAddr()),
				slog.Any("err", err),
				slog.String("subject", s.Subject))
		}),
	)
}

func forQueueSub() {
	const (
		connNum = 10
		subj    = "test.cluster.subj"
		queue   = "QUEUE"
	)

	var closeChan = make(chan struct{})
	var allConn = make([]*nats.Conn, 0, connNum)
	for i := 0; i < connNum; i++ {
		nc, e := createNatsConn(urls, nil)
		if e != nil {
			slog.Error("connect nats error", slog.String("err", e.Error()))
			continue
		}
		allConn = append(allConn, nc)
	}

	var recvCount [connNum]atomic.Int64

	var wg sync.WaitGroup
	wg.Add(connNum)

	if len(allConn) == 0 {
		slog.Error("no nats conn")
		return
	}

	var normalCount atomic.Int64
	var queueCount atomic.Int64

	allConn[0].Subscribe(subj, func(msg *nats.Msg) {
		normalCount.Add(1)
	})

	for idx, nc := range allConn {
		idx := idx
		nc.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
			recvCount[idx].Add(1)
			queueCount.Add(1)
		})
		go func(idx int, nc *nats.Conn) {
			defer wg.Done()
			var sendTicker = time.NewTicker(time.Millisecond * 10)
			defer sendTicker.Stop()
			defer func() {
				nc.Close()
				slog.Info("close nats conn", slog.Int("idx", idx))
			}()
			for {
				select {
				case <-closeChan:
					return
				case <-sendTicker.C:
					err := nc.Publish(subj, []byte("hello world"))
					if err != nil {
						slog.Error("send nats msg error",
							slog.Int("idx", idx),
							slog.String("err", err.Error()),
						)
						return
					}
				}
			}
		}(idx, nc)
	}

	var sigChan = make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	var ticker = time.NewTicker(time.Second)
	defer ticker.Stop()

	var buffer strings.Builder

	writeCount := func(idx int, count *atomic.Int64) int64 {
		var buf [8]byte
		_ = idx
		cnt := count.Swap(0)
		buffer.Write(fmt.Appendf(buf[:0], "%04d", cnt))
		buffer.WriteString("  |  ")
		return cnt
	}

	writeHeader := func(idx int) {
		var buf [8]byte
		buffer.Write(fmt.Appendf(buf[:0], "%4d", idx))
		buffer.WriteString("  |  ")
	}

	writeHeader(-1)
	writeHeader(-2)
	for idx := range recvCount {
		writeHeader(idx)
	}
	fmt.Println(buffer.String())
	buffer.Reset()
	fmt.Println(strings.Repeat("---------", connNum+2))
	for {
		select {
		case <-ticker.C:
			writeCount(-1, &normalCount)
			writeCount(-2, &queueCount)
			for idx := range recvCount {
				count := &recvCount[idx]
				writeCount(idx, count)
			}
			fmt.Println(buffer.String())
			buffer.Reset()
		case <-sigChan:
			close(closeChan)
			wg.Wait()
			slog.Info("receive signal, exit")
			return
		}
	}

}
