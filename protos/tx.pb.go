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
	Data                 *TxData  `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
	Time                 int64    `protobuf:"varint,2,opt,name=Time,proto3" json:"Time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Tx) Reset()         { *m = Tx{} }
func (m *Tx) String() string { return proto.CompactTextString(m) }
func (*Tx) ProtoMessage()    {}
func (*Tx) Descriptor() ([]byte, []int) {
	return fileDescriptor_tx_6b15f031d5d487d1, []int{0}
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
	AccountNonce         uint64   `protobuf:"varint,2,opt,name=AccountNonce,proto3" json:"AccountNonce,omitempty"`
	Recipient            []byte   `protobuf:"bytes,3,opt,name=Recipient,proto3" json:"Recipient,omitempty"`
	Payload              []byte   `protobuf:"bytes,4,opt,name=Payload,proto3" json:"Payload,omitempty"`
	Version              int32    `protobuf:"varint,5,opt,name=Version,proto3" json:"Version,omitempty"`
	Sig                  *TxSig   `protobuf:"bytes,6,opt,name=Sig,proto3" json:"Sig,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TxData) Reset()         { *m = TxData{} }
func (m *TxData) String() string { return proto.CompactTextString(m) }
func (*TxData) ProtoMessage()    {}
func (*TxData) Descriptor() ([]byte, []int) {
	return fileDescriptor_tx_6b15f031d5d487d1, []int{1}
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

func (m *TxData) GetAccountNonce() uint64 {
	if m != nil {
		return m.AccountNonce
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
	return fileDescriptor_tx_6b15f031d5d487d1, []int{2}
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

func init() { proto.RegisterFile("tx.proto", fileDescriptor_tx_6b15f031d5d487d1) }

var fileDescriptor_tx_6b15f031d5d487d1 = []byte{
	// 238 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0xcf, 0x4b, 0x3b, 0x31,
	0x10, 0xc5, 0x49, 0xf6, 0xc7, 0xf7, 0xdb, 0x71, 0x2d, 0x32, 0xa7, 0x1c, 0x04, 0x97, 0x9c, 0xd6,
	0x4b, 0x0f, 0x7a, 0xf5, 0x22, 0xf6, 0x22, 0x05, 0x59, 0xa6, 0xc5, 0x7b, 0x5c, 0x43, 0x0d, 0xd4,
	0xa4, 0x74, 0x23, 0xac, 0x7f, 0x9c, 0xff, 0x9b, 0x64, 0x42, 0x59, 0x3c, 0x65, 0xde, 0xe7, 0x4d,
	0x1e, 0x79, 0x81, 0xff, 0x71, 0x5a, 0x1d, 0x4f, 0x21, 0x06, 0xac, 0xf9, 0x18, 0xf5, 0x03, 0xc8,
	0xdd, 0x84, 0x1a, 0xca, 0xb5, 0x89, 0x46, 0x89, 0x56, 0x74, 0x17, 0x77, 0xcb, 0xbc, 0x33, 0xae,
	0xe2, 0x94, 0x28, 0xb1, 0x87, 0x08, 0xe5, 0xce, 0x7d, 0x5a, 0x25, 0x5b, 0xd1, 0x15, 0xc4, 0xb3,
	0xfe, 0x11, 0x50, 0xe7, 0x25, 0xbc, 0x86, 0xc5, 0xd3, 0x87, 0xf1, 0xde, 0x1e, 0x9e, 0xd7, 0x9c,
	0xb3, 0xa0, 0x19, 0xa0, 0x86, 0xe6, 0x71, 0x18, 0xc2, 0x97, 0x8f, 0x2f, 0xc1, 0x0f, 0x39, 0xa4,
	0xa4, 0x3f, 0x2c, 0x25, 0x90, 0x1d, 0xdc, 0xd1, 0x59, 0x1f, 0x55, 0xd1, 0x8a, 0xae, 0xa1, 0x19,
	0xa0, 0x82, 0x7f, 0xbd, 0xf9, 0x3e, 0x04, 0xf3, 0xae, 0x4a, 0xf6, 0xce, 0x32, 0x39, 0xaf, 0xf6,
	0x34, 0xba, 0xe0, 0x55, 0xd5, 0x8a, 0xae, 0xa2, 0xb3, 0xc4, 0x1b, 0x28, 0xb6, 0x6e, 0xaf, 0x6a,
	0x6e, 0x75, 0x39, 0xb7, 0xda, 0xba, 0x3d, 0x25, 0x47, 0xdf, 0x42, 0xc5, 0x0a, 0x97, 0x20, 0xfb,
	0x0d, 0x3f, 0xbb, 0x21, 0xd9, 0x6f, 0xf0, 0x2a, 0xdf, 0x94, 0x0c, 0xd2, 0xf8, 0x96, 0x3f, 0xec,
	0xfe, 0x37, 0x00, 0x00, 0xff, 0xff, 0xd5, 0x28, 0x09, 0xdd, 0x43, 0x01, 0x00, 0x00,
}
