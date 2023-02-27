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
	"bytes"
	"compress/zlib"
	"context"
	"errors"
	"fmt"

	set "github.com/deckarep/golang-set/v2"
	"github.com/golang/protobuf/proto"
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
	uuidIfc utils.UuidUtilInterface, catalogItemIds []*marinaIfc.CatalogItemId,
	catalogItemTypes []marinaIfc.CatalogItemInfo_CatalogItemType, latest bool,
	catalogItemChan chan []*marinaIfc.CatalogItemInfo, errorChan chan error) {

	catalogItems, err := catalogItem.GetCatalogItems(ctx, cpdbIfc, uuidIfc, catalogItemIds, catalogItemTypes,
		latest, "catalog_item_list")
	catalogItemChan <- catalogItems
	errorChan <- err
}

// GetCatalogItems loads Catalog Items from IDF and returns a list of CatalogItem.
// Returns ([]CatalogItem, nil) on success and (nil, error) on failure.
func (*CatalogItemImpl) GetCatalogItems(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	uuidIfc utils.UuidUtilInterface, catalogItemIds []*marinaIfc.CatalogItemId,
	catalogItemTypes []marinaIfc.CatalogItemInfo_CatalogItemType, latest bool,
	queryName string) ([]*marinaIfc.CatalogItemInfo, error) {

	var predicates []*insights_interface.BooleanExpression
	var whereClause *insights_interface.BooleanExpression
	if len(catalogItemIds) > 0 {
		for _, catalogItemId := range catalogItemIds {
			gcUuid := catalogItemId.GetGlobalCatalogItemUuid()
			if err := uuidIfc.ValidateUUID(gcUuid, "GlobalCatalogItemUUID"); err != nil {
				return nil, err
			}

			gcUuidStr := uuid4.ToUuid4(gcUuid).UuidToString()
			if catalogItemId.Version != nil {
				predicates = append(predicates, AND(EQ(COL(GlobalCatalogItemUuid), STR(gcUuidStr)),
					EQ(COL(CatalogItemVersion), INT64(*catalogItemId.Version))))
			} else {
				predicates = append(predicates, EQ(COL(GlobalCatalogItemUuid), STR(gcUuidStr)))
			}
		}
		whereClause = ANY(predicates...)
	}

	var itemTypeExpression *insights_interface.BooleanExpression
	if len(catalogItemTypes) > 0 {
		var catalogItemTypesStr []string
		for _, catalogItemType := range catalogItemTypes {
			catalogItemTypesStr = append(catalogItemTypesStr, catalogItemType.Enum().String())
		}

		itemTypeExpression = IN(COL(CatalogItemType), STR_LIST(catalogItemTypesStr...))

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
		queryBuilder = queryBuilder.SELECT(CatalogItemVersion).LIMIT(1).ORDER_BY(DESCENDING(CatalogItemVersion))
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

func (*CatalogItemImpl) GetCatalogItem(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	catalogItemUuid *uuid4.Uuid) (*marinaIfc.CatalogItemInfo, error) {

	entity, err := cpdbIfc.GetEntity(db.CatalogItem.ToString(), catalogItemUuid.String(), false)
	if err == insights_interface.ErrNotFound || entity == nil {
		log.Errorf("Catalog Item %s not found", catalogItemUuid.String())
		return nil, marinaError.ErrNotFound

	} else if err != nil {
		errMsg := fmt.Sprintf("Error encountered while getting Catalog Item %s: %v", catalogItemUuid.String(), err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))

	}

	catalogItem := &marinaIfc.CatalogItemInfo{}
	err = entity.DeserializeEntity(catalogItem)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to deserialize catalog item IDF entry: %v", err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return catalogItem, nil
}

func (*CatalogItemImpl) CreateCatalogItemFromCreateSpec(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	protoIfc utils.ProtoUtilInterface, catalogItemUuid *uuid4.Uuid, spec *marinaIfc.CatalogItemCreateSpec,
	clusters []uuid4.Uuid, sourceGroups []*marinaIfc.SourceGroup) error {

	attrParams := make(cpdb.AttributeParams)
	attrParams["uuid"] = catalogItemUuid.String()
	attrParams["name"] = spec.GetName()
	attrParams["annotation"] = spec.GetAnnotation()
	attrParams["item_type"] = spec.GetItemType().String()
	attrParams["version"] = spec.GetVersion()

	catalogItemType := spec.GetItemType()
	catalogItem := &marinaIfc.CatalogItemInfo{
		Name:       proto.String(spec.GetName()),
		Annotation: proto.String(spec.GetAnnotation()),
		ItemType:   &catalogItemType,
		Version:    proto.Int64(spec.GetVersion()),
	}

	ownerClusterUuid := spec.GetOwnerClusterUuid()
	if ownerClusterUuid != nil {
		catalogItem.OwnerClusterUuid = uuid4.ToUuid4(ownerClusterUuid).RawBytes()
		attrParams["owner_cluster_uuid"] = uuid4.ToUuid4(ownerClusterUuid).String()
	}

	gcUuid := spec.GetGlobalCatalogItemUuid()
	if gcUuid != nil {
		catalogItem.GlobalCatalogItemUuid = uuid4.ToUuid4(gcUuid).RawBytes()
		attrParams["global_catalog_item_uuid"] = uuid4.ToUuid4(gcUuid).String()
	}

	if spec.GetOpaque() != nil {
		catalogItem.Opaque = spec.GetOpaque()
		attrParams["opaque"] = spec.GetOpaque()
	}

	if spec.GetCatalogVersion() != nil {
		catalogItem.CatalogVersion = &marinaIfc.CatalogVersion{
			ProductName:    proto.String(spec.GetCatalogVersion().GetProductName()),
			ProductVersion: proto.String(spec.GetCatalogVersion().GetProductVersion()),
		}
	}

	if spec.GetSourceCatalogItemId() != nil {
		catalogItem.SourceCatalogItemId = &marinaIfc.CatalogItemId{
			GlobalCatalogItemUuid: spec.GetSourceCatalogItemId().GetGlobalCatalogItemUuid(),
			Version:               spec.GetSourceCatalogItemId().Version,
		}

	} else {
		catalogItem.SourceCatalogItemId = &marinaIfc.CatalogItemId{
			GlobalCatalogItemUuid: catalogItem.GetGlobalCatalogItemUuid(),
			Version:               catalogItem.Version,
		}
	}

	for i := range clusters {
		location := &marinaIfc.CatalogItemInfo_CatalogItemLocation{
			ClusterUuid: clusters[i].RawBytes(),
		}
		catalogItem.LocationList = append(catalogItem.LocationList, location)
	}

	catalogItem.SourceGroupList = sourceGroups

	marshal, err := protoIfc.Marshal(catalogItem)
	if err != nil {
		errMsg := fmt.Sprintf("Error marshaling the catalog item proto: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	var buffer bytes.Buffer
	writer := zlib.NewWriter(&buffer)
	_, err = writer.Write(marshal)
	if err != nil {
		errMsg := fmt.Sprintf("Error compressing the catalog item proto: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	err = writer.Close()
	if err != nil {
		errMsg := fmt.Sprintf("Error encountered while closing zlib stream: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	attrParams[insights_interface.COMPRESSED_PROTOBUF_ATTR] = buffer.Bytes()
	attrVals, err := cpdbIfc.BuildAttributeDataArgs(&attrParams)
	if err != nil {
		errMsg := fmt.Sprintf("Error encountered while generating IDF attributes: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	_, err = cpdbIfc.UpdateEntity(db.CatalogItem.ToString(), catalogItemUuid.String(), attrVals, nil, false, 0)
	if err != nil {
		errMsg := fmt.Sprintf("Error while creating the IDF entry for catalog item %s: %v", catalogItemUuid.String(), err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}

func (*CatalogItemImpl) CreateCatalogItemFromUpdateSpec(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	protoIfc utils.ProtoUtilInterface, catalogItem *marinaIfc.CatalogItemInfo,
	catalogItemUuid *uuid4.Uuid, spec *marinaIfc.CatalogItemUpdateSpec, clusters []uuid4.Uuid,
	sourceGroups []*marinaIfc.SourceGroup) error {

	attrParams := make(cpdb.AttributeParams)
	attrParams["uuid"] = catalogItemUuid.String()
	attrParams["name"] = spec.GetName()
	attrParams["annotation"] = spec.GetAnnotation()
	attrParams["item_type"] = spec.GetItemType().String()

	newCatalogItem := proto.Clone(catalogItem).(*marinaIfc.CatalogItemInfo)
	newCatalogItem.Uuid = catalogItemUuid.RawBytes()
	newCatalogItem.Name = proto.String(spec.GetName())
	newCatalogItem.Annotation = proto.String(spec.GetAnnotation())
	itemType := marinaIfc.CatalogItemInfo_CatalogItemType(
		marinaIfc.CatalogItemInfo_CatalogItemType_value[spec.GetItemType().String()])
	newCatalogItem.ItemType = &itemType

	attrParams["global_catalog_item_uuid"] = uuid4.ToUuid4(newCatalogItem.GetGlobalCatalogItemUuid()).String()
	attrParams["owner_cluster_uuid"] = uuid4.ToUuid4(newCatalogItem.GetOwnerClusterUuid()).String()

	if spec.Version != nil {
		attrParams["version"] = spec.GetVersion()
		newCatalogItem.Version = spec.Version

	} else {
		attrParams["version"] = newCatalogItem.GetVersion() + 1
		newCatalogItem.Version = proto.Int64(newCatalogItem.GetVersion() + 1)
	}

	if spec.GetOpaque() != nil {
		newCatalogItem.Opaque = spec.GetOpaque()
		attrParams["opaque"] = spec.GetOpaque()
	}

	if spec.GetCatalogVersion() != nil {
		newCatalogItem.CatalogVersion = &marinaIfc.CatalogVersion{
			ProductName:    proto.String(spec.GetCatalogVersion().GetProductName()),
			ProductVersion: proto.String(spec.GetCatalogVersion().GetProductVersion()),
		}
	}

	newCatalogItem.LocationList = nil
	for i := range clusters {
		location := &marinaIfc.CatalogItemInfo_CatalogItemLocation{
			ClusterUuid: clusters[i].RawBytes(),
		}
		newCatalogItem.LocationList = append(newCatalogItem.LocationList, location)
	}

	newCatalogItem.SourceCatalogItemId = &marinaIfc.CatalogItemId{
		GlobalCatalogItemUuid: newCatalogItem.GlobalCatalogItemUuid,
		Version:               newCatalogItem.Version,
	}
	newCatalogItem.SourceGroupList = sourceGroups

	if len(spec.RemoveFileUuidList) > 0 {
		removeFileUuids := set.NewSet[uuid4.Uuid]()
		for _, fileUuid := range spec.RemoveFileUuidList {
			removeFileUuids.Add(*uuid4.ToUuid4(fileUuid))
		}

		var newSourceGroups []*marinaIfc.SourceGroup
		for _, sourceGroup := range newCatalogItem.SourceGroupList {
			var newSources []*marinaIfc.Source
			for _, source := range sourceGroup.SourceList {
				if !removeFileUuids.Contains(*uuid4.ToUuid4(source.FileUuid)) {
					newSources = append(newSources, source)
				}
			}
			sourceGroup.SourceList = newSources

			if len(sourceGroup.SourceList) > 0 {
				newSourceGroups = append(newSourceGroups, sourceGroup)
			}
		}
		newCatalogItem.SourceGroupList = newSourceGroups
	}

	marshal, err := protoIfc.Marshal(newCatalogItem)
	if err != nil {
		errMsg := fmt.Sprintf("Error marshaling the catalog item proto: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	var buffer bytes.Buffer
	writer := zlib.NewWriter(&buffer)
	_, err = writer.Write(marshal)
	if err != nil {
		errMsg := fmt.Sprintf("Error compressing the catalog item proto: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	err = writer.Close()
	if err != nil {
		errMsg := fmt.Sprintf("Error encountered while closing zlib stream: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	attrParams[insights_interface.COMPRESSED_PROTOBUF_ATTR] = buffer.Bytes()
	attrVals, err := cpdbIfc.BuildAttributeDataArgs(&attrParams)
	if err != nil {
		errMsg := fmt.Sprintf("Error encountered while generating IDF attributes: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	_, err = cpdbIfc.UpdateEntity(db.CatalogItem.ToString(), catalogItemUuid.String(), attrVals, nil, false, 0)
	if err != nil {
		errMsg := fmt.Sprintf("Error while creating the IDF entry for catalog item %s: %v", catalogItemUuid.String(), err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}
