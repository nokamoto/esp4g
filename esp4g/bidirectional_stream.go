package main

import "google.golang.org/grpc"

type bidirectionalServerStream struct {
	grpc.ServerStream
}

func (b *bidirectionalServerStream) Send(m *ProxyMessage) error {
	return b.ServerStream.SendMsg(m)
}

func (b *bidirectionalServerStream) Recv() (*ProxyMessage, error) {
	m := NewProxyMessage()
	if err := b.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

type bidirectionalClientStream struct {
	grpc.ClientStream
}

func (b *bidirectionalClientStream)Send(m *ProxyMessage) error {
	return b.ClientStream.SendMsg(m)
}

func (b *bidirectionalClientStream)Recv() (*ProxyMessage, error) {
	m := NewProxyMessage()
	if err := b.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}
