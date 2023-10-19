package main

import (
	"runtime/metrics"
	"testing"
	_ "unsafe"
)

// 假设函数a已经被linkname到了b, 那么如果还需要在别的包中使用a, 需要linkname到b
//
//go:linkname readMetricsName runtime/metrics_test.runtime_readMetricNames
func readMetricsName() []string

func TestMetrics(t *testing.T) {

	// 获取所有的metrics name
	desc := metrics.All()
	var samples = make([]metrics.Sample, len(desc))
	for i := range samples {
		samples[i].Name = desc[i].Name
	}

	// 读取metrics
	metrics.Read(samples)

	for _, s := range samples {

		switch s.Value.Kind() {
		case metrics.KindUint64:
			t.Logf("name:%v, value:%v", s.Name, s.Value.Uint64())
		case metrics.KindFloat64:
			t.Logf("name:%v, value:%v", s.Name, s.Value.Float64())
		case metrics.KindFloat64Histogram:
			t.Logf("name:%v, not support", s.Name)
		default:
			t.Logf("name:%v, value:%v", s.Name, s.Value)
		}
	}
}
