// Code generated by protoc-gen-go. DO NOT EDIT.
// source: p2p.proto

/*
Package p2p is a generated protocol buffer package.

It is generated from these files:
	p2p.proto

It has these top-level messages:
	P2PMessage
	P2PReply
*/
package p2p

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import coredocument "github.com/CentrifugeInc/centrifuge-protobufs/coredocument"

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

// P2PMessage wraps a single CoreDocument to be transferred to another noed
type P2PMessage struct {
	Document *coredocument.CoreDocument `protobuf:"bytes,1,opt,name=document" json:"document,omitempty"`
}

func (m *P2PMessage) Reset()                    { *m = P2PMessage{} }
func (m *P2PMessage) String() string            { return proto.CompactTextString(m) }
func (*P2PMessage) ProtoMessage()               {}
func (*P2PMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *P2PMessage) GetDocument() *coredocument.CoreDocument {
	if m != nil {
		return m.Document
	}
	return nil
}

type P2PReply struct {
	Document *coredocument.CoreDocument `protobuf:"bytes,1,opt,name=document" json:"document,omitempty"`
}

func (m *P2PReply) Reset()                    { *m = P2PReply{} }
func (m *P2PReply) String() string            { return proto.CompactTextString(m) }
func (*P2PReply) ProtoMessage()               {}
func (*P2PReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *P2PReply) GetDocument() *coredocument.CoreDocument {
	if m != nil {
		return m.Document
	}
	return nil
}

func init() {
	proto.RegisterType((*P2PMessage)(nil), "p2p.P2PMessage")
	proto.RegisterType((*P2PReply)(nil), "p2p.P2PReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for P2PService service

type P2PServiceClient interface {
	Transmit(ctx context.Context, in *P2PMessage, opts ...grpc.CallOption) (*P2PReply, error)
}

type p2PServiceClient struct {
	cc *grpc.ClientConn
}

func NewP2PServiceClient(cc *grpc.ClientConn) P2PServiceClient {
	return &p2PServiceClient{cc}
}

func (c *p2PServiceClient) Transmit(ctx context.Context, in *P2PMessage, opts ...grpc.CallOption) (*P2PReply, error) {
	out := new(P2PReply)
	err := grpc.Invoke(ctx, "/p2p.P2PService/Transmit", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for P2PService service

type P2PServiceServer interface {
	Transmit(context.Context, *P2PMessage) (*P2PReply, error)
}

func RegisterP2PServiceServer(s *grpc.Server, srv P2PServiceServer) {
	s.RegisterService(&_P2PService_serviceDesc, srv)
}

func _P2PService_Transmit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(P2PMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServiceServer).Transmit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/p2p.P2PService/Transmit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServiceServer).Transmit(ctx, req.(*P2PMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _P2PService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "p2p.P2PService",
	HandlerType: (*P2PServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Transmit",
			Handler:    _P2PService_Transmit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "p2p.proto",
}

func init() { proto.RegisterFile("p2p.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 209 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2c, 0x30, 0x2a, 0xd0,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2e, 0x30, 0x2a, 0x90, 0x32, 0x49, 0x4e, 0xcd, 0x2b,
	0x29, 0xca, 0x4c, 0x2b, 0x4d, 0x4f, 0xd5, 0x05, 0x0b, 0x27, 0x95, 0xa6, 0x15, 0xeb, 0x27, 0xe7,
	0x17, 0xa5, 0xa6, 0xe4, 0x27, 0x97, 0xe6, 0xa6, 0xe6, 0x95, 0xa0, 0x70, 0x20, 0x5a, 0x95, 0x5c,
	0xb8, 0xb8, 0x02, 0x8c, 0x02, 0x7c, 0x53, 0x8b, 0x8b, 0x13, 0xd3, 0x53, 0x85, 0xcc, 0xb8, 0x38,
	0x60, 0xf2, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xdc, 0x46, 0x52, 0x7a, 0x28, 0x9a, 0x9c, 0xf3, 0x8b,
	0x52, 0x5d, 0xa0, 0x9c, 0x20, 0xb8, 0x5a, 0x25, 0x27, 0x2e, 0x8e, 0x00, 0xa3, 0x80, 0xa0, 0xd4,
	0x82, 0x9c, 0x4a, 0x72, 0xcd, 0x30, 0xb2, 0x02, 0xbb, 0x24, 0x38, 0xb5, 0xa8, 0x2c, 0x33, 0x39,
	0x55, 0x48, 0x87, 0x8b, 0x23, 0xa4, 0x28, 0x31, 0xaf, 0x38, 0x37, 0xb3, 0x44, 0x88, 0x5f, 0x0f,
	0xe4, 0x55, 0x84, 0x33, 0xa5, 0x78, 0x61, 0x02, 0x60, 0x1b, 0x95, 0x18, 0x9c, 0xcc, 0xa3, 0x4c,
	0xd3, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5, 0x9d, 0xe1, 0x01, 0xe1, 0x99,
	0x97, 0xac, 0x9f, 0x9e, 0xaf, 0x8b, 0x08, 0x19, 0x7d, 0x24, 0x66, 0x81, 0x51, 0x41, 0x12, 0x1b,
	0x38, 0x14, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xe6, 0x62, 0xcc, 0x2e, 0x4d, 0x01, 0x00,
	0x00,
}