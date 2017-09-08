// Code generated by protoc-gen-go. DO NOT EDIT.
// source: examples/calc/protobuf/service.proto

/*
Package esp4g_calc is a generated protocol buffer package.

It is generated from these files:
	examples/calc/protobuf/service.proto

It has these top-level messages:
	Operand
	OperandList
	Sum
*/
package esp4g_calc

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

type Operand struct {
	X int64 `protobuf:"varint,1,opt,name=x" json:"x,omitempty"`
	Y int64 `protobuf:"varint,2,opt,name=y" json:"y,omitempty"`
}

func (m *Operand) Reset()                    { *m = Operand{} }
func (m *Operand) String() string            { return proto.CompactTextString(m) }
func (*Operand) ProtoMessage()               {}
func (*Operand) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Operand) GetX() int64 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Operand) GetY() int64 {
	if m != nil {
		return m.Y
	}
	return 0
}

type OperandList struct {
	Operand []*Operand `protobuf:"bytes,1,rep,name=operand" json:"operand,omitempty"`
}

func (m *OperandList) Reset()                    { *m = OperandList{} }
func (m *OperandList) String() string            { return proto.CompactTextString(m) }
func (*OperandList) ProtoMessage()               {}
func (*OperandList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *OperandList) GetOperand() []*Operand {
	if m != nil {
		return m.Operand
	}
	return nil
}

type Sum struct {
	Z int64 `protobuf:"varint,3,opt,name=z" json:"z,omitempty"`
}

func (m *Sum) Reset()                    { *m = Sum{} }
func (m *Sum) String() string            { return proto.CompactTextString(m) }
func (*Sum) ProtoMessage()               {}
func (*Sum) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Sum) GetZ() int64 {
	if m != nil {
		return m.Z
	}
	return 0
}

func init() {
	proto.RegisterType((*Operand)(nil), "esp4g.calc.Operand")
	proto.RegisterType((*OperandList)(nil), "esp4g.calc.OperandList")
	proto.RegisterType((*Sum)(nil), "esp4g.calc.Sum")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for CalcService service

type CalcServiceClient interface {
	AddAll(ctx context.Context, opts ...grpc.CallOption) (CalcService_AddAllClient, error)
	AddDeffered(ctx context.Context, in *OperandList, opts ...grpc.CallOption) (CalcService_AddDefferedClient, error)
	AddAsync(ctx context.Context, opts ...grpc.CallOption) (CalcService_AddAsyncClient, error)
}

type calcServiceClient struct {
	cc *grpc.ClientConn
}

func NewCalcServiceClient(cc *grpc.ClientConn) CalcServiceClient {
	return &calcServiceClient{cc}
}

func (c *calcServiceClient) AddAll(ctx context.Context, opts ...grpc.CallOption) (CalcService_AddAllClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_CalcService_serviceDesc.Streams[0], c.cc, "/esp4g.calc.CalcService/AddAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &calcServiceAddAllClient{stream}
	return x, nil
}

type CalcService_AddAllClient interface {
	Send(*Operand) error
	CloseAndRecv() (*Sum, error)
	grpc.ClientStream
}

type calcServiceAddAllClient struct {
	grpc.ClientStream
}

func (x *calcServiceAddAllClient) Send(m *Operand) error {
	return x.ClientStream.SendMsg(m)
}

func (x *calcServiceAddAllClient) CloseAndRecv() (*Sum, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Sum)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *calcServiceClient) AddDeffered(ctx context.Context, in *OperandList, opts ...grpc.CallOption) (CalcService_AddDefferedClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_CalcService_serviceDesc.Streams[1], c.cc, "/esp4g.calc.CalcService/AddDeffered", opts...)
	if err != nil {
		return nil, err
	}
	x := &calcServiceAddDefferedClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type CalcService_AddDefferedClient interface {
	Recv() (*Sum, error)
	grpc.ClientStream
}

type calcServiceAddDefferedClient struct {
	grpc.ClientStream
}

func (x *calcServiceAddDefferedClient) Recv() (*Sum, error) {
	m := new(Sum)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *calcServiceClient) AddAsync(ctx context.Context, opts ...grpc.CallOption) (CalcService_AddAsyncClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_CalcService_serviceDesc.Streams[2], c.cc, "/esp4g.calc.CalcService/AddAsync", opts...)
	if err != nil {
		return nil, err
	}
	x := &calcServiceAddAsyncClient{stream}
	return x, nil
}

type CalcService_AddAsyncClient interface {
	Send(*Operand) error
	Recv() (*Sum, error)
	grpc.ClientStream
}

type calcServiceAddAsyncClient struct {
	grpc.ClientStream
}

func (x *calcServiceAddAsyncClient) Send(m *Operand) error {
	return x.ClientStream.SendMsg(m)
}

func (x *calcServiceAddAsyncClient) Recv() (*Sum, error) {
	m := new(Sum)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for CalcService service

type CalcServiceServer interface {
	AddAll(CalcService_AddAllServer) error
	AddDeffered(*OperandList, CalcService_AddDefferedServer) error
	AddAsync(CalcService_AddAsyncServer) error
}

func RegisterCalcServiceServer(s *grpc.Server, srv CalcServiceServer) {
	s.RegisterService(&_CalcService_serviceDesc, srv)
}

func _CalcService_AddAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CalcServiceServer).AddAll(&calcServiceAddAllServer{stream})
}

type CalcService_AddAllServer interface {
	SendAndClose(*Sum) error
	Recv() (*Operand, error)
	grpc.ServerStream
}

type calcServiceAddAllServer struct {
	grpc.ServerStream
}

func (x *calcServiceAddAllServer) SendAndClose(m *Sum) error {
	return x.ServerStream.SendMsg(m)
}

func (x *calcServiceAddAllServer) Recv() (*Operand, error) {
	m := new(Operand)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _CalcService_AddDeffered_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(OperandList)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(CalcServiceServer).AddDeffered(m, &calcServiceAddDefferedServer{stream})
}

type CalcService_AddDefferedServer interface {
	Send(*Sum) error
	grpc.ServerStream
}

type calcServiceAddDefferedServer struct {
	grpc.ServerStream
}

func (x *calcServiceAddDefferedServer) Send(m *Sum) error {
	return x.ServerStream.SendMsg(m)
}

func _CalcService_AddAsync_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CalcServiceServer).AddAsync(&calcServiceAddAsyncServer{stream})
}

type CalcService_AddAsyncServer interface {
	Send(*Sum) error
	Recv() (*Operand, error)
	grpc.ServerStream
}

type calcServiceAddAsyncServer struct {
	grpc.ServerStream
}

func (x *calcServiceAddAsyncServer) Send(m *Sum) error {
	return x.ServerStream.SendMsg(m)
}

func (x *calcServiceAddAsyncServer) Recv() (*Operand, error) {
	m := new(Operand)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _CalcService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "esp4g.calc.CalcService",
	HandlerType: (*CalcServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "AddAll",
			Handler:       _CalcService_AddAll_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "AddDeffered",
			Handler:       _CalcService_AddDeffered_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "AddAsync",
			Handler:       _CalcService_AddAsync_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "examples/calc/protobuf/service.proto",
}

func init() { proto.RegisterFile("examples/calc/protobuf/service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 233 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x49, 0xad, 0x48, 0xcc,
	0x2d, 0xc8, 0x49, 0x2d, 0xd6, 0x4f, 0x4e, 0xcc, 0x49, 0xd6, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0x4f,
	0x2a, 0x4d, 0xd3, 0x2f, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0xd5, 0x03, 0x0b, 0x08, 0x71, 0xa5,
	0x16, 0x17, 0x98, 0xa4, 0xeb, 0x81, 0x94, 0x28, 0xa9, 0x72, 0xb1, 0xfb, 0x17, 0xa4, 0x16, 0x25,
	0xe6, 0xa5, 0x08, 0xf1, 0x70, 0x31, 0x56, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x30, 0x07, 0x31, 0x56,
	0x80, 0x78, 0x95, 0x12, 0x4c, 0x10, 0x5e, 0xa5, 0x92, 0x0d, 0x17, 0x37, 0x54, 0x99, 0x4f, 0x66,
	0x71, 0x89, 0x90, 0x2e, 0x17, 0x7b, 0x3e, 0x84, 0x2b, 0xc1, 0xa8, 0xc0, 0xac, 0xc1, 0x6d, 0x24,
	0xac, 0x87, 0x30, 0x53, 0x0f, 0xaa, 0x32, 0x08, 0xa6, 0x46, 0x49, 0x98, 0x8b, 0x39, 0xb8, 0x34,
	0x17, 0x64, 0x64, 0x95, 0x04, 0x33, 0xc4, 0xc8, 0x2a, 0xa3, 0x6d, 0x8c, 0x5c, 0xdc, 0xce, 0x89,
	0x39, 0xc9, 0xc1, 0x10, 0xb7, 0x09, 0x19, 0x71, 0xb1, 0x39, 0xa6, 0xa4, 0x38, 0xe6, 0xe4, 0x08,
	0x61, 0x33, 0x4c, 0x8a, 0x1f, 0x59, 0x30, 0xb8, 0x34, 0x57, 0x89, 0x41, 0x83, 0x51, 0xc8, 0x9a,
	0x8b, 0xdb, 0x31, 0x25, 0xc5, 0x25, 0x35, 0x2d, 0x2d, 0xb5, 0x28, 0x35, 0x45, 0x48, 0x1c, 0x8b,
	0x46, 0x90, 0x7b, 0xb1, 0x68, 0x36, 0x60, 0x14, 0x32, 0xe3, 0xe2, 0x00, 0x59, 0x58, 0x5c, 0x99,
	0x97, 0x4c, 0xbc, 0x95, 0x06, 0x8c, 0x49, 0x6c, 0xe0, 0x50, 0x34, 0x06, 0x04, 0x00, 0x00, 0xff,
	0xff, 0x99, 0x49, 0x9d, 0x1d, 0x6d, 0x01, 0x00, 0x00,
}
