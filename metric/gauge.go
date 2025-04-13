package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

type (
	// GaugeVecOpts is an alias of VectorOpts.
	GaugeVecOpts VectorOpts

	// GaugeVec represents a gauge vector.
	GaugeVec interface {
		// Set sets v to labels.
		Set(v float64, labels ...string)
		// Inc increments labels.
		Inc(labels ...string)
		// Add adds v to labels.
		Add(v float64, labels ...string)
	}

	promGaugeVec struct {
		gauge *prom.GaugeVec
	}
)

// NewGaugeVec returns a GaugeVec.
func NewGaugeVec(reg prom.Registerer, cfg *GaugeVecOpts) GaugeVec {
	if cfg == nil {
		return nil
	}

	vec := prom.NewGaugeVec(
		prom.GaugeOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.Name,
			Help:      cfg.Help,
		}, cfg.Labels)

	reg.MustRegister(vec)

	gv := &promGaugeVec{
		gauge: vec,
	}

	return gv
}

func (gv *promGaugeVec) Inc(labels ...string) {
	gv.gauge.WithLabelValues(labels...).Inc()
}

func (gv *promGaugeVec) Add(v float64, labels ...string) {
	gv.gauge.WithLabelValues(labels...).Add(v)
}

func (gv *promGaugeVec) Set(v float64, labels ...string) {
	gv.gauge.WithLabelValues(labels...).Set(v)
}
