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
	"github.com/nutanix-core/content-management-marina/protos/marina"
)

func getMockEntityMetricAndMetadata(uuid *uuid4.Uuid) (*insights_interface.EntityWithMetric,
	*marina.EntityMetadata) {

	createTime := rand.Int63()
	lastUpdateTime := rand.Int63()
	ownerUuid, _ := uuid4.New()
	userName := fmt.Sprintf("User-%v", rand.Int63())

	cat1, _ := uuid4.New()
	cat2, _ := uuid4.New()
	categoryList := [][]byte{
		cat1.RawBytes(),
		cat2.RawBytes(),
	}

	entityWithMetric := &insights_interface.EntityWithMetric{
		MetricDataList: []*insights_interface.MetricData{
			{
				Name: proto.String(kindID),
				ValueList: []*insights_interface.TimeValuePair{
					{
						Value: &insights_interface.DataValue{
							ValueType: &insights_interface.DataValue_BytesValue{
								BytesValue: uuid.RawBytes(),
							},
						},
					},
				},
			},

			{
				Name: proto.String(createTimeUsecs),
				ValueList: []*insights_interface.TimeValuePair{
					{
						Value: &insights_interface.DataValue{
							ValueType: &insights_interface.DataValue_Int64Value{
								Int64Value: createTime,
							},
						},
					},
				},
			},

			{
				Name: proto.String(lastUpdatedTimeUsecs),
				ValueList: []*insights_interface.TimeValuePair{
					{
						Value: &insights_interface.DataValue{
							ValueType: &insights_interface.DataValue_Int64Value{
								Int64Value: lastUpdateTime,
							},
						},
					},
				},
			},

			{
				Name: proto.String(ownerReference),
				ValueList: []*insights_interface.TimeValuePair{
					{
						Value: &insights_interface.DataValue{
							ValueType: &insights_interface.DataValue_BytesValue{
								BytesValue: ownerUuid.RawBytes(),
							},
						},
					},
				},
			},

			{
				Name: proto.String(ownerUsername),
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
				Name: proto.String(categoryIdList),
				ValueList: []*insights_interface.TimeValuePair{
					{
						Value: &insights_interface.DataValue{
							ValueType: &insights_interface.DataValue_BytesList_{
								BytesList: &insights_interface.DataValue_BytesList{
									ValueList: categoryList,
								},
							},
						},
					},
				},
			},
		},
	}
	metadata := &marina.EntityMetadata{
		CreateTimeUsecs:     &createTime,
		LastUpdateTimeUsecs: &lastUpdateTime,
		OwnerUserUuid:       ownerUuid.RawBytes(),
		OwnerUserName:       &userName,
		CategoriesUuidList:  categoryList,
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
	mockCpdbIfc.On("Query", mock.Anything, mock.Anything).
		Return(nil, errors.InternalError{})

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
	mockCpdbIfc.On("Query", mock.Anything, mock.Anything).
		Return(res, nil)

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
		EntityMetadataByUuid{
			*entityUuid1: metadata1,
			*entityUuid2: metadata2,
		},
		metadataByUuid)
	mockCpdbIfc.AssertExpectations(t)
}
