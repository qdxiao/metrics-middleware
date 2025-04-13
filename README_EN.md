# metrics-middleware

Metrics-middleware is a monitoring metrics middleware that supports metrics monitoring for Gin HTTP servers and gRPC servers/clients. It provides a complete monitoring chain, including metrics collection at the server-side, client-side, and database level, with integrated Grafana dashboard.

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
  - [Installation](#installation)
  - [Gin Web Framework Example](#gin-web-framework-example)
  - [gRPC Server Example](#grpc-server-example)
  - [gRPC Client Example](#grpc-client-example)
- [Project Structure](#project-structure)
  - [Built-in Metrics](#built-in-metrics-pkgmetric)
  - [Middleware Implementation](#middleware-implementation)
  - [Grafana Configuration](#grafana-configuration-pkgmetricgrafana)
- [Custom Metrics](#custom-metrics)
  - [Counter](#counter)
  - [Gauge](#gauge)
  - [Histogram](#histogram)
  - [Global Labels](#global-labels)
  - [Metric Naming Convention](#metric-naming-convention)
- [Dependencies](#dependencies)
- [Contributing](#contributing)
- [License](#license)

## Features

- Support for Gin Web framework metrics monitoring
- Support for gRPC server and client metrics monitoring
- Complete monitoring chain: client, server, and database levels
- Integrated Grafana dashboard configuration
- Using Prometheus as metrics collection system
- Support for custom built-in metrics

## Quick Start

### Installation

```bash
go get github.com/qdxiao/metrics-middleware
```

### Gin Web Framework Example

```go
package main

import (
    "github.com/gin-gonic/gin"
    metric2 "github.com/qdxiao/metrics-middleware/metric"
    "github.com/qdxiao/metrics-middleware/pkg/metric"
    "github.com/qdxiao/metrics-middleware/pkg/metric/gin-metrics/ginmetrics"
)

func main() {
    r := gin.Default()

    // Register default runtime metrics
    metric.NewBuiltinMetric(
        metric2.NewCustomMetricsRegistry(),
        metric.WithEnableGo(true),
        metric.WithEnbaleProcess(true),
    )

    // Get global Monitor object
    m := ginmetrics.NewMonitor(
        "serverName",
        ginmetrics.WithMetricPath("/metrics"), // Specify metrics path
    )

    m.Use(r)

    r.GET("/product/:id", func(ctx *gin.Context) {
        ctx.JSON(200, map[string]string{
            "productId": ctx.Param("id"),
        })
    })

    _ = r.Run()
}
```

### gRPC Server Example

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    serverinterceptors "github.com/qdxiao/metrics-middleware/pkg/metric/rpc-metrics/server"
    "google.golang.org/grpc"
)

func main() {
    s := grpc.NewServer(
        // Add metrics interceptor
        grpc.UnaryInterceptor(
            serverinterceptors.UnaryPrometheusInterceptor(),
        ),
    )

    // Start Prometheus metrics service
    go func() {
        metric.NewBuiltinMetric(
            metric2.NewCustomMetricsRegistry(),
            metric.WithEnableGo(true),
            metric.WithEnbaleProcess(true),
        )

        r := gin.Default()
        r.GET("/metrics", func(ctx *gin.Context) {
            promhttp.HandlerFor(metric2.NewCustomMetricsRegistry(), promhttp.HandlerOpts{
                Registry: metric2.NewCustomMetricsRegistry(),
            }).ServeHTTP(ctx.Writer, ctx.Request)
        })
        r.Run(":8080")
    }()

    // ... Server implementation
}
```

### gRPC Client Example

```go
package main

import (
    "github.com/qdxiao/metrics-middleware/pkg/metric/rpc-metrics/client"
    "google.golang.org/grpc"
)

func main() {
    conn, err := grpc.NewClient(
        *addr,
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        // Add metrics monitoring interceptor
        grpc.WithUnaryInterceptor(client.PrometheusInterceptor()),
    )
    // ... Client implementation
}
```

## Project Structure

### Built-in Metrics (`pkg/metric/`)

1. Basic Metric Definitions (`metrics.go`)
```go
const (
    Namespace = "gfast"
    TypeHttp  = "http"
    TypeGrpc  = "grpc"
    TypeMysql = "mysql"
)
```

2. Built-in Runtime Metrics (`builtin-metrics.go`)
- Go Runtime metrics
- Process metrics
- Support for custom collectors

3. Server Metrics (`server.go`)
- Request handling time histogram: `server_handle_seconds`
    - Labels: method, type, status, peer, peer_host
- Request total counter: `server_handle_total`
    - Labels: method, type, peer, peer_host
- Client request time histogram: `client_handle_seconds`
    - Labels: method, type, server
- Client request total: `client_handle_total`
    - Labels: method, type, server

4. Database Metrics (`database.go`)
- Database operation counter: `database_lib_handle_total`
    - Labels: type, method, name, server

### Middleware Implementation

1. Gin Web Framework Middleware (`pkg/metric/gin-metrics/`)
- Automatic HTTP request metrics collection
- Custom metrics path support
- Integrated Prometheus collection endpoint

2. gRPC Middleware (`pkg/metric/rpc-metrics/`)
- Client interceptor: Automatic gRPC client call metrics collection
- Server interceptor: Automatic gRPC server request metrics collection

### Grafana Configuration (`pkg/metric/grafana/`)

1. Default Dashboard Configuration: `默认指标监控-1744569190739.json`
    - Pre-configured visualization panels
    - One-click import to Grafana

2. Dashboard Preview:
![Grafana Dashboard](pkg/metric/grafana/image.png)

The dashboard includes the following main metrics:
- Traffic distribution metrics
- Runtime system metrics
- GC-related metrics
- Service golden signals (QPS, latency, error rates, etc.)

## Custom Metrics

In addition to built-in monitoring metrics, the middleware supports using the `metric` package to create custom metrics:

### Counter

```go
counter := metric.NewCounterVec(metric.NewCustomMetricsRegistry(), &metric.CounterVecOpts{
    Namespace: "myapp",
    Subsystem: "requests",
    Name:      "total",
    Help:      "Total number of requests",
    Labels:    []string{"method", "path"},
})

// Usage
counter.Inc("GET", "/api/users")  // Increment by 1
counter.Add(2, "POST", "/api/users")  // Add specific value
```

### Gauge

```go
gauge := metric.NewGaugeVec(metric.NewCustomMetricsRegistry(), &metric.GaugeVecOpts{
    Namespace: "myapp",
    Subsystem: "resources",
    Name:      "connections",
    Help:      "Number of active connections",
    Labels:    []string{"type"},
})

// Usage
gauge.Set(100, "tcp")  // Set value
gauge.Inc("tcp")       // Increment by 1
gauge.Add(10, "tcp")   // Add specific value
```

### Histogram

```go
histogram := metric.NewHistogramVec(metric.NewCustomMetricsRegistry(), &metric.HistogramVecOpts{
    Namespace: "myapp",
    Subsystem: "requests",
    Name:      "duration_ms",
    Help:      "Request duration in milliseconds",
    Labels:    []string{"method"},
    Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000}, // Custom distribution intervals
})

// Usage
histogram.Observe(42, "GET")  // Record an observation
```

### Global Labels

All metrics created through `NewCustomMetricsRegistry()` will automatically add the following global labels:
- `app`: Application name
- `host`: Hostname

### Metric Naming Convention

- Namespace: Usually the application name
- Subsystem: Represents subsystem or module
- Name: Specific metric name
- Labels: For multi-dimensional data analysis

## Dependencies

- gin-gonic/gin: Web framework
- prometheus/client_golang: Prometheus client
- google.golang.org/grpc: gRPC framework

## Contributing

If you find any issues or have suggestions for improvements, you're welcome to:

1. Submit an Issue: Describe your problem or suggestion
2. Submit a Pull Request:
   - Fork this repository
   - Create your feature branch (`git checkout -b feature/AmazingFeature`)
   - Commit your changes (`git commit -m 'Add some AmazingFeature'`)
   - Push to the branch (`git push origin feature/AmazingFeature`)
   - Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 