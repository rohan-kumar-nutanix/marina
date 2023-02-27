/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Author: rishabh.gupta@nutanix.com
*
* The implementation for CatalogItemCreate RPC
 */

package tasks

import (
	"errors"
	"fmt"

	set "github.com/deckarep/golang-set/v2"
	"github.com/golang/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/ergon"
	taskUtil "github.com/nutanix-core/acs-aos-go/ergon/task"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	utils "github.com/nutanix-core/content-management-marina/util"
)

type CatalogItemCreateTask struct {
	*CatalogItemBaseTask
	arg                   *marinaIfc.CatalogItemCreateArg
	spec                  *marinaIfc.CatalogItemCreateSpec
	catalogItemUuid       *uuid4.Uuid
	version               int64
	errReasons            []string
	fanoutSuccessClusters []*uuid4.Uuid
}

func NewCatalogItemCreateTask(catalogItemBaseTask *CatalogItemBaseTask) *CatalogItemCreateTask {
	return &CatalogItemCreateTask{
		CatalogItemBaseTask: catalogItemBaseTask,
	}
}

func (task *CatalogItemCreateTask) StartHook() error {
	arg, err := task.getCatalogItemCreateArg()
	if err != nil {
		return err
	}
	task.arg = arg

	spec := arg.GetSpec()
	if spec == nil {
		errMsg := fmt.Sprintf("Spec is missing from the arguments")
		return marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
	}
	task.spec = spec

	if catalogItemUuid := spec.GetUuid(); catalogItemUuid != nil {
		task.catalogItemUuid = uuid4.ToUuid4(catalogItemUuid)

	} else {
		uuid, err := task.InternalInterfaces().UuidIfc().New()
		if err != nil {
			return err
		}
		task.catalogItemUuid = uuid
	}

	if globalCatalogItemUuid := spec.GetGlobalCatalogItemUuid(); globalCatalogItemUuid != nil {
		task.globalCatalogItemUuid = uuid4.ToUuid4(globalCatalogItemUuid)

	} else {
		uuid, err := task.InternalInterfaces().UuidIfc().New()
		if err != nil {
			return err
		}
		task.globalCatalogItemUuid = uuid
	}

	if len(spec.GetSourceGroupSpecList()) > 1 {
		errMsg := fmt.Sprintf("Multiple sources are not supported for creating catalog item")
		return marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
	}

	task.version = spec.GetVersion()
	if spec.GetOwnerClusterUuid() == nil {
		spec.OwnerClusterUuid = task.ExternalInterfaces().ZeusConfig().ClusterUuid().RawBytes()
	}

	task.addCatalogItemEntity(task.catalogItemUuid)

	wal := task.Wal()
	wal.Data = &marinaIfc.PcTaskWalRecordData{
		CatalogItem: &marinaIfc.PcTaskWalRecordCatalogItemData{
			CatalogItemUuid:       task.catalogItemUuid.RawBytes(),
			GlobalCatalogItemUuid: task.globalCatalogItemUuid.RawBytes(),
		},
	}
	return task.SetWal(wal)
}

func (task *CatalogItemCreateTask) RecoverHook() error {
	arg, err := task.getCatalogItemCreateArg()
	if err != nil {
		return err
	}
	task.arg = arg

	spec := arg.GetSpec()
	if spec == nil {
		errMsg := fmt.Sprintf("Spec is missing from the arguments")
		return marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
	}
	task.spec = spec

	wal := task.Wal().GetData().GetCatalogItem()
	task.catalogItemUuid = uuid4.ToUuid4(wal.GetCatalogItemUuid())
	task.globalCatalogItemUuid = uuid4.ToUuid4(wal.GetGlobalCatalogItemUuid())
	task.spec.GlobalCatalogItemUuid = wal.GetGlobalCatalogItemUuid()
	task.version = task.spec.GetVersion()
	if task.spec.GetOwnerClusterUuid() == nil {
		task.spec.OwnerClusterUuid = task.ExternalInterfaces().ZeusConfig().ClusterUuid().RawBytes()
	}
	return nil
}

func (task *CatalogItemCreateTask) Enqueue() {
	serialExecutor := task.ExternalInterfaces().SerialExecutor()
	serialExecutor.SubmitJob(task)
}

func (task *CatalogItemCreateTask) Execute() {
	task.Resume(task)
}

func (task *CatalogItemCreateTask) SerializationID() string {
	return task.globalCatalogItemUuid.String()
}

func (task *CatalogItemCreateTask) Run() error {
	task.taskUuid = uuid4.ToUuid4(task.Proto().Uuid)

	// Assign UUIDs to source specs
	err := task.assignSourceGroupSpecUuid()
	if err != nil {
		return err
	}

	// Create the list of target clusters for creating catalog item
	targetClusters, err := task.targetClusterList()
	if err != nil {
		return err
	}

	// Create a map with CatalogItemCreate arg for each target cluster
	argByCluster, err := task.prepareFanoutArgs(targetClusters)
	if err != nil {
		return err
	}

	// Save the cluster UUIDs to WAL where we want to fanout the request
	walClusters, err := task.walTargetClusters(argByCluster)
	if err != nil {
		return err
	}

	// Check if the catalog item already exists
	ret := &marinaIfc.CatalogItemCreateTaskRet{}
	exists, err := task.checkCatalogItemExists(ret)
	if err != nil {
		return err
	}

	retBytes, err := task.InternalInterfaces().ProtoIfc().Marshal(ret)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to serialize the return object: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	if exists {
		task.Complete(task, retBytes, nil)
	}

	// Fanout the request to target clusters
	task.fanoutSuccessClusters = walClusters
	taskByCluster, err := task.fanoutCreateRequests(walClusters, argByCluster)
	if err != nil {
		task.errReasons = append(task.errReasons, fmt.Sprintf("Catalog item create failed: %v", err))
		err = task.handleCleanup(task.InternalInterfaces().UuidIfc().UuidPointersToUuids(task.fanoutSuccessClusters))
		if err != nil {
			return err
		}
	}

	// Get the clusters where the fanout task succeeded
	var finalClusters []uuid4.Uuid
	finalClusters, err = task.processFanoutResult(taskByCluster)
	if err != nil {
		task.errReasons = append(task.errReasons, fmt.Sprintf("Catalog item create failed: %v", err))
		err = task.handleCleanup(task.InternalInterfaces().UuidIfc().UuidPointersToUuids(task.fanoutSuccessClusters))
		if err != nil {
			return err
		}
	}

	// Create PC file repo entries
	sourceGroups, err := task.createFileRepoEntries()
	if err != nil {
		task.errReasons = append(task.errReasons, fmt.Sprintf("Catalog item create failed: %v", err))
		err = task.handleCleanup(finalClusters)
		if err != nil {
			return err
		}
	}

	// Create PC CatalogItem IDF entry
	err = task.InternalInterfaces().CatalogItemIfc().CreateCatalogItemFromCreateSpec(nil,
		task.ExternalInterfaces().CPDBIfc(), task.InternalInterfaces().ProtoIfc(), task.catalogItemUuid,
		task.spec, finalClusters, sourceGroups)
	if err != nil {
		task.errReasons = append(task.errReasons, fmt.Sprintf("Catalog item create failed: %v", err))
		err = task.handleCleanup(finalClusters)
		if err != nil {
			return err
		}
	}

	task.Complete(task, retBytes, nil)
	return nil
}

func (task *CatalogItemCreateTask) getCatalogItemCreateArg() (*marinaIfc.CatalogItemCreateArg, error) {
	embedded := task.Proto().GetRequest().GetArg().GetEmbedded()
	arg := &marinaIfc.CatalogItemCreateArg{}
	if err := proto.Unmarshal(embedded, arg); err != nil {
		errMsg := fmt.Sprintf("Failed to unmarshal the RPC arguments: %v", err)
		return nil, marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
	}
	return arg, nil
}

func (task *CatalogItemCreateTask) assignSourceGroupSpecUuid() error {
	wal := task.Wal()
	catalogWal := wal.GetData().GetCatalogItem()
	if len(catalogWal.AddSourceGroupList) == 0 {
		for _, spec := range task.spec.SourceGroupSpecList {
			if spec.Uuid == nil {
				specUuid, err := task.InternalInterfaces().UuidIfc().New()
				if err != nil {
					return err
				}
				spec.Uuid = specUuid.RawBytes()
			}

			catalogWal.AddSourceGroupList = append(catalogWal.AddSourceGroupList, spec)
		}

		err := task.SetWal(wal)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to set the task WAL: %v", err)
			return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
		}
		task.Save(task)

	} else {
		task.spec.SourceGroupSpecList = nil
		for _, spec := range catalogWal.GetAddSourceGroupList() {
			task.spec.SourceGroupSpecList = append(task.spec.SourceGroupSpecList, spec)
		}
	}
	return nil
}

func (task *CatalogItemCreateTask) targetClusterList() ([]uuid4.Uuid, error) {
	var finalClusters []uuid4.Uuid
	if len(task.spec.InitialClusterUuidList) == 0 {
		for _, uuid := range task.ExternalInterfaces().ZeusConfig().PeClusterUuids() {
			finalClusters = append(finalClusters, *uuid)
		}

	} else {
		remoteClusters := set.NewSet[uuid4.Uuid]()
		initialClusters := set.NewSet[uuid4.Uuid]()
		for _, uuid := range task.ExternalInterfaces().ZeusConfig().PeClusterUuids() {
			remoteClusters.Add(*uuid)
		}
		for _, uuidBytes := range task.spec.InitialClusterUuidList {
			initialClusters.Add(*uuid4.ToUuid4(uuidBytes))
		}

		finalClusters = remoteClusters.Intersect(initialClusters).ToSlice()

		if finalClusters == nil {
			initialClustersStr := task.InternalInterfaces().UuidIfc().StringFromUuids(initialClusters.ToSlice())
			errMsg := fmt.Sprintf("Invalid clusters uuids provided in initial_cluster_uuid_list %s", initialClustersStr)
			return nil, marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
		}

		clustersStr := task.InternalInterfaces().UuidIfc().StringFromUuids(finalClusters)
		log.Infof("[%s] Creating catalog item %s on remote endpoints: %s.", task.taskUuid.String(),
			task.globalCatalogItemUuid.String(), clustersStr)
	}

	return finalClusters, nil
}

func (task *CatalogItemCreateTask) prepareFanoutArgs(targetClusters []uuid4.Uuid) (
	map[uuid4.Uuid]*marinaIfc.CatalogItemCreateArg, error) {

	// It's possible that the creation arg has either remote or local imports or both. For remote source we need
	// to modify the arg sent to each PE to point to the correct container UUID. We don't need to so for local
	// source where the source of the file is an ADSF file or another catalog file object, in such cases we expect
	// the arg to already have correct container info.
	containerByCluster := task.ExternalInterfaces().ZeusConfig().ClusterSSPContainerUuidMap()
	fileUuidMap, err := task.getClusterFileUuidMap(task.spec.SourceGroupSpecList)
	if err != nil {
		return nil, err
	}

	argByCluster := make(map[uuid4.Uuid]*marinaIfc.CatalogItemCreateArg)
	for _, clusterUuid := range targetClusters {
		var containerUuid *uuid4.Uuid
		if container, ok := containerByCluster[clusterUuid]; ok && container.GetContainerUuid() != "" {
			containerUuid, err = uuid4.StringToUuid4(container.GetContainerUuid())
			if err != nil {
				errMsg := fmt.Sprintf("Error encountered while creating container UUID: %v", err)
				return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
			}
		}

		arg := proto.Clone(task.arg).(*marinaIfc.CatalogItemCreateArg)
		arg.GetSpec().GlobalCatalogItemUuid = task.globalCatalogItemUuid.RawBytes()
		for _, spec := range arg.GetSpec().SourceGroupSpecList {
			task.adjustSourceGroupSpec(spec, fileUuidMap, containerUuid, clusterUuid)
		}
		argByCluster[clusterUuid] = arg
	}

	task.purgeArgWithoutImport(argByCluster)
	return argByCluster, nil
}

func (task *CatalogItemCreateTask) purgeArgWithoutImport(argByCluster map[uuid4.Uuid]*marinaIfc.CatalogItemCreateArg) {
	for clusterUuid, arg := range argByCluster {
		for _, sourceGroupSpec := range arg.GetSpec().SourceGroupSpecList {
			for _, sourceSpec := range sourceGroupSpec.SourceSpecList {
				if len(sourceSpec.GetImportSpec().LocalImportList) == 0 && len(
					sourceSpec.GetImportSpec().RemoteImportList) == 0 {
					delete(argByCluster, clusterUuid)
				}
			}
		}
	}
}

func (task *CatalogItemCreateTask) walTargetClusters(argByCluster map[uuid4.Uuid]*marinaIfc.CatalogItemCreateArg) (
	[]*uuid4.Uuid, error) {

	wal := task.Wal()
	catalogItemWal := wal.GetData().GetCatalogItem()

	// Save the cluster UUIDs we are acting on, in the case that a PE is registered in the middle of a task execution.
	var clusterUuids []*uuid4.Uuid
	if len(catalogItemWal.GetClusterUuidList()) > 0 {
		for _, peUuidBytes := range catalogItemWal.GetClusterUuidList() {
			clusterUuids = append(clusterUuids, uuid4.ToUuid4(peUuidBytes))
		}

	} else {
		// In case there is no source (remote or local) provided for catalog item creation, i.e. it's an empty catalog
		// item then fanout the request to all registered PEs.
		if len(argByCluster) == 0 {
			clusterUuids = task.ExternalInterfaces().ZeusConfig().PeClusterUuids()

		} else {
			for clusterUuid := range argByCluster {
				clusterUuids = append(clusterUuids, uuid4.ToUuid4(clusterUuid.RawBytes()))
			}
		}

		if len(clusterUuids) == 0 {
			errMsg := fmt.Sprintf("No registered clusters found which can be used to create catalog items")
			return nil, marinaError.ErrNotFound.SetCauseAndLog(errors.New(errMsg))
		}

		for i := range clusterUuids {
			catalogItemWal.ClusterUuidList = append(catalogItemWal.ClusterUuidList, clusterUuids[i].RawBytes())
		}
		err := task.SetWal(wal)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to set the task WAL: %v", err)
			return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
		}
		task.Save(task)
	}

	return clusterUuids, nil
}

func (task *CatalogItemCreateTask) checkCatalogItemExists(ret *marinaIfc.CatalogItemCreateTaskRet) (bool, error) {
	catalogItems, err := task.InternalInterfaces().CatalogItemIfc().GetCatalogItems(
		nil, task.ExternalInterfaces().CPDBIfc(), task.InternalInterfaces().UuidIfc(),
		[]*marinaIfc.CatalogItemId{{
			GlobalCatalogItemUuid: task.globalCatalogItemUuid.RawBytes(),
			Version:               proto.Int64(task.version),
		}},
		nil, false, "CatalogItemCreate:GetCatalogItems")
	if err != nil {
		return false, err
	}

	catalogItem, err := task.InternalInterfaces().CatalogItemIfc().GetCatalogItem(nil,
		task.ExternalInterfaces().CPDBIfc(), task.catalogItemUuid)
	if err != nil && err != marinaError.ErrNotFound {
		return false, err
	}

	if len(catalogItems) > 0 {
		if len(catalogItems) != 1 {
			errMsg := fmt.Sprintf("More than 1 version of catalog item %s exists",
				task.globalCatalogItemUuid.String())
			return false, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
		}

		if catalogItem != nil &&
			!task.globalCatalogItemUuid.Equals(uuid4.ToUuid4(catalogItems[0].GlobalCatalogItemUuid)) {

			gcUuid := uuid4.ToUuid4(catalogItems[0].GlobalCatalogItemUuid)
			errMsg := fmt.Sprintf("Catalog Item UUID mismatch between the args (%s) and IDF entry (%s",
				task.globalCatalogItemUuid.String(), gcUuid.String())
			return false, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
		}
	}

	ret.CatalogItemUuid = task.catalogItemUuid.RawBytes()
	ret.GlobalCatalogItemUuid = task.globalCatalogItemUuid.RawBytes()

	return catalogItem != nil || len(catalogItems) != 0, nil
}

func (task *CatalogItemCreateTask) fanoutCreateRequests(remoteClusters []*uuid4.Uuid,
	argByCluster map[uuid4.Uuid]*marinaIfc.CatalogItemCreateArg) (map[uuid4.Uuid]*ergon.Task, error) {

	taskUuidByEndpoint := make(map[uuid4.Uuid]uuid4.Uuid)
	wal := task.Wal()
	tasks := wal.GetData().GetTaskList()
	for _, taskObj := range tasks {
		var remoteClusterUuid uuid4.Uuid
		if clusterUuid := taskObj.GetClusterUuid(); clusterUuid != nil {
			remoteClusterUuid = *uuid4.ToUuid4(clusterUuid)
		}
		taskUuidByEndpoint[remoteClusterUuid] = *uuid4.ToUuid4(taskObj.GetTaskUuid())
	}

	var failureClusters []*uuid4.Uuid
	successTaskToCluster := make(map[uuid4.Uuid]uuid4.Uuid)
	for _, clusterUuid := range remoteClusters {
		var remoteTaskUuid *uuid4.Uuid
		var err error
		if taskUuid, ok := taskUuidByEndpoint[*clusterUuid]; ok {
			remoteTaskUuid = &taskUuid

		} else {
			remoteTaskUuid, err = task.InternalInterfaces().UuidIfc().New()
			if err != nil {
				return nil, err
			}

			endpointTask := &marinaIfc.EndpointTask{TaskUuid: remoteTaskUuid.RawBytes()}
			if *clusterUuid != utils.NilUuid {
				endpointTask.ClusterUuid = clusterUuid.RawBytes()
			}

			wal.GetData().TaskList = append(wal.GetData().TaskList, endpointTask)
			err = task.SetWal(wal)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to set the task WAL: %v", err)
				return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
			}
			task.Save(task)
			taskUuidByEndpoint[*clusterUuid] = *remoteTaskUuid
		}

		fanoutArg := argByCluster[*clusterUuid]
		fanoutArg.TaskUuid = remoteTaskUuid.RawBytes()
		fanoutArg.ParentTaskUuid = task.Proto().Uuid
		fanoutRet := &marinaIfc.CatalogItemCreateRet{}
		client := remoteCatalogService(task.ExternalInterfaces().ZkSession(), task.ExternalInterfaces().ZeusConfig(),
			clusterUuid, nil, nil, nil)
		err = client.SendMsg("CatalogItemCreate", fanoutArg, fanoutRet)
		if err != nil {
			endpointName := utils.GetUserVisibleId(task.ExternalInterfaces().ZeusConfig(), *clusterUuid)
			log.Errorf("[%s] Failed to fanout CatalogItemCreate request (task %s) to %s : %s", task.taskUuid.String(),
				remoteTaskUuid.String(), endpointName, err)
			failureClusters = append(failureClusters, clusterUuid)
			task.errReasons = append(task.errReasons, fmt.Sprintf("%s : %v", endpointName, err))

		} else {
			successTaskToCluster[*remoteTaskUuid] = *clusterUuid
		}
	}

	// Finding the clusters where the fanout was successful
	failedClusters := set.NewSet[uuid4.Uuid]()
	for _, clusterUuid := range failureClusters {
		failedClusters.Add(*clusterUuid)
	}

	var successClusters []*uuid4.Uuid
	for _, clusterUuid := range remoteClusters {
		if !failedClusters.Contains(*clusterUuid) {
			successClusters = append(successClusters, clusterUuid)
		}
	}
	task.fanoutSuccessClusters = successClusters

	if len(successTaskToCluster) <= 0 {
		errMsg := fmt.Sprintf("Failed to fanout CatalogItemCreate request to all remote clusters")
		return nil, marinaError.ErrCatalogTaskForwardError.SetCauseAndLog(errors.New(errMsg))
	}

	taskByCluster, err := task.pollTasks(successTaskToCluster)
	if err != nil {
		return nil, err
	}
	return taskByCluster, nil
}

func (task *CatalogItemCreateTask) pollTasks(successTaskToCluster map[uuid4.Uuid]uuid4.Uuid) (map[uuid4.Uuid]*ergon.Task, error) {
	taskMap := make(map[string]bool)
	for taskUuid := range successTaskToCluster {
		taskMap[taskUuid.String()] = true
	}

	tasks, err := task.PollAll(task, taskMap)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to poll for tasks")
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	if len(tasks) != len(successTaskToCluster) {
		errMsg := fmt.Sprintf("Failed to poll completion for tasks")
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	taskByCluster := make(map[uuid4.Uuid]*ergon.Task)
	for _, taskObj := range tasks {
		taskUuid := uuid4.ToUuid4(taskObj.Uuid)
		taskByCluster[successTaskToCluster[*taskUuid]] = taskObj
	}
	return taskByCluster, nil
}

func (task *CatalogItemCreateTask) processFanoutResult(taskByCluster map[uuid4.Uuid]*ergon.Task) ([]uuid4.Uuid, error) {
	var abortedCLusters, successfulCLusters []uuid4.Uuid
	for clusterUuid, taskObj := range taskByCluster {
		var errMsg string
		endpointName := utils.GetUserVisibleId(task.ExternalInterfaces().ZeusConfig(), clusterUuid)
		taskUuid := uuid4.ToUuid4(taskObj.Uuid)
		response := taskObj.GetResponse()
		if response == nil {
			errMsg = fmt.Sprintf("Failed to retrieve response for task %s (%s)", taskUuid.String(), endpointName)

		} else if taskObj.GetStatus() == ergon.Task_kAborted {
			log.Infof("[%s] Aborted task %s on %s", task.taskUuid.String(), taskUuid.String(), endpointName)
			abortedCLusters = append(abortedCLusters, clusterUuid)

		} else if taskObj.GetStatus() != ergon.Task_kSucceeded {
			errMsg = fmt.Sprintf("Polled task %s (%s) failed: error code = %d, detail = %s",
				taskUuid.String(), endpointName, response.GetErrorCode(), response.GetErrorDetail())

		} else {
			successfulCLusters = append(successfulCLusters, clusterUuid)
		}

		if errMsg != "" {
			log.Errorf("[%s]"+errMsg, task.taskUuid.String())
			task.errReasons = append(task.errReasons, errMsg)
		}

		if len(abortedCLusters) == len(taskByCluster) {
			task.AbortTask(nil, task, nil, nil)
		}
	}

	if len(successfulCLusters) == 0 {
		errMsg := fmt.Sprintf("Request failed on all registered AHV clusters")
		return nil, marinaError.ErrCatalogTaskForwardError.SetCauseAndLog(errors.New(errMsg))
	}

	return successfulCLusters, nil
}

func (task *CatalogItemCreateTask) createFileRepoEntries() ([]*marinaIfc.SourceGroup, error) {
	isContentAddressable := false
	wal := task.Wal()
	catalogWal := wal.GetData().GetCatalogItem()
	sourceGroupSpecs := catalogWal.GetAddSourceGroupList()
	if len(sourceGroupSpecs) > 0 && len(sourceGroupSpecs[0].GetSourceSpecList()) > 0 &&
		sourceGroupSpecs[0].GetSourceSpecList()[0].GetImportSpec() != nil {
		isContentAddressable = sourceGroupSpecs[0].GetSourceSpecList()[0].GetImportSpec().GetIsRequestContentAddressable()
	}

	catalogItemByCLuster, err := task.clusterCatalogItemMap()
	if err != nil {
		return nil, err
	}

	newSourceGroupSpecs := task.generateSourceGroupSpecs(sourceGroupSpecs, catalogItemByCLuster, isContentAddressable)

	var fileUuids []*uuid4.Uuid
	for _, fileUuid := range catalogWal.FileUuidList {
		fileUuids = append(fileUuids, uuid4.ToUuid4(fileUuid))
	}

	allocatedFileUuids := make([]*uuid4.Uuid, len(fileUuids))
	if fileUuids == nil {
		allocatedFileUuids = nil
	} else {
		copy(allocatedFileUuids, fileUuids)
	}

	var sourceGroups []*marinaIfc.SourceGroup
	fileImportSpecByFile := make(map[uuid4.Uuid]*marinaIfc.FileImportSpec)
	for _, sourceGroupSpec := range newSourceGroupSpecs {
		sourceGroup := &marinaIfc.SourceGroup{}
		if sourceGroupSpec.GetUuid() != nil {
			sourceGroup.Uuid = sourceGroupSpec.GetUuid()
		} else {
			uuid, err := task.InternalInterfaces().UuidIfc().New()
			if err != nil {
				return nil, err
			}
			sourceGroup.Uuid = uuid.RawBytes()
		}

		for _, sourceSpec := range sourceGroupSpec.SourceSpecList {
			if sourceSpec.GetCatalogItemUuid() != nil && sourceSpec.GetSnapshotUuid() != nil {
				errMsg := fmt.Sprintf("Snapshot and catalog item source types are currently not supported")
				return nil, marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))

			} else if sourceSpec.GetImportSpec() != nil {
				var fileUuid *uuid4.Uuid
				importSpec := sourceSpec.GetImportSpec()
				if importSpec.GetFileUuid() != nil {
					fileUuid = uuid4.ToUuid4(importSpec.GetFileUuid())

				} else if len(allocatedFileUuids) > 0 {
					fileUuid = allocatedFileUuids[0]
					allocatedFileUuids = allocatedFileUuids[1:]

				} else {
					uuid, err := task.InternalInterfaces().UuidIfc().New()
					if err != nil {
						return nil, err
					}
					fileUuid = uuid
				}

				sourceGroup.SourceList = append(sourceGroup.SourceList, &marinaIfc.Source{FileUuid: fileUuid.RawBytes()})
				sourceSpec.GetImportSpec().FileUuid = fileUuid.RawBytes()
				fileImportSpecByFile[*fileUuid] = sourceSpec.GetImportSpec()
				fileUuids = append(fileUuids, fileUuid)

			} else {
				errMsg := fmt.Sprintf("No source specified")
				return nil, marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
			}
		}
		sourceGroups = append(sourceGroups, sourceGroup)
	}

	if len(catalogWal.FileUuidList) == 0 {
		for _, fileUuid := range fileUuids {
			catalogWal.FileUuidList = append(catalogWal.FileUuidList, fileUuid.RawBytes())
		}
		err = task.SetWal(wal)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to set the task WAL: %v", err)
			return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
		}
		task.Save(task)
	}

	var startedTaskUuids []*uuid4.Uuid
	startedImportFileUuids := set.NewSet[uuid4.Uuid]()
	subtasks, err := task.IterSubtasks(task)
	if err != nil {
		errMsg := fmt.Sprintf("Error encountered while iterating through subtasks: %s", err)
		return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
	}
	for _, subtask := range subtasks {
		if subtask.GetRequest().GetMethodName() != "FileImport" {
			continue
		}

		embedded := subtask.GetRequest().GetArg().GetEmbedded()
		arg := &marinaIfc.FileImportArg{}
		if err = proto.Unmarshal(embedded, arg); err != nil {
			errMsg := fmt.Sprintf("Unable to umarshal argument: %s", err)
			return nil, marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
		}

		fileUuid := uuid4.ToUuid4(arg.GetSpec().FileUuid)
		startedImportFileUuids.Add(*fileUuid)
		startedTaskUuids = append(startedTaskUuids, uuid4.ToUuid4(subtask.Uuid))
	}

	for fileUuid, importSpec := range fileImportSpecByFile {
		if startedImportFileUuids.Contains(fileUuid) {
			continue
		}
		taskUuid, err := task.addFileToLocalRepo(importSpec)
		if err != nil {
			return nil, err
		}
		startedImportFileUuids.Add(fileUuid)
		startedTaskUuids = append(startedTaskUuids, taskUuid)
	}

	subTasks, err := task.RunSubtasks(task, startedTaskUuids, make([]taskUtil.Task, 0), make([]taskUtil.RemoteTask, 0))
	if err != nil {
		return nil, err
	}

	if len(subTasks) != len(startedTaskUuids) {
		errMsg := fmt.Sprintf("Incorrect subtasks found for CatalogItem create task: %s", uuid4.ToUuid4(task.Proto().Uuid).String())
		return nil, marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
	}

	for _, taskObj := range subTasks {
		if *taskObj.Status != ergon.Task_kSucceeded {
			errMsg := fmt.Sprintf("FileImport subtask %s did not succeed: %s",
				uuid4.ToUuid4(taskObj.Uuid).String(), taskObj.Status)
			return nil, marinaError.ErrImportError.SetCauseAndLog(errors.New(errMsg))

		} else {
			// Handle Content Addressable Requests
			arg := &marinaIfc.FileImportArg{}
			embedded := taskObj.GetRequest().GetArg().GetEmbedded()
			if err = proto.Unmarshal(embedded, arg); err != nil {
				errMsg := fmt.Sprintf("Unable to umarshal argument: %s", err)
				return nil, marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
			}

			ret := &marinaIfc.FileImportTaskRet{}
			embedded = taskObj.GetResponse().GetRet().GetEmbedded()
			if err = proto.Unmarshal(embedded, ret); err != nil {
				errMsg := fmt.Sprintf("Unable to umarshal argument: %s", err)
				return nil, marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
			}

			if !arg.GetSpec().GetIsRequestContentAddressable() ||
				uuid4.ToUuid4(arg.GetSpec().GetFileUuid()).Equals(uuid4.ToUuid4(ret.GetFileUuid())) {
				continue
			}

			for _, sourceGroup := range sourceGroups {
				for _, source := range sourceGroup.SourceList {
					if uuid4.ToUuid4(source.GetFileUuid()).Equals(uuid4.ToUuid4(arg.GetSpec().GetFileUuid())) {
						log.Infof("[%s] Replacing file %s with %s for content addressable request",
							task.taskUuid.String(), uuid4.ToUuid4(source.GetFileUuid()), uuid4.ToUuid4(ret.GetFileUuid()))
						source.FileUuid = ret.GetFileUuid()
					}
				}
			}
		}
	}

	return sourceGroups, nil
}

func (task *CatalogItemCreateTask) clusterCatalogItemMap() (map[uuid4.Uuid]*marinaIfc.CatalogItemInfo, error) {
	arg := &marinaIfc.CatalogItemGetArg{
		CatalogItemIdList: []*marinaIfc.CatalogItemId{{GlobalCatalogItemUuid: task.globalCatalogItemUuid.RawBytes()}},
		Latest:            proto.Bool(true),
	}

	catalogItemByCluster := make(map[uuid4.Uuid]*marinaIfc.CatalogItemInfo)
	for i := range task.fanoutSuccessClusters {
		clusterUuid := task.fanoutSuccessClusters[i]
		ret := &marinaIfc.CatalogItemGetRet{}
		client := remoteCatalogService(task.ExternalInterfaces().ZkSession(), task.ExternalInterfaces().ZeusConfig(),
			clusterUuid, nil, nil, nil)
		err := client.SendMsg("CatalogItemGet", arg, ret)
		if err != nil {
			endpointName := utils.GetUserVisibleId(task.ExternalInterfaces().ZeusConfig(), *clusterUuid)
			log.Errorf("[%s] Failed to fetch catalog item %s from %s %s: %s", task.taskUuid.String(),
				task.globalCatalogItemUuid.String(), endpointName, clusterUuid.String(), err)
			errMsg := fmt.Sprintf("Failed to fetch catalog item from %s", endpointName)
			return nil, marinaError.ErrCatalogTaskForwardError.SetCauseAndLog(errors.New(errMsg))
		}

		if len(ret.GetCatalogItemList()) >= 1 {
			catalogItemByCluster[*clusterUuid] = ret.GetCatalogItemList()[0]
		}
	}

	return catalogItemByCluster, nil
}

func (task *CatalogItemCreateTask) cleanupCatalogFanout(finalClusters []uuid4.Uuid) error {
	log.Infof("[%s] Clean-up of %s in progress", task.taskUuid.String(), task.globalCatalogItemUuid.String())
	wal := task.Wal()
	wal.GetData().TaskList = nil

	err := task.SetWal(wal)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to set the task WAL: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	task.Save(task)

	argByCluster := make(map[uuid4.Uuid]*marinaIfc.CatalogItemDeleteArg)
	arg := &marinaIfc.CatalogItemDeleteArg{
		CatalogItemId: &marinaIfc.CatalogItemId{
			GlobalCatalogItemUuid: task.globalCatalogItemUuid.RawBytes(),
			Version:               proto.Int64(task.version),
		},
	}

	var clusterUuids []*uuid4.Uuid
	for _, clusterUuid := range finalClusters {
		argByCluster[clusterUuid] = arg
		clusterUuids = append(clusterUuids, &clusterUuid)
	}

	taskByCluster, err := task.fanoutDeleteRequests(clusterUuids, argByCluster)
	if err != nil {
		return err
	}

	var failedClusters string
	for clusterUuid, taskObj := range taskByCluster {
		if taskObj == nil || *taskObj.Status != ergon.Task_kSucceeded {
			endpointName := utils.GetUserVisibleId(task.ExternalInterfaces().ZeusConfig(), clusterUuid)
			failedClusters += endpointName + ", "
			log.Errorf("[%s] Failed to cleanup Catalog Item %s at %s", task.taskUuid.String(),
				task.globalCatalogItemUuid.String(), endpointName)
		}
	}

	if failedClusters != "" {
		task.errReasons = append(task.errReasons, "Cleanup failed on %s remote clusters", failedClusters)
	}
	return nil
}

func (task *CatalogItemCreateTask) fanoutDeleteRequests(remoteClusters []*uuid4.Uuid,
	argByCluster map[uuid4.Uuid]*marinaIfc.CatalogItemDeleteArg) (map[uuid4.Uuid]*ergon.Task, error) {

	taskUuidByEndpoint := make(map[uuid4.Uuid]uuid4.Uuid)
	wal := task.Wal()
	tasks := wal.GetData().GetTaskList()
	for _, taskObj := range tasks {
		var remoteClusterUuid uuid4.Uuid
		if clusterUuid := taskObj.GetClusterUuid(); clusterUuid != nil {
			remoteClusterUuid = *uuid4.ToUuid4(clusterUuid)
		}
		taskUuidByEndpoint[remoteClusterUuid] = *uuid4.ToUuid4(taskObj.GetTaskUuid())
	}

	var failureClusters []*uuid4.Uuid
	successTaskToCluster := make(map[uuid4.Uuid]uuid4.Uuid)
	for _, clusterUuid := range remoteClusters {
		var remoteTaskUuid *uuid4.Uuid
		var err error
		if taskUuid, ok := taskUuidByEndpoint[*clusterUuid]; ok {
			remoteTaskUuid = &taskUuid

		} else {
			remoteTaskUuid, err = task.InternalInterfaces().UuidIfc().New()
			if err != nil {
				return nil, err
			}

			endpointTask := &marinaIfc.EndpointTask{TaskUuid: remoteTaskUuid.RawBytes()}
			if *clusterUuid != utils.NilUuid {
				endpointTask.ClusterUuid = clusterUuid.RawBytes()
			}

			wal.GetData().TaskList = append(wal.GetData().TaskList, endpointTask)
			err = task.SetWal(wal)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to set the task WAL: %v", err)
				return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
			}
			task.Save(task)
			taskUuidByEndpoint[*clusterUuid] = *remoteTaskUuid
		}

		fanoutArg := argByCluster[*clusterUuid]
		fanoutArg.TaskUuid = remoteTaskUuid.RawBytes()
		fanoutArg.ParentTaskUuid = task.Proto().Uuid
		fanoutRet := &marinaIfc.CatalogItemDeleteRet{}
		client := remoteCatalogService(task.ExternalInterfaces().ZkSession(), task.ExternalInterfaces().ZeusConfig(),
			clusterUuid, nil, nil, nil)
		err = client.SendMsg("CatalogItemDelete", fanoutArg, fanoutRet)
		if err != nil {
			endpointName := utils.GetUserVisibleId(task.ExternalInterfaces().ZeusConfig(), *clusterUuid)
			log.Errorf("Failed to fanout CatalogItemDelete request (task %s) to %s : %s",
				remoteTaskUuid.String(), endpointName, err)
			failureClusters = append(failureClusters, clusterUuid)
			task.errReasons = append(task.errReasons, fmt.Sprintf("%s : %v", endpointName, err))

		} else {
			successTaskToCluster[*remoteTaskUuid] = *clusterUuid
		}
	}

	// Finding the clusters where the fanout was successful
	failedClusters := set.NewSet[uuid4.Uuid]()
	for _, clusterUuid := range failureClusters {
		failedClusters.Add(*clusterUuid)
	}

	var successClusters []*uuid4.Uuid
	for _, clusterUuid := range remoteClusters {
		if !failedClusters.Contains(*clusterUuid) {
			successClusters = append(successClusters, clusterUuid)
		}
	}
	task.fanoutSuccessClusters = successClusters

	if len(successTaskToCluster) <= 0 {
		errMsg := fmt.Sprintf("Failed to fanout CatalogItemDelete request to all remote clusters")
		return nil, marinaError.ErrCatalogTaskForwardError.SetCauseAndLog(errors.New(errMsg))
	}

	taskByCluster, err := task.pollTasks(successTaskToCluster)
	if err != nil {
		return nil, err
	}
	return taskByCluster, nil
}

func (task *CatalogItemCreateTask) handleCleanup(finalClusters []uuid4.Uuid) error {
	err := task.cleanupCatalogFanout(finalClusters)
	if err != nil {
		return err
	}

	err = task.cleanupFileRepoEntries()
	if err != nil {
		return err
	}

	var errMsg string
	for _, errStr := range task.errReasons {
		errMsg += errStr
	}
	return marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
}
