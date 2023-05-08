/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Authors: rajesh.battala@nutanix.com
 *
 */

package scanner_config

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
	"github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/config"
	"github.com/nutanix-core/content-management-marina/protos/apis/common/v1/response"
	utils "github.com/nutanix-core/content-management-marina/util"
)

var (
	scannerConfigDBImpl IScannerConfigDB = nil
	once                sync.Once
)

type ScannerConfigDBImpl struct {
}

func newScannerConfigDBImpl() IScannerConfigDB {
	once.Do(func() {
		scannerConfigDBImpl = &ScannerConfigDBImpl{}
	})
	return scannerConfigDBImpl
}
func (impl *ScannerConfigDBImpl) CreateScannerConfig(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	protoIfc utils.ProtoUtilInterface, scannerConfigUuid *uuid4.Uuid, scannerConfigPB *config.ScannerConfig) error {
	log.Infof("Persisting Scanner Policy Model %v", scannerConfigPB)
	marshal, err := protoIfc.Marshal(scannerConfigPB)
	if err != nil {
		errMsg := fmt.Sprintf("Error marshaling the Scanner Config proto: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	var buffer bytes.Buffer
	writer := zlib.NewWriter(&buffer)
	_, err = writer.Write(marshal)
	if err != nil {
		errMsg := fmt.Sprintf("Error compressing the Scanner Config proto: %v", err)
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

	_, err = cpdbIfc.UpdateEntity(db.ScannerConfig.ToString(), scannerConfigUuid.String(), attrVals, nil, false, 0)
	if err != nil {
		errMsg := fmt.Sprintf("Error while creating the IDF entry for Scanner Config %s: %v", scannerConfigUuid.String(), err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}

func (impl *ScannerConfigDBImpl) GetScannerConfig(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	scannerConfigUuid *uuid4.Uuid) (*config.ScannerConfig, error) {
	entity, err := cpdbIfc.GetEntity(db.ScannerConfig.ToString(), scannerConfigUuid.String(), false)
	if err == insights_interface.ErrNotFound || entity == nil {
		log.Errorf("Scanner Config %s not found", scannerConfigUuid.String())
		return nil, marinaError.ErrNotFound
	} else if err != nil {
		errMsg := fmt.Sprintf("Error encountered while getting Scanner Config %s: %v", scannerConfigUuid.String(), err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	scannerConfig := &config.ScannerConfig{}
	err = entity.DeserializeEntity(scannerConfig)
	if err != nil {
		errMsg := fmt.Sprintf("failed to deserialize Scanner Config IDF entry: %v", err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return scannerConfig, nil
}

func (impl *ScannerConfigDBImpl) DeleteScannerConfig(ctx context.Context, idfIfc db.IdfClientInterface,
	cpdbIfc cpdb.CPDBClientInterface, scannerConfigUuid string) error {
	err := idfIfc.DeleteEntities(context.Background(), cpdbIfc, db.ScannerConfig, []string{scannerConfigUuid}, true)
	if err == insights_interface.ErrNotFound {
		log.Errorf("Scanner Config UUID :%v do not exist in IDF", err)
		return nil

	} else if err != nil {
		errMsg := fmt.Sprintf("Failed to delete the Scanner Config: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}

func (impl *ScannerConfigDBImpl) UpdateScannerConfig(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	protoIfc utils.ProtoUtilInterface, scannerConfigUuid *uuid4.Uuid, scannerConfigPB *config.ScannerConfig) error {

	entity, err := cpdbIfc.GetEntity(db.ScannerConfig.ToString(), scannerConfigUuid.String(), true)
	if err == insights_interface.ErrNotFound || entity == nil {
		log.Errorf("Scanner Config %s not found", scannerConfigUuid.String())
		return marinaError.ErrNotFound
	} else if err != nil {
		errMsg := fmt.Sprintf("Error encountered while getting Scanner Config %s: %v", scannerConfigUuid.String(), err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	marshal, err := protoIfc.Marshal(scannerConfigPB)
	if err != nil {
		errMsg := fmt.Sprintf("Error marshaling the Scanner Config proto: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	var buffer bytes.Buffer
	writer := zlib.NewWriter(&buffer)
	_, err = writer.Write(marshal)
	if err != nil {
		errMsg := fmt.Sprintf("Error compressing the Scanner Config proto: %v", err)
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

	_, err = cpdbIfc.UpdateEntity(db.ScannerConfig.ToString(), scannerConfigUuid.String(), attrVals,
		entity, false, entity.GetCasValue()+1)
	if err != nil {
		errMsg := fmt.Sprintf("Error while creating the IDF entry for Scanner Config %s: %v", scannerConfigUuid.String(), err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}

func (impl *ScannerConfigDBImpl) ListScannerConfigurations(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface) (
	[]*config.ScannerConfig, error) {
	var scannerConfigs []*config.ScannerConfig
	entities, err := cpdbIfc.GetEntitiesOfType(db.ScannerConfig.ToString(), false)
	if err != nil {
		log.Errorf("error occurred while fetching data from IDF %s", err)
		return nil, err
	}
	for _, entity := range entities {
		scannerConfig := &config.ScannerConfig{}
		err := entity.DeserializeEntity(scannerConfig)
		if err != nil {
			log.Errorf("error occurred in deserializing the Scanner Config Entity and skipping it. %v", err)
		}

		// scannerConfig.Base.ExtId = entity.EntityGuid.EntityId
		scannerConfig.Base = &response.ExternalizableAbstractModel{ExtId: entity.EntityGuid.EntityId}
		scannerConfigs = append(scannerConfigs, scannerConfig)
	}
	return scannerConfigs, nil
}
