package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	kcp_benchmark_config "ganyyy.com/go-exp/demo/kcp-go/benchmark/config"
	"github.com/rcrowley/go-metrics"
	"github.com/xtaci/kcp-go/v5"
)

var configPath = flag.String("config", "config.toml", "config file path")

func main() {
	flag.Parse()

	kcp_benchmark_config.MustReadConfig(*configPath)

	var sigChan = make(chan os.Signal, 1)

	signal.Notify(sigChan, os.Interrupt)

	if kcp_benchmark_config.Config.IsKCP {
		{
			listNum := kcp_benchmark_config.Config.ListenNum
			kcpAcceptCounts = make([]metrics.Counter, listNum)
			for i := 0; i < listNum; i++ {
				kcpAcceptCounts[i] = metrics.NewCounter()
				register.Register(fmt.Sprintf("kcpAcceptCount_%d", i), kcpAcceptCounts[i])
			}
		}
		option.AcceptCallback = AcceptKcpSession
		option.Dial = func(addr string) (net.Conn, error) {
			return kcp.DialWithOptions(addr, nil, 0, 0)
		}
		var listenFD uintptr
		option.Listen = func(addr string) (net.Listener, error) {
			var cfg net.ListenConfig
			cfg.Control = func(network, address string, c syscall.RawConn) error {
				return c.Control(func(fd uintptr) {
					slog.Info("Control",
						slog.String("network", network),
						slog.String("address", address),
						slog.Int64("fd", int64(fd)),
					)
					listenFD = fd
					syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)
				})
			}
			conn, err := cfg.ListenPacket(context.Background(), "udp", addr)
			if err != nil {
				return nil, err
			}
			listener, err := kcp.ServeConn(nil, 0, 0, conn)
			if err != nil {
				return nil, err
			}
			return &KcpListener{
				Listener:   listener,
				PacketConn: conn,
				listenFD:   listenFD,
			}, nil
		}
		option.ServerMetrics = &KCPMetrics{}
	} else {
		option.Dial = DailTCP
		option.ServerMetrics = &TCPMetrics{}
		option.Listen = func(s string) (net.Listener, error) {
			return net.Listen("tcp", s)
		}
	}

	slog.Info("start", slog.Any("config", kcp_benchmark_config.Config))

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
