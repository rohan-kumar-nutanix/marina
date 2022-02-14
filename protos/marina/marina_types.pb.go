// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.11.4
// source: protos/marina/marina_types.proto

package marina

import (
	_ "github.com/golang/protobuf/protoc-gen-go/descriptor"
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

// Add the item type for which spec needs to be created on PC during
// migration in gflag catalog_create_spec_item_type_list.
type CatalogItemInfo_CatalogItemType int32

const (
	// Acropolis image.
	CatalogItemInfo_kImage CatalogItemInfo_CatalogItemType = 1
	// Acropolis VM snapshot.
	CatalogItemInfo_kAcropolisVmSnapshot CatalogItemInfo_CatalogItemType = 2
	// VM snapshot (cerebro).
	CatalogItemInfo_kVmSnapshot CatalogItemInfo_CatalogItemType = 3
	// File store item
	CatalogItemInfo_kFile CatalogItemInfo_CatalogItemType = 4
	// LCM module.
	CatalogItemInfo_kLCM CatalogItemInfo_CatalogItemType = 5
	// OVA package.
	CatalogItemInfo_kOVA CatalogItemInfo_CatalogItemType = 6
	// VM Template.
	CatalogItemInfo_kVmTemplate CatalogItemInfo_CatalogItemType = 7
)

// Enum value maps for CatalogItemInfo_CatalogItemType.
var (
	CatalogItemInfo_CatalogItemType_name = map[int32]string{
		1: "kImage",
		2: "kAcropolisVmSnapshot",
		3: "kVmSnapshot",
		4: "kFile",
		5: "kLCM",
		6: "kOVA",
		7: "kVmTemplate",
	}
	CatalogItemInfo_CatalogItemType_value = map[string]int32{
		"kImage":               1,
		"kAcropolisVmSnapshot": 2,
		"kVmSnapshot":          3,
		"kFile":                4,
		"kLCM":                 5,
		"kOVA":                 6,
		"kVmTemplate":          7,
	}
)

func (x CatalogItemInfo_CatalogItemType) Enum() *CatalogItemInfo_CatalogItemType {
	p := new(CatalogItemInfo_CatalogItemType)
	*p = x
	return p
}

func (x CatalogItemInfo_CatalogItemType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CatalogItemInfo_CatalogItemType) Descriptor() protoreflect.EnumDescriptor {
	return file_protos_marina_marina_types_proto_enumTypes[0].Descriptor()
}

func (CatalogItemInfo_CatalogItemType) Type() protoreflect.EnumType {
	return &file_protos_marina_marina_types_proto_enumTypes[0]
}

func (x CatalogItemInfo_CatalogItemType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *CatalogItemInfo_CatalogItemType) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = CatalogItemInfo_CatalogItemType(num)
	return nil
}

// Deprecated: Use CatalogItemInfo_CatalogItemType.Descriptor instead.
func (CatalogItemInfo_CatalogItemType) EnumDescriptor() ([]byte, []int) {
	return file_protos_marina_marina_types_proto_rawDescGZIP(), []int{0, 0}
}

type CatalogItemInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Catalog item UUID.
	Uuid []byte `protobuf:"bytes,1,opt,name=uuid" json:"uuid,omitempty"`
	// Catalog item name.
	Name *string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	// Annotation describing the catalog item.
	Annotation *string `protobuf:"bytes,3,opt,name=annotation" json:"annotation,omitempty"`
	// Catalog item type.
	ItemType *CatalogItemInfo_CatalogItemType `protobuf:"varint,4,opt,name=item_type,json=itemType,enum=protos.marina.CatalogItemInfo_CatalogItemType" json:"item_type,omitempty"`
	// Catalog item version.
	Version *int64 `protobuf:"varint,5,opt,name=version" json:"version,omitempty"`
	// Number of copies of the catalog item to maintain. Note that these copies
	// are spread across availability zones/clusters.
	NumberOfCopies *int32 `protobuf:"varint,6,opt,name=number_of_copies,json=numberOfCopies" json:"number_of_copies,omitempty"`
	// Opaque blob that can be used by the consuming service to deploy an entity
	// from a catalog item.
	Opaque []byte `protobuf:"bytes,7,opt,name=opaque" json:"opaque,omitempty"`
	// List of source groups. Each group defines the set of physical bits that
	// need to be grouped together. Catalog service uses this to figure out how
	// to replicate data across availability zones and clusters.
	SourceGroupList []*SourceGroup `protobuf:"bytes,8,rep,name=source_group_list,json=sourceGroupList" json:"source_group_list,omitempty"`
	// The global UUID groups together versions of the same catalog item. A
	// version of a catalog item can be uniquely identified by 'uuid' or by
	// ('global_uuid', 'logical_timestamp'). In addition, this field is used to
	// group together catalog items that exist on other clusters.
	GlobalCatalogItemUuid []byte `protobuf:"bytes,9,opt,name=global_catalog_item_uuid,json=globalCatalogItemUuid" json:"global_catalog_item_uuid,omitempty"`
	// Human readable version for CatalogItem
	CatalogVersion *CatalogVersion `protobuf:"bytes,10,opt,name=catalog_version,json=catalogVersion" json:"catalog_version,omitempty"`
	// List of locations where the catalog item (metadata) is currently exists.
	// For catalog items with multiple source groups the location of each source
	// group could be different from the catalog item location. This field is
	// strictly maintained by the Catalog and is not exposed via any  mutating
	// API call. Note this field is also different from other attributes on a
	// catalog item in that it can be updated without changing the version of a
	// catalog item. Also, this field should only be set from Xi portal or PC and
	// not on PE.
	LocationList []*CatalogItemInfo_CatalogItemLocation `protobuf:"bytes,11,rep,name=location_list,json=locationList" json:"location_list,omitempty"`
	// Holds uuid of cluster that is allowed to perform CUD operations
	// on this catalog item. For catalog items on PE that are not migrated
	// to PC, this will be the PE uuid. For catalog items on PE that are
	// migrated to PC, this will be the PC uuid.
	OwnerClusterUuid []byte `protobuf:"bytes,12,opt,name=owner_cluster_uuid,json=ownerClusterUuid" json:"owner_cluster_uuid,omitempty"`
	// Source catalog item from which this catalog item was created from.
	// In case of copying catalog items between a pair of PCs we use this
	// field to avoid creating multiple copies of the same catalog item of
	// same version on destination PC. Once the copie'd catalog item gets
	// updated, we erase this field. Since at that point the copied catalog
	// item has divereged and creating another copy of same catalog item
	// again is valid.
	SourceCatalogItemId *CatalogItemId `protobuf:"bytes,13,opt,name=source_catalog_item_id,json=sourceCatalogItemId" json:"source_catalog_item_id,omitempty"`
}

func (x *CatalogItemInfo) Reset() {
	*x = CatalogItemInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_marina_marina_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CatalogItemInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CatalogItemInfo) ProtoMessage() {}

func (x *CatalogItemInfo) ProtoReflect() protoreflect.Message {
	mi := &file_protos_marina_marina_types_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CatalogItemInfo.ProtoReflect.Descriptor instead.
func (*CatalogItemInfo) Descriptor() ([]byte, []int) {
	return file_protos_marina_marina_types_proto_rawDescGZIP(), []int{0}
}

func (x *CatalogItemInfo) GetUuid() []byte {
	if x != nil {
		return x.Uuid
	}
	return nil
}

func (x *CatalogItemInfo) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *CatalogItemInfo) GetAnnotation() string {
	if x != nil && x.Annotation != nil {
		return *x.Annotation
	}
	return ""
}

func (x *CatalogItemInfo) GetItemType() CatalogItemInfo_CatalogItemType {
	if x != nil && x.ItemType != nil {
		return *x.ItemType
	}
	return CatalogItemInfo_kImage
}

func (x *CatalogItemInfo) GetVersion() int64 {
	if x != nil && x.Version != nil {
		return *x.Version
	}
	return 0
}

func (x *CatalogItemInfo) GetNumberOfCopies() int32 {
	if x != nil && x.NumberOfCopies != nil {
		return *x.NumberOfCopies
	}
	return 0
}

func (x *CatalogItemInfo) GetOpaque() []byte {
	if x != nil {
		return x.Opaque
	}
	return nil
}

func (x *CatalogItemInfo) GetSourceGroupList() []*SourceGroup {
	if x != nil {
		return x.SourceGroupList
	}
	return nil
}

func (x *CatalogItemInfo) GetGlobalCatalogItemUuid() []byte {
	if x != nil {
		return x.GlobalCatalogItemUuid
	}
	return nil
}

func (x *CatalogItemInfo) GetCatalogVersion() *CatalogVersion {
	if x != nil {
		return x.CatalogVersion
	}
	return nil
}

func (x *CatalogItemInfo) GetLocationList() []*CatalogItemInfo_CatalogItemLocation {
	if x != nil {
		return x.LocationList
	}
	return nil
}

func (x *CatalogItemInfo) GetOwnerClusterUuid() []byte {
	if x != nil {
		return x.OwnerClusterUuid
	}
	return nil
}

func (x *CatalogItemInfo) GetSourceCatalogItemId() *CatalogItemId {
	if x != nil {
		return x.SourceCatalogItemId
	}
	return nil
}

type CatalogItemId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// global catalog item UUIDs to lookup
	GlobalCatalogItemUuid []byte `protobuf:"bytes,1,opt,name=global_catalog_item_uuid,json=globalCatalogItemUuid" json:"global_catalog_item_uuid,omitempty"`
	// version of the catalog item to lookup
	Version *int64 `protobuf:"varint,2,opt,name=version" json:"version,omitempty"`
}

func (x *CatalogItemId) Reset() {
	*x = CatalogItemId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_marina_marina_types_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CatalogItemId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CatalogItemId) ProtoMessage() {}

func (x *CatalogItemId) ProtoReflect() protoreflect.Message {
	mi := &file_protos_marina_marina_types_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CatalogItemId.ProtoReflect.Descriptor instead.
func (*CatalogItemId) Descriptor() ([]byte, []int) {
	return file_protos_marina_marina_types_proto_rawDescGZIP(), []int{1}
}

func (x *CatalogItemId) GetGlobalCatalogItemUuid() []byte {
	if x != nil {
		return x.GlobalCatalogItemUuid
	}
	return nil
}

func (x *CatalogItemId) GetVersion() int64 {
	if x != nil && x.Version != nil {
		return *x.Version
	}
	return 0
}

type CatalogVersion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the product/version of the catalog item.
	ProductName *string `protobuf:"bytes,1,opt,name=product_name,json=productName" json:"product_name,omitempty"`
	// Version string of the catalog item.
	ProductVersion *string `protobuf:"bytes,2,opt,name=product_version,json=productVersion" json:"product_version,omitempty"`
}

func (x *CatalogVersion) Reset() {
	*x = CatalogVersion{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_marina_marina_types_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CatalogVersion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CatalogVersion) ProtoMessage() {}

func (x *CatalogVersion) ProtoReflect() protoreflect.Message {
	mi := &file_protos_marina_marina_types_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CatalogVersion.ProtoReflect.Descriptor instead.
func (*CatalogVersion) Descriptor() ([]byte, []int) {
	return file_protos_marina_marina_types_proto_rawDescGZIP(), []int{2}
}

func (x *CatalogVersion) GetProductName() string {
	if x != nil && x.ProductName != nil {
		return *x.ProductName
	}
	return ""
}

func (x *CatalogVersion) GetProductVersion() string {
	if x != nil && x.ProductVersion != nil {
		return *x.ProductVersion
	}
	return ""
}

type SourceGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// SourceGroup UUID
	Uuid []byte `protobuf:"bytes,1,opt,name=uuid" json:"uuid,omitempty"`
	// List of UUIDs pointing at the clusters where this source group must
	// reside.
	ClusterUuidList [][]byte `protobuf:"bytes,2,rep,name=cluster_uuid_list,json=clusterUuidList" json:"cluster_uuid_list,omitempty"`
	// List of logical IDs pointing at the availability zones where this source
	// group must reside.
	AvailabilityZoneLogicalIdList []string `protobuf:"bytes,3,rep,name=availability_zone_logical_id_list,json=availabilityZoneLogicalIdList" json:"availability_zone_logical_id_list,omitempty"`
	// List of sources.
	SourceList []*Source `protobuf:"bytes,4,rep,name=source_list,json=sourceList" json:"source_list,omitempty"`
}

func (x *SourceGroup) Reset() {
	*x = SourceGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_marina_marina_types_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SourceGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SourceGroup) ProtoMessage() {}

func (x *SourceGroup) ProtoReflect() protoreflect.Message {
	mi := &file_protos_marina_marina_types_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SourceGroup.ProtoReflect.Descriptor instead.
func (*SourceGroup) Descriptor() ([]byte, []int) {
	return file_protos_marina_marina_types_proto_rawDescGZIP(), []int{3}
}

func (x *SourceGroup) GetUuid() []byte {
	if x != nil {
		return x.Uuid
	}
	return nil
}

func (x *SourceGroup) GetClusterUuidList() [][]byte {
	if x != nil {
		return x.ClusterUuidList
	}
	return nil
}

func (x *SourceGroup) GetAvailabilityZoneLogicalIdList() []string {
	if x != nil {
		return x.AvailabilityZoneLogicalIdList
	}
	return nil
}

func (x *SourceGroup) GetSourceList() []*Source {
	if x != nil {
		return x.SourceList
	}
	return nil
}

type Source struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// UUID pointing at a child catalog item.
	CatalogItemUuid []byte `protobuf:"bytes,1,opt,name=catalog_item_uuid,json=catalogItemUuid" json:"catalog_item_uuid,omitempty"`
	// UUID pointing at a file within the file repository.
	FileUuid []byte `protobuf:"bytes,2,opt,name=file_uuid,json=fileUuid" json:"file_uuid,omitempty"`
	// UUID pointing at a cerebro snapshot.
	SnapshotUuid []byte `protobuf:"bytes,3,opt,name=snapshot_uuid,json=snapshotUuid" json:"snapshot_uuid,omitempty"`
}

func (x *Source) Reset() {
	*x = Source{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_marina_marina_types_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Source) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Source) ProtoMessage() {}

func (x *Source) ProtoReflect() protoreflect.Message {
	mi := &file_protos_marina_marina_types_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Source.ProtoReflect.Descriptor instead.
func (*Source) Descriptor() ([]byte, []int) {
	return file_protos_marina_marina_types_proto_rawDescGZIP(), []int{4}
}

func (x *Source) GetCatalogItemUuid() []byte {
	if x != nil {
		return x.CatalogItemUuid
	}
	return nil
}

func (x *Source) GetFileUuid() []byte {
	if x != nil {
		return x.FileUuid
	}
	return nil
}

func (x *Source) GetSnapshotUuid() []byte {
	if x != nil {
		return x.SnapshotUuid
	}
	return nil
}

type CatalogItemInfo_CatalogItemLocation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Cluster UUID.
	ClusterUuid []byte `protobuf:"bytes,1,opt,name=cluster_uuid,json=clusterUuid" json:"cluster_uuid,omitempty"`
	// Availability zone logical id.
	AvailabilityZoneLogicalId *string `protobuf:"bytes,2,opt,name=availability_zone_logical_id,json=availabilityZoneLogicalId" json:"availability_zone_logical_id,omitempty"`
}

func (x *CatalogItemInfo_CatalogItemLocation) Reset() {
	*x = CatalogItemInfo_CatalogItemLocation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_marina_marina_types_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CatalogItemInfo_CatalogItemLocation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CatalogItemInfo_CatalogItemLocation) ProtoMessage() {}

func (x *CatalogItemInfo_CatalogItemLocation) ProtoReflect() protoreflect.Message {
	mi := &file_protos_marina_marina_types_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CatalogItemInfo_CatalogItemLocation.ProtoReflect.Descriptor instead.
func (*CatalogItemInfo_CatalogItemLocation) Descriptor() ([]byte, []int) {
	return file_protos_marina_marina_types_proto_rawDescGZIP(), []int{0, 0}
}

func (x *CatalogItemInfo_CatalogItemLocation) GetClusterUuid() []byte {
	if x != nil {
		return x.ClusterUuid
	}
	return nil
}

func (x *CatalogItemInfo_CatalogItemLocation) GetAvailabilityZoneLogicalId() string {
	if x != nil && x.AvailabilityZoneLogicalId != nil {
		return *x.AvailabilityZoneLogicalId
	}
	return ""
}

var File_protos_marina_marina_types_proto protoreflect.FileDescriptor

var file_protos_marina_marina_types_proto_rawDesc = []byte{
	0x0a, 0x20, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61, 0x2f,
	0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x6d, 0x61, 0x72, 0x69, 0x6e,
	0x61, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x9a, 0x07, 0x0a, 0x0f, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49,
	0x74, 0x65, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x1e, 0x0a, 0x0a, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x4b, 0x0a, 0x09, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x2e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x6d, 0x61, 0x72, 0x69,
	0x6e, 0x61, 0x2e, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x49, 0x6e,
	0x66, 0x6f, 0x2e, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x08, 0x69, 0x74, 0x65, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x0a, 0x10, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x5f, 0x6f, 0x66, 0x5f, 0x63, 0x6f, 0x70, 0x69, 0x65, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0e, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x4f, 0x66, 0x43, 0x6f, 0x70, 0x69, 0x65, 0x73,
	0x12, 0x16, 0x0a, 0x06, 0x6f, 0x70, 0x61, 0x71, 0x75, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x06, 0x6f, 0x70, 0x61, 0x71, 0x75, 0x65, 0x12, 0x46, 0x0a, 0x11, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x5f, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x08, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x6d, 0x61, 0x72,
	0x69, 0x6e, 0x61, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52,
	0x0f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x4c, 0x69, 0x73, 0x74,
	0x12, 0x37, 0x0a, 0x18, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x5f, 0x63, 0x61, 0x74, 0x61, 0x6c,
	0x6f, 0x67, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x15, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f,
	0x67, 0x49, 0x74, 0x65, 0x6d, 0x55, 0x75, 0x69, 0x64, 0x12, 0x46, 0x0a, 0x0f, 0x63, 0x61, 0x74,
	0x61, 0x6c, 0x6f, 0x67, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x6d, 0x61, 0x72, 0x69,
	0x6e, 0x61, 0x2e, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x0e, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x57, 0x0a, 0x0d, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6c, 0x69,
	0x73, 0x74, 0x18, 0x0b, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x73, 0x2e, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61, 0x2e, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67,
	0x49, 0x74, 0x65, 0x6d, 0x49, 0x6e, 0x66, 0x6f, 0x2e, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67,
	0x49, 0x74, 0x65, 0x6d, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0c, 0x6c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x12, 0x6f, 0x77,
	0x6e, 0x65, 0x72, 0x5f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x75, 0x75, 0x69, 0x64,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x10, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x43, 0x6c, 0x75,
	0x73, 0x74, 0x65, 0x72, 0x55, 0x75, 0x69, 0x64, 0x12, 0x51, 0x0a, 0x16, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x5f, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x5f,
	0x69, 0x64, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x73, 0x2e, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61, 0x2e, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67,
	0x49, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x52, 0x13, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x43, 0x61,
	0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x1a, 0x79, 0x0a, 0x13, 0x43,
	0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x75, 0x75,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x55, 0x75, 0x69, 0x64, 0x12, 0x3f, 0x0a, 0x1c, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62,
	0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x7a, 0x6f, 0x6e, 0x65, 0x5f, 0x6c, 0x6f, 0x67, 0x69, 0x63,
	0x61, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x19, 0x61, 0x76, 0x61,
	0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5a, 0x6f, 0x6e, 0x65, 0x4c, 0x6f, 0x67,
	0x69, 0x63, 0x61, 0x6c, 0x49, 0x64, 0x22, 0x78, 0x0a, 0x0f, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f,
	0x67, 0x49, 0x74, 0x65, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x6b, 0x49, 0x6d,
	0x61, 0x67, 0x65, 0x10, 0x01, 0x12, 0x18, 0x0a, 0x14, 0x6b, 0x41, 0x63, 0x72, 0x6f, 0x70, 0x6f,
	0x6c, 0x69, 0x73, 0x56, 0x6d, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x10, 0x02, 0x12,
	0x0f, 0x0a, 0x0b, 0x6b, 0x56, 0x6d, 0x53, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x10, 0x03,
	0x12, 0x09, 0x0a, 0x05, 0x6b, 0x46, 0x69, 0x6c, 0x65, 0x10, 0x04, 0x12, 0x08, 0x0a, 0x04, 0x6b,
	0x4c, 0x43, 0x4d, 0x10, 0x05, 0x12, 0x08, 0x0a, 0x04, 0x6b, 0x4f, 0x56, 0x41, 0x10, 0x06, 0x12,
	0x0f, 0x0a, 0x0b, 0x6b, 0x56, 0x6d, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x10, 0x07,
	0x22, 0x62, 0x0a, 0x0d, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x49,
	0x64, 0x12, 0x37, 0x0a, 0x18, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x5f, 0x63, 0x61, 0x74, 0x61,
	0x6c, 0x6f, 0x67, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x15, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x43, 0x61, 0x74, 0x61, 0x6c,
	0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x55, 0x75, 0x69, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x76, 0x65, 0x72,
	0x73, 0x69, 0x6f, 0x6e, 0x22, 0x5c, 0x0a, 0x0e, 0x43, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x22, 0xcf, 0x01, 0x0a, 0x0b, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x47, 0x72, 0x6f,
	0x75, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x2a, 0x0a, 0x11, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65,
	0x72, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0c, 0x52, 0x0f, 0x63, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x55, 0x75, 0x69, 0x64, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x48, 0x0a, 0x21, 0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69,
	0x74, 0x79, 0x5f, 0x7a, 0x6f, 0x6e, 0x65, 0x5f, 0x6c, 0x6f, 0x67, 0x69, 0x63, 0x61, 0x6c, 0x5f,
	0x69, 0x64, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x1d, 0x61,
	0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5a, 0x6f, 0x6e, 0x65, 0x4c,
	0x6f, 0x67, 0x69, 0x63, 0x61, 0x6c, 0x49, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x36, 0x0a, 0x0b,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x6d, 0x61, 0x72, 0x69, 0x6e,
	0x61, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x4c, 0x69, 0x73, 0x74, 0x22, 0x76, 0x0a, 0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x2a,
	0x0a, 0x11, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x75,
	0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0f, 0x63, 0x61, 0x74, 0x61, 0x6c,
	0x6f, 0x67, 0x49, 0x74, 0x65, 0x6d, 0x55, 0x75, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69,
	0x6c, 0x65, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x66,
	0x69, 0x6c, 0x65, 0x55, 0x75, 0x69, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x6e, 0x61, 0x70, 0x73,
	0x68, 0x6f, 0x74, 0x5f, 0x75, 0x75, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c,
	0x73, 0x6e, 0x61, 0x70, 0x73, 0x68, 0x6f, 0x74, 0x55, 0x75, 0x69, 0x64, 0x42, 0x0f, 0x5a, 0x0d,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x61,
}

var (
	file_protos_marina_marina_types_proto_rawDescOnce sync.Once
	file_protos_marina_marina_types_proto_rawDescData = file_protos_marina_marina_types_proto_rawDesc
)

func file_protos_marina_marina_types_proto_rawDescGZIP() []byte {
	file_protos_marina_marina_types_proto_rawDescOnce.Do(func() {
		file_protos_marina_marina_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_marina_marina_types_proto_rawDescData)
	})
	return file_protos_marina_marina_types_proto_rawDescData
}

var file_protos_marina_marina_types_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_protos_marina_marina_types_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_protos_marina_marina_types_proto_goTypes = []interface{}{
	(CatalogItemInfo_CatalogItemType)(0),        // 0: protos.marina.CatalogItemInfo.CatalogItemType
	(*CatalogItemInfo)(nil),                     // 1: protos.marina.CatalogItemInfo
	(*CatalogItemId)(nil),                       // 2: protos.marina.CatalogItemId
	(*CatalogVersion)(nil),                      // 3: protos.marina.CatalogVersion
	(*SourceGroup)(nil),                         // 4: protos.marina.SourceGroup
	(*Source)(nil),                              // 5: protos.marina.Source
	(*CatalogItemInfo_CatalogItemLocation)(nil), // 6: protos.marina.CatalogItemInfo.CatalogItemLocation
}
var file_protos_marina_marina_types_proto_depIdxs = []int32{
	0, // 0: protos.marina.CatalogItemInfo.item_type:type_name -> protos.marina.CatalogItemInfo.CatalogItemType
	4, // 1: protos.marina.CatalogItemInfo.source_group_list:type_name -> protos.marina.SourceGroup
	3, // 2: protos.marina.CatalogItemInfo.catalog_version:type_name -> protos.marina.CatalogVersion
	6, // 3: protos.marina.CatalogItemInfo.location_list:type_name -> protos.marina.CatalogItemInfo.CatalogItemLocation
	2, // 4: protos.marina.CatalogItemInfo.source_catalog_item_id:type_name -> protos.marina.CatalogItemId
	5, // 5: protos.marina.SourceGroup.source_list:type_name -> protos.marina.Source
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_protos_marina_marina_types_proto_init() }
func file_protos_marina_marina_types_proto_init() {
	if File_protos_marina_marina_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_marina_marina_types_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CatalogItemInfo); i {
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
		file_protos_marina_marina_types_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CatalogItemId); i {
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
		file_protos_marina_marina_types_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CatalogVersion); i {
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
		file_protos_marina_marina_types_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SourceGroup); i {
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
		file_protos_marina_marina_types_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Source); i {
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
		file_protos_marina_marina_types_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CatalogItemInfo_CatalogItemLocation); i {
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
			RawDescriptor: file_protos_marina_marina_types_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protos_marina_marina_types_proto_goTypes,
		DependencyIndexes: file_protos_marina_marina_types_proto_depIdxs,
		EnumInfos:         file_protos_marina_marina_types_proto_enumTypes,
		MessageInfos:      file_protos_marina_marina_types_proto_msgTypes,
	}.Build()
	File_protos_marina_marina_types_proto = out.File
	file_protos_marina_marina_types_proto_rawDesc = nil
	file_protos_marina_marina_types_proto_goTypes = nil
	file_protos_marina_marina_types_proto_depIdxs = nil
}