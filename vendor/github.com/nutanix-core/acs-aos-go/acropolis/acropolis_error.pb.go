// Code generated by protoc-gen-go. DO NOT EDIT.
// source: acropolis/acropolis_error.proto

package acropolis

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

type AcropolisError_Type int32

const (
	AcropolisError_kNoError                      AcropolisError_Type = 0
	AcropolisError_kCanceled                     AcropolisError_Type = 1
	AcropolisError_kRetry                        AcropolisError_Type = 2
	AcropolisError_kTimeout                      AcropolisError_Type = 3
	AcropolisError_kNotSupported                 AcropolisError_Type = 4
	AcropolisError_kUncaughtException            AcropolisError_Type = 5
	AcropolisError_kInvalidArgument              AcropolisError_Type = 6
	AcropolisError_kInvalidState                 AcropolisError_Type = 7
	AcropolisError_kLogicalTimestampMismatch     AcropolisError_Type = 8
	AcropolisError_kExists                       AcropolisError_Type = 9
	AcropolisError_kNotFound                     AcropolisError_Type = 10
	AcropolisError_kDeviceLocked                 AcropolisError_Type = 11
	AcropolisError_kBusFull                      AcropolisError_Type = 12
	AcropolisError_kBusSlotOccupied              AcropolisError_Type = 13
	AcropolisError_kNoHostResources              AcropolisError_Type = 14
	AcropolisError_kNotManaged                   AcropolisError_Type = 15
	AcropolisError_kNotMaster                    AcropolisError_Type = 16
	AcropolisError_kMigrationFail                AcropolisError_Type = 17
	AcropolisError_kInUse                        AcropolisError_Type = 18
	AcropolisError_kAddressPoolExhausted         AcropolisError_Type = 19
	AcropolisError_kHotPlugFailure               AcropolisError_Type = 20
	AcropolisError_kHostEvacuationFailure        AcropolisError_Type = 21
	AcropolisError_kUploadFailure                AcropolisError_Type = 22
	AcropolisError_kChecksumMismatch             AcropolisError_Type = 23
	AcropolisError_kHostRestartAllVmsFailure     AcropolisError_Type = 24
	AcropolisError_kHostRestoreVmLocalityFailure AcropolisError_Type = 25
	AcropolisError_kTransportError               AcropolisError_Type = 26
	AcropolisError_kHypervisorConnectionError    AcropolisError_Type = 27
	AcropolisError_kNetworkError                 AcropolisError_Type = 28
	AcropolisError_kVmCloneError                 AcropolisError_Type = 29
	AcropolisError_kVgCloneError                 AcropolisError_Type = 30
	AcropolisError_kDelayRetry                   AcropolisError_Type = 31
	AcropolisError_kVNumaPinningFailure          AcropolisError_Type = 32
	AcropolisError_kUpgradeFailure               AcropolisError_Type = 33
	AcropolisError_kCopyFileFailure              AcropolisError_Type = 34
	AcropolisError_kNGTError                     AcropolisError_Type = 35
	AcropolisError_kInternal                     AcropolisError_Type = 36
	AcropolisError_kServiceError                 AcropolisError_Type = 37
	AcropolisError_kVmNotFound                   AcropolisError_Type = 38
	AcropolisError_kTaskNotFound                 AcropolisError_Type = 39
	AcropolisError_kSchedulerError               AcropolisError_Type = 40
	AcropolisError_kSchedulerMissingEntity       AcropolisError_Type = 41
	AcropolisError_kImportVmError                AcropolisError_Type = 42
	AcropolisError_kRemoteShell                  AcropolisError_Type = 43
	AcropolisError_kVNumaPowerOn                 AcropolisError_Type = 44
	AcropolisError_kVNumaMigrate                 AcropolisError_Type = 45
	AcropolisError_kCasError                     AcropolisError_Type = 46
	// Deprecated.
	// kDoNotRestartTask = 47;
	AcropolisError_kLazanDelayRetry                    AcropolisError_Type = 48
	AcropolisError_kImportError                        AcropolisError_Type = 49
	AcropolisError_kInvalidVmState                     AcropolisError_Type = 50
	AcropolisError_kNetworkNotFound                    AcropolisError_Type = 51
	AcropolisError_kMicrosegRuleError                  AcropolisError_Type = 52
	AcropolisError_kCerebroError                       AcropolisError_Type = 54
	AcropolisError_kVmSyncRepError                     AcropolisError_Type = 55
	AcropolisError_kVmSyncRepCannotEnable              AcropolisError_Type = 56
	AcropolisError_kVmSyncRepUpdateDormantVm           AcropolisError_Type = 57
	AcropolisError_kVmSyncRepFailedDiskProtectionError AcropolisError_Type = 58
	AcropolisError_kVmSyncRepFailedToGetVmSpecError    AcropolisError_Type = 59
	AcropolisError_kHypervisorError                    AcropolisError_Type = 61
	AcropolisError_kVmSyncRepMigrateError              AcropolisError_Type = 62
	// If VmCrossClusterLiveMigrate failed and no action is needed from the
	// caller. Possible error scenarios.
	// 1. If prechecks fail.
	// 2. If PrepareVmCrossClusterLiveMigrate failed.
	// 3. If VmCrossClusterLiveMigrate fails while updating the parcel on
	//    source, essentially identifying that the VM is powered off.
	// 4. If VmCrossClusterLiveMigrate fails and VM continues to run on the
	//    source.
	AcropolisError_kVmCrossClusterLiveMigrateError AcropolisError_Type = 64
	// CleanupPrepareVmCrossClusterLiveMigrate needs to be called on the
	// destination cluster. This error code can also be raised if
	// PrepareVmCrossClusterLiveMigrate does not respond.
	// VM is guaranteed to not run on the destination in this case.
	AcropolisError_kVmCrossClusterLiveMigratePrepareCleanupError AcropolisError_Type = 65
	// Error in managing uncommitted VM disks.
	AcropolisError_kUncommittedDiskError AcropolisError_Type = 66
	// Error while failing to delete a file.
	AcropolisError_kDeleteFileFailure AcropolisError_Type = 67
	// VM tasks are disabled for a VM.
	AcropolisError_kVmDisableUpdateError AcropolisError_Type = 68
	// Error in prechecks for cross cluster live migration.
	AcropolisError_kVmCrossClusterLiveMigratePrechecksError AcropolisError_Type = 69
	// Stretch parameters is not valid.
	AcropolisError_kVmSyncRepInvalidStretchParamsError AcropolisError_Type = 70
	// If VmCrossClusterLiveMigrate fails since acropolis cannot figure out
	// whether the VM is running on source or destination, we rely on I/O
	// being blocked on both sides resulting in the VM committing suicide
	// and the VM then getting powered on the source.
	AcropolisError_kVmCrossClusterLiveMigrateUnknownError AcropolisError_Type = 71
	// Error while handling a virtual switch operation
	AcropolisError_kVirtualSwitchError AcropolisError_Type = 72
)

var AcropolisError_Type_name = map[int32]string{
	0:  "kNoError",
	1:  "kCanceled",
	2:  "kRetry",
	3:  "kTimeout",
	4:  "kNotSupported",
	5:  "kUncaughtException",
	6:  "kInvalidArgument",
	7:  "kInvalidState",
	8:  "kLogicalTimestampMismatch",
	9:  "kExists",
	10: "kNotFound",
	11: "kDeviceLocked",
	12: "kBusFull",
	13: "kBusSlotOccupied",
	14: "kNoHostResources",
	15: "kNotManaged",
	16: "kNotMaster",
	17: "kMigrationFail",
	18: "kInUse",
	19: "kAddressPoolExhausted",
	20: "kHotPlugFailure",
	21: "kHostEvacuationFailure",
	22: "kUploadFailure",
	23: "kChecksumMismatch",
	24: "kHostRestartAllVmsFailure",
	25: "kHostRestoreVmLocalityFailure",
	26: "kTransportError",
	27: "kHypervisorConnectionError",
	28: "kNetworkError",
	29: "kVmCloneError",
	30: "kVgCloneError",
	31: "kDelayRetry",
	32: "kVNumaPinningFailure",
	33: "kUpgradeFailure",
	34: "kCopyFileFailure",
	35: "kNGTError",
	36: "kInternal",
	37: "kServiceError",
	38: "kVmNotFound",
	39: "kTaskNotFound",
	40: "kSchedulerError",
	41: "kSchedulerMissingEntity",
	42: "kImportVmError",
	43: "kRemoteShell",
	44: "kVNumaPowerOn",
	45: "kVNumaMigrate",
	46: "kCasError",
	48: "kLazanDelayRetry",
	49: "kImportError",
	50: "kInvalidVmState",
	51: "kNetworkNotFound",
	52: "kMicrosegRuleError",
	54: "kCerebroError",
	55: "kVmSyncRepError",
	56: "kVmSyncRepCannotEnable",
	57: "kVmSyncRepUpdateDormantVm",
	58: "kVmSyncRepFailedDiskProtectionError",
	59: "kVmSyncRepFailedToGetVmSpecError",
	61: "kHypervisorError",
	62: "kVmSyncRepMigrateError",
	64: "kVmCrossClusterLiveMigrateError",
	65: "kVmCrossClusterLiveMigratePrepareCleanupError",
	66: "kUncommittedDiskError",
	67: "kDeleteFileFailure",
	68: "kVmDisableUpdateError",
	69: "kVmCrossClusterLiveMigratePrechecksError",
	70: "kVmSyncRepInvalidStretchParamsError",
	71: "kVmCrossClusterLiveMigrateUnknownError",
	72: "kVirtualSwitchError",
}

var AcropolisError_Type_value = map[string]int32{
	"kNoError":                            0,
	"kCanceled":                           1,
	"kRetry":                              2,
	"kTimeout":                            3,
	"kNotSupported":                       4,
	"kUncaughtException":                  5,
	"kInvalidArgument":                    6,
	"kInvalidState":                       7,
	"kLogicalTimestampMismatch":           8,
	"kExists":                             9,
	"kNotFound":                           10,
	"kDeviceLocked":                       11,
	"kBusFull":                            12,
	"kBusSlotOccupied":                    13,
	"kNoHostResources":                    14,
	"kNotManaged":                         15,
	"kNotMaster":                          16,
	"kMigrationFail":                      17,
	"kInUse":                              18,
	"kAddressPoolExhausted":               19,
	"kHotPlugFailure":                     20,
	"kHostEvacuationFailure":              21,
	"kUploadFailure":                      22,
	"kChecksumMismatch":                   23,
	"kHostRestartAllVmsFailure":           24,
	"kHostRestoreVmLocalityFailure":       25,
	"kTransportError":                     26,
	"kHypervisorConnectionError":          27,
	"kNetworkError":                       28,
	"kVmCloneError":                       29,
	"kVgCloneError":                       30,
	"kDelayRetry":                         31,
	"kVNumaPinningFailure":                32,
	"kUpgradeFailure":                     33,
	"kCopyFileFailure":                    34,
	"kNGTError":                           35,
	"kInternal":                           36,
	"kServiceError":                       37,
	"kVmNotFound":                         38,
	"kTaskNotFound":                       39,
	"kSchedulerError":                     40,
	"kSchedulerMissingEntity":             41,
	"kImportVmError":                      42,
	"kRemoteShell":                        43,
	"kVNumaPowerOn":                       44,
	"kVNumaMigrate":                       45,
	"kCasError":                           46,
	"kLazanDelayRetry":                    48,
	"kImportError":                        49,
	"kInvalidVmState":                     50,
	"kNetworkNotFound":                    51,
	"kMicrosegRuleError":                  52,
	"kCerebroError":                       54,
	"kVmSyncRepError":                     55,
	"kVmSyncRepCannotEnable":              56,
	"kVmSyncRepUpdateDormantVm":           57,
	"kVmSyncRepFailedDiskProtectionError": 58,
	"kVmSyncRepFailedToGetVmSpecError":    59,
	"kHypervisorError":                    61,
	"kVmSyncRepMigrateError":              62,
	"kVmCrossClusterLiveMigrateError":     64,
	"kVmCrossClusterLiveMigratePrepareCleanupError": 65,
	"kUncommittedDiskError":                         66,
	"kDeleteFileFailure":                            67,
	"kVmDisableUpdateError":                         68,
	"kVmCrossClusterLiveMigratePrechecksError":      69,
	"kVmSyncRepInvalidStretchParamsError":           70,
	"kVmCrossClusterLiveMigrateUnknownError":        71,
	"kVirtualSwitchError":                           72,
}

func (x AcropolisError_Type) Enum() *AcropolisError_Type {
	p := new(AcropolisError_Type)
	*p = x
	return p
}

func (x AcropolisError_Type) String() string {
	return proto.EnumName(AcropolisError_Type_name, int32(x))
}

func (x *AcropolisError_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(AcropolisError_Type_value, data, "AcropolisError_Type")
	if err != nil {
		return err
	}
	*x = AcropolisError_Type(value)
	return nil
}

func (AcropolisError_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_9a54a472133860b6, []int{0, 0}
}

type AcropolisError struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AcropolisError) Reset()         { *m = AcropolisError{} }
func (m *AcropolisError) String() string { return proto.CompactTextString(m) }
func (*AcropolisError) ProtoMessage()    {}
func (*AcropolisError) Descriptor() ([]byte, []int) {
	return fileDescriptor_9a54a472133860b6, []int{0}
}

func (m *AcropolisError) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AcropolisError.Unmarshal(m, b)
}
func (m *AcropolisError) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AcropolisError.Marshal(b, m, deterministic)
}
func (m *AcropolisError) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AcropolisError.Merge(m, src)
}
func (m *AcropolisError) XXX_Size() int {
	return xxx_messageInfo_AcropolisError.Size(m)
}
func (m *AcropolisError) XXX_DiscardUnknown() {
	xxx_messageInfo_AcropolisError.DiscardUnknown(m)
}

var xxx_messageInfo_AcropolisError proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("nutanix.acropolis.AcropolisError_Type", AcropolisError_Type_name, AcropolisError_Type_value)
	proto.RegisterType((*AcropolisError)(nil), "nutanix.acropolis.AcropolisError")
}

func init() { proto.RegisterFile("acropolis/acropolis_error.proto", fileDescriptor_9a54a472133860b6) }

var fileDescriptor_9a54a472133860b6 = []byte{
	// 952 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x55, 0x6b, 0x73, 0x5b, 0x35,
	0x13, 0x7e, 0x5f, 0x28, 0x6d, 0xa3, 0xdc, 0x14, 0xe5, 0xd6, 0xa4, 0xa4, 0x69, 0x93, 0xd2, 0x96,
	0xd2, 0x06, 0x0a, 0x0c, 0xe5, 0x32, 0x30, 0x24, 0xb6, 0x73, 0x99, 0x49, 0x52, 0x4f, 0x6c, 0x9f,
	0x0f, 0x7c, 0x61, 0x54, 0x9d, 0x1d, 0x5b, 0xa3, 0xcb, 0x9e, 0xd1, 0x25, 0x89, 0xf9, 0x89, 0xfc,
	0x07, 0xfe, 0x0b, 0x73, 0x24, 0xe7, 0x24, 0x81, 0x29, 0xdf, 0xec, 0x67, 0x57, 0xab, 0x67, 0x9f,
	0x7d, 0x56, 0x87, 0x6c, 0x72, 0xe1, 0xb0, 0x42, 0x2d, 0xfd, 0x97, 0xcd, 0xaf, 0xdf, 0xc1, 0x39,
	0x74, 0x3b, 0x95, 0xc3, 0x80, 0x6c, 0xc1, 0xc6, 0xc0, 0xad, 0xbc, 0xdc, 0x69, 0xc2, 0x5b, 0x7f,
	0xcd, 0x90, 0xb9, 0xdd, 0xab, 0x7f, 0x9d, 0x3a, 0x77, 0xeb, 0xcf, 0x19, 0x72, 0xa7, 0x3f, 0xae,
	0x80, 0xcd, 0x90, 0xfb, 0xea, 0x14, 0x13, 0x48, 0xff, 0xc7, 0x66, 0xc9, 0x94, 0x6a, 0x71, 0x2b,
	0x40, 0x43, 0x49, 0xff, 0xcf, 0x08, 0xb9, 0xab, 0xce, 0x20, 0xb8, 0x31, 0xfd, 0x28, 0x25, 0xf6,
	0xa5, 0x01, 0x8c, 0x81, 0x7e, 0xcc, 0x16, 0xc8, 0xac, 0x3a, 0xc5, 0xd0, 0x8b, 0x55, 0x85, 0x2e,
	0x40, 0x49, 0xef, 0xb0, 0x15, 0xc2, 0xd4, 0xc0, 0x0a, 0x1e, 0x87, 0xa3, 0xd0, 0xb9, 0x14, 0x50,
	0x05, 0x89, 0x96, 0x7e, 0xc2, 0x96, 0x08, 0x55, 0x47, 0xf6, 0x9c, 0x6b, 0x59, 0xee, 0xba, 0x61,
	0x34, 0x60, 0x03, 0xbd, 0x9b, 0x0a, 0x4c, 0xd0, 0x5e, 0xe0, 0x01, 0xe8, 0x3d, 0xb6, 0x41, 0xd6,
	0xd4, 0x31, 0x0e, 0xa5, 0xe0, 0xba, 0xbe, 0xc8, 0x07, 0x6e, 0xaa, 0x13, 0xe9, 0x0d, 0x0f, 0x62,
	0x44, 0xef, 0xb3, 0x69, 0x72, 0x4f, 0x75, 0x2e, 0xa5, 0x0f, 0x9e, 0x4e, 0x25, 0xa2, 0xa7, 0x18,
	0xf6, 0x31, 0xda, 0x92, 0x92, 0x54, 0xad, 0x0d, 0xe7, 0x52, 0xc0, 0x31, 0x0a, 0x05, 0x25, 0x9d,
	0x4e, 0x7c, 0xf7, 0xa2, 0xdf, 0x8f, 0x5a, 0xd3, 0x99, 0x44, 0x62, 0x2f, 0xfa, 0x9e, 0xc6, 0xf0,
	0x4e, 0x88, 0x58, 0x49, 0x28, 0xe9, 0x6c, 0x42, 0x4f, 0xf1, 0x10, 0x7d, 0x38, 0x03, 0x8f, 0xd1,
	0x09, 0xf0, 0x74, 0x8e, 0xcd, 0x93, 0xe9, 0xba, 0xf6, 0x09, 0xb7, 0x7c, 0x08, 0x25, 0x9d, 0x67,
	0x73, 0x84, 0x64, 0xc0, 0x07, 0x70, 0x94, 0x32, 0x46, 0xe6, 0xd4, 0x89, 0x1c, 0x3a, 0x5e, 0x77,
	0xb8, 0xcf, 0xa5, 0xa6, 0x0b, 0x49, 0xaa, 0x23, 0x3b, 0xf0, 0x40, 0x19, 0x5b, 0x23, 0xcb, 0x6a,
	0xb7, 0x2c, 0x1d, 0x78, 0xdf, 0x45, 0xd4, 0x9d, 0xcb, 0x11, 0x8f, 0xbe, 0x16, 0x69, 0x91, 0x2d,
	0x92, 0x79, 0x75, 0x88, 0xa1, 0xab, 0xe3, 0xb0, 0x3e, 0x18, 0x1d, 0xd0, 0x25, 0xb6, 0x4e, 0x56,
	0x54, 0x4d, 0xa2, 0x73, 0xce, 0x45, 0x6c, 0x8a, 0xd6, 0xb1, 0xe5, 0x74, 0xd7, 0xa0, 0xd2, 0xc8,
	0xcb, 0x2b, 0x6c, 0x85, 0x2d, 0x93, 0x05, 0xd5, 0x1a, 0x81, 0x50, 0x3e, 0x9a, 0x46, 0xa0, 0xd5,
	0xa4, 0xdf, 0xa4, 0x97, 0xc0, 0x5d, 0xd8, 0xd5, 0xba, 0x30, 0xfe, 0xea, 0xd4, 0x03, 0xf6, 0x84,
	0x6c, 0x34, 0x61, 0x74, 0x50, 0x98, 0x63, 0x14, 0x5c, 0xcb, 0x30, 0xbe, 0x4a, 0x59, 0x4b, 0xec,
	0xfa, 0x8e, 0x5b, 0x5f, 0x0f, 0x35, 0x7b, 0x62, 0x9d, 0x3d, 0x22, 0xeb, 0xea, 0x70, 0x5c, 0x81,
	0x3b, 0x97, 0x1e, 0x5d, 0x0b, 0xad, 0x05, 0x51, 0x73, 0xcc, 0xf1, 0x87, 0xd9, 0x0a, 0x10, 0x2e,
	0xd0, 0xa9, 0x0c, 0x7d, 0x9a, 0xa0, 0xc2, 0xb4, 0x34, 0x5a, 0xc8, 0xd0, 0x46, 0x86, 0x86, 0x37,
	0xa0, 0x47, 0x49, 0xe7, 0x36, 0x68, 0x3e, 0xce, 0x16, 0xdb, 0x64, 0x0f, 0xc8, 0x92, 0x2a, 0x4e,
	0xa3, 0xe1, 0x5d, 0x69, 0xad, 0xb4, 0x8d, 0x42, 0x8f, 0x13, 0xb1, 0x41, 0x35, 0x74, 0xbc, 0x84,
	0x2b, 0xf0, 0x49, 0x9a, 0x5e, 0x0b, 0xab, 0xf1, 0xbe, 0xd4, 0x0d, 0xba, 0x95, 0x9d, 0x71, 0xd0,
	0xcf, 0x97, 0x6c, 0xa7, 0xbf, 0x47, 0x36, 0x80, 0xb3, 0x5c, 0xd3, 0xa7, 0x89, 0x46, 0xaf, 0x6e,
	0x45, 0x4c, 0x68, 0x7c, 0x96, 0x68, 0x14, 0xa6, 0x31, 0xd3, 0xb3, 0x94, 0xd3, 0xe7, 0xfe, 0xda,
	0x5f, 0xcf, 0xd3, 0xfd, 0x3d, 0x31, 0x82, 0x32, 0x6a, 0x70, 0xf9, 0xe0, 0x0b, 0xf6, 0x90, 0xac,
	0x5e, 0x83, 0x27, 0xd2, 0x7b, 0x69, 0x87, 0x1d, 0x1b, 0x64, 0x18, 0xd3, 0xcf, 0xd3, 0xdc, 0x8e,
	0x4c, 0xad, 0x63, 0x61, 0xf2, 0x81, 0x97, 0x8c, 0x92, 0x19, 0x75, 0x06, 0x06, 0x03, 0xf4, 0x46,
	0xa0, 0x35, 0xfd, 0x22, 0xab, 0x92, 0x3a, 0xc6, 0x0b, 0x70, 0xef, 0x2c, 0x7d, 0x75, 0x0d, 0x65,
	0x87, 0x01, 0x7d, 0x3d, 0xd9, 0xca, 0xbc, 0xb9, 0x74, 0x27, 0xf5, 0x7d, 0xcc, 0xff, 0xe0, 0xf6,
	0x86, 0x78, 0x5f, 0xa5, 0xe2, 0xf9, 0xc2, 0x9c, 0xf7, 0x26, 0x91, 0x9e, 0xac, 0x58, 0x61, 0xf2,
	0x92, 0x7d, 0x9d, 0x2d, 0x9f, 0xa7, 0xd5, 0xf4, 0xf7, 0x4d, 0xda, 0xdd, 0x13, 0x29, 0x1c, 0x7a,
	0x18, 0x9e, 0x45, 0x3d, 0xd1, 0xe6, 0xdb, 0x44, 0xa6, 0x05, 0x0e, 0xde, 0xbb, 0xc9, 0x13, 0xf1,
	0x5d, 0xaa, 0x5a, 0x98, 0xde, 0xd8, 0x8a, 0x33, 0xa8, 0x32, 0xf8, 0x36, 0x39, 0xb8, 0x01, 0x5b,
	0xdc, 0x5a, 0x0c, 0x1d, 0xcb, 0xdf, 0x6b, 0xa0, 0xdf, 0x27, 0x5b, 0x36, 0xb1, 0x41, 0x55, 0xf2,
	0x00, 0x6d, 0x74, 0x86, 0xdb, 0x50, 0x18, 0xfa, 0x03, 0x7b, 0x4e, 0xb6, 0xaf, 0xc3, 0xf5, 0x18,
	0xa1, 0x6c, 0x4b, 0xaf, 0xba, 0x0e, 0xc3, 0x4d, 0x9f, 0xfd, 0xc8, 0x9e, 0x92, 0xc7, 0xff, 0x4c,
	0xec, 0xe3, 0x01, 0x84, 0xc2, 0xf4, 0x2a, 0x10, 0x39, 0xeb, 0xa7, 0xd4, 0xdf, 0xb5, 0x5b, 0x33,
	0xfa, 0xf3, 0x6d, 0x7e, 0x13, 0x61, 0x73, 0xec, 0x17, 0xb6, 0x4d, 0x36, 0x6b, 0xb3, 0x3a, 0xf4,
	0xbe, 0xa5, 0xeb, 0x3d, 0x75, 0xc7, 0xf2, 0x1c, 0x6e, 0x25, 0xfd, 0xca, 0xde, 0x90, 0xd7, 0x1f,
	0x4e, 0xea, 0x3a, 0xa8, 0xb8, 0x83, 0x96, 0x06, 0x6e, 0xe3, 0x44, 0x93, 0xdd, 0xf4, 0x0a, 0x0c,
	0xac, 0x40, 0x63, 0x64, 0x08, 0xb9, 0xa9, 0x1c, 0xda, 0x4b, 0x72, 0xb7, 0x41, 0x43, 0x80, 0x9b,
	0xde, 0x6d, 0xa5, 0x23, 0x85, 0x69, 0x4b, 0x5f, 0x4b, 0x97, 0xa5, 0xca, 0x47, 0xda, 0xec, 0x15,
	0x79, 0xf1, 0x9f, 0x04, 0x44, 0x7a, 0x0f, 0x72, 0x76, 0xe7, 0xb6, 0xa8, 0xcd, 0x33, 0xeb, 0x20,
	0x88, 0x51, 0x97, 0x3b, 0x6e, 0x26, 0x89, 0xfb, 0xec, 0x25, 0x79, 0xf6, 0xe1, 0xb2, 0x03, 0xab,
	0x2c, 0x5e, 0x4c, 0x06, 0x70, 0xc0, 0x56, 0xc9, 0xa2, 0x2a, 0xa4, 0x0b, 0x91, 0xeb, 0xde, 0x85,
	0x0c, 0x62, 0x94, 0x03, 0x87, 0x7b, 0x6f, 0xc9, 0xb2, 0x40, 0xb3, 0xf3, 0xaf, 0x0f, 0xcf, 0xde,
	0xe2, 0xed, 0xaf, 0x4e, 0x3d, 0x53, 0xfc, 0x6d, 0xaa, 0x89, 0xff, 0x1d, 0x00, 0x00, 0xff, 0xff,
	0x10, 0xd9, 0x97, 0x86, 0xcd, 0x06, 0x00, 0x00,
}