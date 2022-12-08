/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 */

package catalog_item

import (
	"context"
	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	"testing"

	"github.com/stretchr/testify/mock"

	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	"github.com/stretchr/testify/assert"
)

func TestGetCatalogItemsDefaultArgs(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType

	catalogItemInfo := &marinaIfc.CatalogItemInfo{}
	entities := getIdfResponse(catalogItemInfo)
	mockIdfIfc.On("Query", mock.Anything, mock.Anything).
		Return(entities, nil).
		Once()

	catalogItemList, err := getCatalogItems(ctx, mockIdfIfc, catalogItemIdList, catalogItemTypeList, false, "test_query")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(catalogItemList))
	assert.Equal(t, catalogItemInfo, catalogItemList[0])
}

func TestGetCatalogItemsTypeList(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	catalogItemTypeList := []marinaIfc.CatalogItemInfo_CatalogItemType{marinaIfc.CatalogItemInfo_kImage}

	catalogItemInfo := &marinaIfc.CatalogItemInfo{}
	entities := getIdfResponse(catalogItemInfo)
	mockIdfIfc.On("Query", mock.Anything, mock.Anything).
		Return(entities, nil).
		Once()

	catalogItemList, err := getCatalogItems(ctx, mockIdfIfc, catalogItemIdList, catalogItemTypeList, false, "test_query")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(catalogItemList))
	assert.Equal(t, catalogItemInfo, catalogItemList[0])
}

func TestGetCatalogItemsAllArgs(t *testing.T) {
	ctx := context.Background()
	catalogItemIdList := []*marinaIfc.CatalogItemId{
		{GlobalCatalogItemUuid: testCatalogItemUuid.RawBytes()},
		{GlobalCatalogItemUuid: testCatalogItemUuid.RawBytes(), Version: &testCatalogItemVersion},
	}
	catalogItemTypeList := []marinaIfc.CatalogItemInfo_CatalogItemType{marinaIfc.CatalogItemInfo_kImage}

	catalogItemInfo := &marinaIfc.CatalogItemInfo{}
	entities := getIdfResponse(catalogItemInfo)
	mockIdfIfc.On("Query", mock.Anything, mock.Anything).
		Return(entities, nil).
		Once()

	catalogItemList, err := getCatalogItems(ctx, mockIdfIfc, catalogItemIdList, catalogItemTypeList, true, "test_query")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(catalogItemList))
	assert.Equal(t, catalogItemInfo, catalogItemList[0])
}

func TestGetCatalogItemsNotFound(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType

	mockIdfIfc.On("Query", mock.Anything, mock.Anything).
		Return(nil, insights_interface.ErrNotFound).
		Once()

	catalogItemList, err := getCatalogItems(ctx, mockIdfIfc, catalogItemIdList, catalogItemTypeList, false, "test_query")
	assert.Nil(t, err)
	assert.Empty(t, catalogItemList)
}

func TestGetCatalogItemsInvalidUuidError(t *testing.T) {
	ctx := context.Background()
	catalogItemIdList := []*marinaIfc.CatalogItemId{{GlobalCatalogItemUuid: []byte("Invalid UUID")}}
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType

	_, err := getCatalogItems(ctx, mockIdfIfc, catalogItemIdList, catalogItemTypeList, false, "test_query")
	assert.NotNil(t, err)
}

func TestGetCatalogItemsInvalidQueryNameError(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType

	_, err := getCatalogItems(ctx, mockIdfIfc, catalogItemIdList, catalogItemTypeList, false, "")
	assert.NotNil(t, err)
}

func TestGetCatalogItemsDeserializeError(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType

	entities := []*insights_interface.EntityWithMetric{
		{
			MetricDataList: []*insights_interface.MetricData{
				{
					Name: &insights_interface.COMPRESSED_PROTOBUF_ATTR,
					ValueList: []*insights_interface.TimeValuePair{
						{
							Value: &insights_interface.DataValue{
								ValueType: &insights_interface.DataValue_BytesValue{
									BytesValue: []byte("Invalid Data"),
								},
							},
						},
					},
				},
			},
		},
	}
	mockIdfIfc.On("Query", mock.Anything, mock.Anything).
		Return(entities, nil).
		Once()

	_, err := getCatalogItems(ctx, mockIdfIfc, catalogItemIdList, catalogItemTypeList, false, "test_query")
	assert.NotNil(t, err)
}
