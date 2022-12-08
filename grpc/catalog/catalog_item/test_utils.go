/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 */

package catalog_item

import (
	"bytes"
	"compress/zlib"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	"github.com/golang/protobuf/proto"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	dbMock "github.com/nutanix-core/content-management-marina/mocks/db"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
)

var (
	mockIdfIfc             = &dbMock.IdfClientInterface{}
	testCatalogItemUuid, _ = uuid4.New()
	testCatalogItemVersion = int64(5)
)

func getIdfResponse(catalogItemInfo *marinaIfc.CatalogItemInfo) []*insights_interface.EntityWithMetric {
	serializedProto, _ := proto.Marshal(catalogItemInfo)
	compressedProtoBuf := &bytes.Buffer{}
	zlibWriter := zlib.NewWriter(compressedProtoBuf)
	zlibWriter.Write(serializedProto)
	zlibWriter.Close()

	entities := []*insights_interface.EntityWithMetric{
		{
			MetricDataList: []*insights_interface.MetricData{
				{
					Name: &insights_interface.COMPRESSED_PROTOBUF_ATTR,
					ValueList: []*insights_interface.TimeValuePair{
						{
							Value: &insights_interface.DataValue{
								ValueType: &insights_interface.DataValue_BytesValue{
									BytesValue: compressedProtoBuf.Bytes(),
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
