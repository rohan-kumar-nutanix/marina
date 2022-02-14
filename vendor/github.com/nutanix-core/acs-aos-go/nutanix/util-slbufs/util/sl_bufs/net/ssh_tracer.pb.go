// Code generated by protoc-gen-go. DO NOT EDIT.
// source: util/sl_bufs/net/ssh_tracer.proto

package net

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Details about a SSH request.
type SSHRequest struct {
	// Seconds to wait before timing out.
	TimeoutSecs *int32 `protobuf:"varint,1,opt,name=timeout_secs,json=timeoutSecs" json:"timeout_secs,omitempty"`
	// Required.
	// Command to be run on host as username.
	ObfuscatedCommand    *string  `protobuf:"bytes,2,opt,name=obfuscated_command,json=obfuscatedCommand" json:"obfuscated_command,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SSHRequest) Reset()         { *m = SSHRequest{} }
func (m *SSHRequest) String() string { return proto.CompactTextString(m) }
func (*SSHRequest) ProtoMessage()    {}
func (*SSHRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_f8e1fccdb6ca6e27, []int{0}
}

func (m *SSHRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SSHRequest.Unmarshal(m, b)
}
func (m *SSHRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SSHRequest.Marshal(b, m, deterministic)
}
func (m *SSHRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SSHRequest.Merge(m, src)
}
func (m *SSHRequest) XXX_Size() int {
	return xxx_messageInfo_SSHRequest.Size(m)
}
func (m *SSHRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SSHRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SSHRequest proto.InternalMessageInfo

func (m *SSHRequest) GetTimeoutSecs() int32 {
	if m != nil && m.TimeoutSecs != nil {
		return *m.TimeoutSecs
	}
	return 0
}

func (m *SSHRequest) GetObfuscatedCommand() string {
	if m != nil && m.ObfuscatedCommand != nil {
		return *m.ObfuscatedCommand
	}
	return ""
}

// Details of the response from a SSH request.
type SSHResponse struct {
	// Exit code after obfuscated_command is run on host.
	ReturnValue *int32 `protobuf:"varint,1,opt,name=return_value,json=returnValue" json:"return_value,omitempty"`
	// stdout output.
	Stdout *string `protobuf:"bytes,2,opt,name=stdout" json:"stdout,omitempty"`
	// stderr output.
	Stderr               *string  `protobuf:"bytes,3,opt,name=stderr" json:"stderr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SSHResponse) Reset()         { *m = SSHResponse{} }
func (m *SSHResponse) String() string { return proto.CompactTextString(m) }
func (*SSHResponse) ProtoMessage()    {}
func (*SSHResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_f8e1fccdb6ca6e27, []int{1}
}

func (m *SSHResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SSHResponse.Unmarshal(m, b)
}
func (m *SSHResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SSHResponse.Marshal(b, m, deterministic)
}
func (m *SSHResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SSHResponse.Merge(m, src)
}
func (m *SSHResponse) XXX_Size() int {
	return xxx_messageInfo_SSHResponse.Size(m)
}
func (m *SSHResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SSHResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SSHResponse proto.InternalMessageInfo

func (m *SSHResponse) GetReturnValue() int32 {
	if m != nil && m.ReturnValue != nil {
		return *m.ReturnValue
	}
	return 0
}

func (m *SSHResponse) GetStdout() string {
	if m != nil && m.Stdout != nil {
		return *m.Stdout
	}
	return ""
}

func (m *SSHResponse) GetStderr() string {
	if m != nil && m.Stderr != nil {
		return *m.Stderr
	}
	return ""
}

// Metadata to be be logged along with SSHRequest and SSHResponse protobufs when
// writing binary log records.
type SSHMetadata struct {
	// Required.
	// Unique id assigned to a SSH request-response pair.
	SshId *int64 `protobuf:"varint,1,opt,name=ssh_id,json=sshId" json:"ssh_id,omitempty"`
	// DNS name or IP address of the node to SSH to.
	Host *string `protobuf:"bytes,2,opt,name=host" json:"host,omitempty"`
	// SSH username.
	Username *string `protobuf:"bytes,3,opt,name=username" json:"username,omitempty"`
	// Local address to bind to.
	BindAddress          *string  `protobuf:"bytes,4,opt,name=bind_address,json=bindAddress" json:"bind_address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SSHMetadata) Reset()         { *m = SSHMetadata{} }
func (m *SSHMetadata) String() string { return proto.CompactTextString(m) }
func (*SSHMetadata) ProtoMessage()    {}
func (*SSHMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_f8e1fccdb6ca6e27, []int{2}
}

func (m *SSHMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SSHMetadata.Unmarshal(m, b)
}
func (m *SSHMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SSHMetadata.Marshal(b, m, deterministic)
}
func (m *SSHMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SSHMetadata.Merge(m, src)
}
func (m *SSHMetadata) XXX_Size() int {
	return xxx_messageInfo_SSHMetadata.Size(m)
}
func (m *SSHMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_SSHMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_SSHMetadata proto.InternalMessageInfo

func (m *SSHMetadata) GetSshId() int64 {
	if m != nil && m.SshId != nil {
		return *m.SshId
	}
	return 0
}

func (m *SSHMetadata) GetHost() string {
	if m != nil && m.Host != nil {
		return *m.Host
	}
	return ""
}

func (m *SSHMetadata) GetUsername() string {
	if m != nil && m.Username != nil {
		return *m.Username
	}
	return ""
}

func (m *SSHMetadata) GetBindAddress() string {
	if m != nil && m.BindAddress != nil {
		return *m.BindAddress
	}
	return ""
}

func init() {
	proto.RegisterType((*SSHRequest)(nil), "nutanix.net.SSHRequest")
	proto.RegisterType((*SSHResponse)(nil), "nutanix.net.SSHResponse")
	proto.RegisterType((*SSHMetadata)(nil), "nutanix.net.SSHMetadata")
}

func init() { proto.RegisterFile("util/sl_bufs/net/ssh_tracer.proto", fileDescriptor_f8e1fccdb6ca6e27) }

var fileDescriptor_f8e1fccdb6ca6e27 = []byte{
	// 293 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0xbf, 0x4e, 0x33, 0x31,
	0x10, 0xc4, 0x95, 0x2f, 0x7f, 0xf4, 0xe1, 0xa3, 0xc1, 0x12, 0xe8, 0x44, 0x95, 0xa4, 0x4a, 0x93,
	0xbb, 0x07, 0x40, 0x14, 0x40, 0x03, 0x05, 0xcd, 0x45, 0xa2, 0xa0, 0xc0, 0xf8, 0xec, 0x4d, 0xce,
	0xd2, 0x9d, 0x1d, 0xbc, 0x6b, 0xe0, 0xf1, 0x91, 0x1d, 0x27, 0x48, 0x74, 0xbb, 0xbf, 0x91, 0x66,
	0xec, 0x1d, 0xb6, 0x08, 0x64, 0xfa, 0x1a, 0x7b, 0xd1, 0x86, 0x2d, 0xd6, 0x16, 0xa8, 0x46, 0xec,
	0x04, 0x79, 0xa9, 0xc0, 0x57, 0x7b, 0xef, 0xc8, 0xf1, 0xc2, 0x06, 0x92, 0xd6, 0x7c, 0x57, 0x16,
	0x68, 0xf9, 0xc6, 0xd8, 0x66, 0xf3, 0xd8, 0xc0, 0x47, 0x00, 0x24, 0xbe, 0x60, 0xe7, 0x64, 0x06,
	0x70, 0x81, 0x04, 0x82, 0xc2, 0x72, 0x34, 0x1f, 0xad, 0xa6, 0x4d, 0x91, 0xd9, 0x06, 0x14, 0xf2,
	0x35, 0xe3, 0xae, 0xdd, 0x06, 0x54, 0x92, 0x40, 0x0b, 0xe5, 0x86, 0x41, 0x5a, 0x5d, 0xfe, 0x9b,
	0x8f, 0x56, 0x67, 0xcd, 0xc5, 0xaf, 0xf2, 0x70, 0x10, 0x96, 0xef, 0xac, 0x48, 0xfe, 0xb8, 0x77,
	0x16, 0x21, 0x06, 0x78, 0xa0, 0xe0, 0xad, 0xf8, 0x94, 0x7d, 0x80, 0x63, 0xc0, 0x81, 0xbd, 0x44,
	0xc4, 0xaf, 0xd8, 0x0c, 0x49, 0xbb, 0x40, 0xd9, 0x34, 0x6f, 0x99, 0x83, 0xf7, 0xe5, 0xf8, 0xc4,
	0xc1, 0xfb, 0xe5, 0x57, 0x4a, 0x78, 0x06, 0x92, 0x5a, 0x92, 0xe4, 0x97, 0x6c, 0x16, 0x7f, 0x6c,
	0x74, 0xf2, 0x1e, 0x37, 0x53, 0xc4, 0xee, 0x49, 0x73, 0xce, 0x26, 0x9d, 0xc3, 0xa3, 0x67, 0x9a,
	0xf9, 0x35, 0xfb, 0x1f, 0x10, 0xbc, 0x95, 0x03, 0x64, 0xcf, 0xd3, 0x1e, 0x1f, 0xda, 0x1a, 0xab,
	0x85, 0xd4, 0xda, 0x03, 0x62, 0x39, 0x49, 0x7a, 0x11, 0xd9, 0xdd, 0x01, 0xdd, 0xdf, 0xbe, 0xde,
	0xec, 0x0c, 0x75, 0xa1, 0xad, 0xc0, 0xee, 0xaa, 0xe3, 0x51, 0x95, 0x1b, 0xea, 0x3c, 0xd7, 0xb1,
	0x8b, 0x35, 0xf6, 0xa9, 0x8a, 0xbf, 0xbd, 0xfc, 0x04, 0x00, 0x00, 0xff, 0xff, 0xcc, 0x21, 0xd1,
	0xc5, 0xaa, 0x01, 0x00, 0x00,
}