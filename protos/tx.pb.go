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
	return fileDescriptor_tx_5f0c88a3f5aa22e5, []int{0}
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
	return fileDescriptor_tx_5f0c88a3f5aa22e5, []int{1}
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
	return fileDescriptor_tx_5f0c88a3f5aa22e5, []int{2}
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

func init() { proto.RegisterFile("tx.proto", fileDescriptor_tx_5f0c88a3f5aa22e5) }

var fileDescriptor_tx_5f0c88a3f5aa22e5 = []byte{
	// 268 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x50, 0x4f, 0x4b, 0xfb, 0x40,
	0x10, 0x65, 0x37, 0x7f, 0xda, 0xce, 0xaf, 0xbf, 0x22, 0x83, 0x87, 0x3d, 0x08, 0x86, 0x9c, 0x72,
	0xea, 0x41, 0xcf, 0x1e, 0xc4, 0x5c, 0x4a, 0x55, 0xc2, 0xb4, 0xf4, 0xbe, 0xd6, 0x25, 0x2e, 0xc4,
	0x6c, 0x69, 0x22, 0xc4, 0x4f, 0xea, 0xd7, 0x91, 0x99, 0x50, 0x73, 0xda, 0xf7, 0xde, 0xec, 0xbc,
	0xe1, 0x3d, 0x98, 0xf7, 0xc3, 0xfa, 0x74, 0x0e, 0x7d, 0xc0, 0x54, 0x9e, 0x2e, 0x7f, 0x06, 0xbd,
	0x1f, 0x70, 0x05, 0x7a, 0x53, 0x1a, 0x95, 0xa9, 0x62, 0x41, 0x7a, 0x53, 0x62, 0x0e, 0x71, 0x69,
	0x7b, 0x6b, 0x74, 0xa6, 0x8a, 0x7f, 0x77, 0xab, 0x71, 0xa7, 0x5b, 0xf7, 0x03, 0xab, 0x24, 0x33,
	0x44, 0x88, 0xf7, 0xfe, 0xd3, 0x99, 0x28, 0x53, 0x45, 0x44, 0x82, 0xf3, 0x1f, 0x05, 0xe9, 0xf8,
	0x09, 0x6f, 0x60, 0xf1, 0xf4, 0x61, 0xdb, 0xd6, 0x35, 0x7f, 0xce, 0x93, 0x80, 0xd7, 0x90, 0xbc,
	0x86, 0xf6, 0xe8, 0xe4, 0x42, 0x4c, 0x23, 0xe1, 0x1d, 0x72, 0x47, 0x7f, 0xf2, 0xae, 0xed, 0xc5,
	0x77, 0x49, 0x93, 0x80, 0x06, 0x66, 0x95, 0xfd, 0x6e, 0x82, 0x7d, 0x37, 0xb1, 0xcc, 0x2e, 0x94,
	0xdd, 0x0e, 0xb6, 0xf9, 0x72, 0x26, 0x19, 0xdd, 0x84, 0xe0, 0x15, 0x44, 0x2f, 0x5d, 0x6d, 0x52,
	0xb9, 0xcd, 0x90, 0x1d, 0x0e, 0xee, 0xdc, 0xf9, 0xd0, 0x9a, 0x59, 0xa6, 0x8a, 0x84, 0x2e, 0x14,
	0x6f, 0x21, 0xda, 0xf9, 0xda, 0xcc, 0x25, 0xef, 0xff, 0x29, 0xef, 0xce, 0xd7, 0xc4, 0x93, 0xfc,
	0x01, 0x12, 0x61, 0x5c, 0x55, 0xb5, 0x95, 0x40, 0x4b, 0xd2, 0xd5, 0x96, 0xaf, 0xf0, 0xa6, 0x16,
	0x81, 0x21, 0x17, 0xf3, 0xd8, 0xd4, 0x41, 0x02, 0x24, 0x24, 0xf8, 0x6d, 0xac, 0xfb, 0xfe, 0x37,
	0x00, 0x00, 0xff, 0xff, 0xd7, 0xe0, 0x3b, 0xf8, 0x81, 0x01, 0x00, 0x00,
}
