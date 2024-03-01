package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	kcp_benchmark_config "ganyyy.com/go-exp/demo/kcp-go/benchmark/config"
	prometheusmetrics "github.com/deathowl/go-metrics-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rcrowley/go-metrics"
	"github.com/xtaci/kcp-go/v5"
)

var (
	register = metrics.NewRegistry()

	echoTimer metrics.Timer

	connGauge    metrics.Gauge     // 连接计数
	retransGauge metrics.Gauge     // 重传计数
	inSegments   metrics.Gauge     // 接收数据包计数
	outSegments  metrics.Gauge     // 发送数据包计数
	inBytes      metrics.Gauge     // 接收字节数
	outBytes     metrics.Gauge     // 发送字节数
	rtoHistogrm  metrics.Histogram // RTO分布

	kcpAcceptCounts []metrics.Counter
)

func InitMetrics(isServer bool) {

	var subSystem = "echo"
	var listenAddr = ":2999"
	if isServer {
		subSystem = "server"
		listenAddr = ":3999"
	}

	var client = prometheusmetrics.NewPrometheusProvider(
		register, "kcp", subSystem,
		prometheus.DefaultRegisterer, 1*time.Second,
	)
	http.Handle("/metrics", promhttp.Handler())
	go client.UpdatePrometheusMetrics()
	go func() {
		err := http.ListenAndServe(listenAddr, nil)
		if err != nil {
			kcp_benchmark_config.LogAndExit(fmt.Errorf("start metrics error: %v", err))
		}
	}()
}

func GetTotalFrame() int64 {
	frame := 1000 / kcp_benchmark_config.Config.EchoInterval
	return int64(frame * kcp_benchmark_config.Config.ClientNum)
}

func OpenClientMetrics() {
	InitMetrics(kcp_benchmark_config.Config.IsServer)
	sampleNum := int(GetTotalFrame())
	slog.Info("sample", slog.Int("Num", sampleNum))
	var sample = metrics.NewExpDecaySample(sampleNum, 0.015)
	var histogram = metrics.NewHistogram(sample)
	echoTimer = metrics.NewCustomTimer(histogram, metrics.NewMeter())
	register.Register("echo", echoTimer)
}

type IServerMetrics interface {
	Reset()
	Connections() int64
	AddConn(int)
	DelConn(int)
}

const (
	MRetransSegs = iota
	MInSegs
	MOutSegs
	MInBytes
	MOutBytes
	MCount
)

type MetricsValue struct{ Val uint64 }

func (m *MetricsValue) Diff(newVal uint64) int64 {
	diff := newVal - m.Val
	m.Val = newVal
	return int64(diff)
}

type KCPMetrics struct {
	last [MCount]MetricsValue
	snmp kcp.Snmp
}

func (m *KCPMetrics) Reset()             { m.snmp = *kcp.DefaultSnmp.Copy() }
func (m *KCPMetrics) Connections() int64 { return int64(m.snmp.CurrEstab) }
func (m *KCPMetrics) AddConn(idx int) {
	if uint(idx) < uint(len(kcpAcceptCounts)) {
		kcpAcceptCounts[idx].Inc(1)
	}
}
func (m *KCPMetrics) DelConn(idx int) {
	if uint(idx) < uint(len(kcpAcceptCounts)) {
		kcpAcceptCounts[idx].Dec(1)
	}
}

// Retrans
func (m *KCPMetrics) Retrans() int64 { return m.last[MRetransSegs].Diff(m.snmp.RetransSegs) }

// InSegs
func (m *KCPMetrics) InSegs() int64 { return m.last[MInSegs].Diff(m.snmp.InSegs) }

// OutSegs
func (m *KCPMetrics) OutSegs() int64 { return m.last[MOutSegs].Diff(m.snmp.OutSegs) }

// InBytes
func (m *KCPMetrics) InBytes() int64 { return m.last[MInBytes].Diff(m.snmp.InBytes) }

// OutBytes
func (m *KCPMetrics) OutBytes() int64 { return m.last[MOutBytes].Diff(m.snmp.OutBytes) }

type TCPMetrics struct{ conn atomic.Int64 }

func (m *TCPMetrics) Reset()             {}
func (m *TCPMetrics) Connections() int64 { return m.conn.Load() }
func (m *TCPMetrics) Retrans() int64     { return 0 }
func (m *TCPMetrics) AddConn(int)        { m.conn.Add(1) }
func (m *TCPMetrics) DelConn(int)        { m.conn.Add(-1) }

func OpenServerMetrics() {
	InitMetrics(kcp_benchmark_config.Config.IsServer)

	connGauge = metrics.NewGauge()
	register.Register("conn", connGauge)

	retransGauge = metrics.NewGauge()
	register.Register("retrans", retransGauge)

	inSegments = metrics.NewGauge()
	register.Register("inSegments", inSegments)

	outSegments = metrics.NewGauge()
	register.Register("outSegments", outSegments)

	inBytes = metrics.NewGauge()
	register.Register("inBytes", inBytes)

	outBytes = metrics.NewGauge()
	register.Register("outBytes", outBytes)

	sample := int(GetTotalFrame())
	rtoHistogrm = metrics.NewHistogram(metrics.NewExpDecaySample(sample, 0.015))
	register.Register("rto", rtoHistogrm)

	slog.Info("server rto sample", slog.Int("Num", sample))

	serverMetrics := option.ServerMetrics

	var kcpMetrics, _ = serverMetrics.(*KCPMetrics)

	go func() {
		var ticker = time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			serverMetrics.Reset()
			connGauge.Update(serverMetrics.Connections())

			if kcpMetrics != nil {
				retransGauge.Update(kcpMetrics.Retrans())
				inSegments.Update(kcpMetrics.InSegs())
				outSegments.Update(kcpMetrics.OutSegs())
				inBytes.Update(kcpMetrics.InBytes())
				outBytes.Update(kcpMetrics.OutBytes())
			}
		}
	}()
}
