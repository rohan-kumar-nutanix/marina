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

	"github.com/nutanix-core/content-management-marina/db"
	"github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"
	utils "github.com/nutanix-core/content-management-marina/util"
)

type IWarehouseDB interface {
	CreateWarehouse(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
		protoIfc utils.ProtoUtilInterface, warehouseUuid *uuid4.Uuid, warehousePB *content.Warehouse) error
	GetWarehouse(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, warehouseUuid *uuid4.Uuid) (
		*content.Warehouse, error)
	DeleteWarehouse(ctx context.Context, idfIfc db.IdfClientInterface,
		cpdbIfc cpdb.CPDBClientInterface, warehouseUuid string) error
	UpdateWarehouse(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
		protoIfc utils.ProtoUtilInterface, warehouseUuid *uuid4.Uuid, warehousePB *content.Warehouse) error
	ListWarehouses(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface) ([]*content.Warehouse, error)

	// WarehouseItem Methods
	CreateWarehouseItem(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
		protoIfc utils.ProtoUtilInterface, warehouseUuid *uuid4.Uuid, warehouseItemUuid *uuid4.Uuid, warehouseItemPB *content.WarehouseItem) error
	GetWarehouseItem(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, warehouseItemUuid *uuid4.Uuid) (
		*content.WarehouseItem, error)
	ListWarehouseItems(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, warehouseUuid *uuid4.Uuid) ([]*content.WarehouseItem, error)
	DeleteWarehouseItem(ctx context.Context, idfIfc db.IdfClientInterface, cpdbIfc cpdb.CPDBClientInterface,
		warehouseUuid string, warehouseItemUuid string) error
	UpdateWarehouseItem(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, protoIfc utils.ProtoUtilInterface,
		warehouseUuid *uuid4.Uuid, warehouseItemUuid *uuid4.Uuid, warehouseItemPB *content.WarehouseItem) error
}

type IWarehouseStorage interface {
	CreateWarehouseBucket(ctx context.Context, warehouseUuid string) error
	DeleteWarehouseBucket(ctx context.Context, warehouseUuid string) error
	UploadFileToWarehouseBucket(ctx context.Context, warehouseUuid string, pathToFile string, data []byte) error
	DeleteFileFromWarehouseBucket(ctx context.Context, warehouseUuid string, pathToFile string) error
	DeleteAllFilesFromWarehouseBucket(ctx context.Context, warehouseUuid string) error
	UpdateFileInWarehouseBucket(ctx context.Context, warehouseUuid string, pathToFile string, data []byte) error
}
