package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"

	kcp_benchmark_config "ganyyy.com/go-exp/demo/kcp-go/benchmark/config"
)

var configPath = flag.String("config", "config.toml", "config file path")

func main() {
	flag.Parse()

	kcp_benchmark_config.MustReadConfig(*configPath)

	var sigChan = make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)

	go func() {
		if kcp_benchmark_config.Config.IsServer {
			StartServer()
		} else {
			StartClient()
		}
	}()
	<-sigChan
	slog.Info("exit")
}
