package main

import (
	"fmt"
	"log/slog"

	kcp_benchmark_config "ganyyy.com/go-exp/demo/kcp-go/benchmark/config"
)

func StartServer() {
	slog.Info("StartServer", slog.Bool("writeDelay", kcp_benchmark_config.Config.WriteDelay))

	kcp_benchmark_config.OpenPProf()
	OpenServerMetrics()

	isKcp := kcp_benchmark_config.Config.IsKCP
	listenNum := kcp_benchmark_config.Config.ListenNum

	if !isKcp {
		listenNum = 1
	}

	for i := 0; i < listenNum; i++ {
		listener, err := option.Listen(kcp_benchmark_config.Config.ServerAddr)
		if err != nil {
			kcp_benchmark_config.LogAndExit(fmt.Errorf("listen error: %v", err))
		}
		go NewListener(i, listener, option.AcceptCallback).Start()
	}

}
