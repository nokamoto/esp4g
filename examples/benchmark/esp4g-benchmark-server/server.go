package main

import (
	"flag"
	"net"
	"fmt"
	"log"
	benchmark "github.com/nokamoto/esp4g/examples/benchmark/protobuf"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"github.com/golang/protobuf/ptypes/empty"
)

type UnaryServer struct {}

func (UnaryServer)Send(_ context.Context, unary *benchmark.Unary) (*benchmark.Unary, error) {
	return unary, nil
}

func (UnaryServer)Check(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

var (
	port = flag.Int("p", 8000, "The gRPC server port")
)

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("listen %v port", *port)
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	benchmark.RegisterUnaryServiceServer(server, UnaryServer{})
	benchmark.RegisterHealthCheckServiceServer(server, UnaryServer{})

	log.Println("start esp4g-benchmark-server...")
	server.Serve(lis)
}
