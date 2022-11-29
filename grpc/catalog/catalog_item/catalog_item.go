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

const (
	CatalogItemTable      = "catalog_item_info"
	Annotation            = "annotation"
	GlobalCatalogItemUuid = "global_catalog_item_uuid"
	CatalogItemType       = "item_type"
	CatalogName           = "name"
	CatalogItemUuid       = "uuid"
	CatalogVersion        = "version"
)

var CatalogItemAttributes = []interface{}{
	Annotation,
	GlobalCatalogItemUuid,
	CatalogItemType,
	CatalogName,
	CatalogItemUuid,
	CatalogVersion,
}

// GetCatalogItemsChan pushes catalog items and error object to respective channels.
func GetCatalogItemsChan(ctx context.Context, catalogItemIdList []*marinaIfc.CatalogItemId,
	catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType, latest bool,
	catalogItemChan chan []*marinaIfc.CatalogItemInfo, errorChan chan error) {
	catalogItemList, err := getCatalogItems(ctx, catalogItemIdList, catalogItemTypeList, latest)
	catalogItemChan <- catalogItemList
	errorChan <- err
}

// GetCatalogItems loads Catalog Items from IDF and returns a list of CatalogItem.
// Returns (CatalogItem, nil) on success and (nil, error) on failure.
func getCatalogItems(ctx context.Context, catalogItemIdList []*marinaIfc.CatalogItemId,
	catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType, latest bool) ([]*marinaIfc.CatalogItemInfo, error) {
	var predicateList []*insights_interface.BooleanExpression
	var whereClause *insights_interface.BooleanExpression
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

	queryBuilder := QUERY("catalog_item_list").SELECT(CatalogItemAttributes...).FROM(CatalogItemTable)
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
			case GlobalCatalogItemUuid:
				catalogItem.GlobalCatalogItemUuid = []byte(metricData.ValueList[0].Value.GetStrValue())
			case CatalogItemUuid:
				catalogItem.Uuid = []byte(metricData.ValueList[0].Value.GetStrValue())
			case CatalogVersion:
				version := metricData.ValueList[0].Value.GetInt64Value()
				catalogItem.Version = &version
			case CatalogName:
				name := metricData.ValueList[0].Value.GetStrValue()
				catalogItem.Name = &name
			case Annotation:
				annotation := metricData.ValueList[0].Value.GetStrValue()
				catalogItem.Name = &annotation
			case CatalogItemType:
				catalogItem.ItemType = util.GetCatalogItemTypeEnum(*metricData.Name)
			}
		}
		catalogItems = append(catalogItems, catalogItem)
	}

	return catalogItems, nil
}
