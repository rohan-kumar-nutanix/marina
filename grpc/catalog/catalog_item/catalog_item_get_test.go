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

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	dbMock "github.com/nutanix-core/content-management-marina/mocks/db"
	catalogItemMock "github.com/nutanix-core/content-management-marina/mocks/grpc/catalog/catalog_item"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
)

func TestCatalogItemGetEmptyArg(t *testing.T) {
	ctx := context.Background()
	arg := &marinaIfc.CatalogItemGetArg{}

	mockCatalogItemIfc := &catalogItemMock.CatalogItemInterface{}
	catalogItemInfo := &marinaIfc.CatalogItemInfo{}
	catalogItemInfoList := []*marinaIfc.CatalogItemInfo{catalogItemInfo}
	mockCatalogItemIfc.On("GetCatalogItemsChan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			catalogItemChan := args.Get(5).(chan []*marinaIfc.CatalogItemInfo)
			catalogItemChan <- catalogItemInfoList
			errChan := args.Get(6).(chan error)
			errChan <- nil
		}).
		Once()

	mockIdfIfc := &dbMock.IdfClientInterface{}
	ret, err := CatalogItemGet(ctx, arg, mockCatalogItemIfc, mockIdfIfc)

	assert.Nil(t, err)
	catalogItemList := ret.GetCatalogItemList()
	assert.Equal(t, 1, len(catalogItemList))
	assert.Equal(t, catalogItemInfo, catalogItemList[0])
	mockCatalogItemIfc.AssertExpectations(t)
}

func TestCatalogItemGetNonEmptyArg(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	for i := 0; i < *CatalogIdfQueryChunkSize+1; i++ {
		catalogItemIdList = append(catalogItemIdList,
			&marinaIfc.CatalogItemId{GlobalCatalogItemUuid: testCatalogItemUuid.RawBytes()})
	}
	arg := &marinaIfc.CatalogItemGetArg{CatalogItemIdList: catalogItemIdList}

	mockCatalogItemIfc := &catalogItemMock.CatalogItemInterface{}
	catalogItemInfo := &marinaIfc.CatalogItemInfo{}
	catalogItemInfoList := []*marinaIfc.CatalogItemInfo{catalogItemInfo}
	mockCatalogItemIfc.On("GetCatalogItemsChan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			catalogItemChan := args.Get(5).(chan []*marinaIfc.CatalogItemInfo)
			catalogItemChan <- catalogItemInfoList
			errChan := args.Get(6).(chan error)
			errChan <- nil
		}).
		Twice()

	mockIdfIfc := &dbMock.IdfClientInterface{}
	ret, err := CatalogItemGet(ctx, arg, mockCatalogItemIfc, mockIdfIfc)

	assert.Nil(t, err)
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
		mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			catalogItemChan := args.Get(5).(chan []*marinaIfc.CatalogItemInfo)
			catalogItemChan <- nil
			errChan := args.Get(6).(chan error)
			errChan <- insights_interface.ErrInternalError
		}).
		Once()

	mockIdfIfc := &dbMock.IdfClientInterface{}
	_, err := CatalogItemGet(ctx, arg, mockCatalogItemIfc, mockIdfIfc)

	assert.NotNil(t, err)
	mockCatalogItemIfc.AssertExpectations(t)
}

func TestCatalogItemGetNonEmptyArgError(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	for i := 0; i < *CatalogIdfQueryChunkSize+1; i++ {
		catalogItemIdList = append(catalogItemIdList,
			&marinaIfc.CatalogItemId{GlobalCatalogItemUuid: testCatalogItemUuid.RawBytes()})
	}
	arg := &marinaIfc.CatalogItemGetArg{CatalogItemIdList: catalogItemIdList}

	mockCatalogItemIfc := &catalogItemMock.CatalogItemInterface{}
	mockCatalogItemIfc.On("GetCatalogItemsChan", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			catalogItemChan := args.Get(5).(chan []*marinaIfc.CatalogItemInfo)
			catalogItemChan <- nil
			errChan := args.Get(6).(chan error)
			errChan <- insights_interface.ErrInternalError
		}).
		Twice()

	mockIdfIfc := &dbMock.IdfClientInterface{}
	_, err := CatalogItemGet(ctx, arg, mockCatalogItemIfc, mockIdfIfc)

	assert.NotNil(t, err)
}
