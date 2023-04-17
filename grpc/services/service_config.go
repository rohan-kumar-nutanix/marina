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
	"github.com/nutanix-core/content-management-marina/grpc/warehouse"
	"github.com/nutanix-core/content-management-marina/task/base"
)

const (
	CatalogItemCreate = "CatalogItemCreate"
	CatalogItemDelete = "CatalogItemDelete"
	CatalogItemUpdate = "CatalogItemUpdate"
	CatalogMigratePc  = "CatalogMigratePc"
	CreateWarehouse   = "CreateWarehouse"
	DeleteWarehouse   = "DeleteWarehouse"
	UpdateWarehouse   = "UpdateWarehouseMetadata"
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
	case UpdateWarehouse:
		return warehouse.NewMarinaWarehouseUpdateTask(
			warehouse.NewMarinaBaseWarehouseTask(base.NewMarinaBaseTask(taskProto)))
	default:
		log.Errorf("Unknown gRPC method %s received", taskProto.Request.GetMethodName())
	}
	return nil
}
