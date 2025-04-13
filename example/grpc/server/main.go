package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	pb "github.com/qdxiao/metrics-middleware/example/grpc/proto"
	metric2 "github.com/qdxiao/metrics-middleware/metric"
	"github.com/qdxiao/metrics-middleware/pkg/metric"
	serverinterceptors "github.com/qdxiao/metrics-middleware/pkg/metric/rpc-metrics/server"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello" + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		// 添加指标拦截器
		grpc.UnaryInterceptor(
			serverinterceptors.UnaryPrometheusInterceptor(),
		),
	)

	pb.RegisterGreeterServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	go func() {
		// 添加一些 runtime 指标
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

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
