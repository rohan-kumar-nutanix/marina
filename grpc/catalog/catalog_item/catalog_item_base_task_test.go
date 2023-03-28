/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Author: rishabh.gupta@nutanix.com
 */

package catalog_item

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mockErgon "github.com/nutanix-core/acs-aos-go/ergon/client/mocks"
	mockCpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db/mocks"
	mockSerialExecutor "github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc/serial_executor/mocks"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/acs-aos-go/zeus"

	"github.com/nutanix-core/content-management-marina/interface/external"
	internal "github.com/nutanix-core/content-management-marina/interface/local"
	mockAuthz "github.com/nutanix-core/content-management-marina/mocks/authz"
	mockDb "github.com/nutanix-core/content-management-marina/mocks/db"
	mockCatalogItem "github.com/nutanix-core/content-management-marina/mocks/grpc/catalog/catalog_item"
	mockFile "github.com/nutanix-core/content-management-marina/mocks/grpc/catalog/file_repo"
	mockImage "github.com/nutanix-core/content-management-marina/mocks/grpc/catalog/image"
	mockMetadata "github.com/nutanix-core/content-management-marina/mocks/metadata"
	mockUtils "github.com/nutanix-core/content-management-marina/mocks/util"
	mockZeus "github.com/nutanix-core/content-management-marina/mocks/zeus"
	"github.com/nutanix-core/content-management-marina/task/base"
)

var (
	authzIfc          = new(mockAuthz.AuthzInterface)
	catalogItemIfc    = new(mockCatalogItem.CatalogItemInterface)
	cpdbIfc           = new(mockCpdb.CPDBClientInterface)
	ergonIfc          = new(mockErgon.Ergon)
	fileRepoIfc       = new(mockFile.FileRepoInterface)
	idfIfc            = new(mockDb.IdfClientInterface)
	imageIfc          = new(mockImage.ImageInterface)
	protoIfc          = new(mockUtils.ProtoUtilInterface)
	metadataIfc       = new(mockMetadata.EntityMetadataInterface)
	serialExecutorIfc = new(mockSerialExecutor.SerialExecutorIfc)
	uuidIfc           = new(mockUtils.UuidUtilInterface)
	configIfc         = new(mockZeus.ConfigCache)
	zkSession         = new(zeus.ZookeeperSession)

	testCatalogItemUuid, _ = uuid4.New()
	testCatalogItemVersion = int64(5)
)

func mockExternalInterfaces() external.MarinaExternalInterfaces {
	return external.GetSingletonServiceWithParams(cpdbIfc, ergonIfc, idfIfc, serialExecutorIfc, configIfc, zkSession)
}
func mockInternalInterfaces() internal.MarinaInternalInterfaces {
	return internal.GetSingletonServiceWithParams(authzIfc, metadataIfc, protoIfc, uuidIfc)
}

func TestNewCatalogItemBaseTask(t *testing.T) {
	baseTask := &base.MarinaBaseTask{}
	catalogItemBaseTask := NewCatalogItemBaseTask(baseTask)
	catalogItemBaseTask.catalogItemIfc = catalogItemIfc
	assert.NotNil(t, catalogItemBaseTask)
	assert.NotNil(t, catalogItemBaseTask.catalogItemIfc)
	assert.Equal(t, catalogItemBaseTask.MarinaBaseTask, baseTask)

}
