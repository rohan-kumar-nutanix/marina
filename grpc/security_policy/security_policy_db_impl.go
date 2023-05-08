/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Authors: rajesh.battala@nutanix.com
 *
 */

package security_policy

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
	utils "github.com/nutanix-core/content-management-marina/util"
)

var (
	securityPolicyDBImpl ISecurityPolicyDB = nil
	once                 sync.Once
)

type SecurityPolicyDBImpl struct {
}

func newSecurityPolicyDBImpl() ISecurityPolicyDB {
	once.Do(func() {
		securityPolicyDBImpl = &SecurityPolicyDBImpl{}
	})
	return securityPolicyDBImpl
}
func (warehouseDBImpl *SecurityPolicyDBImpl) CreateSecurityPolicy(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	protoIfc utils.ProtoUtilInterface, securityPolicyUuid *uuid4.Uuid, securityPolicyPB *config.SecurityPolicy) error {
	log.Infof("Persisting Security Policy Model %v", securityPolicyPB)
	marshal, err := protoIfc.Marshal(securityPolicyPB)
	if err != nil {
		errMsg := fmt.Sprintf("Error marshaling the Security Policy proto: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	var buffer bytes.Buffer
	writer := zlib.NewWriter(&buffer)
	_, err = writer.Write(marshal)
	if err != nil {
		errMsg := fmt.Sprintf("Error compressing the Security Policy proto: %v", err)
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

	_, err = cpdbIfc.UpdateEntity(db.SecurityPolicy.ToString(), securityPolicyUuid.String(), attrVals, nil, false, 0)
	if err != nil {
		errMsg := fmt.Sprintf("Error while creating the IDF entry for Security Policy %s: %v", securityPolicyUuid.String(), err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}

func (warehouseDBImpl *SecurityPolicyDBImpl) GetSecurityPolicy(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	securityPolicyUuid *uuid4.Uuid) (*config.SecurityPolicy, error) {
	entity, err := cpdbIfc.GetEntity(db.SecurityPolicy.ToString(), securityPolicyUuid.String(), false)
	if err == insights_interface.ErrNotFound || entity == nil {
		log.Errorf("Security Policy %s not found", securityPolicyUuid.String())
		return nil, marinaError.ErrNotFound
	} else if err != nil {
		errMsg := fmt.Sprintf("Error encountered while getting Security Policy %s: %v", securityPolicyUuid.String(), err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	securityPolicy := &config.SecurityPolicy{}
	err = entity.DeserializeEntity(securityPolicy)
	if err != nil {
		errMsg := fmt.Sprintf("failed to deserialize Security Policy IDF entry: %v", err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return securityPolicy, nil
}

func (warehouseDBImpl *SecurityPolicyDBImpl) DeleteSecurityPolicy(ctx context.Context, idfIfc db.IdfClientInterface,
	cpdbIfc cpdb.CPDBClientInterface, securityPolicyUuid string) error {
	err := idfIfc.DeleteEntities(context.Background(), cpdbIfc, db.SecurityPolicy, []string{securityPolicyUuid}, true)
	if err == insights_interface.ErrNotFound {
		log.Errorf("Security Policy UUID :%v do not exist in IDF", err)
		return nil

	} else if err != nil {
		errMsg := fmt.Sprintf("Failed to delete the Security Policy: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}

func (warehouseDBImpl *SecurityPolicyDBImpl) UpdateSecurityPolicy(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface,
	protoIfc utils.ProtoUtilInterface, securityPolicyUuid *uuid4.Uuid, securityPolicyPB *config.SecurityPolicy) error {

	entity, err := cpdbIfc.GetEntity(db.SecurityPolicy.ToString(), securityPolicyUuid.String(), true)
	if err == insights_interface.ErrNotFound || entity == nil {
		log.Errorf("Security Policy %s not found", securityPolicyUuid.String())
		return marinaError.ErrNotFound
	} else if err != nil {
		errMsg := fmt.Sprintf("Error encountered while getting Security Policy %s: %v", securityPolicyUuid.String(), err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	marshal, err := protoIfc.Marshal(securityPolicyPB)
	if err != nil {
		errMsg := fmt.Sprintf("Error marshaling the Security Policy proto: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	var buffer bytes.Buffer
	writer := zlib.NewWriter(&buffer)
	_, err = writer.Write(marshal)
	if err != nil {
		errMsg := fmt.Sprintf("Error compressing the Security Policy proto: %v", err)
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

	_, err = cpdbIfc.UpdateEntity(db.SecurityPolicy.ToString(), securityPolicyUuid.String(), attrVals,
		entity, false, entity.GetCasValue()+1)
	if err != nil {
		errMsg := fmt.Sprintf("Error while creating the IDF entry for Security Policy %s: %v", securityPolicyUuid.String(), err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}

func (warehouseDBImpl *SecurityPolicyDBImpl) ListSecurityPolices(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface) (
	[]*config.SecurityPolicy, error) {
	var securityPolicies []*config.SecurityPolicy
	entities, err := cpdbIfc.GetEntitiesOfType(db.SecurityPolicy.ToString(), false)
	if err != nil {
		return nil, err
	}
	for _, entity := range entities {
		securityPolicy := &config.SecurityPolicy{}
		err := entity.DeserializeEntity(securityPolicy)
		if err != nil {
			log.Errorf("error occurred in deserializing the Security Policy Entity and skipping it. %v", err)
		}
		securityPolicies = append(securityPolicies, securityPolicy)
	}
	return securityPolicies, nil
}
