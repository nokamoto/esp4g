// Code generated by protoc-gen-go. DO NOT EDIT.
// source: examples/ping/protobuf/service.proto

/*
Package eps4g_ping is a generated protocol buffer package.

It is generated from these files:
	examples/ping/protobuf/service.proto

It has these top-level messages:
	Ping
	Pong
*/
package eps4g_ping

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

type Ping struct {
	X int64 `protobuf:"varint,1,opt,name=x" json:"x,omitempty"`
}

func (m *Ping) Reset()                    { *m = Ping{} }
func (m *Ping) String() string            { return proto.CompactTextString(m) }
func (*Ping) ProtoMessage()               {}
func (*Ping) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Ping) GetX() int64 {
	if m != nil {
		return m.X
	}
	return 0
}

type Pong struct {
	Y int64 `protobuf:"varint,1,opt,name=y" json:"y,omitempty"`
}

func (m *Pong) Reset()                    { *m = Pong{} }
func (m *Pong) String() string            { return proto.CompactTextString(m) }
func (*Pong) ProtoMessage()               {}
func (*Pong) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Pong) GetY() int64 {
	if m != nil {
		return m.Y
	}
	return 0
}

func init() {
	proto.RegisterType((*Ping)(nil), "eps4g.ping.Ping")
	proto.RegisterType((*Pong)(nil), "eps4g.ping.Pong")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for PingService service

type PingServiceClient interface {
	Send(ctx context.Context, in *Ping, opts ...grpc.CallOption) (*Pong, error)
	Unavailable(ctx context.Context, in *Ping, opts ...grpc.CallOption) (*Pong, error)
}

type pingServiceClient struct {
	cc *grpc.ClientConn
}

func NewPingServiceClient(cc *grpc.ClientConn) PingServiceClient {
	return &pingServiceClient{cc}
}

func (c *pingServiceClient) Send(ctx context.Context, in *Ping, opts ...grpc.CallOption) (*Pong, error) {
	out := new(Pong)
	err := grpc.Invoke(ctx, "/eps4g.ping.PingService/Send", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pingServiceClient) Unavailable(ctx context.Context, in *Ping, opts ...grpc.CallOption) (*Pong, error) {
	out := new(Pong)
	err := grpc.Invoke(ctx, "/eps4g.ping.PingService/Unavailable", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for PingService service

type PingServiceServer interface {
	Send(context.Context, *Ping) (*Pong, error)
	Unavailable(context.Context, *Ping) (*Pong, error)
}

func RegisterPingServiceServer(s *grpc.Server, srv PingServiceServer) {
	s.RegisterService(&_PingService_serviceDesc, srv)
}

func _PingService_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ping)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PingServiceServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eps4g.ping.PingService/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PingServiceServer).Send(ctx, req.(*Ping))
	}
	return interceptor(ctx, in, info, handler)
}

func _PingService_Unavailable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ping)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PingServiceServer).Unavailable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eps4g.ping.PingService/Unavailable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PingServiceServer).Unavailable(ctx, req.(*Ping))
	}
	return interceptor(ctx, in, info, handler)
}

var _PingService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "eps4g.ping.PingService",
	HandlerType: (*PingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _PingService_Send_Handler,
		},
		{
			MethodName: "Unavailable",
			Handler:    _PingService_Unavailable_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "examples/ping/protobuf/service.proto",
}

func init() { proto.RegisterFile("examples/ping/protobuf/service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 155 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x49, 0xad, 0x48, 0xcc,
	0x2d, 0xc8, 0x49, 0x2d, 0xd6, 0x2f, 0xc8, 0xcc, 0x4b, 0xd7, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0x4f,
	0x2a, 0x4d, 0xd3, 0x2f, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0xd5, 0x03, 0x0b, 0x08, 0x71, 0xa5,
	0x16, 0x14, 0x9b, 0xa4, 0xeb, 0x81, 0x94, 0x28, 0x89, 0x70, 0xb1, 0x04, 0x64, 0xe6, 0xa5, 0x0b,
	0xf1, 0x70, 0x31, 0x56, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x30, 0x07, 0x31, 0x56, 0x80, 0x45, 0xf3,
	0x21, 0xa2, 0x95, 0x30, 0xd1, 0x4a, 0xa3, 0x02, 0x2e, 0x6e, 0x90, 0xda, 0x60, 0x88, 0x61, 0x42,
	0x3a, 0x5c, 0x2c, 0xc1, 0xa9, 0x79, 0x29, 0x42, 0x02, 0x7a, 0x08, 0xf3, 0xf4, 0x40, 0x0a, 0xa4,
	0x50, 0x45, 0xf2, 0xf3, 0xd2, 0x95, 0x18, 0x84, 0x8c, 0xb9, 0xb8, 0x43, 0xf3, 0x12, 0xcb, 0x12,
	0x33, 0x73, 0x12, 0x93, 0x72, 0x52, 0x89, 0xd3, 0x94, 0xc4, 0x06, 0x76, 0xb0, 0x31, 0x20, 0x00,
	0x00, 0xff, 0xff, 0x7a, 0x0e, 0x57, 0x24, 0xd8, 0x00, 0x00, 0x00,
}
