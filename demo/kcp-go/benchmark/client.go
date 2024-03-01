package main

import (
	"log/slog"

	kcp_benchmark_config "ganyyy.com/go-exp/demo/kcp-go/benchmark/config"
)

func StartClient() {
	slog.Info("StartClient")
	OpenClientMetrics()

	clientNum := kcp_benchmark_config.Config.ClientNum

	RunClients(clientNum, kcp_benchmark_config.Config.ServerAddr, option.Dial)
}
