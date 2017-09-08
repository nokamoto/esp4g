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

	ProxyBidirectionalStreaming(method string, stream *bidirectionalServerStream, desc *grpc.StreamDesc) error
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

func (p *proxyServer)ProxyBidirectionalStreaming(method string, stream *bidirectionalServerStream, desc *grpc.StreamDesc) error {
	log.Printf("%s", method)

	cs, err := grpc.NewClientStream(context.Background(), desc, p.con, method)
	if err != nil {
		return err
	}

	proxy := bidirectionalClientStream{cs}

	downstream := make(chan error, 1)
	upstream := make(chan error, 1)

	go func() {
		for {
			m, err := stream.Recv()
			if err != nil {
				downstream <- err
				break
			}
			if err = proxy.Send(m); err != nil {
				downstream <- err
				break
			}
		}
	}()

	go func() {
		for {
			m, err := proxy.Recv()
			if err != nil {
				upstream <- err
				break
			}
			if err = stream.Send(m); err != nil {
				downstream <- err
				break
			}
		}
	}()

	select {
	case err := <- downstream:
		if err == io.EOF {
			if err = proxy.CloseSend(); err != nil {
				return err
			}
		} else if err != nil {
			if err = proxy.CloseSend(); err != nil {
				return nil
			}
		}

	case err := <- upstream:
		if err != io.EOF {
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
