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
	return fileDescriptor_tx_b7f74a5a9235be2e, []int{0}
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
	return fileDescriptor_tx_b7f74a5a9235be2e, []int{1}
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
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxSig) Reset()         { *m = TxSig{} }
func (m *TxSig) String() string { return proto.CompactTextString(m) }
func (*TxSig) ProtoMessage()    {}
func (*TxSig) Descriptor() ([]byte, []int) {
	return fileDescriptor_tx_b7f74a5a9235be2e, []int{2}
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

func init() {
	proto.RegisterType((*Tx)(nil), "protos.Tx")
	proto.RegisterType((*TxData)(nil), "protos.txData")
	proto.RegisterType((*TxSig)(nil), "protos.txSig")
}

func init() { proto.RegisterFile("tx.proto", fileDescriptor_tx_b7f74a5a9235be2e) }

var fileDescriptor_tx_b7f74a5a9235be2e = []byte{
	// 257 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x50, 0xcb, 0x4e, 0xc3, 0x30,
	0x10, 0x94, 0x9d, 0x47, 0xdb, 0xa5, 0x54, 0x68, 0xc5, 0xc1, 0x07, 0x24, 0xa2, 0x9c, 0xc2, 0xa5,
	0x07, 0xf8, 0x04, 0x72, 0xa9, 0x0a, 0x28, 0x72, 0xab, 0xde, 0x4d, 0xb1, 0x82, 0xa5, 0x10, 0x57,
	0x8d, 0x91, 0xc2, 0x97, 0xf2, 0x3b, 0x68, 0xd7, 0x2a, 0x39, 0x79, 0x67, 0xd6, 0x33, 0xa3, 0x59,
	0x98, 0x87, 0x71, 0x7d, 0x3a, 0xfb, 0xe0, 0x31, 0xe7, 0x67, 0x28, 0x5f, 0x40, 0xee, 0x47, 0x5c,
	0x81, 0xdc, 0xd4, 0x4a, 0x14, 0xa2, 0x5a, 0x68, 0xb9, 0xa9, 0xb1, 0x84, 0xb4, 0x36, 0xc1, 0x28,
	0x59, 0x88, 0xea, 0xea, 0x71, 0x15, 0x35, 0xc3, 0x3a, 0x8c, 0xc4, 0x6a, 0xde, 0x21, 0x42, 0xba,
	0x77, 0x5f, 0x56, 0x25, 0x85, 0xa8, 0x12, 0xcd, 0x73, 0xf9, 0x2b, 0x20, 0x8f, 0x9f, 0xf0, 0x0e,
	0x16, 0xcf, 0x9f, 0xa6, 0xef, 0x6d, 0xf7, 0xef, 0x3c, 0x11, 0x78, 0x0b, 0xd9, 0x9b, 0xef, 0x8f,
	0x96, 0x13, 0x52, 0x1d, 0x01, 0x69, 0xb4, 0x3d, 0xba, 0x93, 0xb3, 0x7d, 0x60, 0xdf, 0xa5, 0x9e,
	0x08, 0x54, 0x30, 0x6b, 0xcc, 0x4f, 0xe7, 0xcd, 0x87, 0x4a, 0x79, 0x77, 0x81, 0xe4, 0x76, 0x30,
	0xdd, 0xb7, 0x55, 0x59, 0x74, 0x63, 0x80, 0x37, 0x90, 0xbc, 0x0e, 0xad, 0xca, 0x39, 0x9b, 0x46,
	0x72, 0x38, 0xd8, 0xf3, 0xe0, 0x7c, 0xaf, 0x66, 0x85, 0xa8, 0x32, 0x7d, 0x81, 0x78, 0x0f, 0xc9,
	0xce, 0xb5, 0x6a, 0xce, 0x7d, 0xaf, 0xa7, 0xbe, 0x3b, 0xd7, 0x6a, 0xda, 0x94, 0x0f, 0x90, 0x31,
	0xa2, 0x53, 0x35, 0x5b, 0x2e, 0xb4, 0xd4, 0xb2, 0xd9, 0x52, 0x0a, 0x29, 0x25, 0x13, 0x34, 0xbe,
	0xc7, 0xd3, 0x3e, 0xfd, 0x05, 0x00, 0x00, 0xff, 0xff, 0x9f, 0x8f, 0x9d, 0x4c, 0x6d, 0x01, 0x00,
	0x00,
}
