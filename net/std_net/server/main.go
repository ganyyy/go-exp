package main

import (
	"io"
	"log/slog"
	"net"
	"os"

	"ganyyy.com/go-exp/helper"
)

const ADDR = "0.0.0.0:9999"

func main() {

	logger := helper.InitSlog()

	listen, err := net.Listen("tcp", ADDR)
	if err != nil {
		logger.Error("listen error", slog.Any("err", err))
		os.Exit(1)
	}
	defer listen.Close()

	for i := 0; i < 100; i++ {
		go func(i int) {
			var total int
			for j := 1000000 * i; j < 1000000*(i+1); j++ {
				total += j
			}
		}(i)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			logger.Error("accept error", slog.Any("err", err))
			break
		}
		go func() {
			defer conn.Close()
			n, cpErr := io.Copy(conn, conn)
			if cpErr != nil {
				logger.Error("copy error", slog.Any("err", cpErr))
			} else {
				logger.Info("copy success", slog.Int64("n", n))
			}

		}()
	}
}
