/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Author: shreyash.turkar@nutanix.com
 *
 * IAMv2 interface implementation for Marina.
 */

package authz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	. "github.com/nutanix-core/acs-aos-go/insights/insights_interface/query"
	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/authz"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/authz/authz_cache"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/tracer"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-slbufs/util/sl_bufs/net"

	marinaError "github.com/nutanix-core/content-management-marina/errors"
)

const (
	all          = "*"
	categoryUuid = "category_uuid"
	ownerUuid    = "owner_uuid"
	uuid         = "uuid"

	eCapTable         = "abac_entity_capability"
	eCapCategoryIdCol = "category_id_list"
	eCapOwnerUuid     = "owner_reference"
	eCapKindID        = "kind_id"
	eCapKind          = "kind"
)

type AuthorizedUuids struct {
	AllAuthorized bool
	UuidSet       uuid4.UuidSet
}

type AuthzUtil struct {
}

// GetAuthorizedUuids Returns set of authorized entity uuids for the given user for the given operation type.
func (authzUtil *AuthzUtil) GetAuthorizedUuids(ctx context.Context, requestContext net.RpcRequestContext,
	iamClient authz_cache.IamClientIfc, cpdbClient cpdb.CPDBClientInterface, iamObjectType string, idfKind string, operation string) (
	*AuthorizedUuids, error) {

	span, ctx := tracer.StartSpan(ctx, "Get Authorized UUIDs")
	defer span.Finish()

	filter, err := getPermissionFilters(ctx, requestContext, iamClient, iamObjectType, operation)
	if err != nil {
		return nil, marinaError.ErrInternal.SetCauseAndLog(err)
	}
	return getEntityUuids(ctx, requestContext, cpdbClient, filter, idfKind)
}

// getPermissionFilters Calls on IAMv2 to get permission filters for the operation type on entity
// for the user.
func getPermissionFilters(ctx context.Context, requestContext net.RpcRequestContext, iamClient authz_cache.IamClientIfc,
	iamObjectType string, operation string) ([]*authz.AuthzAndFilter, error) {

	span, ctx := tracer.StartSpan(ctx, "IAM-request")
	defer span.Finish()

	userUuid := uuid4.ToUuid4(requestContext.GetUserUuid())
	userId := authz.UserIdentity{
		Type: "User",
		Uuid: userUuid.UuidToString(),
	}

	reqID := fmt.Sprintf("Marina_%s", operation)
	reqList := authz.AuthFilterRequest_FilterRequest{
		EntityType: iamObjectType,
		Operation:  operation,
		ReqId:      reqID,
	}
	authFilterReq := authz.AuthFilterRequest{
		User:    &userId,
		Env:     nil,
		ReqList: []*authz.AuthFilterRequest_FilterRequest{&reqList},
	}

	res := iamClient.GetPermissionFiltersWithContext(ctx, &authFilterReq)
	if res.Error != "" {
		errMsg := fmt.Sprintf("Iam request failed with error: %s, Http Response Code: %d",
			res.Error, res.HttpStatusCode)
		return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
	}

	resList := res.GetResList()
	if len(resList) != 1 {
		errMsg := fmt.Sprintf("Received incorrect IAMv2 response: %v.", resList)
		return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
	}

	filters := resList[0].GetFilters()
	return filters.GetAndFilters(), nil
}

// queryECapTable Query ECap table with predicates on the given entity. Returns entity uuids in string slice.
func queryECapTable(ctx context.Context, cpdbClient cpdb.CPDBClientInterface,
	predicates []*insights_interface.BooleanExpression, kind string, queryName string) (
	[]string, error) {

	span, ctx := tracer.StartSpan(ctx, "Query-ECap")
	defer span.Finish()

	whereClause := AND(EQ(COL(eCapKind), STR(kind)), ANY(predicates...))

	query, err := QUERY(queryName).
		SELECT(eCapKindID).
		FROM(eCapTable).
		WHERE(whereClause).Proto()
	if err != nil {
		errMsg := fmt.Sprintf("Error encountered while building the IDF query %s: %v", queryName, err)
		return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
	}

	arg := insights_interface.GetEntitiesWithMetricsArg{Query: query}
	resp, err := cpdbClient.Query(&arg)
	if err != nil {
		errMsg := fmt.Sprintf("CPDB query failed: %v", err)
		return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
	}

	var uuidList []string
	for _, entity := range resp {
		for _, metric := range entity.MetricDataList {
			if values := metric.GetValueList(); len(values) == 1 {
				uuidList = append(uuidList, values[0].Value.GetStrValue())
			}
		}
	}

	return uuidList, nil
}

// getEntityUuids Parses authz filters to retrieve authorized entity uuids for a user. Query ECap table if needed.
// Returns nil, error: if some issue occurs.
// Returns AuthorizedUuids{AllAuthorized: true}: if all entities are authorized.
// Returns AuthorizedUuids{AllAuthorized: false, List: []}: if subset of entities is/are authorized. List is a set of uuids.
func getEntityUuids(ctx context.Context, requestContext net.RpcRequestContext, cpdbClient cpdb.CPDBClientInterface,
	filters []*authz.AuthzAndFilter, idfKind string) (*AuthorizedUuids, error) {

	if len(filters) == 0 {
		return &AuthorizedUuids{}, nil
	}

	var uuidStrList []string
	var predicates []*insights_interface.BooleanExpression
	for _, filter := range filters {

		attrFilters := filter.GetAttrFilters()
		if len(attrFilters) != 1 {
			log.Warningf("Skipping incorrect attribute filter: %s", attrFilters)
			continue
		}

		attrFilter := attrFilters[0]
		switch attrFilter.Attr {
		case all:
			return &AuthorizedUuids{AllAuthorized: true}, nil
		case categoryUuid:
			var uuids []string
			err := json.Unmarshal(attrFilter.Value, &uuids)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to unmarshal json response for category UUID: %v", err)
				return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
			}
			if len(uuids) > 0 {
				predicates = append(predicates, IN(COL(eCapCategoryIdCol), STR_LIST(uuids...)))
			}
		case ownerUuid:
			userUuid := uuid4.ToUuid4(requestContext.GetUserUuid())
			predicates = append(predicates, EQ(COL(eCapOwnerUuid), STR(userUuid.UuidToString())))
		case uuid:
			var attrValue []string
			err := json.Unmarshal(attrFilter.Value, &attrValue)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to unmarshal json response for entity UUID: %v", err)
				return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
			}
			uuidStrList = append(uuidStrList, attrValue...)
		}
	}

	if predicates != nil {
		queryName := fmt.Sprintf("marina:authorized_uuids_%s", idfKind)
		uuids, err := queryECapTable(ctx, cpdbClient, predicates, idfKind, queryName)
		if err != nil {
			return nil, marinaError.ErrInternal.SetCauseAndLog(err)
		}
		uuidStrList = append(uuidStrList, uuids...)
	}

	var uuidList []uuid4.Uuid
	for _, uuidStr := range uuidStrList {
		uuid, err := uuid4.StringToUuid4(uuidStr)
		if err != nil {
			return nil, marinaError.ErrInternal.SetCauseAndLog(err)
		}
		uuidList = append(uuidList, *uuid)
	}
	return &AuthorizedUuids{UuidSet: uuid4.UuidSetFromList(uuidList)}, nil
}
