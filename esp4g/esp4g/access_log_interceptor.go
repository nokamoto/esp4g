package esp4g

import (
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"time"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	proto "github.com/nokamoto/esp4g/protobuf"
	"github.com/nokamoto/esp4g/esp4g-utils"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/nokamoto/esp4g/esp4g-extension"
)

type accessLogInterceptor struct {
	con *grpc.ClientConn

	service *extension.AccessLogService
}

func newAccessLogInterceptor(address string, fds *descriptor.FileDescriptorSet, yml string) *accessLogInterceptor {
	if len(address) != 0 {
		opts := []grpc.DialOption{grpc.WithInsecure()}

		con, err := grpc.Dial(address, opts...)
		if err != nil {
			utils.Logger.Fatalw("failed to create gRPC dial", "err", err)
		}

		return &accessLogInterceptor{con: con}
	}

	buf, err := ioutil.ReadFile(yml)
	if err != nil {
		utils.Logger.Fatalw("failed to read yaml", "yaml", yml, "err", err)
	}

	var config extension.Config
	if err = yaml.Unmarshal(buf, &config); err != nil {
		utils.Logger.Fatalw("failed to unmarshal", "err", err)
	}

	return &accessLogInterceptor{service:extension.NewAccessLogService(config, fds)}
}

func (a *accessLogInterceptor)doAccessLog(method string, responseTime time.Duration, stat codes.Code, in int, out int) error {
	unary := proto.UnaryAccessLog{
		Method: method,
		ResponseTime: utils.ConvertDuration(responseTime),
		Status: stat.String(),
		RequestSize: int64(in),
		ResponseSize: int64(out),
	}

	var err error

	if a.con != nil {
		client := proto.NewAccessLogServiceClient(a.con)
		_, err = client.UnaryAccess(context.Background(), &unary)
	}

	if a.service != nil {
		_, err = a.service.UnaryAccess(context.Background(), &unary)
	}

	return err
}

func (a *accessLogInterceptor)doStreamAccessLog(method string, responseTime time.Duration, stat codes.Code) error {
	stream := proto.StreamAccessLog{
		Method: method,
		ResponseTime: utils.ConvertDuration(responseTime),
		Status: stat.String(),
	}
	var err error

	if a.con != nil {
		client := proto.NewAccessLogServiceClient(a.con)
		_, err = client.StreamAccess(context.Background(), &stream)
	}

	if a.service != nil {
		_, err = a.service.StreamAccess(context.Background(), &stream)
	}
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
		if m, ok := req.(*proxyMessage); ok {
			inBytes = len(m.bytes)
		}
		if m, ok := res.(*proxyMessage); ok {
			outBytes = len(m.bytes)
		}

		if skipErr := a.doAccessLog(info.FullMethod, elapsed, code, inBytes, outBytes); skipErr != nil {
			utils.Logger.Infow("access log failed", "err", skipErr)
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
			utils.Logger.Infow("access log failed", "err", skipErr)
		}

		return err
	}

	i := grpc.StreamServerInterceptor(f)

	return &i
}
