package esp4g_test

import (
	ping "github.com/nokamoto/esp4g/examples/ping/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type PingService struct {
	lastRequest *ping.Ping
	lastResponse *ping.Pong
}

func (p *PingService)Send(_ context.Context, req *ping.Ping) (*ping.Pong, error) {
	p.lastRequest = req
	p.lastResponse = &ping.Pong{Y: req.GetX()}
	return p.lastResponse, nil
}

func (p *PingService)Unavailable(context.Context, *ping.Ping) (*ping.Pong, error) {
	return nil, status.Error(codes.Unavailable, "unavailable")
}
