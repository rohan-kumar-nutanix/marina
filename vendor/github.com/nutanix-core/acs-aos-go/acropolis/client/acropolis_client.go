/*
 * Copyrigright (c) 2016 Nutanix Inc. All rights reserved.
 *
 * Defines client interface for Acropolis service
 *
 */

package acropolis_client

import (
	"bytes"
	"errors"
	"time"

	glog "github.com/golang/glog"
	proto "github.com/golang/protobuf/proto"

	acropolis_ifc "github.com/nutanix-core/acs-aos-go/acropolis"
	pithos_ifc "github.com/nutanix-core/acs-aos-go/pithos"

	ntnx_errors "github.com/nutanix-core/acs-aos-go/nutanix/util-go/errors"
	util_misc "github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
	util_net "github.com/nutanix-core/acs-aos-go/nutanix/util-go/net"

	util_base "github.com/nutanix-core/acs-aos-go/nutanix/util-go/base"
)

type FuncArguments map[string]interface{}

const (
	acropolisServiceName             = "nutanix.acropolis.AcropolisRpcSvc"
	defaultAcropolisAddr             = "127.0.0.1"
	defaultAcropolisPort             = 2030
	clientRetryInitialDelayMilliSecs = 250
	clientRetryMaxDelaySecs          = 2
	clientRetryTimeoutSecs           = 30
)

// Acropolis Error defn.
type AcropolisError_ struct {
	*ntnx_errors.NtnxError
}

func (e *AcropolisError_) TypeOfError() int {
	return ntnx_errors.AcropolisErrorType
}

func (e *AcropolisError_) SetCause(err error) error {
	return ntnx_errors.NewNtnxErrorRef(err, e.GetErrorDetail(),
		e.GetErrorCode(), e)
}

func (e *AcropolisError_) Equals(err interface{}) bool {
	if obj, ok := err.(*AcropolisError_); ok {
		return e == obj
	}
	return false
}

func AcropolisError(errMsg string, errCode int) *AcropolisError_ {
	return &AcropolisError_{ntnx_errors.NewNtnxError(errMsg, errCode)}
}

var (
	ErrNoError                                      = AcropolisError("NoError", 0)
	ErrCanceled                                     = AcropolisError("Canceled", 1)
	ErrRetry                                        = AcropolisError("Retry", 2)
	ErrTimeout                                      = AcropolisError("Timeout", 3)
	ErrNotSupported                                 = AcropolisError("NotSupported", 4)
	ErrUncaughtException                            = AcropolisError("UncaughtException", 5)
	ErrInvalidArgument                              = AcropolisError("InvalidArgument", 6)
	ErrInvalidState                                 = AcropolisError("InvalidState", 7)
	ErrLogicalTimestampMismatch                     = AcropolisError("LogicalTimestampMismatch", 8)
	ErrExists                                       = AcropolisError("Exists", 9)
	ErrNotFound                                     = AcropolisError("NotFound", 10)
	ErrDeviceLocked                                 = AcropolisError("DeviceLocked", 11)
	ErrBusFull                                      = AcropolisError("BusFull", 12)
	ErrBusSlotOccupied                              = AcropolisError("BusSlotOccupied", 13)
	ErrNoHostResources                              = AcropolisError("NoHostResources", 14)
	ErrNotManaged                                   = AcropolisError("NotManaged", 15)
	ErrNotMaster                                    = AcropolisError("NotMaster", 16)
	ErrMigrationFail                                = AcropolisError("MigrationFail", 17)
	ErrInUse                                        = AcropolisError("InUse", 18)
	ErrAddressPoolExhausted                         = AcropolisError("AddressPoolExhausted", 19)
	ErrHotPlugFailure                               = AcropolisError("HotPlugFailure", 20)
	ErrHostEvacuationFailure                        = AcropolisError("HostEvacuationFailure", 21)
	ErrUploadFailure                                = AcropolisError("UploadFailure", 22)
	ErrChecksumMismatch                             = AcropolisError("ChecksumMismatch", 23)
	ErrHostRestartAllVmsFailure                     = AcropolisError("HostRestartAllVmsFailure", 24)
	ErrHostRestoreVmLocalityFailure                 = AcropolisError("HostRestoreVmLocalityFailure", 25)
	ErrTransportError                               = AcropolisError("TransportError", 26)
	ErrHypervisorConnection                         = AcropolisError("HypervisorConnectionError", 27)
	ErrNetworkError                                 = AcropolisError("NetworkError", 28)
	ErrVmCloneError                                 = AcropolisError("VmCloneError", 29)
	ErrVgCloneError                                 = AcropolisError("VgCloneError", 30)
	ErrDelayRetry                                   = AcropolisError("DelayRetry", 31)
	ErrVNumaPinningFailure                          = AcropolisError("VNumaPinningFailure", 32)
	ErrUpgradeFailure                               = AcropolisError("UpgradeFailure", 33)
	ErrCopyFileFailure                              = AcropolisError("CopyFileFailure", 34)
	ErrNGTError                                     = AcropolisError("NGTError", 35)
	ErrInternal                                     = AcropolisError("Internal", 36)
	ErrServiceError                                 = AcropolisError("ServiceError", 37)
	ErrVmNotFound                                   = AcropolisError("VmNotFound", 38)
	ErrTaskNotFound                                 = AcropolisError("TaskNotFound", 39)
	ErrSchedulerError                               = AcropolisError("SchedulerError", 40)
	ErrSchedulerMissingEntity                       = AcropolisError("SchedulerMissingEntity", 41)
	ErrImportVmError                                = AcropolisError("ImportVmError", 42)
	ErrRemoteShell                                  = AcropolisError("RemoteShell", 43)
	ErrVNumaPowerOn                                 = AcropolisError("VNumaPowerOn", 44)
	ErrVNumaMigrate                                 = AcropolisError("VNumaMigrate", 45)
	ErrCasError                                     = AcropolisError("CasError", 46)
	ErrLazanDelayRetry                              = AcropolisError("LazanDelayRetry", 48)
	ErrImportError                                  = AcropolisError("ImportError", 49)
	ErrInvalidVmState                               = AcropolisError("InvalidVmState", 50)
	ErrNetworkNotFound                              = AcropolisError("NetworkNotFound", 51)
	ErrMicrosegRuleError                            = AcropolisError("MicrosegRuleError", 52)
	ErrCerebroError                                 = AcropolisError("CerebroError", 54)
	ErrVmSyncRepError                               = AcropolisError("VmSyncRepError", 55)
	ErrVmSyncRepCannotEnable                        = AcropolisError("VmSyncRepCannotEnable", 56)
	ErrVmSyncRepUpdateDormantVm                     = AcropolisError("VmSyncRepUpdateDormantVm", 57)
	ErrVmSyncRepFailedDiskProtectionError           = AcropolisError("VmSyncRepFailedDiskProtectionError", 58)
	ErrVmSyncRepFailedToGetVmSpecError              = AcropolisError("VmSyncRepFailedToGetVmSpecError", 59)
	ErrHypervisorError                              = AcropolisError("HypervisorError", 61)
	ErrVmSyncRepMigrateError                        = AcropolisError("VmSyncRepMigrateError", 62)
	ErrVmCrossClusterLiveMigrateError               = AcropolisError("VmCrossClusterLiveMigrateError", 64)
	ErrVmCrossClusterLiveMigratePrepareCleanupError = AcropolisError("VmCrossClusterLiveMigratePrepareCleanupError", 65)
	ErrUncommittedDiskError                         = AcropolisError("UncommittedDiskError", 66)
	ErrDeleteFileFailure                            = AcropolisError("DeleteFileFailure", 67)
	ErrVmDisableUpdateError                         = AcropolisError("VmDisableUpdateError", 68)
	ErrVmCrossClusterLiveMigratePrechecksError      = AcropolisError("VmCrossClusterLiveMigratePrechecksError", 69)
	ErrVmSyncRepInvalidStretchParamsError           = AcropolisError("VmSyncRepInvalidStretchParamsError", 70)
	ErrVmCrossClusterLiveMigrateUnknownError        = AcropolisError("VmCrossClusterLiveMigrateUnknownError", 71)
	ErrVirtualSwitchError                           = AcropolisError("VirtualSwitchError", 72)
)

var (
	ErrInvalidResponse = AcropolisError("Invalid RPC Response", 101)
	ErrInternalError   = AcropolisError("AcropolisInternalError", 102)
)

//Acropolis unmapped exception error
var (
	ErrUnmappedResponseCode = AcropolisError("ErrUnmappedResponseCode", 1001)
)

var Errors = map[int]*AcropolisError_{
	0:  ErrNoError,
	1:  ErrCanceled,
	2:  ErrRetry,
	3:  ErrTimeout,
	4:  ErrNotSupported,
	5:  ErrUncaughtException,
	6:  ErrInvalidArgument,
	7:  ErrInvalidState,
	8:  ErrLogicalTimestampMismatch,
	9:  ErrExists,
	10: ErrNotFound,
	11: ErrDeviceLocked,
	12: ErrBusFull,
	13: ErrBusSlotOccupied,
	14: ErrNoHostResources,
	15: ErrNotManaged,
	16: ErrNotMaster,
	17: ErrMigrationFail,
	18: ErrInUse,
	19: ErrAddressPoolExhausted,
	20: ErrHotPlugFailure,
	21: ErrHostEvacuationFailure,
	22: ErrUploadFailure,
	23: ErrChecksumMismatch,
	24: ErrHostRestartAllVmsFailure,
	25: ErrHostRestoreVmLocalityFailure,
	26: ErrTransportError,
	27: ErrHypervisorConnection,
	28: ErrNetworkError,
	29: ErrVmCloneError,
	30: ErrVgCloneError,
	31: ErrDelayRetry,
	32: ErrVNumaPinningFailure,
	33: ErrUpgradeFailure,
	34: ErrCopyFileFailure,
	35: ErrNGTError,
	36: ErrInternal,
	37: ErrServiceError,
	38: ErrVmNotFound,
	39: ErrTaskNotFound,
	40: ErrSchedulerError,
	41: ErrSchedulerMissingEntity,
	42: ErrImportVmError,
	43: ErrRemoteShell,
	44: ErrVNumaPowerOn,
	45: ErrVNumaMigrate,
	46: ErrCasError,
	48: ErrLazanDelayRetry,
	49: ErrImportError,
	50: ErrInvalidVmState,
	51: ErrNetworkNotFound,
	52: ErrMicrosegRuleError,
	54: ErrCerebroError,
	55: ErrVmSyncRepError,
	56: ErrVmSyncRepCannotEnable,
	57: ErrVmSyncRepUpdateDormantVm,
	58: ErrVmSyncRepFailedDiskProtectionError,
	59: ErrVmSyncRepFailedToGetVmSpecError,
	61: ErrHypervisorError,
	62: ErrVmSyncRepMigrateError,
	64: ErrVmCrossClusterLiveMigrateError,
	65: ErrVmCrossClusterLiveMigratePrepareCleanupError,
	66: ErrUncommittedDiskError,
	67: ErrDeleteFileFailure,
	68: ErrVmDisableUpdateError,
	69: ErrVmCrossClusterLiveMigratePrechecksError,
	70: ErrVmSyncRepInvalidStretchParamsError,
	71: ErrVmCrossClusterLiveMigrateUnknownError,
	72: ErrVirtualSwitchError,
}

var NewProtobufRPCClient = util_net.NewProtobufRPCClientIfc

const (
	DiskClone    = 0
	DiskCreate   = 1
	DiskExisting = 2
)

type Acropolis struct {
	client util_net.ProtobufRPCClientIfc
}

type AcropolisClientInterface interface {
	GetAllImages(includeDiskSize bool, includeDiskPath bool) (
		[]*acropolis_ifc.ImageInfo, error)
	GetImages(uuidList *[][]byte, includeDiskSize bool, includeDiskPath bool) (
		[]*acropolis_ifc.ImageInfo, error)
	CreateDiskImageFromDisk(name, annotation string, srcVmdiskUuid []byte,
		dstStorageUuid []byte, parentTask []byte) ([]byte, error)
	UpdateVmDisk(uuid []byte, diskSpec *acropolis_ifc.VmDiskUpdateArg_VmDiskUpdateSpec,
		parentUuid []byte) ([]byte, error)
	UpdateImage(uuid []byte, name, annotation string, parentTaskUuid []byte) (
		[]byte, error)
	DeleteImage(uuid []byte, parentTaskUuid []byte) ([]byte, error)
	GetNetworks(uuidList [][]byte) ([]*acropolis_ifc.NetworkConfig, error)
	GetVmsOnHost(uuid []byte) ([][]byte, error)
	VmNicGet(vmUuid []byte, macAddrList [][]byte) ([]*acropolis_ifc.VmNicConfig, error)
	GetNetworkAddressTable(networkUuid []byte) (*acropolis_ifc.NetworkAddressAssignmentTable, error)
	ReserveIp(uuid []byte, numIpAddress uint64, cookie string) ([][]byte, error)
	UnreserveIp(uuid []byte, ipAddrList [][]byte, cookie string) error
}

func DefaultAcropolisService() *Acropolis {
	return NewAcropolisService(defaultAcropolisAddr, defaultAcropolisPort)
}

func NewAcropolisService(serviceIP string, servicePort uint16) *Acropolis {
	glog.Infoln("Initializing Acropolis Service")
	return &Acropolis{
		client: NewProtobufRPCClient(serviceIP, servicePort),
	}
}

func (svc *Acropolis) sendMsg(service string, request, response proto.Message) error {

	retryWait := util_misc.NewExponentialBackoff(1*time.Second, 30*time.Second,
		10)
	var done bool = false
	for !done {
		err := svc.client.CallMethodSync(acropolisServiceName, service, request,
			response, 0)
		if err != nil {
			// If App error is set in RPC - extract the code and get the relevant
			// acropolis error.

			if obj, ok := ntnx_errors.TypeAssert(err, ntnx_errors.AppErrorType); ok {
				if errObj, ok := obj.(*util_net.AppError_); ok {
					errCode := errObj.GetErrorCode()
					if val, ok := Errors[errCode]; ok {
						return val
					}
					if bd := util_base.BuildInfo(util_base.GetBuildVersion); bd.IsDebugBuild() {
						glog.Fatalf("Unmapped response code: %d", errCode)
					} else {
						glog.Errorf("Unmapped response code: %d", errCode)
						return ErrUnmappedResponseCode
					}
				} else {
					return ErrInternalError
				}
			}
			if obj, ok := ntnx_errors.TypeAssert(err, ntnx_errors.RpcErrorType); ok {
				if rpcErr, ok := obj.(*util_net.RpcError_); ok {
					errCode := rpcErr.GetErrorCode()
					glog.Errorf("RpcError in %s: %d", service, errCode)
					if util_net.ErrRpcTransport.Equals(rpcErr) {
						glog.Infof("Failure to send %s, will retry later", service)
						waited := retryWait.Backoff()
						if waited != util_misc.Stop {
							glog.Infof("Retrying, msg: %s", service)
							continue
						} else {
							done = true
						}
					}
				}
			}
		}
		return err
	}
	return nil
}

func (svc *Acropolis) CreateVmSpec(name string, uuid []byte, numVcpus uint64,
	numCoresPerVcpu int64, memMb uint64, descr string) (
	*acropolis_ifc.VmCreateSpec, error) {
	//
	spec := &acropolis_ifc.VmCreateSpec{
		Name:            proto.String(name),
		Uuid:            uuid,
		Annotation:      proto.String(descr),
		NumVcpus:        proto.Uint64(numVcpus),
		MemoryMb:        proto.Uint64(memMb),
		NumCoresPerVcpu: proto.Int64(numCoresPerVcpu),
	}

	return spec, nil
}

// Get all networks
// If network uuid list is not specified, it will return info for all networks.
func (svc *Acropolis) GetNetworks(uuidList [][]byte) ([]*acropolis_ifc.NetworkConfig, error) {

	if len(uuidList) == 0 {
		listArg := &acropolis_ifc.NetworkListArg{}
		listRet := &acropolis_ifc.NetworkListRet{}

		err := svc.sendMsg("NetworkList", listArg, listRet)

		if err != nil {
			return nil, err
		}
		uuidList = listRet.GetNetworkUuidList()
		if len(uuidList) == 0 {
			return nil, nil
		}
	}
	networkArg := &acropolis_ifc.NetworkGetArg{
		NetworkUuidList: uuidList,
	}
	networkRet := &acropolis_ifc.NetworkGetRet{}
	err := svc.sendMsg("NetworkGet", networkArg, networkRet)
	if err != nil {
		return nil, err
	}
	networkConfigs := networkRet.GetNetworkConfigList()

	return networkConfigs, nil
}

func (svc *Acropolis) VmNicSpec(macAddr []byte, netUuid []byte,
	model string, ipAddr []byte, vlanType *acropolis_ifc.VmNicVlanType,
	trunkedNetworks []int32, requestIp bool) (*acropolis_ifc.VmNicSpec, error) {
	//
	spec := &acropolis_ifc.VmNicSpec{
		MacAddr:            macAddr,
		NetworkUuid:        netUuid,
		Model:              proto.String(model),
		RequestedIpAddress: ipAddr,
		VlanMode:           vlanType,
		TrunkedNetworks:    trunkedNetworks,
		RequestIp:          proto.Bool(requestIp),
	}

	return spec, nil
}

func (svc *Acropolis) VmDiskSpecCloneFromVmDisk(srcUuid []byte, ctrId int64,
	ctrUuid []byte, isCdrom bool, bus string, index uint32) (*acropolis_ifc.VmDiskSpec, error) {

	spec := &acropolis_ifc.VmDiskSpecClone{
		VmdiskUuid:  srcUuid,
		ContainerId: proto.Int64(ctrId),
	}
	if len(ctrUuid) > 0 {
		spec.ContainerUuid = ctrUuid
	}
	diskSpec := &acropolis_ifc.VmDiskSpec{
		Addr: &acropolis_ifc.VmDiskAddr{
			Bus:   proto.String(bus),
			Index: proto.Uint32(index),
		},
		Cdrom: proto.Bool(isCdrom),
		Clone: spec,
	}

	return diskSpec, nil
}

func (svc *Acropolis) VmDiskSpecCloneFromNfs(nfsPath string, ctrId int64,
	ctrUuid []byte, isCdrom bool, bus string, index uint32) (*acropolis_ifc.VmDiskSpec, error) {

	spec := &acropolis_ifc.VmDiskSpecClone{
		NfsPath:       proto.String(nfsPath),
		ContainerId:   proto.Int64(ctrId),
		ContainerUuid: ctrUuid,
	}

	diskSpec := &acropolis_ifc.VmDiskSpec{
		Addr: &acropolis_ifc.VmDiskAddr{
			Bus:   proto.String(bus),
			Index: proto.Uint32(index),
		},
		Cdrom: proto.Bool(isCdrom),
		Clone: spec,
	}

	return diskSpec, nil
}

func (svc *Acropolis) VmDiskSpecCreateId(size uint64, containerId uint64,
	isCdrom bool, bus string, index uint32) (*acropolis_ifc.VmDiskSpec, error) {

	spec := &acropolis_ifc.VmDiskSpecCreate{
		Size:        proto.Uint64(size),
		ContainerId: proto.Uint64(containerId),
	}

	diskSpec := &acropolis_ifc.VmDiskSpec{
		Addr: &acropolis_ifc.VmDiskAddr{
			Bus:   proto.String(bus),
			Index: proto.Uint32(index),
		},
		Cdrom:  proto.Bool(isCdrom),
		Create: spec,
	}

	return diskSpec, nil
}

func (svc *Acropolis) VmDiskSpecCreateName(size uint64, containerName string,
	isCdrom bool, bus string, index uint32) (*acropolis_ifc.VmDiskSpec, error) {

	spec := &acropolis_ifc.VmDiskSpecCreate{
		Size:          proto.Uint64(size),
		ContainerName: proto.String(containerName),
	}
	diskSpec := &acropolis_ifc.VmDiskSpec{
		Addr: &acropolis_ifc.VmDiskAddr{
			Bus:   proto.String(bus),
			Index: proto.Uint32(index),
		},
		Cdrom:  proto.Bool(isCdrom),
		Create: spec,
	}

	return diskSpec, nil
}

func (svc *Acropolis) VmDiskSpecCreateUuid(size uint64, containerUuid []byte,
	isCdrom bool, bus string, index uint32) (*acropolis_ifc.VmDiskSpec, error) {

	spec := &acropolis_ifc.VmDiskSpecCreate{
		Size:          proto.Uint64(size),
		ContainerUuid: containerUuid,
	}
	diskSpec := &acropolis_ifc.VmDiskSpec{
		Addr: &acropolis_ifc.VmDiskAddr{
			Bus:   proto.String(bus),
			Index: proto.Uint32(index),
		},
		Cdrom:  proto.Bool(isCdrom),
		Create: spec,
	}
	return diskSpec, nil
}

func (svc *Acropolis) CreateVmBootConfig(bus string, index uint32) (
	*acropolis_ifc.VmBootConfig, error) {

	bootConfig := &acropolis_ifc.VmBootConfig{
		Device: &acropolis_ifc.VmBootConfig_VmBootDevice{
			DiskAddr: &acropolis_ifc.VmDiskAddr{
				Bus:   proto.String(bus),
				Index: proto.Uint32(index),
			},
		},
	}
	return bootConfig, nil
}

func (svc *Acropolis) CreateVmBootConfigSpec(bus string, index uint32) (
	*acropolis_ifc.VmBootConfigSpec, error) {

	bootConfig := &acropolis_ifc.VmBootConfigSpec{
		Device: &acropolis_ifc.VmBootConfigSpec_VmBootDeviceSpec{
			DiskAddr: &acropolis_ifc.VmDiskAddr{
				Bus:   proto.String(bus),
				Index: proto.Uint32(index),
			},
		},
	}
	return bootConfig, nil
}

func (svc *Acropolis) VmCreate(vmSpec *acropolis_ifc.VmCreateSpec,
	taskUuid, parentTaskUuid []byte) ([]byte, error) {

	vmArg := &acropolis_ifc.VmCreateArg{
		Spec:           vmSpec,
		TaskUuid:       taskUuid,
		ParentTaskUuid: parentTaskUuid,
	}
	vmRet := &acropolis_ifc.VmCreateRet{}

	err := svc.sendMsg("VmCreate", vmArg, vmRet)
	if err != nil {
		return nil, err
	}

	return vmRet.GetTaskUuid(), nil
}

func (svc *Acropolis) VmUpdateBootDevice(vmUuid []byte,
	bootConfig *acropolis_ifc.VmBootConfig) (
	[]byte, error) {

	vmArg := &acropolis_ifc.VmUpdateArg{
		VmUuid: vmUuid,
		Boot:   bootConfig,
	}
	vmRet := &acropolis_ifc.VmCreateRet{}

	err := svc.sendMsg("VmUpdate", vmArg, vmRet)
	if err != nil {
		return nil, err
	}

	return vmRet.GetTaskUuid(), nil
}

// add a Nic from the network to the specified VM.
func (svc *Acropolis) VmNicAddFromNetwork(vm_uuid, nw_uuid, parentUuid []byte) (
	[]byte, error) {

	nicSpec := &acropolis_ifc.VmNicSpec{
		NetworkUuid: nw_uuid,
	}

	var nicSpecList []*acropolis_ifc.VmNicSpec
	nicSpecList = append(nicSpecList, nicSpec)

	vmArg := &acropolis_ifc.VmNicCreateArg{
		VmUuid:         vm_uuid,
		SpecList:       nicSpecList,
		ParentTaskUuid: parentUuid,
	}

	vmRet := &acropolis_ifc.VmNicCreateRet{}

	err := svc.sendMsg("VmNicCreate", vmArg, vmRet)
	if err != nil {
		return nil, err
	}

	return vmRet.GetTaskUuid(), nil
}

func (svc *Acropolis) CreateGenericQueryRequest(filter string, sortCriteria string,
	limit uint64) (*acropolis_ifc.GenericQueryRequest, error) {

	req := &acropolis_ifc.GenericQueryRequest{
		FilterCriteria: proto.String(filter),
		SortCriteria:   proto.String(sortCriteria),
		Limit:          proto.Uint64(limit),
	}

	return req, nil
}

func (svc *Acropolis) VmGet(req *acropolis_ifc.GenericQueryRequest,
	uuidList [][]byte, includeNonExistent bool, includeSizes bool,
	includeAddress bool, includeDiskPath bool) ([]*acropolis_ifc.VmInfo,
	[]byte, error) {

	getArg := &acropolis_ifc.VmGetArg{
		Request:                   req,
		IncludeNonExistent:        proto.Bool(includeNonExistent),
		VmUuidList:                uuidList,
		IncludeVmdiskSizes:        proto.Bool(includeSizes),
		IncludeAddressAssignments: proto.Bool(includeAddress),
		IncludeVmdiskPaths:        proto.Bool(includeDiskPath),
	}

	getRet := &acropolis_ifc.VmGetRet{}
	// Send the VmGet Rpc
	err := svc.sendMsg("VmGet", getArg, getRet)
	if err != nil {
		return nil, nil, err
	}

	return getRet.GetVmInfoList(), getRet.GetResponse().GetCursor(), nil
}

func (svc *Acropolis) VmDelete(uuid []byte, parentTaskUuid []byte,
	deleteSnapshots bool) (
	[]byte, error) {
	delArg := &acropolis_ifc.VmDeleteArg{
		VmUuid:          uuid,
		DeleteSnapshots: proto.Bool(deleteSnapshots),
		ParentTaskUuid:  parentTaskUuid,
	}

	delRet := &acropolis_ifc.VmDeleteRet{}

	err := svc.sendMsg("VmDelete", delArg, delRet)

	if err != nil {
		return nil, err
	}
	return delRet.GetTaskUuid(), nil
}

// Set the power state of a Vm
func (svc *Acropolis) VmSetPowerState(uuid []byte,
	state acropolis_ifc.VmStateTransition_Transition,
	parentUuid []byte) ([]byte, error) {

	stateArg := &acropolis_ifc.VmSetPowerStateArg{
		VmUuid:         uuid,
		Transition:     state.Enum(),
		ParentTaskUuid: parentUuid,
	}

	stateRet := &acropolis_ifc.VmSetPowerStateRet{}

	err := svc.sendMsg("VmSetPowerState", stateArg, stateRet)
	if err != nil {
		return nil, err
	}
	return stateRet.GetTaskUuid(), nil
}

func (svc *Acropolis) VmDiskCreate(vmUuid []byte,
	specList []*acropolis_ifc.VmDiskSpec) ([]byte, error) {

	diskArg := &acropolis_ifc.VmDiskCreateArg{
		VmUuid:   vmUuid,
		SpecList: specList,
	}

	diskRet := &acropolis_ifc.VmDiskCreateRet{}

	err := svc.sendMsg("VmDiskCreate", diskArg, diskRet)
	if err != nil {
		return nil, err
	}

	return diskRet.GetTaskUuid(), nil
}

// Attach /Detach / VolumeAdd/ VolumeRemove operation.
func (svc *Acropolis) VmAttachVg(vgUuid []byte, vmUuid []byte,
	op *acropolis_ifc.VmAttachVgArg_Operation, args FuncArguments) (
	[]byte, error) {

	vgArg := &acropolis_ifc.VmAttachVgArg{
		VolumeGroupUuid: vgUuid,
		VmUuid:          vmUuid,
		Operation:       op,
	}

	index, ok := args["Index"].(uint32)
	if ok {
		vgArg.Index = proto.Uint32(index)
	}

	volDiskIndexList, ok := args["VolumeDiskIndexList"].([]uint32)
	if ok {
		vgArg.VolumeDiskIndexList = volDiskIndexList
	}

	ts, ok := args["LogicalTimestamp"].(int64)
	if ok {
		vgArg.LogicalTimestamp = proto.Int64(ts)
	}

	vgRet := &acropolis_ifc.VmAttachVgRet{}

	err := svc.sendMsg("VmAttachVg", vgArg, vgRet)

	if err != nil {
		return nil, err
	}
	return vgRet.GetTaskUuid(), nil
}

// Create a volume group
func (svc *Acropolis) VgCreate(uuid, parentTaskUuid []byte, name string, desc string,
	diskList []*acropolis_ifc.VolumeDiskCreateSpec, shared bool,
	iScsiTargetName string, extIqnList []string,
	fsType *acropolis_ifc.VolumeGroupCreateSpec_FileSystemType) ([]byte, error) {

	vgArg := &acropolis_ifc.VolumeGroupCreateArg{
		ParentTaskUuid: parentTaskUuid,
		Spec: &acropolis_ifc.VolumeGroupCreateSpec{
			Uuid:                     uuid,
			Name:                     proto.String(name),
			Annotation:               proto.String(desc),
			DiskList:                 diskList,
			Shared:                   proto.Bool(shared),
			IscsiTargetName:          proto.String(iScsiTargetName),
			ExternalInitiatorIqnList: extIqnList,
			FileSystemType:           fsType,
		},
	}

	vgRet := &acropolis_ifc.VolumeGroupCreateRet{}

	err := svc.sendMsg("VolumeGroupCreate", vgArg, vgRet)

	if err != nil {
		return nil, err
	}
	return vgRet.GetTaskUuid(), nil
}

func (svc *Acropolis) VgGet(uuidList [][]byte, includeDiskSize bool,
	includeDiskPath bool, attachedVms [][]byte,
	attachedInitiatorList []string) ([]*pithos_ifc.VolumeGroupConfig, error) {

	vgArg := &acropolis_ifc.VolumeGroupGetArg{
		VolumeGroupUuidList:       uuidList,
		IncludeVmdiskSizes:        proto.Bool(includeDiskSize),
		IncludeVmdiskPaths:        proto.Bool(includeDiskPath),
		AttachedVmUuidList:        attachedVms,
		AttachedInitiatorNameList: attachedInitiatorList,
	}

	vgRet := &acropolis_ifc.VolumeGroupGetRet{}

	err := svc.sendMsg("VolumeGroupGet", vgArg, vgRet)
	if err != nil {
		return nil, err
	}

	return vgRet.GetVolumeGroupList(), nil
}

// Add external attachment to a VG
func (svc *Acropolis) VgAddExternalInitiator(uuid []byte,
	iqnList []string) ([]byte, error) {

	vgArg := &acropolis_ifc.VolumeGroupUpdateArg{
		VolumeGroupUuid: uuid,
		ExternalInitiatorIqnList: &acropolis_ifc.VolumeGroupUpdateArg_IqnList{
			IqnList: iqnList,
		},
	}

	vgRet := &acropolis_ifc.VolumeGroupUpdateRet{}

	err := svc.sendMsg("VolumeGroupUpdate", vgArg, vgRet)
	if err != nil {
		return nil, err
	}
	return vgRet.GetTaskUuid(), nil
}

func (svc *Acropolis) VgUpdate(uuid []byte, args FuncArguments) (
	[]byte, error) {

	vgArgs := &acropolis_ifc.VolumeGroupUpdateArg{
		VolumeGroupUuid: uuid,
	}

	name, ok := args["Name"].(string)
	if ok {
		vgArgs.Name = proto.String(name)
	}
	desc, ok := args["Annotation"].(string)
	if ok {
		vgArgs.Annotation = proto.String(desc)
	}
	ts, ok := args["LogicalTimestamp"].(int64)
	if ok {
		vgArgs.LogicalTimestamp = proto.Int64(ts)
	}
	shared, ok := args["Shared"].(bool)
	if ok {
		vgArgs.Shared = proto.Bool(shared)
	}
	isIscsi, ok := args["IscsiTargetName"].(string)
	if ok {
		vgArgs.IscsiTargetName = proto.String(isIscsi)
	}
	extIqnList, ok :=
		args["ExternalInitiatorIqnList"].(*acropolis_ifc.VolumeGroupUpdateArg_IqnList)
	if ok {
		vgArgs.ExternalInitiatorIqnList = extIqnList
	}

	vgRet := &acropolis_ifc.VolumeGroupUpdateRet{}

	err := svc.sendMsg("VolumeGroupUpdateRet", vgArgs, vgRet)

	if err != nil {
		return nil, err
	}
	return vgRet.GetTaskUuid(), nil
}

func (svc *Acropolis) VgDelete(uuid, parentTaskUuid []byte) ([]byte, error) {
	vgArg := &acropolis_ifc.VolumeGroupDeleteArg{
		VolumeGroupUuid: uuid,
		ParentTaskUuid:  parentTaskUuid,
	}

	vgRet := &acropolis_ifc.VolumeGroupDeleteRet{}

	err := svc.sendMsg("VolumeGroupDelete", vgArg, vgRet)
	if err != nil {
		return nil, err
	}
	return vgRet.GetTaskUuid(), nil
}

func (svc *Acropolis) VgDiskCreate(vgUuid []byte,
	diskSpecList []*acropolis_ifc.VolumeDiskCreateSpec, args FuncArguments) (
	[]byte, error) {

	vgArg := &acropolis_ifc.VolumeDiskCreateArg{
		VolumeGroupUuid: vgUuid,
		SpecList:        diskSpecList,
	}

	if len(args) > 0 {
		ts, ok := args["LogicalTimestamp"].(int64)
		if ok {
			vgArg.LogicalTimestamp = proto.Int64(ts)
		}
	}

	vgRet := &acropolis_ifc.VolumeDiskCreateRet{}

	err := svc.sendMsg("VolumeDiskCreate", vgArg, vgRet)
	if err != nil {
		return nil, err
	}

	return vgRet.GetTaskUuid(), nil
}

func (svc *Acropolis) VgAttachExternal(vgUuid []byte, initiatorName string,
	args FuncArguments) error {

	op := acropolis_ifc.VolumeGroupAttachExternalArg_kAttach
	vgArg := &acropolis_ifc.VolumeGroupAttachExternalArg{
		VolumeGroupUuid: vgUuid,
		InitiatorName:   proto.String(initiatorName),
		Operation:       &op,
	}

	if len(args) > 0 {
		ts, ok := args["LogicalTimestamp"].(int64)
		if ok {
			vgArg.LogicalTimestamp = proto.Int64(ts)
		}
	}

	vgRet := &acropolis_ifc.VolumeGroupAttachExternalRet{}

	err := svc.sendMsg("VolumeGroupAttachExternal", vgArg, vgRet)
	if err != nil {
		return err
	}
	return nil
}

func (svc *Acropolis) VgDetachExternal(vgUuid []byte, initiatorName string,
	args FuncArguments) error {

	op := acropolis_ifc.VolumeGroupAttachExternalArg_kDetach
	vgArg := &acropolis_ifc.VolumeGroupAttachExternalArg{
		VolumeGroupUuid: vgUuid,
		InitiatorName:   proto.String(initiatorName),
		Operation:       &op,
	}

	if len(args) > 0 {
		ts, ok := args["LogicalTimestamp"].(int64)
		if ok {
			vgArg.LogicalTimestamp = proto.Int64(ts)
		}
	}

	vgRet := &acropolis_ifc.VolumeGroupAttachExternalRet{}
	err := svc.sendMsg("VolumeGroupAttachExternal", vgArg, vgRet)
	if err != nil {
		return err
	}
	return nil
}

func (svc *Acropolis) processTaskResponse(taskRet proto.Message,
	response *acropolis_ifc.MetaResponse) error {

	retErr := response.GetErrorCode()

	retVal := int(retErr)
	if retVal < 0 && retVal >= len(Errors) {
		return ErrInvalidResponse.SetCause(errors.New(
			"Incorrect Response Code"))
	}
	if retErr != acropolis_ifc.AcropolisError_kNoError {
		return Errors[retVal]
	}
	if taskRet != nil {
		retData := response.GetRet()
		if retData == nil {
			return ErrInvalidResponse.SetCause(errors.New(
				"Missing Return data in response"))
		}

		embedData := retData.GetEmbedded()
		if embedData == nil {
			return ErrInvalidResponse.SetCause(errors.New("Missing embedded data"))
		}

		err := proto.Unmarshal(embedData, taskRet)
		if err != nil {
			return ErrInvalidResponse.SetCause(err)
		}
	}
	return nil
}

// Poll for a submitted task to get the status. Call blocks till
// task is completed.
func (svc *Acropolis) TaskPoll(taskId []byte, timeout int64, interval int64,
	taskProto proto.Message) error {

	taskArg := &acropolis_ifc.TaskPollArg{
		TimeoutSec:   proto.Int64(timeout),
		TaskUuidList: [][]byte{taskId},
	}
	taskRet := &acropolis_ifc.TaskPollRet{}

	done := false
	svc.client.SetRequestTimeout(timeout)
	for !done {
		err := svc.sendMsg("TaskPoll", taskArg, taskRet)
		if err != nil {
			// task timed out.
			return err
		}
		if taskRet.GetTimedOut() {
			glog.Infoln("Task poll timed out - retrying")
			time.Sleep(time.Duration(interval) * time.Second)
			continue
		}
		taskList := taskRet.GetReadyTaskList()
		found := false
		for i := 0; i < len(taskList) && !found; i++ {
			id := taskList[i].GetUuid()
			if bytes.Compare(id[:], taskId[:]) == 0 {
				found = true
				if taskList[i].GetPercentageComplete() < 100 {
					// wait for some time before continue
					glog.Infof("Task still in progress - %d completed\n",
						taskList[i].GetPercentageComplete())
					time.Sleep(time.Duration(interval) * time.Second)
					continue
				}
				// Task Completed.
				done = true
				resp := taskList[i].GetResponse()
				if resp == nil {
					return ErrInvalidResponse.SetCause(errors.New(
						"Missing Embedded Response"))
				}
				err := svc.processTaskResponse(taskProto, resp)
				if err != nil {
					return err
				}
				return nil
			}
		}

	}
	// Dont know what happened.
	return ErrInternalError.SetCause(errors.New("Invalid exit from Task Poll"))
}

// Get all images from image service
func (svc *Acropolis) GetAllImages(includeDiskSize bool, includeDiskPath bool) (
	[]*acropolis_ifc.ImageInfo, error) {

	var uuidlist [][]byte
	return svc.GetImages(&uuidlist, includeDiskSize, includeDiskPath)
}

// Get images specified by Uuid list
func (svc *Acropolis) GetImages(uuidList *[][]byte, includeDiskSize bool,
	includeDiskPath bool) ([]*acropolis_ifc.ImageInfo, error) {

	imgArg := &acropolis_ifc.ImageGetArg{
		ImageUuidList:      *uuidList,
		IncludeVmdiskSizes: proto.Bool(includeDiskSize),
		IncludeVmdiskPaths: proto.Bool(includeDiskPath),
	}

	imgRet := &acropolis_ifc.ImageGetRet{}

	err := svc.sendMsg("ImageGet", imgArg, imgRet)
	if err != nil {
		return nil, err
	}
	return imgRet.GetImageInfoList(), nil
}

// Import an image from a file into image service as disk image
func (svc *Acropolis) CreateDiskImageFromFile(name string, annotation string,
	filePath string, ctrUuid []byte, parentTask []byte, hidden bool) ([]byte, error) {

	url := "nfs://127.0.0.1" + filePath
	imgArg := &acropolis_ifc.ImageCreateArg{
		Spec: &acropolis_ifc.ImageCreateSpec{
			Name:       proto.String(name),
			Annotation: proto.String(annotation),
			ImageType:  acropolis_ifc.ImageInfo_kDiskImage.Enum(),
			Importfrom: &acropolis_ifc.ImageImportSpec{
				Url:           proto.String(url),
				ContainerUuid: ctrUuid,
			},
		},
		ParentTaskUuid: parentTask,
	}
	glog.Infof("Url - %s\n", url)

	imgRet := &acropolis_ifc.ImageCreateRet{}

	err := svc.sendMsg("ImageCreate", imgArg, imgRet)
	if err != nil {
		return nil, err
	}

	return imgRet.GetTaskUuid(), nil
}

func (svc *Acropolis) CreateDiskImageFromDisk(name, annotation string,
	srcVmdiskUuid []byte, dstStorageUuid []byte, parentTask []byte) (
	[]byte, error) {

	imgArg := &acropolis_ifc.ImageCreateArg{
		Spec: &acropolis_ifc.ImageCreateSpec{
			Name:       proto.String(name),
			Annotation: proto.String(annotation),
			ImageType:  acropolis_ifc.ImageInfo_kDiskImage.Enum(),
			Clone: &acropolis_ifc.VmDiskSpecClone{
				VmdiskUuid:    srcVmdiskUuid,
				ContainerUuid: dstStorageUuid,
			},
		},
		ParentTaskUuid: parentTask,
	}

	imgRet := &acropolis_ifc.ImageCreateRet{}

	err := svc.sendMsg("ImageCreate", imgArg, imgRet)
	if err != nil {
		return nil, err
	}

	return imgRet.GetTaskUuid(), nil
}

// Update a VM's disk
func (svc *Acropolis) UpdateVmDisk(uuid []byte,
	diskSpec *acropolis_ifc.VmDiskUpdateArg_VmDiskUpdateSpec,
	parentUuid []byte) ([]byte, error) {

	var updateList []*acropolis_ifc.VmDiskUpdateArg_VmDiskUpdateSpec
	updateList = append(updateList, diskSpec)
	updateArg := &acropolis_ifc.VmDiskUpdateArg{
		VmUuid:         uuid,
		UpdateList:     updateList,
		ParentTaskUuid: parentUuid,
	}

	updateRet := &acropolis_ifc.VmDiskUpdateRet{}

	err := svc.sendMsg("VmDiskUpdate", updateArg, updateRet)

	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	return updateRet.GetTaskUuid(), nil
}

func (svc *Acropolis) UpdateImage(uuid []byte, name, annotation string,
	parentTaskUuid []byte) ([]byte, error) {

	arg := &acropolis_ifc.ImageUpdateArg{
		Spec: &acropolis_ifc.ImageCreateSpec{
			Uuid:       uuid,
			Name:       proto.String(name),
			Annotation: proto.String(annotation),
		},
		ParentTaskUuid: parentTaskUuid,
	}

	ret := &acropolis_ifc.ImageUpdateRet{}

	err := svc.sendMsg("ImageUpdate", arg, ret)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	return ret.GetTaskUuid(), nil

}

func (svc *Acropolis) DeleteImage(uuid []byte, parentTaskUuid []byte) (
	[]byte, error) {

	arg := &acropolis_ifc.ImageDeleteArg{
		ImageUuid:      uuid,
		ParentTaskUuid: parentTaskUuid,
	}

	ret := &acropolis_ifc.ImageDeleteRet{}

	err := svc.sendMsg("ImageDelete", arg, ret)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	return ret.GetTaskUuid(), nil
}

// Get all hosts
func (svc *Acropolis) GetVmsOnHost(uuid []byte) ([][]byte, error) {

	listArg := &acropolis_ifc.LazanGetVMsOnHostArg{HostUuid: uuid}
	listRet := &acropolis_ifc.LazanGetVMsOnHostRet{}

	err := svc.sendMsg("LazanGetVMsOnHost", listArg, listRet)

	if err != nil {
		return nil, err
	}
	vmUuidList := listRet.GetVmUuidList()
	if len(vmUuidList) == 0 {
		return nil, nil
	}
	return vmUuidList, nil
}

func (svc *Acropolis) VmNicGet(vmUuid []byte, macAddrList [][]byte) ([]*acropolis_ifc.VmNicConfig, error) {
	listArg := &acropolis_ifc.VmNicGetArg{VmUuid: vmUuid, MacAddrList: macAddrList, IncludeAddressAssignments: proto.Bool(true)}

	listRet := &acropolis_ifc.VmNicGetRet{}

	err := svc.sendMsg("VmNicGet", listArg, listRet)
	if err != nil {
		return nil, err
	}

	return listRet.NicList, nil
}

// Get address table for a given network
func (svc *Acropolis) GetNetworkAddressTable(networkUuid []byte) (*acropolis_ifc.NetworkAddressAssignmentTable, error) {
	networkAddressTablGetArg := &acropolis_ifc.NetworkAddressTableGetArg{
		NetworkUuid: networkUuid,
	}
	networkAddressTableGetRet := &acropolis_ifc.NetworkAddressTableGetRet{}
	err := svc.sendMsg("NetworkAddressTableGet", networkAddressTablGetArg, networkAddressTableGetRet)
	if err != nil {
		return nil, err
	}
	addressAssignmentTable := networkAddressTableGetRet.GetAddressAssignmentTable()

	return addressAssignmentTable, nil
}
func (svc *Acropolis) ReserveIp(uuid []byte, numIpAddress uint64, cookie string) ([][]byte, error) {

	op := acropolis_ifc.NetworkReserveIpArg_kReserve
	arg := &acropolis_ifc.NetworkReserveIpArg{
		NetworkUuid:    uuid,
		NumIpAddresses: &numIpAddress,
		Operation:      &op,
		Cookie:         proto.String(cookie),
	}
	ret := &acropolis_ifc.NetworkReserveIpRet{}
	err := svc.sendMsg("NetworkReserveIp", arg, ret)
	if err != nil {
		return nil, err
	}
	return ret.IpAddressList, nil
}

func (svc *Acropolis) UnreserveIp(uuid []byte, ipAddrList [][]byte, cookie string) error {

	op := acropolis_ifc.NetworkReserveIpArg_kUnreserve
	arg := &acropolis_ifc.NetworkReserveIpArg{
		NetworkUuid:   uuid,
		IpAddressList: ipAddrList,
		Operation:     &op,
		Cookie:        proto.String(cookie),
	}
	ret := &acropolis_ifc.NetworkReserveIpRet{}
	err := svc.sendMsg("NetworkReserveIp", arg, ret)
	if err != nil {
		return err
	}
	return nil
}
