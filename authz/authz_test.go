/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Author: shreyash.turkar@nutanix.com
 *
 */

package authz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	"github.com/nutanix-core/acs-aos-go/insights/insights_interface/query"
	cpdbMock "github.com/nutanix-core/acs-aos-go/nusights/util/db/mocks"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/authz"
	iamMocks "github.com/nutanix-core/acs-aos-go/nutanix/util-go/authz/authz_cache/mocks"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
)

func TestGetPermissionFiltersRequestContextError(t *testing.T) {
	ctx := context.Background()
	mockIamClient := new(iamMocks.IamClientIfc)
	iamObjectType := "mockEntity"
	operation := "mockOperation"

	filter, err := getPermissionFilters(ctx, mockIamClient, iamObjectType, operation)

	assert.Nil(t, filter)
	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	mockIamClient.AssertExpectations(t)
}

func TestGetPermissionFiltersIamError(t *testing.T) {
	userUuid, _ := uuid4.New()
	userId := &authz.UserIdentity{
		Type: userType,
		Uuid: userUuid.UuidToString(),
	}
	iamObjectType := "mockEntity"
	operation := "mockOperation"
	reqID := fmt.Sprintf("Marina_%s", operation)
	reqList := &authz.AuthFilterRequest_FilterRequest{
		EntityType: iamObjectType,
		Operation:  operation,
		ReqId:      reqID,
	}
	req := &authz.AuthFilterRequest{
		User:    userId,
		Env:     nil,
		ReqList: []*authz.AuthFilterRequest_FilterRequest{reqList},
	}
	res := &authz.AuthFilterResponse{Error: "Hey! An error occurred.", HttpStatusCode: 404}
	mockIamClient := new(iamMocks.IamClientIfc)
	mockIamClient.On("GetPermissionFiltersWithContext", mock.Anything, req).Return(res).Once()

	ctx := context.Background()
	token := fmt.Sprintf("{\"userUUID\":\"%s\",\"username\":\"admin\"}", userUuid.String())
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("token_claims", token))

	filter, err := getPermissionFilters(ctx, mockIamClient, iamObjectType, operation)

	assert.Nil(t, filter)
	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	mockIamClient.AssertExpectations(t)
}

func TestGetPermissionFiltersEmptyFilter(t *testing.T) {
	userUuid, _ := uuid4.New()
	userId := &authz.UserIdentity{
		Type: userType,
		Uuid: userUuid.UuidToString(),
	}
	iamObjectType := "mockEntity"
	operation := "mockOperation"
	reqID := fmt.Sprintf("Marina_%s", operation)
	reqList := &authz.AuthFilterRequest_FilterRequest{
		EntityType: iamObjectType,
		Operation:  operation,
		ReqId:      reqID,
	}
	req := &authz.AuthFilterRequest{
		User:    userId,
		Env:     nil,
		ReqList: []*authz.AuthFilterRequest_FilterRequest{reqList},
	}
	res := &authz.AuthFilterResponse{}
	mockIamClient := new(iamMocks.IamClientIfc)
	mockIamClient.On("GetPermissionFiltersWithContext", mock.Anything, req).Return(res).Once()

	ctx := context.Background()
	token := fmt.Sprintf("{\"userUUID\":\"%s\",\"username\":\"admin\"}", userUuid.String())
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("token_claims", token))

	_, err := getPermissionFilters(ctx, mockIamClient, iamObjectType, operation)

	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	mockIamClient.AssertExpectations(t)
}

func TestGetPermissionFilters(t *testing.T) {
	userUuid, _ := uuid4.New()
	userId := &authz.UserIdentity{
		Type: userType,
		Uuid: userUuid.UuidToString(),
	}
	iamObjectType := "mockEntity"
	operation := "mockOperation"
	reqID := fmt.Sprintf("Marina_%s", operation)
	reqList := &authz.AuthFilterRequest_FilterRequest{
		EntityType: iamObjectType,
		Operation:  operation,
		ReqId:      reqID,
	}
	req := &authz.AuthFilterRequest{
		User:    userId,
		Env:     nil,
		ReqList: []*authz.AuthFilterRequest_FilterRequest{reqList},
	}
	andFilter := []*authz.AuthzAndFilter{{}}
	filter := &authz.AuthzOrFilter{AndFilters: andFilter}
	resList := &authz.AuthFilterResponse_FilterResponse{Filters: filter}
	res := &authz.AuthFilterResponse{ResList: []*authz.AuthFilterResponse_FilterResponse{resList}}
	mockIamClient := new(iamMocks.IamClientIfc)
	mockIamClient.On("GetPermissionFiltersWithContext", mock.Anything, req).Return(res).Once()

	ctx := context.Background()
	token := fmt.Sprintf("{\"userUUID\":\"%s\",\"username\":\"admin\"}", userUuid.String())
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("token_claims", token))

	ret, err := getPermissionFilters(ctx, mockIamClient, iamObjectType, operation)

	assert.Nil(t, err)
	assert.Equal(t, ret, andFilter)
	mockIamClient.AssertExpectations(t)
}

func TestQueryECapTableQueryBuildErr(t *testing.T) {
	cpdbClient := &cpdbMock.CPDBClientInterface{}
	predicates := []*insights_interface.BooleanExpression{query.EQ(query.COL("mockCol"), query.STR("mockName"))}

	_, err := queryECapTable(cpdbClient, predicates, "mock", "")

	assert.NotNil(t, err)
	cpdbClient.AssertExpectations(t)
}

func TestQueryECapTableIdfQueryErr(t *testing.T) {
	predicates := []*insights_interface.BooleanExpression{query.EQ(query.COL("mockCol"), query.STR("mockName"))}

	retErr := marinaError.ErrInvalidArgument
	cpdbClient := &cpdbMock.CPDBClientInterface{}
	cpdbClient.On("Query", mock.Anything).Return(nil, retErr)

	_, err := queryECapTable(cpdbClient, predicates, "mock", "mockQuery")

	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	cpdbClient.AssertExpectations(t)
}

func TestQueryECapTable(t *testing.T) {
	predicates := []*insights_interface.BooleanExpression{query.EQ(query.COL("mockCol"), query.STR("mockName"))}
	value := &insights_interface.TimeValuePair{
		Value: &insights_interface.DataValue{
			ValueType: &insights_interface.DataValue_StrValue{StrValue: "mock-uuid-mock-uuid-mock"},
		},
	}
	metricData1 := &insights_interface.MetricData{ValueList: []*insights_interface.TimeValuePair{value}}
	metricData2 := &insights_interface.MetricData{}
	entityWithMetric := &insights_interface.EntityWithMetric{
		MetricDataList: []*insights_interface.MetricData{metricData1, metricData2},
	}
	res := []*insights_interface.EntityWithMetric{entityWithMetric}
	cpdbClient := &cpdbMock.CPDBClientInterface{}
	cpdbClient.On("Query", mock.Anything).Return(res, nil)

	uuids, err := queryECapTable(cpdbClient, predicates, "mock", "mockQuery")

	assert.Nil(t, err)
	assert.Equal(t, uuids, []string{"mock-uuid-mock-uuid-mock"})
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuidsNoAccess(t *testing.T) {
	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("token_claims", "{}"))
	cpdbClient := &cpdbMock.CPDBClientInterface{}
	var filter []*authz.AuthzAndFilter
	idfKind := "mockKind"

	authorized, err := getEntityUuids(ctx, cpdbClient, filter, idfKind)

	assert.Nil(t, err)
	assert.Equal(t, authorized.All, false)
	assert.Equal(t, authorized.Uuids.Cardinality(), 0)
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuidsAllAuthorised(t *testing.T) {
	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("token_claims", "{}"))
	cpdbClient := &cpdbMock.CPDBClientInterface{}
	idfKind := "mockKind"
	attrFilter := &authz.AuthzAttrFilter{Attr: allFilter}
	filter := &authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{attrFilter},
	}
	filters := []*authz.AuthzAndFilter{filter}

	authorized, err := getEntityUuids(ctx, cpdbClient, filters, idfKind)

	assert.Nil(t, err)
	assert.Equal(t, authorized.All, true)
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuidsClusterUuidFilter(t *testing.T) {
	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("token_claims", "{}"))
	cpdbClient := &cpdbMock.CPDBClientInterface{}
	attrFilter := &authz.AuthzAttrFilter{Attr: clusterUuidFilter}
	filter := &authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{attrFilter},
	}
	filters := []*authz.AuthzAndFilter{filter}
	idfKind := "mockKind"

	authorized, err := getEntityUuids(ctx, cpdbClient, filters, idfKind)

	assert.Nil(t, err)
	assert.Equal(t, authorized.All, false)
	assert.Equal(t, authorized.Uuids.Cardinality(), 0)
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuidsCategoryUnmarshalErr(t *testing.T) {
	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("token_claims", "{}"))
	cpdbClient := &cpdbMock.CPDBClientInterface{}
	idfKind := "mockKind"
	attrFilter := &authz.AuthzAttrFilter{Attr: categoryUuidFilter}
	filter := &authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{attrFilter},
	}
	filters := []*authz.AuthzAndFilter{filter}

	authorized, err := getEntityUuids(ctx, cpdbClient, filters, idfKind)

	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	assert.Nil(t, authorized)
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuidsEntityUuidUnmarshalErr(t *testing.T) {
	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("token_claims", "{}"))
	cpdbClient := &cpdbMock.CPDBClientInterface{}
	idfKind := "mockKind"
	attrFilter := &authz.AuthzAttrFilter{Attr: uuidFilter}
	filter := &authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{attrFilter},
	}
	filters := []*authz.AuthzAndFilter{filter}

	authorized, err := getEntityUuids(ctx, cpdbClient, filters, idfKind)

	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	assert.Nil(t, authorized)
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuidsNoPredicates(t *testing.T) {
	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("token_claims", "{}"))
	cpdbClient := &cpdbMock.CPDBClientInterface{}
	idfKind := "mockKind"
	attrFilter := &authz.AuthzAttrFilter{Attr: uuidFilter}
	filter := &authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{attrFilter},
	}
	filters := []*authz.AuthzAndFilter{filter}

	authorized, err := getEntityUuids(ctx, cpdbClient, filters, idfKind)

	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	assert.Nil(t, authorized)
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuidsECapQueryErr(t *testing.T) {
	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("token_claims", "{}"))
	idfKind := "mockKind"
	uuid, _ := uuid4.New()
	uuids := []string{uuid.String()}
	value, _ := json.Marshal(uuids)
	attrFilter := &authz.AuthzAttrFilter{Attr: categoryUuidFilter, Value: value}
	filter := &authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{attrFilter},
	}
	filters := []*authz.AuthzAndFilter{filter}
	retErr := marinaError.ErrInternal.SetCauseAndLog(errors.New("Something went wrong."))
	cpdbClient := &cpdbMock.CPDBClientInterface{}
	cpdbClient.On("Query", mock.Anything, mock.Anything).Return(nil, retErr).Once()

	authorized, err := getEntityUuids(ctx, cpdbClient, filters, idfKind)

	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	assert.Nil(t, authorized)
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuidsUuidConvertErr(t *testing.T) {
	ctx := context.Background()
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("token_claims", "{}"))
	cpdbClient := &cpdbMock.CPDBClientInterface{}
	idfKind := "mockKind"
	uuids := []string{"invalid-uuid-mock-invalid"}
	value, _ := json.Marshal(uuids)
	attrFilter := &authz.AuthzAttrFilter{Attr: uuidFilter, Value: value}
	filter := &authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{attrFilter},
	}
	filters := []*authz.AuthzAndFilter{filter}

	authorized, err := getEntityUuids(ctx, cpdbClient, filters, idfKind)

	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	assert.Nil(t, authorized)
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuidsRequestContextError(t *testing.T) {
	ctx := context.Background()
	cpdbClient := &cpdbMock.CPDBClientInterface{}
	attrFilter := &authz.AuthzAttrFilter{Attr: ownerUuidFilter}
	filter := &authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{attrFilter},
	}
	filters := []*authz.AuthzAndFilter{filter}
	idfKind := "mockKind"

	authorized, err := getEntityUuids(ctx, cpdbClient, filters, idfKind)

	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	assert.Nil(t, authorized)
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuids(t *testing.T) {
	ctx := context.Background()
	userUuid, _ := uuid4.New()
	token := fmt.Sprintf("{\"userUUID\":\"%s\",\"username\":\"admin\"}", userUuid.String())
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("token_claims", token))
	idfKind := "mockKind"
	catUuid, _ := uuid4.New()
	catUuids := []string{catUuid.String()}
	value, _ := json.Marshal(catUuids)
	attrFilter := &authz.AuthzAttrFilter{Attr: categoryUuidFilter, Value: value}
	filter1 := &authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{attrFilter},
	}

	entityUuid1, _ := uuid4.New()
	entityUuid2, _ := uuid4.New()
	uuids := []string{entityUuid1.String(), entityUuid2.String()}
	value, _ = json.Marshal(uuids)
	attrFilter2 := &authz.AuthzAttrFilter{Attr: uuidFilter, Value: value}
	filter2 := &authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{attrFilter2},
	}

	attrFilter3 := &authz.AuthzAttrFilter{Attr: ownerUuidFilter}
	filter3 := &authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{attrFilter3},
	}

	filter4 := &authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{},
	}

	filters := []*authz.AuthzAndFilter{filter1, filter2, filter3, filter4}

	entityUuid3, _ := uuid4.New()
	entityUuidValue1 := &insights_interface.TimeValuePair{
		Value: &insights_interface.DataValue{
			ValueType: &insights_interface.DataValue_StrValue{StrValue: entityUuid1.String()},
		},
	}
	metricData1 := &insights_interface.MetricData{ValueList: []*insights_interface.TimeValuePair{entityUuidValue1}}
	entityUuidValue2 := &insights_interface.TimeValuePair{
		Value: &insights_interface.DataValue{
			ValueType: &insights_interface.DataValue_StrValue{StrValue: entityUuid3.String()},
		},
	}
	metricData2 := &insights_interface.MetricData{ValueList: []*insights_interface.TimeValuePair{entityUuidValue2}}
	entityWithMetric := &insights_interface.EntityWithMetric{
		MetricDataList: []*insights_interface.MetricData{metricData1, metricData2},
	}
	res := []*insights_interface.EntityWithMetric{entityWithMetric}

	cpdbClient := &cpdbMock.CPDBClientInterface{}
	cpdbClient.On("Query", mock.Anything, mock.Anything).Return(res, nil).Once()

	authorized, err := getEntityUuids(ctx, cpdbClient, filters, idfKind)

	assert.Nil(t, err)
	assert.Equal(t, authorized.Uuids.Cardinality(), 3)
	cpdbClient.AssertExpectations(t)
}

func TestGetAuthorizedUuidsErr(t *testing.T) {
	userUuid, _ := uuid4.New()
	ctx := context.Background()
	token := fmt.Sprintf("{\"userUUID\":\"%s\",\"username\":\"admin\"}", userUuid.String())
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("token_claims", token))
	iamObjectType := "mockEntity"
	operation := "mockOperation"
	idfKind := "mockIdfKind"
	res := &authz.AuthFilterResponse{Error: "Error", HttpStatusCode: 404}
	mockIamClient := new(iamMocks.IamClientIfc)
	mockIamClient.On("GetPermissionFiltersWithContext", mock.Anything, mock.Anything).Return(res).Once()
	cpdbClient := &cpdbMock.CPDBClientInterface{}

	authzUtil := AuthzUtil{}
	authorized, err := authzUtil.GetAuthorizedUuids(ctx, mockIamClient, cpdbClient, iamObjectType, idfKind, operation)

	assert.Nil(t, authorized)
	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	mockIamClient.AssertExpectations(t)
	cpdbClient.AssertExpectations(t)
}

func TestGetAuthorizedUuids(t *testing.T) {
	userUuid, _ := uuid4.New()
	ctx := context.Background()
	token := fmt.Sprintf("{\"userUUID\":\"%s\",\"username\":\"admin\"}", userUuid.String())
	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("token_claims", token))
	iamObjectType := "mockEntity"
	operation := "mockOperation"
	idfKind := "mockIdfKind"
	andFilter := []*authz.AuthzAndFilter{{}}
	filter := &authz.AuthzOrFilter{AndFilters: andFilter}
	resList := &authz.AuthFilterResponse_FilterResponse{Filters: filter}
	res := &authz.AuthFilterResponse{ResList: []*authz.AuthFilterResponse_FilterResponse{resList}}
	mockIamClient := new(iamMocks.IamClientIfc)
	mockIamClient.On("GetPermissionFiltersWithContext", mock.Anything, mock.Anything).Return(res).Once()
	cpdbClient := &cpdbMock.CPDBClientInterface{}

	authzUtil := AuthzUtil{}
	authorized, err := authzUtil.GetAuthorizedUuids(ctx, mockIamClient, cpdbClient, iamObjectType, idfKind, operation)

	assert.Nil(t, err)
	assert.NotNil(t, authorized)
	mockIamClient.AssertExpectations(t)
	cpdbClient.AssertExpectations(t)
}
