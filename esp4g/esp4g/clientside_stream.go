package esp4g

import "google.golang.org/grpc"

type clientSideServerStream struct {
	grpc.ServerStream
}

func (c *clientSideServerStream) SendAndClose(m *proxyMessage) error {
	return c.ServerStream.SendMsg(m)
}

func (c *clientSideServerStream) Recv() (*proxyMessage, error) {
	m := newProxyMessage()
	if err := c.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

type clientSideClientStream struct {
	grpc.ClientStream
}

func (c *clientSideClientStream)Send(m *proxyMessage) error {
	return c.ClientStream.SendMsg(m)
}

func (c *clientSideClientStream)CloseAndRecv() (*proxyMessage, error) {
	if err := c.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := newProxyMessage()
	if err := c.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}
