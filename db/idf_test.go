/*
* Copyright (c) 2021 Nutanix Inc. All rights reserved.
*
* Author: rajesh.battala@nutanix.com
*
 */

package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	mockIdf "github.com/nutanix-core/acs-aos-go/insights/insights_interface/mock"
	mockCpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db/mocks"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
)

func TestIdfClientWithRetry(t *testing.T) {
	mockIdfIfc := &mockIdf.InsightsServiceInterface{}
	idfClient := IdfClientWithRetry(mockIdfIfc)
	idfObj := idfClient.(*IdfClient)
	assert.NotNil(t, idfObj)
	assert.NotNil(t, idfObj.IdfSvc)
	assert.Nil(t, idfObj.Retry)
}

func TestIdfClientWithoutRetry(t *testing.T) {
	mockIdfIfc := &mockIdf.InsightsServiceInterface{}
	idfClient := IdfClientWithoutRetry(mockIdfIfc)
	idfObj := idfClient.(*IdfClient)
	assert.NotNil(t, idfObj)
	assert.NotNil(t, idfObj.IdfSvc)
	assert.NotNil(t, idfObj.Retry)
}

func TestDeleteEntities(t *testing.T) {
	mockIdfIfc := &mockIdf.InsightsServiceInterface{}
	idfClient := IdfClient{IdfSvc: mockIdfIfc}
	cpdbMock := &mockCpdb.CPDBClientInterface{}
	entityUuid, _ := uuid4.New()
	entityUuidStr := entityUuid.String()
	cas := uint64(1)
	entities := []*insights_interface.Entity{
		{
			EntityGuid: &insights_interface.EntityGuid{
				EntityId: &entityUuidStr,
			},
			CasValue: &cas,
		},
	}
	cpdbMock.On("GetEntities", mock.Anything, mock.Anything).Return(entities, nil).Once()
	mockIdfIfc.On("SendMsg", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()

	err := idfClient.DeleteEntities(context.Background(), cpdbMock, CatalogItem, []string{entityUuidStr}, true)

	assert.NoError(t, err)
	cpdbMock.AssertExpectations(t)
	mockIdfIfc.AssertExpectations(t)
}

func TestDeleteEntitiesPartial(t *testing.T) {
	mockIdfIfc := &mockIdf.InsightsServiceInterface{}
	idfClient := IdfClient{IdfSvc: mockIdfIfc}
	cpdbMock := &mockCpdb.CPDBClientInterface{}
	entityUuid, _ := uuid4.New()
	entityUuidStr := entityUuid.String()
	cas := uint64(1)
	entities := []*insights_interface.Entity{
		{
			EntityGuid: &insights_interface.EntityGuid{
				EntityId: &entityUuidStr,
			},
			CasValue: &cas,
		},
	}
	cpdbMock.On("GetEntities", mock.Anything, mock.Anything).Return(entities, nil).Once()
	noError := insights_interface.InsightsErrorProto_kNoError
	mockIdfIfc.On("SendMsg", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			ret := args.Get(2).(*insights_interface.BatchDeleteEntitiesRet)
			ret.RetList = append(ret.RetList, &insights_interface.BatchDeleteEntitiesRet_RetElem{
				Status: &insights_interface.InsightsErrorProto{ErrorType: &noError},
			})
			internalError := insights_interface.InsightsErrorProto_kInternalError
			ret.RetList = append(ret.RetList, &insights_interface.BatchDeleteEntitiesRet_RetElem{
				Entity: entities[0],
				Status: &insights_interface.InsightsErrorProto{ErrorType: &internalError},
			})
		}).Return(insights_interface.ErrPartial).Once()
	mockIdfIfc.On("SendMsg", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()

	err := idfClient.DeleteEntities(context.Background(), cpdbMock, CatalogItem, []string{entityUuidStr}, true)

	assert.NoError(t, err)
	cpdbMock.AssertExpectations(t)
	mockIdfIfc.AssertExpectations(t)
}

func TestDeleteEntitiesNoOp(t *testing.T) {
	mockIdfIfc := &mockIdf.InsightsServiceInterface{}
	idfClient := IdfClient{IdfSvc: mockIdfIfc}
	cpdbMock := &mockCpdb.CPDBClientInterface{}
	entityUuid, _ := uuid4.New()
	var entities []*insights_interface.Entity
	cpdbMock.On("GetEntities", mock.Anything, mock.Anything).Return(entities, nil).Once()

	err := idfClient.DeleteEntities(context.Background(), cpdbMock, CatalogItem, []string{entityUuid.String()}, true)

	assert.NoError(t, err)
	cpdbMock.AssertExpectations(t)
}

func TestDeleteEntitiesGetError(t *testing.T) {
	mockIdfIfc := &mockIdf.InsightsServiceInterface{}
	idfClient := IdfClient{IdfSvc: mockIdfIfc}
	cpdbMock := &mockCpdb.CPDBClientInterface{}
	entityUuid, _ := uuid4.New()
	cpdbMock.On("GetEntities", mock.Anything, mock.Anything).Return(nil, marinaError.ErrInternalError()).Once()

	err := idfClient.DeleteEntities(context.Background(), cpdbMock, CatalogItem, []string{entityUuid.String()}, true)

	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
	cpdbMock.AssertExpectations(t)
}

func TestDeleteEntitiesDeleteError(t *testing.T) {
	mockIdfIfc := &mockIdf.InsightsServiceInterface{}
	idfClient := IdfClient{IdfSvc: mockIdfIfc}
	cpdbMock := &mockCpdb.CPDBClientInterface{}
	entityUuid, _ := uuid4.New()
	entityUuidStr := entityUuid.String()
	cas := uint64(1)
	entities := []*insights_interface.Entity{
		{
			EntityGuid: &insights_interface.EntityGuid{
				EntityId: &entityUuidStr,
			},
			CasValue: &cas,
		},
	}
	cpdbMock.On("GetEntities", mock.Anything, mock.Anything).Return(entities, nil).Once()
	mockIdfIfc.On("SendMsg", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(insights_interface.ErrInternalError).Once()

	err := idfClient.DeleteEntities(context.Background(), cpdbMock, CatalogItem, []string{entityUuidStr}, true)

	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
	cpdbMock.AssertExpectations(t)
	mockIdfIfc.AssertExpectations(t)
}

func TestDeleteEntitiesPartialError(t *testing.T) {
	mockIdfIfc := &mockIdf.InsightsServiceInterface{}
	idfClient := IdfClient{IdfSvc: mockIdfIfc}
	cpdbMock := &mockCpdb.CPDBClientInterface{}
	entityUuid, _ := uuid4.New()
	entityUuidStr := entityUuid.String()
	cas := uint64(1)
	entities := []*insights_interface.Entity{
		{
			EntityGuid: &insights_interface.EntityGuid{
				EntityId: &entityUuidStr,
			},
			CasValue: &cas,
		},
	}
	cpdbMock.On("GetEntities", mock.Anything, mock.Anything).Return(entities, nil).Once()
	mockIdfIfc.On("SendMsg", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(insights_interface.ErrPartial).Once()

	err := idfClient.DeleteEntities(context.Background(), cpdbMock, CatalogItem, []string{entityUuidStr}, true)

	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
	cpdbMock.AssertExpectations(t)
	mockIdfIfc.AssertExpectations(t)
}

func TestDeleteEntitiesRetryError(t *testing.T) {
	mockIdfIfc := &mockIdf.InsightsServiceInterface{}
	idfClient := IdfClient{IdfSvc: mockIdfIfc}
	cpdbMock := &mockCpdb.CPDBClientInterface{}
	entityUuid, _ := uuid4.New()
	entityUuidStr := entityUuid.String()
	cas := uint64(1)
	entities := []*insights_interface.Entity{
		{
			EntityGuid: &insights_interface.EntityGuid{
				EntityId: &entityUuidStr,
			},
			CasValue: &cas,
		},
	}
	cpdbMock.On("GetEntities", mock.Anything, mock.Anything).Return(entities, nil).Once()
	mockIdfIfc.On("SendMsg", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			internalError := insights_interface.InsightsErrorProto_kInternalError
			ret := args.Get(2).(*insights_interface.BatchDeleteEntitiesRet)
			ret.RetList = append(ret.RetList, &insights_interface.BatchDeleteEntitiesRet_RetElem{
				Status: &insights_interface.InsightsErrorProto{ErrorType: &internalError},
			})
		}).Return(insights_interface.ErrPartial)

	actualMaxRetries := idfMaxRetries
	idfMaxRetries = 1
	err := idfClient.DeleteEntities(context.Background(), cpdbMock, CatalogItem, []string{entityUuidStr}, true)
	idfMaxRetries = actualMaxRetries

	assert.Error(t, err)
	assert.IsType(t, marinaError.ErrRetry, err)
	cpdbMock.AssertExpectations(t)
	mockIdfIfc.AssertExpectations(t)
}
