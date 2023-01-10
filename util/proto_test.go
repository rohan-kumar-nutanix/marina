/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Author: rishabh.gupta@nutanix.com
 */

package utils

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	marinaError "github.com/nutanix-core/content-management-marina/errors"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
)

func TestMarshal(t *testing.T) {
	protoUtil := ProtoUtil{}
	bytes, err := protoUtil.Marshal(&marinaIfc.CatalogItemInfo{})
	assert.NotNil(t, bytes)
	assert.NoError(t, err)
}

func TestMarshalError(t *testing.T) {
	protoUtil := ProtoUtil{}
	var msg proto.Message
	bytes, err := protoUtil.Marshal(msg)
	assert.Nil(t, bytes)
	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
}
