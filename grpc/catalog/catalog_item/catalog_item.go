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
	"errors"
	"fmt"

	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	. "github.com/nutanix-core/acs-aos-go/insights/insights_interface/query"
	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/content-management-marina/db"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	utils "github.com/nutanix-core/content-management-marina/util"
)

const (
	GlobalCatalogItemUuid = "global_catalog_item_uuid"
	CatalogItemType       = "item_type"
	CatalogItemVersion    = "version"
)

var catalogItemAttributes = []interface{}{
	insights_interface.COMPRESSED_PROTOBUF_ATTR,
}

type CatalogItemImpl struct {
}

// GetCatalogItemsChan pushes catalog items and error object to respective channels.
func (catalogItem *CatalogItemImpl) GetCatalogItemsChan(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	uuidIfc utils.UuidUtilInterface, catalogItemIdList []*marinaIfc.CatalogItemId,
	catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType, latest bool,
	catalogItemChan chan []*marinaIfc.CatalogItemInfo, errorChan chan error) {

	catalogItemList, err := catalogItem.GetCatalogItems(ctx, cpdbIfc, uuidIfc, catalogItemIdList, catalogItemTypeList,
		latest, "catalog_item_list")
	catalogItemChan <- catalogItemList
	errorChan <- err
}

// GetCatalogItems loads Catalog Items from IDF and returns a list of CatalogItem.
// Returns ([]CatalogItem, nil) on success and (nil, error) on failure.
func (*CatalogItemImpl) GetCatalogItems(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	uuidIfc utils.UuidUtilInterface, catalogItemIdList []*marinaIfc.CatalogItemId,
	catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType, latest bool,
	queryName string) ([]*marinaIfc.CatalogItemInfo, error) {

	var predicateList []*insights_interface.BooleanExpression
	var whereClause *insights_interface.BooleanExpression
	if len(catalogItemIdList) > 0 {
		var catalogItemUuidStrList []string
		for _, catalogItemId := range catalogItemIdList {
			gcUuid := catalogItemId.GetGlobalCatalogItemUuid()

			gcUuidStr := uuid4.ToUuid4(gcUuid).UuidToString()
			if err := uuidIfc.ValidateUUID(gcUuid, "GlobalCatalogItemUUID"); err != nil {
				log.Errorf("Invalid UUID: %s", gcUuidStr)
				return nil, err
			}

			catalogItemUuidStrList = append(catalogItemUuidStrList, gcUuidStr)

			if catalogItemId.Version != nil {
				predicateList = append(predicateList, AND(EQ(COL(GlobalCatalogItemUuid), STR(gcUuidStr)),
					EQ(COL(CatalogItemVersion), INT64(*catalogItemId.Version))))
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

	queryBuilder := QUERY(queryName).SELECT(catalogItemAttributes...).FROM(db.CatalogItem.ToString())
	if whereClause != nil {
		queryBuilder = queryBuilder.WHERE(whereClause)
	}

	if latest {
		queryBuilder = queryBuilder.LIMIT(1).ORDER_BY(DESCENDING(CatalogItemVersion))
	}

	query, err := queryBuilder.Proto()
	if err != nil {
		errMsg := fmt.Sprintf("Error encountered while building the IDF query %s: %v", queryName, err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	arg := insights_interface.GetEntitiesWithMetricsArg{Query: query}
	idfResponse, err := cpdbIfc.Query(&arg)

	var catalogItems []*marinaIfc.CatalogItemInfo
	if err == insights_interface.ErrNotFound {
		return catalogItems, nil
	} else if err != nil {
		errMsg := fmt.Sprintf("Failed to fetch the catalog item(s): %v", err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	for _, entityWithMetric := range idfResponse {
		catalogItem := &marinaIfc.CatalogItemInfo{}
		err = entityWithMetric.DeserializeEntity(catalogItem)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to deserialize catalog item IDF entry: %v", err)
			return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
		}
		catalogItems = append(catalogItems, catalogItem)
	}

	return catalogItems, nil
}
