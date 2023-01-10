/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Author: rishabh.gupta@nutanix.com
*
* The implementation for Proto utils
 */

package utils

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"

	marinaError "github.com/nutanix-core/content-management-marina/errors"
)

type ProtoUtilInterface interface {
	Marshal(msg proto.Message) ([]byte, error)
}

type ProtoUtil struct {
}

func (*ProtoUtil) Marshal(msg proto.Message) ([]byte, error) {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to serialize the proto: %v", err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return bytes, nil
}
