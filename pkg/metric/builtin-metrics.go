package metric

import (
	"regexp"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// 指标监控
type builtinMetric struct {
	enableGo      bool
	enbaleProcess bool

	// 自定义内置指标
	coll []prom.Collector

	pidFn        func() (int, error)
	nameSpace    string
	reportErrors bool
}

type Option func(metric *builtinMetric)

func NewBuiltinMetric(reg prom.Registerer, opts ...Option) *builtinMetric {
	bm := &builtinMetric{
		enableGo:      false,
		enbaleProcess: false,
		nameSpace:     "gfast",
	}
	for _, opt := range opts {
		opt(bm)
	}
	var coll []prom.Collector
	// 添加用户传入的指标
	coll = append(coll, bm.coll...)

	if bm.enableGo {
		coll = append(coll, collectors.NewGoCollector(
			collectors.WithGoCollectorRuntimeMetrics(
				collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/.*")},
			),
		))
	}

	if bm.enbaleProcess {
		coll = append(coll, collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
			PidFn:        bm.pidFn,
			Namespace:    bm.nameSpace,
			ReportErrors: bm.reportErrors,
		}))
	}

	reg.MustRegister(
		coll...,
	)

	return bm
}
func WithEnableGo(enableGo bool) Option {
	return func(metric *builtinMetric) {
		metric.enableGo = enableGo
	}
}
func WithEnbaleProcess(enbaleProcess bool) Option {
	return func(metric *builtinMetric) {
		metric.enbaleProcess = enbaleProcess
	}
}
func WithColl(coll []prom.Collector) Option {
	return func(metric *builtinMetric) {
		copy(coll, metric.coll)
	}
}
