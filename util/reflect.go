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
	log "k8s.io/klog/v2"
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

// SetReflectBytes takes a reflect value, and set a specific field of that
// reflect value with a provided value.
func SetReflectBytes(reflectValue reflect.Value, fieldName string, value []byte) error {
	valueElem := reflectValue.Elem()
	field := valueElem.FieldByName(fieldName)
	if !field.IsValid() || !field.CanSet() {
		return marinaError.ErrInternal.SetCause(
			fmt.Errorf("reflect value with field name %s cannot be set", fieldName))
	}
	field.SetBytes(value)
	return nil
}

// GetReflectBytes returns the specified field value for a reflect value.
func GetReflectBytes(reflectValue reflect.Value, fieldName string) ([]byte, error) {
	valueElem := reflectValue.Elem()
	field := valueElem.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, marinaError.ErrInternal.SetCause(
			fmt.Errorf("reflect value doesn't has field name %s", fieldName))
	}
	return field.Bytes(), nil
}

// SetReflectBool takes a reflect value, and set a specific boolean field of
// that reflect value with a provided value.
func SetReflectBool(reflectValue reflect.Value, fieldName string, value bool) error {
	valueElem := reflectValue.Elem()
	field := valueElem.FieldByName(fieldName)
	if !field.IsValid() || !field.CanSet() {
		return marinaError.ErrInternal.SetCause(
			fmt.Errorf("reflect value with field name %s cannot be set", fieldName))
	}
	field.Set(reflect.ValueOf(proto.Bool(value)))
	return nil
}

// ScreenUnknownField does a BFS traverse of the input proto, and returns an error if
// any nested message type contains an unknown field, or the input proto is
// not a valid proto Message.
func ScreenUnknownField(p proto.Message) error {
	// All Message types will be a pointer to struct. We will need to
	// dereference the pointer to get the struct type.
	log.V(5).Infof("Proto: %s.", proto.MarshalTextString(p))
	reflectValue := reflect.ValueOf(p)
	if reflectValue.Kind() != reflect.Ptr {
		return fmt.Errorf("not a pointer type: %v", reflectValue.Kind())
	}
	reflectElem := reflectValue.Elem()
	if reflectElem.Kind() != reflect.Struct {
		return fmt.Errorf("not a struct type: %v", reflectElem.Kind())
	}

	// Feed the queue with the starting point.
	current := []reflect.Value{reflectElem}
	next := []reflect.Value{}
	for len(current) != 0 {
		for i := 0; i < len(current); i++ {
			reflectElem := current[i]
			protoMessage := reflectElem.Addr().Interface().(proto.Message)
			// Check if there is any unknown field in the proto message.
			if len(proto.MessageReflect(protoMessage).GetUnknown()) != 0 {
				return fmt.Errorf("newer proto than supported: %v", protoMessage)
			}

			// Iterate through all fields, and add those not-nil message type to
			// the 'next' queue.
			for j := 0; j < reflectElem.NumField(); j++ {
				field := reflectElem.Field(j)
				switch field.Kind() {
				case reflect.Ptr:
					next = appendMessageField(next, field)
				case reflect.Slice:
					if field.Len() <= 0 {
						break
					}
					if field.Index(0).Kind() != reflect.Ptr {
						// If the entry of the slice is not a pointer, it is
						// not a message type, so we can just skip it.
						break
					}
					for i := 0; i < field.Len(); i++ {
						next = appendMessageField(next, field.Index(i))
					}
				}
			}
		}
		// Swap 'current' and 'next'.
		current = next
		next = []reflect.Value{}
	}
	return nil
}

func appendMessageField(list []reflect.Value, field reflect.Value) []reflect.Value {
	if field.IsNil() {
		return list
	}
	indirect := reflect.Indirect(field)
	if indirect.Kind() == reflect.Struct {
		list = append(list, indirect)
	}
	return list
}
