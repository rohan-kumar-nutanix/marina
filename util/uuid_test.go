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
