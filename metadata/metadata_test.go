/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: shreyash.turkar@nutanix.com
 *
 */

package metadata

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	cpdbMock "github.com/nutanix-core/acs-aos-go/nusights/util/db/mocks"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/content-management-marina/errors"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
)

func getMockEntityMetricAndMetadata(uuid *uuid4.Uuid) (*insights_interface.EntityWithMetric, *marinaIfc.EntityMetadata) {
	createTime := rand.Uint64()
	lastUpdateTime := rand.Uint64()
	ownerUuid, _ := uuid4.New()
	userName := fmt.Sprintf("User-%v", rand.Int63())

	cat1, _ := uuid4.New()
	cat2, _ := uuid4.New()
	categoryStrs := []string{
		cat1.String(),
		cat2.String(),
	}
	categoryBytes := [][]byte{
		cat1.RawBytes(),
		cat2.RawBytes(),
	}

	entityWithMetric := &insights_interface.EntityWithMetric{
		MetricDataList: []*insights_interface.MetricData{
			{
				Name: proto.String(kindIdCol),
				ValueList: []*insights_interface.TimeValuePair{
					{
						Value: &insights_interface.DataValue{
							ValueType: &insights_interface.DataValue_StrValue{
								StrValue: uuid.String(),
							},
						},
					},
				},
			},
			{
				Name: proto.String(createTimeUsecsCol),
				ValueList: []*insights_interface.TimeValuePair{
					{
						Value: &insights_interface.DataValue{
							ValueType: &insights_interface.DataValue_Uint64Value{
								Uint64Value: createTime,
							},
						},
					},
				},
			},
			{
				Name: proto.String(lastUpdatedTimeUsecsCol),
				ValueList: []*insights_interface.TimeValuePair{
					{
						Value: &insights_interface.DataValue{
							ValueType: &insights_interface.DataValue_Uint64Value{
								Uint64Value: lastUpdateTime,
							},
						},
					},
				},
			},
			{
				Name: proto.String(ownerReferenceCol),
				ValueList: []*insights_interface.TimeValuePair{
					{
						Value: &insights_interface.DataValue{
							ValueType: &insights_interface.DataValue_StrValue{
								StrValue: ownerUuid.String(),
							},
						},
					},
				},
			},
			{
				Name: proto.String(ownerUsernameCol),
				ValueList: []*insights_interface.TimeValuePair{
					{
						Value: &insights_interface.DataValue{
							ValueType: &insights_interface.DataValue_StrValue{
								StrValue: userName,
							},
						},
					},
				},
			},
			{
				Name: proto.String(categoryIdListCol),
				ValueList: []*insights_interface.TimeValuePair{
					{
						Value: &insights_interface.DataValue{
							ValueType: &insights_interface.DataValue_StrList_{
								StrList: &insights_interface.DataValue_StrList{
									ValueList: categoryStrs,
								},
							},
						},
					},
				},
			},
		},
	}
	metadata := &marinaIfc.EntityMetadata{
		CreateTimeUsecs:     proto.Int64(int64(createTime)),
		LastUpdateTimeUsecs: proto.Int64(int64(lastUpdateTime)),
		OwnerUserUuid:       ownerUuid.RawBytes(),
		OwnerUserName:       proto.String(userName),
		CategoriesUuidList:  categoryBytes,
	}
	return entityWithMetric, metadata
}

func TestGetEntityMetadataQueryBuilderError(t *testing.T) {
	ctx := context.Background()

	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}

	kindID1, _ := uuid4.New()
	kindID2, _ := uuid4.New()
	kindIDs := []string{
		kindID1.UuidToString(),
		kindID2.UuidToString(),
	}

	idfKind := "mockKind"
	queryName := ""

	metadataUtil := EntityMetadataUtil{}
	metadataByUuid, err := metadataUtil.GetEntityMetadata(ctx, mockCpdbIfc, kindIDs, idfKind, queryName)

	assert.Nil(t, metadataByUuid)
	assert.Error(t, err)
	assert.IsType(t, new(errors.InternalError), err)
	mockCpdbIfc.AssertNotCalled(t, "Query")
}

func TestGetEntityMetadataQueryError(t *testing.T) {
	ctx := context.Background()

	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	mockCpdbIfc.On("Query", mock.Anything, mock.Anything).Return(nil, errors.InternalError{})

	kindID1, _ := uuid4.New()
	kindID2, _ := uuid4.New()
	kindIDs := []string{
		kindID1.UuidToString(),
		kindID2.UuidToString(),
	}

	idfKind := "mockKind"
	queryName := "mockQuery"

	metadataUtil := EntityMetadataUtil{}
	metadataByUuid, err := metadataUtil.GetEntityMetadata(ctx, mockCpdbIfc, kindIDs, idfKind, queryName)

	assert.Nil(t, metadataByUuid)
	assert.Error(t, err)
	assert.IsType(t, new(errors.InternalError), err)
	mockCpdbIfc.AssertExpectations(t)
}

func TestGetEntityMetadataKindIdError(t *testing.T) {
	ctx := context.Background()

	res := []*insights_interface.EntityWithMetric{
		{
			MetricDataList: []*insights_interface.MetricData{
				{
					Name:      proto.String("foo"),
					ValueList: []*insights_interface.TimeValuePair{},
				},
				{
					Name: proto.String(kindIdCol),
					ValueList: []*insights_interface.TimeValuePair{
						{
							Value: &insights_interface.DataValue{
								ValueType: &insights_interface.DataValue_StrValue{
									StrValue: "foo",
								},
							},
						},
					},
				},
			},
		},
	}
	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	mockCpdbIfc.On("Query", mock.Anything, mock.Anything).Return(res, nil)

	kindID1, _ := uuid4.New()
	kindID2, _ := uuid4.New()
	kindIDs := []string{
		kindID1.UuidToString(),
		kindID2.UuidToString(),
	}

	idfKind := "mockKind"
	queryName := "mockQuery"

	metadataUtil := EntityMetadataUtil{}
	metadataByUuid, err := metadataUtil.GetEntityMetadata(ctx, mockCpdbIfc, kindIDs, idfKind, queryName)

	assert.Nil(t, metadataByUuid)
	assert.Error(t, err)
	assert.IsType(t, new(errors.InternalError), err)
	mockCpdbIfc.AssertExpectations(t)
}

func TestGetEntityMetadataOwnerReferenceError(t *testing.T) {
	ctx := context.Background()

	res := []*insights_interface.EntityWithMetric{
		{
			MetricDataList: []*insights_interface.MetricData{
				{
					Name:      proto.String("foo"),
					ValueList: []*insights_interface.TimeValuePair{},
				},
				{
					Name: proto.String(ownerReferenceCol),
					ValueList: []*insights_interface.TimeValuePair{
						{
							Value: &insights_interface.DataValue{
								ValueType: &insights_interface.DataValue_StrValue{
									StrValue: "foo",
								},
							},
						},
					},
				},
			},
		},
	}
	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	mockCpdbIfc.On("Query", mock.Anything, mock.Anything).Return(res, nil)

	kindID1, _ := uuid4.New()
	kindID2, _ := uuid4.New()
	kindIDs := []string{
		kindID1.UuidToString(),
		kindID2.UuidToString(),
	}

	idfKind := "mockKind"
	queryName := "mockQuery"

	metadataUtil := EntityMetadataUtil{}
	metadataByUuid, err := metadataUtil.GetEntityMetadata(ctx, mockCpdbIfc, kindIDs, idfKind, queryName)

	assert.Nil(t, metadataByUuid)
	assert.Error(t, err)
	assert.IsType(t, new(errors.InternalError), err)
	mockCpdbIfc.AssertExpectations(t)
}

func TestGetEntityMetadataCategoryIdListError(t *testing.T) {
	ctx := context.Background()

	res := []*insights_interface.EntityWithMetric{
		{
			MetricDataList: []*insights_interface.MetricData{
				{
					Name:      proto.String("foo"),
					ValueList: []*insights_interface.TimeValuePair{},
				},
				{
					Name: proto.String(categoryIdListCol),
					ValueList: []*insights_interface.TimeValuePair{
						{
							Value: &insights_interface.DataValue{
								ValueType: &insights_interface.DataValue_StrList_{
									StrList: &insights_interface.DataValue_StrList{
										ValueList: []string{"foo"},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	mockCpdbIfc.On("Query", mock.Anything, mock.Anything).Return(res, nil)

	kindID1, _ := uuid4.New()
	kindID2, _ := uuid4.New()
	kindIDs := []string{
		kindID1.UuidToString(),
		kindID2.UuidToString(),
	}

	idfKind := "mockKind"
	queryName := "mockQuery"

	metadataUtil := EntityMetadataUtil{}
	metadataByUuid, err := metadataUtil.GetEntityMetadata(ctx, mockCpdbIfc, kindIDs, idfKind, queryName)

	assert.Nil(t, metadataByUuid)
	assert.Error(t, err)
	assert.IsType(t, new(errors.InternalError), err)
	mockCpdbIfc.AssertExpectations(t)
}

func TestGetEntityMetadata(t *testing.T) {
	ctx := context.Background()

	entityUuid1, _ := uuid4.New()
	entityUuid2, _ := uuid4.New()

	entityWithMetric1, metadata1 := getMockEntityMetricAndMetadata(entityUuid1)
	entityWithMetric2, metadata2 := getMockEntityMetricAndMetadata(entityUuid2)

	res := []*insights_interface.EntityWithMetric{entityWithMetric1, entityWithMetric2}

	mockCpdbIfc := &cpdbMock.CPDBClientInterface{}
	mockCpdbIfc.On("Query", mock.Anything, mock.Anything).Return(res, nil)

	kindIDs := []string{
		entityUuid1.UuidToString(),
		entityUuid2.UuidToString(),
	}

	idfKind := "mockKind"
	queryName := "mockQuery"

	metadataUtil := EntityMetadataUtil{}
	metadataByUuid, err := metadataUtil.GetEntityMetadata(ctx, mockCpdbIfc, kindIDs, idfKind, queryName)

	assert.NoError(t, err)
	assert.Equal(t,
		map[uuid4.Uuid]*marinaIfc.EntityMetadata{*entityUuid1: metadata1, *entityUuid2: metadata2}, metadataByUuid)
	mockCpdbIfc.AssertExpectations(t)
}
