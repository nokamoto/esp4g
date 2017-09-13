package esp4g

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getMetadata(ctx context.Context) (metadata.MD, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "failed to read metadata")
	}
	return md, nil
}

func safeMetadata(md metadata.MD, k string) []string {
	res, ok := md[k]
	if !ok {
		return []string{}
	}
	return res
}
