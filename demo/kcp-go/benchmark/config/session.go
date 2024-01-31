package kcp_benchmark_config

import "github.com/xtaci/kcp-go/v5"

func InitKcpSession(session *kcp.UDPSession) {
	session.SetWriteDelay(Config.WriteDelay)
	session.SetNoDelay(1, Config.Intervals, 2, 1)
}
