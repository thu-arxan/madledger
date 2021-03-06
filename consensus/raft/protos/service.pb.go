// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service.proto

package protos

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

type RaftTX struct {
	Tx                   []byte   `protobuf:"bytes,1,opt,name=Tx,proto3" json:"Tx,omitempty"`
	Caller               uint64   `protobuf:"varint,2,opt,name=Caller,proto3" json:"Caller,omitempty"`
	Channel              string   `protobuf:"bytes,3,opt,name=Channel,proto3" json:"Channel,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RaftTX) Reset()         { *m = RaftTX{} }
func (m *RaftTX) String() string { return proto.CompactTextString(m) }
func (*RaftTX) ProtoMessage()    {}
func (*RaftTX) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_85a69614227ce52d, []int{0}
}
func (m *RaftTX) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RaftTX.Unmarshal(m, b)
}
func (m *RaftTX) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RaftTX.Marshal(b, m, deterministic)
}
func (dst *RaftTX) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RaftTX.Merge(dst, src)
}
func (m *RaftTX) XXX_Size() int {
	return xxx_messageInfo_RaftTX.Size(m)
}
func (m *RaftTX) XXX_DiscardUnknown() {
	xxx_messageInfo_RaftTX.DiscardUnknown(m)
}

var xxx_messageInfo_RaftTX proto.InternalMessageInfo

func (m *RaftTX) GetTx() []byte {
	if m != nil {
		return m.Tx
	}
	return nil
}

func (m *RaftTX) GetCaller() uint64 {
	if m != nil {
		return m.Caller
	}
	return 0
}

func (m *RaftTX) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

type None struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *None) Reset()         { *m = None{} }
func (m *None) String() string { return proto.CompactTextString(m) }
func (*None) ProtoMessage()    {}
func (*None) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_85a69614227ce52d, []int{1}
}
func (m *None) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_None.Unmarshal(m, b)
}
func (m *None) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_None.Marshal(b, m, deterministic)
}
func (dst *None) XXX_Merge(src proto.Message) {
	xxx_messageInfo_None.Merge(dst, src)
}
func (m *None) XXX_Size() int {
	return xxx_messageInfo_None.Size(m)
}
func (m *None) XXX_DiscardUnknown() {
	xxx_messageInfo_None.DiscardUnknown(m)
}

var xxx_messageInfo_None proto.InternalMessageInfo

func init() {
	proto.RegisterType((*RaftTX)(nil), "protos.RaftTX")
	proto.RegisterType((*None)(nil), "protos.None")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BlockChainClient is the client API for BlockChain service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BlockChainClient interface {
	AddTx(ctx context.Context, in *RaftTX, opts ...grpc.CallOption) (*None, error)
}

type blockChainClient struct {
	cc *grpc.ClientConn
}

func NewBlockChainClient(cc *grpc.ClientConn) BlockChainClient {
	return &blockChainClient{cc}
}

func (c *blockChainClient) AddTx(ctx context.Context, in *RaftTX, opts ...grpc.CallOption) (*None, error) {
	out := new(None)
	err := c.cc.Invoke(ctx, "/protos.BlockChain/AddTx", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BlockChainServer is the server API for BlockChain service.
type BlockChainServer interface {
	AddTx(context.Context, *RaftTX) (*None, error)
}

func RegisterBlockChainServer(s *grpc.Server, srv BlockChainServer) {
	s.RegisterService(&_BlockChain_serviceDesc, srv)
}

func _BlockChain_AddTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RaftTX)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BlockChainServer).AddTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.BlockChain/AddTx",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BlockChainServer).AddTx(ctx, req.(*RaftTX))
	}
	return interceptor(ctx, in, info, handler)
}

var _BlockChain_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.BlockChain",
	HandlerType: (*BlockChainServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddTx",
			Handler:    _BlockChain_AddTx_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

func init() { proto.RegisterFile("service.proto", fileDescriptor_service_85a69614227ce52d) }

var fileDescriptor_service_85a69614227ce52d = []byte{
	// 159 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x03, 0x53, 0xc5, 0x4a, 0x5e,
	0x5c, 0x6c, 0x41, 0x89, 0x69, 0x25, 0x21, 0x11, 0x42, 0x7c, 0x5c, 0x4c, 0x21, 0x15, 0x12, 0x8c,
	0x0a, 0x8c, 0x1a, 0x3c, 0x41, 0x4c, 0x21, 0x15, 0x42, 0x62, 0x5c, 0x6c, 0xce, 0x89, 0x39, 0x39,
	0xa9, 0x45, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x2c, 0x41, 0x50, 0x9e, 0x90, 0x04, 0x17, 0xbb, 0x73,
	0x46, 0x62, 0x5e, 0x5e, 0x6a, 0x8e, 0x04, 0xb3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x8c, 0xab, 0xc4,
	0xc6, 0xc5, 0xe2, 0x97, 0x9f, 0x97, 0x6a, 0x64, 0xca, 0xc5, 0xe5, 0x94, 0x93, 0x9f, 0x9c, 0xed,
	0x9c, 0x91, 0x98, 0x99, 0x27, 0xa4, 0xce, 0xc5, 0xea, 0x98, 0x92, 0x12, 0x52, 0x21, 0xc4, 0x07,
	0xb1, 0xba, 0x58, 0x0f, 0x62, 0xa1, 0x14, 0x0f, 0x8c, 0x0f, 0xd2, 0xa4, 0xc4, 0x90, 0x04, 0x71,
	0x92, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x5e, 0xff, 0x07, 0x55, 0xaa, 0x00, 0x00, 0x00,
}
