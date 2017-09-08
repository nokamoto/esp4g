package main

import "google.golang.org/grpc"

type clientSideServerStream struct {
	grpc.ServerStream
}

func (c *clientSideServerStream) SendAndClose(m *ProxyMessage) error {
	return c.ServerStream.SendMsg(m)
}

func (c *clientSideServerStream) Recv() (*ProxyMessage, error) {
	m := NewProxyMessage()
	if err := c.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

type clientSideClientStream struct {
	grpc.ClientStream
}

func (c *clientSideClientStream)Send(m *ProxyMessage) error {
	return c.ClientStream.SendMsg(m)
}

func (c *clientSideClientStream)CloseAndRecv() (*ProxyMessage, error) {
	if err := c.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := NewProxyMessage()
	if err := c.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}
