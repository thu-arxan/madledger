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

// However, this is not contains sig now, but this is necessary
// if we want to verify the permission.
// TODO: add sig to identity.
type FetchBlockRequest struct {
	ChannelID            string   `protobuf:"bytes,1,opt,name=ChannelID,proto3" json:"ChannelID,omitempty"`
	Number               uint64   `protobuf:"varint,2,opt,name=Number,proto3" json:"Number,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FetchBlockRequest) Reset()         { *m = FetchBlockRequest{} }
func (m *FetchBlockRequest) String() string { return proto.CompactTextString(m) }
func (*FetchBlockRequest) ProtoMessage()    {}
func (*FetchBlockRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_a015a6b48025b526, []int{0}
}
func (m *FetchBlockRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FetchBlockRequest.Unmarshal(m, b)
}
func (m *FetchBlockRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FetchBlockRequest.Marshal(b, m, deterministic)
}
func (dst *FetchBlockRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FetchBlockRequest.Merge(dst, src)
}
func (m *FetchBlockRequest) XXX_Size() int {
	return xxx_messageInfo_FetchBlockRequest.Size(m)
}
func (m *FetchBlockRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FetchBlockRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FetchBlockRequest proto.InternalMessageInfo

func (m *FetchBlockRequest) GetChannelID() string {
	if m != nil {
		return m.ChannelID
	}
	return ""
}

func (m *FetchBlockRequest) GetNumber() uint64 {
	if m != nil {
		return m.Number
	}
	return 0
}

type ListChannelsRequest struct {
	// If system channel are included
	System               bool     `protobuf:"varint,1,opt,name=System,proto3" json:"System,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListChannelsRequest) Reset()         { *m = ListChannelsRequest{} }
func (m *ListChannelsRequest) String() string { return proto.CompactTextString(m) }
func (*ListChannelsRequest) ProtoMessage()    {}
func (*ListChannelsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_a015a6b48025b526, []int{1}
}
func (m *ListChannelsRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListChannelsRequest.Unmarshal(m, b)
}
func (m *ListChannelsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListChannelsRequest.Marshal(b, m, deterministic)
}
func (dst *ListChannelsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListChannelsRequest.Merge(dst, src)
}
func (m *ListChannelsRequest) XXX_Size() int {
	return xxx_messageInfo_ListChannelsRequest.Size(m)
}
func (m *ListChannelsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListChannelsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListChannelsRequest proto.InternalMessageInfo

func (m *ListChannelsRequest) GetSystem() bool {
	if m != nil {
		return m.System
	}
	return false
}

// ChannelInfos contains ChannelInfo
type ChannelInfos struct {
	Channels             []*ChannelInfo `protobuf:"bytes,1,rep,name=Channels,proto3" json:"Channels,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *ChannelInfos) Reset()         { *m = ChannelInfos{} }
func (m *ChannelInfos) String() string { return proto.CompactTextString(m) }
func (*ChannelInfos) ProtoMessage()    {}
func (*ChannelInfos) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_a015a6b48025b526, []int{2}
}
func (m *ChannelInfos) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChannelInfos.Unmarshal(m, b)
}
func (m *ChannelInfos) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChannelInfos.Marshal(b, m, deterministic)
}
func (dst *ChannelInfos) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChannelInfos.Merge(dst, src)
}
func (m *ChannelInfos) XXX_Size() int {
	return xxx_messageInfo_ChannelInfos.Size(m)
}
func (m *ChannelInfos) XXX_DiscardUnknown() {
	xxx_messageInfo_ChannelInfos.DiscardUnknown(m)
}

var xxx_messageInfo_ChannelInfos proto.InternalMessageInfo

func (m *ChannelInfos) GetChannels() []*ChannelInfo {
	if m != nil {
		return m.Channels
	}
	return nil
}

// ChannelInfo includes some infomations of a channel
type ChannelInfo struct {
	ChannelID            string   `protobuf:"bytes,1,opt,name=ChannelID,proto3" json:"ChannelID,omitempty"`
	BlockSize            uint64   `protobuf:"varint,2,opt,name=BlockSize,proto3" json:"BlockSize,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChannelInfo) Reset()         { *m = ChannelInfo{} }
func (m *ChannelInfo) String() string { return proto.CompactTextString(m) }
func (*ChannelInfo) ProtoMessage()    {}
func (*ChannelInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_a015a6b48025b526, []int{3}
}
func (m *ChannelInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChannelInfo.Unmarshal(m, b)
}
func (m *ChannelInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChannelInfo.Marshal(b, m, deterministic)
}
func (dst *ChannelInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChannelInfo.Merge(dst, src)
}
func (m *ChannelInfo) XXX_Size() int {
	return xxx_messageInfo_ChannelInfo.Size(m)
}
func (m *ChannelInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_ChannelInfo.DiscardUnknown(m)
}

var xxx_messageInfo_ChannelInfo proto.InternalMessageInfo

func (m *ChannelInfo) GetChannelID() string {
	if m != nil {
		return m.ChannelID
	}
	return ""
}

func (m *ChannelInfo) GetBlockSize() uint64 {
	if m != nil {
		return m.BlockSize
	}
	return 0
}

// AddChannelRequest includes the profile of Channel
type AddChannelRequest struct {
	ChannelID            string   `protobuf:"bytes,1,opt,name=ChannelID,proto3" json:"ChannelID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddChannelRequest) Reset()         { *m = AddChannelRequest{} }
func (m *AddChannelRequest) String() string { return proto.CompactTextString(m) }
func (*AddChannelRequest) ProtoMessage()    {}
func (*AddChannelRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_a015a6b48025b526, []int{4}
}
func (m *AddChannelRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddChannelRequest.Unmarshal(m, b)
}
func (m *AddChannelRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddChannelRequest.Marshal(b, m, deterministic)
}
func (dst *AddChannelRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddChannelRequest.Merge(dst, src)
}
func (m *AddChannelRequest) XXX_Size() int {
	return xxx_messageInfo_AddChannelRequest.Size(m)
}
func (m *AddChannelRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddChannelRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddChannelRequest proto.InternalMessageInfo

func (m *AddChannelRequest) GetChannelID() string {
	if m != nil {
		return m.ChannelID
	}
	return ""
}

type AddTxRequest struct {
	Tx                   *Tx      `protobuf:"bytes,1,opt,name=Tx,proto3" json:"Tx,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddTxRequest) Reset()         { *m = AddTxRequest{} }
func (m *AddTxRequest) String() string { return proto.CompactTextString(m) }
func (*AddTxRequest) ProtoMessage()    {}
func (*AddTxRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_a015a6b48025b526, []int{5}
}
func (m *AddTxRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddTxRequest.Unmarshal(m, b)
}
func (m *AddTxRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddTxRequest.Marshal(b, m, deterministic)
}
func (dst *AddTxRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddTxRequest.Merge(dst, src)
}
func (m *AddTxRequest) XXX_Size() int {
	return xxx_messageInfo_AddTxRequest.Size(m)
}
func (m *AddTxRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddTxRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddTxRequest proto.InternalMessageInfo

func (m *AddTxRequest) GetTx() *Tx {
	if m != nil {
		return m.Tx
	}
	return nil
}

// TxStatus include nothing now.
type TxStatus struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxStatus) Reset()         { *m = TxStatus{} }
func (m *TxStatus) String() string { return proto.CompactTextString(m) }
func (*TxStatus) ProtoMessage()    {}
func (*TxStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_a015a6b48025b526, []int{6}
}
func (m *TxStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TxStatus.Unmarshal(m, b)
}
func (m *TxStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TxStatus.Marshal(b, m, deterministic)
}
func (dst *TxStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxStatus.Merge(dst, src)
}
func (m *TxStatus) XXX_Size() int {
	return xxx_messageInfo_TxStatus.Size(m)
}
func (m *TxStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_TxStatus.DiscardUnknown(m)
}

var xxx_messageInfo_TxStatus proto.InternalMessageInfo

func init() {
	proto.RegisterType((*FetchBlockRequest)(nil), "protos.FetchBlockRequest")
	proto.RegisterType((*ListChannelsRequest)(nil), "protos.ListChannelsRequest")
	proto.RegisterType((*ChannelInfos)(nil), "protos.ChannelInfos")
	proto.RegisterType((*ChannelInfo)(nil), "protos.ChannelInfo")
	proto.RegisterType((*AddChannelRequest)(nil), "protos.AddChannelRequest")
	proto.RegisterType((*AddTxRequest)(nil), "protos.AddTxRequest")
	proto.RegisterType((*TxStatus)(nil), "protos.TxStatus")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// OrdererClient is the client API for Orderer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OrdererClient interface {
	FetchBlock(ctx context.Context, in *FetchBlockRequest, opts ...grpc.CallOption) (*Block, error)
	ListChannels(ctx context.Context, in *ListChannelsRequest, opts ...grpc.CallOption) (*ChannelInfos, error)
	AddChannel(ctx context.Context, in *AddChannelRequest, opts ...grpc.CallOption) (*ChannelInfo, error)
	AddTx(ctx context.Context, in *AddTxRequest, opts ...grpc.CallOption) (*TxStatus, error)
}

type ordererClient struct {
	cc *grpc.ClientConn
}

func NewOrdererClient(cc *grpc.ClientConn) OrdererClient {
	return &ordererClient{cc}
}

func (c *ordererClient) FetchBlock(ctx context.Context, in *FetchBlockRequest, opts ...grpc.CallOption) (*Block, error) {
	out := new(Block)
	err := c.cc.Invoke(ctx, "/protos.Orderer/FetchBlock", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ordererClient) ListChannels(ctx context.Context, in *ListChannelsRequest, opts ...grpc.CallOption) (*ChannelInfos, error) {
	out := new(ChannelInfos)
	err := c.cc.Invoke(ctx, "/protos.Orderer/ListChannels", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ordererClient) AddChannel(ctx context.Context, in *AddChannelRequest, opts ...grpc.CallOption) (*ChannelInfo, error) {
	out := new(ChannelInfo)
	err := c.cc.Invoke(ctx, "/protos.Orderer/AddChannel", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ordererClient) AddTx(ctx context.Context, in *AddTxRequest, opts ...grpc.CallOption) (*TxStatus, error) {
	out := new(TxStatus)
	err := c.cc.Invoke(ctx, "/protos.Orderer/AddTx", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrdererServer is the server API for Orderer service.
type OrdererServer interface {
	FetchBlock(context.Context, *FetchBlockRequest) (*Block, error)
	ListChannels(context.Context, *ListChannelsRequest) (*ChannelInfos, error)
	AddChannel(context.Context, *AddChannelRequest) (*ChannelInfo, error)
	AddTx(context.Context, *AddTxRequest) (*TxStatus, error)
}

func RegisterOrdererServer(s *grpc.Server, srv OrdererServer) {
	s.RegisterService(&_Orderer_serviceDesc, srv)
}

func _Orderer_FetchBlock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchBlockRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrdererServer).FetchBlock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Orderer/FetchBlock",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrdererServer).FetchBlock(ctx, req.(*FetchBlockRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orderer_ListChannels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListChannelsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrdererServer).ListChannels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Orderer/ListChannels",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrdererServer).ListChannels(ctx, req.(*ListChannelsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orderer_AddChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrdererServer).AddChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Orderer/AddChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrdererServer).AddChannel(ctx, req.(*AddChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orderer_AddTx_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddTxRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrdererServer).AddTx(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Orderer/AddTx",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrdererServer).AddTx(ctx, req.(*AddTxRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Orderer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Orderer",
	HandlerType: (*OrdererServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FetchBlock",
			Handler:    _Orderer_FetchBlock_Handler,
		},
		{
			MethodName: "ListChannels",
			Handler:    _Orderer_ListChannels_Handler,
		},
		{
			MethodName: "AddChannel",
			Handler:    _Orderer_AddChannel_Handler,
		},
		{
			MethodName: "AddTx",
			Handler:    _Orderer_AddTx_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

// PeerClient is the client API for Peer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PeerClient interface {
}

type peerClient struct {
	cc *grpc.ClientConn
}

func NewPeerClient(cc *grpc.ClientConn) PeerClient {
	return &peerClient{cc}
}

// PeerServer is the server API for Peer service.
type PeerServer interface {
}

func RegisterPeerServer(s *grpc.Server, srv PeerServer) {
	s.RegisterService(&_Peer_serviceDesc, srv)
}

var _Peer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Peer",
	HandlerType: (*PeerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "service.proto",
}

func init() { proto.RegisterFile("service.proto", fileDescriptor_service_a015a6b48025b526) }

var fileDescriptor_service_a015a6b48025b526 = []byte{
	// 342 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xc1, 0x4e, 0xc2, 0x40,
	0x10, 0x86, 0x5b, 0xc4, 0x5a, 0xa6, 0x90, 0xc8, 0x42, 0x08, 0xae, 0x1c, 0xc8, 0x9e, 0x88, 0x89,
	0x18, 0xea, 0xc5, 0x93, 0x06, 0x31, 0x26, 0x24, 0x46, 0x4d, 0xe9, 0x0b, 0x00, 0x1d, 0x03, 0x11,
	0xa8, 0xee, 0x2e, 0xa6, 0xfa, 0xea, 0x5e, 0x8c, 0xdb, 0x2e, 0xdb, 0x84, 0x26, 0x7a, 0x6a, 0xe6,
	0x9f, 0x7f, 0xa6, 0x33, 0xdf, 0x2c, 0xd4, 0x04, 0xf2, 0x8f, 0xe5, 0x1c, 0xfb, 0x6f, 0x3c, 0x96,
	0x31, 0x71, 0xd4, 0x47, 0x50, 0x57, 0x26, 0xa9, 0x42, 0xbd, 0xd9, 0x2a, 0x9e, 0xbf, 0xa6, 0x01,
	0x1b, 0x43, 0xfd, 0x1e, 0xe5, 0x7c, 0x71, 0xfb, 0xab, 0x05, 0xf8, 0xbe, 0x45, 0x21, 0x49, 0x07,
	0x2a, 0xa3, 0xc5, 0x74, 0xb3, 0xc1, 0xd5, 0xf8, 0xae, 0x6d, 0x77, 0xed, 0x5e, 0x25, 0x30, 0x02,
	0x69, 0x81, 0xf3, 0xb8, 0x5d, 0xcf, 0x90, 0xb7, 0x4b, 0x5d, 0xbb, 0x57, 0x0e, 0xb2, 0x88, 0x9d,
	0x43, 0xe3, 0x61, 0x29, 0x64, 0x66, 0x14, 0xba, 0x59, 0x0b, 0x9c, 0xc9, 0xa7, 0x90, 0xb8, 0x56,
	0x9d, 0xdc, 0x20, 0x8b, 0xd8, 0x0d, 0x54, 0x75, 0xcf, 0xcd, 0x4b, 0x2c, 0xc8, 0x05, 0xb8, 0xba,
	0xb4, 0x6d, 0x77, 0x0f, 0x7a, 0x9e, 0xdf, 0x48, 0x67, 0x14, 0xfd, 0x9c, 0x2f, 0xd8, 0x99, 0xd8,
	0x18, 0xbc, 0x5c, 0xe2, 0x8f, 0xa1, 0x3b, 0x50, 0x51, 0x2b, 0x4e, 0x96, 0x5f, 0x98, 0xcd, 0x6d,
	0x04, 0x36, 0x80, 0xfa, 0x30, 0x8a, 0x32, 0xf7, 0xbf, 0x28, 0xb0, 0x33, 0xa8, 0x0e, 0xa3, 0x28,
	0x4c, 0xb4, 0x9b, 0x42, 0x29, 0x4c, 0x94, 0xcd, 0xf3, 0x41, 0x0f, 0x1e, 0x26, 0x41, 0x29, 0x4c,
	0x18, 0x80, 0x1b, 0x26, 0x13, 0x39, 0x95, 0x5b, 0xe1, 0x7f, 0xdb, 0x70, 0xf4, 0xc4, 0x23, 0xe4,
	0xc8, 0xc9, 0x15, 0x80, 0x81, 0x4f, 0x4e, 0x74, 0xd5, 0xde, 0x41, 0x68, 0x4d, 0xa7, 0x94, 0xca,
	0x2c, 0x32, 0x82, 0x6a, 0x9e, 0x35, 0x39, 0xd5, 0x86, 0x82, 0x0b, 0xd0, 0x66, 0x01, 0x47, 0xc1,
	0x2c, 0x72, 0x0d, 0x60, 0xb6, 0x36, 0xbf, 0xdf, 0x23, 0x41, 0x8b, 0x0e, 0xc1, 0x2c, 0x32, 0x80,
	0x43, 0x85, 0x80, 0x34, 0x73, 0xa5, 0x3b, 0x22, 0xf4, 0xd8, 0x50, 0x48, 0x77, 0x67, 0x96, 0xef,
	0x40, 0xf9, 0x19, 0x91, 0xcf, 0xd2, 0x57, 0x79, 0xf9, 0x13, 0x00, 0x00, 0xff, 0xff, 0xbc, 0xbf,
	0x4e, 0x0b, 0xad, 0x02, 0x00, 0x00,
}
