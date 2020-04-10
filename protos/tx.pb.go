// Code generated by protoc-gen-go. DO NOT EDIT.
// source: tx.proto

package protos

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Tx is the transaction, which structure is not decided yet
type Tx struct {
	ID                   string   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Data                 *TxData  `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data,omitempty"`
	Time                 int64    `protobuf:"varint,3,opt,name=Time,proto3" json:"Time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Tx) Reset()         { *m = Tx{} }
func (m *Tx) String() string { return proto.CompactTextString(m) }
func (*Tx) ProtoMessage()    {}
func (*Tx) Descriptor() ([]byte, []int) {
	return fileDescriptor_tx_71f4be22132760d7, []int{0}
}
func (m *Tx) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Tx.Unmarshal(m, b)
}
func (m *Tx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Tx.Marshal(b, m, deterministic)
}
func (dst *Tx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Tx.Merge(dst, src)
}
func (m *Tx) XXX_Size() int {
	return xxx_messageInfo_Tx.Size(m)
}
func (m *Tx) XXX_DiscardUnknown() {
	xxx_messageInfo_Tx.DiscardUnknown(m)
}

var xxx_messageInfo_Tx proto.InternalMessageInfo

func (m *Tx) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *Tx) GetData() *TxData {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Tx) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

// txData is the data of Tx
type TxData struct {
	ChannelID            string   `protobuf:"bytes,1,opt,name=ChannelID,proto3" json:"ChannelID,omitempty"`
	Nonce                uint64   `protobuf:"varint,2,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
	Recipient            []byte   `protobuf:"bytes,3,opt,name=Recipient,proto3" json:"Recipient,omitempty"`
	Payload              []byte   `protobuf:"bytes,4,opt,name=Payload,proto3" json:"Payload,omitempty"`
	Value                uint64   `protobuf:"varint,5,opt,name=Value,proto3" json:"Value,omitempty"`
	Msg                  string   `protobuf:"bytes,6,opt,name=Msg,proto3" json:"Msg,omitempty"`
	Version              int32    `protobuf:"varint,7,opt,name=Version,proto3" json:"Version,omitempty"`
	Sig                  *TxSig   `protobuf:"bytes,8,opt,name=Sig,proto3" json:"Sig,omitempty"`
	Gas                  uint64   `protobuf:"varint,9,opt,name=Gas,proto3" json:"Gas,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxData) Reset()         { *m = TxData{} }
func (m *TxData) String() string { return proto.CompactTextString(m) }
func (*TxData) ProtoMessage()    {}
func (*TxData) Descriptor() ([]byte, []int) {
	return fileDescriptor_tx_71f4be22132760d7, []int{1}
}
func (m *TxData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TxData.Unmarshal(m, b)
}
func (m *TxData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TxData.Marshal(b, m, deterministic)
}
func (dst *TxData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxData.Merge(dst, src)
}
func (m *TxData) XXX_Size() int {
	return xxx_messageInfo_TxData.Size(m)
}
func (m *TxData) XXX_DiscardUnknown() {
	xxx_messageInfo_TxData.DiscardUnknown(m)
}

var xxx_messageInfo_TxData proto.InternalMessageInfo

func (m *TxData) GetChannelID() string {
	if m != nil {
		return m.ChannelID
	}
	return ""
}

func (m *TxData) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *TxData) GetRecipient() []byte {
	if m != nil {
		return m.Recipient
	}
	return nil
}

func (m *TxData) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *TxData) GetValue() uint64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *TxData) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func (m *TxData) GetVersion() int32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *TxData) GetSig() *TxSig {
	if m != nil {
		return m.Sig
	}
	return nil
}

func (m *TxData) GetGas() uint64 {
	if m != nil {
		return m.Gas
	}
	return 0
}

// txSig is the sig of tx
// However, it has a gap with the struct of txSig,
// because there is no big int in the protobuf
type TxSig struct {
	PK                   []byte   `protobuf:"bytes,1,opt,name=PK,proto3" json:"PK,omitempty"`
	Sig                  []byte   `protobuf:"bytes,2,opt,name=Sig,proto3" json:"Sig,omitempty"`
	Algo                 int32    `protobuf:"varint,3,opt,name=Algo,proto3" json:"Algo,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxSig) Reset()         { *m = TxSig{} }
func (m *TxSig) String() string { return proto.CompactTextString(m) }
func (*TxSig) ProtoMessage()    {}
func (*TxSig) Descriptor() ([]byte, []int) {
	return fileDescriptor_tx_71f4be22132760d7, []int{2}
}
func (m *TxSig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TxSig.Unmarshal(m, b)
}
func (m *TxSig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TxSig.Marshal(b, m, deterministic)
}
func (dst *TxSig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxSig.Merge(dst, src)
}
func (m *TxSig) XXX_Size() int {
	return xxx_messageInfo_TxSig.Size(m)
}
func (m *TxSig) XXX_DiscardUnknown() {
	xxx_messageInfo_TxSig.DiscardUnknown(m)
}

var xxx_messageInfo_TxSig proto.InternalMessageInfo

func (m *TxSig) GetPK() []byte {
	if m != nil {
		return m.PK
	}
	return nil
}

func (m *TxSig) GetSig() []byte {
	if m != nil {
		return m.Sig
	}
	return nil
}

func (m *TxSig) GetAlgo() int32 {
	if m != nil {
		return m.Algo
	}
	return 0
}

func init() {
	proto.RegisterType((*Tx)(nil), "protos.Tx")
	proto.RegisterType((*TxData)(nil), "protos.txData")
	proto.RegisterType((*TxSig)(nil), "protos.txSig")
}

func init() { proto.RegisterFile("tx.proto", fileDescriptor_tx_71f4be22132760d7) }

var fileDescriptor_tx_71f4be22132760d7 = []byte{
	// 275 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x51, 0x4d, 0x4b, 0xc3, 0x40,
	0x10, 0x65, 0x37, 0x1f, 0x6d, 0xc6, 0x5a, 0x64, 0xf0, 0xb0, 0x07, 0xc1, 0x90, 0x53, 0x4e, 0x3d,
	0xe8, 0xd9, 0x83, 0x18, 0x90, 0x52, 0x95, 0xb0, 0x2d, 0xbd, 0xaf, 0x75, 0x89, 0x0b, 0x31, 0x5b,
	0x9a, 0x08, 0xf1, 0x27, 0xfb, 0x2f, 0x64, 0x66, 0xa9, 0x39, 0xe5, 0x7d, 0xcc, 0xbc, 0xe1, 0x65,
	0x61, 0x3e, 0x8c, 0xab, 0xe3, 0xc9, 0x0f, 0x1e, 0x53, 0xfe, 0xf4, 0xc5, 0x0b, 0xc8, 0xdd, 0x88,
	0x4b, 0x90, 0xeb, 0x4a, 0x89, 0x5c, 0x94, 0x99, 0x96, 0xeb, 0x0a, 0x0b, 0x88, 0x2b, 0x33, 0x18,
	0x25, 0x73, 0x51, 0x5e, 0xdc, 0x2d, 0xc3, 0x4e, 0xbf, 0x1a, 0x46, 0x52, 0x35, 0x7b, 0x88, 0x10,
	0xef, 0xdc, 0x97, 0x55, 0x51, 0x2e, 0xca, 0x48, 0x33, 0x2e, 0x7e, 0x05, 0xa4, 0x61, 0x08, 0x6f,
	0x20, 0x7b, 0xfa, 0x34, 0x5d, 0x67, 0xdb, 0xff, 0xe4, 0x49, 0xc0, 0x6b, 0x48, 0xde, 0x7c, 0x77,
	0xb0, 0x7c, 0x21, 0xd6, 0x81, 0xd0, 0x8e, 0xb6, 0x07, 0x77, 0x74, 0xb6, 0x1b, 0x38, 0x77, 0xa1,
	0x27, 0x01, 0x15, 0xcc, 0x6a, 0xf3, 0xd3, 0x7a, 0xf3, 0xa1, 0x62, 0xf6, 0xce, 0x94, 0xd2, 0xf6,
	0xa6, 0xfd, 0xb6, 0x2a, 0x09, 0x69, 0x4c, 0xf0, 0x0a, 0xa2, 0xd7, 0xbe, 0x51, 0x29, 0xdf, 0x26,
	0x48, 0x09, 0x7b, 0x7b, 0xea, 0x9d, 0xef, 0xd4, 0x2c, 0x17, 0x65, 0xa2, 0xcf, 0x14, 0x6f, 0x21,
	0xda, 0xba, 0x46, 0xcd, 0xb9, 0xef, 0xe5, 0xd4, 0x77, 0xeb, 0x1a, 0x4d, 0x0e, 0x85, 0x3d, 0x9b,
	0x5e, 0x65, 0x7c, 0x80, 0x60, 0xf1, 0x00, 0x09, 0xfb, 0xf4, 0xf3, 0xea, 0x0d, 0x57, 0x5c, 0x68,
	0x59, 0x6f, 0x68, 0x94, 0xb2, 0x24, 0x0b, 0xbc, 0x8c, 0x10, 0x3f, 0xb6, 0x8d, 0xe7, 0x4a, 0x89,
	0x66, 0xfc, 0x1e, 0x1e, 0xe0, 0xfe, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x86, 0x81, 0x31, 0xe9, 0x93,
	0x01, 0x00, 0x00,
}
