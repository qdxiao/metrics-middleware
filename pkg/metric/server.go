package metric

import "github.com/qdxiao/metrics-middleware/metric"

// 参数说明：
// method: 请求方式+请求方法
// type: 协议类型：http grpc
// status: 状态码
// host: 服务地址
// app: 服务名称
// server: 对端的服务地址
// peer: 对端的服务名称
// peer_host: 对端的服务 ip

// '两个基本指标：
// '1.每个请求的耗时
// '2.每个请求的状态计数器

// 接口维度
var (
	ServerHandleHistogram = metric.NewHistogramVec(metric.NewCustomMetricsRegistry(), &metric.HistogramVecOpts{
		Namespace: Namespace,
		Subsystem: "request",
		Name:      "server_handle_seconds",
		Help:      "server handle seconds.",
		Labels:    []string{"method", "type", "status", "peer", "peer_host"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	ServerHandleCounter = metric.NewCounterVec(metric.NewCustomMetricsRegistry(), &metric.CounterVecOpts{
		Namespace: Namespace,
		Subsystem: "request",
		Name:      "server_handle_total",
		Help:      "server handle total.",
		Labels:    []string{"method", "type", "peer", "peer_host"},
	})

	ClientHandleHistogram = metric.NewHistogramVec(metric.NewCustomMetricsRegistry(), &metric.HistogramVecOpts{
		Namespace: Namespace,
		Subsystem: "request",
		Name:      "client_handle_seconds",
		Help:      "client handle seconds.",
		Labels:    []string{"method", "type", "server"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	ClientHandleCounter = metric.NewCounterVec(metric.NewCustomMetricsRegistry(), &metric.CounterVecOpts{
		Namespace: Namespace,
		Subsystem: "request",
		Name:      "client_handle_total",
		Help:      "client handle total.",
		Labels:    []string{"method", "type", "server"},
	})
)
