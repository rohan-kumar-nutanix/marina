//
// Copyright (c) 2022 Nutanix Inc. All rights reserved.
//
// Author: rajesh.battala@nutanix.com

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.11.4
// source: protos/marina/marina_interface.proto

package marina

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CatalogItemGetArg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of catalog item types to use as a filter.
	CatalogItemTypeList []CatalogItemInfo_CatalogItemType `protobuf:"varint,1,rep,name=catalog_item_type_list,json=catalogItemTypeList,enum=protos.marina.CatalogItemInfo_CatalogItemType" json:"catalog_item_type_list,omitempty"`
	// List of catalog item query details.
	CatalogItemIdList []*CatalogItemId `protobuf:"bytes,2,rep,name=catalog_item_id_list,json=catalogItemIdList" json:"catalog_item_id_list,omitempty"`
	// Get only latest version of the catalog item
	Latest *bool `protobuf:"varint,3,opt,name=latest,def=0" json:"latest,omitempty"`
}

// Default values for CatalogItemGetArg fields.
const (
	Default_CatalogItemGetArg_Latest = bool(false)
)

func (x *CatalogItemGetArg) Reset() {
	*x = CatalogItemGetArg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_marina_marina_interface_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CatalogItemGetArg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CatalogItemGetArg) ProtoMessage() {}

func (x *CatalogItemGetArg) ProtoReflect() protoreflect.Message {
	mi := &file_protos_marina_marina_interface_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CatalogItemGetArg.ProtoReflect.Descriptor instead.
func (*CatalogItemGetArg) Descriptor() ([]byte, []int) {
	return file_protos_marina_marina_interface_proto_rawDescGZIP(), []int{0}
}

func (x *CatalogItemGetArg) GetCatalogItemTypeList() []CatalogItemInfo_CatalogItemType {
	if x != nil {
		return x.CatalogItemTypeList
	}
	return nil
}

func (x *CatalogItemGetArg) GetCatalogItemIdList() []*CatalogItemId {
	if x != nil {
		return x.CatalogItemIdList
	}
	return nil
}

func (x *CatalogItemGetArg) GetLatest() bool {
	if x != nil && x.Latest != nil {
		return *x.Latest
	}
	return Default_CatalogItemGetArg_Latest
}

type CatalogItemGetRet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of catalog item objects.
	CatalogItemList []*CatalogItemInfo `protobuf:"bytes,1,rep,name=catalog_item_list,json=catalogItemList" json:"catalog_item_list,omitempty"`
}

func (x *CatalogItemGetRet) Reset() {
	*x = CatalogItemGetRet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_marina_marina_interface_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CatalogItemGetRet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CatalogItemGetRet) ProtoMessage() {}

func (x *CatalogItemGetRet) ProtoReflect() protoreflect.Message {
	mi := &file_protos_marina_marina_interface_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CatalogItemGetRet.ProtoReflect.Descriptor instead.
func (*CatalogItemGetRet) Descriptor() ([]byte, []int) {
	return file_protos_marina_marina_interface_proto_rawDescGZIP(), []int{1}
}

func (x *CatalogItemGetRet) GetCatalogItemList() []*CatalogItemInfo {
	if x != nil {
		return x.CatalogItemList
	}
	return nil
}

var File_protos_marina_marina_interface_proto protoreflect.FileDescriptor

var file_protos_marina_marina_interface_proto_rawDesc = []byte{
	0x0a, 0x24, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61, 0x2f,
	0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x66, 0x61, 0x63, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x6d,
	0x61, 0x72, 0x69, 0x6e, 0x61, 0x1a, 0x20, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6d, 0x61,
	0x72, 0x69, 0x6e, 0x61, 0x2f, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe6, 0x01, 0x0a, 0x11, 0x43, 0x61, 0x74, 0x61,
	0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x47, 0x65, 0x74, 0x41, 0x72, 0x67, 0x12, 0x63, 0x0a,
	0x16, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x2e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61, 0x2e, 0x43, 0x61,
	0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x43, 0x61,
	0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x52, 0x13, 0x63,
	0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x4d, 0x0a, 0x14, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x5f, 0x69, 0x74,
	0x65, 0x6d, 0x5f, 0x69, 0x64, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61,
	0x2e, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x52, 0x11,
	0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x4c, 0x69, 0x73,
	0x74, 0x12, 0x1d, 0x0a, 0x06, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x08, 0x3a, 0x05, 0x66, 0x61, 0x6c, 0x73, 0x65, 0x52, 0x06, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74,
	0x22, 0x5f, 0x0a, 0x11, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x47,
	0x65, 0x74, 0x52, 0x65, 0x74, 0x12, 0x4a, 0x0a, 0x11, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67,
	0x5f, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61,
	0x2e, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x0f, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x4c, 0x69, 0x73,
	0x74, 0x32, 0x5e, 0x0a, 0x06, 0x4d, 0x61, 0x72, 0x69, 0x6e, 0x61, 0x12, 0x54, 0x0a, 0x0e, 0x43,
	0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x47, 0x65, 0x74, 0x12, 0x20, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61, 0x2e, 0x43, 0x61,
	0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x47, 0x65, 0x74, 0x41, 0x72, 0x67, 0x1a,
	0x20, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61, 0x2e,
	0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x47, 0x65, 0x74, 0x52, 0x65,
	0x74, 0x42, 0x0f, 0x5a, 0x0d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6d, 0x61, 0x72, 0x69,
	0x6e, 0x61,
}

var (
	file_protos_marina_marina_interface_proto_rawDescOnce sync.Once
	file_protos_marina_marina_interface_proto_rawDescData = file_protos_marina_marina_interface_proto_rawDesc
)

func file_protos_marina_marina_interface_proto_rawDescGZIP() []byte {
	file_protos_marina_marina_interface_proto_rawDescOnce.Do(func() {
		file_protos_marina_marina_interface_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_marina_marina_interface_proto_rawDescData)
	})
	return file_protos_marina_marina_interface_proto_rawDescData
}

var file_protos_marina_marina_interface_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_protos_marina_marina_interface_proto_goTypes = []interface{}{
	(*CatalogItemGetArg)(nil),            // 0: protos.marina.CatalogItemGetArg
	(*CatalogItemGetRet)(nil),            // 1: protos.marina.CatalogItemGetRet
	(CatalogItemInfo_CatalogItemType)(0), // 2: protos.marina.CatalogItemInfo.CatalogItemType
	(*CatalogItemId)(nil),                // 3: protos.marina.CatalogItemId
	(*CatalogItemInfo)(nil),              // 4: protos.marina.CatalogItemInfo
}
var file_protos_marina_marina_interface_proto_depIdxs = []int32{
	2, // 0: protos.marina.CatalogItemGetArg.catalog_item_type_list:type_name -> protos.marina.CatalogItemInfo.CatalogItemType
	3, // 1: protos.marina.CatalogItemGetArg.catalog_item_id_list:type_name -> protos.marina.CatalogItemId
	4, // 2: protos.marina.CatalogItemGetRet.catalog_item_list:type_name -> protos.marina.CatalogItemInfo
	0, // 3: protos.marina.Marina.CatalogItemGet:input_type -> protos.marina.CatalogItemGetArg
	1, // 4: protos.marina.Marina.CatalogItemGet:output_type -> protos.marina.CatalogItemGetRet
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_protos_marina_marina_interface_proto_init() }
func file_protos_marina_marina_interface_proto_init() {
	if File_protos_marina_marina_interface_proto != nil {
		return
	}
	file_protos_marina_marina_types_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_protos_marina_marina_interface_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CatalogItemGetArg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_protos_marina_marina_interface_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CatalogItemGetRet); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protos_marina_marina_interface_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_marina_marina_interface_proto_goTypes,
		DependencyIndexes: file_protos_marina_marina_interface_proto_depIdxs,
		MessageInfos:      file_protos_marina_marina_interface_proto_msgTypes,
	}.Build()
	File_protos_marina_marina_interface_proto = out.File
	file_protos_marina_marina_interface_proto_rawDesc = nil
	file_protos_marina_marina_interface_proto_goTypes = nil
	file_protos_marina_marina_interface_proto_depIdxs = nil
}
