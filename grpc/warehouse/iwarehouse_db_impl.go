package warehouse

import (
	"bytes"
	"compress/zlib"
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	log "k8s.io/klog/v2"

	"github.com/nutanix-core/content-management-marina/db"
	warehousePB "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"

	marinaError "github.com/nutanix-core/content-management-marina/errors"
	utils "github.com/nutanix-core/content-management-marina/util"
)

type WarehouseDBImpl struct {
}

func (warehouseDBImpl *WarehouseDBImpl) CreateWarehouse(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	protoIfc utils.ProtoUtilInterface, warehouseUuid *uuid4.Uuid, warehouseBody *warehousePB.Warehouse) error {
	marshal, err := protoIfc.Marshal(warehouseBody)
	if err != nil {
		errMsg := fmt.Sprintf("Error marshaling the Warehouse proto: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	var buffer bytes.Buffer
	writer := zlib.NewWriter(&buffer)
	_, err = writer.Write(marshal)
	if err != nil {
		errMsg := fmt.Sprintf("Error compressing the Warehouse proto: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	err = writer.Close()
	if err != nil {
		errMsg := fmt.Sprintf("Error encountered while closing zlib stream: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	attrParams := make(cpdb.AttributeParams)
	attrParams[insights_interface.COMPRESSED_PROTOBUF_ATTR] = buffer.Bytes()
	attrVals, err := cpdbIfc.BuildAttributeDataArgs(&attrParams)
	if err != nil {
		errMsg := fmt.Sprintf("Error encountered while generating IDF attributes: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	_, err = cpdbIfc.UpdateEntity(db.Warehouse.ToString(), warehouseUuid.String(), attrVals, nil, false, 0)
	if err != nil {
		errMsg := fmt.Sprintf("Error while creating the IDF entry for Warehouse %s: %v", warehouseUuid.String(), err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}

/*func getIdfAttributesFromWarehouse(warehouse *warehousePB.Warehouse, entityUuid string) []*insights_interface.AttributeDataArg {
	zproto, _ := getZProtoFromEntity(warehouse)

	/*attrs := []*insights_interface.AttributeDataArg{
		// nameProtoAttr(proto.String("name"), warehouse.Name),
		// nameProtoAttr(proto.String("uuid"), entityUuid),
		bytesProtoAttr(proto.String("__zprotobuf__"), zproto.Bytes()),
	}
	attrParams := make(cpdb.AttributeParams)
	attrParams[insights_interface.COMPRESSED_PROTOBUF_ATTR] = zproto.Bytes()
	attrVals, err := cpdbIfc.BuildAttributeDataArgs(&attrParams)

	// attrs = append(attrs, bytesProtoAttr(proto.String("__zprotobuf__"), zproto.Bytes()))
	return attrs
}*/
/*func GetIdfAttributesFromWarehouseItem(warehouseItem *warehousePB.WarehouseItem, entityUuid *string, warehouseUuid *string) []*insights_interface.AttributeDataArg {
	zproto, _ := getWarehouseItemZProtoFromEntity(warehouseItem)

	attrs := []*insights_interface.AttributeDataArg{
		nameProtoAttr(proto.String("name"), warehouseItem.Name),
		nameProtoAttr(proto.String("uuid"), entityUuid),
		nameProtoAttr(proto.String("warehouse_uuid"), warehouseUuid),
		nameProtoAttr(proto.String("object_source_key"), warehouseItem.GetCloudObjectSourceSourceReference().Value.ObjectKey),
		bytesProtoAttr(proto.String("__zprotobuf__"), zproto.Bytes()),
	}
	return attrs
}*/

// AddCompressedProtobufToEntityFromWarehouse Helper function to populate compressed protobuf attribute of the given entity
// with the given task.
func AddCompressedProtobufToEntityFromWarehouse(entity *insights_interface.Entity,
	warehouse *warehousePB.Warehouse) error {
	serializedProto, err := proto.Marshal(warehouse)
	if err != nil {
		log.Error("Failed to serialize task proto : ", err)
		return err
	}
	compressedProtoBuf := &bytes.Buffer{}
	zlibWriter := zlib.NewWriter(compressedProtoBuf)
	_, err = zlibWriter.Write(serializedProto)
	if err != nil {
		return err
	}
	zlibWriter.Close()
	entity.AttributeDataMap = append(entity.AttributeDataMap,
		&insights_interface.NameTimeValuePair{
			Value: &insights_interface.DataValue{
				ValueType: &insights_interface.DataValue_BytesValue{
					BytesValue: compressedProtoBuf.Bytes(),
				},
			},
			Name: proto.String("__zprotobuf__"),
		})
	return nil
}

/*func getZProtoFromEntity(msg proto.Message) (*bytes.Buffer, error) {
	serializedProto, err := internal.Interfaces().ProtoIfc().Marshal(msg)
	if err != nil {
		log.Error("Failed to serialize task proto : ", err)
		return nil, err
	}
	compressedProtoBuf := &bytes.Buffer{}
	zlibWriter := zlib.NewWriter(compressedProtoBuf)
	zlibWriter.Write(serializedProto)
	zlibWriter.Close()
	return compressedProtoBuf, nil
}*/
