/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Author: rishabh.gupta@nutanix.com
*
* The implementation for CatalogItemDelete RPC
 */

package tasks

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	. "github.com/nutanix-core/acs-aos-go/insights/insights_interface/query"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/content-management-marina/db"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
	"github.com/nutanix-core/content-management-marina/grpc/catalog/catalog_item"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	utils "github.com/nutanix-core/content-management-marina/util"
	catalogClient "github.com/nutanix-core/content-management-marina/util/catalog/client"
)

var remoteCatalogService = catalogClient.NewRemoteCatalogService

var catalogItemDeleteAttributes = []interface{}{
	"uuid",
}

type CatalogItemDeleteTask struct {
	*CatalogItemBaseTask
	arg *marinaIfc.CatalogItemDeleteArg
}

func NewCatalogItemDeleteTask(catalogItemBaseTask *CatalogItemBaseTask) *CatalogItemDeleteTask {
	return &CatalogItemDeleteTask{
		CatalogItemBaseTask: catalogItemBaseTask,
	}
}

func (task *CatalogItemDeleteTask) StartHook() error {
	arg, err := task.getCatalogItemDeleteArg()
	if err != nil {
		return err
	}
	task.arg = arg

	uuid, err := task.getCatalogItemUuid()
	if err != nil {
		return err
	}
	task.globalCatalogItemUuid = uuid

	wal := task.Wal()
	wal.Data = &marinaIfc.PcTaskWalRecordData{CatalogItem: &marinaIfc.PcTaskWalRecordCatalogItemData{}}
	return task.SetWal(wal)
}

func (task *CatalogItemDeleteTask) RecoverHook() error {
	arg, err := task.getCatalogItemDeleteArg()
	if err != nil {
		return err
	}
	task.arg = arg

	uuid, err := task.getCatalogItemUuid()
	if err != nil {
		return err
	}
	task.globalCatalogItemUuid = uuid
	return nil
}

func (task *CatalogItemDeleteTask) Enqueue() {
	serialExecutor := task.ExternalInterfaces().SerialExecutor()
	serialExecutor.SubmitJob(task)
}

func (task *CatalogItemDeleteTask) Execute() {
	task.Resume(task)
}

func (task *CatalogItemDeleteTask) SerializationID() string {
	return task.globalCatalogItemUuid.String()
}

func (task *CatalogItemDeleteTask) Run() error {
	wal := task.Wal()
	catalogItemWal := wal.GetData().GetCatalogItem()

	var peUuidList []*uuid4.Uuid
	if len(catalogItemWal.GetClusterUuidList()) == 0 {
		// Save the cluster UUIDs we are acting on, in the case that a PE
		// is registered in the middle of a task execution.
		peUuidList = task.ExternalInterfaces().ZeusConfig().PeClusterUuids()
		for _, peUuid := range peUuidList {
			catalogItemWal.ClusterUuidList = append(catalogItemWal.GetClusterUuidList(), peUuid.RawBytes())
		}
		err := task.SetWal(wal)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to set the task WAL: %v", err)
			return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
		}
		task.Save(task)

	} else {
		for _, peUuidBytes := range catalogItemWal.GetClusterUuidList() {
			peUuidList = append(peUuidList, uuid4.ToUuid4(peUuidBytes))
		}
	}

	sourceClusterUuid := task.ExternalInterfaces().ZeusConfig().ClusterUuid()
	var remoteEndpointList []utils.RemoteEndpoint
	for _, peUuid := range peUuidList {
		remoteEndpointList = append(remoteEndpointList, utils.RemoteEndpoint{
			RemoteClusterUuid: *peUuid,
			SourceClusterUuid: *sourceClusterUuid,
		})
	}

	if len(remoteEndpointList) > 0 {
		err := task.fanoutCatalogRequests(remoteEndpointList, sourceClusterUuid)
		if err != nil {
			return err
		}
	}

	err := task.deleteCatalogItem("catalog_item_delete")
	if err != nil {
		return err
	}

	ret := &marinaIfc.CatalogItemDeleteRet{}
	retBytes, err := task.InternalInterfaces().ProtoIfc().Marshal(ret)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to serialize the return object: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	task.Complete(task, retBytes, nil)
	return nil
}

func (task *CatalogItemDeleteTask) getCatalogItemUuid() (*uuid4.Uuid, error) {
	uuid := uuid4.ToUuid4(task.arg.GetCatalogItemId().GetGlobalCatalogItemUuid())
	if uuid == nil {
		errMsg := fmt.Sprintf("Catalog Item UUID missing in the request")
		return nil, marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
	}
	return uuid, nil
}

func (task *CatalogItemDeleteTask) getCatalogItemDeleteArg() (*marinaIfc.CatalogItemDeleteArg, error) {
	embedded := task.Proto().GetRequest().GetArg().GetEmbedded()
	arg := &marinaIfc.CatalogItemDeleteArg{}
	if err := proto.Unmarshal(embedded, arg); err != nil {
		errMsg := fmt.Sprintf("Failed to unmarshal the RPC arguments: %v", err)
		return nil, marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
	}
	return arg, nil
}

func (task *CatalogItemDeleteTask) deleteCatalogItem(queryName string) error {
	version := task.arg.GetCatalogItemId().Version
	where := EQ(COL(catalog_item.GlobalCatalogItemUuid), STR(task.globalCatalogItemUuid.String()))
	if version != nil {
		where = AND(where, EQ(COL(catalog_item.CatalogItemVersion), INT64(*version)))
	}

	query, err := QUERY(queryName).
		FROM(db.CatalogItem.ToString()).
		SELECT(catalogItemDeleteAttributes...).
		WHERE(where).
		Proto()
	if err != nil {
		errMsg := fmt.Sprintf("Error encountered while building the IDF query %s: %v", queryName, err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	cpdbIfc := task.ExternalInterfaces().CPDBIfc()
	arg := &insights_interface.GetEntitiesWithMetricsArg{Query: query}
	entities, err := cpdbIfc.Query(arg)
	if err == insights_interface.ErrNotFound {
		log.Errorf("IDF query did not yield any result: %v", err)
		return nil

	} else if err != nil {
		errMsg := fmt.Sprintf("Failed to fetch the catalog item(s): %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	var catalogItemUuidList []string
	for _, entity := range entities {
		catalogItemUuidList = append(catalogItemUuidList, entity.GetEntityGuid().GetEntityId())
	}

	err = task.ExternalInterfaces().IdfIfc().DeleteEntities(context.Background(), cpdbIfc, db.CatalogItem,
		catalogItemUuidList, true)
	if err == insights_interface.ErrNotFound {
		log.Errorf("Provided catalog item(s) do not exist in IDF: %v", err)
		return nil

	} else if err != nil {
		errMsg := fmt.Sprintf("Failed to delete the catalog item(s): %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}

func (task *CatalogItemDeleteTask) fanoutCatalogRequests(remoteEndpointList []utils.RemoteEndpoint,
	sourceClusterUuid *uuid4.Uuid) error {
	taskUuidByEndpoint := make(map[utils.RemoteEndpoint]*uuid4.Uuid)
	wal := task.Wal()
	taskList := wal.GetData().GetTaskList()
	for _, taskObj := range taskList {
		remoteEndpoint := utils.RemoteEndpoint{SourceClusterUuid: *sourceClusterUuid}
		if clusterUuid := taskObj.GetClusterUuid(); clusterUuid != nil {
			remoteEndpoint.RemoteClusterUuid = *uuid4.ToUuid4(clusterUuid)
		}
		taskUuidByEndpoint[remoteEndpoint] = uuid4.ToUuid4(taskObj.GetTaskUuid())
	}

	successTaskUuidByEndpoint := make(map[utils.RemoteEndpoint]*uuid4.Uuid)
	for _, remoteEndpoint := range remoteEndpointList {
		var remoteTaskUuid *uuid4.Uuid
		var err error
		if taskUuid, ok := taskUuidByEndpoint[remoteEndpoint]; ok {
			remoteTaskUuid = taskUuid

		} else {
			remoteTaskUuid, err = task.InternalInterfaces().UuidIfc().New()
			if err != nil {
				return err
			}

			endpointTask := &marinaIfc.EndpointTask{TaskUuid: remoteTaskUuid.RawBytes()}
			if clusterUuid := remoteEndpoint.RemoteClusterUuid; clusterUuid != utils.NilUuid {
				endpointTask.ClusterUuid = clusterUuid.RawBytes()
			}

			wal.GetData().TaskList = append(wal.GetData().TaskList, endpointTask)
			err = task.SetWal(wal)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to set the task WAL: %v", err)
				return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
			}
			task.Save(task)
			taskUuidByEndpoint[remoteEndpoint] = remoteTaskUuid
		}

		fanoutArg := proto.Clone(task.arg).(*marinaIfc.CatalogItemDeleteArg)
		fanoutArg.TaskUuid = remoteTaskUuid.RawBytes()
		fanoutArg.ParentTaskUuid = task.Proto().Uuid
		fanoutRet := &marinaIfc.CatalogItemDeleteRet{}
		client := remoteCatalogService(task.ExternalInterfaces().ZkSession(), task.ExternalInterfaces().ZeusConfig(),
			&remoteEndpoint.RemoteClusterUuid, nil, nil, nil)
		err = client.SendMsg("CatalogItemDelete", fanoutArg, fanoutRet)
		if err != nil {
			endpointName := remoteEndpoint.GetUserVisibleId(task.ExternalInterfaces().ZeusConfig())

			if endpointName == nil {
				log.Errorf("Failed to fanout CatalogItemDelete request (task %s): %s",
					remoteTaskUuid.String(), err)
			} else {
				log.Errorf("Failed to fanout CatalogItemDelete request (task %s) to %s : %s",
					remoteTaskUuid.String(), *endpointName, err)
			}

		} else {
			successTaskUuidByEndpoint[remoteEndpoint] = remoteTaskUuid
		}
	}

	if len(successTaskUuidByEndpoint) <= 0 {
		errMsg := fmt.Sprintf("Failed to fanout CatalogItemDelete request to all remote clusters")
		return marinaError.ErrCatalogTaskForwardError.SetCauseAndLog(errors.New(errMsg))
	}

	task.InternalInterfaces().FanoutTaskPollerIfc().PollAllRemoteTasks(task.ExternalInterfaces().ZkSession(),
		&successTaskUuidByEndpoint)
	return nil
}
