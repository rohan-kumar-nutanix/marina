/*
 *
 * UUID Utils code.
 *
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: Rajesh Battala <rajesh.battala@nutanix.com>
 *
 */

package utils

import (
	"errors"
	"fmt"

	set "github.com/deckarep/golang-set/v2"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
)

var NilUuid uuid4.Uuid

type UuidUtilInterface interface {
	New() (*uuid4.Uuid, error)
	ValidateUUID(uuidValue []byte, fieldName string) error
	StringFromUuids(uuids []uuid4.Uuid) string
	StringFromUuidPointers(uuids []*uuid4.Uuid) string
	UuidPointersToUuids(uuids []*uuid4.Uuid) []uuid4.Uuid
	UuidBytesToUuidPointers(uuids [][]byte) []*uuid4.Uuid
	UuidBytesToUuids(uuids [][]byte) []uuid4.Uuid
	UuidsToUuidBytes(uuids []uuid4.Uuid) [][]byte
	Difference(uuidListA []*uuid4.Uuid, uuidListB []*uuid4.Uuid) []*uuid4.Uuid
}

type UuidUtil struct {
}

func (*UuidUtil) New() (*uuid4.Uuid, error) {
	uuid, err := uuid4.New()
	if err != nil {
		errMsg := fmt.Sprintf("Error while creating a UUID: %v", err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return uuid, nil
}

func (*UuidUtil) ValidateUUID(uuidValue []byte, fieldName string) error {
	if uuid := uuid4.ToUuid4(uuidValue); uuid == nil {
		errMsg := fmt.Sprintf("Invalid '%s' (%s). UUID must be exactly 16 bytes string",
			fieldName, string(uuidValue))
		return marinaError.ErrMarinaInvalidUuid(string(uuidValue)).SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}

func (*UuidUtil) StringFromUuids(uuids []uuid4.Uuid) string {
	var uuidsStr string
	for i, uuid := range uuids {
		if i == len(uuids)-1 {
			uuidsStr += uuid.String()
		} else {
			uuidsStr += uuid.String() + ", "
		}
	}
	return uuidsStr
}

func (*UuidUtil) StringFromUuidPointers(uuids []*uuid4.Uuid) string {
	var uuidsStr string
	for i, uuid := range uuids {
		if i == len(uuids)-1 {
			uuidsStr += uuid.String()
		} else {
			uuidsStr += uuid.String() + ", "
		}
	}
	return uuidsStr
}

func (*UuidUtil) UuidPointersToUuids(uuids []*uuid4.Uuid) []uuid4.Uuid {
	var ret []uuid4.Uuid
	for _, uuid := range uuids {
		ret = append(ret, *uuid)
	}
	return ret
}

func (*UuidUtil) UuidBytesToUuidPointers(uuids [][]byte) []*uuid4.Uuid {
	var ret []*uuid4.Uuid
	for _, uuid := range uuids {
		ret = append(ret, uuid4.ToUuid4(uuid))
	}
	return ret
}

func (*UuidUtil) UuidBytesToUuids(uuids [][]byte) []uuid4.Uuid {
	var ret []uuid4.Uuid
	for _, uuid := range uuids {
		ret = append(ret, *uuid4.ToUuid4(uuid))
	}
	return ret
}

func (*UuidUtil) UuidsToUuidBytes(uuids []uuid4.Uuid) [][]byte {
	var ret [][]byte
	for _, uuid := range uuids {
		uuid := uuid
		ret = append(ret, uuid.RawBytes())
	}
	return ret
}

// Difference - Return the difference A - B
func (*UuidUtil) Difference(uuidsA []*uuid4.Uuid, uuidsB []*uuid4.Uuid) []*uuid4.Uuid {
	uuidsBSet := set.NewSet[uuid4.Uuid]()
	for _, clusterUuid := range uuidsB {
		uuidsBSet.Add(*clusterUuid)
	}

	var ret []*uuid4.Uuid
	for _, clusterUuid := range uuidsA {
		if !uuidsBSet.Contains(*clusterUuid) {
			ret = append(ret, clusterUuid)
		}
	}

	return ret
}
