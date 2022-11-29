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
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	. "github.com/nutanix-core/acs-aos-go/insights/insights_interface/query"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/content-management-marina/common"
	"github.com/nutanix-core/content-management-marina/db"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	util "github.com/nutanix-core/content-management-marina/util"
)

type CatalogItem struct {
	catalogItemId *marinaIfc.CatalogItemId
	// This should be nil if 'exists' is false.
	catalogItemInfo *marinaIfc.CatalogItemInfo
	exists          bool
}

func GetCatalogItemsChan(ctx context.Context, catalogItemIdList []*marinaIfc.CatalogItemId,
	catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType, latest bool,
	catalogItemChan chan []*marinaIfc.CatalogItemInfo, errorChan chan error) {
	catalogItemList, err := getCatalogItems(ctx, catalogItemIdList, catalogItemTypeList, latest)
	catalogItemChan <- catalogItemList
	errorChan <- err
}

// GetCatalogItems Loads Catalog Items from IDF and returns a list of CatalogItem.
// Returns (CatalogItem, nil) on success and (nil, error) on failure.
func getCatalogItems(ctx context.Context, catalogItemIdList []*marinaIfc.CatalogItemId,
	catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType, latest bool) ([]*marinaIfc.CatalogItemInfo, error) {
	var predicateList []*insights_interface.BooleanExpression
	if len(catalogItemIdList) > 0 {
		var catalogItemUuidStrList []string
		for _, catalogItemId := range catalogItemIdList {
			gcUuid := catalogItemId.GetGlobalCatalogItemUuid()

			if err := common.ValidateUUID(gcUuid, "GlobalCatalogItemUUID"); err != nil {
				return nil, err
			}

			gcUuidStr := uuid4.ToUuid4(gcUuid).UuidToString()
			catalogItemUuidStrList = append(catalogItemUuidStrList, gcUuidStr)

			if catalogItemId.Version != nil {
				predicateList = append(predicateList, AND(EQ(COL(db.GlobalCatalogItemUuid), STR(gcUuidStr)),
					EQ(COL(db.CatalogVersion), INT64(*catalogItemId.Version))))
			} else {
				predicateList = append(predicateList, EQ(COL(db.GlobalCatalogItemUuid), STR(gcUuidStr)))
			}
		}
	}

	var itemTypeExpression, whereClause *insights_interface.BooleanExpression
	if len(catalogItemTypeList) > 0 {
		var catalogItemTypeStrList []string
		for _, catalogItemType := range catalogItemTypeList {
			catalogItemTypeStrList = append(catalogItemTypeStrList, catalogItemType.Enum().String())
		}

		itemTypeExpression = IN(COL(db.CatalogItemType), STR_LIST(catalogItemTypeStrList...))

		if predicateList == nil {
			predicateList = append(predicateList, itemTypeExpression)
			whereClause = ANY(predicateList...)
		} else {
			whereClause = AND(ANY(predicateList...), itemTypeExpression)
		}
	}

	queryBuilder := QUERY("latest_catalog_item_list").SELECT(db.CatalogItemAttributes...).FROM(db.CatalogItemTable)
	if whereClause != nil {
		queryBuilder = queryBuilder.WHERE(whereClause)
	}

	if latest {
		queryBuilder = queryBuilder.LIMIT(1).ORDER_BY(DESCENDING(db.CatalogVersion))
	}

	query, err := queryBuilder.Proto()
	if err != nil {
		return nil, err
	}

	idfQueryArg := &insights_interface.GetEntitiesWithMetricsArg{Query: query}
	idfResponse, err := db.IdfClientWithRetry().Query(ctx, idfQueryArg)
	if err != nil {
		return nil, err
	}

	var catalogItems []*marinaIfc.CatalogItemInfo
	for _, entityWithMetric := range idfResponse {
		catalogItem := &marinaIfc.CatalogItemInfo{}
		for _, metricData := range entityWithMetric.MetricDataList {

			if len(metricData.ValueList) == 0 {
				log.Infof("Could not find value for metric %v", *metricData.Name)
				continue
			}

			switch *metricData.Name {
			case db.GlobalCatalogItemUuid:
				catalogItem.GlobalCatalogItemUuid = []byte(metricData.ValueList[0].Value.GetStrValue())
			case db.CatalogItemUuid:
				catalogItem.Uuid = []byte(metricData.ValueList[0].Value.GetStrValue())
			case db.CatalogVersion:
				version := metricData.ValueList[0].Value.GetInt64Value()
				catalogItem.Version = &version
			case db.CatalogName:
				name := metricData.ValueList[0].Value.GetStrValue()
				catalogItem.Name = &name
			case db.Annotation:
				annotation := metricData.ValueList[0].Value.GetStrValue()
				catalogItem.Name = &annotation
			case db.CatalogItemType:
				catalogItem.ItemType = util.GetCatalogItemTypeEnum(*metricData.Name)
			}
		}
		catalogItems = append(catalogItems, catalogItem)
	}

	return catalogItems, nil
}
