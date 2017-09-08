package main

import (
	"flag"
	"io/ioutil"
	"log"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/grpc"
	"fmt"
	"net"
)

var (
	pb = flag.String("d", "descriptor.pb", "FileDescriptorSet protocol buffer file")
	port = flag.Int("p", 9000, "The gRPC server port")
	proxy = flag.Int("proxy", 8000, "The gRPC proxy port")
	accessLog = flag.Int("log", 10000, "The gRPC access log service port")
	accessControl = flag.Int("control", 10000, "The gRPC access control service port")
)

func main() {
	flag.Parse()

	data, err := ioutil.ReadFile(*pb)
	if err != nil {
		log.Fatal(err)
	}

	fds := &descriptor.FileDescriptorSet{}

	proto.Unmarshal(data, fds)

	services := CreateProxyServiceDesc(fds)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("listen %v port", *port)
	}

	logInterceptor := NewAccessLogInterceptor(accessLog)
	controlInterceptor := NewAccessControlInterceptor(accessControl)

	opts := []grpc.ServerOption{}

	{
		f := controlInterceptor.createApiKeyInterceptor(nil)
		g := logInterceptor.createAccessLogInterceptor(f)
		opts = append(opts, grpc.UnaryInterceptor(*g))
	}

	{
		f := controlInterceptor.createStreamApiKeyInterceptor(nil)
		g := logInterceptor.createStreamAccessLogInterceptor(f)
		opts = append(opts, grpc.StreamInterceptor(*g))
	}

	server := grpc.NewServer(opts...)

	proxy, err := NewProxyServer(*proxy)
	if err != nil {
		log.Fatal(err)
	}

	for _, desc := range services {
		log.Printf("register %v", desc)
		server.RegisterService(&desc, proxy)
	}

	log.Println("start esp server...")
	server.Serve(lis)
}
