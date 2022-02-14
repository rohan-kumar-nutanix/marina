/*
 * Copyright (c) 2016 Nutanix Inc. All rights reserved.
 *
 * This module implements a base error handling mechanism to be used in GO
 * It allows a service to embed an underlying cause in a predefined error
 * The caller/user of the service can look into the nested errors by
 * using the TypeAssert.
 *
 */
package errors

import (
	"fmt"
)

type INtnxError interface {
	TypeOfError() int
	GetErrorCode() int
	GetErrorDetail() string
	Error() string
}

type NtnxError struct {
	_error error
	msg    string
	code   int
	_ref   INtnxError
	//_ref   error
}

const (
	NtnxErrorType = iota
	RpcErrorType
	AppErrorType
	AcropolisErrorType
	ErgonErrorType
	UhuraErrorType
	InsightErrorType
	OrionErrorType
	AcsErrorType
	ArithmosErrorType
	NDFSErrorType
	MetropolisErrorType
	RegistryErrorType
	NgtErrorType
	TaskErrorType
	ConfigErrorType
	LazanErrorType
	CatalogErrorType
	LogCollectorErrorType
	AndurilErrorType
	BracketsErrorType
	PolluxErrorType
	StatsGatewayErrorType
	CastorErrorType
	MinervaErrorType
	DpmErrorType
	MagnetoErrorType
	VCenterErrorType
	SecurityErrorType
)

func (e *NtnxError) Error() string {
	if e._error == nil {
		return fmt.Sprintf("%s: %d", e.msg, e.code)
	} else {
		return fmt.Sprintf("%s: %d\n  :%s", e.msg, e.code, e._error)
	}
}

func (e *NtnxError) TypeOfError() int {
	return NtnxErrorType
}

func NewNtnxError(errMsg string, errCode int) *NtnxError {
	return &NtnxError{msg: errMsg, code: errCode}
}

func (e *NtnxError) SetCause(err error) error {
	return &NtnxError{err, e.msg, e.code, e}
}

func NewNtnxErrorRef(err error, msg string, code int, ref INtnxError) *NtnxError {
	return &NtnxError{err, msg, code, ref}
}

func (e *NtnxError) GetErrorCode() int {
	return e.code
}

func (e *NtnxError) GetErrorDetail() string {
	return e.msg
}

func (e *NtnxError) Equals(err interface{}) bool {
	if obj, ok := err.(*NtnxError); ok {
		return e == obj
	}
	return false
}

//
// Return true and the appropriate interface if the error is
// an INtnxError of type errorType
func TypeAssert(e error, errorType int) (INtnxError, bool) {

	if ntnxErr, ok := e.(INtnxError); ok {
		for ntnxErr != nil {
			if ntnxErr.TypeOfError() == errorType {
				return ntnxErr, true
			} else {
				valErr, ok := ntnxErr.(*NtnxError)
				if ok {
					ntnxErr = valErr._ref
				} else {
					ntnxErr = nil
				}
			}
		}
	}
	return nil, false
}
