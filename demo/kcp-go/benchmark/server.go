package main

import (
	"fmt"
	"log/slog"
	"net"

	kcp_benchmark_config "ganyyy.com/go-exp/demo/kcp-go/benchmark/config"
	"github.com/xtaci/kcp-go/v5"
)

func StartServer() {
	slog.Info("StartServer", slog.Bool("writeDelay", kcp_benchmark_config.Config.WriteDelay))

	kcp_benchmark_config.OpenPProf()

	listener, err := kcp.Listen(kcp_benchmark_config.Config.ServerAddr)
	if err != nil {
		kcp_benchmark_config.LogAndExit(fmt.Errorf("listen error: %v", err))
	}

	lis := NewListener(listener, func(conn net.Conn) error {
		c, ok := conn.(*kcp.UDPSession)
		if !ok {
			return nil
		}
		kcp_benchmark_config.InitKcpSession(c)
		return nil
	})
	lis.Start()
}
