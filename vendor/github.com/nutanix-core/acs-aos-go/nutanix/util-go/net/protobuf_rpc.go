//
// Copyright (c) 2016 Nutanix Inc. All rights reserved.
// Author: akshay@nutanix.com

package net

import (
	"github.com/golang/protobuf/proto"
	ntnx_errors "github.com/nutanix-core/acs-aos-go/nutanix/util-go/errors"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-slbufs/util/sl_bufs/net"
)

// RPC status values.
const (
	NO_ERROR uint32 = iota
	METHOD_ERROR
	APP_ERROR
	CANCELED
	TIMEOUT
	TRANSPORT_ERROR
)

type RpcError_ struct {
	*ntnx_errors.NtnxError
}

func (e *RpcError_) TypeOfError() int {
	return ntnx_errors.RpcErrorType
}

func RpcError(errMsg string, errCode int) *RpcError_ {
	return &RpcError_{ntnx_errors.NewNtnxError(errMsg, errCode)}
}

func (e *RpcError_) SetCause(err error) error {
	return ntnx_errors.NewNtnxErrorRef(err, e.GetErrorDetail(),
		e.GetErrorCode(), e)
}

func (e *RpcError_) Equals(err interface{}) bool {
	if obj, ok := err.(*RpcError_); ok {
		return e == obj
	}
	return false
}

func (e *RpcError_) IsRetryError() bool {
	return ErrRpcTransport.Equals(e)
}

type AppError_ struct {
	*ntnx_errors.NtnxError
}

func (e *AppError_) TypeOfError() int {
	return ntnx_errors.AppErrorType
}

func (e *AppError_) SetCause(err error) error {
	return ntnx_errors.NewNtnxErrorRef(err, e.GetErrorDetail(),
		e.GetErrorCode(), e)
}

func (e *AppError_) Equals(err interface{}) bool {
	if obj, ok := err.(*AppError_); ok {
		return e == obj
	}
	return false
}

func AppError(errMsg string, errCode int) *AppError_ {
	return &AppError_{ntnx_errors.NewNtnxError(errMsg, errCode)}
}

var (
	ErrMethod      = RpcError("Method Error", 1)
	ErrApplication = RpcError("App Error", 2)
	ErrCanceled    = RpcError("Rpc Canceled", 3)
	ErrTimeout     = RpcError("Rpc Timedout", 4)
	ErrTransport   = RpcError("Rpc Transport Error", 5)
)

var Errors = map[int]error{
	0: nil,
	1: ErrMethod,
	2: ErrApplication,
	3: ErrCanceled,
	4: ErrTimeout,
	5: ErrTransport,
}

var (
	ErrSendingRpc         = RpcError("Invalid Rpc Request", 101)
	ErrInvalidRpcResponse = RpcError("Invalid Rpc response", 102)
	ErrRpcTransport       = RpcError("Error in sending RPC", 103)
)

func ExtractRpcError(err error) (*RpcError_, bool) {
	if obj, ok := ntnx_errors.TypeAssert(err, ntnx_errors.RpcErrorType); ok {
		if rpcErr, ok := obj.(*RpcError_); ok {
			return rpcErr, true
		}
	}
	return nil, false
}

func IsRpcTransportError(err error) bool {
	rpcErr, ok := ExtractRpcError(err)
	return ok && ErrRpcTransport.Equals(rpcErr)
}

// Struct to hold Protobuf RPC related buffers and parameters.
type ProtobufRpc struct {
	RequestHeader   net.RpcRequestHeader
	ResponseHeader  net.RpcResponseHeader
	RequestPayload  []byte
	ResponsePayload []byte
	Done            chan bool
	TimeoutSecs     int64
	RpcError        error
	Status          uint32
}

// Use this method to create new ProtobufRpc.
func NewProtobufRPC() *ProtobufRpc {
	return &ProtobufRpc{Done: make(chan bool), RpcError: nil}
}

func (rpc *ProtobufRpc) SetAppError(appError int32, errorDetail string) {
	rpc.ResponseHeader.AppError = proto.Int32(appError)
	rpc.ResponseHeader.ErrorDetail = proto.String(errorDetail)
	rpc.ResponseHeader.RpcStatus = net.RpcResponseHeader_kAppError.Enum()

}

func ExtractAppError(err error) (*AppError_, bool) {
	if obj, ok := ntnx_errors.TypeAssert(err, ntnx_errors.AppErrorType); ok {
		if errObj, ok := obj.(*AppError_); ok {
			return errObj, true
		}
	}
	return nil, false
}
