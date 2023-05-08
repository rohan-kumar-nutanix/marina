/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Authors: rajesh.battala@nutanix.com
 *
 */

package scanner_config

import (
	"github.com/nutanix-core/acs-aos-go/ergon"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	"github.com/nutanix-core/content-management-marina/task/base"
)

type MarinaBaseWarehouseTaskInterface interface {
	AddAnEntity(entityType *ergon.EntityId_Entity, entityUuid *uuid4.Uuid)
	AddScannerConfigEntity()
	GetScannerConfigUuid() *uuid4.Uuid
	SetScannerConfigUuid(uuid *uuid4.Uuid)
}

type MarinaBaseScannerConfigTask struct {
	*base.MarinaBaseTask
	ScannerConfigUuid *uuid4.Uuid
	IScannerConfigDB
}

func NewMarinaBaseScannerConfigTask(
	marinaBaseTask *base.MarinaBaseTask) *MarinaBaseScannerConfigTask {
	return &MarinaBaseScannerConfigTask{
		MarinaBaseTask:   marinaBaseTask,
		IScannerConfigDB: newScannerConfigDBImpl(),
	}
}

func (t *MarinaBaseScannerConfigTask) GetScannerConfigUuid() *uuid4.Uuid {
	return t.ScannerConfigUuid
}

func (t *MarinaBaseScannerConfigTask) AddAnEntity(entityType *ergon.EntityId_Entity, entityUuid *uuid4.Uuid) {
	// TODO: Add ScannerConfig Entity in Ergon.
	entity := &ergon.EntityId{
		EntityType: ergon.EntityId_kCatalogItem.Enum(),
		EntityId:   entityUuid.RawBytes(),
	}
	t.AddEntity(t.Proto(), entity)
}

func (t *MarinaBaseScannerConfigTask) AddScannerConfigEntity() {
	// TODO: Add ScannerConfig to Ergon Types.
	t.AddAnEntity(ergon.EntityId_kUnknown.Enum(), t.GetScannerConfigUuid())
}

func (t *MarinaBaseScannerConfigTask) SetScannerConfigUuid(uuid *uuid4.Uuid) {
	t.ScannerConfigUuid = uuid
}

func (t *MarinaBaseScannerConfigTask) RecoverHook() error {
	return nil
}
