package main

import (
	"google.golang.org/grpc"
	"fmt"
	"golang.org/x/net/context"
	"log"
)

type ProxyServer interface {
	Proxy(ctx context.Context, method string, req *ProxyMessage) (interface{}, error)
}

type proxyServer struct {
	con *grpc.ClientConn
}

func (p *proxyServer)Proxy(ctx context.Context, method string, req *ProxyMessage) (interface{}, error) {
	log.Printf("%s", method)
	out := NewProxyMessage()
	err := grpc.Invoke(ctx, method, req, out, p.con)
	return out, err
}

func NewProxyServer(port int) (*proxyServer, error) {
	opts := []grpc.DialOption{grpc.WithInsecure()}

	con, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), opts...)
	if err != nil {
		return nil, err
	}

	return &proxyServer{con: con}, nil
}
