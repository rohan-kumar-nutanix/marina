package warehouse

import (
	"bytes"
	"compress/zlib"
	"context"
	"errors"
	"fmt"
	"sync"

	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	"github.com/nutanix-core/content-management-marina/db"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
	"github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"
	"github.com/nutanix-core/content-management-marina/protos/apis/common/v1/response"
	utils "github.com/nutanix-core/content-management-marina/util"
)

var (
	warehouseDBImpl IWarehouseDB = nil
	once            sync.Once
)

type WarehouseDBImpl struct {
}

func newWarehouseDBImpl() IWarehouseDB {
	once.Do(func() {
		warehouseDBImpl = &WarehouseDBImpl{}
	})
	return warehouseDBImpl
}
func (warehouseDBImpl *WarehouseDBImpl) CreateWarehouse(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	protoIfc utils.ProtoUtilInterface, warehouseUuid *uuid4.Uuid, warehouseBody *content.Warehouse) error {
	log.Infof("Warehouse body getting persisted to IDF %v", warehouseBody)
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

func (warehouseDBImpl *WarehouseDBImpl) GetWarehouse(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	warehouseUuid *uuid4.Uuid) (*content.Warehouse, error) {
	entity, err := cpdbIfc.GetEntity(db.Warehouse.ToString(), warehouseUuid.String(), false)
	if err == insights_interface.ErrNotFound || entity == nil {
		log.Errorf("Warehouse %s not found", warehouseUuid.String())
		return nil, marinaError.ErrNotFound
	} else if err != nil {
		errMsg := fmt.Sprintf("Error encountered while getting Warehouse %s: %v", warehouseUuid.String(), err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	warehouse := &content.Warehouse{}
	err = entity.DeserializeEntity(warehouse)
	if err != nil {
		errMsg := fmt.Sprintf("failed to deserialize Warehouse IDF entry: %v", err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	if warehouse.GetBase() == nil {
		warehouse.Base = &response.ExternalizableAbstractModel{}
	}
	warehouse.Base.ExtId = entity.EntityGuid.EntityId
	return warehouse, nil
}

func (warehouseDBImpl *WarehouseDBImpl) DeleteWarehouse(ctx context.Context, idfIfc db.IdfClientInterface,
	cpdbIfc cpdb.CPDBClientInterface, warehouseUuid string) error {
	err := idfIfc.DeleteEntities(context.Background(), cpdbIfc, db.Warehouse, []string{warehouseUuid}, true)
	if err == insights_interface.ErrNotFound {
		log.Errorf("Warehouse UUID :%v do not exist in IDF", err)
		return nil

	} else if err != nil {
		errMsg := fmt.Sprintf("Failed to delete the Warehouse: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}

func (warehouseDBImpl *WarehouseDBImpl) UpdateWarehouse(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	protoIfc utils.ProtoUtilInterface, warehouseUuid *uuid4.Uuid, warehouseBody *content.Warehouse) error {

	entity, err := cpdbIfc.GetEntity(db.Warehouse.ToString(), warehouseUuid.String(), true)
	if err == insights_interface.ErrNotFound || entity == nil {
		log.Errorf("Warehouse %s not found", warehouseUuid.String())
		return marinaError.ErrNotFound
	} else if err != nil {
		errMsg := fmt.Sprintf("Error encountered while getting Warehouse %s: %v", warehouseUuid.String(), err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

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

	_, err = cpdbIfc.UpdateEntity(db.Warehouse.ToString(), warehouseUuid.String(), attrVals,
		entity, false, entity.GetCasValue()+1)
	if err != nil {
		errMsg := fmt.Sprintf("Error while creating the IDF entry for Warehouse %s: %v", warehouseUuid.String(), err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}

func (warehouseDBImpl *WarehouseDBImpl) ListWarehouses(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface) (
	[]*content.Warehouse, error) {
	var warehouses []*content.Warehouse
	entities, err := cpdbIfc.GetEntitiesOfType(db.Warehouse.ToString(), false)
	if err != nil {
		return nil, err
	}
	for _, entity := range entities {
		warehouse := &content.Warehouse{}
		log.Infof("entity %v", entity)
		err := entity.DeserializeEntity(warehouse)
		if err != nil {
			log.Errorf("error occurred in deserializing the Warehouse Entity and skipping it. %v", err)
		}
		// TODO Remove this hack
		if warehouse.GetBase() == nil {
			warehouse.Base = &response.ExternalizableAbstractModel{}
		}
		warehouse.Base.ExtId = entity.EntityGuid.EntityId

		warehouses = append(warehouses, warehouse)
	}
	return warehouses, nil
}

func (warehouseDBImpl *WarehouseDBImpl) CreateWarehouseItem(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	protoIfc utils.ProtoUtilInterface, warehouseUuid *uuid4.Uuid, warehouseItemUuid *uuid4.Uuid,
	warehouseItemPB *content.WarehouseItem) error {
	log.Infof("WarehouseItem paylod getting persisted to IDF %v", warehouseItemPB)
	marshal, err := protoIfc.Marshal(warehouseItemPB)
	if err != nil {
		errMsg := fmt.Sprintf("Error marshaling the WarehouseItem proto: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	var buffer bytes.Buffer
	writer := zlib.NewWriter(&buffer)
	_, err = writer.Write(marshal)
	if err != nil {
		errMsg := fmt.Sprintf("Error compressing the WarehouseItem proto: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	err = writer.Close()
	if err != nil {
		errMsg := fmt.Sprintf("Error encountered while closing zlib stream: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	attrParams := make(cpdb.AttributeParams)
	attrParams[insights_interface.COMPRESSED_PROTOBUF_ATTR] = buffer.Bytes()
	attrParams["warehouse_uuid"] = warehouseUuid.String()
	attrVals, err := cpdbIfc.BuildAttributeDataArgs(&attrParams)
	if err != nil {
		errMsg := fmt.Sprintf("Error encountered while generating IDF attributes: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	_, err = cpdbIfc.UpdateEntity(db.WarehouseItem.ToString(), warehouseItemUuid.String(), attrVals, nil, false, 0)
	if err != nil {
		errMsg := fmt.Sprintf("Error while creating the IDF entry for WarehouseItem %s: %v", warehouseItemUuid.String(), err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return nil

}

func (warehouseDBImpl *WarehouseDBImpl) GetWarehouseItem(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	warehouseItemUuid *uuid4.Uuid) (*content.WarehouseItem, error) {
	entity, err := cpdbIfc.GetEntity(db.WarehouseItem.ToString(), warehouseItemUuid.String(), false)
	if err == insights_interface.ErrNotFound || entity == nil {
		log.Errorf("WarehouseItem %s not found", warehouseItemUuid.String())
		return nil, marinaError.ErrNotFound
	} else if err != nil {
		errMsg := fmt.Sprintf("Error encountered while getting WarehouseItem %s: %v", warehouseItemUuid.String(), err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	warehouseItem := &content.WarehouseItem{}
	err = entity.DeserializeEntity(warehouseItem)
	if err != nil {
		errMsg := fmt.Sprintf("failed to deserialize WarehouseItem IDF entry: %v", err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return warehouseItem, nil
}

func (warehouseDBImpl *WarehouseDBImpl) ListWarehouseItems(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, warehouseUUID *uuid4.Uuid) (
	[]*content.WarehouseItem, error) {
	var warehouseItems []*content.WarehouseItem
	entities, err := cpdbIfc.GetEntitiesByAttributeValue(db.WarehouseItem.ToString(),
		"warehouse_uuid", warehouseUUID.String(), []string{insights_interface.COMPRESSED_PROTOBUF_ATTR})
	if err == insights_interface.ErrNotFound || entities == nil {
		log.Errorf("WarehouseItems for Warehouse %s not found", warehouseUUID.String())
		return warehouseItems, marinaError.ErrNotFound
	} else if err != nil {
		errMsg := fmt.Sprintf("Error encountered while getting Warehouse %s: %v", warehouseUUID.String(), err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	for _, entity := range entities {
		log.Infof("WarehouseItem Entity %v", entity)
		warehouseItem := &content.WarehouseItem{}
		err := entity.DeserializeEntity(warehouseItem)
		if err != nil {
			log.Errorf("error occurred in deserializing the WarehouseItem Entity and skipping it. %v", err)
		}
		warehouseItems = append(warehouseItems, warehouseItem)
	}
	return warehouseItems, nil
}

func (warehouseDBImpl *WarehouseDBImpl) DeleteWarehouseItem(ctx context.Context, idfIfc db.IdfClientInterface,
	cpdbIfc cpdb.CPDBClientInterface, warehouseUuid string, warehouseItemUuid string) error {
	err := idfIfc.DeleteEntities(context.Background(), cpdbIfc, db.WarehouseItem, []string{warehouseItemUuid}, true)
	if err == insights_interface.ErrNotFound {
		log.Errorf("WarehouseItemUuid UUID :%v do not exist in IDF", err)
		return nil

	} else if err != nil {
		errMsg := fmt.Sprintf("Failed to delete the WarehouseItem: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	log.Infof("WarehouseItem UUID : %s is deleted", warehouseUuid)
	return nil
}

func (warehouseDBImpl *WarehouseDBImpl) UpdateWarehouseItem(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	protoIfc utils.ProtoUtilInterface, warehouseUuid *uuid4.Uuid, warehouseItemUuid *uuid4.Uuid,
	warehouseItemPB *content.WarehouseItem) error {

	entity, err := cpdbIfc.GetEntity(db.WarehouseItem.ToString(), warehouseItemUuid.String(), false)
	if err == insights_interface.ErrNotFound || entity == nil {
		log.Errorf("WarehouseItem %s not found", warehouseUuid.String())
		return marinaError.ErrNotFound
	} else if err != nil {
		errMsg := fmt.Sprintf("Error encountered while getting WarehouseItem %s: %v", warehouseUuid.String(), err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	marshal, err := protoIfc.Marshal(warehouseItemPB)
	if err != nil {
		errMsg := fmt.Sprintf("Error marshaling the WarehouseItem proto: %v", err)
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

	_, err = cpdbIfc.UpdateEntity(db.WarehouseItem.ToString(), warehouseItemUuid.String(), attrVals,
		entity, false, entity.GetCasValue()+1)
	if err != nil {
		errMsg := fmt.Sprintf("Error while updating the IDF entry for WarehouseItem %s: %v", warehouseUuid.String(), err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}
