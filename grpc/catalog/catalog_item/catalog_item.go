/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 * Wrapper around the Catalog Item DB entry. Includes libraries that will
 * interact with IDF and query for Catalog Items.
 */

package catalog_item

import (
	"context"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	. "github.com/nutanix-core/acs-aos-go/insights/insights_interface/query"
	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	utils "github.com/nutanix-core/content-management-marina/util"
)

const (
	CatalogItemTable      = "catalog_item_info"
	GlobalCatalogItemUuid = "global_catalog_item_uuid"
	CatalogItemType       = "item_type"
	CatalogVersion        = "version"
)

var catalogItemAttributes = []interface{}{
	insights_interface.COMPRESSED_PROTOBUF_ATTR,
}

type CatalogItemImpl struct {
}

// GetCatalogItemsChan pushes catalog items and error object to respective channels.
func (catalogItem *CatalogItemImpl) GetCatalogItemsChan(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	catalogItemIdList []*marinaIfc.CatalogItemId, catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType,
	latest bool, catalogItemChan chan []*marinaIfc.CatalogItemInfo, errorChan chan error) {
	catalogItemList, err := catalogItem.GetCatalogItems(ctx, cpdbIfc, catalogItemIdList, catalogItemTypeList, latest, "catalog_item_list")
	catalogItemChan <- catalogItemList
	errorChan <- err
}

// GetCatalogItems loads Catalog Items from IDF and returns a list of CatalogItem.
// Returns ([]CatalogItem, nil) on success and (nil, error) on failure.
func (*CatalogItemImpl) GetCatalogItems(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	catalogItemIdList []*marinaIfc.CatalogItemId, catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType,
	latest bool, queryName string) ([]*marinaIfc.CatalogItemInfo, error) {
	var predicateList []*insights_interface.BooleanExpression
	var whereClause *insights_interface.BooleanExpression
	if len(catalogItemIdList) > 0 {
		var catalogItemUuidStrList []string
		for _, catalogItemId := range catalogItemIdList {
			gcUuid := catalogItemId.GetGlobalCatalogItemUuid()

			if err := utils.ValidateUUID(gcUuid, "GlobalCatalogItemUUID"); err != nil {
				return nil, err
			}

			gcUuidStr := uuid4.ToUuid4(gcUuid).UuidToString()
			catalogItemUuidStrList = append(catalogItemUuidStrList, gcUuidStr)

			if catalogItemId.Version != nil {
				predicateList = append(predicateList, AND(EQ(COL(GlobalCatalogItemUuid), STR(gcUuidStr)),
					EQ(COL(CatalogVersion), INT64(*catalogItemId.Version))))
			} else {
				predicateList = append(predicateList, EQ(COL(GlobalCatalogItemUuid), STR(gcUuidStr)))
			}
		}
		whereClause = ANY(predicateList...)
	}

	var itemTypeExpression *insights_interface.BooleanExpression
	if len(catalogItemTypeList) > 0 {
		var catalogItemTypeStrList []string
		for _, catalogItemType := range catalogItemTypeList {
			catalogItemTypeStrList = append(catalogItemTypeStrList, catalogItemType.Enum().String())
		}

		itemTypeExpression = IN(COL(CatalogItemType), STR_LIST(catalogItemTypeStrList...))

		if whereClause == nil {
			whereClause = itemTypeExpression
		} else {
			whereClause = AND(whereClause, itemTypeExpression)
		}
	}

	queryBuilder := QUERY(queryName).SELECT(catalogItemAttributes...).FROM(CatalogItemTable)
	if whereClause != nil {
		queryBuilder = queryBuilder.WHERE(whereClause)
	}

	if latest {
		queryBuilder = queryBuilder.LIMIT(1).ORDER_BY(DESCENDING(CatalogVersion))
	}

	query, err := queryBuilder.Proto()
	if err != nil {
		return nil, err
	}

	arg := insights_interface.GetEntitiesWithMetricsArg{Query: query}
	idfResponse, err := cpdbIfc.Query(&arg)

	var catalogItems []*marinaIfc.CatalogItemInfo
	if err == insights_interface.ErrNotFound {
		return catalogItems, nil
	} else if err != nil {
		return nil, err
	}

	for _, entityWithMetric := range idfResponse {
		catalogItem := &marinaIfc.CatalogItemInfo{}
		err = entityWithMetric.DeserializeEntity(catalogItem)
		if err != nil {
			return nil, err
		}
		catalogItems = append(catalogItems, catalogItem)
	}

	return catalogItems, nil
}
