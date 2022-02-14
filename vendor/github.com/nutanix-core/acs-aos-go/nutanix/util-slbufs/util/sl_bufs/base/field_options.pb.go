// Code generated by protoc-gen-go. DO NOT EDIT.
// source: util/sl_bufs/base/field_options.proto

package base

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
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

var E_ExtendedType = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         80001,
	Name:          "nutanix.extended_type",
	Tag:           "bytes,80001,opt,name=extended_type",
	Filename:      "util/sl_bufs/base/field_options.proto",
}

func init() {
	proto.RegisterExtension(E_ExtendedType)
}

func init() {
	proto.RegisterFile("util/sl_bufs/base/field_options.proto", fileDescriptor_ff7d93c7da386766)
}

var fileDescriptor_ff7d93c7da386766 = []byte{
	// 196 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2d, 0x2d, 0xc9, 0xcc,
	0xd1, 0x2f, 0xce, 0x89, 0x4f, 0x2a, 0x4d, 0x2b, 0xd6, 0x4f, 0x4a, 0x2c, 0x4e, 0xd5, 0x4f, 0xcb,
	0x4c, 0xcd, 0x49, 0x89, 0xcf, 0x2f, 0x28, 0xc9, 0xcc, 0xcf, 0x2b, 0xd6, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x17, 0x62, 0xcf, 0x2b, 0x2d, 0x49, 0xcc, 0xcb, 0xac, 0x90, 0x52, 0x48, 0xcf, 0xcf, 0x4f,
	0xcf, 0x49, 0xd5, 0x07, 0x0b, 0x27, 0x95, 0xa6, 0xe9, 0xa7, 0xa4, 0x16, 0x27, 0x17, 0x65, 0x16,
	0x94, 0xe4, 0x17, 0x41, 0x94, 0x5a, 0xb9, 0x70, 0xf1, 0xa6, 0x56, 0x94, 0xa4, 0xe6, 0xa5, 0xa4,
	0xa6, 0xc4, 0x97, 0x54, 0x16, 0xa4, 0x0a, 0xc9, 0xea, 0x41, 0xf4, 0xe8, 0xc1, 0xf4, 0xe8, 0xb9,
	0x81, 0x6c, 0xf0, 0x87, 0x58, 0x20, 0xd1, 0xf8, 0x91, 0x45, 0x81, 0x51, 0x83, 0x33, 0x88, 0x07,
	0xa6, 0x2b, 0xa4, 0xb2, 0x20, 0xd5, 0x29, 0x8e, 0x8b, 0x3b, 0x39, 0x3f, 0x57, 0x0f, 0x6a, 0xad,
	0x93, 0x20, 0xb2, 0x96, 0x00, 0x90, 0x39, 0x51, 0x36, 0xe9, 0x99, 0x25, 0x19, 0xa5, 0x49, 0x7a,
	0xa9, 0x79, 0xe9, 0x30, 0x65, 0x7a, 0xc9, 0xf9, 0xb9, 0xfa, 0x50, 0xb6, 0x3e, 0xc8, 0x63, 0xba,
	0xc5, 0x39, 0x60, 0x7f, 0x61, 0x78, 0x12, 0x10, 0x00, 0x00, 0xff, 0xff, 0x67, 0x8a, 0x40, 0x8d,
	0xf8, 0x00, 0x00, 0x00,
}
