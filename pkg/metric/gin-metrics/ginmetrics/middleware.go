package ginmetrics

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metric2 "github.com/qdxiao/metrics-middleware/metric"
	"github.com/qdxiao/metrics-middleware/pkg/metric"
)

// Use gin-metrics中间件
func (m *Monitor) Use(r gin.IRoutes) {

	r.Use(m.monitorInterceptor)
	r.GET(m.metricPath, func(ctx *gin.Context) {
		promhttp.HandlerFor(
			metric2.NewCustomMetricsRegistry(),
			promhttp.HandlerOpts{
				Registry: metric2.NewCustomMetricsRegistry(),
			},
		).ServeHTTP(ctx.Writer, ctx.Request)
	})
}

// UseWithoutExposingEndpoint 用于将监视器拦截器添加到 gin 路由器。可以多次调用它来拦截多个杜松子酒。IRoutes http 路径未设置，为此请使用 Expose 函数
func (m *Monitor) UseWithoutExposingEndpoint(r gin.IRoutes) {
	r.Use(m.monitorInterceptor)
}

// Expose 为给定的路由器添加度量路径。路由器可以与传递给 UseWithoutExposingEndpoint 的路由器不同。这允许在不同端口上公开指标。
func (m *Monitor) Expose(r gin.IRoutes) {
	r.GET(m.metricPath, func(ctx *gin.Context) {
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	})
}

func (m *Monitor) includesMetadata() bool {
	return len(m.metadata) > 0
}

func (m *Monitor) getMetadata() ([]string, []string) {
	var metadataLabels []string
	var metadataValues []string

	for v := range m.metadata {
		metadataLabels = append(metadataLabels, v)
		metadataValues = append(metadataValues, m.metadata[v])
	}

	return metadataLabels, metadataValues
}

// monitorInterceptor 作为 gin 监控中间件
func (m *Monitor) monitorInterceptor(ctx *gin.Context) {
	// 排除不应该监控的路径
	if ctx.Request.URL.Path == m.metricPath ||
		slices.Contains(m.excludePaths, ctx.Request.URL.Path) {
		ctx.Next()
		return
	}
	startTime := time.Now()

	// execute normal process.
	ctx.Next()
	// after request
	m.ginMetricHandle(ctx, startTime)
}

func (m *Monitor) ginMetricHandle(ctx *gin.Context, start time.Time) {
	r := ctx.Request
	w := ctx.Writer

	metric.ServerHandleHistogram.Observe(int64(time.Since(start)/time.Millisecond), fmt.Sprintf("%s->%s", r.Method, ctx.FullPath()), metric.TypeHttp, strconv.Itoa(w.Status()), "", "")
	metric.ServerHandleCounter.Inc(fmt.Sprintf("%s->%s", r.Method, ctx.FullPath()), metric.TypeHttp, "", "")
}
