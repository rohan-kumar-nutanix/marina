/*
 * Copyright (c) 2017 Nutanix Inc. All rights reserved.
 *
 * Author: chaitanya.karlekar@nutanix.com
 *
 * Defines client interface for Catalog service.
 * XXX: This code should be generated.
 */

package catalog_client

import (
	"errors"
	"flag"
	"fmt"
	"time"

	glog "github.com/golang/glog"
	proto "github.com/golang/protobuf/proto"

	catalog_ifc "github.com/nutanix-core/acs-aos-go/catalog"

	ntnx_errors "github.com/nutanix-core/acs-aos-go/nutanix/util-go/errors"
	util_misc "github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
	util_net "github.com/nutanix-core/acs-aos-go/nutanix/util-go/net"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
)

const (
	catalogServiceName = "nutanix.catalog.CatalogRpcService"
)

var (
	defaultCatalogAddr string
	defaultCatalogPort uint16 = 2007
)

func init() {
	flag.StringVar(&defaultCatalogAddr,
		"catalog_ip",
		"127.0.0.1",
		"IP address of Catalog service to connect.")

	flag.Var(&util_misc.PortNumber{&defaultCatalogPort},
		"catalog_port",
		"Port number of the Catalog service to connect.")
}

// Catalog Error defn.
type CatalogError_ struct {
	*ntnx_errors.NtnxError
}

func (e *CatalogError_) TypeOfError() int {
	return ntnx_errors.CatalogErrorType
}

func (e *CatalogError_) SetCause(err error) error {
	return ntnx_errors.NewNtnxErrorRef(err, e.GetErrorDetail(),
		e.GetErrorCode(), e)
}

func (e *CatalogError_) Equals(err interface{}) bool {
	if obj, ok := err.(*CatalogError_); ok {
		return e == obj
	}
	return false
}

func CatalogError(errMsg string, errCode int) *CatalogError_ {
	return &CatalogError_{ntnx_errors.NewNtnxError(errMsg, errCode)}
}

// The error code offsets are based on error codes defined in
// catalog_client/catalog_error.proto.client
var (
	ErrNoError                     = CatalogError("NoError", 0)
	ErrCanceled                    = CatalogError("Canceled", 1)
	ErrRetry                       = CatalogError("Retry", 2)
	ErrTimeout                     = CatalogError("Timeout", 3)
	ErrUncaughtException           = CatalogError("UncaughtException", 4)
	ErrInvalidArgument             = CatalogError("InvalidArgument", 5)
	ErrLogicalTimestampMismatch    = CatalogError("LogicalTimestampMismatch", 6)
	ErrNotFound                    = CatalogError("NotFound", 7)
	ErrNotMaster                   = CatalogError("NotMaster", 8)
	ErrExists                      = CatalogError("Exists", 9)
	ErrImportError                 = CatalogError("ImportError", 10)
	ErrImageRpcForwardError        = CatalogError("ImageRpcForwardError", 11)
	ErrCatalogUploadFailure        = CatalogError("CatalogUploadFailure", 12)
	ErrCatalogTaskForwardError     = CatalogError("CatalogTaskForwardError", 13)
	ErrCatalogFileLockError        = CatalogError("CatalogFileLockError", 14)
	ErrCatalogFileUnlockError      = CatalogError("CatalogFileUnlockError", 15)
	ErrCatalogItemCheckoutError    = CatalogError("CatalogItemCheckoutError", 16)
	ErrCatalogItemMigrateError     = CatalogError("CatalogItemMigrateError", 17)
	ErrCatalogUnregCleanupError    = CatalogError("CatalogUnregCleanupError", 18)
	ErrNotSupported                = CatalogError("NotSupported", 19)
	ErrCasError                    = CatalogError("CasError", 20)
	ErrInternal                    = CatalogError("Internal", 21)
	ErrCatalogFileCheckoutError    = CatalogError("CatalogFileCheckoutError", 22)
	ErrCatalogRemoteSeedingError   = CatalogError("CatalogRemoteSeedingError", 23)
	ErrCatalogPlacementPolicyError = CatalogError("CatalogPlacementPolicyError", 24)
	ErrCatalogItemUncheckoutError  = CatalogError("CatalogItemUncheckoutError", 25)
	ErrImageCheckoutError          = CatalogError("ImageCheckoutError", 26)
	ErrCatalogItemRegUpdateError   = CatalogError("CatalogItemRegUpdateError", 27)
	ErrCatalogRateLimitError       = CatalogError("CatalogRateLimitError", 28)
)

var Errors = map[int]*CatalogError_{
	0:  ErrNoError,
	1:  ErrCanceled,
	2:  ErrRetry,
	3:  ErrTimeout,
	4:  ErrUncaughtException,
	5:  ErrInvalidArgument,
	6:  ErrLogicalTimestampMismatch,
	7:  ErrNotFound,
	8:  ErrNotMaster,
	9:  ErrExists,
	10: ErrImportError,
	11: ErrImageRpcForwardError,
	12: ErrCatalogUploadFailure,
	13: ErrCatalogTaskForwardError,
	14: ErrCatalogFileLockError,
	15: ErrCatalogFileUnlockError,
	16: ErrCatalogItemCheckoutError,
	17: ErrCatalogItemMigrateError,
	18: ErrCatalogUnregCleanupError,
	19: ErrNotSupported,
	20: ErrCasError,
	21: ErrInternal,
	22: ErrCatalogFileCheckoutError,
	23: ErrCatalogRemoteSeedingError,
	24: ErrCatalogPlacementPolicyError,
	25: ErrCatalogItemUncheckoutError,
	26: ErrImageCheckoutError,
	27: ErrCatalogItemRegUpdateError,
	28: ErrCatalogRateLimitError,
}

type Catalog struct {
	serverIp   string
	serverPort uint16
	client     util_net.ProtobufRPCClientIfc
}

type CatalogClientInterface interface {
	CatalogItemCheckout(*catalog_ifc.CatalogItemCheckoutArg) (
		*catalog_ifc.CatalogItemCheckoutRet, error)
	CatalogItemCreate(*catalog_ifc.CatalogItemCreateArg) (
		*catalog_ifc.CatalogItemCreateRet, error)
	CatalogItemDelete(*catalog_ifc.CatalogItemDeleteArg) (
		*catalog_ifc.CatalogItemDeleteRet, error)
	CatalogItemGet(*catalog_ifc.CatalogItemGetArg) (
		*catalog_ifc.CatalogItemGetRet, error)
	CatalogItemUpdate(*catalog_ifc.CatalogItemUpdateArg) (
		*catalog_ifc.CatalogItemUpdateRet, error)
	FileGet(*catalog_ifc.FileGetArg) (
		*catalog_ifc.FileGetRet, error)
	ImageCheckout(*catalog_ifc.ImageCheckoutArg) (
		*catalog_ifc.ImageCheckoutRet, error)
}

func DefaultCatalogService() *Catalog {
	return NewCatalogService(defaultCatalogAddr, defaultCatalogPort)
}

func NewCatalogService(serviceIP string, servicePort uint16) *Catalog {
	return &Catalog{
		client:     util_net.NewProtobufRPCClient(serviceIP, servicePort),
		serverIp:   serviceIP,
		serverPort: servicePort,
	}
}

type RemoteCatalogTask struct {
	MethodName string
	Client     *Catalog
	Uuid       *uuid4.Uuid
	Arg        interface{}
}

func (task *RemoteCatalogTask) GetUuid() (*uuid4.Uuid, error) {
	if task.Uuid == nil {
		uuid, err := uuid4.New()
		if err != nil {
			return nil, err
		}
		task.Uuid = uuid
	}
	return task.Uuid, nil
}

func (task *RemoteCatalogTask) CallRpc(RpcName string, arg interface{},
	SubtaskSequenceId uint64, ParentTaskUuid []byte) ([]byte, error) {
	var err error
	var taskUuid []byte
	switch catalogArg := arg.(type) {
	case *catalog_ifc.CatalogItemCheckoutArg:
		catalogRet := &catalog_ifc.CatalogItemCheckoutRet{}
		catalogArg.TaskUuid = task.Uuid.RawBytes()
		catalogArg.SubtaskSequenceId = proto.Uint64(SubtaskSequenceId)
		catalogArg.ParentTaskUuid = ParentTaskUuid
		err = task.Client.sendMsg("CatalogItemCheckout", catalogArg, catalogRet)
		taskUuid = catalogRet.TaskUuid
	case *catalog_ifc.ImageCheckoutArg:
		catalogRet := &catalog_ifc.ImageCheckoutRet{}
		catalogArg.TaskUuid = task.Uuid.RawBytes()
		catalogArg.SubtaskSequenceId = proto.Uint64(SubtaskSequenceId)
		catalogArg.ParentTaskUuid = ParentTaskUuid
		err = task.Client.sendMsg("ImageCheckout", catalogArg, catalogRet)
		taskUuid = catalogRet.TaskUuid
	default:
		return nil, ErrNotFound.SetCause(errors.New(
			"Remote subtask is not supported"))
	}

	if err != nil {
		return nil, err
	}

	taskUuidStr := uuid4.ToUuid4(taskUuid).String()
	if task.Uuid != nil && taskUuidStr != task.Uuid.String() {
		return nil, ErrNotFound.SetCause(errors.New(
			fmt.Sprintf("Task's uuid %s is not as expected", taskUuidStr)))
	}
	return taskUuid, nil
}

func (task *RemoteCatalogTask) GetArg() interface{} {
	return task.Arg
}

func (task *RemoteCatalogTask) RpcName() string {
	return task.MethodName
}

func (svc *Catalog) sendMsg(
	service string, request, response proto.Message) error {

	retryWait := util_misc.NewExponentialBackoff(1*time.Second, 30*time.Second,
		10)
	var done bool = false
	for !done {
		err := svc.client.CallMethodSync(catalogServiceName, service, request,
			response, 0)
		if err != nil {
			// If App error is set in RPC - extract the code and get the relevant
			// catalog error.

			if obj, ok := ntnx_errors.TypeAssert(err, ntnx_errors.AppErrorType); ok {
				if errObj, ok := obj.(*util_net.AppError_); ok {
					errCode := errObj.GetErrorCode()
					glog.Errorf("AppError: %d", errCode)
					return Errors[errCode]
				} else {
					return ErrInvalidArgument
				}
			}

			if obj, ok := ntnx_errors.TypeAssert(err, ntnx_errors.RpcErrorType); ok {
				if rpcErr, ok := obj.(*util_net.RpcError_); ok {
					errCode := rpcErr.GetErrorCode()
					glog.Errorf("RpcError: %d", errCode)
					if util_net.ErrRpcTransport.Equals(rpcErr) {
						glog.Infof("Failure to send msg, will retry later")
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

func (svc *Catalog) CatalogItemCheckout(
	arg *catalog_ifc.CatalogItemCheckoutArg) (*catalog_ifc.CatalogItemCheckoutRet, error) {
	ret := &catalog_ifc.CatalogItemCheckoutRet{}
	err := svc.sendMsg("CatalogItemCheckout", arg, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (svc *Catalog) CatalogItemCreate(arg *catalog_ifc.CatalogItemCreateArg) (
	*catalog_ifc.CatalogItemCreateRet, error) {
	ret := &catalog_ifc.CatalogItemCreateRet{}
	err := svc.sendMsg("CatalogItemCreate", arg, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (svc *Catalog) CatalogItemDelete(arg *catalog_ifc.CatalogItemDeleteArg) (
	*catalog_ifc.CatalogItemDeleteRet, error) {
	ret := &catalog_ifc.CatalogItemDeleteRet{}
	err := svc.sendMsg("CatalogItemDelete", arg, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (svc *Catalog) CatalogItemGet(arg *catalog_ifc.CatalogItemGetArg) (
	*catalog_ifc.CatalogItemGetRet, error) {
	ret := &catalog_ifc.CatalogItemGetRet{}
	err := svc.sendMsg("CatalogItemGet", arg, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (svc *Catalog) CatalogItemUpdate(arg *catalog_ifc.CatalogItemUpdateArg) (
	*catalog_ifc.CatalogItemUpdateRet, error) {
	ret := &catalog_ifc.CatalogItemUpdateRet{}
	err := svc.sendMsg("CatalogItemUpdate", arg, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (svc *Catalog) FileGet(arg *catalog_ifc.FileGetArg) (
	*catalog_ifc.FileGetRet, error) {
	ret := &catalog_ifc.FileGetRet{}
	err := svc.sendMsg("FileGet", arg, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (svc *Catalog) ImageCheckout(
	arg *catalog_ifc.ImageCheckoutArg) (*catalog_ifc.ImageCheckoutRet, error) {
	ret := &catalog_ifc.ImageCheckoutRet{}
	err := svc.sendMsg("ImageCheckout", arg, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
