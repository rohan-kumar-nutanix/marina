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
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	. "github.com/nutanix-core/acs-aos-go/insights/insights_interface/query"
	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/content-management-marina/db"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
)

const (
	kindCol                 = "kind"
	kindIdCol               = "kind_id"
	ownerReferenceCol       = "owner_reference"
	ownerUsernameCol        = "owner_username"
	createTimeUsecsCol      = "create_time_usecs"
	lastUpdatedTimeUsecsCol = "last_updated_time_usecs"
	categoryIdListCol       = "category_id_list"
)

var entityMetadataAttribute = []interface{}{
	kindIdCol,
	ownerReferenceCol,
	ownerUsernameCol,
	createTimeUsecsCol,
	lastUpdatedTimeUsecsCol,
	categoryIdListCol,
}

type EntityMetadataUtil struct {
}

// GetEntityMetadata Query the ECap table to fetch metadata of the given entities.
// Returns map of Entity Metadata by Entity UUIDs.
func (*EntityMetadataUtil) GetEntityMetadata(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, kindIDs []string,
	idfKind string, queryName string) (map[uuid4.Uuid]*marinaIfc.EntityMetadata, error) {

	query, err := QUERY(queryName).
		SELECT(entityMetadataAttribute...).
		FROM(db.EntityCapability.ToString()).
		WHERE(AND(
			EQ(COL(kindCol), STR(idfKind)),
			IN(COL(kindIdCol), STR_LIST(kindIDs...)))).
		Proto()
	if err != nil {
		errMsg := fmt.Sprintf("Failed to create the IDF query %s: %v", queryName, err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	arg := &insights_interface.GetEntitiesWithMetricsArg{Query: query}
	entities, err := cpdbIfc.Query(arg)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to fetch the ecap entry for %s: %v", idfKind, err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	metadataByUuid := make(map[uuid4.Uuid]*marinaIfc.EntityMetadata)
	for _, entityWithMetric := range entities {
		var entityUuid uuid4.Uuid
		entityMetadata := &marinaIfc.EntityMetadata{}
		for _, metricData := range entityWithMetric.MetricDataList {
			if len(metricData.ValueList) == 0 {
				continue
			}

			value := metricData.ValueList[0].Value
			switch *metricData.Name {
			case kindIdCol:
				uuid, err := uuid4.StringToUuid4(value.GetStrValue())
				if err != nil {
					errMsg := fmt.Sprintf("Failed to convert string %s to UUID", value.GetStrValue())
					return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
				}
				entityUuid = *uuid

			case createTimeUsecsCol:
				entityMetadata.CreateTimeUsecs = proto.Int64(int64(value.GetUint64Value()))

			case lastUpdatedTimeUsecsCol:
				entityMetadata.LastUpdateTimeUsecs = proto.Int64(int64(value.GetUint64Value()))

			case ownerReferenceCol:
				uuid, err := uuid4.StringToUuid4(value.GetStrValue())
				if err != nil {
					errMsg := fmt.Sprintf("Failed to convert string %s to UUID", value.GetStrValue())
					return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
				}
				entityMetadata.OwnerUserUuid = uuid.RawBytes()

			case ownerUsernameCol:
				entityMetadata.OwnerUserName = proto.String(value.GetStrValue())

			case categoryIdListCol:
				for _, uuidStr := range value.GetStrList().GetValueList() {
					uuid, err := uuid4.StringToUuid4(uuidStr)
					if err != nil {
						errMsg := fmt.Sprintf("Failed to convert string %s to UUID", value.GetStrValue())
						return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
					}
					entityMetadata.CategoriesUuidList = append(entityMetadata.CategoriesUuidList, uuid.RawBytes())
				}
			}
		}
		metadataByUuid[entityUuid] = entityMetadata
	}
	return metadataByUuid, nil
}
