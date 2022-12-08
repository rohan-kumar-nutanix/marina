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
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
)

func TestCatalogItemGetEmptyArg(t *testing.T) {
	ctx := context.Background()
	arg := &marinaIfc.CatalogItemGetArg{}

	catalogItemInfo := &marinaIfc.CatalogItemInfo{}
	entities := getIdfResponse(catalogItemInfo)
	mockIdfIfc.On("Query", mock.Anything, mock.Anything).
		Return(entities, nil).
		Once()

	ret, err := CatalogItemGet(ctx, arg, mockIdfIfc)
	assert.Nil(t, err)

	catalogItemList := ret.GetCatalogItemList()
	assert.Equal(t, 1, len(catalogItemList))
	assert.Equal(t, catalogItemInfo, catalogItemList[0])
}

func TestCatalogItemGetNonEmptyArg(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	for i := 0; i < *CatalogIdfQueryChunkSize+1; i++ {
		catalogItemIdList = append(catalogItemIdList,
			&marinaIfc.CatalogItemId{GlobalCatalogItemUuid: testCatalogItemUuid.RawBytes()})
	}
	arg := &marinaIfc.CatalogItemGetArg{CatalogItemIdList: catalogItemIdList}

	catalogItemInfo := &marinaIfc.CatalogItemInfo{}
	entities := getIdfResponse(catalogItemInfo)
	mockIdfIfc.On("Query", mock.Anything, mock.Anything).
		Return(entities, nil).
		Twice()

	ret, err := CatalogItemGet(ctx, arg, mockIdfIfc)
	assert.Nil(t, err)

	catalogItemList := ret.GetCatalogItemList()
	assert.Equal(t, 2, len(catalogItemList))
	assert.Equal(t, catalogItemInfo, catalogItemList[0])
	assert.Equal(t, catalogItemInfo, catalogItemList[1])
}

func TestCatalogItemGetEmptyArgError(t *testing.T) {
	ctx := context.Background()
	arg := &marinaIfc.CatalogItemGetArg{}
	mockIdfIfc.On("Query", mock.Anything, mock.Anything).
		Return(nil, insights_interface.ErrInternalError).
		Once()

	_, err := CatalogItemGet(ctx, arg, mockIdfIfc)
	assert.NotNil(t, err)
}

func TestCatalogItemGetNonEmptyArgError(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	for i := 0; i < *CatalogIdfQueryChunkSize+1; i++ {
		catalogItemIdList = append(catalogItemIdList,
			&marinaIfc.CatalogItemId{GlobalCatalogItemUuid: testCatalogItemUuid.RawBytes()})
	}
	arg := &marinaIfc.CatalogItemGetArg{CatalogItemIdList: catalogItemIdList}
	mockIdfIfc.On("Query", mock.Anything, mock.Anything).
		Return(nil, insights_interface.ErrInternalError).
		Twice()

	_, err := CatalogItemGet(ctx, arg, mockIdfIfc)
	assert.NotNil(t, err)
}
