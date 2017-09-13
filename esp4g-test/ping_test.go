package esp4g_test

import (
	ping "github.com/nokamoto/esp4g/examples/ping/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"sync"
)

type PingService struct {
	_lastRequest *ping.Ping
	_lastResponse *ping.Pong

	mu *sync.Mutex
}

func (p *PingService)lastRequest() *ping.Ping {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p._lastRequest
}

func (p *PingService)lastResponse() *ping.Pong {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p._lastResponse
}

func (p *PingService)Send(_ context.Context, req *ping.Ping) (*ping.Pong, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p._lastRequest = req
	p._lastResponse = &ping.Pong{Y: req.GetX()}
	return p._lastResponse, nil
}

func (p *PingService)Unavailable(context.Context, *ping.Ping) (*ping.Pong, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return nil, status.Error(codes.Unavailable, "unavailable")
}
