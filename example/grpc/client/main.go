package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/qdxiao/metrics-middleware/example/grpc/proto"
	"github.com/qdxiao/metrics-middleware/pkg/metric/rpc-metrics/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", "defaultName", "Name to greet")
)

func main() {
	flag.Parse()

	conn, err := grpc.NewClient(
		*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// 指标监控拦截器
		grpc.WithUnaryInterceptor(client.PrometheusInterceptor()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		panic(err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
