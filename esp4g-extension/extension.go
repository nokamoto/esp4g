package main

import (
	"flag"
	"log"
	"google.golang.org/grpc"
	"net"
	"fmt"
	extension "github.com/nokamoto/esp4g/protobuf"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

var (
	yml = flag.String("c", "./config.yaml", "The application config file path")
)


func main() {
	flag.Parse()

	buf, err := ioutil.ReadFile(*yml)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	if err = yaml.Unmarshal(buf, &config); err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Server.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("listen %v port", config.Server.Port)
	}

	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)

	extension.RegisterAccessControlServiceServer(server, NewAccessControlService(config))
	extension.RegisterAccessLogServiceServer(server, NewAccessLogService(config))

	log.Println("start esp server...")
	server.Serve(lis)
}
