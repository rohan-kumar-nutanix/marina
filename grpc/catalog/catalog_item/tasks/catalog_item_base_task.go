/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Author: rishabh.gupta@nutanix.com
*
* The implementation for CatalogItem Base Task
 */

package tasks

import (
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/content-management-marina/task/base"
)

type CatalogItemBaseTask struct {
	*base.MarinaBaseTask
	globalCatalogItemUuid *uuid4.Uuid
}

func NewCatalogItemBaseTask(marinaBaseTask *base.MarinaBaseTask) *CatalogItemBaseTask {
	return &CatalogItemBaseTask{
		MarinaBaseTask: marinaBaseTask,
	}
}
