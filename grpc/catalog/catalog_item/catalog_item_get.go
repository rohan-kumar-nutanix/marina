/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 * The underlying grpc server that exports the various services. Services may
 * be added in the registerServices() implementation.
 */
package catalog_item

import (
	"context"
	"errors"
	"fmt"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	. "github.com/nutanix-core/acs-aos-go/insights/insights_interface/query"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	common "github.com/nutanix-core/content-management-marina/common"
	marinaError "github.com/nutanix-core/content-management-marina/error"
	marina_pb "github.com/nutanix-core/content-management-marina/protos/marina"
	util "github.com/nutanix-core/content-management-marina/util"
	"github.com/nutanix-core/content-management-marina/util/idf"
	"github.com/nutanix-core/ntnx-api-utils-go/tracer"
	log "k8s.io/klog/v2"
)

func validateCatalogItemGetArg(arg *marina_pb.CatalogItemGetArg) marinaError.MarinaErrorInterface {
	for _, catalogItemId := range arg.CatalogItemIdList {
		if err := common.ValidateUUID(catalogItemId.GlobalCatalogItemUuid, "GlobalCatalogItem"); err != nil {
			return err
		}
	}
	return nil
}

func GetCatalogItems(ctx context.Context, arg *marina_pb.CatalogItemGetArg) (*marina_pb.CatalogItemGetRet, error) {
	// TODO: Add Tracer support.
	// TODO: Add Authz support.
	// TODO: leverage channel and go routine for fetching results.

	span, ctx := tracer.StartSpan(ctx, "GetCatalogItems")
	defer span.Finish()

	ret := &marina_pb.CatalogItemGetRet{}
	if err := validateCatalogItemGetArg(arg); err != nil {
		log.Error("Error occured : ", err)
		return nil, util.NewGrpcStatusUtil().BuildGrpcError(err)
	}

	var expTwo, whereClause *insights_interface.BooleanExpression
	var predicateList []*insights_interface.BooleanExpression = nil
	var orderByColumn string
	var query *insights_interface.Query
	whereClause = nil
	if len(arg.GetCatalogItemIdList()) > 0 {
		var citemUuidListHex []string
		for _, catalogItemId := range arg.CatalogItemIdList {
			str := uuid4.ToUuid4(catalogItemId.GetGlobalCatalogItemUuid()).UuidToString()
			citemUuidListHex = append(citemUuidListHex, str)
			if catalogItemId.Version != nil {
				predicateList = append(predicateList, AND(EQ(COL(idf.GlobalCatalogItemUuid), STR(str)),
					EQ(COL(idf.CatalogVersion), INT64(*catalogItemId.Version))))
			} else {
				predicateList = append(predicateList, EQ(COL(idf.GlobalCatalogItemUuid), STR(str)))
			}
		}
	}

	if len(arg.GetCatalogItemTypeList()) > 0 {
		var citemTypes []string
		for _, catalogItemType := range arg.GetCatalogItemTypeList() {
			log.Infof("Item type %s ", catalogItemType.Enum().String())
			citemTypes = append(citemTypes, catalogItemType.Enum().String())
		}
		expTwo = IN(COL(idf.CatalogItemType), STR_LIST(citemTypes...))

		if predicateList == nil {
			predicateList = append(predicateList, expTwo)
			whereClause = ANY(predicateList...)
		} else {
			whereClause = AND(ANY(predicateList...), expTwo)
			log.Info("Adding guid and catalog item types.")
		}
	}

	if whereClause != nil {
		orderByColumn = idf.CatalogVersion
		query, _ = QUERY("marina_latest_catalog_item_list1").SELECT(
			idf.CatalogItemAttributes...,
		).FROM(idf.CatalogDB).
			WHERE(whereClause).
			ORDER_BY(DESCENDING(orderByColumn)).Proto()
	} else {
		query, _ = QUERY("marina_latest_catalog_item_list2").SELECT(
			idf.CatalogItemAttributes...,
		).FROM(idf.CatalogDB).Proto()
	}
	// var idfClient = idf.NewIdfClient()
	idfQueryArg := &insights_interface.GetEntitiesWithMetricsArg{Query: query}
	idfResponse, err := idf.NewIdfClient().Query(ctx, idfQueryArg)


	if err != nil {
		log.Errorf("IDF query failed because of error - %s\n", err)
		errMsg := fmt.Sprintf("Error while fetching CatalogItems list: %v", err)
		return nil, marinaError.ErrInternal.SetCause(errors.New(errMsg))
	}
	var catalogItems []*marina_pb.CatalogItemInfo

	for _, entityWithMetric := range idfResponse {
		log.Infof("Entity ID: %v", entityWithMetric.EntityGuid.GetEntityId())
		c_item := &marina_pb.CatalogItemInfo{}
		for _, metricData := range entityWithMetric.MetricDataList {

			if len(metricData.ValueList) == 0 {
				log.Infof("Could not find value for metric %v", *metricData.Name)
				continue
			}

			switch *metricData.Name {
			case idf.GlobalCatalogItemUuid:
				c_item.GlobalCatalogItemUuid = []byte(metricData.ValueList[0].Value.GetStrValue())
			case idf.CatalogItemUuid:
				c_item.Uuid = []byte(metricData.ValueList[0].Value.GetStrValue())
			case idf.CatalogVersion:
				version := metricData.ValueList[0].Value.GetInt64Value()
				c_item.Version = &version
			case idf.CatalogName:
				name := metricData.ValueList[0].Value.GetStrValue()
				c_item.Name = &name
			case idf.Annotation:
				annotation := metricData.ValueList[0].Value.GetStrValue()
				c_item.Name = &annotation
			case idf.CatalogItemType:
				c_item.ItemType = util.GetCatalogItemTypeEnum(*metricData.Name)
			}
		}
		catalogItems = append(catalogItems, c_item)
	}
	ret.CatalogItemList = catalogItems
	return ret, nil
}
