/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 * Catalog Item interface
 *
 */

package catalog_item

import (
	"context"

	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	utils "github.com/nutanix-core/content-management-marina/util"
)

type CatalogItemInterface interface {
	GetCatalogItemsChan(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, uuidIfc utils.UuidUtilInterface,
		catalogItemIds []*marinaIfc.CatalogItemId, catalogItemTypes []marinaIfc.CatalogItemInfo_CatalogItemType,
		latest bool, catalogItemChan chan []*marinaIfc.CatalogItemInfo, errorChan chan error)
	GetCatalogItems(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, uuidIfc utils.UuidUtilInterface,
		catalogItemIds []*marinaIfc.CatalogItemId, catalogItemTypes []marinaIfc.CatalogItemInfo_CatalogItemType,
		latest bool, queryName string) ([]*marinaIfc.CatalogItemInfo, error)
	GetCatalogItem(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, catalogItemUuid *uuid4.Uuid) (
		*marinaIfc.CatalogItemInfo, error)
	CreateCatalogItemFromCreateSpec(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
		protoIfc utils.ProtoUtilInterface, catalogItemUuid *uuid4.Uuid, spec *marinaIfc.CatalogItemCreateSpec,
		clusters []uuid4.Uuid, sourceGroups []*marinaIfc.SourceGroup) error
	CreateCatalogItemFromUpdateSpec(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
		protoIfc utils.ProtoUtilInterface, catalogItem *marinaIfc.CatalogItemInfo, catalogItemUuid *uuid4.Uuid, spec *marinaIfc.CatalogItemUpdateSpec,
		clusters []uuid4.Uuid, sourceGroups []*marinaIfc.SourceGroup) error
}
