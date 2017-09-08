package main

import (
	"google.golang.org/grpc"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"io"
)

type ProxyServer interface {
	Proxy(ctx context.Context, method string, req *ProxyMessage) (interface{}, error)

	ProxyClientSideStreaming(method string, stream *clientSideServerStream, desc *grpc.StreamDesc) error

	ProxyServerSideStreaming(method string, req *ProxyMessage, stream *serverSideServerStream, desc *grpc.StreamDesc) error
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

func (p *proxyServer)ProxyClientSideStreaming(method string, stream *clientSideServerStream, desc *grpc.StreamDesc) error {
	log.Printf("%s", method)

	cs, err := grpc.NewClientStream(context.Background(), desc, p.con, method)
	if err != nil {
		return err
	}

	proxy := clientSideClientStream{cs}

	for {
		m, err := stream.Recv()
		if err == io.EOF {
			res, err := proxy.CloseAndRecv()
			if err != nil {
				return err
			}
			stream.SendAndClose(res)
			break
		}
		if err != nil {
			return err
		}
		if err = proxy.Send(m); err != nil {
			return err
		}
	}

	return nil
}

func (p *proxyServer)ProxyServerSideStreaming(method string, req *ProxyMessage, stream *serverSideServerStream, desc *grpc.StreamDesc) error {
	log.Printf("%s", method)

	cs, err := grpc.NewClientStream(context.Background(), desc, p.con, method)
	if err != nil {
		return err
	}

	proxy := serverSideClientStream{cs}
	if err = proxy.ClientStream.SendMsg(req); err != nil {
		return err
	}
	if err = proxy.ClientStream.CloseSend(); err != nil {
		return err
	}

	for {
		m, err := proxy.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err = stream.Send(m); err != nil {
			return err
		}
	}

	return nil
}

func NewProxyServer(port int) (*proxyServer, error) {
	opts := []grpc.DialOption{grpc.WithInsecure()}

	con, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), opts...)
	if err != nil {
		return nil, err
	}

	return &proxyServer{con: con}, nil
}
