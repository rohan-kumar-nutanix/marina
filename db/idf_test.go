/*
* Copyright (c) 2021 Nutanix Inc. All rights reserved.
*
* Author: rajesh.battala@nutanix.com
*
* IDF utilities for Marina Service.
 */

package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	mockIdfIfc "github.com/nutanix-core/acs-aos-go/insights/insights_interface/mock"
	. "github.com/nutanix-core/acs-aos-go/insights/insights_interface/query"
)

func TestIdfClientWithRetry(t *testing.T) {
	idfClient := IdfClientWithRetry()
	idfObj := idfClient.(*IdfClient)
	assert.NotNil(t, idfObj)
	assert.NotNil(t, idfObj.IdfSvc)
	assert.Nil(t, idfObj.Retry)
}

func TestIdfClientWithoutRetry(t *testing.T) {
	idfClient := IdfClientWithoutRetry()
	idfObj := idfClient.(*IdfClient)
	assert.NotNil(t, idfObj)
	assert.NotNil(t, idfObj.IdfSvc)
	assert.NotNil(t, idfObj.Retry)
}

func TestQuery(t *testing.T) {
	idfMockSvc := new(mockIdfIfc.InsightsServiceInterface)
	idfClient := IdfClient{
		IdfSvc: idfMockSvc,
	}

	tableName := "catalog_item_info"
	query, _ := QUERY("test_query").
		FROM(tableName).
		Proto()
	arg := &insights_interface.GetEntitiesWithMetricsArg{Query: query}
	ret := &insights_interface.GetEntitiesWithMetricsRet{}
	dummyEntityWithMetric := &insights_interface.EntityWithMetric{
		EntityGuid: &insights_interface.EntityGuid{
			EntityTypeName: &tableName,
		},
	}
	var rawResults []*insights_interface.EntityWithMetric
	rawResults = append(rawResults, dummyEntityWithMetric)
	validGroupResultList := &insights_interface.QueryGroupResult{
		RawResults: rawResults,
	}
	idfMockSvc.On("SendMsg", "GetEntitiesWithMetrics", arg, ret, idfClient.Retry).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(2).(*insights_interface.GetEntitiesWithMetricsRet)
			arg.GroupResultsList = make([]*insights_interface.QueryGroupResult, 0)
			arg.GroupResultsList = append(arg.GroupResultsList, validGroupResultList)
		}).Once()

	ctx := context.Background()
	entities, err := idfClient.Query(ctx, query)
	assert.Equal(t, entities, rawResults)
	assert.Nil(t, err)
	idfMockSvc.AssertExpectations(t)
}

func TestQueryEmptyGroupResult(t *testing.T) {
	idfMockSvc := new(mockIdfIfc.InsightsServiceInterface)
	idfClient := IdfClient{
		IdfSvc: idfMockSvc,
	}

	query, _ := QUERY("test_query").
		FROM("catalog_item_info").
		Proto()
	arg := &insights_interface.GetEntitiesWithMetricsArg{Query: query}
	ret := &insights_interface.GetEntitiesWithMetricsRet{}
	idfMockSvc.On("SendMsg", "GetEntitiesWithMetrics", arg, ret, idfClient.Retry).
		Return(nil).
		Once()

	ctx := context.Background()
	entities, err := idfClient.Query(ctx, query)

	assert.Nil(t, entities)
	assert.Equal(t, err, insights_interface.ErrNotFound)
	idfMockSvc.AssertExpectations(t)
}

func TestQueryEmptyResult(t *testing.T) {
	idfMockSvc := new(mockIdfIfc.InsightsServiceInterface)
	idfClient := IdfClient{
		IdfSvc: idfMockSvc,
	}

	query, _ := QUERY("test_query").
		FROM("catalog_item_info").
		Proto()
	arg := &insights_interface.GetEntitiesWithMetricsArg{Query: query}
	ret := &insights_interface.GetEntitiesWithMetricsRet{}
	emptyGroupResultList := &insights_interface.QueryGroupResult{}
	idfMockSvc.On("SendMsg", "GetEntitiesWithMetrics", arg, ret, idfClient.Retry).
		Return(nil).
		Run(func(args mock.Arguments) {
			arg := args.Get(2).(*insights_interface.GetEntitiesWithMetricsRet)
			arg.GroupResultsList = make([]*insights_interface.QueryGroupResult, 0)
			arg.GroupResultsList = append(arg.GroupResultsList, emptyGroupResultList)
		}).Once()

	ctx := context.Background()
	entities, err := idfClient.Query(ctx, query)

	assert.Nil(t, entities)
	assert.Equal(t, err, insights_interface.ErrNotFound)
	idfMockSvc.AssertExpectations(t)
}

func TestQueryError(t *testing.T) {
	idfMockSvc := new(mockIdfIfc.InsightsServiceInterface)
	idfClient := IdfClient{
		IdfSvc: idfMockSvc,
	}

	query, _ := QUERY("test_query").
		FROM("catalog_item_info").
		Proto()
	arg := &insights_interface.GetEntitiesWithMetricsArg{Query: query}
	ret := &insights_interface.GetEntitiesWithMetricsRet{}
	idfMockSvc.On("SendMsg", "GetEntitiesWithMetrics", arg, ret, idfClient.Retry).
		Return(insights_interface.ErrRetry).
		Once()

	ctx := context.Background()
	entities, err := idfClient.Query(ctx, query)

	assert.Nil(t, entities)
	assert.Equal(t, err, insights_interface.ErrRetry)
	idfMockSvc.AssertExpectations(t)
}
