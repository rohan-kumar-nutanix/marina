/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Authors: rajesh.battala@nutanix.com
*
* This is the implementation of the Warehouse Base task for Warehouse Entity.
 */

package warehouse

import (
	"github.com/nutanix-core/acs-aos-go/ergon"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	"github.com/nutanix-core/content-management-marina/task/base"
)

type MarinaBaseWarehouseTaskInterface interface {
	AddAnEntity(entityType *ergon.EntityId_Entity, entityUuid *uuid4.Uuid)
	AddWarehouseEntity()
	GetWarehouseUuid() *uuid4.Uuid
	SetLogicalTimestamp(*int64)
	GetLogicalTimestamp() *int64
	SetWarehouseUuid(uuid *uuid4.Uuid)
}

type MarinaBaseWarehouseTask struct {
	*base.MarinaBaseTask
	WarehouseUuid             *uuid4.Uuid
	WarehouseLogicalTimestamp *int64
	IWarehouseDB
}

func NewMarinaBaseWarehouseTask(
	marinaBaseTask *base.MarinaBaseTask) *MarinaBaseWarehouseTask {
	return &MarinaBaseWarehouseTask{
		MarinaBaseTask: marinaBaseTask,
		IWarehouseDB:   &WarehouseDBImpl{},
	}
}

func (t *MarinaBaseWarehouseTask) GetWarehouseUuid() *uuid4.Uuid {
	return t.WarehouseUuid
}

func (t *MarinaBaseWarehouseTask) AddAnEntity(entityType *ergon.EntityId_Entity, entityUuid *uuid4.Uuid) {
	// TODO: Add Warehouse Entity in Ergon.
	entity := &ergon.EntityId{
		EntityType: ergon.EntityId_kCatalogItem.Enum(),
		EntityId:   entityUuid.RawBytes(),
	}
	t.AddEntity(t.Proto(), entity)
}

func (t *MarinaBaseWarehouseTask) AddWarehouseEntity() {
	// TODO: Add Warehouse Entity to Ergon Types.
	t.AddAnEntity(ergon.EntityId_kUnknown.Enum(), t.GetWarehouseUuid())
}

func (t *MarinaBaseWarehouseTask) SetLogicalTimestamp(warehouseLogicalTimestamp *int64) {
	t.WarehouseLogicalTimestamp = warehouseLogicalTimestamp
}

func (t *MarinaBaseWarehouseTask) GetLogicalTimestamp() *int64 {
	return t.WarehouseLogicalTimestamp
}

func (t *MarinaBaseWarehouseTask) SetWarehouseUuid(uuid *uuid4.Uuid) {
	t.WarehouseUuid = uuid
}

/*func (task *MarinaBaseWarehouseTask) RecoverHook() error {
	return nil
}*/

/*func newMarinaBaseWarehouseTaskUtil(
	marinaBaseTask *base.MarinaBaseTask) *MarinaBaseWarehouseTask {
	return &MarinaBaseWarehouseTask{
		MarinaBaseTask: marinaBaseTask,
	}
}*/
