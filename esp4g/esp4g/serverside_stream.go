package esp4g

import "google.golang.org/grpc"

type serverSideServerStream struct {
	grpc.ServerStream
}

func (s *serverSideServerStream)Send(m *proxyMessage) error {
	return s.ServerStream.SendMsg(m)
}

type serverSideClientStream struct {
	grpc.ClientStream
}

func (x *serverSideClientStream) Recv() (*proxyMessage, error) {
	m := newProxyMessage()
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}
