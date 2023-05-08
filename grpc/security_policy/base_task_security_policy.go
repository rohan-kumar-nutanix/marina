/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Authors: rajesh.battala@nutanix.com
 *
 */

package security_policy

import (
	"github.com/nutanix-core/acs-aos-go/ergon"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	"github.com/nutanix-core/content-management-marina/task/base"
)

type MarinaBaseWarehouseTaskInterface interface {
	AddAnEntity(entityType *ergon.EntityId_Entity, entityUuid *uuid4.Uuid)
	AddSecurityPolicyEntity()
	GetSecurityPolicyUuid() *uuid4.Uuid
	SetSecurityPolicyUuid(uuid *uuid4.Uuid)
}

type MarinaBaseSecurityPolicyTask struct {
	*base.MarinaBaseTask
	SecurityPolicyUuid *uuid4.Uuid
	ISecurityPolicyDB
}

func NewMarinaBaseSecurityPolicyTask(
	marinaBaseTask *base.MarinaBaseTask) *MarinaBaseSecurityPolicyTask {
	return &MarinaBaseSecurityPolicyTask{
		MarinaBaseTask:    marinaBaseTask,
		ISecurityPolicyDB: newSecurityPolicyDBImpl(),
	}
}

func (t *MarinaBaseSecurityPolicyTask) GetSecurityPolicyUuid() *uuid4.Uuid {
	return t.SecurityPolicyUuid
}

func (t *MarinaBaseSecurityPolicyTask) AddAnEntity(entityType *ergon.EntityId_Entity, entityUuid *uuid4.Uuid) {
	// TODO: Add SecurityPolicy Entity in Ergon.
	entity := &ergon.EntityId{
		EntityType: ergon.EntityId_kCatalogItem.Enum(),
		EntityId:   entityUuid.RawBytes(),
	}
	t.AddEntity(t.Proto(), entity)
}

func (t *MarinaBaseSecurityPolicyTask) AddSecurityPolicyEntity() {
	// TODO: Add SecurityPolicy to Ergon Types.
	t.AddAnEntity(ergon.EntityId_kUnknown.Enum(), t.GetSecurityPolicyUuid())
}

func (t *MarinaBaseSecurityPolicyTask) SetSecurityPolicyUuid(uuid *uuid4.Uuid) {
	t.SecurityPolicyUuid = uuid
}

func (t *MarinaBaseSecurityPolicyTask) RecoverHook() error {
	return nil
}
