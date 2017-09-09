package main

import (
	"flag"
	"log"
	"net"
	"fmt"
	"github.com/nokamoto/esp4g/esp4g-extension/esp4g-extension"
)


func main() {
	var (
		port = flag.Int("p", 10000, "The gRPC server port")
		yml = flag.String("c", "./config.yaml", "The application config file path")
	)

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Printf("listen %v port", *port)
	}

	server := extension.NewGrpcServer(*yml)

	log.Println("start esp extension server...")
	server.Serve(lis)
}
