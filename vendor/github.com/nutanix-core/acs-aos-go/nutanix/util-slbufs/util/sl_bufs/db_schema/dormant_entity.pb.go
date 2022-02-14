// Code generated by protoc-gen-go. DO NOT EDIT.
// source: util/sl_bufs/db_schema/dormant_entity.proto

package db_schema

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

//------------------------------------------------------------------------------
// Dormant VMs used for AHV sync rep enabled VMs.
type DormantVm struct {
	// The UUID of the VM being backed by this dormant VM.
	VmUuid []byte `protobuf:"bytes,1,opt,name=vm_uuid,json=vmUuid" json:"vm_uuid,omitempty"`
	// The current logical timestamp of the VM corresponding to this dormant VM.
	LogicalTimestamp *int64 `protobuf:"varint,2,opt,name=logical_timestamp,json=logicalTimestamp" json:"logical_timestamp,omitempty"`
	// VM spec.
	VmSpec *string `protobuf:"bytes,3,opt,name=vm_spec,json=vmSpec" json:"vm_spec,omitempty"`
	// For cloned VM, this is the UUID of the VM from which the VM is cloned from.
	SourceVmUuid []byte `protobuf:"bytes,4,opt,name=source_vm_uuid,json=sourceVmUuid" json:"source_vm_uuid,omitempty"`
	// Specify source VM disk UUIDs and NFS paths for cloned virtual disks.
	SourceDisks  []*DormantVm_SourceDisk  `protobuf:"bytes,5,rep,name=source_disks,json=sourceDisks" json:"source_disks,omitempty"`
	DiskUuidMaps []*DormantVm_DiskUuidMap `protobuf:"bytes,6,rep,name=disk_uuid_maps,json=diskUuidMaps" json:"disk_uuid_maps,omitempty"`
	// Store the acropolis "VmInfo" proto as a serialized string.
	AcropolisVmInfo []byte `protobuf:"bytes,7,opt,name=acropolis_vm_info,json=acropolisVmInfo" json:"acropolis_vm_info,omitempty"`
	// HA priority.
	HaPriority *int64 `protobuf:"varint,8,opt,name=ha_priority,json=haPriority" json:"ha_priority,omitempty"`
	// Incarnation IDs of IDF entities (VM and its child entities)
	VmIncarnation *DormantVm_VmIncarnation `protobuf:"bytes,9,opt,name=vm_incarnation,json=vmIncarnation" json:"vm_incarnation,omitempty"`
	// Entity Version of the VM (from the meatadata section). This is captured
	// from the current VM's metadata and will be used to create a new VM
	// from this DormantVm.
	VmEntityVersion      *string  `protobuf:"bytes,10,opt,name=vm_entity_version,json=vmEntityVersion" json:"vm_entity_version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DormantVm) Reset()         { *m = DormantVm{} }
func (m *DormantVm) String() string { return proto.CompactTextString(m) }
func (*DormantVm) ProtoMessage()    {}
func (*DormantVm) Descriptor() ([]byte, []int) {
	return fileDescriptor_bcb21020b4a1876c, []int{0}
}

func (m *DormantVm) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DormantVm.Unmarshal(m, b)
}
func (m *DormantVm) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DormantVm.Marshal(b, m, deterministic)
}
func (m *DormantVm) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DormantVm.Merge(m, src)
}
func (m *DormantVm) XXX_Size() int {
	return xxx_messageInfo_DormantVm.Size(m)
}
func (m *DormantVm) XXX_DiscardUnknown() {
	xxx_messageInfo_DormantVm.DiscardUnknown(m)
}

var xxx_messageInfo_DormantVm proto.InternalMessageInfo

func (m *DormantVm) GetVmUuid() []byte {
	if m != nil {
		return m.VmUuid
	}
	return nil
}

func (m *DormantVm) GetLogicalTimestamp() int64 {
	if m != nil && m.LogicalTimestamp != nil {
		return *m.LogicalTimestamp
	}
	return 0
}

func (m *DormantVm) GetVmSpec() string {
	if m != nil && m.VmSpec != nil {
		return *m.VmSpec
	}
	return ""
}

func (m *DormantVm) GetSourceVmUuid() []byte {
	if m != nil {
		return m.SourceVmUuid
	}
	return nil
}

func (m *DormantVm) GetSourceDisks() []*DormantVm_SourceDisk {
	if m != nil {
		return m.SourceDisks
	}
	return nil
}

func (m *DormantVm) GetDiskUuidMaps() []*DormantVm_DiskUuidMap {
	if m != nil {
		return m.DiskUuidMaps
	}
	return nil
}

func (m *DormantVm) GetAcropolisVmInfo() []byte {
	if m != nil {
		return m.AcropolisVmInfo
	}
	return nil
}

func (m *DormantVm) GetHaPriority() int64 {
	if m != nil && m.HaPriority != nil {
		return *m.HaPriority
	}
	return 0
}

func (m *DormantVm) GetVmIncarnation() *DormantVm_VmIncarnation {
	if m != nil {
		return m.VmIncarnation
	}
	return nil
}

func (m *DormantVm) GetVmEntityVersion() string {
	if m != nil && m.VmEntityVersion != nil {
		return *m.VmEntityVersion
	}
	return ""
}

// Specify the source VM disk UUID and NFS path for a cloned virtual disk.
type DormantVm_SourceDisk struct {
	// The device UUID of the cloned virtual disk.
	DeviceUuid []byte `protobuf:"bytes,1,opt,name=device_uuid,json=deviceUuid" json:"device_uuid,omitempty"`
	// The UUID of the source VM disk from which to clone the virtual disk.
	SourceVmdiskUuid []byte `protobuf:"bytes,2,opt,name=source_vmdisk_uuid,json=sourceVmdiskUuid" json:"source_vmdisk_uuid,omitempty"`
	// The NFS path from which to clone the virtual disk.
	SourceNfsPath        *string  `protobuf:"bytes,3,opt,name=source_nfs_path,json=sourceNfsPath" json:"source_nfs_path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DormantVm_SourceDisk) Reset()         { *m = DormantVm_SourceDisk{} }
func (m *DormantVm_SourceDisk) String() string { return proto.CompactTextString(m) }
func (*DormantVm_SourceDisk) ProtoMessage()    {}
func (*DormantVm_SourceDisk) Descriptor() ([]byte, []int) {
	return fileDescriptor_bcb21020b4a1876c, []int{0, 0}
}

func (m *DormantVm_SourceDisk) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DormantVm_SourceDisk.Unmarshal(m, b)
}
func (m *DormantVm_SourceDisk) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DormantVm_SourceDisk.Marshal(b, m, deterministic)
}
func (m *DormantVm_SourceDisk) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DormantVm_SourceDisk.Merge(m, src)
}
func (m *DormantVm_SourceDisk) XXX_Size() int {
	return xxx_messageInfo_DormantVm_SourceDisk.Size(m)
}
func (m *DormantVm_SourceDisk) XXX_DiscardUnknown() {
	xxx_messageInfo_DormantVm_SourceDisk.DiscardUnknown(m)
}

var xxx_messageInfo_DormantVm_SourceDisk proto.InternalMessageInfo

func (m *DormantVm_SourceDisk) GetDeviceUuid() []byte {
	if m != nil {
		return m.DeviceUuid
	}
	return nil
}

func (m *DormantVm_SourceDisk) GetSourceVmdiskUuid() []byte {
	if m != nil {
		return m.SourceVmdiskUuid
	}
	return nil
}

func (m *DormantVm_SourceDisk) GetSourceNfsPath() string {
	if m != nil && m.SourceNfsPath != nil {
		return *m.SourceNfsPath
	}
	return ""
}

// Specify a mapping from device UUIDs to VM disk UUIDs.
type DormantVm_DiskUuidMap struct {
	DeviceUuid           []byte   `protobuf:"bytes,1,opt,name=device_uuid,json=deviceUuid" json:"device_uuid,omitempty"`
	VmdiskUuid           []byte   `protobuf:"bytes,2,opt,name=vmdisk_uuid,json=vmdiskUuid" json:"vmdisk_uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DormantVm_DiskUuidMap) Reset()         { *m = DormantVm_DiskUuidMap{} }
func (m *DormantVm_DiskUuidMap) String() string { return proto.CompactTextString(m) }
func (*DormantVm_DiskUuidMap) ProtoMessage()    {}
func (*DormantVm_DiskUuidMap) Descriptor() ([]byte, []int) {
	return fileDescriptor_bcb21020b4a1876c, []int{0, 1}
}

func (m *DormantVm_DiskUuidMap) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DormantVm_DiskUuidMap.Unmarshal(m, b)
}
func (m *DormantVm_DiskUuidMap) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DormantVm_DiskUuidMap.Marshal(b, m, deterministic)
}
func (m *DormantVm_DiskUuidMap) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DormantVm_DiskUuidMap.Merge(m, src)
}
func (m *DormantVm_DiskUuidMap) XXX_Size() int {
	return xxx_messageInfo_DormantVm_DiskUuidMap.Size(m)
}
func (m *DormantVm_DiskUuidMap) XXX_DiscardUnknown() {
	xxx_messageInfo_DormantVm_DiskUuidMap.DiscardUnknown(m)
}

var xxx_messageInfo_DormantVm_DiskUuidMap proto.InternalMessageInfo

func (m *DormantVm_DiskUuidMap) GetDeviceUuid() []byte {
	if m != nil {
		return m.DeviceUuid
	}
	return nil
}

func (m *DormantVm_DiskUuidMap) GetVmdiskUuid() []byte {
	if m != nil {
		return m.VmdiskUuid
	}
	return nil
}

// This specifies the incarnation IDs for IDF entities related to this VM.
// One use case for this structure is for VM migration from one cluster to
// another. All VM related IDF entities will retain the same entity UUIDs,
// and the destination cluster needs to have AttachEntity() invoked to take
// ownership of all those VM related IDF entities.
type DormantVm_VmIncarnation struct {
	// VM incarnation ID.
	VmIncarnationId *uint64 `protobuf:"varint,1,opt,name=vm_incarnation_id,json=vmIncarnationId" json:"vm_incarnation_id,omitempty"`
	// A list of VM disk incarnation IDs.
	DiskIncarnations []*DormantVm_VmIncarnation_DiskIncarnation `protobuf:"bytes,2,rep,name=disk_incarnations,json=diskIncarnations" json:"disk_incarnations,omitempty"`
	// A list of VM NIC incarnation IDs.
	NicIncarnations      []*DormantVm_VmIncarnation_NicIncarnation `protobuf:"bytes,3,rep,name=nic_incarnations,json=nicIncarnations" json:"nic_incarnations,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                  `json:"-"`
	XXX_unrecognized     []byte                                    `json:"-"`
	XXX_sizecache        int32                                     `json:"-"`
}

func (m *DormantVm_VmIncarnation) Reset()         { *m = DormantVm_VmIncarnation{} }
func (m *DormantVm_VmIncarnation) String() string { return proto.CompactTextString(m) }
func (*DormantVm_VmIncarnation) ProtoMessage()    {}
func (*DormantVm_VmIncarnation) Descriptor() ([]byte, []int) {
	return fileDescriptor_bcb21020b4a1876c, []int{0, 2}
}

func (m *DormantVm_VmIncarnation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DormantVm_VmIncarnation.Unmarshal(m, b)
}
func (m *DormantVm_VmIncarnation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DormantVm_VmIncarnation.Marshal(b, m, deterministic)
}
func (m *DormantVm_VmIncarnation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DormantVm_VmIncarnation.Merge(m, src)
}
func (m *DormantVm_VmIncarnation) XXX_Size() int {
	return xxx_messageInfo_DormantVm_VmIncarnation.Size(m)
}
func (m *DormantVm_VmIncarnation) XXX_DiscardUnknown() {
	xxx_messageInfo_DormantVm_VmIncarnation.DiscardUnknown(m)
}

var xxx_messageInfo_DormantVm_VmIncarnation proto.InternalMessageInfo

func (m *DormantVm_VmIncarnation) GetVmIncarnationId() uint64 {
	if m != nil && m.VmIncarnationId != nil {
		return *m.VmIncarnationId
	}
	return 0
}

func (m *DormantVm_VmIncarnation) GetDiskIncarnations() []*DormantVm_VmIncarnation_DiskIncarnation {
	if m != nil {
		return m.DiskIncarnations
	}
	return nil
}

func (m *DormantVm_VmIncarnation) GetNicIncarnations() []*DormantVm_VmIncarnation_NicIncarnation {
	if m != nil {
		return m.NicIncarnations
	}
	return nil
}

type DormantVm_VmIncarnation_VmDiskAddr struct {
	// The device bus type.
	// Acceptable values for AHV: scsi, ide, pci, sata, spapr(only ppc).
	Bus *string `protobuf:"bytes,1,opt,name=bus" json:"bus,omitempty"`
	// Device index on the bus.
	Index                *uint32  `protobuf:"varint,2,opt,name=index" json:"index,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DormantVm_VmIncarnation_VmDiskAddr) Reset()         { *m = DormantVm_VmIncarnation_VmDiskAddr{} }
func (m *DormantVm_VmIncarnation_VmDiskAddr) String() string { return proto.CompactTextString(m) }
func (*DormantVm_VmIncarnation_VmDiskAddr) ProtoMessage()    {}
func (*DormantVm_VmIncarnation_VmDiskAddr) Descriptor() ([]byte, []int) {
	return fileDescriptor_bcb21020b4a1876c, []int{0, 2, 0}
}

func (m *DormantVm_VmIncarnation_VmDiskAddr) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DormantVm_VmIncarnation_VmDiskAddr.Unmarshal(m, b)
}
func (m *DormantVm_VmIncarnation_VmDiskAddr) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DormantVm_VmIncarnation_VmDiskAddr.Marshal(b, m, deterministic)
}
func (m *DormantVm_VmIncarnation_VmDiskAddr) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DormantVm_VmIncarnation_VmDiskAddr.Merge(m, src)
}
func (m *DormantVm_VmIncarnation_VmDiskAddr) XXX_Size() int {
	return xxx_messageInfo_DormantVm_VmIncarnation_VmDiskAddr.Size(m)
}
func (m *DormantVm_VmIncarnation_VmDiskAddr) XXX_DiscardUnknown() {
	xxx_messageInfo_DormantVm_VmIncarnation_VmDiskAddr.DiscardUnknown(m)
}

var xxx_messageInfo_DormantVm_VmIncarnation_VmDiskAddr proto.InternalMessageInfo

func (m *DormantVm_VmIncarnation_VmDiskAddr) GetBus() string {
	if m != nil && m.Bus != nil {
		return *m.Bus
	}
	return ""
}

func (m *DormantVm_VmIncarnation_VmDiskAddr) GetIndex() uint32 {
	if m != nil && m.Index != nil {
		return *m.Index
	}
	return 0
}

type DormantVm_VmIncarnation_DiskIncarnation struct {
	// Identify a VM disk.
	Addr *DormantVm_VmIncarnation_VmDiskAddr `protobuf:"bytes,1,opt,name=addr" json:"addr,omitempty"`
	// VM disk incarnation ID.
	DiskIncarnationId    *uint64  `protobuf:"varint,2,opt,name=disk_incarnation_id,json=diskIncarnationId" json:"disk_incarnation_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DormantVm_VmIncarnation_DiskIncarnation) Reset() {
	*m = DormantVm_VmIncarnation_DiskIncarnation{}
}
func (m *DormantVm_VmIncarnation_DiskIncarnation) String() string { return proto.CompactTextString(m) }
func (*DormantVm_VmIncarnation_DiskIncarnation) ProtoMessage()    {}
func (*DormantVm_VmIncarnation_DiskIncarnation) Descriptor() ([]byte, []int) {
	return fileDescriptor_bcb21020b4a1876c, []int{0, 2, 1}
}

func (m *DormantVm_VmIncarnation_DiskIncarnation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DormantVm_VmIncarnation_DiskIncarnation.Unmarshal(m, b)
}
func (m *DormantVm_VmIncarnation_DiskIncarnation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DormantVm_VmIncarnation_DiskIncarnation.Marshal(b, m, deterministic)
}
func (m *DormantVm_VmIncarnation_DiskIncarnation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DormantVm_VmIncarnation_DiskIncarnation.Merge(m, src)
}
func (m *DormantVm_VmIncarnation_DiskIncarnation) XXX_Size() int {
	return xxx_messageInfo_DormantVm_VmIncarnation_DiskIncarnation.Size(m)
}
func (m *DormantVm_VmIncarnation_DiskIncarnation) XXX_DiscardUnknown() {
	xxx_messageInfo_DormantVm_VmIncarnation_DiskIncarnation.DiscardUnknown(m)
}

var xxx_messageInfo_DormantVm_VmIncarnation_DiskIncarnation proto.InternalMessageInfo

func (m *DormantVm_VmIncarnation_DiskIncarnation) GetAddr() *DormantVm_VmIncarnation_VmDiskAddr {
	if m != nil {
		return m.Addr
	}
	return nil
}

func (m *DormantVm_VmIncarnation_DiskIncarnation) GetDiskIncarnationId() uint64 {
	if m != nil && m.DiskIncarnationId != nil {
		return *m.DiskIncarnationId
	}
	return 0
}

type DormantVm_VmIncarnation_NicIncarnation struct {
	// Identify a VM NIC.
	MacAddr []byte `protobuf:"bytes,1,opt,name=mac_addr,json=macAddr" json:"mac_addr,omitempty"`
	// VM NIC incarnation ID.
	NicIncarnationId     *uint64  `protobuf:"varint,2,opt,name=nic_incarnation_id,json=nicIncarnationId" json:"nic_incarnation_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DormantVm_VmIncarnation_NicIncarnation) Reset() {
	*m = DormantVm_VmIncarnation_NicIncarnation{}
}
func (m *DormantVm_VmIncarnation_NicIncarnation) String() string { return proto.CompactTextString(m) }
func (*DormantVm_VmIncarnation_NicIncarnation) ProtoMessage()    {}
func (*DormantVm_VmIncarnation_NicIncarnation) Descriptor() ([]byte, []int) {
	return fileDescriptor_bcb21020b4a1876c, []int{0, 2, 2}
}

func (m *DormantVm_VmIncarnation_NicIncarnation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DormantVm_VmIncarnation_NicIncarnation.Unmarshal(m, b)
}
func (m *DormantVm_VmIncarnation_NicIncarnation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DormantVm_VmIncarnation_NicIncarnation.Marshal(b, m, deterministic)
}
func (m *DormantVm_VmIncarnation_NicIncarnation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DormantVm_VmIncarnation_NicIncarnation.Merge(m, src)
}
func (m *DormantVm_VmIncarnation_NicIncarnation) XXX_Size() int {
	return xxx_messageInfo_DormantVm_VmIncarnation_NicIncarnation.Size(m)
}
func (m *DormantVm_VmIncarnation_NicIncarnation) XXX_DiscardUnknown() {
	xxx_messageInfo_DormantVm_VmIncarnation_NicIncarnation.DiscardUnknown(m)
}

var xxx_messageInfo_DormantVm_VmIncarnation_NicIncarnation proto.InternalMessageInfo

func (m *DormantVm_VmIncarnation_NicIncarnation) GetMacAddr() []byte {
	if m != nil {
		return m.MacAddr
	}
	return nil
}

func (m *DormantVm_VmIncarnation_NicIncarnation) GetNicIncarnationId() uint64 {
	if m != nil && m.NicIncarnationId != nil {
		return *m.NicIncarnationId
	}
	return 0
}

func init() {
	proto.RegisterType((*DormantVm)(nil), "nutanix.db_schema.DormantVm")
	proto.RegisterType((*DormantVm_SourceDisk)(nil), "nutanix.db_schema.DormantVm.SourceDisk")
	proto.RegisterType((*DormantVm_DiskUuidMap)(nil), "nutanix.db_schema.DormantVm.DiskUuidMap")
	proto.RegisterType((*DormantVm_VmIncarnation)(nil), "nutanix.db_schema.DormantVm.VmIncarnation")
	proto.RegisterType((*DormantVm_VmIncarnation_VmDiskAddr)(nil), "nutanix.db_schema.DormantVm.VmIncarnation.VmDiskAddr")
	proto.RegisterType((*DormantVm_VmIncarnation_DiskIncarnation)(nil), "nutanix.db_schema.DormantVm.VmIncarnation.DiskIncarnation")
	proto.RegisterType((*DormantVm_VmIncarnation_NicIncarnation)(nil), "nutanix.db_schema.DormantVm.VmIncarnation.NicIncarnation")
}

func init() {
	proto.RegisterFile("util/sl_bufs/db_schema/dormant_entity.proto", fileDescriptor_bcb21020b4a1876c)
}

var fileDescriptor_bcb21020b4a1876c = []byte{
	// 647 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xcf, 0x6f, 0xd3, 0x30,
	0x14, 0x56, 0xd7, 0xee, 0x47, 0x5f, 0xfa, 0xd3, 0x80, 0x08, 0xbd, 0xac, 0x42, 0x08, 0xaa, 0x6d,
	0xa4, 0xd2, 0x04, 0x07, 0xb8, 0x6d, 0x1a, 0x87, 0x22, 0x31, 0x46, 0x06, 0x95, 0xe0, 0x62, 0xb9,
	0x76, 0xda, 0x5a, 0xab, 0x9d, 0x28, 0x76, 0xa2, 0xed, 0xcc, 0x95, 0xbf, 0x84, 0xbf, 0x91, 0x03,
	0xb2, 0x93, 0xa6, 0x69, 0x37, 0x0d, 0x76, 0xf3, 0x7b, 0x79, 0xef, 0xfb, 0xbe, 0xf7, 0x3d, 0xc7,
	0x70, 0x98, 0x68, 0xbe, 0x18, 0xaa, 0x05, 0x9e, 0x24, 0x53, 0x35, 0x64, 0x13, 0xac, 0xe8, 0x3c,
	0x10, 0x64, 0xc8, 0xc2, 0x58, 0x10, 0xa9, 0x71, 0x20, 0x35, 0xd7, 0x37, 0x5e, 0x14, 0x87, 0x3a,
	0x44, 0x5d, 0x99, 0x68, 0x22, 0xf9, 0xb5, 0x57, 0xd4, 0x3d, 0xff, 0x53, 0x87, 0xfa, 0x59, 0x56,
	0x3b, 0x16, 0xe8, 0x29, 0xec, 0xa6, 0x02, 0x27, 0x09, 0x67, 0x6e, 0xa5, 0x5f, 0x19, 0x34, 0xfc,
	0x9d, 0x54, 0x7c, 0x4b, 0x38, 0x43, 0x87, 0xd0, 0x5d, 0x84, 0x33, 0x4e, 0xc9, 0x02, 0x6b, 0x2e,
	0x02, 0xa5, 0x89, 0x88, 0xdc, 0xad, 0x7e, 0x65, 0x50, 0xf5, 0x3b, 0xf9, 0x87, 0xaf, 0xcb, 0x7c,
	0x8e, 0xa2, 0xa2, 0x80, 0xba, 0xd5, 0x7e, 0x65, 0x50, 0x37, 0x28, 0x97, 0x51, 0x40, 0xd1, 0x0b,
	0x68, 0xa9, 0x30, 0x89, 0x69, 0x80, 0x97, 0x2c, 0x35, 0xcb, 0xd2, 0xc8, 0xb2, 0xe3, 0x8c, 0xeb,
	0x23, 0xe4, 0x31, 0x66, 0x5c, 0x5d, 0x29, 0x77, 0xbb, 0x5f, 0x1d, 0x38, 0xc7, 0xaf, 0xbc, 0x5b,
	0xe2, 0xbd, 0x42, 0xb8, 0x77, 0x69, 0x1b, 0xce, 0xb8, 0xba, 0xf2, 0x1d, 0x55, 0x9c, 0x15, 0x3a,
	0x87, 0x96, 0x01, 0xb1, 0x64, 0x58, 0x90, 0x48, 0xb9, 0x3b, 0x16, 0x6d, 0x70, 0x2f, 0x9a, 0xe9,
	0x35, 0x52, 0x3e, 0x91, 0xc8, 0x6f, 0xb0, 0x55, 0xa0, 0xd0, 0x01, 0x74, 0x09, 0x8d, 0xc3, 0x28,
	0x5c, 0x70, 0x65, 0x86, 0xe0, 0x72, 0x1a, 0xba, 0xbb, 0x76, 0x88, 0x76, 0xf1, 0x61, 0x2c, 0x46,
	0x72, 0x1a, 0xa2, 0x7d, 0x70, 0xe6, 0x04, 0x47, 0x31, 0x0f, 0x63, 0xae, 0x6f, 0xdc, 0x3d, 0xeb,
	0x16, 0xcc, 0xc9, 0x45, 0x9e, 0x41, 0x5f, 0xa0, 0x65, 0x21, 0x28, 0x89, 0x25, 0xd1, 0x3c, 0x94,
	0x6e, 0xbd, 0x5f, 0x19, 0x38, 0xc7, 0x07, 0xf7, 0x8a, 0x33, 0xe8, 0x45, 0x87, 0xdf, 0x4c, 0xcb,
	0xa1, 0xd1, 0x97, 0x8a, 0x7c, 0xe9, 0x38, 0x0d, 0x62, 0x65, 0x50, 0xc1, 0x2e, 0xa1, 0x9d, 0x8a,
	0x0f, 0x36, 0x3f, 0xce, 0xd2, 0xbd, 0x9f, 0x15, 0x80, 0x95, 0x6f, 0x46, 0x2e, 0x0b, 0x52, 0x4e,
	0x83, 0xf2, 0xfe, 0x21, 0x4b, 0xd9, 0xbd, 0x1c, 0x01, 0x2a, 0xb6, 0x57, 0x98, 0x6a, 0x2f, 0x41,
	0xc3, 0xef, 0x2c, 0x37, 0xb8, 0x74, 0x0b, 0xbd, 0x84, 0x76, 0x5e, 0x2d, 0xa7, 0x0a, 0x47, 0x44,
	0xcf, 0xf3, 0xcb, 0xd0, 0xcc, 0xd2, 0xe7, 0x53, 0x75, 0x41, 0xf4, 0xbc, 0xf7, 0x19, 0x9c, 0x92,
	0xdd, 0xff, 0x56, 0xb1, 0x0f, 0xce, 0x6d, 0x7a, 0x48, 0x0b, 0xe2, 0xde, 0xef, 0x1a, 0x34, 0xc7,
	0x77, 0x98, 0x52, 0xf2, 0x19, 0xe7, 0xc8, 0x35, 0x63, 0x4a, 0xa9, 0x72, 0xc4, 0xd0, 0x0c, 0xba,
	0x16, 0xbc, 0x54, 0xad, 0xdc, 0x2d, 0x7b, 0x67, 0xde, 0xff, 0xff, 0x5a, 0xec, 0x0d, 0x2a, 0xaf,
	0xa9, 0xc3, 0xd6, 0x13, 0x0a, 0x31, 0xe8, 0x48, 0x4e, 0xd7, 0x79, 0xaa, 0x96, 0xe7, 0xdd, 0x03,
	0x78, 0xce, 0x39, 0x2d, 0xd3, 0xb4, 0xe5, 0x5a, 0xac, 0x7a, 0x6f, 0x00, 0xc6, 0xc2, 0x88, 0x39,
	0x61, 0x2c, 0x46, 0x1d, 0xa8, 0x4e, 0x12, 0x65, 0x47, 0xaf, 0xfb, 0xe6, 0x88, 0x1e, 0xc3, 0x36,
	0x97, 0x2c, 0xb8, 0xb6, 0x3e, 0x36, 0xfd, 0x2c, 0xe8, 0xfd, 0xaa, 0x40, 0x7b, 0x63, 0x02, 0x34,
	0x82, 0x1a, 0x61, 0x2c, 0xb6, 0xcd, 0xce, 0xf1, 0xdb, 0x07, 0x68, 0x5c, 0x09, 0xf0, 0x2d, 0x04,
	0xf2, 0xe0, 0xd1, 0xa6, 0xc7, 0x38, 0x5f, 0x65, 0xcd, 0xef, 0x6e, 0x38, 0x35, 0x62, 0xbd, 0xef,
	0xd0, 0x5a, 0x9f, 0x13, 0x3d, 0x83, 0x3d, 0x41, 0x28, 0x2e, 0x04, 0x35, 0xfc, 0x5d, 0x41, 0xa8,
	0x9d, 0xf1, 0x08, 0xd0, 0x86, 0xaf, 0x2b, 0xec, 0xce, 0xba, 0x3d, 0x23, 0x76, 0x2a, 0xe1, 0x09,
	0x0d, 0xc5, 0xed, 0x61, 0x4e, 0x51, 0x3e, 0x4d, 0xf6, 0xcb, 0x5c, 0x98, 0xe7, 0xf3, 0xc7, 0xc9,
	0x8c, 0xeb, 0x79, 0x32, 0xf1, 0x02, 0x39, 0x2b, 0x3a, 0x68, 0x28, 0x86, 0xf9, 0x79, 0x68, 0x9e,
	0xe2, 0xd7, 0x6a, 0x61, 0x5f, 0xe2, 0xbb, 0x9f, 0xe5, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xd4,
	0x97, 0xca, 0xda, 0xaf, 0x05, 0x00, 0x00,
}