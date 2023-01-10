/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 * This file is to make RPC requests to Catalog RPC Server
 *
 */

package catalog_client

import (
	"time"

	"github.com/golang/protobuf/proto"
	log "k8s.io/klog/v2"

	ntnx_errors "github.com/nutanix-core/acs-aos-go/nutanix/util-go/errors"
	util_misc "github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
	util_net "github.com/nutanix-core/acs-aos-go/nutanix/util-go/net"
	"github.com/nutanix-core/content-management-marina/common"
	utils "github.com/nutanix-core/content-management-marina/util"
)

const (
	catalogServiceName = "nutanix.catalog.CatalogRpcService"
)

// Catalog Error definition.
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

func DefaultCatalogService() *Catalog {
	return NewCatalogService(utils.HostAddr, uint16(*common.CatalogPort))
}

func NewCatalogService(serviceIP string, servicePort uint16) *Catalog {
	return &Catalog{
		client:     util_net.NewProtobufRPCClient(serviceIP, servicePort),
		serverIp:   serviceIP,
		serverPort: servicePort,
	}
}

// TODO Implement SendMsg with Context once the ProtobufRPCClientIfc lib has the methods.

func (svc *Catalog) SendMsg(
	service string, request, response proto.Message) error {
	log.Info("Sending RPC request to Catalog Service for method : ", service)
	return svc.sendMsg(service, request, response)
}

func (svc *Catalog) sendMsg(
	service string, request, response proto.Message) error {

	retryWait := util_misc.NewExponentialBackoff(1*time.Second, 30*time.Second,
		10)
	var done bool = false
	for !done {
		err := svc.client.CallMethodSync(catalogServiceName, service, request,
			response, utils.RpcTimeoutSec)
		if err != nil {
			// If App error is set in RPC - extract the code and get the relevant
			// catalog error.

			if obj, ok := ntnx_errors.TypeAssert(err, ntnx_errors.AppErrorType); ok {
				if errObj, ok := obj.(*util_net.AppError_); ok {
					errCode := errObj.GetErrorCode()
					log.Errorf("AppError: %d", errCode)
					return Errors[errCode]
				} else {
					return ErrInvalidArgument
				}
			}

			if obj, ok := ntnx_errors.TypeAssert(err, ntnx_errors.RpcErrorType); ok {
				if rpcErr, ok := obj.(*util_net.RpcError_); ok {
					errCode := rpcErr.GetErrorCode()
					log.Errorf("RpcError: %d", errCode)
					if util_net.ErrRpcTransport.Equals(rpcErr) {
						log.Infof("Failure to send msg, will retry later")
						waited := retryWait.Backoff()
						if waited != util_misc.Stop {
							log.Infof("Retrying, msg: %s", service)
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
