// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protobuf/service.proto

/*
Package eps4g_extension is a generated protocol buffer package.

It is generated from these files:
	protobuf/service.proto

It has these top-level messages:
	UnaryAccessLog
	StreamAccessLog
	AccessIdentity
	AccessControl
*/
package eps4g_extension

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/empty"
import google_protobuf1 "github.com/golang/protobuf/ptypes/duration"

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

type AccessPolicy int32

const (
	AccessPolicy_ALLOW AccessPolicy = 0
	AccessPolicy_DENY  AccessPolicy = 1
)

var AccessPolicy_name = map[int32]string{
	0: "ALLOW",
	1: "DENY",
}
var AccessPolicy_value = map[string]int32{
	"ALLOW": 0,
	"DENY":  1,
}

func (x AccessPolicy) String() string {
	return proto.EnumName(AccessPolicy_name, int32(x))
}
func (AccessPolicy) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type UnaryAccessLog struct {
	Method       string                     `protobuf:"bytes,1,opt,name=method" json:"method,omitempty"`
	ResponseTime *google_protobuf1.Duration `protobuf:"bytes,2,opt,name=response_time,json=responseTime" json:"response_time,omitempty"`
	Status       string                     `protobuf:"bytes,3,opt,name=status" json:"status,omitempty"`
	RequestSize  int64                      `protobuf:"varint,4,opt,name=request_size,json=requestSize" json:"request_size,omitempty"`
	ResponseSize int64                      `protobuf:"varint,5,opt,name=response_size,json=responseSize" json:"response_size,omitempty"`
}

func (m *UnaryAccessLog) Reset()                    { *m = UnaryAccessLog{} }
func (m *UnaryAccessLog) String() string            { return proto.CompactTextString(m) }
func (*UnaryAccessLog) ProtoMessage()               {}
func (*UnaryAccessLog) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *UnaryAccessLog) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *UnaryAccessLog) GetResponseTime() *google_protobuf1.Duration {
	if m != nil {
		return m.ResponseTime
	}
	return nil
}

func (m *UnaryAccessLog) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *UnaryAccessLog) GetRequestSize() int64 {
	if m != nil {
		return m.RequestSize
	}
	return 0
}

func (m *UnaryAccessLog) GetResponseSize() int64 {
	if m != nil {
		return m.ResponseSize
	}
	return 0
}

type StreamAccessLog struct {
	Method       string                     `protobuf:"bytes,1,opt,name=method" json:"method,omitempty"`
	ResponseTime *google_protobuf1.Duration `protobuf:"bytes,2,opt,name=response_time,json=responseTime" json:"response_time,omitempty"`
	Status       string                     `protobuf:"bytes,3,opt,name=status" json:"status,omitempty"`
}

func (m *StreamAccessLog) Reset()                    { *m = StreamAccessLog{} }
func (m *StreamAccessLog) String() string            { return proto.CompactTextString(m) }
func (*StreamAccessLog) ProtoMessage()               {}
func (*StreamAccessLog) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *StreamAccessLog) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *StreamAccessLog) GetResponseTime() *google_protobuf1.Duration {
	if m != nil {
		return m.ResponseTime
	}
	return nil
}

func (m *StreamAccessLog) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type AccessIdentity struct {
	Method string   `protobuf:"bytes,1,opt,name=method" json:"method,omitempty"`
	ApiKey []string `protobuf:"bytes,2,rep,name=api_key,json=apiKey" json:"api_key,omitempty"`
}

func (m *AccessIdentity) Reset()                    { *m = AccessIdentity{} }
func (m *AccessIdentity) String() string            { return proto.CompactTextString(m) }
func (*AccessIdentity) ProtoMessage()               {}
func (*AccessIdentity) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *AccessIdentity) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *AccessIdentity) GetApiKey() []string {
	if m != nil {
		return m.ApiKey
	}
	return nil
}

type AccessControl struct {
	Policy AccessPolicy `protobuf:"varint,1,opt,name=policy,enum=eps4g.extension.AccessPolicy" json:"policy,omitempty"`
}

func (m *AccessControl) Reset()                    { *m = AccessControl{} }
func (m *AccessControl) String() string            { return proto.CompactTextString(m) }
func (*AccessControl) ProtoMessage()               {}
func (*AccessControl) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *AccessControl) GetPolicy() AccessPolicy {
	if m != nil {
		return m.Policy
	}
	return AccessPolicy_ALLOW
}

func init() {
	proto.RegisterType((*UnaryAccessLog)(nil), "eps4g.extension.UnaryAccessLog")
	proto.RegisterType((*StreamAccessLog)(nil), "eps4g.extension.StreamAccessLog")
	proto.RegisterType((*AccessIdentity)(nil), "eps4g.extension.AccessIdentity")
	proto.RegisterType((*AccessControl)(nil), "eps4g.extension.AccessControl")
	proto.RegisterEnum("eps4g.extension.AccessPolicy", AccessPolicy_name, AccessPolicy_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for AccessLogService service

type AccessLogServiceClient interface {
	UnaryAccess(ctx context.Context, in *UnaryAccessLog, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	StreamAccess(ctx context.Context, in *StreamAccessLog, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
}

type accessLogServiceClient struct {
	cc *grpc.ClientConn
}

func NewAccessLogServiceClient(cc *grpc.ClientConn) AccessLogServiceClient {
	return &accessLogServiceClient{cc}
}

func (c *accessLogServiceClient) UnaryAccess(ctx context.Context, in *UnaryAccessLog, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/eps4g.extension.AccessLogService/UnaryAccess", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accessLogServiceClient) StreamAccess(ctx context.Context, in *StreamAccessLog, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/eps4g.extension.AccessLogService/StreamAccess", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AccessLogService service

type AccessLogServiceServer interface {
	UnaryAccess(context.Context, *UnaryAccessLog) (*google_protobuf.Empty, error)
	StreamAccess(context.Context, *StreamAccessLog) (*google_protobuf.Empty, error)
}

func RegisterAccessLogServiceServer(s *grpc.Server, srv AccessLogServiceServer) {
	s.RegisterService(&_AccessLogService_serviceDesc, srv)
}

func _AccessLogService_UnaryAccess_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnaryAccessLog)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessLogServiceServer).UnaryAccess(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eps4g.extension.AccessLogService/UnaryAccess",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessLogServiceServer).UnaryAccess(ctx, req.(*UnaryAccessLog))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccessLogService_StreamAccess_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StreamAccessLog)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessLogServiceServer).StreamAccess(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eps4g.extension.AccessLogService/StreamAccess",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessLogServiceServer).StreamAccess(ctx, req.(*StreamAccessLog))
	}
	return interceptor(ctx, in, info, handler)
}

var _AccessLogService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "eps4g.extension.AccessLogService",
	HandlerType: (*AccessLogServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UnaryAccess",
			Handler:    _AccessLogService_UnaryAccess_Handler,
		},
		{
			MethodName: "StreamAccess",
			Handler:    _AccessLogService_StreamAccess_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protobuf/service.proto",
}

// Client API for AccessControlService service

type AccessControlServiceClient interface {
	Access(ctx context.Context, in *AccessIdentity, opts ...grpc.CallOption) (*AccessControl, error)
}

type accessControlServiceClient struct {
	cc *grpc.ClientConn
}

func NewAccessControlServiceClient(cc *grpc.ClientConn) AccessControlServiceClient {
	return &accessControlServiceClient{cc}
}

func (c *accessControlServiceClient) Access(ctx context.Context, in *AccessIdentity, opts ...grpc.CallOption) (*AccessControl, error) {
	out := new(AccessControl)
	err := grpc.Invoke(ctx, "/eps4g.extension.AccessControlService/Access", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AccessControlService service

type AccessControlServiceServer interface {
	Access(context.Context, *AccessIdentity) (*AccessControl, error)
}

func RegisterAccessControlServiceServer(s *grpc.Server, srv AccessControlServiceServer) {
	s.RegisterService(&_AccessControlService_serviceDesc, srv)
}

func _AccessControlService_Access_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccessIdentity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessControlServiceServer).Access(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/eps4g.extension.AccessControlService/Access",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessControlServiceServer).Access(ctx, req.(*AccessIdentity))
	}
	return interceptor(ctx, in, info, handler)
}

var _AccessControlService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "eps4g.extension.AccessControlService",
	HandlerType: (*AccessControlServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Access",
			Handler:    _AccessControlService_Access_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protobuf/service.proto",
}

func init() { proto.RegisterFile("protobuf/service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 414 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x92, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x86, 0xb3, 0x4d, 0x6b, 0xe8, 0xc4, 0x4d, 0xa3, 0x15, 0x0a, 0xc6, 0x88, 0x62, 0xdc, 0x4b,
	0xc4, 0xc1, 0x95, 0x0c, 0x5c, 0x91, 0x22, 0x5a, 0x04, 0x24, 0x02, 0xe4, 0x80, 0x10, 0xa7, 0xc8,
	0x71, 0x06, 0xb3, 0x22, 0xf6, 0x9a, 0xdd, 0x35, 0xc2, 0xb9, 0xf1, 0x36, 0x3c, 0x0a, 0x8f, 0x85,
	0xec, 0x75, 0x22, 0x27, 0x91, 0xcf, 0x3d, 0xee, 0xcc, 0xa7, 0xff, 0x9f, 0x9d, 0xf9, 0x61, 0x98,
	0x09, 0xae, 0xf8, 0x22, 0xff, 0x76, 0x25, 0x51, 0xfc, 0x62, 0x11, 0x7a, 0x55, 0x81, 0x9e, 0x63,
	0x26, 0x9f, 0xc7, 0x1e, 0xfe, 0x56, 0x98, 0x4a, 0xc6, 0x53, 0xfb, 0x61, 0xcc, 0x79, 0xbc, 0xc2,
	0xab, 0x2d, 0x8f, 0x49, 0xa6, 0x0a, 0x4d, 0xdb, 0x17, 0xfb, 0xcd, 0x65, 0x2e, 0x42, 0xc5, 0x78,
	0xaa, 0xfb, 0xee, 0x3f, 0x02, 0xfd, 0xcf, 0x69, 0x28, 0x8a, 0x71, 0x14, 0xa1, 0x94, 0x53, 0x1e,
	0xd3, 0x21, 0x18, 0x09, 0xaa, 0xef, 0x7c, 0x69, 0x11, 0x87, 0x8c, 0x4e, 0x83, 0xfa, 0x45, 0x5f,
	0xc2, 0x99, 0x40, 0x99, 0xf1, 0x54, 0xe2, 0x5c, 0xb1, 0x04, 0xad, 0x23, 0x87, 0x8c, 0x7a, 0xfe,
	0x03, 0x4f, 0x5b, 0x78, 0x1b, 0x0b, 0xef, 0xba, 0xb6, 0x08, 0xcc, 0x0d, 0xff, 0x89, 0x25, 0x58,
	0xea, 0x4a, 0x15, 0xaa, 0x5c, 0x5a, 0x5d, 0xad, 0xab, 0x5f, 0xf4, 0x09, 0x98, 0x02, 0x7f, 0xe6,
	0x28, 0xd5, 0x5c, 0xb2, 0x35, 0x5a, 0xc7, 0x0e, 0x19, 0x75, 0x83, 0x5e, 0x5d, 0x9b, 0xb1, 0x35,
	0xd2, 0xcb, 0x86, 0x75, 0xc5, 0x9c, 0x54, 0xcc, 0x56, 0xbf, 0x84, 0xdc, 0x3f, 0x04, 0xce, 0x67,
	0x4a, 0x60, 0x98, 0xdc, 0xda, 0x5f, 0xdc, 0x31, 0xf4, 0xb5, 0xf9, 0xdb, 0x25, 0xa6, 0x8a, 0xa9,
	0xa2, 0x75, 0x82, 0xfb, 0x70, 0x27, 0xcc, 0xd8, 0xfc, 0x07, 0x16, 0xd6, 0x91, 0xd3, 0x2d, 0x1b,
	0x61, 0xc6, 0x26, 0x58, 0xb8, 0xaf, 0xe1, 0x4c, 0x4b, 0xbc, 0xe2, 0xa9, 0x12, 0x7c, 0x45, 0x5f,
	0x80, 0x91, 0xf1, 0x15, 0x8b, 0x8a, 0x4a, 0xa1, 0xef, 0x3f, 0xf2, 0xf6, 0x12, 0xe0, 0x69, 0xfe,
	0x63, 0x05, 0x05, 0x35, 0xfc, 0xf4, 0x12, 0xcc, 0x66, 0x9d, 0x9e, 0xc2, 0xc9, 0x78, 0x3a, 0xfd,
	0xf0, 0x65, 0xd0, 0xa1, 0x77, 0xe1, 0xf8, 0xfa, 0xe6, 0xfd, 0xd7, 0x01, 0xf1, 0xff, 0x12, 0x18,
	0x6c, 0xb7, 0x35, 0xd3, 0x39, 0xa3, 0x6f, 0xa0, 0xd7, 0x88, 0x04, 0x7d, 0x7c, 0xe0, 0xb7, 0x1b,
	0x18, 0x7b, 0x78, 0xb0, 0xb5, 0x9b, 0x32, 0x81, 0x6e, 0x87, 0xbe, 0x03, 0xb3, 0x79, 0x11, 0xea,
	0x1c, 0x48, 0xed, 0x1d, 0xac, 0x5d, 0xcb, 0x8f, 0xe0, 0xde, 0xce, 0x5e, 0x36, 0xd3, 0x4e, 0xc0,
	0x68, 0x1d, 0x74, 0xf7, 0x16, 0xf6, 0x45, 0x0b, 0x50, 0x2b, 0xba, 0x9d, 0x85, 0x51, 0xd9, 0x3e,
	0xfb, 0x1f, 0x00, 0x00, 0xff, 0xff, 0x4f, 0xe7, 0xcf, 0x41, 0x7d, 0x03, 0x00, 0x00,
}