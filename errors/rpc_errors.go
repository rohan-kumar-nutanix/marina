/*
 *
 * Marina service errors.
 *
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: Rajesh Battala <rajesh.battala@nutanix.com>
 *
 * Marina Error and Interface Types
 */

package errors

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/errors"

	ntnxCmsApiError "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/error"
)

// MarinaErrorInterface shall be embedded in each of the app error struct defined below.
type MarinaErrorInterface interface {
	errors.INtnxError
	SetCauseAndLog(err error) MarinaErrorInterface
	SetCause(err error) MarinaErrorInterface
	ConvertToAppMessagePb() *ntnxCmsApiError.AppMessage
}

type MarinaError struct {
	*errors.NtnxError
}

var (
	ErrNoError                  = err("NO_ERROR", 0)
	ErrCanceled                 = err("CANCELED", 1)
	ErrRetry                    = err("RETRY", 2)
	ErrTimeout                  = err("TIMEOUT", 3)
	ErrUncaughtException        = err("UNCAUGHT_EXCEPTION", 4)
	ErrInvalidArgument          = err("INVALID_ARGUMENT", 5)
	ErrLogicalTimestampMismatch = err("LOGICAL", 6)
	ErrNotFound                 = err("NOT_FOUND", 7)
	ErrNotMaster                = err("NOT_MASTER", 8)
	ErrExists                   = err("EXISTS", 9)
	ErrImportError              = err("IMPORT_ERROR", 10)
	ErrImageRpcForwardError     = err("IMAGE_RPC_FORWARD_ERROR", 11)
	ErrCatalogUploadFailure     = err("CATALOG_UPLOAD_FAILURE", 12)
	ErrCatalogTaskForwardError  = err("CATALOG_TASK_FORWARD_ERROR", 13)
	ErrNotSupported             = err("NOT_SUPPORTED", 19)
	ErrCasError                 = err("CAS_ERROR", 20)
	ErrInternal                 = err("INTERNAL", 21)
	ErrMarinaInternal           = err("INTERNAL", 21)
	ErrUnimplemented            = err("UNIMPLEMENTED", 31)
)

var Errors = map[int]*MarinaError{
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
	// 14: ErrCatalogFileLockError,
	// 15: ErrCatalogFileUnlockError,
	// 16: ErrCatalogItemCheckoutError,
	// 17: ErrCatalogItemMigrateError,
	// 18: ErrCatalogUnregCleanupError,
	19: ErrNotSupported,
	20: ErrCasError,
	21: ErrInternal,
	// 22: ErrCatalogFileCheckoutError,
	// 23: ErrCatalogRemoteSeedingError,
	// 24: ErrCatalogPlacementPolicyError,
	// 25: ErrCatalogItemUncheckoutError,
	// 26: ErrImageCheckoutError,
	// 27: ErrCatalogItemRegUpdateError,
	// 28: ErrCatalogRateLimitError,
	// 29: ErrCatalogTaskAbortedError,
	// 30: ErrCatalogClusterLockExtensionError,
}

// TODO: Map / Define these error codes in ntnx-api-catalog repo.
var marinaAppErrorCodeGrpcCodeMapping = map[int]codes.Code{
	20201: codes.Internal,
	20301: codes.InvalidArgument,
}

var marinaErrorCodeMapping = map[int]MarinaErrorInterface{
	20201: ErrMarinaInternalError(),
	20301: ErrMarinaInvalidArgument("", "", ""),
}

var marinaErrorCodeGrpcCodeMapping = map[int]codes.Code{
	0:     codes.OK,
	1:     codes.Canceled,
	5:     codes.InvalidArgument,
	7:     codes.NotFound,
	9:     codes.AlreadyExists,
	21:    codes.Internal,
	10001: codes.Unimplemented,
	10002: codes.Unavailable,
	20201: codes.Internal,
	20301: codes.InvalidArgument,
}

func GetGrpcCodeFromMarinaError(marinaError MarinaErrorInterface) codes.Code {
	if grpcCode, found := marinaErrorCodeGrpcCodeMapping[marinaError.GetErrorCode()]; found {
		return grpcCode
	} else {
		log.Warningf("Unable to map Marina error '%s' to gRPC code.", marinaError)
		return codes.Internal
	}
}

// TypeOfError returns the Marina error type.
func (e *MarinaError) TypeOfError() int {
	// TODO: Add MarinaErrorType to NtnxError util master.
	return errors.CatalogErrorType
}

// SetCause sets the cause of the provided error.
func (e *MarinaError) SetCause(err error) MarinaErrorInterface {
	return &MarinaError{e.NtnxError.SetCause(err).(*errors.NtnxError)}
}

// SetCauseAndLog sets the cause of the provided error, and logs the err message.
func (e *MarinaError) SetCauseAndLog(err error) MarinaErrorInterface {
	_, fullPath, line, _ := runtime.Caller(1)
	paths := strings.Split(fullPath, "/")
	pathLen := len(paths)
	message := fmt.Sprintf("%s/%s:%d %s", paths[pathLen-2], paths[pathLen-1],
		line, err.Error())
	log.Error(message)
	return &MarinaError{e.NtnxError.SetCause(fmt.Errorf("%s", message)).(*errors.NtnxError)}
}

func (e *MarinaError) ConvertToAppMessagePb() *ntnxCmsApiError.AppMessage {
	return &ntnxCmsApiError.AppMessage{
		Code:    proto.String(strconv.Itoa(e.GetErrorCode())),
		Message: proto.String(e.GetErrorDetail()),
		// TODO Add arg map
		ArgumentsMap: nil,
	}
}

// err creates a new Marina error.
func err(errMsg string, errCode int) *MarinaError {
	return &MarinaError{
		NtnxError: errors.NewNtnxError(errMsg, errCode),
	}
}

// InternalError - Internal error.
type InternalError struct {
	*MarinaError
}

func ErrInternalError() *InternalError {
	return &InternalError{
		MarinaError: err("Marina Server Internal Error", 21),
	}
}

func (e *InternalError) SetCauseAndLog(err error) MarinaErrorInterface {
	e.MarinaError.SetCauseAndLog(err)
	return e
}

func (e *InternalError) SetCause(err error) MarinaErrorInterface {
	e.MarinaError.SetCause(err)
	return e
}

// MarinaInternalError - Marina Internal error.
type MarinaInternalError struct {
	*MarinaError
}

func ErrMarinaInternalError() *InternalError {
	return &InternalError{
		MarinaError: err("Operation Timeout Error", 21),
	}
}

func (e *MarinaInternalError) SetCauseAndLog(err error) MarinaErrorInterface {
	e.MarinaError.SetCauseAndLog(err)
	return e
}

func (e *MarinaInternalError) SetCause(err error) MarinaErrorInterface {
	e.MarinaError.SetCause(err)
	return e
}

// MarinaInvalidArgumentError
type MarinaInvalidArgumentError struct {
	*MarinaError
	catalogItemUuid string
	argumentKey     string
	argumentValue   string
}

func ErrMarinaInvalidArgument(uuid, argumentKey, argumentValue string) *MarinaInvalidArgumentError {
	return &MarinaInvalidArgumentError{
		MarinaError:     err("Catalog Invalid Argument Error", 5),
		catalogItemUuid: uuid,
		argumentKey:     argumentKey,
		argumentValue:   argumentValue,
	}
}

func (e *MarinaInvalidArgumentError) SetCauseAndLog(err error) MarinaErrorInterface {
	e.MarinaError.SetCauseAndLog(err)
	return e
}

func (e *MarinaInvalidArgumentError) SetCause(err error) MarinaErrorInterface {
	e.MarinaError.SetCause(err)
	return e
}

// MarinaInvalidUuidError.
type MarinaInvalidUuidError struct {
	*MarinaError
	uuid string
}

func ErrMarinaInvalidUuid(uuid string) *MarinaInvalidUuidError {
	return &MarinaInvalidUuidError{
		MarinaError: err(fmt.Sprintf("Catalog Invalid UUID (%s) Argument Error", uuid), 5),
		uuid:        uuid,
	}
}

func (e *MarinaInvalidUuidError) SetCauseAndLog(err error) MarinaErrorInterface {
	e.MarinaError.SetCauseAndLog(err)
	return e
}

func (e *MarinaInvalidUuidError) SetCause(err error) MarinaErrorInterface {
	e.MarinaError.SetCause(err)
	return e
}

// MarinaNotSupportedError - Marina Not supported error.
type MarinaNotSupportedError struct {
	*MarinaError
	operation string
}

func ErrMarinaNotSupportedError(operation string) *MarinaNotSupportedError {
	return &MarinaNotSupportedError{
		MarinaError: err(fmt.Sprintf("Operation %s not supported by Marina", operation), 19),
		operation:   operation,
	}
}

func (e *MarinaNotSupportedError) SetCauseAndLog(err error) MarinaErrorInterface {
	e.MarinaError.SetCauseAndLog(err)
	return e
}

func (e *MarinaNotSupportedError) SetCause(err error) MarinaErrorInterface {
	e.MarinaError.SetCause(err)
	return e
}
