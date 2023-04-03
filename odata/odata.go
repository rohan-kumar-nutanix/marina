/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 * Odata interface implementation for Marina.
 */

package odata

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	set "github.com/deckarep/golang-set/v2"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	. "github.com/nutanix-core/acs-aos-go/insights/insights_interface/query"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/content-management-marina/authz"
	"github.com/nutanix-core/content-management-marina/db"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
)

var uuidCol = "uuid"

type OdataUtil struct {
}

func (*OdataUtil) GetLimitAndOffset(ctx context.Context, idfParams *marinaIfc.IdfParameters) (*int64, *int64, error) {
	if idfParams.PaginationInfo == nil {
		return nil, nil, nil
	}

	limit := idfParams.GetPaginationInfo().GetPageSize()
	if limit <= 0 {
		errMsg := fmt.Sprintf("Page size %d <= 0", limit)
		return nil, nil, marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
	}

	pageNumber := idfParams.GetPaginationInfo().GetPageNumber()
	if pageNumber < 0 {
		errMsg := fmt.Sprintf("Page number %d < 0", limit)
		return nil, nil, marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
	}

	return &limit, &pageNumber, nil
}

func (odataUtil *OdataUtil) QueryWithAuthorizedUuids(ctx context.Context, idfIfc db.IdfClientInterface,
	authorized *authz.Authorized, argUuids [][]byte, idfParams *marinaIfc.IdfParameters, queryName string,
	entityType db.EntityType, tableCols ...interface{}) ([]*insights_interface.EntityWithMetric, int64, error) {

	var entities []*insights_interface.EntityWithMetric
	if !authorized.All && authorized.Uuids.Cardinality() == 0 {
		return entities, 0, nil
	}

	var fetchUuids []uuid4.Uuid
	if len(argUuids) == 0 && authorized.Uuids.Cardinality() > 0 {
		fetchUuids = authorized.Uuids.ToSlice()

	} else {
		requestedUuids := set.NewSet[uuid4.Uuid]()
		for _, uuidByte := range argUuids {
			requestedUuids.Add(*uuid4.ToUuid4(uuidByte))
		}

		if authorized.All {
			fetchUuids = requestedUuids.ToSlice()

		} else if authorized.Uuids.Cardinality() > 0 {
			fetchUuids = authorized.Uuids.Intersect(requestedUuids).ToSlice()
			if len(fetchUuids) == 0 {
				return entities, 0, nil
			}
		}
	}

	var predicates []*insights_interface.BooleanExpression
	queryBuilder := QUERY(queryName).SELECT(tableCols...).FROM(entityType.ToString())
	if idfParams != nil {
		limit, pageNumber, err := odataUtil.GetLimitAndOffset(ctx, idfParams)
		if err != nil {
			return nil, 0, err
		}

		if limit != nil && pageNumber != nil {
			queryBuilder = queryBuilder.LIMIT(*limit).SKIP(*pageNumber)
		}

		if idfParams.GetBooleanExpression() != nil {
			f := reflect.ValueOf(idfParams.BooleanExpression)
			v := f.Interface().(*insights_interface.BooleanExpression)
			predicates = append(predicates, v)
		}

		if len(idfParams.GetQueryOrderBy()) > 0 {
			if queryOrderBy := idfParams.QueryOrderBy[0]; queryOrderBy != nil {
				f := reflect.ValueOf(queryOrderBy)
				v := f.Interface().(*insights_interface.QueryOrderBy)
				queryBuilder = queryBuilder.SELECT(v.GetSortColumn()).ORDER_BY(v)
			}
		}
	}

	if len(fetchUuids) > 0 {
		var uuids [][]byte
		for _, uuid := range fetchUuids {
			uuid := uuid
			uuids = append(uuids, uuid.RawBytes())
		}
		predicates = append(predicates, IN(COL(uuidCol), BYTES_LIST(uuids...)))
	}

	if len(predicates) > 0 {
		queryBuilder = queryBuilder.WHERE(ALL(predicates...))
	}

	query, err := queryBuilder.Proto()
	if err != nil {
		errMsg := fmt.Sprintf("Error encountered while building the IDF query %s: %v", queryName, err)
		return nil, 0, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	arg := &insights_interface.GetEntitiesWithMetricsArg{Query: query}
	ret := &insights_interface.GetEntitiesWithMetricsRet{}
	err = idfIfc.GetEntitiesWithMetrics(ctx, arg, ret)
	if err != nil {
		return nil, 0, err
	}

	groupResults := ret.GetGroupResultsList()
	if len(groupResults) == 0 {
		return entities, 0, nil
	}

	entities = groupResults[0].GetRawResults()
	totalCount := groupResults[0].GetTotalEntityCount()
	if len(entities) == 0 {
		return entities, totalCount, nil
	}

	return entities, totalCount, nil
}
