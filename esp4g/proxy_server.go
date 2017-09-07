package main

import (
	"google.golang.org/grpc"
	"fmt"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type ProxyServer interface {
	Proxy() (interface{}, error)
}

type proxyServer struct {
	con *grpc.ClientConn
}

func (*proxyServer)Proxy() (interface{}, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented yet")
}

func NewProxyServer(port int) (*proxyServer, error) {
	opts := []grpc.DialOption{grpc.WithInsecure()}

	con, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), opts...)
	if err != nil {
		return nil, err
	}

	return &proxyServer{con: con}, nil
}
