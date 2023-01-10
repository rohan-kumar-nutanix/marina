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
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	utils "github.com/nutanix-core/content-management-marina/util"
)

type CatalogItemInterface interface {
	GetCatalogItemsChan(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, uuidIfc utils.UuidUtilInterface,
		catalogItemIdList []*marinaIfc.CatalogItemId, catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType,
		latest bool, catalogItemChan chan []*marinaIfc.CatalogItemInfo, errorChan chan error)
	GetCatalogItems(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, cuuidIfc utils.UuidUtilInterface,
		atalogItemIdList []*marinaIfc.CatalogItemId, catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType,
		latest bool, queryName string) ([]*marinaIfc.CatalogItemInfo, error)
}
