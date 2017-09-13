package esp4g_test

import (
	"golang.org/x/net/context"
	"github.com/golang/protobuf/ptypes/empty"
)

type HealthCheckService struct {}

func (HealthCheckService)Check(_ context.Context, _ *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
