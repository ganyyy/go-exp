package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
)

type natsWarpper struct {
	*nats.Conn
	log *slog.Logger
}

func forNatsCallback() {

	const (
		connNum = 30
	)

	var allConn = make([]natsWarpper, 0, connNum)
	for i := 0; i < connNum; i++ {
		idxLog := slog.Default().With(slog.Int("index", i))
		conn, err := createNatsConn(urls, idxLog)
		if err != nil {
			idxLog.Error("create nats conn error", slog.String("err", err.Error()))
			continue
		}
		idxLog.Info("create nats conn success")
		allConn = append(allConn, natsWarpper{
			Conn: conn,
			log:  idxLog,
		})
	}

	defer func() {
		for _, conn := range allConn {
			conn.Close()
			conn.log.Info("close nats conn")
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
}
