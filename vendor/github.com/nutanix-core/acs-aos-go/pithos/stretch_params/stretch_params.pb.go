// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pithos/stretch_params.proto

package pithos_stretch_params

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

// Supported entity types in entity centric stretch.
type StretchParams_EntityType int32

const (
	StretchParams_kVM StretchParams_EntityType = 1
	StretchParams_kVG StretchParams_EntityType = 2
)

var StretchParams_EntityType_name = map[int32]string{
	1: "kVM",
	2: "kVG",
}

var StretchParams_EntityType_value = map[string]int32{
	"kVM": 1,
	"kVG": 2,
}

func (x StretchParams_EntityType) Enum() *StretchParams_EntityType {
	p := new(StretchParams_EntityType)
	*p = x
	return p
}

func (x StretchParams_EntityType) String() string {
	return proto.EnumName(StretchParams_EntityType_name, int32(x))
}

func (x *StretchParams_EntityType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(StretchParams_EntityType_value, data, "StretchParams_EntityType")
	if err != nil {
		return err
	}
	*x = StretchParams_EntityType(value)
	return nil
}

func (StretchParams_EntityType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e51be5f912c7850a, []int{0, 0}
}

// Supported migration types.
// 1. kNone - indicates that no migration is in progress.
// 2. kOffline - indicates that PFO is in progress.
// 3. kLive - indicates live migration is in progress.
type StretchParams_MigrationType int32

const (
	StretchParams_kNone    StretchParams_MigrationType = 1
	StretchParams_kOffline StretchParams_MigrationType = 2
	StretchParams_kLive    StretchParams_MigrationType = 3
)

var StretchParams_MigrationType_name = map[int32]string{
	1: "kNone",
	2: "kOffline",
	3: "kLive",
}

var StretchParams_MigrationType_value = map[string]int32{
	"kNone":    1,
	"kOffline": 2,
	"kLive":    3,
}

func (x StretchParams_MigrationType) Enum() *StretchParams_MigrationType {
	p := new(StretchParams_MigrationType)
	*p = x
	return p
}

func (x StretchParams_MigrationType) String() string {
	return proto.EnumName(StretchParams_MigrationType_name, int32(x))
}

func (x *StretchParams_MigrationType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(StretchParams_MigrationType_value, data, "StretchParams_MigrationType")
	if err != nil {
		return err
	}
	*x = StretchParams_MigrationType(value)
	return nil
}

func (StretchParams_MigrationType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_e51be5f912c7850a, []int{0, 1}
}

// Contains information about stretching an entity (vdisk/vstore) to one or
// more clusters.
type StretchParams struct {
	// Universally identifies a stretched entity (vdisk/vstore) across clusters.
	// This id is constructed by the originating cluster, which sets the local
	// cluster id and cluster incarnation id. The entity id must be gauranteed
	// to be unique on the originating cluster. A stretched entity will use the
	// stretch params id across its constituent clusters. Once created, this must
	// not change.
	StretchParamsId *StretchParams_UniversalId `protobuf:"bytes,1,req,name=stretch_params_id,json=stretchParamsId" json:"stretch_params_id,omitempty"`
	// Used for identifying the current version of the stretched parameters. If
	// an event happens to the stretched entity (e.g., site partition, entity
	// takeover) then the initiating site updates the version's cluster id and
	// cluster incarnation id to itself, and increments the entity id. If the
	// version doesn't match the primary site's, then it means the sites have
	// become desynchronized.
	//
	// In the case of a transient network partition that causes the primary to
	// progress ahead of the secondary, the primary will snapshot its state at
	// its current version, and then increment the version's entity id. Upon site
	// reconnect, two cases will arise:
	// (1) The secondary remained idle, in which case the version's cluster ids
	//     remain that of the primary. The secondary is at a point in the past,
	//     and the version can be used for determining the differences between
	//     snapshots.
	// (2) The secondary performed a takeover, in which case the version's
	//     cluster ids differ. This implies a split-brain scenario has occurred,
	//     and an administrator must specify which site should reclaim ownership
	//     as the primary.
	Version *StretchParams_UniversalId `protobuf:"bytes,2,req,name=version" json:"version,omitempty"`
	// If the stretched entity is a vstore, then this contains the local
	// cluster's id of the stretched vstore.
	VstoreId *int64 `protobuf:"varint,3,opt,name=vstore_id,json=vstoreId" json:"vstore_id,omitempty"`
	// Name of the remote site this cluster should forward client-initiated
	// requests to. If not set, then this site must be the primary for the
	// stretched entity.
	//
	// Note this, if set, it may not correspond to the primary site when stretch
	// clusters are in a chain formation, where an intermediate secondary site
	// may exist between the local site and primary.
	//
	// It's always the case that at least one site will believe it's the primary
	// (more than one in the case of split-brain scenario), although the primary
	// site may be offline. In this case 'version's cluster id and cluster
	// incarnation id will be set to correspond to the local cluster.
	//
	// If the remote site doesn't exist in the Zeus config, the same logic
	// applies as if the site were down.
	ForwardRemoteName *string `protobuf:"bytes,4,opt,name=forward_remote_name,json=forwardRemoteName" json:"forward_remote_name,omitempty"`
	// If the stretched entity is a vstore (indicated by 'vstore_id' being set),
	// then this contains the id of corresponding vstore on the primary site.
	//
	// This will only be set if 'forward_remote_name' and 'vstore_id' are set.
	ForwardVstoreId *int64 `protobuf:"varint,5,opt,name=forward_vstore_id,json=forwardVstoreId" json:"forward_vstore_id,omitempty"`
	// Name of the remote site that the local cluster should replicate data to.
	// If set, then 'forward_remote_name' must not be. If the remote site doesn't
	// exist in the Zeus config, the same logic applies as if the site were down.
	ReplicateRemoteName *string `protobuf:"bytes,6,opt,name=replicate_remote_name,json=replicateRemoteName" json:"replicate_remote_name,omitempty"`
	// If the stretched entity is a vstore (indicated by 'vstore_id' being set),
	// then this contains the id of the corresponding vstore on the remote site
	// that the local cluster should replicate its data to.
	//
	// This will only be set if 'replicate_remote_name' and 'vstore_id' are set.
	ReplicateVstoreId *int64 `protobuf:"varint,7,opt,name=replicate_vstore_id,json=replicateVstoreId" json:"replicate_vstore_id,omitempty"`
	// Whether the stretched entity is undergoing data resynchronization. If
	// false, then the secondary site is up-to-date and it's possible for the
	// site to take over the stretched entity. Otherwise, the secondary site is
	// undergoing data resynchronization with the primary, and therefore doesn't
	// have the entire data set resident.
	//
	// This will be set on both the primary and secondary site to indicate that
	// resynchronization is outstanding. However, the primary will always be
	// considered to be fully synchronized.
	Resynchronizing *bool `protobuf:"varint,8,opt,name=resynchronizing" json:"resynchronizing,omitempty"`
	// Whether all replication activity should be paused. If true, it suggests
	// that an event be imminent that would cause the version to change, and
	// therefore all write IO should be quiesced and halted.
	ReadOnly *bool `protobuf:"varint,9,opt,name=read_only,json=readOnly" json:"read_only,omitempty"`
	// If true, indicates an intent to remove the stretch parameters. This should
	// only be set if it's guaranteed that all referencing vdisks and vstores are
	// also marked for removal.
	ToRemove *bool `protobuf:"varint,10,opt,name=to_remove,json=toRemove" json:"to_remove,omitempty"`
	// The last remote name that was either set in field 'forward_remote_name' or
	// 'replicate_remote_name'. This is used as a hint only, and therefore does
	// not have an affect on the current stretch parameters state.
	//
	// The last remote name can be useful for debugging, but its primary purpose
	// is for simplifying administration, where a client may wish to break and
	// restart stretching without having to re-specify the target remote.
	LastRemoteName *string `protobuf:"bytes,11,opt,name=last_remote_name,json=lastRemoteName" json:"last_remote_name,omitempty"`
	// The duration, in milliseconds, that the primary will allow a secondary
	// to be unavailable before automatically breaking replication with it. If
	// not set, then no automatic breaking of replication will be performed.
	BreakReplicationTimeoutMsecs *int64 `protobuf:"varint,12,opt,name=break_replication_timeout_msecs,json=breakReplicationTimeoutMsecs" json:"break_replication_timeout_msecs,omitempty"`
	// Group id of the stretch relation in the global witness server. This field
	// is present only for a witnessed stretch. It is a 128 bit UUID. The group
	// id is represented in different formats by different components. Pithos and
	// WAL use binary representation; Cerebro uses boost uuid object; Witness
	// client uses UUID string representation. The conversion of these
	// representations are as follows:
	//
	//        To-> |      Binary       |    boost::uuid     |    UUID string
	// From        |                   |                    |
	// ------------+-------------------+--------------------+--------------------
	// Binary      |                   | Uuid::Convert      | Uuid::ConvertUuid
	//             |                   | BinaryToUuid       | BinaryToText
	// ------------+-------------------+--------------------+--------------------
	// boost::uuid | Uuid::ConvertUuid |                    | boost::uuids::
	//             | ToBinary          |                    | to_string
	// ------------+-------------------+--------------------+--------------------
	// UUID string | Uuid::ConvertUuid | boost::uuids::     |
	//             | TextToBinary      | string_generator() |
	//
	GroupId []byte `protobuf:"bytes,13,opt,name=group_id,json=groupId" json:"group_id,omitempty"`
	// Entity version of global witness entry.
	EntityVersion *uint64 `protobuf:"varint,14,opt,name=entity_version,json=entityVersion" json:"entity_version,omitempty"`
	// This field will be set at the end of phase 1 prepare resync at secondary
	// site. This serves as a marker indicating that secondary has performed
	// clean of container, as part of action intiated by a primary change stretch
	// mode executor whose version was 'primary_site_version'.
	PrimarySiteVersion *StretchParams_UniversalId `protobuf:"bytes,15,opt,name=primary_site_version,json=primarySiteVersion" json:"primary_site_version,omitempty"`
	// If the stretched entities are VM/VG, then this field is a list of the
	// entities being stretched.
	EntityVec []*StretchParams_Entity `protobuf:"bytes,16,rep,name=entity_vec,json=entityVec" json:"entity_vec,omitempty"`
	// If the stretch entity is a vstore, this field indicates whether the
	// stretch relationship is in a decoupled state.
	// 'Decoupled' means that the local site is still in active stretch
	// relationship, but the remote site is no more in active stretch
	// relationship with this site. That is, the local site vstore is stretching
	// the IOs to the remote site but the remote site vstore is operating as a
	// standalone site.
	// NOTE: 'decoupled' is applicable only when 'replicate_remote_name' is set.
	Decoupled *bool `protobuf:"varint,17,opt,name=decoupled" json:"decoupled,omitempty"`
	// This field is populated when migration for the stretched entity is going
	// on. No other stretch transitions like break, disk additions are allowed if
	// this is set. This is applicable only in the case of AHV entity centric
	// stretch scenarios.
	MigrationType        *StretchParams_MigrationType `protobuf:"varint,18,opt,name=migration_type,json=migrationType,enum=nutanix.pithos.StretchParams_MigrationType" json:"migration_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *StretchParams) Reset()         { *m = StretchParams{} }
func (m *StretchParams) String() string { return proto.CompactTextString(m) }
func (*StretchParams) ProtoMessage()    {}
func (*StretchParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_e51be5f912c7850a, []int{0}
}

func (m *StretchParams) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StretchParams.Unmarshal(m, b)
}
func (m *StretchParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StretchParams.Marshal(b, m, deterministic)
}
func (m *StretchParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StretchParams.Merge(m, src)
}
func (m *StretchParams) XXX_Size() int {
	return xxx_messageInfo_StretchParams.Size(m)
}
func (m *StretchParams) XXX_DiscardUnknown() {
	xxx_messageInfo_StretchParams.DiscardUnknown(m)
}

var xxx_messageInfo_StretchParams proto.InternalMessageInfo

func (m *StretchParams) GetStretchParamsId() *StretchParams_UniversalId {
	if m != nil {
		return m.StretchParamsId
	}
	return nil
}

func (m *StretchParams) GetVersion() *StretchParams_UniversalId {
	if m != nil {
		return m.Version
	}
	return nil
}

func (m *StretchParams) GetVstoreId() int64 {
	if m != nil && m.VstoreId != nil {
		return *m.VstoreId
	}
	return 0
}

func (m *StretchParams) GetForwardRemoteName() string {
	if m != nil && m.ForwardRemoteName != nil {
		return *m.ForwardRemoteName
	}
	return ""
}

func (m *StretchParams) GetForwardVstoreId() int64 {
	if m != nil && m.ForwardVstoreId != nil {
		return *m.ForwardVstoreId
	}
	return 0
}

func (m *StretchParams) GetReplicateRemoteName() string {
	if m != nil && m.ReplicateRemoteName != nil {
		return *m.ReplicateRemoteName
	}
	return ""
}

func (m *StretchParams) GetReplicateVstoreId() int64 {
	if m != nil && m.ReplicateVstoreId != nil {
		return *m.ReplicateVstoreId
	}
	return 0
}

func (m *StretchParams) GetResynchronizing() bool {
	if m != nil && m.Resynchronizing != nil {
		return *m.Resynchronizing
	}
	return false
}

func (m *StretchParams) GetReadOnly() bool {
	if m != nil && m.ReadOnly != nil {
		return *m.ReadOnly
	}
	return false
}

func (m *StretchParams) GetToRemove() bool {
	if m != nil && m.ToRemove != nil {
		return *m.ToRemove
	}
	return false
}

func (m *StretchParams) GetLastRemoteName() string {
	if m != nil && m.LastRemoteName != nil {
		return *m.LastRemoteName
	}
	return ""
}

func (m *StretchParams) GetBreakReplicationTimeoutMsecs() int64 {
	if m != nil && m.BreakReplicationTimeoutMsecs != nil {
		return *m.BreakReplicationTimeoutMsecs
	}
	return 0
}

func (m *StretchParams) GetGroupId() []byte {
	if m != nil {
		return m.GroupId
	}
	return nil
}

func (m *StretchParams) GetEntityVersion() uint64 {
	if m != nil && m.EntityVersion != nil {
		return *m.EntityVersion
	}
	return 0
}

func (m *StretchParams) GetPrimarySiteVersion() *StretchParams_UniversalId {
	if m != nil {
		return m.PrimarySiteVersion
	}
	return nil
}

func (m *StretchParams) GetEntityVec() []*StretchParams_Entity {
	if m != nil {
		return m.EntityVec
	}
	return nil
}

func (m *StretchParams) GetDecoupled() bool {
	if m != nil && m.Decoupled != nil {
		return *m.Decoupled
	}
	return false
}

func (m *StretchParams) GetMigrationType() StretchParams_MigrationType {
	if m != nil && m.MigrationType != nil {
		return *m.MigrationType
	}
	return StretchParams_kNone
}

// Used to identify an entity uniquely across clusters.
type StretchParams_UniversalId struct {
	// The cluster's id, incarnation id, and an implementation-defined id
	// associated with an entity.
	ClusterId            *int64   `protobuf:"varint,1,req,name=cluster_id,json=clusterId" json:"cluster_id,omitempty"`
	ClusterIncarnationId *int64   `protobuf:"varint,2,req,name=cluster_incarnation_id,json=clusterIncarnationId" json:"cluster_incarnation_id,omitempty"`
	EntityId             *int64   `protobuf:"varint,3,req,name=entity_id,json=entityId" json:"entity_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StretchParams_UniversalId) Reset()         { *m = StretchParams_UniversalId{} }
func (m *StretchParams_UniversalId) String() string { return proto.CompactTextString(m) }
func (*StretchParams_UniversalId) ProtoMessage()    {}
func (*StretchParams_UniversalId) Descriptor() ([]byte, []int) {
	return fileDescriptor_e51be5f912c7850a, []int{0, 0}
}

func (m *StretchParams_UniversalId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StretchParams_UniversalId.Unmarshal(m, b)
}
func (m *StretchParams_UniversalId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StretchParams_UniversalId.Marshal(b, m, deterministic)
}
func (m *StretchParams_UniversalId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StretchParams_UniversalId.Merge(m, src)
}
func (m *StretchParams_UniversalId) XXX_Size() int {
	return xxx_messageInfo_StretchParams_UniversalId.Size(m)
}
func (m *StretchParams_UniversalId) XXX_DiscardUnknown() {
	xxx_messageInfo_StretchParams_UniversalId.DiscardUnknown(m)
}

var xxx_messageInfo_StretchParams_UniversalId proto.InternalMessageInfo

func (m *StretchParams_UniversalId) GetClusterId() int64 {
	if m != nil && m.ClusterId != nil {
		return *m.ClusterId
	}
	return 0
}

func (m *StretchParams_UniversalId) GetClusterIncarnationId() int64 {
	if m != nil && m.ClusterIncarnationId != nil {
		return *m.ClusterIncarnationId
	}
	return 0
}

func (m *StretchParams_UniversalId) GetEntityId() int64 {
	if m != nil && m.EntityId != nil {
		return *m.EntityId
	}
	return 0
}

type StretchParams_Entity struct {
	// Entity type of the entity supported in entity centric stretch.
	EntityType *StretchParams_EntityType `protobuf:"varint,1,opt,name=entity_type,json=entityType,enum=nutanix.pithos.StretchParams_EntityType" json:"entity_type,omitempty"`
	// Hypervisor agnostic UUID of the entity.
	EntityUuid           []byte   `protobuf:"bytes,2,opt,name=entity_uuid,json=entityUuid" json:"entity_uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StretchParams_Entity) Reset()         { *m = StretchParams_Entity{} }
func (m *StretchParams_Entity) String() string { return proto.CompactTextString(m) }
func (*StretchParams_Entity) ProtoMessage()    {}
func (*StretchParams_Entity) Descriptor() ([]byte, []int) {
	return fileDescriptor_e51be5f912c7850a, []int{0, 1}
}

func (m *StretchParams_Entity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StretchParams_Entity.Unmarshal(m, b)
}
func (m *StretchParams_Entity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StretchParams_Entity.Marshal(b, m, deterministic)
}
func (m *StretchParams_Entity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StretchParams_Entity.Merge(m, src)
}
func (m *StretchParams_Entity) XXX_Size() int {
	return xxx_messageInfo_StretchParams_Entity.Size(m)
}
func (m *StretchParams_Entity) XXX_DiscardUnknown() {
	xxx_messageInfo_StretchParams_Entity.DiscardUnknown(m)
}

var xxx_messageInfo_StretchParams_Entity proto.InternalMessageInfo

func (m *StretchParams_Entity) GetEntityType() StretchParams_EntityType {
	if m != nil && m.EntityType != nil {
		return *m.EntityType
	}
	return StretchParams_kVM
}

func (m *StretchParams_Entity) GetEntityUuid() []byte {
	if m != nil {
		return m.EntityUuid
	}
	return nil
}

func init() {
	proto.RegisterEnum("nutanix.pithos.StretchParams_EntityType", StretchParams_EntityType_name, StretchParams_EntityType_value)
	proto.RegisterEnum("nutanix.pithos.StretchParams_MigrationType", StretchParams_MigrationType_name, StretchParams_MigrationType_value)
	proto.RegisterType((*StretchParams)(nil), "nutanix.pithos.StretchParams")
	proto.RegisterType((*StretchParams_UniversalId)(nil), "nutanix.pithos.StretchParams.UniversalId")
	proto.RegisterType((*StretchParams_Entity)(nil), "nutanix.pithos.StretchParams.Entity")
}

func init() { proto.RegisterFile("pithos/stretch_params.proto", fileDescriptor_e51be5f912c7850a) }

var fileDescriptor_e51be5f912c7850a = []byte{
	// 654 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x4d, 0x4f, 0xdc, 0x3a,
	0x14, 0x55, 0x66, 0x80, 0x99, 0xdc, 0x61, 0xbe, 0x0c, 0x3c, 0xe5, 0x01, 0xef, 0x11, 0xa1, 0xf7,
	0xa4, 0xb4, 0x95, 0x52, 0x69, 0xda, 0x7d, 0x25, 0x10, 0xaa, 0x22, 0x95, 0x0f, 0x99, 0x8f, 0x45,
	0xbb, 0x88, 0xd2, 0xd8, 0x80, 0x35, 0x89, 0x1d, 0x39, 0xce, 0xb4, 0xe9, 0xa6, 0x3f, 0xbd, 0x95,
	0x9d, 0x64, 0x66, 0xc2, 0x02, 0x95, 0x5d, 0x7c, 0xce, 0xb9, 0xc7, 0xe7, 0xde, 0xd8, 0x86, 0x83,
	0x8c, 0xa9, 0x47, 0x91, 0xbf, 0xcd, 0x95, 0xa4, 0x2a, 0x7e, 0x0c, 0xb3, 0x48, 0x46, 0x69, 0xee,
	0x67, 0x52, 0x28, 0x81, 0x46, 0xbc, 0x50, 0x11, 0x67, 0xdf, 0xfd, 0x4a, 0x74, 0xfc, 0xcb, 0x86,
	0xe1, 0x75, 0x25, 0xbc, 0x32, 0x3a, 0x74, 0x0b, 0xd3, 0x76, 0x65, 0xc8, 0x88, 0x63, 0xb9, 0x1d,
	0x6f, 0x30, 0x7b, 0xe5, 0xb7, 0xab, 0xfd, 0x56, 0xa5, 0x7f, 0xcb, 0xd9, 0x82, 0xca, 0x3c, 0x4a,
	0x02, 0x82, 0xc7, 0xf9, 0x3a, 0x15, 0x10, 0x74, 0x0a, 0x3d, 0x4d, 0x32, 0xc1, 0x9d, 0xce, 0x4b,
	0xcd, 0x9a, 0x4a, 0x74, 0x00, 0xf6, 0x22, 0x57, 0x42, 0x52, 0x9d, 0xa9, 0xeb, 0x5a, 0x5e, 0x17,
	0xf7, 0x2b, 0x20, 0x20, 0xc8, 0x87, 0x9d, 0x7b, 0x21, 0xbf, 0x45, 0x92, 0x84, 0x92, 0xa6, 0x42,
	0xd1, 0x90, 0x47, 0x29, 0x75, 0x36, 0x5c, 0xcb, 0xb3, 0xf1, 0xb4, 0xa6, 0xb0, 0x61, 0x2e, 0xa2,
	0x94, 0xa2, 0xd7, 0xd0, 0x80, 0xe1, 0xca, 0x74, 0xd3, 0x98, 0x8e, 0x6b, 0xe2, 0xae, 0xf1, 0x9e,
	0xc1, 0x9e, 0xa4, 0x59, 0xc2, 0xe2, 0x48, 0xd1, 0x96, 0xfb, 0x96, 0x71, 0xdf, 0x59, 0x92, 0x6b,
	0xfe, 0x3e, 0xac, 0xe0, 0xb5, 0x1d, 0x7a, 0x66, 0x87, 0xe9, 0x92, 0x5a, 0xee, 0xe1, 0xc1, 0x58,
	0xd2, 0xbc, 0xe4, 0xf1, 0xa3, 0x14, 0x9c, 0xfd, 0x60, 0xfc, 0xc1, 0xe9, 0xbb, 0x96, 0xd7, 0xc7,
	0x4f, 0x61, 0x3d, 0x06, 0x49, 0x23, 0x12, 0x0a, 0x9e, 0x94, 0x8e, 0x6d, 0x34, 0x7d, 0x0d, 0x5c,
	0xf2, 0xa4, 0xd4, 0xa4, 0x12, 0x26, 0xe3, 0x82, 0x3a, 0x50, 0x91, 0x4a, 0x60, 0xb3, 0x46, 0x1e,
	0x4c, 0x92, 0x28, 0x57, 0xad, 0x16, 0x06, 0xa6, 0x85, 0x91, 0xc6, 0xd7, 0xd2, 0x9f, 0xc1, 0xd1,
	0x57, 0x49, 0xa3, 0x79, 0xd8, 0x04, 0x65, 0x82, 0x87, 0x8a, 0xa5, 0x54, 0x14, 0x2a, 0x4c, 0x73,
	0x1a, 0xe7, 0xce, 0xb6, 0xe9, 0xe4, 0xd0, 0xc8, 0xf0, 0x4a, 0x75, 0x53, 0x89, 0xce, 0xb5, 0x06,
	0xfd, 0x0d, 0xfd, 0x07, 0x29, 0x8a, 0x4c, 0x77, 0x3e, 0x74, 0x2d, 0x6f, 0x1b, 0xf7, 0xcc, 0x3a,
	0x20, 0xe8, 0x7f, 0x18, 0x51, 0xae, 0x98, 0x2a, 0xc3, 0xe6, 0x60, 0x8c, 0x5c, 0xcb, 0xdb, 0xc0,
	0xc3, 0x0a, 0xbd, 0xab, 0xff, 0xf9, 0x17, 0xd8, 0xcd, 0x24, 0x4b, 0x23, 0x59, 0x86, 0x39, 0xd3,
	0x93, 0xac, 0xc5, 0x63, 0xd7, 0x7a, 0xd9, 0x29, 0x42, 0xb5, 0xcd, 0x35, 0x53, 0xb4, 0x31, 0x3f,
	0x05, 0x58, 0x66, 0x88, 0x9d, 0x89, 0xdb, 0xf5, 0x06, 0xb3, 0xff, 0x9e, 0xb7, 0x3c, 0x33, 0x7a,
	0x6c, 0x37, 0x29, 0x63, 0x74, 0x08, 0x36, 0xa1, 0xb1, 0x28, 0xb2, 0x84, 0x12, 0x67, 0x6a, 0x26,
	0xbe, 0x02, 0x10, 0x86, 0x51, 0xca, 0x1e, 0x64, 0x3d, 0xc0, 0x32, 0xa3, 0x0e, 0x72, 0x2d, 0x6f,
	0x34, 0x7b, 0xf3, 0xfc, 0x36, 0xe7, 0x4d, 0xcd, 0x4d, 0x99, 0x51, 0x3c, 0x4c, 0xd7, 0x97, 0xfb,
	0x3f, 0x61, 0xb0, 0xd6, 0x19, 0xfa, 0x07, 0x20, 0x4e, 0x8a, 0x5c, 0x51, 0xd9, 0xdc, 0xd5, 0x2e,
	0xb6, 0x6b, 0x24, 0x20, 0xe8, 0x3d, 0xfc, 0xb5, 0xa4, 0x79, 0x1c, 0x49, 0x5e, 0x65, 0x61, 0xc4,
	0xdc, 0xc4, 0x2e, 0xde, 0x6d, 0xa4, 0x2b, 0x32, 0x20, 0xfa, 0x1c, 0xd5, 0xa3, 0x31, 0x77, 0x4d,
	0x0b, 0xfb, 0x15, 0x10, 0x90, 0x7d, 0x05, 0x5b, 0xd5, 0x1c, 0x50, 0x00, 0x83, 0x5a, 0x66, 0x7a,
	0xb3, 0x4c, 0x6f, 0xde, 0x9f, 0x8c, 0xd0, 0x34, 0x56, 0x8f, 0x5f, 0x7f, 0xa3, 0xa3, 0xa5, 0x55,
	0x51, 0x98, 0x70, 0xfa, 0xb8, 0xd4, 0x82, 0xdb, 0x82, 0x91, 0xe3, 0x7f, 0x01, 0x56, 0xa5, 0xa8,
	0x07, 0xdd, 0xf9, 0xdd, 0xf9, 0xc4, 0xaa, 0x3e, 0x3e, 0x4e, 0x3a, 0xc7, 0xef, 0x60, 0xd8, 0x1a,
	0x1b, 0xb2, 0x61, 0x73, 0x7e, 0x21, 0x38, 0x9d, 0x58, 0x68, 0x1b, 0xfa, 0xf3, 0xcb, 0xfb, 0xfb,
	0x84, 0x71, 0x3a, 0xe9, 0x18, 0xe2, 0x13, 0x5b, 0xd0, 0x49, 0xf7, 0xe4, 0x03, 0xa0, 0x58, 0xa4,
	0x4f, 0x02, 0x9f, 0xa0, 0x56, 0xe2, 0x2b, 0xfd, 0x76, 0x7e, 0xde, 0xab, 0x9b, 0x69, 0x3f, 0x8f,
	0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xcc, 0xeb, 0x63, 0x81, 0x70, 0x05, 0x00, 0x00,
}
