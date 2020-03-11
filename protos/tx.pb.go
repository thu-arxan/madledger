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
	return fileDescriptor_tx_ce29446be3a85ab1, []int{0}
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
	return fileDescriptor_tx_ce29446be3a85ab1, []int{1}
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
	return fileDescriptor_tx_ce29446be3a85ab1, []int{2}
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

func init() { proto.RegisterFile("tx.proto", fileDescriptor_tx_ce29446be3a85ab1) }

var fileDescriptor_tx_ce29446be3a85ab1 = []byte{
	// 265 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0xcd, 0x4a, 0xc3, 0x40,
	0x10, 0xc7, 0xd9, 0xcd, 0x47, 0x9b, 0xb1, 0x16, 0x19, 0x3c, 0xcc, 0x41, 0x30, 0xe4, 0x14, 0x2f,
	0x3d, 0xe8, 0x23, 0x18, 0x90, 0x52, 0x95, 0xb0, 0x2d, 0xbd, 0xaf, 0x75, 0x89, 0x0b, 0x31, 0x29,
	0x4d, 0x84, 0xf8, 0xc8, 0xbe, 0x85, 0xcc, 0x84, 0x36, 0xa7, 0xfd, 0x7f, 0xec, 0xfe, 0x86, 0x59,
	0x98, 0xf7, 0xc3, 0xea, 0x78, 0x6a, 0xfb, 0x16, 0x63, 0x39, 0xba, 0xec, 0x15, 0xf4, 0x6e, 0xc0,
	0x25, 0xe8, 0x75, 0x41, 0x2a, 0x55, 0x79, 0x62, 0xf4, 0xba, 0xc0, 0x0c, 0xc2, 0xc2, 0xf6, 0x96,
	0x74, 0xaa, 0xf2, 0xab, 0xc7, 0xe5, 0xf8, 0xa6, 0x5b, 0xf5, 0x03, 0xa7, 0x46, 0x3a, 0x44, 0x08,
	0x77, 0xfe, 0xdb, 0x51, 0x90, 0xaa, 0x3c, 0x30, 0xa2, 0xb3, 0x3f, 0x05, 0xf1, 0x78, 0x09, 0xef,
	0x20, 0x79, 0xfe, 0xb2, 0x4d, 0xe3, 0xea, 0x0b, 0x79, 0x0a, 0xf0, 0x16, 0xa2, 0xf7, 0xb6, 0x39,
	0x38, 0x99, 0x10, 0x9a, 0xd1, 0xf0, 0x1b, 0xe3, 0x0e, 0xfe, 0xe8, 0x5d, 0xd3, 0x0b, 0x77, 0x61,
	0xa6, 0x00, 0x09, 0x66, 0xa5, 0xfd, 0xad, 0x5b, 0xfb, 0x49, 0xa1, 0x74, 0x67, 0xcb, 0xb4, 0xbd,
	0xad, 0x7f, 0x1c, 0x45, 0x23, 0x4d, 0x0c, 0xde, 0x40, 0xf0, 0xd6, 0x55, 0x14, 0xcb, 0x6c, 0x96,
	0x4c, 0xd8, 0xbb, 0x53, 0xe7, 0xdb, 0x86, 0x66, 0xa9, 0xca, 0x23, 0x73, 0xb6, 0x78, 0x0f, 0xc1,
	0xd6, 0x57, 0x34, 0x97, 0x7d, 0xaf, 0xa7, 0x7d, 0xb7, 0xbe, 0x32, 0xdc, 0x30, 0xec, 0xc5, 0x76,
	0x94, 0xc8, 0x00, 0x96, 0xd9, 0x03, 0x44, 0xd2, 0xf3, 0xe7, 0x95, 0x1b, 0x59, 0x71, 0x61, 0x74,
	0xb9, 0xe1, 0xab, 0xcc, 0xd2, 0x12, 0xb0, 0xfc, 0x18, 0x3f, 0xfb, 0xe9, 0x3f, 0x00, 0x00, 0xff,
	0xff, 0xed, 0x2f, 0xeb, 0x62, 0x7f, 0x01, 0x00, 0x00,
}
