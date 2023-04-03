/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Author: shreyash.turkar@nutanix.com
 *
 */

package authz

import (
	"context"

	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/authz/authz_cache"
)

type AuthzInterface interface {
	GetAuthorizedUuids(ctx context.Context, iamClient authz_cache.IamClientIfc, cpdbClient cpdb.CPDBClientInterface,
		iamObjectType string, idfKind string, operation string) (*Authorized, error)
}
