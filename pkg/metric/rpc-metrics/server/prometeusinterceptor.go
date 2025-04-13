package serverinterceptors

import (
	"context"
	"strconv"
	"time"

	"github.com/qdxiao/metrics-middleware/pkg/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnaryPrometheusInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := time.Now()

		resp, err = handler(ctx, req)

		// 记录耗时
		metric.ServerHandleHistogram.Observe(int64(time.Since(startTime)/time.Millisecond), info.FullMethod, metric.TypeGrpc, strconv.Itoa(int(status.Code(err))), "peer", "peer_host")
		// 记录请求状态
		metric.ServerHandleCounter.Inc(info.FullMethod, metric.TypeGrpc, "peer", "peer_host")
		return resp, err
	}
}
