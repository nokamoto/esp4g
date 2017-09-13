package esp4g

import (
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	proto "github.com/nokamoto/esp4g/protobuf"
	"github.com/nokamoto/esp4g/esp4g-utils"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/nokamoto/esp4g/esp4g-extension"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/nokamoto/esp4g/esp4g-extension/config"
)

const API_KEY_HEADER = "x-api-key"

type accessControlInterceptor struct {
	con *grpc.ClientConn
	service *extension.AccessControlService
}

func newAccessControlInterceptor(address string, fds *descriptor.FileDescriptorSet, yml string) *accessControlInterceptor {
	if len(address) != 0 {
		opts := []grpc.DialOption{grpc.WithInsecure()}

		con, err := grpc.Dial(address, opts...)
		if err != nil {
			utils.Logger.Fatalw("failed to create gRPC dial", "err", err)
		}

		return &accessControlInterceptor{con: con}
	}

	buf, err := ioutil.ReadFile(yml)
	if err != nil {
		utils.Logger.Fatalw("failed to read yaml", "yaml", yml, "err", err)
	}

	var cfg config.ExtensionConfig
	if err = yaml.Unmarshal(buf, &cfg); err != nil {
		utils.Logger.Fatalw("failed to unmarshal", "err", err)
	}

	return &accessControlInterceptor{service: extension.NewAccessControlService(cfg, fds)}
}

func (a *accessControlInterceptor)callAccessControl(method string, keys []string) (proto.AccessPolicy, error) {
	id := proto.AccessIdentity{
		Method: method,
		ApiKey: keys,
	}

	var ctl *proto.AccessControl
	var err error

	if a.con != nil {
		client := proto.NewAccessControlServiceClient(a.con)
		ctl, err = client.Access(context.Background(), &id)
	}

	if a.service != nil {
		ctl, err = a.service.Access(context.Background(), &id)
	}

	if err != nil {
		return proto.AccessPolicy_DENY, err
	}
	return ctl.Policy, nil
}

func fetchApiKey(ctx context.Context) ([]string, error) {
	md, err := getMetadata(ctx)
	if err != nil {
		return nil, err
	}
	return safeMetadata(md, API_KEY_HEADER), nil
}

func (a *accessControlInterceptor)accessControl(ctx context.Context, method string) error {
	apiKey, err := fetchApiKey(ctx)
	if err != nil {
		return err
	}

	policy, err := a.callAccessControl(method, apiKey)
	if err != nil {
		utils.Logger.Infow("access control failed", "method", method, "err", err)
		return status.Error(codes.Unavailable, "proxy server error")
	}

	if policy == proto.AccessPolicy_ALLOW {
		return nil
	}

	return status.Error(codes.Unauthenticated, "access denied")
}

func (a *accessControlInterceptor)createApiKeyInterceptor(next *grpc.UnaryServerInterceptor) *grpc.UnaryServerInterceptor {
	f := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := a.accessControl(ctx, info.FullMethod); err != nil {
			return nil, err
		}
		if next != nil {
			return (*next)(ctx, req, info, handler)
		}
		return handler(ctx, req)
	}

	i := grpc.UnaryServerInterceptor(f)

	return &i
}

func (a *accessControlInterceptor)createStreamApiKeyInterceptor(next *grpc.StreamServerInterceptor) *grpc.StreamServerInterceptor {
	f := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := a.accessControl(ss.Context(), info.FullMethod); err != nil {
			return err
		}
		if next != nil {
			return (*next)(srv, ss, info, handler)
		}
		return handler(srv, ss)
	}

	i := grpc.StreamServerInterceptor(f)

	return &i
}
