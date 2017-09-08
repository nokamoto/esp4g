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
	"time"
	"google.golang.org/grpc/codes"
)

var (
	pb = flag.String("d", "descriptor.pb", "FileDescriptorSet protocol buffer file")
	port = flag.Int("p", 9000, "The gRPC server port")
	proxy = flag.Int("proxy", 8000, "The gRPC proxy port")
)

func doAccessLog(method string, responseTime time.Duration, stat codes.Code, in int, out int) {
	fmt.Println(method, responseTime, stat, in, out)
}

func doStreamAccessLog(method string, responseTime time.Duration, stat codes.Code) {
	fmt.Println(method, responseTime, stat)
}

func doAccessControl(method string, keys []string) Policy {
	return ALLOW
}

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

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(*createAccessLogInterceptor(doAccessLog, createApiKeyInterceptor(doAccessControl, nil))),
		grpc.StreamInterceptor(*createStreamAccessLogInterceptor(doStreamAccessLog, createStreamApiKeyInterceptor(doAccessControl, nil))),
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
