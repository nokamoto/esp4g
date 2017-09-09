package esp4g

import (
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"fmt"
	"log"
	extension "github.com/nokamoto/esp4g/protobuf"
)

const API_KEY_HEADER = "x-api-key"

type accessControlInterceptor struct {
	con *grpc.ClientConn
}

func newAccessControlInterceptor(port int) *accessControlInterceptor {
	opts := []grpc.DialOption{grpc.WithInsecure()}

	con, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), opts...)
	if err != nil {
		log.Fatal(err)
	}

	return &accessControlInterceptor{con: con}
}

func (a *accessControlInterceptor)doAccessControl(method string, keys []string) (extension.AccessPolicy, error) {
	client := extension.NewAccessControlServiceClient(a.con)
	id := extension.AccessIdentity{
		Method: method,
		ApiKey: keys,
	}
	ctl, err := client.Access(context.Background(), &id)
	if err != nil {
		return extension.AccessPolicy_DENY, err
	}
	return ctl.Policy, nil
}

func (a *accessControlInterceptor)createApiKeyInterceptor(next *grpc.UnaryServerInterceptor) *grpc.UnaryServerInterceptor {
	f := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "failed to read metadata")
		}

		apiKey, ok := md[API_KEY_HEADER]

		policy, err := a.doAccessControl(info.FullMethod, apiKey)
		if err != nil {
			log.Println(err)
			return nil, status.Error(codes.Unavailable, "proxy server error")
		}
		if policy == extension.AccessPolicy_ALLOW {
			if next != nil {
				return (*next)(ctx, req, info, handler)
			}
			return handler(ctx, req)
		}

		return nil, status.Error(codes.Unauthenticated, "unauthenticated request")
	}

	i := grpc.UnaryServerInterceptor(f)

	return &i
}

func (a *accessControlInterceptor)createStreamApiKeyInterceptor(next *grpc.StreamServerInterceptor) *grpc.StreamServerInterceptor {
	f := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		md, ok := metadata.FromIncomingContext(ss.Context())
		if !ok {
			return status.Error(codes.Internal, "failed to read metadata")
		}

		apiKey, ok := md[API_KEY_HEADER]

		policy, err := a.doAccessControl(info.FullMethod, apiKey)
		if err != nil {
			log.Println(err)
			return status.Error(codes.Unavailable, "proxy server error")
		}
		if policy == extension.AccessPolicy_ALLOW {
			if next != nil {
				return (*next)(srv, ss, info, handler)
			}
			return handler(srv, ss)
		}

		return status.Error(codes.Unauthenticated, "unauthenticated request")
	}

	i := grpc.StreamServerInterceptor(f)

	return &i
}
