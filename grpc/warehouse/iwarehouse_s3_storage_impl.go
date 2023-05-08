/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Authors: rajesh.battala@nutanix.com
 *
 */

package warehouse

import (
	"context"
	"sync"

	log "k8s.io/klog/v2"

	"github.com/nutanix-core/content-management-marina/storageadapters"
)

var (
	warehouseS3StorageImpl IWarehouseStorage = nil
	clientOnce             sync.Once
)

type WarehouseS3StorageImpl struct {
	storageadapters.AwsS3Impl
}

func newWarehouseS3StorageImpl() IWarehouseStorage {
	clientOnce.Do(func() {
		warehouseS3StorageImpl = &WarehouseS3StorageImpl{}
	})
	return warehouseS3StorageImpl
}

func (storageImpl WarehouseS3StorageImpl) CreateWarehouseBucket(ctx context.Context, warehouseUuid string) error {
	log.Infof("Creating Bucket For Warehouse %s", warehouseUuid)

	err := storageImpl.AwsS3Impl.CreateWarehouseBucket(ctx, warehouseUuid)
	if err != nil {
		log.Errorf("Error Occurred while the Creating Warehouse Bucket %s", err)

	}
	return err
}

func (storageImpl WarehouseS3StorageImpl) DeleteWarehouseBucket(ctx context.Context, warehouseUuid string) error {
	
	if len(warehouseUuid) == 0 {
		log.Errorf("Warehouse UUID is empty")
		return nil
	}
	err := storageImpl.AwsS3Impl.DeleteWarehouseBucket(ctx, warehouseUuid)
	if err != nil {
		log.Errorf("Error Occurred while Deleting the Warehouse Bucket %s", err)

	}
	return err
}
