/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: shreyash.turkar@nutanix.com
 *
 * Metadata interface implementation for Marina.
 */

package metadata

import (
	"context"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	. "github.com/nutanix-core/acs-aos-go/insights/insights_interface/query"
	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/content-management-marina/protos/marina"
)

const (
	tableNames           = "abac_entity_capability"
	kind                 = "kind"
	kindID               = "kind_id"
	ownerReference       = "owner_reference"
	ownerUsername        = "owner_username"
	createTimeUsecs      = "create_time_usecs"
	lastUpdatedTimeUsecs = "last_updated_time_usecs"
	categoryIdList       = "category_id_list"
)

var entityMetadataAttribute = []interface{}{
	kindID,
	ownerReference,
	ownerUsername,
	createTimeUsecs,
	lastUpdatedTimeUsecs,
	categoryIdList,
}

type EntityMetadataByUuid map[uuid4.Uuid]*marina.EntityMetadata

type EntityMetadataUtil struct {
}

// GetEntityMetadata Query the ECap table to fetch metadata of the given entities.
// Returns map of Entity Metadata by Entity UUIDs.
func (*EntityMetadataUtil) GetEntityMetadata(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, kindIDs []string,
	idfKind string, queryName string) (EntityMetadataByUuid, error) {

	idfQuery, err := QUERY(queryName).
		SELECT(entityMetadataAttribute...).
		FROM(tableNames).
		WHERE(AND(
			EQ(COL(kind), STR(idfKind)),
			IN(COL(kindID), STR_LIST(kindIDs...)))).
		Proto()
	if err != nil {
		return nil, err
	}
	idfQueryArg := &insights_interface.GetEntitiesWithMetricsArg{Query: idfQuery}
	idfResponse, err := cpdbIfc.Query(idfQueryArg)
	if err != nil {
		return nil, err
	}

	metadataByUuid := EntityMetadataByUuid{}
	for _, entityWithMetric := range idfResponse {
		entityMetadata := &marina.EntityMetadata{}
		var entityUuid uuid4.Uuid
		for _, metricData := range entityWithMetric.MetricDataList {
			value := metricData.ValueList[0].Value
			switch *metricData.Name {
			case kindID:
				entityUuid = *uuid4.ToUuid4(value.GetBytesValue())
			case createTimeUsecs:
				createTime := value.GetInt64Value()
				entityMetadata.CreateTimeUsecs = &createTime
			case lastUpdatedTimeUsecs:
				lastUpdateTime := value.GetInt64Value()
				entityMetadata.LastUpdateTimeUsecs = &lastUpdateTime
			case ownerReference:
				entityMetadata.OwnerUserUuid = value.GetBytesValue()
			case ownerUsername:
				userName := value.GetStrValue()
				entityMetadata.OwnerUserName = &userName
			case categoryIdList:
				entityMetadata.CategoriesUuidList = value.GetBytesList().GetValueList()
			}
		}
		metadataByUuid[entityUuid] = entityMetadata
	}
	return metadataByUuid, nil
}
