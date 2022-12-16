/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 */

package catalog_item

import (
	"bytes"
	"compress/zlib"
	"context"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	dbMock "github.com/nutanix-core/content-management-marina/mocks/db"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
)

func getIdfResponse(catalogItemInfo *marinaIfc.CatalogItemInfo) []*insights_interface.EntityWithMetric {
	serializedProto, _ := proto.Marshal(catalogItemInfo)
	compressedProtoBuf := &bytes.Buffer{}
	zlibWriter := zlib.NewWriter(compressedProtoBuf)
	zlibWriter.Write(serializedProto)
	zlibWriter.Close()
	return entityListFromBytes(compressedProtoBuf.Bytes())
}

func entityListFromBytes(value []byte) []*insights_interface.EntityWithMetric {
	entities := []*insights_interface.EntityWithMetric{
		{
			MetricDataList: []*insights_interface.MetricData{
				{
					Name: &insights_interface.COMPRESSED_PROTOBUF_ATTR,
					ValueList: []*insights_interface.TimeValuePair{
						{
							Value: &insights_interface.DataValue{
								ValueType: &insights_interface.DataValue_BytesValue{
									BytesValue: value,
								},
							},
						},
					},
				},
			},
		},
	}
	return entities
}

func TestGetCatalogItemsChan(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType

	mockIdfIfc := &dbMock.IdfClientInterface{}
	catalogItemInfo := &marinaIfc.CatalogItemInfo{}
	entities := getIdfResponse(catalogItemInfo)
	mockIdfIfc.On("Query", mock.Anything, mock.Anything).
		Return(entities, nil).
		Once()

	catalogItemChan := make(chan []*marinaIfc.CatalogItemInfo)
	errChan := make(chan error)

	catalogItemService := new(CatalogItemImpl)
	go catalogItemService.GetCatalogItemsChan(ctx, mockIdfIfc, catalogItemIdList, catalogItemTypeList,
		false, catalogItemChan, errChan)

	catalogItemList := <-catalogItemChan
	err := <-errChan
	assert.Equal(t, 1, len(catalogItemList))
	assert.Equal(t, catalogItemInfo, catalogItemList[0])
	assert.Nil(t, err)
	mockIdfIfc.AssertExpectations(t)
}

func TestGetCatalogItemsDefaultArgs(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType

	mockIdfIfc := &dbMock.IdfClientInterface{}
	catalogItemInfo := &marinaIfc.CatalogItemInfo{}
	entities := getIdfResponse(catalogItemInfo)
	mockIdfIfc.On("Query", mock.Anything, mock.Anything).
		Return(entities, nil).
		Once()

	catalogItemService := new(CatalogItemImpl)
	catalogItemList, err := catalogItemService.GetCatalogItems(ctx, mockIdfIfc, catalogItemIdList, catalogItemTypeList,
		false, "test_query")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(catalogItemList))
	assert.Equal(t, catalogItemInfo, catalogItemList[0])
	mockIdfIfc.AssertExpectations(t)
}

func TestGetCatalogItemsTypeList(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	catalogItemTypeList := []marinaIfc.CatalogItemInfo_CatalogItemType{marinaIfc.CatalogItemInfo_kImage}

	mockIdfIfc := &dbMock.IdfClientInterface{}
	catalogItemInfo := &marinaIfc.CatalogItemInfo{}
	entities := getIdfResponse(catalogItemInfo)
	mockIdfIfc.On("Query", mock.Anything, mock.Anything).
		Return(entities, nil).
		Once()

	catalogItemService := new(CatalogItemImpl)
	catalogItemList, err := catalogItemService.GetCatalogItems(ctx, mockIdfIfc, catalogItemIdList, catalogItemTypeList,
		false, "test_query")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(catalogItemList))
	assert.Equal(t, catalogItemInfo, catalogItemList[0])
	mockIdfIfc.AssertExpectations(t)
}

func TestGetCatalogItemsAllArgs(t *testing.T) {
	ctx := context.Background()
	catalogItemIdList := []*marinaIfc.CatalogItemId{
		{GlobalCatalogItemUuid: testCatalogItemUuid.RawBytes()},
		{GlobalCatalogItemUuid: testCatalogItemUuid.RawBytes(), Version: &testCatalogItemVersion},
	}
	catalogItemTypeList := []marinaIfc.CatalogItemInfo_CatalogItemType{marinaIfc.CatalogItemInfo_kImage}

	mockIdfIfc := &dbMock.IdfClientInterface{}
	catalogItemInfo := &marinaIfc.CatalogItemInfo{}
	entities := getIdfResponse(catalogItemInfo)
	mockIdfIfc.On("Query", mock.Anything, mock.Anything).
		Return(entities, nil).
		Once()

	catalogItemService := new(CatalogItemImpl)
	catalogItemList, err := catalogItemService.GetCatalogItems(ctx, mockIdfIfc, catalogItemIdList, catalogItemTypeList,
		true, "test_query")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(catalogItemList))
	assert.Equal(t, catalogItemInfo, catalogItemList[0])
	mockIdfIfc.AssertExpectations(t)
}

func TestGetCatalogItemsNotFound(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType

	mockIdfIfc := &dbMock.IdfClientInterface{}
	mockIdfIfc.On("Query", mock.Anything, mock.Anything).
		Return(nil, insights_interface.ErrNotFound).
		Once()

	catalogItemService := new(CatalogItemImpl)
	catalogItemList, err := catalogItemService.GetCatalogItems(ctx, mockIdfIfc, catalogItemIdList, catalogItemTypeList,
		false, "test_query")
	assert.Nil(t, err)
	assert.Empty(t, catalogItemList)
	mockIdfIfc.AssertExpectations(t)
}

func TestGetCatalogItemsInvalidUuidError(t *testing.T) {
	ctx := context.Background()
	catalogItemIdList := []*marinaIfc.CatalogItemId{{GlobalCatalogItemUuid: []byte("Invalid UUID")}}
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType

	mockIdfIfc := &dbMock.IdfClientInterface{}
	catalogItemService := new(CatalogItemImpl)
	_, err := catalogItemService.GetCatalogItems(ctx, mockIdfIfc, catalogItemIdList, catalogItemTypeList, false, "test_query")
	assert.NotNil(t, err)
}

func TestGetCatalogItemsInvalidQueryNameError(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType

	mockIdfIfc := &dbMock.IdfClientInterface{}
	catalogItemService := new(CatalogItemImpl)
	_, err := catalogItemService.GetCatalogItems(ctx, mockIdfIfc, catalogItemIdList, catalogItemTypeList, false, "")
	assert.NotNil(t, err)
}

func TestGetCatalogItemsIdfInternalError(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType

	mockIdfIfc := &dbMock.IdfClientInterface{}
	mockIdfIfc.On("Query", mock.Anything, mock.Anything).
		Return(nil, insights_interface.ErrInternalError).
		Once()

	catalogItemService := new(CatalogItemImpl)
	_, err := catalogItemService.GetCatalogItems(ctx, mockIdfIfc, catalogItemIdList, catalogItemTypeList, false, "test_query")

	assert.NotNil(t, err)
	mockIdfIfc.AssertExpectations(t)
}

func TestGetCatalogItemsDeserializeError(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType

	mockIdfIfc := &dbMock.IdfClientInterface{}
	entities := entityListFromBytes([]byte("Invalid Data"))
	mockIdfIfc.On("Query", mock.Anything, mock.Anything).
		Return(entities, nil).
		Once()

	catalogItemService := new(CatalogItemImpl)
	_, err := catalogItemService.GetCatalogItems(ctx, mockIdfIfc, catalogItemIdList, catalogItemTypeList, false, "test_query")

	assert.NotNil(t, err)
	mockIdfIfc.AssertExpectations(t)
}
