package main

import (
	"github.com/gin-gonic/gin"
	metric2 "github.com/qdxiao/metrics-middleware/metric"
	"github.com/qdxiao/metrics-middleware/pkg/metric"
	"github.com/qdxiao/metrics-middleware/pkg/metric/gin-metrics/ginmetrics"
)

func main() {
	r := gin.Default()

	// 注册默认 runtime 指标
	metric.NewBuiltinMetric(
		metric2.NewCustomMetricsRegistry(),
		metric.WithEnableGo(true),
		metric.WithEnbaleProcess(true),
	)

	// get global Monitor object
	m := ginmetrics.NewMonitor(
		"serverName",
		ginmetrics.WithMetricPath("/metrics"), // 指定的 metrics 路径
	)

	m.Use(r)

	r.GET("/product/:id", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]string{
			"productId": ctx.Param("id"),
		})
	})

	_ = r.Run()
}
