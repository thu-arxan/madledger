// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/api/context.proto

package serviceconfig // import "google.golang.org/genproto/googleapis/api/serviceconfig"

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

// `Context` defines which contexts an API requests.
//
// Example:
//
//     context:
//       rules:
//       - selector: "*"
//         requested:
//         - google.rpc.context.ProjectContext
//         - google.rpc.context.OriginContext
//
// The above specifies that all methods in the API request
// `google.rpc.context.ProjectContext` and
// `google.rpc.context.OriginContext`.
//
// Available context types are defined in package
// `google.rpc.context`.
type Context struct {
	// A list of RPC context rules that apply to individual API methods.
	//
	// **NOTE:** All service configuration rules follow "last one wins" order.
	Rules                []*ContextRule `protobuf:"bytes,1,rep,name=rules,proto3" json:"rules,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Context) Reset()         { *m = Context{} }
func (m *Context) String() string { return proto.CompactTextString(m) }
func (*Context) ProtoMessage()    {}
func (*Context) Descriptor() ([]byte, []int) {
	return fileDescriptor_context_aced50768e79506a, []int{0}
}
func (m *Context) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Context.Unmarshal(m, b)
}
func (m *Context) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Context.Marshal(b, m, deterministic)
}
func (dst *Context) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Context.Merge(dst, src)
}
func (m *Context) XXX_Size() int {
	return xxx_messageInfo_Context.Size(m)
}
func (m *Context) XXX_DiscardUnknown() {
	xxx_messageInfo_Context.DiscardUnknown(m)
}

var xxx_messageInfo_Context proto.InternalMessageInfo

func (m *Context) GetRules() []*ContextRule {
	if m != nil {
		return m.Rules
	}
	return nil
}

// A context rule provides information about the context for an individual API
// element.
type ContextRule struct {
	// Selects the methods to which this rule applies.
	//
	// Refer to [selector][google.api.DocumentationRule.selector] for syntax details.
	Selector string `protobuf:"bytes,1,opt,name=selector,proto3" json:"selector,omitempty"`
	// A list of full type names of requested contexts.
	Requested []string `protobuf:"bytes,2,rep,name=requested,proto3" json:"requested,omitempty"`
	// A list of full type names of provided contexts.
	Provided             []string `protobuf:"bytes,3,rep,name=provided,proto3" json:"provided,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ContextRule) Reset()         { *m = ContextRule{} }
func (m *ContextRule) String() string { return proto.CompactTextString(m) }
func (*ContextRule) ProtoMessage()    {}
func (*ContextRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_context_aced50768e79506a, []int{1}
}
func (m *ContextRule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ContextRule.Unmarshal(m, b)
}
func (m *ContextRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ContextRule.Marshal(b, m, deterministic)
}
func (dst *ContextRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ContextRule.Merge(dst, src)
}
func (m *ContextRule) XXX_Size() int {
	return xxx_messageInfo_ContextRule.Size(m)
}
func (m *ContextRule) XXX_DiscardUnknown() {
	xxx_messageInfo_ContextRule.DiscardUnknown(m)
}

var xxx_messageInfo_ContextRule proto.InternalMessageInfo

func (m *ContextRule) GetSelector() string {
	if m != nil {
		return m.Selector
	}
	return ""
}

func (m *ContextRule) GetRequested() []string {
	if m != nil {
		return m.Requested
	}
	return nil
}

func (m *ContextRule) GetProvided() []string {
	if m != nil {
		return m.Provided
	}
	return nil
}

func init() {
	proto.RegisterType((*Context)(nil), "google.api.Context")
	proto.RegisterType((*ContextRule)(nil), "google.api.ContextRule")
}

func init() { proto.RegisterFile("google/api/context.proto", fileDescriptor_context_aced50768e79506a) }

var fileDescriptor_context_aced50768e79506a = []byte{
	// 231 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0x4f, 0x4b, 0xc4, 0x30,
	0x14, 0xc4, 0xe9, 0xd6, 0x7f, 0x7d, 0x2b, 0x1e, 0x7a, 0x31, 0x88, 0x87, 0xb2, 0xa7, 0x5e, 0x4c,
	0x41, 0x2f, 0x82, 0x27, 0x57, 0x44, 0xbc, 0x95, 0x1e, 0xbd, 0xc5, 0xf4, 0x19, 0x02, 0x31, 0x2f,
	0x26, 0xe9, 0xe2, 0xe7, 0xf1, 0x93, 0xca, 0x26, 0x65, 0xff, 0x1c, 0x67, 0x7e, 0x33, 0x24, 0xf3,
	0x80, 0x29, 0x22, 0x65, 0xb0, 0x13, 0x4e, 0x77, 0x92, 0x6c, 0xc4, 0xdf, 0xc8, 0x9d, 0xa7, 0x48,
	0x35, 0x64, 0xc2, 0x85, 0xd3, 0xab, 0x47, 0x38, 0x7f, 0xc9, 0xb0, 0xbe, 0x83, 0x53, 0x3f, 0x19,
	0x0c, 0xac, 0x68, 0xca, 0x76, 0x79, 0x7f, 0xcd, 0xf7, 0x31, 0x3e, 0x67, 0x86, 0xc9, 0xe0, 0x90,
	0x53, 0x2b, 0x09, 0xcb, 0x03, 0xb7, 0xbe, 0x81, 0x8b, 0x80, 0x06, 0x65, 0x24, 0xcf, 0x8a, 0xa6,
	0x68, 0xab, 0x61, 0xa7, 0xeb, 0x5b, 0xa8, 0x3c, 0xfe, 0x4c, 0x18, 0x22, 0x8e, 0x6c, 0xd1, 0x94,
	0x6d, 0x35, 0xec, 0x8d, 0x6d, 0xd3, 0x79, 0xda, 0xe8, 0x11, 0x47, 0x56, 0x26, 0xb8, 0xd3, 0x6b,
	0x0b, 0x57, 0x92, 0xbe, 0x0f, 0x7e, 0xb2, 0xbe, 0x9c, 0x1f, 0xed, 0xb7, 0x53, 0xfa, 0xe2, 0xe3,
	0x75, 0x66, 0x8a, 0x8c, 0xb0, 0x8a, 0x93, 0x57, 0x9d, 0x42, 0x9b, 0x86, 0x76, 0x19, 0x09, 0xa7,
	0x43, 0xba, 0x42, 0x40, 0xbf, 0xd1, 0x12, 0x25, 0xd9, 0x2f, 0xad, 0x9e, 0x8e, 0xd4, 0xdf, 0xe2,
	0xe4, 0xed, 0xb9, 0x7f, 0xff, 0x3c, 0x4b, 0xc5, 0x87, 0xff, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb5,
	0x18, 0x98, 0x7a, 0x3d, 0x01, 0x00, 0x00,
}
