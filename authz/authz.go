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

	set "github.com/deckarep/golang-set/v2"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	. "github.com/nutanix-core/acs-aos-go/insights/insights_interface/query"
	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/authz"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/authz/authz_cache"
	grpcUtils "github.com/nutanix-core/acs-aos-go/nutanix/util-go/net/grpc"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/content-management-marina/db"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
)

const (
	userType = "User"

	allFilter          = "*"
	categoryUuidFilter = "category_uuid"
	ownerUuidFilter    = "owner_uuid"
	uuidFilter         = "uuid"
	clusterUuidFilter  = "cluster_uuid"

	eCapCategoryIdCol = "category_id_list"
	eCapOwnerUuidCol  = "owner_reference"
	eCapKindIdCol     = "kind_id"
	eCapKindCol       = "kind"
)

type Authorized struct {
	All   bool
	Uuids set.Set[uuid4.Uuid]
}

type AuthzUtil struct {
}

// GetAuthorizedUuids Returns set of authorized entity uuids for the given user for the given operation type.
func (authzUtil *AuthzUtil) GetAuthorizedUuids(ctx context.Context, iamClient authz_cache.IamClientIfc,
	cpdbClient cpdb.CPDBClientInterface, iamObjectType string, idfKind string, operation string) (*Authorized, error) {

	filters, err := getPermissionFilters(ctx, iamClient, iamObjectType, operation)
	if err != nil {
		return nil, err
	}
	return getEntityUuids(ctx, cpdbClient, filters, idfKind)
}

// getPermissionFilters Calls on IAMv2 to get permission filters for the operation type on entity
// for the user.
func getPermissionFilters(ctx context.Context, iamClient authz_cache.IamClientIfc, iamObjectType string,
	operation string) ([]*authz.AuthzAndFilter, error) {

	requestContext, err := grpcUtils.GetRequestContextFromContext(ctx)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get request context: %v", err)
		return nil, marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
	}

	userUuid := uuid4.ToUuid4(requestContext.GetUserUuid())
	userId := authz.UserIdentity{
		Type: userType,
		Uuid: userUuid.UuidToString(),
	}

	reqID := fmt.Sprintf("Marina_%s", operation)
	requests := &authz.AuthFilterRequest_FilterRequest{
		EntityType: iamObjectType,
		Operation:  operation,
		ReqId:      reqID,
	}
	authFilterReq := &authz.AuthFilterRequest{
		User:    &userId,
		Env:     nil,
		ReqList: []*authz.AuthFilterRequest_FilterRequest{requests},
	}

	res := iamClient.GetPermissionFiltersWithContext(ctx, authFilterReq)
	if res.Error != "" {
		errMsg := fmt.Sprintf("IAM request failed with error: %s, HTTP Response Code: %d",
			res.Error, res.HttpStatusCode)
		return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
	}

	responses := res.GetResList()
	if len(responses) != 1 {
		errMsg := fmt.Sprintf("Received incorrect IAMv2 response: %v.", responses)
		return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
	}

	filters := responses[0].GetFilters().GetAndFilters()
	return filters, nil
}

// queryECapTable Query ECap table with predicates on the given entity. Returns entity uuids in string slice.
func queryECapTable(cpdbClient cpdb.CPDBClientInterface, predicates []*insights_interface.BooleanExpression,
	kind string, queryName string) ([]string, error) {

	whereClause := AND(EQ(COL(eCapKindCol), STR(kind)), ANY(predicates...))
	query, err := QUERY(queryName).SELECT(eCapKindIdCol).FROM(db.EntityCapability.ToString()).WHERE(whereClause).Proto()
	if err != nil {
		errMsg := fmt.Sprintf("Error encountered while building the IDF query %s: %v", queryName, err)
		return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
	}

	arg := &insights_interface.GetEntitiesWithMetricsArg{Query: query}
	ret, err := cpdbClient.Query(arg)
	if err != nil && err != insights_interface.ErrNotFound {
		errMsg := fmt.Sprintf("CPDB query failed: %v", err)
		return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
	}

	var uuids []string
	for _, entity := range ret {
		for _, metric := range entity.MetricDataList {
			if values := metric.GetValueList(); len(values) == 1 {
				uuids = append(uuids, values[0].Value.GetStrValue())
			}
		}
	}
	return uuids, nil
}

// getEntityUuids Parses authz filters to retrieve authorized entity uuids for a user. Query ECap table if needed.
// Returns nil, error: if some issue occurs.
// Returns AuthorizedUuids{AllAuthorized: true}: if all entities are authorized.
// Returns AuthorizedUuids{AllAuthorized: false, List: []}: if subset of entities is/are authorized. List is a set of uuids.
func getEntityUuids(ctx context.Context, cpdbClient cpdb.CPDBClientInterface, filters []*authz.AuthzAndFilter,
	idfKind string) (*Authorized, error) {

	if len(filters) == 0 {
		return &Authorized{All: false, Uuids: set.NewSet[uuid4.Uuid]()}, nil
	}

	var uuidStrs []string
	var predicates []*insights_interface.BooleanExpression
	for _, filter := range filters {
		attrFilters := filter.GetAttrFilters()
		if len(attrFilters) != 1 {
			log.Warningf("Skipping incorrect attribute filter: %s", attrFilters)
			continue
		}

		attrFilter := attrFilters[0]
		switch attrFilter.Attr {
		case allFilter:
			return &Authorized{All: true, Uuids: set.NewSet[uuid4.Uuid]()}, nil

		case categoryUuidFilter:
			var uuids []string
			err := json.Unmarshal(attrFilter.Value, &uuids)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to unmarshal json response for category UUID: %v", err)
				return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
			}
			if len(uuids) > 0 {
				predicates = append(predicates, IN(COL(eCapCategoryIdCol), STR_LIST(uuids...)))
			}

		case ownerUuidFilter:
			requestContext, err := grpcUtils.GetRequestContextFromContext(ctx)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to get request context: %v", err)
				return nil, marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
			}
			userUuid := uuid4.ToUuid4(requestContext.GetUserUuid())
			predicates = append(predicates, EQ(COL(eCapOwnerUuidCol), STR(userUuid.UuidToString())))

		case uuidFilter:
			var attrValue []string
			err := json.Unmarshal(attrFilter.Value, &attrValue)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to unmarshal json response for entity UUID: %v", err)
				return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
			}
			uuidStrs = append(uuidStrs, attrValue...)

		case clusterUuidFilter:
			// This is field is used by Prism Admin role. We should not act on this filter.
			continue
		}
	}

	if len(predicates) > 0 {
		queryName := fmt.Sprintf("authorized_uuids_%s", idfKind)
		uuids, err := queryECapTable(cpdbClient, predicates, idfKind, queryName)
		if err != nil {
			return nil, err
		}
		uuidStrs = append(uuidStrs, uuids...)
	}

	uuids := set.NewSet[uuid4.Uuid]()
	for _, uuidStr := range uuidStrs {
		uuid, err := uuid4.StringToUuid4(uuidStr)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to convert string to UUID")
			return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
		}
		uuids.Add(*uuid)
	}
	return &Authorized{All: false, Uuids: uuids}, nil
}
