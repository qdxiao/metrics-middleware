package metric

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	ioprometheusclient "github.com/prometheus/client_model/go"
)

var (
	instanceOnce sync.Once
	instance     *customMetricsRegistry
)

func retLabels() map[string]string {
	// 设置全局的自定义的标签
	return map[string]string{
		"app":  "host.AppName()",
		"host": "host.HostName()",
	}
}

type customMetricsRegistry struct {
	*prometheus.Registry
	customLabels []*ioprometheusclient.LabelPair
}

func NewCustomMetricsRegistry() *customMetricsRegistry {
	instanceOnce.Do(func() {
		instance = &customMetricsRegistry{
			Registry: prometheus.NewRegistry(),
		}

		for s, s2 := range retLabels() {
			kCp := s
			vCp := s2
			instance.customLabels = append(instance.customLabels, &ioprometheusclient.LabelPair{
				Name:  &kCp,
				Value: &vCp,
			})
		}
	})

	return instance
}

func (c *customMetricsRegistry) Gather() ([]*ioprometheusclient.MetricFamily, error) {
	metricFamilies, err := c.Registry.Gather()

	for _, family := range metricFamilies {
		metrics := family.Metric
		for _, metric := range metrics {
			// 给每个指标添加自定义指标
			metric.Label = append(metric.Label, c.customLabels...)
		}
	}
	return metricFamilies, err
}

// A VectorOpts is a general configuration.
type VectorOpts struct {
	// 命名空间
	Namespace string
	// 系统
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}
