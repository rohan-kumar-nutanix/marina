/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Author: rishabh.gupta@nutanix.com
*
* The implementation for CatalogItem Base Task
 */

package catalog_item

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/ergon"
	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	"github.com/nutanix-core/content-management-marina/db"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
	"github.com/nutanix-core/content-management-marina/grpc/catalog/file_repo"
	"github.com/nutanix-core/content-management-marina/grpc/catalog/image"

	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	"github.com/nutanix-core/content-management-marina/task/base"
	catalogClient "github.com/nutanix-core/content-management-marina/util/catalog/client"
)

var remoteCatalogService = catalogClient.NewRemoteCatalogService

type CatalogItemBaseTask struct {
	*base.MarinaBaseTask
	globalCatalogItemUuid *uuid4.Uuid
	taskUuid              *uuid4.Uuid
	catalogItemIfc        CatalogItemInterface
	// TODO: Remove fileRepoIfc & imageIfc once File Entity is added.
	fileRepoIfc file_repo.FileRepoInterface
	imageIfc    image.ImageInterface
}

func NewCatalogItemBaseTask(marinaBaseTask *base.MarinaBaseTask) *CatalogItemBaseTask {
	return &CatalogItemBaseTask{
		MarinaBaseTask: marinaBaseTask,
		catalogItemIfc: newCatalogItemImpl(),
		fileRepoIfc:    file_repo.NewFileRepoImpl(),
		imageIfc:       image.NewImageInterface(),
	}
}

func (task *CatalogItemBaseTask) addCatalogItemEntity(uuid *uuid4.Uuid) {
	entity := &ergon.EntityId{
		EntityType: ergon.EntityId_kCatalogItem.Enum(),
		EntityId:   uuid.RawBytes(),
	}
	task.AddEntity(task.Proto(), entity)
}

// getClusterFileUuidMap returns a map with following mappings [pcFileUuid][peUuid] --> peFileUuid
func (task *CatalogItemBaseTask) getClusterFileUuidMap(sourceGroupSpecs []*marinaIfc.SourceGroupSpec) (
	map[uuid4.Uuid]map[uuid4.Uuid]uuid4.Uuid, error) {

	// Find the cluster corresponding file uuid for each file. This is used to replace the file uuid in local
	// import part of the spec specific to each cluster where the request will be fanned out.
	clusterFileUuidMap := make(map[uuid4.Uuid]map[uuid4.Uuid]uuid4.Uuid)
	for _, groupSpec := range sourceGroupSpecs {
		for _, sourceSpec := range groupSpec.SourceSpecList {
			if sourceSpec.GetImportSpec() != nil {
				for _, localImport := range sourceSpec.GetImportSpec().LocalImportList {
					if localImport.GetFileUuid() == nil {
						// This is when local import points to a ADSF file on any PE.
						continue
					}

					fileUuid := uuid4.ToUuid4(localImport.GetFileUuid())
					clusterFileUuidMap[*fileUuid] = make(map[uuid4.Uuid]uuid4.Uuid)
					file, err := task.fileRepoIfc.GetFile(context.TODO(), task.ExternalInterfaces().CPDBIfc(), fileUuid)
					if err != nil {
						return nil, err
					}

					for _, location := range file.LocationList {
						if location.ClusterUuid != nil {
							clusterUuid := uuid4.ToUuid4(location.ClusterUuid)
							clusterFileUuid := uuid4.ToUuid4(location.FileUuid)
							clusterFileUuidMap[*fileUuid][*clusterUuid] = *clusterFileUuid
						}
					}
				}
			}
		}
	}
	return clusterFileUuidMap, nil
}

func (task *CatalogItemBaseTask) adjustSourceGroupSpec(sourceGroupSpec *marinaIfc.SourceGroupSpec,
	fileUuidMap map[uuid4.Uuid]map[uuid4.Uuid]uuid4.Uuid, containerUuid *uuid4.Uuid, clusterUuid uuid4.Uuid) {

	for _, sourceSpec := range sourceGroupSpec.SourceSpecList {
		// Iterate over remote import list. If cluster specific container uuid is not available we
		// delete the entire remote_import_list.
		if containerUuid == nil {
			sourceSpec.GetImportSpec().RemoteImportList = nil

		} else {
			for _, remoteImport := range sourceSpec.GetImportSpec().RemoteImportList {
				remoteImport.ContainerUuid = containerUuid.RawBytes()
			}
		}

		// Iterate over local import list. If the file is not present on the remote endpoint cluster then the
		// respective file uuid will not be present, and thus we delete the local import from the cluster specific
		// arg's source_spec's local_import_list. Currently, we can support only one local_import where we import
		// from another file.
		var localImports []*marinaIfc.LocalImportSpec
		for _, localImport := range sourceSpec.GetImportSpec().LocalImportList {
			if localImport.FileUuid != nil {
				fileUuid := uuid4.ToUuid4(localImport.FileUuid)
				if clusterFileUuid, ok := fileUuidMap[*fileUuid][clusterUuid]; ok {
					sourceSpec.GetImportSpec().FileUuid = clusterFileUuid.RawBytes()
					localImport.ClusterUuid = nil
					localImport.FileUuid = clusterFileUuid.RawBytes()
					localImport.ContainerUuid = containerUuid.RawBytes()
					localImports = append(localImports, localImport)
				}

			} else if localImport.GetAdsfPath() != "" && localImport.ClusterUuid != nil {
				localImportClusterUuid := uuid4.ToUuid4(localImport.ClusterUuid)
				if clusterUuid.Equals(localImportClusterUuid) {
					localImports = append(localImports, localImport)
				}

			} else if localImport.GetAdsfPath() == "" {
				localImports = append(localImports, localImport)
			}
		}

		sourceSpec.GetImportSpec().LocalImportList = localImports
	}
}

func (task *CatalogItemBaseTask) generateSourceGroupSpecs(sourceGroupSpecs []*marinaIfc.SourceGroupSpec,
	catalogItemByCluster map[uuid4.Uuid]*marinaIfc.CatalogItemInfo, isContentAddressable bool) []*marinaIfc.SourceGroupSpec {

	newSourceGroupSpecByUuid := make(map[uuid4.Uuid]*marinaIfc.SourceGroupSpec)
	for _, sourceGroupSpec := range sourceGroupSpecs {
		sourceGroupUuid := uuid4.ToUuid4(sourceGroupSpec.GetUuid())
		newSourceGroupSpecByUuid[*sourceGroupUuid] = &marinaIfc.SourceGroupSpec{Uuid: sourceGroupUuid.RawBytes()}
	}

	for clusterUuid, catalogItem := range catalogItemByCluster {
		for _, sourceGroup := range catalogItem.SourceGroupList {
			sourceGroupUuid := uuid4.ToUuid4(sourceGroup.GetUuid())
			newSourceGroupSpec, ok := newSourceGroupSpecByUuid[*sourceGroupUuid]
			if !ok {
				if len(sourceGroupSpecs) == 1 {
					sourceGroupUuid = uuid4.ToUuid4(sourceGroupSpecs[0].GetUuid())
					newSourceGroupSpec = newSourceGroupSpecByUuid[*sourceGroupUuid]

				} else {
					continue
				}
			}

			for i, source := range sourceGroup.SourceList {
				var newSourceSpec *marinaIfc.SourceSpec
				if i <= len(newSourceGroupSpec.SourceSpecList)-1 {
					newSourceSpec = newSourceGroupSpec.SourceSpecList[i]
				} else {
					newSourceSpec = &marinaIfc.SourceSpec{}
					newSourceGroupSpec.SourceSpecList = append(newSourceGroupSpec.SourceSpecList, newSourceSpec)
				}

				if newSourceSpec.GetImportSpec() == nil {
					newSourceSpec.ImportSpec = &marinaIfc.FileImportSpec{}
				}
				importSpec := newSourceSpec.GetImportSpec()
				log.Infof("[%s] Marking the PC FileImport task Content Addressability as %v",
					task.taskUuid.String(), isContentAddressable)
				importSpec.IsRequestContentAddressable = proto.Bool(isContentAddressable)
				localImport := &marinaIfc.LocalImportSpec{
					ClusterUuid: uuid4.ToUuid4(clusterUuid.RawBytes()).RawBytes(),
					FileUuid:    source.FileUuid,
				}
				importSpec.LocalImportList = append(importSpec.LocalImportList, localImport)
			}
		}
	}

	var newSourceGroupSpecs []*marinaIfc.SourceGroupSpec
	for _, sourceGroupSpec := range newSourceGroupSpecByUuid {
		newSourceGroupSpecs = append(newSourceGroupSpecs, sourceGroupSpec)
	}

	return newSourceGroupSpecs
}

func (task *CatalogItemBaseTask) addFileToLocalRepo(importSpec *marinaIfc.FileImportSpec) (*uuid4.Uuid, error) {
	arg := &marinaIfc.FileImportArg{}
	arg.Spec = proto.Clone(importSpec).(*marinaIfc.FileImportSpec)
	uuid, err := task.InternalInterfaces().UuidIfc().New()
	if err != nil {
		return nil, err
	}
	arg.TaskUuid = uuid.RawBytes()
	arg.ParentTaskUuid = task.Proto().Uuid
	ret := &marinaIfc.FileImportRet{}
	// TODO: Once FileImport is implemented in Marina, directly create the task object
	client := catalogClient.DefaultCatalogService()
	err = client.SendMsg("FileImport", arg, ret)
	if err != nil {
		return nil, err
	}
	return uuid4.ToUuid4(ret.TaskUuid), nil
}

func (task *CatalogItemBaseTask) cleanupFileRepoEntries() error {
	var fileUuids []string
	for _, fileUuid := range task.Wal().GetData().GetCatalogItem().FileUuidList {
		fileUuids = append(fileUuids, uuid4.ToUuid4(fileUuid).String())
	}

	err := task.ExternalInterfaces().IdfIfc().DeleteEntities(context.TODO(), task.ExternalInterfaces().CPDBIfc(),
		db.File, fileUuids, true)
	if err == insights_interface.ErrNotFound {
		log.Errorf("[%s] Provided file(s) do not exist in IDF: %v", task.taskUuid.String(), err)
		return nil

	} else if err != nil {
		errMsg := fmt.Sprintf("Failed to delete the file(s): %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	return nil
}
