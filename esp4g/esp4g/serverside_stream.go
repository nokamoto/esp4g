package esp4g

import "google.golang.org/grpc"

type serverSideServerStream struct {
	grpc.ServerStream
}

func (s *serverSideServerStream)Send(m *ProxyMessage) error {
	return s.ServerStream.SendMsg(m)
}

type serverSideClientStream struct {
	grpc.ClientStream
}

func (x *serverSideClientStream) Recv() (*ProxyMessage, error) {
	m := NewProxyMessage()
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}
