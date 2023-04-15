/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Authors: rajesh.battala@nutanix.com
 *
 */

package warehouse

import (
	"context"

	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	"github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"
	utils "github.com/nutanix-core/content-management-marina/util"
)

type IWarehouseDB interface {
	CreateWarehouse(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
		protoIfc utils.ProtoUtilInterface, warehouseUuid *uuid4.Uuid, warehousePB *content.Warehouse) error
}

type IWarehouseStorage interface {
	CreateWarehouseBucket(ctx context.Context, warehouseUuid *uuid4.Uuid)
}
