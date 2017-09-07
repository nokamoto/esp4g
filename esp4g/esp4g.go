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
	pb = flag.String("descriptor", "descriptor.pb", "FileDescriptorSet protocol buffer file")
	port = flag.Int("p", 9000, "The gRPC proxy server port")
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

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	proxy, err := NewProxyServer(*port)
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
