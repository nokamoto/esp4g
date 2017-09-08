package main

import (
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"time"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"github.com/golang/protobuf/ptypes/duration"
	extension "github.com/nokamoto/esp4g/protobuf"
	"log"
	"fmt"
)

type accessLogInterceptor struct {
	con *grpc.ClientConn
}

func NewAccessLogInterceptor(port *int) *accessLogInterceptor {
	opts := []grpc.DialOption{grpc.WithInsecure()}

	con, err := grpc.Dial(fmt.Sprintf("localhost:%d", *port), opts...)
	if err != nil {
		log.Fatal(err)
	}

	return &accessLogInterceptor{con: con}
}

func convert(d time.Duration) *duration.Duration {
	return &duration.Duration{
		Seconds: d.Nanoseconds() / time.Second.Nanoseconds(),
		Nanos: int32(d.Nanoseconds() % time.Second.Nanoseconds()),
	}
}

func (a *accessLogInterceptor)doAccessLog(method string, responseTime time.Duration, stat codes.Code, in int, out int) error {
	client := extension.NewAccessLogServiceClient(a.con)
	unary := extension.UnaryAccessLog{
		Method: method,
		ResponseTime: convert(responseTime),
		Status: stat.String(),
		RequestSize: int64(in),
		ResponseSize: int64(out),
	}
	_, err := client.UnaryAccess(context.Background(), &unary)
	return err
}

func (a *accessLogInterceptor)doStreamAccessLog(method string, responseTime time.Duration, stat codes.Code) error {
	client := extension.NewAccessLogServiceClient(a.con)
	stream := extension.StreamAccessLog{
		Method: method,
		ResponseTime: convert(responseTime),
		Status: stat.String(),
	}
	_, err := client.StreamAccess(context.Background(), &stream)
	return err
}


func (a *accessLogInterceptor)createAccessLogInterceptor(next *grpc.UnaryServerInterceptor) *grpc.UnaryServerInterceptor {
	f := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		var res interface{}
		var err error

		if next != nil {
			res, err = (*next)(ctx, req, info, handler)
		} else {
			res, err = handler(ctx, req)
		}

		elapsed := time.Since(start)

		code := codes.Unknown
		if stat, ok := status.FromError(err); ok {
			code = stat.Code()
		}

		inBytes := -1
		outBytes := -1
		if m, ok := req.(*ProxyMessage); ok {
			inBytes = len(m.bytes)
		}
		if m, ok := res.(*ProxyMessage); ok {
			outBytes = len(m.bytes)
		}

		if skipErr := a.doAccessLog(info.FullMethod, elapsed, code, inBytes, outBytes); skipErr != nil {
			log.Println(skipErr)
		}

		return res, err
	}

	i := grpc.UnaryServerInterceptor(f)

	return &i
}

func (a *accessLogInterceptor)createStreamAccessLogInterceptor(next *grpc.StreamServerInterceptor) *grpc.StreamServerInterceptor {
	f := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		var err error
		if next != nil {
			err = (*next)(srv, ss, info, handler)
		} else {
			err = handler(srv, ss)
		}

		elapsed := time.Since(start)

		code := codes.Unknown
		if stat, ok := status.FromError(err); ok {
			code = stat.Code()
		}

		if skipErr := a.doStreamAccessLog(info.FullMethod, elapsed, code); skipErr != nil {
			log.Println(skipErr)
		}

		return err
	}

	i := grpc.StreamServerInterceptor(f)

	return &i
}
