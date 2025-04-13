package ginmetrics

type MetricType int

const (
	defaultMetricPath = "/debug/metrics"
)

var (
	defaultExcludePaths []string
	//monitor             *Monitor
)

type Options func(monitor *Monitor)

// Monitor 是一个用于设置gin服务器监视器的对象
type Monitor struct {
	metricPath   string
	excludePaths []string
	metadata     map[string]string
	serverName   string
}

// NewMonitor 用于获取全局Monitor对象，此函数返回一个单例对象。
func NewMonitor(serverName string, opts ...Options) *Monitor {
	monitor := &Monitor{
		metricPath:   defaultMetricPath,
		excludePaths: defaultExcludePaths,
		metadata:     make(map[string]string),
		serverName:   serverName,
	}
	for _, opt := range opts {
		opt(monitor)
	}
	return monitor
}

// WithMetricPath 设置 metricPath 属性。metricPath用于 Prometheus 获取 gin 服务器监控数据。
func WithMetricPath(path string) Options {
	return func(monitor *Monitor) {
		monitor.metricPath = path
	}
}

// WithExcludePaths 设置排除不应报告的路径（例如/ping、/healthz…）
func WithExcludePaths(paths []string) Options {
	return func(monitor *Monitor) {
		copy(paths, monitor.excludePaths)
	}
}
