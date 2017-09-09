package esp4g

import "google.golang.org/grpc"

type bidirectionalServerStream struct {
	grpc.ServerStream
}

func (b *bidirectionalServerStream) Send(m *proxyMessage) error {
	return b.ServerStream.SendMsg(m)
}

func (b *bidirectionalServerStream) Recv() (*proxyMessage, error) {
	m := newProxyMessage()
	if err := b.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

type bidirectionalClientStream struct {
	grpc.ClientStream
}

func (b *bidirectionalClientStream)Send(m *proxyMessage) error {
	return b.ClientStream.SendMsg(m)
}

func (b *bidirectionalClientStream)Recv() (*proxyMessage, error) {
	m := newProxyMessage()
	if err := b.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}
