/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: shreyash.turkar@nutanix.com
 *
 */

package metadata

import (
	"context"

	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
)

type EntityMetadataInterface interface {
	GetEntityMetadata(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, kindIDs []string, idfKind string,
		queryName string) (map[uuid4.Uuid]*marinaIfc.EntityMetadata, error)
}
