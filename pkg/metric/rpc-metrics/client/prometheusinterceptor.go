package client

import (
	"context"
	"time"

	"github.com/qdxiao/metrics-middleware/pkg/metric"
	"google.golang.org/grpc"
)

func PrometheusInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		startTime := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		// 记录耗时
		metric.ClientHandleHistogram.Observe(int64(time.Since(startTime)/time.Millisecond), method, metric.TypeGrpc, cc.Target())
		// 记录请求状态
		metric.ClientHandleCounter.Inc(method, metric.TypeGrpc, cc.Target())
		return err
	}
}
