/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Authors: rajesh.battala@nutanix.com
 *
 */

package services

import (
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/ergon"
	ergonTask "github.com/nutanix-core/acs-aos-go/ergon/task"

	"github.com/nutanix-core/content-management-marina/grpc/catalog/catalog_item"
	"github.com/nutanix-core/content-management-marina/grpc/scanner_config"
	"github.com/nutanix-core/content-management-marina/grpc/security_policy"
	"github.com/nutanix-core/content-management-marina/grpc/warehouse"
	"github.com/nutanix-core/content-management-marina/task/base"
)

const (
	CatalogItemCreate              = "CatalogItemCreate"
	CatalogItemDelete              = "CatalogItemDelete"
	CatalogItemUpdate              = "CatalogItemUpdate"
	CatalogMigratePc               = "CatalogMigratePc"
	CreateWarehouse                = "CreateWarehouse"
	DeleteWarehouse                = "DeleteWarehouse"
	UpdateWarehouse                = "UpdateWarehouseMetadata"
	SyncWarehouse                  = "SyncWarehouseMetadata"
	AddItemToWarehouse             = "AddItemToWarehouse"
	DeleteWarehouseItem            = "DeleteWarehouseItem"
	UpdateWarehouseItemMetadata    = "UpdateWarehouseItemMetadata"
	CreateSecurityPolicy           = "CreateSecurityPolicy"
	DeleteSecurityPolicyByExtId    = "DeleteSecurityPolicyByExtId"
	UpdateSecurityPolicyByExtId    = "UpdateSecurityPolicyByExtId"
	CreateScannerToolConfig        = "CreateScannerToolConfig"
	DeleteScannerToolConfigByExtId = "DeleteScannerToolConfigByExtId"
	UpdateScannerToolConfigByExtId = "UpdateScannerToolConfigByExtId"
)

func GetErgonFullTaskByProto(taskProto *ergon.Task) ergonTask.FullTask {

	switch taskProto.Request.GetMethodName() {
	case CatalogItemCreate:
		return catalog_item.NewCatalogItemCreateTask(
			catalog_item.NewCatalogItemBaseTask(base.NewMarinaBaseTask(taskProto)))
	case CatalogItemDelete:
		return catalog_item.NewCatalogItemDeleteTask(
			catalog_item.NewCatalogItemBaseTask(base.NewMarinaBaseTask(taskProto)))
	case CatalogItemUpdate:
		return catalog_item.NewCatalogItemUpdateTask(
			catalog_item.NewCatalogItemBaseTask(base.NewMarinaBaseTask(taskProto)))
	case CatalogMigratePc:
		return catalog_item.NewCatalogMigratePcTask(
			catalog_item.NewCatalogItemBaseTask(base.NewMarinaBaseTask(taskProto)))
	case CreateWarehouse:
		return warehouse.NewMarinaWarehouseCreateTask(
			warehouse.NewMarinaBaseWarehouseTask(base.NewMarinaBaseTask(taskProto)))
	case DeleteWarehouse:
		return warehouse.NewMarinaWarehouseDeleteTask(
			warehouse.NewMarinaBaseWarehouseTask(base.NewMarinaBaseTask(taskProto)))
	case SyncWarehouse:
		return warehouse.NewMarinaWarehouseSyncTask(
			warehouse.NewMarinaBaseWarehouseTask(base.NewMarinaBaseTask(taskProto)))
	case UpdateWarehouse:
		return warehouse.NewMarinaWarehouseUpdateTask(
			warehouse.NewMarinaBaseWarehouseTask(base.NewMarinaBaseTask(taskProto)))
	case AddItemToWarehouse:
		return warehouse.NewMarinaWarehouseItemCreateTask(
			warehouse.NewMarinaBaseWarehouseTask(base.NewMarinaBaseTask(taskProto)))
	case DeleteWarehouseItem:
		return warehouse.NewMarinaWarehouseItemDeleteTask(
			warehouse.NewMarinaBaseWarehouseTask(base.NewMarinaBaseTask(taskProto)))
	case UpdateWarehouseItemMetadata:
		return warehouse.NewMarinaWarehouseItemUpdateTask(
			warehouse.NewMarinaBaseWarehouseTask(base.NewMarinaBaseTask(taskProto)))
	case CreateSecurityPolicy:
		return security_policy.NewMarinaSecurityPolicyCreateTask(
			security_policy.NewMarinaBaseSecurityPolicyTask(base.NewMarinaBaseTask(taskProto)))
	case DeleteSecurityPolicyByExtId:
		return security_policy.NewSecurityPolicyDeleteTask(
			security_policy.NewMarinaBaseSecurityPolicyTask(base.NewMarinaBaseTask(taskProto)))
	case UpdateSecurityPolicyByExtId:
		return security_policy.NewMarinaSecurityPolicyUpdateTask(
			security_policy.NewMarinaBaseSecurityPolicyTask(base.NewMarinaBaseTask(taskProto)))
	case CreateScannerToolConfig:
		return scanner_config.NewMarinaScannerConfigCreateTask(
			scanner_config.NewMarinaBaseScannerConfigTask(base.NewMarinaBaseTask(taskProto)))
	case DeleteScannerToolConfigByExtId:
		return scanner_config.NewMarinaScannerConfigDeleteTask(
			scanner_config.NewMarinaBaseScannerConfigTask(base.NewMarinaBaseTask(taskProto)))
	case UpdateScannerToolConfigByExtId:
		return scanner_config.NewMarinaScannerConfigUpdateTask(
			scanner_config.NewMarinaBaseScannerConfigTask(base.NewMarinaBaseTask(taskProto)))
	default:
		log.Errorf("Unknown gRPC method %s received", taskProto.Request.GetMethodName())
	}
	return nil
}
