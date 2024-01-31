package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	kcp_benchmark_config "ganyyy.com/go-exp/demo/kcp-go/benchmark/config"
	prometheusmetrics "github.com/deathowl/go-metrics-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rcrowley/go-metrics"
)

var (
	echoTimer metrics.Timer
)

func OpenMetrics() {

	frame := 1000 / kcp_benchmark_config.Config.EchoInterval
	sampleNum := int(float64(frame*kcp_benchmark_config.Config.ClientNum) * 0.8)
	slog.Info("sample", slog.Int("Num", sampleNum))
	var sample = metrics.NewExpDecaySample(sampleNum, 0.015)
	var histogram = metrics.NewHistogram(sample)
	echoTimer = metrics.NewCustomTimer(histogram, metrics.NewMeter())
	var register = metrics.NewRegistry()
	var client = prometheusmetrics.NewPrometheusProvider(
		register, "kcp", "echo",
		prometheus.DefaultRegisterer, 1*time.Second,
	)
	register.Register("echo", echoTimer)
	go client.UpdatePrometheusMetrics()
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":2999", nil)
		if err != nil {
			kcp_benchmark_config.LogAndExit(fmt.Errorf("start metrics error: %v", err))
		}
	}()
}
