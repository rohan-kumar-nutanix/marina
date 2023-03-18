/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Author: rajesh.battala@nutanix.com
*
* The implementation for CatalogItemCreate RPC
*
 */

package catalog_item

import (
	"context"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	marinaErrors "github.com/nutanix-core/content-management-marina/errors"
	mockCatalogItem "github.com/nutanix-core/content-management-marina/mocks/grpc/catalog/catalog_item"
	mockFile "github.com/nutanix-core/content-management-marina/mocks/grpc/catalog/file_repo"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	"github.com/nutanix-core/content-management-marina/task/base"
)

var (
	catalogItemInfo = &marinaIfc.CatalogItemInfo{
		GlobalCatalogItemUuid: testCatalogItemUuid.RawBytes(),
		Version:               proto.Int64(0),
	}
)

func TestCatalogItemCreateTask_checkCatalogItemExistsError(t *testing.T) {
	baseTask := &base.MarinaBaseTask{
		ExternalSingletonInterface: mockExternalInterfaces(),
		InternalSingletonInterface: mockInternalInterfaces(),
	}
	citemIfc := new(mockCatalogItem.CatalogItemInterface)
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemBaseTask.globalCatalogItemUuid, _ = uuid4.New()
	catalogItemBaseTask.taskUuid, _ = uuid4.New()
	catalogItemBaseTask.catalogItemIfc = citemIfc

	ret := &marinaIfc.CatalogItemCreateTaskRet{}
	citemCreateTask := NewCatalogItemCreateTask(catalogItemBaseTask)

	citemIfc.On("GetCatalogItems",
		context.TODO(), mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything).Return(
		nil, marinaErrors.ErrInternal)

	exists, err := citemCreateTask.checkCatalogItemExists(ret)
	assert.Error(t, err)
	assert.IsType(t, marinaErrors.ErrInternal, err)
	assert.False(t, exists)
	citemIfc.AssertExpectations(t)
}

func TestCatalogItemCreateTask_checkCatalogItemExists(t *testing.T) {

	baseTask := &base.MarinaBaseTask{
		ExternalSingletonInterface: mockExternalInterfaces(),
		InternalSingletonInterface: mockInternalInterfaces(),
	}
	citemIfc := new(mockCatalogItem.CatalogItemInterface)
	fileRepoIfc = new(mockFile.FileRepoInterface)
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemBaseTask.globalCatalogItemUuid = testCatalogItemUuid
	catalogItemBaseTask.taskUuid, _ = uuid4.New()
	catalogItemBaseTask.catalogItemIfc = citemIfc

	ret := &marinaIfc.CatalogItemCreateTaskRet{}
	citemCreateTask := NewCatalogItemCreateTask(catalogItemBaseTask)
	catalogItems := []*marinaIfc.CatalogItemInfo{catalogItemInfo}
	citemCreateTask.catalogItemUuid = testCatalogItemUuid

	citemIfc.On("GetCatalogItems",
		context.TODO(), cpdbIfc, uuidIfc, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything).Return(
		catalogItems, nil)
	citemIfc.On("GetCatalogItem", mock.Anything, mock.Anything,
		mock.Anything).Return(catalogItemInfo, nil)

	exists, err := citemCreateTask.checkCatalogItemExists(ret)
	assert.Nil(t, err)
	assert.True(t, exists)
	citemIfc.AssertExpectations(t)
}

func TestCatalogItemCreateTask_checkCatalogItemComplete(t *testing.T) {

	baseTask := &base.MarinaBaseTask{
		ExternalSingletonInterface: mockExternalInterfaces(),
		InternalSingletonInterface: mockInternalInterfaces(),
	}
	citemIfc := new(mockCatalogItem.CatalogItemInterface)
	fileRepoIfc = new(mockFile.FileRepoInterface)
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemBaseTask.globalCatalogItemUuid = testCatalogItemUuid
	catalogItemBaseTask.taskUuid, _ = uuid4.New()
	catalogItemBaseTask.catalogItemIfc = citemIfc

	ret := &marinaIfc.CatalogItemCreateTaskRet{}
	citemCreateTask := NewCatalogItemCreateTask(catalogItemBaseTask)
	catalogItems := []*marinaIfc.CatalogItemInfo{catalogItemInfo}
	citemCreateTask.catalogItemUuid = testCatalogItemUuid

	// Case 1: IDF Error
	citemIfc.On("GetCatalogItems",
		context.TODO(), mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything).Return(
		nil, marinaErrors.ErrInternal).Once()
	// Run Case 1
	exists, err := citemCreateTask.checkCatalogItemExists(ret)
	assert.Error(t, err)

	// Case 2: Valid case.
	citemIfc.On("GetCatalogItems",
		context.TODO(), cpdbIfc, uuidIfc, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything).Return(
		catalogItems, nil).Once()
	citemIfc.On("GetCatalogItem", mock.Anything, mock.Anything,
		mock.Anything).Return(catalogItemInfo, nil).Once()
	// Run Case 2
	exists, err = citemCreateTask.checkCatalogItemExists(ret)
	assert.Nil(t, err)
	assert.True(t, exists)
	citemIfc.AssertExpectations(t)

	// Case 3: Multiple version of CatalogItem.
	catalogItems2 := []*marinaIfc.CatalogItemInfo{
		catalogItemInfo, &marinaIfc.CatalogItemInfo{}}

	citemIfc.On("GetCatalogItems",
		context.TODO(), cpdbIfc, uuidIfc, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything).Return(
		catalogItems2, nil)
	citemIfc.On("GetCatalogItem", mock.Anything, mock.Anything,
		mock.Anything).Return(catalogItemInfo, nil)
	// Run Case 3
	exists, err = citemCreateTask.checkCatalogItemExists(ret)
	assert.NotNil(t, err)
	assert.Error(t, marinaErrors.ErrInternalError(), err)
	assert.False(t, exists)
	citemIfc.AssertExpectations(t)

	citemIfc.On("GetCatalogItems",
		context.TODO(), cpdbIfc, uuidIfc, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything).Return(
		catalogItems2, nil)
	citemIfc.On("GetCatalogItem", mock.Anything, mock.Anything,
		mock.Anything).Return(catalogItemInfo, nil)
	// Run Case 3
	exists, err = citemCreateTask.checkCatalogItemExists(ret)
	assert.NotNil(t, err)
	assert.Error(t, marinaErrors.ErrInternalError(), err)
	assert.False(t, exists)
	citemIfc.AssertExpectations(t)
}
