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

// Behavior defines the behavior
type Behavior int32

const (
	// Fail right away if not exist
	Behavior_FAIL_IF_NOT_READY Behavior = 0
	// Return until ready
	Behavior_RETURN_UNTIL_READY Behavior = 1
)

var Behavior_name = map[int32]string{
	0: "FAIL_IF_NOT_READY",
	1: "RETURN_UNTIL_READY",
}
var Behavior_value = map[string]int32{
	"FAIL_IF_NOT_READY":  0,
	"RETURN_UNTIL_READY": 1,
}

func (x Behavior) String() string {
	return proto.EnumName(Behavior_name, int32(x))
}
func (Behavior) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_service_6c32c62bd3c6bf16, []int{0}
}

// Identity defines the identity in the channel
type Identity int32

const (
	Identity_MEMBER   Identity = 0
	Identity_ADMIN    Identity = 1
	Identity_OUTSIDER Identity = 2
)

var Identity_name = map[int32]string{
	0: "MEMBER",
	1: "ADMIN",
	2: "OUTSIDER",
}
var Identity_value = map[string]int32{
	"MEMBER":   0,
	"ADMIN":    1,
	"OUTSIDER": 2,
}

func (x Identity) String() string {
	return proto.EnumName(Identity_name, int32(x))
}
func (Identity) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_service_6c32c62bd3c6bf16, []int{1}
}

// However, this is not contains sig now, but this is necessary
// if we want to verify the permission.
// TODO: add sig to identity.
type FetchBlockRequest struct {
	ChannelID            string   `protobuf:"bytes,1,opt,name=ChannelID,proto3" json:"ChannelID,omitempty"`
	Number               uint64   `protobuf:"varint,2,opt,name=Number,proto3" json:"Number,omitempty"`
	Behavior             Behavior `protobuf:"varint,3,opt,name=Behavior,proto3,enum=protos.Behavior" json:"Behavior,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FetchBlockRequest) Reset()         { *m = FetchBlockRequest{} }
func (m *FetchBlockRequest) String() string { return proto.CompactTextString(m) }
func (*FetchBlockRequest) ProtoMessage()    {}
func (*FetchBlockRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_6c32c62bd3c6bf16, []int{0}
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

func (m *FetchBlockRequest) GetBehavior() Behavior {
	if m != nil {
		return m.Behavior
	}
	return Behavior_FAIL_IF_NOT_READY
}

// TODO: This should contain signature.
type ListChannelsRequest struct {
	// If system channel are included
	System               bool     `protobuf:"varint,1,opt,name=System,proto3" json:"System,omitempty"`
	PK                   []byte   `protobuf:"bytes,2,opt,name=PK,proto3" json:"PK,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListChannelsRequest) Reset()         { *m = ListChannelsRequest{} }
func (m *ListChannelsRequest) String() string { return proto.CompactTextString(m) }
func (*ListChannelsRequest) ProtoMessage()    {}
func (*ListChannelsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_6c32c62bd3c6bf16, []int{1}
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

func (m *ListChannelsRequest) GetPK() []byte {
	if m != nil {
		return m.PK
	}
	return nil
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
	return fileDescriptor_service_6c32c62bd3c6bf16, []int{2}
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
	Identity             Identity `protobuf:"varint,3,opt,name=Identity,proto3,enum=protos.Identity" json:"Identity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChannelInfo) Reset()         { *m = ChannelInfo{} }
func (m *ChannelInfo) String() string { return proto.CompactTextString(m) }
func (*ChannelInfo) ProtoMessage()    {}
func (*ChannelInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_6c32c62bd3c6bf16, []int{3}
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

func (m *ChannelInfo) GetIdentity() Identity {
	if m != nil {
		return m.Identity
	}
	return Identity_MEMBER
}

// CreateChannelRequest include a special tx which create a channel.
type CreateChannelRequest struct {
	Tx                   *Tx      `protobuf:"bytes,1,opt,name=Tx,proto3" json:"Tx,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateChannelRequest) Reset()         { *m = CreateChannelRequest{} }
func (m *CreateChannelRequest) String() string { return proto.CompactTextString(m) }
func (*CreateChannelRequest) ProtoMessage()    {}
func (*CreateChannelRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_6c32c62bd3c6bf16, []int{4}
}
func (m *CreateChannelRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateChannelRequest.Unmarshal(m, b)
}
func (m *CreateChannelRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateChannelRequest.Marshal(b, m, deterministic)
}
func (dst *CreateChannelRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateChannelRequest.Merge(dst, src)
}
func (m *CreateChannelRequest) XXX_Size() int {
	return xxx_messageInfo_CreateChannelRequest.Size(m)
}
func (m *CreateChannelRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateChannelRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateChannelRequest proto.InternalMessageInfo

func (m *CreateChannelRequest) GetTx() *Tx {
	if m != nil {
		return m.Tx
	}
	return nil
}

// CreateChannelTxPayload is the payload of create channel tx
type CreateChannelTxPayload struct {
	ChannelID            string   `protobuf:"bytes,1,opt,name=ChannelID,proto3" json:"ChannelID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateChannelTxPayload) Reset()         { *m = CreateChannelTxPayload{} }
func (m *CreateChannelTxPayload) String() string { return proto.CompactTextString(m) }
func (*CreateChannelTxPayload) ProtoMessage()    {}
func (*CreateChannelTxPayload) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_6c32c62bd3c6bf16, []int{5}
}
func (m *CreateChannelTxPayload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateChannelTxPayload.Unmarshal(m, b)
}
func (m *CreateChannelTxPayload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateChannelTxPayload.Marshal(b, m, deterministic)
}
func (dst *CreateChannelTxPayload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateChannelTxPayload.Merge(dst, src)
}
func (m *CreateChannelTxPayload) XXX_Size() int {
	return xxx_messageInfo_CreateChannelTxPayload.Size(m)
}
func (m *CreateChannelTxPayload) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateChannelTxPayload.DiscardUnknown(m)
}

var xxx_messageInfo_CreateChannelTxPayload proto.InternalMessageInfo

func (m *CreateChannelTxPayload) GetChannelID() string {
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
	return fileDescriptor_service_6c32c62bd3c6bf16, []int{6}
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
	Err                  string   `protobuf:"bytes,1,opt,name=Err,proto3" json:"Err,omitempty"`
	BlockNumber          uint64   `protobuf:"varint,2,opt,name=BlockNumber,proto3" json:"BlockNumber,omitempty"`
	BlockIndex           int32    `protobuf:"varint,3,opt,name=BlockIndex,proto3" json:"BlockIndex,omitempty"`
	Output               []byte   `protobuf:"bytes,4,opt,name=Output,proto3" json:"Output,omitempty"`
	ContractAddress      string   `protobuf:"bytes,5,opt,name=ContractAddress,proto3" json:"ContractAddress,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxStatus) Reset()         { *m = TxStatus{} }
func (m *TxStatus) String() string { return proto.CompactTextString(m) }
func (*TxStatus) ProtoMessage()    {}
func (*TxStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_6c32c62bd3c6bf16, []int{7}
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

func (m *TxStatus) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

func (m *TxStatus) GetBlockNumber() uint64 {
	if m != nil {
		return m.BlockNumber
	}
	return 0
}

func (m *TxStatus) GetBlockIndex() int32 {
	if m != nil {
		return m.BlockIndex
	}
	return 0
}

func (m *TxStatus) GetOutput() []byte {
	if m != nil {
		return m.Output
	}
	return nil
}

func (m *TxStatus) GetContractAddress() string {
	if m != nil {
		return m.ContractAddress
	}
	return ""
}

type GetTxStatusRequest struct {
	ChannelID            string   `protobuf:"bytes,1,opt,name=ChannelID,proto3" json:"ChannelID,omitempty"`
	TxID                 string   `protobuf:"bytes,2,opt,name=TxID,proto3" json:"TxID,omitempty"`
	Behavior             Behavior `protobuf:"varint,3,opt,name=Behavior,proto3,enum=protos.Behavior" json:"Behavior,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTxStatusRequest) Reset()         { *m = GetTxStatusRequest{} }
func (m *GetTxStatusRequest) String() string { return proto.CompactTextString(m) }
func (*GetTxStatusRequest) ProtoMessage()    {}
func (*GetTxStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_6c32c62bd3c6bf16, []int{8}
}
func (m *GetTxStatusRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTxStatusRequest.Unmarshal(m, b)
}
func (m *GetTxStatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTxStatusRequest.Marshal(b, m, deterministic)
}
func (dst *GetTxStatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTxStatusRequest.Merge(dst, src)
}
func (m *GetTxStatusRequest) XXX_Size() int {
	return xxx_messageInfo_GetTxStatusRequest.Size(m)
}
func (m *GetTxStatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTxStatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetTxStatusRequest proto.InternalMessageInfo

func (m *GetTxStatusRequest) GetChannelID() string {
	if m != nil {
		return m.ChannelID
	}
	return ""
}

func (m *GetTxStatusRequest) GetTxID() string {
	if m != nil {
		return m.TxID
	}
	return ""
}

func (m *GetTxStatusRequest) GetBehavior() Behavior {
	if m != nil {
		return m.Behavior
	}
	return Behavior_FAIL_IF_NOT_READY
}

type ListTxHistoryRequest struct {
	Address              []byte   `protobuf:"bytes,1,opt,name=Address,proto3" json:"Address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListTxHistoryRequest) Reset()         { *m = ListTxHistoryRequest{} }
func (m *ListTxHistoryRequest) String() string { return proto.CompactTextString(m) }
func (*ListTxHistoryRequest) ProtoMessage()    {}
func (*ListTxHistoryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_6c32c62bd3c6bf16, []int{9}
}
func (m *ListTxHistoryRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListTxHistoryRequest.Unmarshal(m, b)
}
func (m *ListTxHistoryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListTxHistoryRequest.Marshal(b, m, deterministic)
}
func (dst *ListTxHistoryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListTxHistoryRequest.Merge(dst, src)
}
func (m *ListTxHistoryRequest) XXX_Size() int {
	return xxx_messageInfo_ListTxHistoryRequest.Size(m)
}
func (m *ListTxHistoryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListTxHistoryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListTxHistoryRequest proto.InternalMessageInfo

func (m *ListTxHistoryRequest) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

// TxHistory includes all txs
type TxHistory struct {
	// repeated string Txs = 1;
	Txs                  map[string]*StringList `protobuf:"bytes,1,rep,name=Txs,proto3" json:"Txs,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *TxHistory) Reset()         { *m = TxHistory{} }
func (m *TxHistory) String() string { return proto.CompactTextString(m) }
func (*TxHistory) ProtoMessage()    {}
func (*TxHistory) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_6c32c62bd3c6bf16, []int{10}
}
func (m *TxHistory) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TxHistory.Unmarshal(m, b)
}
func (m *TxHistory) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TxHistory.Marshal(b, m, deterministic)
}
func (dst *TxHistory) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxHistory.Merge(dst, src)
}
func (m *TxHistory) XXX_Size() int {
	return xxx_messageInfo_TxHistory.Size(m)
}
func (m *TxHistory) XXX_DiscardUnknown() {
	xxx_messageInfo_TxHistory.DiscardUnknown(m)
}

var xxx_messageInfo_TxHistory proto.InternalMessageInfo

func (m *TxHistory) GetTxs() map[string]*StringList {
	if m != nil {
		return m.Txs
	}
	return nil
}

type GetAccountInfoRequest struct {
	Address              []byte   `protobuf:"bytes,1,opt,name=Address,proto3" json:"Address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetAccountInfoRequest) Reset()         { *m = GetAccountInfoRequest{} }
func (m *GetAccountInfoRequest) String() string { return proto.CompactTextString(m) }
func (*GetAccountInfoRequest) ProtoMessage()    {}
func (*GetAccountInfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_6c32c62bd3c6bf16, []int{11}
}
func (m *GetAccountInfoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetAccountInfoRequest.Unmarshal(m, b)
}
func (m *GetAccountInfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetAccountInfoRequest.Marshal(b, m, deterministic)
}
func (dst *GetAccountInfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetAccountInfoRequest.Merge(dst, src)
}
func (m *GetAccountInfoRequest) XXX_Size() int {
	return xxx_messageInfo_GetAccountInfoRequest.Size(m)
}
func (m *GetAccountInfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetAccountInfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetAccountInfoRequest proto.InternalMessageInfo

func (m *GetAccountInfoRequest) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

type AccountInfo struct {
	Balance              uint64   `protobuf:"varint,1,opt,name=Balance,proto3" json:"Balance,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AccountInfo) Reset()         { *m = AccountInfo{} }
func (m *AccountInfo) String() string { return proto.CompactTextString(m) }
func (*AccountInfo) ProtoMessage()    {}
func (*AccountInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_service_6c32c62bd3c6bf16, []int{12}
}
func (m *AccountInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountInfo.Unmarshal(m, b)
}
func (m *AccountInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountInfo.Marshal(b, m, deterministic)
}
func (dst *AccountInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountInfo.Merge(dst, src)
}
func (m *AccountInfo) XXX_Size() int {
	return xxx_messageInfo_AccountInfo.Size(m)
}
func (m *AccountInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountInfo.DiscardUnknown(m)
}

var xxx_messageInfo_AccountInfo proto.InternalMessageInfo

func (m *AccountInfo) GetBalance() uint64 {
	if m != nil {
		return m.Balance
	}
	return 0
}

func init() {
	proto.RegisterType((*FetchBlockRequest)(nil), "protos.FetchBlockRequest")
	proto.RegisterType((*ListChannelsRequest)(nil), "protos.ListChannelsRequest")
	proto.RegisterType((*ChannelInfos)(nil), "protos.ChannelInfos")
	proto.RegisterType((*ChannelInfo)(nil), "protos.ChannelInfo")
	proto.RegisterType((*CreateChannelRequest)(nil), "protos.CreateChannelRequest")
	proto.RegisterType((*CreateChannelTxPayload)(nil), "protos.CreateChannelTxPayload")
	proto.RegisterType((*AddTxRequest)(nil), "protos.AddTxRequest")
	proto.RegisterType((*TxStatus)(nil), "protos.TxStatus")
	proto.RegisterType((*GetTxStatusRequest)(nil), "protos.GetTxStatusRequest")
	proto.RegisterType((*ListTxHistoryRequest)(nil), "protos.ListTxHistoryRequest")
	proto.RegisterType((*TxHistory)(nil), "protos.TxHistory")
	proto.RegisterMapType((map[string]*StringList)(nil), "protos.TxHistory.TxsEntry")
	proto.RegisterType((*GetAccountInfoRequest)(nil), "protos.GetAccountInfoRequest")
	proto.RegisterType((*AccountInfo)(nil), "protos.AccountInfo")
	proto.RegisterEnum("protos.Behavior", Behavior_name, Behavior_value)
	proto.RegisterEnum("protos.Identity", Identity_name, Identity_value)
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
	CreateChannel(ctx context.Context, in *CreateChannelRequest, opts ...grpc.CallOption) (*ChannelInfo, error)
	AddTx(ctx context.Context, in *AddTxRequest, opts ...grpc.CallOption) (*TxStatus, error)
	GetAccountInfo(ctx context.Context, in *GetAccountInfoRequest, opts ...grpc.CallOption) (*AccountInfo, error)
	GetTxStatus(ctx context.Context, in *GetTxStatusRequest, opts ...grpc.CallOption) (*TxStatus, error)
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

func (c *ordererClient) CreateChannel(ctx context.Context, in *CreateChannelRequest, opts ...grpc.CallOption) (*ChannelInfo, error) {
	out := new(ChannelInfo)
	err := c.cc.Invoke(ctx, "/protos.Orderer/CreateChannel", in, out, opts...)
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

func (c *ordererClient) GetAccountInfo(ctx context.Context, in *GetAccountInfoRequest, opts ...grpc.CallOption) (*AccountInfo, error) {
	out := new(AccountInfo)
	err := c.cc.Invoke(ctx, "/protos.Orderer/GetAccountInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ordererClient) GetTxStatus(ctx context.Context, in *GetTxStatusRequest, opts ...grpc.CallOption) (*TxStatus, error) {
	out := new(TxStatus)
	err := c.cc.Invoke(ctx, "/protos.Orderer/GetTxStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrdererServer is the server API for Orderer service.
type OrdererServer interface {
	FetchBlock(context.Context, *FetchBlockRequest) (*Block, error)
	ListChannels(context.Context, *ListChannelsRequest) (*ChannelInfos, error)
	CreateChannel(context.Context, *CreateChannelRequest) (*ChannelInfo, error)
	AddTx(context.Context, *AddTxRequest) (*TxStatus, error)
	GetAccountInfo(context.Context, *GetAccountInfoRequest) (*AccountInfo, error)
	GetTxStatus(context.Context, *GetTxStatusRequest) (*TxStatus, error)
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

func _Orderer_CreateChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrdererServer).CreateChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Orderer/CreateChannel",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrdererServer).CreateChannel(ctx, req.(*CreateChannelRequest))
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

func _Orderer_GetAccountInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAccountInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrdererServer).GetAccountInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Orderer/GetAccountInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrdererServer).GetAccountInfo(ctx, req.(*GetAccountInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Orderer_GetTxStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTxStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrdererServer).GetTxStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Orderer/GetTxStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrdererServer).GetTxStatus(ctx, req.(*GetTxStatusRequest))
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
			MethodName: "CreateChannel",
			Handler:    _Orderer_CreateChannel_Handler,
		},
		{
			MethodName: "AddTx",
			Handler:    _Orderer_AddTx_Handler,
		},
		{
			MethodName: "GetAccountInfo",
			Handler:    _Orderer_GetAccountInfo_Handler,
		},
		{
			MethodName: "GetTxStatus",
			Handler:    _Orderer_GetTxStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

// PeerClient is the client API for Peer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PeerClient interface {
	GetTxStatus(ctx context.Context, in *GetTxStatusRequest, opts ...grpc.CallOption) (*TxStatus, error)
	ListTxHistory(ctx context.Context, in *ListTxHistoryRequest, opts ...grpc.CallOption) (*TxHistory, error)
}

type peerClient struct {
	cc *grpc.ClientConn
}

func NewPeerClient(cc *grpc.ClientConn) PeerClient {
	return &peerClient{cc}
}

func (c *peerClient) GetTxStatus(ctx context.Context, in *GetTxStatusRequest, opts ...grpc.CallOption) (*TxStatus, error) {
	out := new(TxStatus)
	err := c.cc.Invoke(ctx, "/protos.Peer/GetTxStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *peerClient) ListTxHistory(ctx context.Context, in *ListTxHistoryRequest, opts ...grpc.CallOption) (*TxHistory, error) {
	out := new(TxHistory)
	err := c.cc.Invoke(ctx, "/protos.Peer/ListTxHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PeerServer is the server API for Peer service.
type PeerServer interface {
	GetTxStatus(context.Context, *GetTxStatusRequest) (*TxStatus, error)
	ListTxHistory(context.Context, *ListTxHistoryRequest) (*TxHistory, error)
}

func RegisterPeerServer(s *grpc.Server, srv PeerServer) {
	s.RegisterService(&_Peer_serviceDesc, srv)
}

func _Peer_GetTxStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTxStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerServer).GetTxStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Peer/GetTxStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerServer).GetTxStatus(ctx, req.(*GetTxStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Peer_ListTxHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTxHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PeerServer).ListTxHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Peer/ListTxHistory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PeerServer).ListTxHistory(ctx, req.(*ListTxHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Peer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Peer",
	HandlerType: (*PeerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTxStatus",
			Handler:    _Peer_GetTxStatus_Handler,
		},
		{
			MethodName: "ListTxHistory",
			Handler:    _Peer_ListTxHistory_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

func init() { proto.RegisterFile("service.proto", fileDescriptor_service_6c32c62bd3c6bf16) }

var fileDescriptor_service_6c32c62bd3c6bf16 = []byte{
	// 767 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0x5f, 0x6f, 0xda, 0x56,
	0x14, 0xb7, 0x1d, 0x48, 0xe1, 0x18, 0x32, 0x72, 0x9a, 0x22, 0xe6, 0x65, 0x13, 0xba, 0x2f, 0x43,
	0x51, 0x95, 0xae, 0x4c, 0x9a, 0xba, 0x49, 0xd5, 0xc4, 0xdf, 0xce, 0x6b, 0x02, 0xe8, 0xe2, 0x3c,
	0xec, 0x09, 0x39, 0xf6, 0xdd, 0x82, 0x02, 0x76, 0x77, 0x7d, 0x49, 0xcd, 0xde, 0xf7, 0xb0, 0x6f,
	0xb1, 0xaf, 0xb3, 0x6f, 0x35, 0xf9, 0xda, 0x17, 0x9b, 0x14, 0xad, 0xdd, 0x13, 0x3e, 0x3f, 0xff,
	0x7c, 0x7e, 0xe7, 0x3f, 0x50, 0x8f, 0x18, 0x7f, 0x58, 0x7a, 0xec, 0xf2, 0x1d, 0x0f, 0x45, 0x88,
	0xc7, 0xf2, 0x27, 0xb2, 0x6a, 0x5e, 0xb8, 0x5e, 0x87, 0x41, 0x8a, 0x5a, 0x15, 0x11, 0x67, 0x4f,
	0xe6, 0xed, 0x2a, 0xf4, 0xee, 0x53, 0x83, 0xbc, 0x87, 0xd3, 0x31, 0x13, 0xde, 0x5d, 0x3f, 0xc1,
	0x28, 0xfb, 0x7d, 0xc3, 0x22, 0x81, 0xe7, 0x50, 0x1d, 0xdc, 0xb9, 0x41, 0xc0, 0x56, 0xf6, 0xb0,
	0xa5, 0xb7, 0xf5, 0x4e, 0x95, 0xe6, 0x00, 0x36, 0xe1, 0x78, 0xb2, 0x59, 0xdf, 0x32, 0xde, 0x32,
	0xda, 0x7a, 0xa7, 0x44, 0x33, 0x0b, 0x9f, 0x43, 0xa5, 0xcf, 0xee, 0xdc, 0x87, 0x65, 0xc8, 0x5b,
	0x47, 0x6d, 0xbd, 0x73, 0xd2, 0x6d, 0xa4, 0x22, 0xd1, 0xa5, 0xc2, 0xe9, 0x8e, 0x41, 0x5e, 0xc3,
	0xd3, 0xab, 0x65, 0x24, 0x32, 0xb7, 0x91, 0x92, 0x6e, 0xc2, 0xf1, 0x7c, 0x1b, 0x09, 0xb6, 0x96,
	0xba, 0x15, 0x9a, 0x59, 0x78, 0x02, 0xc6, 0xec, 0xad, 0x14, 0xac, 0x51, 0x63, 0xf6, 0x96, 0xfc,
	0x08, 0x35, 0x15, 0x51, 0xf0, 0x6b, 0x18, 0xe1, 0x0b, 0xa8, 0x28, 0x57, 0x2d, 0xbd, 0x7d, 0xd4,
	0x31, 0xbb, 0x4f, 0x95, 0x78, 0x81, 0x47, 0x77, 0x24, 0xf2, 0x1e, 0xcc, 0xc2, 0x8b, 0x8f, 0xa4,
	0x7c, 0x0e, 0x55, 0x59, 0xa0, 0xf9, 0xf2, 0x0f, 0x96, 0x65, 0x9d, 0x03, 0x49, 0xe2, 0xb6, 0xcf,
	0x02, 0xb1, 0x14, 0xdb, 0xc7, 0x89, 0x2b, 0x9c, 0xee, 0x18, 0xa4, 0x0b, 0x67, 0x03, 0xce, 0x5c,
	0xc1, 0x32, 0xf7, 0x2a, 0x73, 0x0b, 0x0c, 0x27, 0x96, 0xd2, 0x66, 0x17, 0xd4, 0xf7, 0x4e, 0x4c,
	0x0d, 0x27, 0x26, 0xdf, 0x41, 0x73, 0xef, 0x1b, 0x27, 0x9e, 0xb9, 0xdb, 0x55, 0xe8, 0xfa, 0xff,
	0x1d, 0x37, 0xb9, 0x80, 0x5a, 0xcf, 0xf7, 0x9d, 0xf8, 0x53, 0x34, 0xfe, 0xd6, 0xa1, 0xe2, 0xc4,
	0x73, 0xe1, 0x8a, 0x4d, 0x84, 0x0d, 0x38, 0x1a, 0x71, 0x9e, 0x39, 0x4c, 0x1e, 0xb1, 0x0d, 0xa6,
	0xcc, 0x78, 0xaf, 0xf5, 0x45, 0x08, 0xbf, 0x02, 0x90, 0xa6, 0x1d, 0xf8, 0x2c, 0x96, 0x85, 0x28,
	0xd3, 0x02, 0x92, 0xb4, 0x76, 0xba, 0x11, 0xef, 0x36, 0xa2, 0x55, 0x92, 0x6d, 0xcc, 0x2c, 0xec,
	0xc0, 0x67, 0x83, 0x30, 0x10, 0xdc, 0xf5, 0x44, 0xcf, 0xf7, 0x39, 0x8b, 0xa2, 0x56, 0x59, 0xea,
	0x3e, 0x86, 0x89, 0x00, 0x7c, 0xc3, 0x84, 0x0a, 0xf2, 0xd3, 0xa6, 0x15, 0xa1, 0xe4, 0xc4, 0xf6,
	0x50, 0x06, 0x5c, 0xa5, 0xf2, 0xf9, 0x7f, 0x4e, 0xea, 0x37, 0x70, 0x96, 0x4c, 0xaa, 0x13, 0xff,
	0xb4, 0x8c, 0x44, 0xc8, 0xb7, 0x4a, 0xb7, 0x05, 0x4f, 0x54, 0xbc, 0xba, 0x4c, 0x48, 0x99, 0xe4,
	0x4f, 0x1d, 0xaa, 0x3b, 0x3a, 0x3e, 0x87, 0x23, 0x27, 0x56, 0x53, 0x69, 0xe5, 0x55, 0xcf, 0xde,
	0x5f, 0x3a, 0x71, 0x34, 0x0a, 0x04, 0xdf, 0xd2, 0x84, 0x66, 0xfd, 0x9c, 0x74, 0x21, 0x05, 0x92,
	0x2e, 0xdc, 0xb3, 0xad, 0xea, 0xc2, 0x3d, 0xdb, 0x62, 0x07, 0xca, 0x0f, 0xee, 0x6a, 0x93, 0x0e,
	0xa1, 0xd9, 0x45, 0xe5, 0x6d, 0x2e, 0xf8, 0x32, 0xf8, 0x2d, 0x09, 0x93, 0xa6, 0x84, 0x1f, 0x8c,
	0x57, 0x3a, 0x79, 0x09, 0xcf, 0xde, 0x30, 0xd1, 0xf3, 0xbc, 0x70, 0x13, 0x08, 0x39, 0xff, 0x1f,
	0x0d, 0xfd, 0x6b, 0x30, 0x0b, 0xfc, 0x84, 0xd8, 0x77, 0x57, 0x6e, 0xe0, 0x31, 0x49, 0x2c, 0x51,
	0x65, 0x5e, 0x7c, 0x9f, 0xd7, 0x10, 0x9f, 0xc1, 0xe9, 0xb8, 0x67, 0x5f, 0x2d, 0xec, 0xf1, 0x62,
	0x32, 0x75, 0x16, 0x74, 0xd4, 0x1b, 0xfe, 0xd2, 0xd0, 0xb0, 0x09, 0x48, 0x47, 0xce, 0x0d, 0x9d,
	0x2c, 0x6e, 0x26, 0x8e, 0x7d, 0x95, 0xe1, 0xfa, 0xc5, 0x8b, 0x7c, 0x5f, 0x10, 0xe0, 0xf8, 0x7a,
	0x74, 0xdd, 0x1f, 0xd1, 0x86, 0x86, 0x55, 0x28, 0xf7, 0x86, 0xd7, 0xf6, 0xa4, 0xa1, 0x63, 0x0d,
	0x2a, 0xd3, 0x1b, 0x67, 0x6e, 0x0f, 0x47, 0xb4, 0x61, 0x74, 0xff, 0x31, 0xe0, 0xc9, 0x94, 0xfb,
	0x8c, 0x33, 0x8e, 0xaf, 0x00, 0xf2, 0x83, 0x85, 0x9f, 0xab, 0x02, 0x7c, 0x70, 0xc4, 0xac, 0xfa,
	0xae, 0xa5, 0x09, 0x4a, 0x34, 0x1c, 0x40, 0xad, 0x78, 0x71, 0xf0, 0x0b, 0x45, 0x38, 0x70, 0x87,
	0xac, 0xb3, 0x03, 0xd7, 0x23, 0x22, 0x1a, 0x0e, 0xa1, 0xbe, 0xb7, 0x89, 0x78, 0xbe, 0x23, 0x1e,
	0x58, 0x6a, 0xeb, 0xd0, 0x11, 0x22, 0x1a, 0xbe, 0x84, 0xb2, 0xdc, 0x4b, 0xdc, 0xc9, 0x14, 0xd7,
	0xd4, 0x6a, 0xe4, 0x43, 0x92, 0x8e, 0x3a, 0xd1, 0x70, 0x0c, 0x27, 0xfb, 0xbd, 0xc4, 0x2f, 0x15,
	0xeb, 0x60, 0x8f, 0x73, 0xe9, 0xc2, 0x3b, 0xa2, 0x75, 0xff, 0xd2, 0xa1, 0x34, 0x63, 0x8c, 0xe3,
	0x6b, 0x30, 0x0b, 0xcb, 0x84, 0x56, 0xc1, 0xdb, 0xa3, 0x0d, 0x3b, 0x18, 0x4f, 0x1f, 0xea, 0x7b,
	0x5b, 0x91, 0x17, 0xe2, 0xd0, 0xb2, 0x58, 0xa7, 0x1f, 0xcc, 0x3d, 0xd1, 0x6e, 0xd3, 0x7f, 0xaa,
	0x6f, 0xff, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x67, 0xd0, 0xeb, 0x5f, 0xc1, 0x06, 0x00, 0x00,
}
