/*
 *
 * Common Utils code.
 *
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: Rajesh Battala <rajesh.battala@nutanix.com>
 *
 */

package utils

import (
	"fmt"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
)

func ValidateUUID(uuidValue []byte, fieldName string) marinaError.MarinaErrorInterface {
	if err := uuid4.ToUuid4(uuidValue); err == nil {
		return marinaError.ErrMarinaInvalidUuid(string(uuidValue)).SetCauseAndLog(
			fmt.Errorf("invalid '%s' (%s). UUID must be exactly 16 bytes string",
				fieldName, string(uuidValue)))
	}
	return nil
}
