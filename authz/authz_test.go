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

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	"github.com/nutanix-core/acs-aos-go/insights/insights_interface/query"
	cpdbMock "github.com/nutanix-core/acs-aos-go/nusights/util/db/mocks"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/authz"
	iamMocks "github.com/nutanix-core/acs-aos-go/nutanix/util-go/authz/authz_cache/mocks"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-slbufs/util/sl_bufs/net"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
)

func TestGetPermissionFiltersResponseErr(t *testing.T) {

	mockIamClient := new(iamMocks.IamClientIfc)

	userUuid, _ := uuid4.New()

	ctx := context.Background()
	rpcReqCtx := net.RpcRequestContext{
		UserUuid: userUuid.RawBytes(),
	}
	userId := authz.UserIdentity{
		Type: "User",
		Uuid: userUuid.UuidToString()}

	iamObjectType := "mockEntity"
	operation := "mockOperation"

	reqID := fmt.Sprintf("Marina_%s", operation)
	reqList := authz.AuthFilterRequest_FilterRequest{
		EntityType: iamObjectType,
		Operation:  operation,
		ReqId:      reqID,
	}

	req := authz.AuthFilterRequest{
		User:    &userId,
		Env:     nil,
		ReqList: []*authz.AuthFilterRequest_FilterRequest{&reqList},
	}
	res := authz.AuthFilterResponse{Error: "Hey! An error occurred.", HttpStatusCode: 404}
	mockIamClient.On("GetPermissionFiltersWithContext", mock.Anything, &req).Return(&res).Once()

	filter, err := getPermissionFilters(ctx, rpcReqCtx, mockIamClient, iamObjectType, operation)
	assert.Nil(t, filter)
	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	mockIamClient.AssertExpectations(t)
}

func TestGetPermissionFiltersEmptyOrFilter(t *testing.T) {

	mockIamClient := new(iamMocks.IamClientIfc)

	ctx := context.Background()

	userUuid, _ := uuid4.New()

	rpcReqCtx := net.RpcRequestContext{
		UserUuid: userUuid.RawBytes(),
	}
	userId := authz.UserIdentity{
		Type: "User",
		Uuid: userUuid.UuidToString()}

	iamObjectType := "mockEntity"
	operation := "mockOperation"

	reqID := fmt.Sprintf("Marina_%s", operation)
	reqList := authz.AuthFilterRequest_FilterRequest{
		EntityType: iamObjectType,
		Operation:  operation,
		ReqId:      reqID,
	}

	req := authz.AuthFilterRequest{
		User:    &userId,
		Env:     nil,
		ReqList: []*authz.AuthFilterRequest_FilterRequest{&reqList},
	}
	res := authz.AuthFilterResponse{}
	mockIamClient.On("GetPermissionFiltersWithContext", mock.Anything, &req).Return(&res).Once()

	_, err := getPermissionFilters(ctx, rpcReqCtx, mockIamClient, iamObjectType, operation)
	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	mockIamClient.AssertExpectations(t)
}

func TestGetPermissionFilters(t *testing.T) {

	mockIamClient := new(iamMocks.IamClientIfc)

	ctx := context.Background()

	userUuid, _ := uuid4.New()

	rpcReqCtx := net.RpcRequestContext{
		UserUuid: userUuid.RawBytes(),
	}
	userId := authz.UserIdentity{
		Type: "User",
		Uuid: userUuid.UuidToString()}

	iamObjectType := "mockEntity"
	operation := "mockOperation"

	reqID := fmt.Sprintf("Marina_%s", operation)
	reqList := authz.AuthFilterRequest_FilterRequest{
		EntityType: iamObjectType,
		Operation:  operation,
		ReqId:      reqID,
	}

	req := authz.AuthFilterRequest{
		User:    &userId,
		Env:     nil,
		ReqList: []*authz.AuthFilterRequest_FilterRequest{&reqList},
	}

	andFilter := []*authz.AuthzAndFilter{{}}
	filter := authz.AuthzOrFilter{AndFilters: andFilter}
	resList := authz.AuthFilterResponse_FilterResponse{Filters: &filter}
	res := authz.AuthFilterResponse{ResList: []*authz.AuthFilterResponse_FilterResponse{&resList}}
	mockIamClient.On("GetPermissionFiltersWithContext", mock.Anything, &req).Return(&res).Once()

	ret, err := getPermissionFilters(ctx, rpcReqCtx, mockIamClient, iamObjectType, operation)
	assert.Nil(t, err)
	assert.Equal(t, ret, andFilter)
	mockIamClient.AssertExpectations(t)
}

func TestQueryECapTableQueryBuildErr(t *testing.T) {
	ctx := context.Background()
	cpdbClient := cpdbMock.CPDBClientInterface{}
	predicates := []*insights_interface.BooleanExpression{query.EQ(query.COL("mockCol"), query.STR("mockName"))}

	_, err := queryECapTable(ctx, &cpdbClient, predicates, "mock", "")
	assert.NotNil(t, err)
	cpdbClient.AssertExpectations(t)
}

func TestQueryECapTableIdfQueryErr(t *testing.T) {
	ctx := context.Background()
	cpdbClient := cpdbMock.CPDBClientInterface{}
	predicates := []*insights_interface.BooleanExpression{query.EQ(query.COL("mockCol"), query.STR("mockName"))}

	retErr := marinaError.ErrInvalidArgument
	cpdbClient.On("Query", mock.Anything).
		Return(nil, retErr)

	_, err := queryECapTable(ctx, &cpdbClient, predicates, "mock", "mockQuery")

	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	cpdbClient.AssertExpectations(t)
}

func TestQueryECapTable(t *testing.T) {
	ctx := context.Background()
	cpdbClient := cpdbMock.CPDBClientInterface{}
	predicates := []*insights_interface.BooleanExpression{query.EQ(query.COL("mockCol"), query.STR("mockName"))}

	value := insights_interface.TimeValuePair{
		Value: &insights_interface.DataValue{
			ValueType: &insights_interface.DataValue_StrValue{StrValue: "mock-uuid-mock-uuid-mock"}}}
	metricData1 := insights_interface.MetricData{
		ValueList: []*insights_interface.TimeValuePair{&value}}
	metricData2 := insights_interface.MetricData{}
	entityWithMetric := insights_interface.EntityWithMetric{
		MetricDataList: []*insights_interface.MetricData{&metricData1, &metricData2}}
	res := []*insights_interface.EntityWithMetric{&entityWithMetric}

	cpdbClient.On("Query", mock.Anything).Return(res, nil)

	uuids, err := queryECapTable(ctx, &cpdbClient, predicates, "mock", "mockQuery")

	assert.Nil(t, err)
	assert.Equal(t, uuids, []string{"mock-uuid-mock-uuid-mock"})
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuidsNoAccess(t *testing.T) {
	ctx := context.Background()
	rpcReqCtx := net.RpcRequestContext{}
	cpdbClient := cpdbMock.CPDBClientInterface{}
	var filter []*authz.AuthzAndFilter
	idfKind := "mockKind"

	authorizedUuids, err := getEntityUuids(ctx, rpcReqCtx, &cpdbClient, filter, idfKind)

	assert.Nil(t, err)
	assert.Equal(t, authorizedUuids.AllAuthorized, false)
	assert.Equal(t, len(uuid4.UuidListFromSet(authorizedUuids.UuidSet)), 0)
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuidsAllAuthorised(t *testing.T) {
	ctx := context.Background()
	rpcReqCtx := net.RpcRequestContext{}
	cpdbClient := cpdbMock.CPDBClientInterface{}
	idfKind := "mockKind"

	attrFilter := authz.AuthzAttrFilter{Attr: all}
	filter := authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{&attrFilter},
	}
	filters := []*authz.AuthzAndFilter{&filter}

	authorizedUuids, err := getEntityUuids(ctx, rpcReqCtx, &cpdbClient, filters, idfKind)

	assert.Nil(t, err)
	assert.Equal(t, authorizedUuids.AllAuthorized, true)
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuidsCategoryUnmarshalErr(t *testing.T) {
	ctx := context.Background()
	rpcReqCtx := net.RpcRequestContext{}
	cpdbClient := cpdbMock.CPDBClientInterface{}
	idfKind := "mockKind"

	attrFilter := authz.AuthzAttrFilter{Attr: categoryUuid}
	filter := authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{&attrFilter},
	}
	filters := []*authz.AuthzAndFilter{&filter}

	authorizedUuids, err := getEntityUuids(ctx, rpcReqCtx, &cpdbClient, filters, idfKind)

	assert.NotNil(t, err)
	assert.Nil(t, authorizedUuids)
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuidsEntityUuidUnmarshalErr(t *testing.T) {
	ctx := context.Background()
	rpcReqCtx := net.RpcRequestContext{}
	cpdbClient := cpdbMock.CPDBClientInterface{}
	idfKind := "mockKind"

	attrFilter := authz.AuthzAttrFilter{Attr: uuid}
	filter := authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{&attrFilter},
	}
	filters := []*authz.AuthzAndFilter{&filter}

	authorizedUuids, err := getEntityUuids(ctx, rpcReqCtx, &cpdbClient, filters, idfKind)

	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	assert.Nil(t, authorizedUuids)
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuidsNoPredicates(t *testing.T) {
	ctx := context.Background()
	rpcReqCtx := net.RpcRequestContext{}
	cpdbClient := cpdbMock.CPDBClientInterface{}
	idfKind := "mockKind"

	attrFilter := authz.AuthzAttrFilter{Attr: uuid}
	filter := authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{&attrFilter},
	}
	filters := []*authz.AuthzAndFilter{&filter}

	authorizedUuids, err := getEntityUuids(ctx, rpcReqCtx, &cpdbClient, filters, idfKind)

	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	assert.Nil(t, authorizedUuids)
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuidsECapQueryErr(t *testing.T) {
	ctx := context.Background()
	rpcReqCtx := net.RpcRequestContext{}
	idfKind := "mockKind"

	uuid, _ := uuid4.New()
	uuids := []string{uuid.String()}
	value, _ := json.Marshal(uuids)
	attrFilter := authz.AuthzAttrFilter{Attr: categoryUuid, Value: value}
	filter := authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{&attrFilter},
	}
	filters := []*authz.AuthzAndFilter{&filter}

	retErr := marinaError.ErrInternal.SetCauseAndLog(errors.New("Something went wrong."))
	cpdbClient := cpdbMock.CPDBClientInterface{}
	cpdbClient.On("Query", mock.Anything, mock.Anything).
		Return(nil, retErr).Once()

	authorizedUuids, err := getEntityUuids(ctx, rpcReqCtx, &cpdbClient, filters, idfKind)

	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	assert.Nil(t, authorizedUuids)
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuidsUuidConvertErr(t *testing.T) {
	ctx := context.Background()
	rpcReqCtx := net.RpcRequestContext{}
	cpdbClient := cpdbMock.CPDBClientInterface{}
	idfKind := "mockKind"

	uuids := []string{"invalid-uuid-mock-invalid"}
	value, _ := json.Marshal(uuids)
	attrFilter := authz.AuthzAttrFilter{Attr: uuid, Value: value}
	filter := authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{&attrFilter},
	}
	filters := []*authz.AuthzAndFilter{&filter}

	authorizedUuids, err := getEntityUuids(ctx, rpcReqCtx, &cpdbClient, filters, idfKind)

	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	assert.Nil(t, authorizedUuids)
	cpdbClient.AssertExpectations(t)
}

func TestGetEntityUuids(t *testing.T) {
	ctx := context.Background()
	userUuid, _ := uuid4.New()
	rpcReqCtx := net.RpcRequestContext{
		UserUuid: userUuid.RawBytes(),
	}
	idfKind := "mockKind"

	catUuid, _ := uuid4.New()
	catUuids := []string{catUuid.String()}
	value, err := json.Marshal(catUuids)
	assert.Nil(t, err)
	attrFilter := authz.AuthzAttrFilter{Attr: categoryUuid, Value: value}
	filter1 := authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{&attrFilter},
	}

	entityUuid1, _ := uuid4.New()
	entityUuid2, _ := uuid4.New()
	uuids := []string{entityUuid1.String(), entityUuid2.String()}
	value, err = json.Marshal(uuids)
	assert.Nil(t, err)
	attrFilter2 := authz.AuthzAttrFilter{Attr: uuid, Value: value}
	filter2 := authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{&attrFilter2},
	}

	attrFilter3 := authz.AuthzAttrFilter{Attr: ownerUuid}
	filter3 := authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{&attrFilter3},
	}

	filter4 := authz.AuthzAndFilter{
		AttrFilters: []*authz.AuthzAttrFilter{},
	}

	filters := []*authz.AuthzAndFilter{&filter1, &filter2, &filter3, &filter4}

	entityUuid3, _ := uuid4.New()
	entityUuidValue1 := insights_interface.TimeValuePair{
		Value: &insights_interface.DataValue{
			ValueType: &insights_interface.DataValue_StrValue{StrValue: entityUuid1.String()}}}
	metricData1 := insights_interface.MetricData{
		ValueList: []*insights_interface.TimeValuePair{&entityUuidValue1}}
	entityUuidValue2 := insights_interface.TimeValuePair{
		Value: &insights_interface.DataValue{
			ValueType: &insights_interface.DataValue_StrValue{StrValue: entityUuid3.String()}}}
	metricData2 := insights_interface.MetricData{
		ValueList: []*insights_interface.TimeValuePair{&entityUuidValue2}}
	entityWithMetric := insights_interface.EntityWithMetric{
		MetricDataList: []*insights_interface.MetricData{&metricData1, &metricData2}}
	res := []*insights_interface.EntityWithMetric{&entityWithMetric}

	cpdbClient := cpdbMock.CPDBClientInterface{}
	cpdbClient.On("Query", mock.Anything, mock.Anything).
		Return(res, nil).Once()

	authorizedUuids, err := getEntityUuids(ctx, rpcReqCtx, &cpdbClient, filters, idfKind)

	assert.Nil(t, err)
	assert.Equal(t, len(authorizedUuids.UuidSet), 3)
	cpdbClient.AssertExpectations(t)
}

func TestGetAuthorizedUuidsErr(t *testing.T) {

	userUuid, _ := uuid4.New()

	rpcReqCtx := net.RpcRequestContext{
		UserUuid: userUuid.RawBytes(),
	}

	iamObjectType := "mockEntity"
	operation := "mockOperation"
	idfKind := "mockIdfKind"
	mockIamClient := new(iamMocks.IamClientIfc)
	res := authz.AuthFilterResponse{Error: "Error", HttpStatusCode: 404}
	mockIamClient.On("GetPermissionFiltersWithContext", mock.Anything, mock.Anything).Return(&res).Once()

	authzUtil := AuthzUtil{}

	ctx := context.Background()
	cpdbClient := cpdbMock.CPDBClientInterface{}
	_, err := authzUtil.GetAuthorizedUuids(ctx, rpcReqCtx, mockIamClient, &cpdbClient, iamObjectType, idfKind,
		operation)

	assert.Error(t, err)
	assert.IsType(t, &marinaError.MarinaError{}, err)
	mockIamClient.AssertExpectations(t)
	cpdbClient.AssertExpectations(t)
}

func TestGetAuthorizedUuids(t *testing.T) {

	userUuid, _ := uuid4.New()

	rpcReqCtx := net.RpcRequestContext{
		UserUuid: userUuid.RawBytes(),
	}

	iamObjectType := "mockEntity"
	operation := "mockOperation"
	idfKind := "mockIdfKind"
	mockIamClient := new(iamMocks.IamClientIfc)

	andFilter := []*authz.AuthzAndFilter{{}}
	filter := authz.AuthzOrFilter{AndFilters: andFilter}
	resList := authz.AuthFilterResponse_FilterResponse{Filters: &filter}
	res := authz.AuthFilterResponse{ResList: []*authz.AuthFilterResponse_FilterResponse{&resList}}
	mockIamClient.On("GetPermissionFiltersWithContext", mock.Anything, mock.Anything).Return(&res).Once()

	authzUtil := AuthzUtil{}

	ctx := context.Background()
	cpdbClient := cpdbMock.CPDBClientInterface{}
	uuidSet, err := authzUtil.GetAuthorizedUuids(ctx, rpcReqCtx, mockIamClient, &cpdbClient, iamObjectType, idfKind,
		operation)

	assert.Nil(t, err)
	assert.NotNil(t, uuidSet)
	mockIamClient.AssertExpectations(t)
	cpdbClient.AssertExpectations(t)
}
