package main

import (
	"math/rand"
	"net/http"
	"time"

	prometheusmetrics "github.com/deathowl/go-metrics-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rcrowley/go-metrics"
)

var (
	registry = metrics.NewRegistry()
	client   = prometheusmetrics.NewPrometheusProvider(
		registry, "server", "main", prometheus.DefaultRegisterer, 1*time.Second,
	)
)

func metricsGauge() {
	var gauge = metrics.NewGauge()
	registry.Register("gauge312313", gauge)
	for {
		time.Sleep(time.Second * time.Duration(rand.Intn(10)+5))
		gauge.Update(rand.Int63n(100))
	}
}

func metricsCounter() {
	var counter = metrics.NewCounter()
	registry.Register("counter123213", counter)
	for {
		time.Sleep(time.Second * time.Duration(rand.Intn(10)+5))
		counter.Inc(rand.Int63n(100))
	}
}

func metricsMeter() {
	var meter = metrics.NewMeter()
	registry.Register("meter13123", meter)
	for {
		time.Sleep(time.Second * time.Duration(rand.Intn(10)+5))
		meter.Mark(rand.Int63n(100))
	}
}

func main() {

	go metricsCounter()
	go metricsGauge()
	go metricsMeter()

	go client.UpdatePrometheusMetrics()

	http.Handle("/metrics", promhttp.Handler())

	_ = http.ListenAndServe(":2112", nil)
}
