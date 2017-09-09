// Code generated by protoc-gen-go. DO NOT EDIT.
// source: examples/benchmark/protobuf/service.proto

/*
Package esp4g_benchmark is a generated protocol buffer package.

It is generated from these files:
	examples/benchmark/protobuf/service.proto

It has these top-level messages:
	Unary
*/
package esp4g_benchmark

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Unary struct {
	Text string `protobuf:"bytes,1,opt,name=text" json:"text,omitempty"`
}

func (m *Unary) Reset()                    { *m = Unary{} }
func (m *Unary) String() string            { return proto.CompactTextString(m) }
func (*Unary) ProtoMessage()               {}
func (*Unary) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Unary) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func init() {
	proto.RegisterType((*Unary)(nil), "esp4g.benchmark.Unary")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for UnaryService service

type UnaryServiceClient interface {
	Send(ctx context.Context, in *Unary, opts ...grpc.CallOption) (*Unary, error)
}

type unaryServiceClient struct {
	cc *grpc.ClientConn
}

func NewUnaryServiceClient(cc *grpc.ClientConn) UnaryServiceClient {
	return &unaryServiceClient{cc}
}

func (c *unaryServiceClient) Send(ctx context.Context, in *Unary, opts ...grpc.CallOption) (*Unary, error) {
	out := new(Unary)
	err := grpc.Invoke(ctx, "/esp4g.benchmark.UnaryService/Send", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for UnaryService service

type UnaryServiceServer interface {
	Send(context.Context, *Unary) (*Unary, error)
}

func RegisterUnaryServiceServer(s *grpc.Server, srv UnaryServiceServer) {
	s.RegisterService(&_UnaryService_serviceDesc, srv)
}

func _UnaryService_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Unary)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UnaryServiceServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/esp4g.benchmark.UnaryService/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UnaryServiceServer).Send(ctx, req.(*Unary))
	}
	return interceptor(ctx, in, info, handler)
}

var _UnaryService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "esp4g.benchmark.UnaryService",
	HandlerType: (*UnaryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _UnaryService_Send_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "examples/benchmark/protobuf/service.proto",
}

func init() { proto.RegisterFile("examples/benchmark/protobuf/service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 135 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4c, 0xad, 0x48, 0xcc,
	0x2d, 0xc8, 0x49, 0x2d, 0xd6, 0x4f, 0x4a, 0xcd, 0x4b, 0xce, 0xc8, 0x4d, 0x2c, 0xca, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0x4f, 0x2a, 0x4d, 0xd3, 0x2f, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0xd5,
	0x03, 0x0b, 0x08, 0xf1, 0xa7, 0x16, 0x17, 0x98, 0xa4, 0xeb, 0xc1, 0xd5, 0x29, 0x49, 0x73, 0xb1,
	0x86, 0xe6, 0x25, 0x16, 0x55, 0x0a, 0x09, 0x71, 0xb1, 0x94, 0xa4, 0x56, 0x94, 0x48, 0x30, 0x2a,
	0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x46, 0x1e, 0x5c, 0x3c, 0x60, 0xc9, 0x60, 0x88, 0x19, 0x42,
	0x16, 0x5c, 0x2c, 0xc1, 0xa9, 0x79, 0x29, 0x42, 0x62, 0x7a, 0x68, 0xc6, 0xe8, 0x81, 0x95, 0x49,
	0xe1, 0x10, 0x57, 0x62, 0x48, 0x62, 0x03, 0x5b, 0x6f, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xba,
	0x0e, 0x1d, 0xb7, 0xab, 0x00, 0x00, 0x00,
}
