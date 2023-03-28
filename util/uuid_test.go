/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Author: rishabh.gupta@nutanix.com
 */

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
)

func TestNew(t *testing.T) {
	uuidUtil := UuidUtil{}
	uuid, err := uuidUtil.New()
	assert.NotNil(t, uuid)
	assert.NoError(t, err)
}

func TestValidateUUID(t *testing.T) {
	uuidUtil := UuidUtil{}
	uuid, _ := uuid4.New()
	err := uuidUtil.ValidateUUID(uuid.RawBytes(), "FooBar")
	assert.NoError(t, err)
}

func TestValidateUUIDError(t *testing.T) {
	uuidUtil := UuidUtil{}
	err := uuidUtil.ValidateUUID([]byte("Invalid UUID"), "FooBar")
	assert.Error(t, err)
	assert.IsType(t, new(marinaError.MarinaInvalidUuidError), err)
}

func TestStringFromUuids(t *testing.T) {
	uuidUtil := UuidUtil{}
	uuid1, _ := uuid4.New()
	uuid2, _ := uuid4.New()
	uuidStr := uuidUtil.StringFromUuids([]uuid4.Uuid{*uuid1, *uuid2})
	assert.Equal(t, uuid1.String()+", "+uuid2.String(), uuidStr)
}

func TestStringFromUuidPointers(t *testing.T) {
	uuidUtil := UuidUtil{}
	uuid1, _ := uuid4.New()
	uuid2, _ := uuid4.New()
	uuidStr := uuidUtil.StringFromUuidPointers([]*uuid4.Uuid{uuid1, uuid2})
	assert.Equal(t, uuid1.String()+", "+uuid2.String(), uuidStr)
}

func TestUuidPointersToUuids(t *testing.T) {
	uuidUtil := UuidUtil{}
	uuid1, _ := uuid4.New()
	uuid2, _ := uuid4.New()
	uuids := uuidUtil.UuidPointersToUuids([]*uuid4.Uuid{uuid1, uuid2})
	assert.Equal(t, []uuid4.Uuid{*uuid1, *uuid2}, uuids)
}

func TestUuidBytesToUuidPointers(t *testing.T) {
	uuidUtil := UuidUtil{}
	uuid1, _ := uuid4.New()
	uuid2, _ := uuid4.New()
	uuids := uuidUtil.UuidBytesToUuidPointers([][]byte{uuid1.RawBytes(), uuid2.RawBytes()})
	assert.Equal(t, []*uuid4.Uuid{uuid1, uuid2}, uuids)
}

func TestUuidBytesToUuids(t *testing.T) {
	uuidUtil := UuidUtil{}
	uuid1, _ := uuid4.New()
	uuid2, _ := uuid4.New()
	uuids := uuidUtil.UuidBytesToUuids([][]byte{uuid1.RawBytes(), uuid2.RawBytes()})
	assert.Equal(t, []uuid4.Uuid{*uuid1, *uuid2}, uuids)
}

func TestUuidToUuidBytes(t *testing.T) {
	uuidUtil := UuidUtil{}
	uuid1, _ := uuid4.New()
	uuid2, _ := uuid4.New()
	uuids := uuidUtil.UuidsToUuidBytes([]uuid4.Uuid{*uuid1, *uuid2})
	assert.Equal(t, [][]byte{uuid1.RawBytes(), uuid2.RawBytes()}, uuids)
}

func TestDifference(t *testing.T) {
	uuidUtil := UuidUtil{}
	uuid1, _ := uuid4.New()
	uuid2, _ := uuid4.New()
	uuids := uuidUtil.Difference([]*uuid4.Uuid{uuid1, uuid2}, []*uuid4.Uuid{uuid2})
	assert.Equal(t, []*uuid4.Uuid{uuid1}, uuids)
}
