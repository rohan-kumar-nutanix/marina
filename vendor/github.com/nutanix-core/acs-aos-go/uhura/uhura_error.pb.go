// Code generated by protoc-gen-go. DO NOT EDIT.
// source: uhura/uhura_error.proto

package uhura

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

type UhuraError_Type int32

const (
	// next index=90
	UhuraError_kNoError UhuraError_Type = 0
	// Error codes for internal system communication issues.
	UhuraError_kCanceled       UhuraError_Type = 1
	UhuraError_kTimeout        UhuraError_Type = 2
	UhuraError_kTransportError UhuraError_Type = 3
	// Application logic specific error codes.
	UhuraError_kAddressPoolExhausted            UhuraError_Type = 4
	UhuraError_kBusFull                         UhuraError_Type = 5
	UhuraError_kBusSlotOccupied                 UhuraError_Type = 6
	UhuraError_kChecksumMismatch                UhuraError_Type = 7
	UhuraError_kDeviceLocked                    UhuraError_Type = 8
	UhuraError_kExists                          UhuraError_Type = 9
	UhuraError_kHostEvacuationFailure           UhuraError_Type = 10
	UhuraError_kHostRestartAllVmsFailure        UhuraError_Type = 11
	UhuraError_kHostRestoreVmLocalityFailure    UhuraError_Type = 12
	UhuraError_kHotPlugFailure                  UhuraError_Type = 13
	UhuraError_kInUse                           UhuraError_Type = 14
	UhuraError_kInvalidArgument                 UhuraError_Type = 15
	UhuraError_kInvalidNode                     UhuraError_Type = 16
	UhuraError_kInvalidState                    UhuraError_Type = 17
	UhuraError_kNotManaged                      UhuraError_Type = 18
	UhuraError_kNotMaster                       UhuraError_Type = 19
	UhuraError_kVirtualDiskNotFound             UhuraError_Type = 20
	UhuraError_kVirtualDiskUpdateFailed         UhuraError_Type = 21
	UhuraError_kVmNotFound                      UhuraError_Type = 22
	UhuraError_kHostConnectionFailure           UhuraError_Type = 34
	UhuraError_kDeviceSlotAlreadyTaken          UhuraError_Type = 42
	UhuraError_kClusterReconfigRuntimeError     UhuraError_Type = 43
	UhuraError_kGuestToolsNotAvaliableError     UhuraError_Type = 44
	UhuraError_kGuestToolsRuntimeError          UhuraError_Type = 45
	UhuraError_kAffinityRuleNotFoundError       UhuraError_Type = 46
	UhuraError_kCustomizationFaultError         UhuraError_Type = 47
	UhuraError_kInvalidDatastoreError           UhuraError_Type = 48
	UhuraError_kTaskInProgressError             UhuraError_Type = 49
	UhuraError_kVmConfigFaultError              UhuraError_Type = 50
	UhuraError_kFileFaultError                  UhuraError_Type = 51
	UhuraError_kMigrationFaultError             UhuraError_Type = 52
	UhuraError_kInsufficientResourcesFaultError UhuraError_Type = 53
	UhuraError_kGuestToolsUpgradeFaultError     UhuraError_Type = 59
	UhuraError_kInvalidLogin                    UhuraError_Type = 62
	UhuraError_kVmUnregisterRuntimeError        UhuraError_Type = 65
	UhuraError_kVmAlreadyExistsError            UhuraError_Type = 75
	UhuraError_kNetworkError                    UhuraError_Type = 77
	UhuraError_kRemoteShellConnectionError      UhuraError_Type = 78
	UhuraError_kRemoteShellMaxLimitReachedError UhuraError_Type = 79
	UhuraError_kVmSnapshotNotFoundError         UhuraError_Type = 80
	UhuraError_kUuidMappingDuplicateValueError  UhuraError_Type = 82
	UhuraError_kStateTransitionTimeoutError     UhuraError_Type = 85
	UhuraError_kVmAlreadyUpgradedError          UhuraError_Type = 87
	// Other miscellaneous error codes.
	UhuraError_kHypervisorNotSupported          UhuraError_Type = 23
	UhuraError_kOperationNotSupported           UhuraError_Type = 24
	UhuraError_kRetry                           UhuraError_Type = 25
	UhuraError_kUncaughtException               UhuraError_Type = 26
	UhuraError_kNotSupported                    UhuraError_Type = 27
	UhuraError_kVmCustomizationIsoCreateFailure UhuraError_Type = 28
	UhuraError_kInvalidOperation                UhuraError_Type = 29
	UhuraError_kNotFound                        UhuraError_Type = 30
	UhuraError_kNoHostResources                 UhuraError_Type = 31
	UhuraError_kInternalTaskCreation            UhuraError_Type = 32
	UhuraError_kDbError                         UhuraError_Type = 33
	UhuraError_kVmFlashModeWatcherError         UhuraError_Type = 35
	UhuraError_kOperationFailed                 UhuraError_Type = 36
	UhuraError_kHyperintClientError             UhuraError_Type = 37
	UhuraError_kRpcClientError                  UhuraError_Type = 38
	UhuraError_kVmInoperableError               UhuraError_Type = 39
	UhuraError_kDbEntityNotFoundError           UhuraError_Type = 40
	UhuraError_kDbNotUpdatedError               UhuraError_Type = 41
	UhuraError_kRestrictedOperationError        UhuraError_Type = 60
	UhuraError_kNetworkNotFound                 UhuraError_Type = 81
	UhuraError_kFlashModeUpdateError            UhuraError_Type = 88
	// Vcenter specific error codes.
	UhuraError_kManagementServerConnectionFailure  UhuraError_Type = 54
	UhuraError_kManagementServerCertificateError   UhuraError_Type = 55
	UhuraError_kManagementServerCredentialsError   UhuraError_Type = 56
	UhuraError_kManagementServerRegisterError      UhuraError_Type = 57
	UhuraError_kManagementServerUnregisterError    UhuraError_Type = 58
	UhuraError_kManagementServerNotRegisteredError UhuraError_Type = 61
	UhuraError_kMultipleManagementServersError     UhuraError_Type = 64
	UhuraError_kManagementServerConnectionError    UhuraError_Type = 83
	UhuraError_kManagementServerNoPermissionError  UhuraError_Type = 84
	UhuraError_kOrphanedVmExistsError              UhuraError_Type = 86
	UhuraError_kVcenterStaleVmError                UhuraError_Type = 89
	// NGT client error.
	UhuraError_kNgtClientError UhuraError_Type = 63
	// Cloud specific error codes.
	UhuraError_kCloudException                   UhuraError_Type = 66
	UhuraError_kCloudResourceLimitExceeded       UhuraError_Type = 67
	UhuraError_kCloudServiceInternalError        UhuraError_Type = 68
	UhuraError_kBootVolumeTypeNotSupported       UhuraError_Type = 69
	UhuraError_kCloudServiceUnAvailable          UhuraError_Type = 70
	UhuraError_kCloudElasticIPMigrateError       UhuraError_Type = 71
	UhuraError_kUnknownVolumeTypeError           UhuraError_Type = 72
	UhuraError_kCloudResourceIncorrectStateError UhuraError_Type = 73
	UhuraError_kCloudConnectionFailure           UhuraError_Type = 74
	UhuraError_kInstanceTagNotCreatedError       UhuraError_Type = 76
)

var UhuraError_Type_name = map[int32]string{
	0:  "kNoError",
	1:  "kCanceled",
	2:  "kTimeout",
	3:  "kTransportError",
	4:  "kAddressPoolExhausted",
	5:  "kBusFull",
	6:  "kBusSlotOccupied",
	7:  "kChecksumMismatch",
	8:  "kDeviceLocked",
	9:  "kExists",
	10: "kHostEvacuationFailure",
	11: "kHostRestartAllVmsFailure",
	12: "kHostRestoreVmLocalityFailure",
	13: "kHotPlugFailure",
	14: "kInUse",
	15: "kInvalidArgument",
	16: "kInvalidNode",
	17: "kInvalidState",
	18: "kNotManaged",
	19: "kNotMaster",
	20: "kVirtualDiskNotFound",
	21: "kVirtualDiskUpdateFailed",
	22: "kVmNotFound",
	34: "kHostConnectionFailure",
	42: "kDeviceSlotAlreadyTaken",
	43: "kClusterReconfigRuntimeError",
	44: "kGuestToolsNotAvaliableError",
	45: "kGuestToolsRuntimeError",
	46: "kAffinityRuleNotFoundError",
	47: "kCustomizationFaultError",
	48: "kInvalidDatastoreError",
	49: "kTaskInProgressError",
	50: "kVmConfigFaultError",
	51: "kFileFaultError",
	52: "kMigrationFaultError",
	53: "kInsufficientResourcesFaultError",
	59: "kGuestToolsUpgradeFaultError",
	62: "kInvalidLogin",
	65: "kVmUnregisterRuntimeError",
	75: "kVmAlreadyExistsError",
	77: "kNetworkError",
	78: "kRemoteShellConnectionError",
	79: "kRemoteShellMaxLimitReachedError",
	80: "kVmSnapshotNotFoundError",
	82: "kUuidMappingDuplicateValueError",
	85: "kStateTransitionTimeoutError",
	87: "kVmAlreadyUpgradedError",
	23: "kHypervisorNotSupported",
	24: "kOperationNotSupported",
	25: "kRetry",
	26: "kUncaughtException",
	27: "kNotSupported",
	28: "kVmCustomizationIsoCreateFailure",
	29: "kInvalidOperation",
	30: "kNotFound",
	31: "kNoHostResources",
	32: "kInternalTaskCreation",
	33: "kDbError",
	35: "kVmFlashModeWatcherError",
	36: "kOperationFailed",
	37: "kHyperintClientError",
	38: "kRpcClientError",
	39: "kVmInoperableError",
	40: "kDbEntityNotFoundError",
	41: "kDbNotUpdatedError",
	60: "kRestrictedOperationError",
	81: "kNetworkNotFound",
	88: "kFlashModeUpdateError",
	54: "kManagementServerConnectionFailure",
	55: "kManagementServerCertificateError",
	56: "kManagementServerCredentialsError",
	57: "kManagementServerRegisterError",
	58: "kManagementServerUnregisterError",
	61: "kManagementServerNotRegisteredError",
	64: "kMultipleManagementServersError",
	83: "kManagementServerConnectionError",
	84: "kManagementServerNoPermissionError",
	86: "kOrphanedVmExistsError",
	89: "kVcenterStaleVmError",
	63: "kNgtClientError",
	66: "kCloudException",
	67: "kCloudResourceLimitExceeded",
	68: "kCloudServiceInternalError",
	69: "kBootVolumeTypeNotSupported",
	70: "kCloudServiceUnAvailable",
	71: "kCloudElasticIPMigrateError",
	72: "kUnknownVolumeTypeError",
	73: "kCloudResourceIncorrectStateError",
	74: "kCloudConnectionFailure",
	76: "kInstanceTagNotCreatedError",
}

var UhuraError_Type_value = map[string]int32{
	"kNoError":                            0,
	"kCanceled":                           1,
	"kTimeout":                            2,
	"kTransportError":                     3,
	"kAddressPoolExhausted":               4,
	"kBusFull":                            5,
	"kBusSlotOccupied":                    6,
	"kChecksumMismatch":                   7,
	"kDeviceLocked":                       8,
	"kExists":                             9,
	"kHostEvacuationFailure":              10,
	"kHostRestartAllVmsFailure":           11,
	"kHostRestoreVmLocalityFailure":       12,
	"kHotPlugFailure":                     13,
	"kInUse":                              14,
	"kInvalidArgument":                    15,
	"kInvalidNode":                        16,
	"kInvalidState":                       17,
	"kNotManaged":                         18,
	"kNotMaster":                          19,
	"kVirtualDiskNotFound":                20,
	"kVirtualDiskUpdateFailed":            21,
	"kVmNotFound":                         22,
	"kHostConnectionFailure":              34,
	"kDeviceSlotAlreadyTaken":             42,
	"kClusterReconfigRuntimeError":        43,
	"kGuestToolsNotAvaliableError":        44,
	"kGuestToolsRuntimeError":             45,
	"kAffinityRuleNotFoundError":          46,
	"kCustomizationFaultError":            47,
	"kInvalidDatastoreError":              48,
	"kTaskInProgressError":                49,
	"kVmConfigFaultError":                 50,
	"kFileFaultError":                     51,
	"kMigrationFaultError":                52,
	"kInsufficientResourcesFaultError":    53,
	"kGuestToolsUpgradeFaultError":        59,
	"kInvalidLogin":                       62,
	"kVmUnregisterRuntimeError":           65,
	"kVmAlreadyExistsError":               75,
	"kNetworkError":                       77,
	"kRemoteShellConnectionError":         78,
	"kRemoteShellMaxLimitReachedError":    79,
	"kVmSnapshotNotFoundError":            80,
	"kUuidMappingDuplicateValueError":     82,
	"kStateTransitionTimeoutError":        85,
	"kVmAlreadyUpgradedError":             87,
	"kHypervisorNotSupported":             23,
	"kOperationNotSupported":              24,
	"kRetry":                              25,
	"kUncaughtException":                  26,
	"kNotSupported":                       27,
	"kVmCustomizationIsoCreateFailure":    28,
	"kInvalidOperation":                   29,
	"kNotFound":                           30,
	"kNoHostResources":                    31,
	"kInternalTaskCreation":               32,
	"kDbError":                            33,
	"kVmFlashModeWatcherError":            35,
	"kOperationFailed":                    36,
	"kHyperintClientError":                37,
	"kRpcClientError":                     38,
	"kVmInoperableError":                  39,
	"kDbEntityNotFoundError":              40,
	"kDbNotUpdatedError":                  41,
	"kRestrictedOperationError":           60,
	"kNetworkNotFound":                    81,
	"kFlashModeUpdateError":               88,
	"kManagementServerConnectionFailure":  54,
	"kManagementServerCertificateError":   55,
	"kManagementServerCredentialsError":   56,
	"kManagementServerRegisterError":      57,
	"kManagementServerUnregisterError":    58,
	"kManagementServerNotRegisteredError": 61,
	"kMultipleManagementServersError":     64,
	"kManagementServerConnectionError":    83,
	"kManagementServerNoPermissionError":  84,
	"kOrphanedVmExistsError":              86,
	"kVcenterStaleVmError":                89,
	"kNgtClientError":                     63,
	"kCloudException":                     66,
	"kCloudResourceLimitExceeded":         67,
	"kCloudServiceInternalError":          68,
	"kBootVolumeTypeNotSupported":         69,
	"kCloudServiceUnAvailable":            70,
	"kCloudElasticIPMigrateError":         71,
	"kUnknownVolumeTypeError":             72,
	"kCloudResourceIncorrectStateError":   73,
	"kCloudConnectionFailure":             74,
	"kInstanceTagNotCreatedError":         76,
}

func (x UhuraError_Type) Enum() *UhuraError_Type {
	p := new(UhuraError_Type)
	*p = x
	return p
}

func (x UhuraError_Type) String() string {
	return proto.EnumName(UhuraError_Type_name, int32(x))
}

func (x *UhuraError_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(UhuraError_Type_value, data, "UhuraError_Type")
	if err != nil {
		return err
	}
	*x = UhuraError_Type(value)
	return nil
}

func (UhuraError_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_3a40fcaea3d916ef, []int{0, 0}
}

type UhuraError struct {
	// Required. Type of error encountered.
	ErrorType *UhuraError_Type `protobuf:"varint,1,opt,name=error_type,json=errorType,enum=nutanix.uhura.UhuraError_Type,def=0" json:"error_type,omitempty"`
	// Descriptive string about the error encountered.
	ErrorDetails         *string  `protobuf:"bytes,2,opt,name=error_details,json=errorDetails" json:"error_details,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UhuraError) Reset()         { *m = UhuraError{} }
func (m *UhuraError) String() string { return proto.CompactTextString(m) }
func (*UhuraError) ProtoMessage()    {}
func (*UhuraError) Descriptor() ([]byte, []int) {
	return fileDescriptor_3a40fcaea3d916ef, []int{0}
}

func (m *UhuraError) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UhuraError.Unmarshal(m, b)
}
func (m *UhuraError) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UhuraError.Marshal(b, m, deterministic)
}
func (m *UhuraError) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UhuraError.Merge(m, src)
}
func (m *UhuraError) XXX_Size() int {
	return xxx_messageInfo_UhuraError.Size(m)
}
func (m *UhuraError) XXX_DiscardUnknown() {
	xxx_messageInfo_UhuraError.DiscardUnknown(m)
}

var xxx_messageInfo_UhuraError proto.InternalMessageInfo

const Default_UhuraError_ErrorType UhuraError_Type = UhuraError_kNoError

func (m *UhuraError) GetErrorType() UhuraError_Type {
	if m != nil && m.ErrorType != nil {
		return *m.ErrorType
	}
	return Default_UhuraError_ErrorType
}

func (m *UhuraError) GetErrorDetails() string {
	if m != nil && m.ErrorDetails != nil {
		return *m.ErrorDetails
	}
	return ""
}

func init() {
	proto.RegisterEnum("nutanix.uhura.UhuraError_Type", UhuraError_Type_name, UhuraError_Type_value)
	proto.RegisterType((*UhuraError)(nil), "nutanix.uhura.UhuraError")
}

func init() { proto.RegisterFile("uhura/uhura_error.proto", fileDescriptor_3a40fcaea3d916ef) }

var fileDescriptor_3a40fcaea3d916ef = []byte{
	// 1303 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x56, 0x69, 0x53, 0x1c, 0x39,
	0x12, 0x5d, 0xbc, 0xbe, 0x90, 0xc1, 0x08, 0x61, 0x1b, 0x7c, 0x61, 0x0c, 0xbe, 0x76, 0xbd, 0x8b,
	0x77, 0xbd, 0x3b, 0x97, 0xe7, 0x6c, 0xba, 0x1b, 0xd3, 0x33, 0x74, 0x9b, 0xe9, 0xcb, 0x33, 0xf3,
	0xc5, 0x21, 0xaa, 0x92, 0x6e, 0x45, 0xa9, 0xa4, 0x0a, 0x1d, 0x98, 0x9e, 0x4f, 0xf3, 0xc7, 0x27,
	0x62, 0x42, 0xa9, 0xea, 0x0b, 0xfc, 0x85, 0xa0, 0x33, 0x9f, 0x52, 0xd2, 0xcb, 0x97, 0xaf, 0x44,
	0xd6, 0xfd, 0xd0, 0x1b, 0xfe, 0x0a, 0xff, 0x7e, 0x00, 0x63, 0xb4, 0xd9, 0x2d, 0x8c, 0x76, 0x9a,
	0x2d, 0x2b, 0xef, 0xb8, 0x12, 0x67, 0xbb, 0x98, 0xda, 0xfe, 0x73, 0x8d, 0x90, 0x5e, 0xf8, 0xaf,
	0x1e, 0x30, 0xec, 0x2d, 0x21, 0x08, 0xfe, 0xe0, 0x46, 0x05, 0x6c, 0x2c, 0x6c, 0x2d, 0xbc, 0xb8,
	0xf9, 0x7a, 0x73, 0x77, 0x6e, 0xc9, 0xee, 0x14, 0xbe, 0xdb, 0x1d, 0x15, 0xf0, 0xe6, 0x7a, 0xd6,
	0xd2, 0xf8, 0xb3, 0xbd, 0x88, 0x6b, 0x43, 0x90, 0xed, 0x90, 0xe5, 0x58, 0x28, 0x05, 0xc7, 0x85,
	0xb4, 0x1b, 0x97, 0xb6, 0x16, 0x5e, 0x2c, 0xb6, 0x97, 0x30, 0x58, 0x8b, 0xb1, 0xed, 0x3f, 0xd6,
	0xc8, 0x65, 0x44, 0x2f, 0x91, 0x49, 0x11, 0xfa, 0x37, 0xb6, 0x4c, 0x16, 0xb3, 0x2a, 0x57, 0x09,
	0x48, 0x48, 0xe9, 0x02, 0x26, 0xbb, 0x22, 0x07, 0xed, 0x1d, 0xbd, 0xc4, 0xd6, 0xc8, 0x4a, 0xd6,
	0x35, 0x5c, 0xd9, 0x42, 0x1b, 0x17, 0x57, 0xfc, 0x9d, 0xdd, 0x25, 0xb7, 0xb3, 0x4a, 0x9a, 0x1a,
	0xb0, 0xf6, 0x48, 0x6b, 0x59, 0x3f, 0x1b, 0x72, 0x6f, 0x1d, 0xa4, 0xf4, 0x32, 0xae, 0xde, 0xf3,
	0x76, 0xdf, 0x4b, 0x49, 0xaf, 0xb0, 0x5b, 0x84, 0x86, 0x5f, 0x1d, 0xa9, 0xdd, 0xbb, 0x24, 0xf1,
	0x85, 0x80, 0x94, 0x5e, 0x65, 0xb7, 0xc9, 0x6a, 0x56, 0x1d, 0x42, 0x92, 0x59, 0x9f, 0x37, 0x85,
	0xcd, 0xb9, 0x4b, 0x86, 0xf4, 0x1a, 0x5b, 0x25, 0xcb, 0x59, 0x0d, 0x4e, 0x45, 0x02, 0x87, 0x3a,
	0xc9, 0x20, 0xa5, 0xd7, 0xd9, 0x0d, 0x72, 0x2d, 0xab, 0x9f, 0x09, 0xeb, 0x2c, 0x5d, 0x64, 0xf7,
	0xc8, 0x9d, 0xec, 0x40, 0x5b, 0x57, 0x3f, 0xe5, 0x89, 0xe7, 0x4e, 0x68, 0xb5, 0xcf, 0x85, 0xf4,
	0x06, 0x28, 0x61, 0x0f, 0xc9, 0x5d, 0xcc, 0xb5, 0xc1, 0x3a, 0x6e, 0x5c, 0x45, 0xca, 0x7e, 0x6e,
	0xc7, 0xe9, 0x1b, 0xec, 0x31, 0x79, 0x38, 0x49, 0x6b, 0x03, 0xfd, 0xfc, 0x50, 0x27, 0x5c, 0x0a,
	0x37, 0x1a, 0x43, 0x96, 0xf0, 0xa2, 0x07, 0xda, 0x1d, 0x49, 0x3f, 0x18, 0x07, 0x97, 0x19, 0x21,
	0x57, 0xb3, 0x86, 0xea, 0x59, 0xa0, 0x37, 0xf1, 0x2e, 0x0d, 0x75, 0xca, 0xa5, 0x48, 0x2b, 0x66,
	0xe0, 0x73, 0x50, 0x8e, 0xae, 0x30, 0x4a, 0x96, 0xc6, 0xd1, 0x96, 0x4e, 0x81, 0x52, 0xbc, 0x46,
	0x19, 0xe9, 0x38, 0xee, 0x80, 0xae, 0xb2, 0x15, 0x72, 0x23, 0x6b, 0x69, 0xd7, 0xe4, 0x8a, 0x0f,
	0x20, 0xa5, 0x8c, 0xdd, 0x24, 0x24, 0x06, 0xac, 0x03, 0x43, 0xd7, 0xd8, 0x06, 0xb9, 0x95, 0xf5,
	0x85, 0x71, 0x9e, 0xcb, 0x9a, 0xb0, 0x21, 0xb5, 0xaf, 0xbd, 0x4a, 0xe9, 0x2d, 0xf6, 0x80, 0x6c,
	0xcc, 0x66, 0x7a, 0x45, 0xca, 0x1d, 0x84, 0x03, 0x42, 0x4a, 0x6f, 0x63, 0xe1, 0x7e, 0x3e, 0x81,
	0xdf, 0x99, 0x70, 0x54, 0xd5, 0x4a, 0x41, 0x32, 0xcb, 0xd1, 0x36, 0xbb, 0x4f, 0xd6, 0x4b, 0x7e,
	0x43, 0x3f, 0x2a, 0xd2, 0x00, 0x4f, 0x47, 0x5d, 0x9e, 0x81, 0xa2, 0xff, 0x64, 0x5b, 0xe4, 0x41,
	0x56, 0x95, 0xa1, 0x8b, 0xa6, 0x0d, 0x89, 0x56, 0x27, 0x62, 0xd0, 0xf6, 0xca, 0x89, 0x1c, 0x62,
	0xd3, 0x5f, 0x22, 0xe2, 0xad, 0x07, 0xeb, 0xba, 0x5a, 0x4b, 0xdb, 0xd2, 0xae, 0x12, 0xee, 0xc8,
	0x8f, 0x65, 0x89, 0xf8, 0x17, 0x6e, 0x30, 0x45, 0xcc, 0x2d, 0xff, 0x37, 0xdb, 0x24, 0xf7, 0xb2,
	0xca, 0xc9, 0x89, 0x50, 0xc2, 0x8d, 0xda, 0x5e, 0xc2, 0xf8, 0xd0, 0x31, 0xbf, 0x8b, 0x17, 0xad,
	0x7a, 0xeb, 0x74, 0x2e, 0x7e, 0x2f, 0x9b, 0xeb, 0x65, 0xa9, 0xb8, 0x57, 0x78, 0xaf, 0x92, 0xd4,
	0x1a, 0x77, 0x1c, 0xbb, 0x18, 0x73, 0xff, 0x41, 0xf2, 0xba, 0xdc, 0x66, 0x0d, 0x75, 0x64, 0xf4,
	0x20, 0x88, 0x32, 0x66, 0xfe, 0xcb, 0xd6, 0xc9, 0x5a, 0xd6, 0xcf, 0xab, 0x78, 0x9b, 0x99, 0x72,
	0xaf, 0xb1, 0xd9, 0xfb, 0x42, 0xc2, 0x4c, 0xf0, 0x7f, 0x58, 0xa7, 0x29, 0x06, 0xe6, 0xfc, 0xee,
	0xff, 0x67, 0x4f, 0xc8, 0x56, 0xd6, 0x50, 0xd6, 0x9f, 0x9c, 0x88, 0x44, 0x80, 0x0a, 0x32, 0xd2,
	0xde, 0x24, 0x60, 0x67, 0x50, 0x9f, 0x9d, 0x23, 0xa8, 0x57, 0x0c, 0x0c, 0x4f, 0x67, 0x77, 0xf8,
	0x7a, 0x56, 0x1a, 0x87, 0x7a, 0x20, 0x14, 0xfd, 0x0e, 0x85, 0xdb, 0xcf, 0x7b, 0xca, 0xc0, 0x40,
	0x20, 0xf9, 0xb3, 0xac, 0x55, 0x70, 0xd2, 0xfa, 0x79, 0xd9, 0xab, 0x38, 0x09, 0x31, 0xf5, 0x13,
	0x16, 0x6b, 0x81, 0xfb, 0xa8, 0x4d, 0x16, 0x43, 0x4d, 0xf6, 0x88, 0xdc, 0xcf, 0xda, 0x90, 0x6b,
	0x07, 0x9d, 0x21, 0x48, 0x39, 0x15, 0x41, 0x04, 0xb4, 0xf0, 0x22, 0x33, 0x80, 0x26, 0x3f, 0x3b,
	0x14, 0xb9, 0x70, 0x6d, 0xe0, 0xc9, 0x10, 0xca, 0x56, 0xbc, 0x8b, 0x9a, 0xcb, 0x3b, 0x8a, 0x17,
	0x76, 0xa8, 0xdd, 0x7c, 0xa3, 0x8e, 0xd8, 0x0e, 0x79, 0x94, 0xf5, 0xbc, 0x48, 0x9b, 0xbc, 0x28,
	0x84, 0x1a, 0xd4, 0x7c, 0x21, 0x45, 0xc2, 0x1d, 0xf4, 0xb9, 0xf4, 0xe5, 0xb9, 0xdb, 0xc8, 0x05,
	0xaa, 0x1f, 0xbd, 0x43, 0x84, 0x33, 0x94, 0x9e, 0x12, 0x11, 0x3d, 0x14, 0xcb, 0xe4, 0x66, 0x25,
	0x59, 0xe5, 0x1e, 0xef, 0x31, 0x79, 0x30, 0x2a, 0xc0, 0x9c, 0x0a, 0xab, 0x4d, 0x4b, 0xbb, 0x8e,
	0x2f, 0x82, 0x01, 0x41, 0x4a, 0xd7, 0x51, 0x0b, 0xef, 0x0a, 0x88, 0x7d, 0x9a, 0xcb, 0x6d, 0xe0,
	0xc0, 0xb6, 0xc1, 0x99, 0x11, 0xbd, 0xcb, 0xee, 0x10, 0x96, 0xf5, 0x54, 0xc2, 0xfd, 0x60, 0xe8,
	0xea, 0x67, 0x09, 0x14, 0x01, 0x4f, 0xef, 0x45, 0xe2, 0x66, 0x97, 0xdd, 0x47, 0x5e, 0xfa, 0xf9,
	0x9c, 0xfc, 0x1a, 0x56, 0x57, 0x0d, 0x94, 0xc3, 0x16, 0x06, 0xe8, 0x01, 0xfa, 0x56, 0xd9, 0xbe,
	0xc9, 0xfe, 0xf4, 0x21, 0xfa, 0xe7, 0x64, 0x04, 0x37, 0xd1, 0x27, 0x5a, 0xba, 0x74, 0x9b, 0x28,
	0x13, 0xfa, 0x08, 0x1b, 0xd9, 0x50, 0x0e, 0x8c, 0xe2, 0x32, 0x88, 0x15, 0x6b, 0x87, 0xf5, 0x5b,
	0x68, 0x99, 0xb5, 0xe3, 0x78, 0xf5, 0xc7, 0x25, 0xf9, 0xfb, 0x92, 0xdb, 0x61, 0x53, 0xa7, 0xf0,
	0x3e, 0x98, 0x23, 0x98, 0x98, 0xdd, 0xc1, 0xe2, 0x93, 0xbd, 0x4b, 0x1b, 0x78, 0x82, 0xca, 0x45,
	0xba, 0x84, 0x72, 0x55, 0x19, 0x14, 0x1a, 0xf1, 0x4f, 0x51, 0xe8, 0xed, 0x22, 0x99, 0x0d, 0x3e,
	0x43, 0x62, 0xfa, 0x79, 0x43, 0xe9, 0x50, 0x68, 0x32, 0xbf, 0xcf, 0x91, 0xd8, 0xda, 0x71, 0x5d,
	0x39, 0xe1, 0x46, 0xf3, 0x5d, 0x7f, 0x81, 0x6b, 0x6a, 0xc7, 0x2d, 0xed, 0xa2, 0x03, 0x95, 0xf1,
	0x7f, 0xa0, 0x7e, 0x83, 0xab, 0x1a, 0x91, 0x38, 0x98, 0xd2, 0x12, 0xd3, 0xdf, 0x44, 0x32, 0xa2,
	0x48, 0x27, 0x14, 0xfd, 0x8c, 0x64, 0x4c, 0x6e, 0x18, 0x0b, 0xc6, 0x05, 0xbf, 0xb0, 0x67, 0x64,
	0x3b, 0x8b, 0x3e, 0x19, 0x0c, 0xb6, 0x03, 0xe6, 0x14, 0xcc, 0x45, 0x33, 0xfb, 0x9c, 0x3d, 0x25,
	0x8f, 0x2f, 0xe2, 0xc0, 0x38, 0x71, 0x82, 0x62, 0x8c, 0xe5, 0xbe, 0xf8, 0x34, 0xcc, 0x40, 0x0a,
	0xca, 0x09, 0x2e, 0xcb, 0x59, 0xfa, 0x92, 0x6d, 0x93, 0xcd, 0x0b, 0xb0, 0x76, 0x39, 0x91, 0x11,
	0xf3, 0x15, 0x6a, 0xe4, 0x3c, 0x66, 0x3a, 0xb7, 0x11, 0xf5, 0x86, 0x3d, 0x27, 0x3b, 0x17, 0x50,
	0x2d, 0xed, 0xc6, 0xc5, 0xc6, 0xc4, 0x7d, 0x8b, 0x63, 0xd4, 0xf4, 0xd2, 0x89, 0x42, 0xc2, 0xf9,
	0x05, 0xe5, 0xb9, 0x7e, 0xf8, 0xe4, 0x9e, 0xe7, 0xa7, 0xba, 0xf3, 0x49, 0xce, 0x5a, 0xfa, 0x08,
	0x4c, 0x2e, 0xac, 0x9d, 0xe0, 0xba, 0x71, 0x70, 0x4c, 0x31, 0xe4, 0x0a, 0xd2, 0x7e, 0x3e, 0xeb,
	0x26, 0xfd, 0xf8, 0x05, 0x4a, 0x20, 0x28, 0xb4, 0xe3, 0xb8, 0x84, 0x7e, 0x1e, 0x33, 0xbf, 0xa2,
	0x84, 0x5a, 0x83, 0x39, 0x5d, 0x7d, 0x8f, 0xc1, 0xaa, 0xd4, 0x3e, 0x9d, 0x0e, 0xd6, 0x1e, 0xda,
	0x0f, 0x06, 0xc7, 0xc2, 0x47, 0x73, 0x09, 0x08, 0x48, 0x21, 0xa5, 0x55, 0xfc, 0x06, 0x20, 0x20,
	0x9c, 0x51, 0x24, 0x30, 0x1e, 0x88, 0x58, 0xb5, 0x86, 0x05, 0xf6, 0xb4, 0x76, 0x7d, 0x2d, 0x7d,
	0x0e, 0xe1, 0xa9, 0x32, 0x37, 0xa7, 0xf5, 0xf8, 0x91, 0x98, 0x29, 0xd0, 0x53, 0x95, 0x53, 0x2e,
	0x64, 0x10, 0x31, 0xdd, 0x9f, 0xee, 0x5f, 0x97, 0xdc, 0x3a, 0x91, 0x34, 0x8e, 0xa2, 0x9d, 0x97,
	0x6a, 0x78, 0x8b, 0xb6, 0xd2, 0x53, 0x99, 0xd2, 0x1f, 0xd5, 0x74, 0x8b, 0x98, 0x3c, 0x40, 0xa9,
	0xcc, 0x9d, 0xbe, 0xa1, 0x12, 0x6d, 0x0c, 0x24, 0x0e, 0x8d, 0x2c, 0xc2, 0x1a, 0x58, 0x03, 0x61,
	0x17, 0x55, 0xf9, 0x23, 0x9e, 0xa0, 0xa1, 0xac, 0x0b, 0xaf, 0xa9, 0x2e, 0x1f, 0xb4, 0xb4, 0x8b,
	0x1e, 0x52, 0x76, 0xfd, 0x70, 0xef, 0x25, 0x59, 0x4d, 0x74, 0x3e, 0xff, 0xc2, 0xdb, 0x5b, 0x99,
	0x3e, 0xf1, 0x8e, 0xc2, 0xa3, 0xf1, 0xb7, 0x2b, 0x18, 0xff, 0x2b, 0x00, 0x00, 0xff, 0xff, 0xd5,
	0x79, 0x3b, 0xce, 0x55, 0x0a, 0x00, 0x00,
}