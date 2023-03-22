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
	cpdbMock "github.com/nutanix-core/acs-aos-go/nusights/util/db/mocks"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	marinaError "github.com/nutanix-core/content-management-marina/errors"
	utilsMock "github.com/nutanix-core/content-management-marina/mocks/util"
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

	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	catalogItemInfo := &marinaIfc.CatalogItemInfo{}
	entities := getIdfResponse(catalogItemInfo)
	mockCpdbIfc.On("Query", mock.Anything).
		Return(entities, nil).
		Once()

	catalogItemChan := make(chan []*marinaIfc.CatalogItemInfo)
	errChan := make(chan error)
	mockUuidIfc := &utilsMock.UuidUtilInterface{}

	catalogItemService := new(CatalogItemImpl)
	go catalogItemService.GetCatalogItemsChan(ctx, mockCpdbIfc, mockUuidIfc, catalogItemIdList, catalogItemTypeList,
		false, catalogItemChan, errChan)

	catalogItemList := <-catalogItemChan
	err := <-errChan
	assert.Equal(t, 1, len(catalogItemList))
	assert.Equal(t, catalogItemInfo, catalogItemList[0])
	assert.NoError(t, err)
	mockCpdbIfc.AssertExpectations(t)
}

func TestGetCatalogItemsDefaultArgs(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType

	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	catalogItemInfo := &marinaIfc.CatalogItemInfo{}
	entities := getIdfResponse(catalogItemInfo)
	mockCpdbIfc.On("Query", mock.Anything).
		Return(entities, nil).
		Once()
	mockUuidIfc := &utilsMock.UuidUtilInterface{}

	catalogItemService := new(CatalogItemImpl)
	catalogItemList, err := catalogItemService.GetCatalogItems(ctx, mockCpdbIfc, mockUuidIfc, catalogItemIdList,
		catalogItemTypeList, false, "test_query")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(catalogItemList))
	assert.Equal(t, catalogItemInfo, catalogItemList[0])
	mockCpdbIfc.AssertExpectations(t)
}

func TestGetCatalogItemsTypeList(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	catalogItemTypeList := []marinaIfc.CatalogItemInfo_CatalogItemType{marinaIfc.CatalogItemInfo_kImage}

	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	catalogItemInfo := &marinaIfc.CatalogItemInfo{}
	entities := getIdfResponse(catalogItemInfo)
	mockCpdbIfc.On("Query", mock.Anything).
		Return(entities, nil).
		Once()
	mockUuidIfc := &utilsMock.UuidUtilInterface{}

	catalogItemService := new(CatalogItemImpl)
	catalogItemList, err := catalogItemService.GetCatalogItems(ctx, mockCpdbIfc, mockUuidIfc, catalogItemIdList,
		catalogItemTypeList, false, "test_query")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(catalogItemList))
	assert.Equal(t, catalogItemInfo, catalogItemList[0])
	mockCpdbIfc.AssertExpectations(t)
}

func TestGetCatalogItemsAllArgs(t *testing.T) {
	ctx := context.Background()
	uuid, _ := uuid4.New()
	catalogItemIdList := []*marinaIfc.CatalogItemId{
		{GlobalCatalogItemUuid: uuid.RawBytes()},
		{GlobalCatalogItemUuid: uuid.RawBytes(), Version: proto.Int64(5)},
	}
	catalogItemTypeList := []marinaIfc.CatalogItemInfo_CatalogItemType{marinaIfc.CatalogItemInfo_kImage}

	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	catalogItemInfo := &marinaIfc.CatalogItemInfo{}
	entities := getIdfResponse(catalogItemInfo)
	mockCpdbIfc.On("Query", mock.Anything, mock.Anything).
		Return(entities, nil).
		Once()
	mockUuidIfc := &utilsMock.UuidUtilInterface{}
	mockUuidIfc.On("ValidateUUID", mock.Anything, mock.Anything).Return(nil).Twice()

	catalogItemService := new(CatalogItemImpl)
	catalogItemList, err := catalogItemService.GetCatalogItems(ctx, mockCpdbIfc, mockUuidIfc, catalogItemIdList,
		catalogItemTypeList, true, "test_query")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(catalogItemList))
	assert.Equal(t, catalogItemInfo, catalogItemList[0])
	mockCpdbIfc.AssertExpectations(t)
	mockUuidIfc.AssertExpectations(t)
}

func TestGetCatalogItemsNotFound(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType

	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	mockCpdbIfc.On("Query", mock.Anything).
		Return(nil, insights_interface.ErrNotFound).
		Once()
	mockUuidIfc := &utilsMock.UuidUtilInterface{}

	catalogItemService := new(CatalogItemImpl)
	catalogItemList, err := catalogItemService.GetCatalogItems(ctx, mockCpdbIfc, mockUuidIfc, catalogItemIdList,
		catalogItemTypeList, false, "test_query")
	assert.NoError(t, err)
	assert.Empty(t, catalogItemList)
	mockCpdbIfc.AssertExpectations(t)
}

func TestGetCatalogItemsInvalidUuidError(t *testing.T) {
	ctx := context.Background()
	catalogItemIdList := []*marinaIfc.CatalogItemId{{GlobalCatalogItemUuid: []byte("Invalid UUID")}}
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType
	mockUuidIfc := &utilsMock.UuidUtilInterface{}
	mockUuidIfc.On("ValidateUUID", mock.Anything, mock.Anything).Return(marinaError.ErrInternalError()).Once()

	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	catalogItemService := new(CatalogItemImpl)
	_, err := catalogItemService.GetCatalogItems(ctx, mockCpdbIfc, mockUuidIfc, catalogItemIdList, catalogItemTypeList,
		false, "test_query")
	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
	mockUuidIfc.AssertExpectations(t)
}

func TestGetCatalogItemsInvalidQueryNameError(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType
	mockUuidIfc := &utilsMock.UuidUtilInterface{}

	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	catalogItemService := new(CatalogItemImpl)
	_, err := catalogItemService.GetCatalogItems(ctx, mockCpdbIfc, mockUuidIfc, catalogItemIdList, catalogItemTypeList,
		false, "")
	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
}

func TestGetCatalogItemsIdfInternalError(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType

	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	mockCpdbIfc.On("Query", mock.Anything).
		Return(nil, insights_interface.ErrInternalError).
		Once()
	mockUuidIfc := &utilsMock.UuidUtilInterface{}

	catalogItemService := new(CatalogItemImpl)
	_, err := catalogItemService.GetCatalogItems(ctx, mockCpdbIfc, mockUuidIfc, catalogItemIdList, catalogItemTypeList,
		false, "test_query")

	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
	mockCpdbIfc.AssertExpectations(t)
}

func TestGetCatalogItemsDeserializeError(t *testing.T) {
	ctx := context.Background()
	var catalogItemIdList []*marinaIfc.CatalogItemId
	var catalogItemTypeList []marinaIfc.CatalogItemInfo_CatalogItemType

	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	entities := entityListFromBytes([]byte("Invalid Data"))
	mockCpdbIfc.On("Query", mock.Anything).
		Return(entities, nil).
		Once()
	mockUuidIfc := &utilsMock.UuidUtilInterface{}

	catalogItemService := new(CatalogItemImpl)
	_, err := catalogItemService.GetCatalogItems(ctx, mockCpdbIfc, mockUuidIfc, catalogItemIdList, catalogItemTypeList,
		false, "test_query")

	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
	mockCpdbIfc.AssertExpectations(t)
}
