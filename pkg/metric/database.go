package metric

import "github.com/qdxiao/metrics-middleware/metric"

// 数据库层面，目前暂时内置 mysql,redis 两种

var (
	MetricLibraryCounter = metric.NewCounterVec(metric.NewCustomMetricsRegistry(), &metric.CounterVecOpts{
		Namespace: Namespace,
		Subsystem: "database",
		Name:      "lib_handle_total",
		Help:      "database situation",
		Labels:    []string{"type", "method", "name", "server"},
	})
)
