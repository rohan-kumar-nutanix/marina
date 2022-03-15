/*
 Copyright (c) 2022 Nutanix Inc. All rights reserved.

 Authors: rajesh.battala@nutanix.com

 This file implements utility function related to the use of reflect.
 These methods are used For ProxyServer Implementation.ß
*/

package utils

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	marinaError "github.com/nutanix-core/content-management-marina/error"
	"reflect"
)

// UnmarshalReflectBytes unmarshal a given reflect byte value into the specified
// reflect type.
func UnmarshalReflectBytes(byteValue []byte, reflectValue reflect.Value) error {
	args := []reflect.Value{
		reflect.ValueOf(byteValue),
		reflectValue,
	}
	unmarshalMethod := reflect.ValueOf(proto.Unmarshal)
	unmarshalError := unmarshalMethod.Call(args)[0]
	if !unmarshalError.IsNil() {
		return marinaError.ErrMarinaInternal.SetCause(
			fmt.Errorf("failed to unmarshal byte value"))
	}
	return nil
}

// MarshalReflectBytes marshal a given reflect value into bytes.
func MarshalReflectBytes(reflectValue reflect.Value) ([]byte, error) {
	values := []reflect.Value{reflectValue}
	marshalMethod := reflect.ValueOf(proto.Marshal)
	responseValues := marshalMethod.Call(values)
	if !responseValues[1].IsNil() {
		return nil, marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to marshal reflect value"))
	}
	return responseValues[0].Interface().([]byte), nil
}
