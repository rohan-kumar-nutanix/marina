//
// Copyright (c) 2022 Nutanix Inc. All rights reserved.
//
// Author: rajesh.battala@nutanix.com

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.11.4
// source: protos/marina/marina_wal.proto

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

type PcTaskWalRecordCatalogItemData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Catalog item UUID.
	CatalogItemUuid []byte `protobuf:"bytes,1,opt,name=catalog_item_uuid,json=catalogItemUuid" json:"catalog_item_uuid,omitempty"`
	// Catalog item's logical timestamp
	Version *int64 `protobuf:"varint,2,opt,name=version,def=-1" json:"version,omitempty"`
	// List of UUIDs of files created in the file repository as the result of a
	// task.
	FileUuidList [][]byte `protobuf:"bytes,3,rep,name=file_uuid_list,json=fileUuidList" json:"file_uuid_list,omitempty"`
	// Catalog item UUID.
	GlobalCatalogItemUuid []byte `protobuf:"bytes,4,opt,name=global_catalog_item_uuid,json=globalCatalogItemUuid" json:"global_catalog_item_uuid,omitempty"`
	// Cluster UUIDs where the image RPC is fanned out. Used by the PC version
	// of the forward image task.
	ClusterUuidList [][]byte `protobuf:"bytes,5,rep,name=cluster_uuid_list,json=clusterUuidList" json:"cluster_uuid_list,omitempty"`
	// New Catalog item UUID.
	NewCatalogItemUuid []byte `protobuf:"bytes,6,opt,name=new_catalog_item_uuid,json=newCatalogItemUuid" json:"new_catalog_item_uuid,omitempty"`
}

// Default values for PcTaskWalRecordCatalogItemData fields.
const (
	Default_PcTaskWalRecordCatalogItemData_Version = int64(-1)
)

func (x *PcTaskWalRecordCatalogItemData) Reset() {
	*x = PcTaskWalRecordCatalogItemData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_marina_marina_wal_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PcTaskWalRecordCatalogItemData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PcTaskWalRecordCatalogItemData) ProtoMessage() {}

func (x *PcTaskWalRecordCatalogItemData) ProtoReflect() protoreflect.Message {
	mi := &file_protos_marina_marina_wal_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PcTaskWalRecordCatalogItemData.ProtoReflect.Descriptor instead.
func (*PcTaskWalRecordCatalogItemData) Descriptor() ([]byte, []int) {
	return file_protos_marina_marina_wal_proto_rawDescGZIP(), []int{0}
}

func (x *PcTaskWalRecordCatalogItemData) GetCatalogItemUuid() []byte {
	if x != nil {
		return x.CatalogItemUuid
	}
	return nil
}

func (x *PcTaskWalRecordCatalogItemData) GetVersion() int64 {
	if x != nil && x.Version != nil {
		return *x.Version
	}
	return Default_PcTaskWalRecordCatalogItemData_Version
}

func (x *PcTaskWalRecordCatalogItemData) GetFileUuidList() [][]byte {
	if x != nil {
		return x.FileUuidList
	}
	return nil
}

func (x *PcTaskWalRecordCatalogItemData) GetGlobalCatalogItemUuid() []byte {
	if x != nil {
		return x.GlobalCatalogItemUuid
	}
	return nil
}

func (x *PcTaskWalRecordCatalogItemData) GetClusterUuidList() [][]byte {
	if x != nil {
		return x.ClusterUuidList
	}
	return nil
}

func (x *PcTaskWalRecordCatalogItemData) GetNewCatalogItemUuid() []byte {
	if x != nil {
		return x.NewCatalogItemUuid
	}
	return nil
}

type PcTaskWalRecordData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// WAL data for catalog item related tasks.
	CatalogItem *PcTaskWalRecordCatalogItemData `protobuf:"bytes,1,opt,name=catalog_item,json=catalogItem" json:"catalog_item,omitempty"`
}

func (x *PcTaskWalRecordData) Reset() {
	*x = PcTaskWalRecordData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_marina_marina_wal_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PcTaskWalRecordData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PcTaskWalRecordData) ProtoMessage() {}

func (x *PcTaskWalRecordData) ProtoReflect() protoreflect.Message {
	mi := &file_protos_marina_marina_wal_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PcTaskWalRecordData.ProtoReflect.Descriptor instead.
func (*PcTaskWalRecordData) Descriptor() ([]byte, []int) {
	return file_protos_marina_marina_wal_proto_rawDescGZIP(), []int{1}
}

func (x *PcTaskWalRecordData) GetCatalogItem() *PcTaskWalRecordCatalogItemData {
	if x != nil {
		return x.CatalogItem
	}
	return nil
}

// Marina task WAL.
type PcTaskWalRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Additional data logged by operations.
	Data *PcTaskWalRecordData `protobuf:"bytes,2,opt,name=data" json:"data,omitempty"`
	// Whether the entry has been deleted. Used by the refactored python task
	// library.
	Deleted *bool `protobuf:"varint,3,opt,name=deleted" json:"deleted,omitempty"`
	// The proxy task WAL.
	ProxyTaskWal *ProxyTaskWal `protobuf:"bytes,4,opt,name=proxy_task_wal,json=proxyTaskWal" json:"proxy_task_wal,omitempty"`
}

func (x *PcTaskWalRecord) Reset() {
	*x = PcTaskWalRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_marina_marina_wal_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PcTaskWalRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PcTaskWalRecord) ProtoMessage() {}

func (x *PcTaskWalRecord) ProtoReflect() protoreflect.Message {
	mi := &file_protos_marina_marina_wal_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PcTaskWalRecord.ProtoReflect.Descriptor instead.
func (*PcTaskWalRecord) Descriptor() ([]byte, []int) {
	return file_protos_marina_marina_wal_proto_rawDescGZIP(), []int{2}
}

func (x *PcTaskWalRecord) GetData() *PcTaskWalRecordData {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *PcTaskWalRecord) GetDeleted() bool {
	if x != nil && x.Deleted != nil {
		return *x.Deleted
	}
	return false
}

func (x *PcTaskWalRecord) GetProxyTaskWal() *ProxyTaskWal {
	if x != nil {
		return x.ProxyTaskWal
	}
	return nil
}

// Proxy task WAL.
type ProxyTaskWal struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The UUID of the Catalog task.
	TaskUuid []byte `protobuf:"bytes,1,opt,name=task_uuid,json=taskUuid" json:"task_uuid,omitempty"`
	// The serialization token of the task.
	SerializationToken *string `protobuf:"bytes,2,opt,name=serialization_token,json=serializationToken" json:"serialization_token,omitempty"`
}

func (x *ProxyTaskWal) Reset() {
	*x = ProxyTaskWal{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_marina_marina_wal_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProxyTaskWal) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProxyTaskWal) ProtoMessage() {}

func (x *ProxyTaskWal) ProtoReflect() protoreflect.Message {
	mi := &file_protos_marina_marina_wal_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProxyTaskWal.ProtoReflect.Descriptor instead.
func (*ProxyTaskWal) Descriptor() ([]byte, []int) {
	return file_protos_marina_marina_wal_proto_rawDescGZIP(), []int{3}
}

func (x *ProxyTaskWal) GetTaskUuid() []byte {
	if x != nil {
		return x.TaskUuid
	}
	return nil
}

func (x *ProxyTaskWal) GetSerializationToken() string {
	if x != nil && x.SerializationToken != nil {
		return *x.SerializationToken
	}
	return ""
}

var File_protos_marina_marina_wal_proto protoreflect.FileDescriptor

var file_protos_marina_marina_wal_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61, 0x2f,
	0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61, 0x5f, 0x77, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61, 0x22,
	0xa8, 0x02, 0x0a, 0x1e, 0x50, 0x63, 0x54, 0x61, 0x73, 0x6b, 0x57, 0x61, 0x6c, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x44, 0x61,
	0x74, 0x61, 0x12, 0x2a, 0x0a, 0x11, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x5f, 0x69, 0x74,
	0x65, 0x6d, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f, 0x63,
	0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x55, 0x75, 0x69, 0x64, 0x12, 0x1c,
	0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x3a,
	0x02, 0x2d, 0x31, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x24, 0x0a, 0x0e,
	0x66, 0x69, 0x6c, 0x65, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0c, 0x52, 0x0c, 0x66, 0x69, 0x6c, 0x65, 0x55, 0x75, 0x69, 0x64, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x37, 0x0a, 0x18, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x5f, 0x63, 0x61, 0x74,
	0x61, 0x6c, 0x6f, 0x67, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x15, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x43, 0x61, 0x74, 0x61,
	0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x55, 0x75, 0x69, 0x64, 0x12, 0x2a, 0x0a, 0x11, 0x63,
	0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x5f, 0x6c, 0x69, 0x73, 0x74,
	0x18, 0x05, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x0f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x55,
	0x75, 0x69, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x31, 0x0a, 0x15, 0x6e, 0x65, 0x77, 0x5f, 0x63,
	0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x75, 0x75, 0x69, 0x64,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x12, 0x6e, 0x65, 0x77, 0x43, 0x61, 0x74, 0x61, 0x6c,
	0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x55, 0x75, 0x69, 0x64, 0x22, 0x67, 0x0a, 0x13, 0x50, 0x63,
	0x54, 0x61, 0x73, 0x6b, 0x57, 0x61, 0x6c, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x44, 0x61, 0x74,
	0x61, 0x12, 0x50, 0x0a, 0x0c, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x5f, 0x69, 0x74, 0x65,
	0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x2e, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61, 0x2e, 0x50, 0x63, 0x54, 0x61, 0x73, 0x6b, 0x57, 0x61,
	0x6c, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74,
	0x65, 0x6d, 0x44, 0x61, 0x74, 0x61, 0x52, 0x0b, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49,
	0x74, 0x65, 0x6d, 0x22, 0xa6, 0x01, 0x0a, 0x0f, 0x50, 0x63, 0x54, 0x61, 0x73, 0x6b, 0x57, 0x61,
	0x6c, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x36, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x6d,
	0x61, 0x72, 0x69, 0x6e, 0x61, 0x2e, 0x50, 0x63, 0x54, 0x61, 0x73, 0x6b, 0x57, 0x61, 0x6c, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12,
	0x18, 0x0a, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x12, 0x41, 0x0a, 0x0e, 0x70, 0x72, 0x6f,
	0x78, 0x79, 0x5f, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x77, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x6d, 0x61, 0x72, 0x69, 0x6e,
	0x61, 0x2e, 0x50, 0x72, 0x6f, 0x78, 0x79, 0x54, 0x61, 0x73, 0x6b, 0x57, 0x61, 0x6c, 0x52, 0x0c,
	0x70, 0x72, 0x6f, 0x78, 0x79, 0x54, 0x61, 0x73, 0x6b, 0x57, 0x61, 0x6c, 0x22, 0x5c, 0x0a, 0x0c,
	0x50, 0x72, 0x6f, 0x78, 0x79, 0x54, 0x61, 0x73, 0x6b, 0x57, 0x61, 0x6c, 0x12, 0x1b, 0x0a, 0x09,
	0x74, 0x61, 0x73, 0x6b, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x08, 0x74, 0x61, 0x73, 0x6b, 0x55, 0x75, 0x69, 0x64, 0x12, 0x2f, 0x0a, 0x13, 0x73, 0x65, 0x72,
	0x69, 0x61, 0x6c, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x69, 0x7a,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x42, 0x0f, 0x5a, 0x0d, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61,
}

var (
	file_protos_marina_marina_wal_proto_rawDescOnce sync.Once
	file_protos_marina_marina_wal_proto_rawDescData = file_protos_marina_marina_wal_proto_rawDesc
)

func file_protos_marina_marina_wal_proto_rawDescGZIP() []byte {
	file_protos_marina_marina_wal_proto_rawDescOnce.Do(func() {
		file_protos_marina_marina_wal_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_marina_marina_wal_proto_rawDescData)
	})
	return file_protos_marina_marina_wal_proto_rawDescData
}

var file_protos_marina_marina_wal_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_protos_marina_marina_wal_proto_goTypes = []interface{}{
	(*PcTaskWalRecordCatalogItemData)(nil), // 0: protos.marina.PcTaskWalRecordCatalogItemData
	(*PcTaskWalRecordData)(nil),            // 1: protos.marina.PcTaskWalRecordData
	(*PcTaskWalRecord)(nil),                // 2: protos.marina.PcTaskWalRecord
	(*ProxyTaskWal)(nil),                   // 3: protos.marina.ProxyTaskWal
}
var file_protos_marina_marina_wal_proto_depIdxs = []int32{
	0, // 0: protos.marina.PcTaskWalRecordData.catalog_item:type_name -> protos.marina.PcTaskWalRecordCatalogItemData
	1, // 1: protos.marina.PcTaskWalRecord.data:type_name -> protos.marina.PcTaskWalRecordData
	3, // 2: protos.marina.PcTaskWalRecord.proxy_task_wal:type_name -> protos.marina.ProxyTaskWal
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_protos_marina_marina_wal_proto_init() }
func file_protos_marina_marina_wal_proto_init() {
	if File_protos_marina_marina_wal_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_marina_marina_wal_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PcTaskWalRecordCatalogItemData); i {
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
		file_protos_marina_marina_wal_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PcTaskWalRecordData); i {
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
		file_protos_marina_marina_wal_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PcTaskWalRecord); i {
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
		file_protos_marina_marina_wal_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProxyTaskWal); i {
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
			RawDescriptor: file_protos_marina_marina_wal_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protos_marina_marina_wal_proto_goTypes,
		DependencyIndexes: file_protos_marina_marina_wal_proto_depIdxs,
		MessageInfos:      file_protos_marina_marina_wal_proto_msgTypes,
	}.Build()
	File_protos_marina_marina_wal_proto = out.File
	file_protos_marina_marina_wal_proto_rawDesc = nil
	file_protos_marina_marina_wal_proto_goTypes = nil
	file_protos_marina_marina_wal_proto_depIdxs = nil
}