package extension

import (
	"golang.org/x/net/context"
	proto "github.com/nokamoto/esp4g/protobuf"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"log"
	"time"
	"github.com/golang/protobuf/ptypes/duration"
)

type accessLogService struct {
	logger *zap.SugaredLogger
}

func convert(d *duration.Duration) time.Duration {
	return time.Duration(d.Seconds) * time.Second + time.Duration(d.Nanos)
}

func (a *accessLogService)UnaryAccess(_ context.Context, unary *proto.UnaryAccessLog) (*empty.Empty, error) {
	a.logger.Infow("unary",
		"method", unary.GetMethod(),
		"status", unary.GetStatus(),
		"response_time", convert(unary.GetResponseTime()),
		"request_size", unary.GetRequestSize(),
		"response_size", unary.GetResponseSize(),
	)
	return &empty.Empty{}, nil
}

func (a *accessLogService)StreamAccess(_ context.Context, stream *proto.StreamAccessLog) (*empty.Empty, error) {
	a.logger.Infow("stream",
		"method", stream.GetMethod(),
		"status", stream.GetStatus(),
		"response_time", convert(stream.GetResponseTime()),
	)
	return &empty.Empty{}, nil
}

func NewAccessLogService(_ Config) *accessLogService {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	return &accessLogService{logger: logger.Sugar()}
}
