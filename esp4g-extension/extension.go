package main

import (
	"flag"
	"log"
	"google.golang.org/grpc"
	"net"
	"fmt"
	extension "github.com/nokamoto/esp4g/protobuf"
)

var (
	port = flag.Int("p", 10000, "The gRPC server port")
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

	extension.RegisterAccessControlServiceServer(server, NewAccessControlService())
	extension.RegisterAccessLogServiceServer(server, NewAccessLogService())

	log.Println("start esp server...")
	server.Serve(lis)
}
