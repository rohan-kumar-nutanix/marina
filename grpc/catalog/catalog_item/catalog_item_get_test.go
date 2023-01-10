/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 */

package catalog_item

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	cpdbMock "github.com/nutanix-core/acs-aos-go/nusights/util/db/mocks"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
	catalogItemMock "github.com/nutanix-core/content-management-marina/mocks/grpc/catalog/catalog_item"
	utilsMock "github.com/nutanix-core/content-management-marina/mocks/util"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
)

func TestCatalogItemGetEmptyArg(t *testing.T) {
	ctx := context.Background()
	arg := &marinaIfc.CatalogItemGetArg{}

	mockCatalogItemIfc := &catalogItemMock.CatalogItemInterface{}
	catalogItemInfo := &marinaIfc.CatalogItemInfo{}
	catalogItemInfoList := []*marinaIfc.CatalogItemInfo{catalogItemInfo}
	mockCatalogItemIfc.On("GetCatalogItemsChan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			catalogItemChan := args.Get(6).(chan []*marinaIfc.CatalogItemInfo)
			catalogItemChan <- catalogItemInfoList
			errChan := args.Get(7).(chan error)
			errChan <- nil
		}).
		Once()

	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	mockUuidIfc := &utilsMock.UuidUtilInterface{}
	ret, err := CatalogItemGet(ctx, arg, mockCatalogItemIfc, mockCpdbIfc, mockUuidIfc)

	assert.NoError(t, err)
	catalogItemList := ret.GetCatalogItemList()
	assert.Equal(t, 1, len(catalogItemList))
	assert.Equal(t, catalogItemInfo, catalogItemList[0])
	mockCatalogItemIfc.AssertExpectations(t)
}

func TestCatalogItemGetNonEmptyArg(t *testing.T) {
	ctx := context.Background()
	uuid, _ := uuid4.New()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	for i := 0; i < *CatalogIdfQueryChunkSize+1; i++ {
		catalogItemIdList = append(catalogItemIdList,
			&marinaIfc.CatalogItemId{GlobalCatalogItemUuid: uuid.RawBytes()})
	}
	arg := &marinaIfc.CatalogItemGetArg{CatalogItemIdList: catalogItemIdList}

	mockCatalogItemIfc := &catalogItemMock.CatalogItemInterface{}
	catalogItemInfo := &marinaIfc.CatalogItemInfo{}
	catalogItemInfoList := []*marinaIfc.CatalogItemInfo{catalogItemInfo}
	mockCatalogItemIfc.On("GetCatalogItemsChan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			catalogItemChan := args.Get(6).(chan []*marinaIfc.CatalogItemInfo)
			catalogItemChan <- catalogItemInfoList
			errChan := args.Get(7).(chan error)
			errChan <- nil
		}).
		Twice()

	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	mockUuidIfc := &utilsMock.UuidUtilInterface{}
	ret, err := CatalogItemGet(ctx, arg, mockCatalogItemIfc, mockCpdbIfc, mockUuidIfc)

	assert.NoError(t, err)
	catalogItemList := ret.GetCatalogItemList()
	assert.Equal(t, 2, len(catalogItemList))
	assert.Equal(t, catalogItemInfo, catalogItemList[0])
	assert.Equal(t, catalogItemInfo, catalogItemList[1])
	mockCatalogItemIfc.AssertExpectations(t)
}

func TestCatalogItemGetEmptyArgError(t *testing.T) {
	ctx := context.Background()
	arg := &marinaIfc.CatalogItemGetArg{}

	mockCatalogItemIfc := &catalogItemMock.CatalogItemInterface{}
	mockCatalogItemIfc.On("GetCatalogItemsChan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			catalogItemChan := args.Get(6).(chan []*marinaIfc.CatalogItemInfo)
			catalogItemChan <- nil
			errChan := args.Get(7).(chan error)
			errChan <- marinaError.ErrInternalError()
		}).
		Once()

	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	mockUuidIfc := &utilsMock.UuidUtilInterface{}
	_, err := CatalogItemGet(ctx, arg, mockCatalogItemIfc, mockCpdbIfc, mockUuidIfc)

	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
	mockCatalogItemIfc.AssertExpectations(t)
}

func TestCatalogItemGetNonEmptyArgError(t *testing.T) {
	ctx := context.Background()
	uuid, _ := uuid4.New()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	for i := 0; i < *CatalogIdfQueryChunkSize+1; i++ {
		catalogItemIdList = append(catalogItemIdList,
			&marinaIfc.CatalogItemId{GlobalCatalogItemUuid: uuid.RawBytes()})
	}
	arg := &marinaIfc.CatalogItemGetArg{CatalogItemIdList: catalogItemIdList}

	mockCatalogItemIfc := &catalogItemMock.CatalogItemInterface{}
	mockCatalogItemIfc.On("GetCatalogItemsChan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			catalogItemChan := args.Get(6).(chan []*marinaIfc.CatalogItemInfo)
			catalogItemChan <- nil
			errChan := args.Get(7).(chan error)
			errChan <- marinaError.ErrInternalError()
		}).
		Twice()

	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	mockUuidIfc := &utilsMock.UuidUtilInterface{}
	_, err := CatalogItemGet(ctx, arg, mockCatalogItemIfc, mockCpdbIfc, mockUuidIfc)

	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
}
