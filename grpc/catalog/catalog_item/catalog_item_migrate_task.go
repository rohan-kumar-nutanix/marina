/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Author: rishabh.gupta@nutanix.com
*
* The implementation for CatalogMigratePc RPC
 */

package catalog_item

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"strings"
	"time"

	set "github.com/deckarep/golang-set/v2"
	"github.com/golang/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/acropolis"
	"github.com/nutanix-core/acs-aos-go/ergon"
	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	. "github.com/nutanix-core/acs-aos-go/insights/insights_interface/query"
	miscUtil "github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	utils "github.com/nutanix-core/content-management-marina/util"
)

const (
	intentSpecTable         = "intent_spec"
	kindColumn              = "kind"
	masterClusterUuidColumn = "_master_cluster_uuid_"
	imageKind               = "image"
)

type CatalogMigratePcTask struct {
	*CatalogItemBaseTask
	arg              *marinaIfc.CatalogMigratePcArg
	peClusterUuid    *uuid4.Uuid
	catalogItemIds   []*marinaIfc.CatalogItemId
	catalogItemTypes []marinaIfc.CatalogItemInfo_CatalogItemType
}

func NewCatalogMigratePcTask(catalogItemBaseTask *CatalogItemBaseTask) *CatalogMigratePcTask {
	return &CatalogMigratePcTask{
		CatalogItemBaseTask: catalogItemBaseTask,
	}
}

func (task *CatalogMigratePcTask) StartHook() error {
	arg, err := task.getCatalogMigratePcArg()
	if err != nil {
		return err
	}
	task.arg = arg

	if task.arg.ClusterUuid == nil {
		errMsg := fmt.Sprintf("Migrate arg is missing cluster uuid")
		return marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
	}
	task.peClusterUuid = uuid4.ToUuid4(task.arg.GetClusterUuid())

	// We accept CatalogItemId as input but reject migrating specific versions. If we have to support migrating
	// specific versions in future then accepting CatalogItemId makes that an additive change in terms or RPC.
	for _, catalogItemId := range task.arg.CatalogItemIdList {
		if catalogItemId.Version != nil {
			errMsg := fmt.Sprintf("Individual catalog item migration is currently not supported")
			return marinaError.ErrNotSupported.SetCauseAndLog(errors.New(errMsg))
		}

		if catalogItemId.GlobalCatalogItemUuid == nil {
			errMsg := fmt.Sprintf("Migrate arg is missing global catalog item uuid")
			return marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
		}
	}

	peExists := false
	for _, clusterUuid := range task.ExternalInterfaces().ZeusConfig().PeClusterUuids() {
		if task.peClusterUuid.Equals(clusterUuid) {
			peExists = true
		}
	}
	if !peExists {
		errMsg := fmt.Sprintf("PE cluster %s does not exist or is not running on top of AHV",
			task.peClusterUuid.String())
		return marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
	}

	task.catalogItemIds = task.arg.GetCatalogItemIdList()
	task.catalogItemTypes = task.arg.GetCatalogItemTypeList()
	task.AddClusterEntity(task.peClusterUuid)

	// Since catalog entity is not user visible, display name needs to be changed for UI to be in sync with image operation.
	task.Proto().DisplayName = proto.String("Import Images")

	// In case of migration of a subset of catalog items
	err = task.checkRemoteCatalogItemsExist()
	if err != nil {
		return err
	}

	wal := task.Wal()
	phase := marinaIfc.PcTaskWalRecordMigrateTaskData_kChangeOwnership
	wal.Data = &marinaIfc.PcTaskWalRecordData{
		Migrate: &marinaIfc.PcTaskWalRecordMigrateTaskData{
			TaskPhase: &phase,
		},
	}
	return task.SetWal(wal)
}

func (task *CatalogMigratePcTask) RecoverHook() error {
	arg, err := task.getCatalogMigratePcArg()
	if err != nil {
		return err
	}
	task.arg = arg

	if task.arg.ClusterUuid == nil {
		errMsg := fmt.Sprintf("Migrate arg is missing cluster uuid")
		return marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
	}
	task.peClusterUuid = uuid4.ToUuid4(task.arg.GetClusterUuid())

	task.catalogItemIds = task.arg.GetCatalogItemIdList()
	task.catalogItemTypes = task.arg.GetCatalogItemTypeList()
	return nil
}

func (task *CatalogMigratePcTask) Enqueue() {
	serialExecutor := task.ExternalInterfaces().SerialExecutor()
	serialExecutor.SubmitJob(task)
}

func (task *CatalogMigratePcTask) Execute() {
	task.Resume(task)
}

func (task *CatalogMigratePcTask) SerializationID() string {
	return task.peClusterUuid.String()
}

func (task *CatalogMigratePcTask) Run() error {
	task.taskUuid = uuid4.ToUuid4(task.Proto().Uuid)

	err := task.checkClusterReadiness()
	if err != nil {
		return err
	}

	wal := task.Wal()
	if wal.GetData().GetMigrate().GetTaskPhase() == marinaIfc.PcTaskWalRecordMigrateTaskData_kChangeOwnership {
		err = task.changeRemoteOwnership()
		if err != nil {
			return err
		}

		err = task.changePhase(marinaIfc.PcTaskWalRecordMigrateTaskData_kCreateEntries)
		if err != nil {
			return err
		}
	}

	if wal.GetData().GetMigrate().GetTaskPhase() == marinaIfc.PcTaskWalRecordMigrateTaskData_kCreateEntries {
		catalogItems, err := task.getCatalogItemsToMigrate()
		if err != nil {
			return err
		}

		err = task.createLocalCatalogEntries(catalogItems)
		if err != nil {
			return err
		}

		err = task.changePhase(marinaIfc.PcTaskWalRecordMigrateTaskData_kCommitOwnership)
		if err != nil {
			return err
		}
	}

	catalogItemByImage := make(map[uuid4.Uuid]*marinaIfc.CatalogItemInfo)
	if wal.GetData().GetMigrate().GetTaskPhase() == marinaIfc.PcTaskWalRecordMigrateTaskData_kCommitOwnership {
		catalogItemByImage, err = task.changeOwnershipForPe()
		if err != nil {
			return err
		}

		err = task.changePhase(marinaIfc.PcTaskWalRecordMigrateTaskData_kPublishImages)
		if err != nil {
			return err
		}
	}

	if wal.GetData().GetMigrate().GetTaskPhase() == marinaIfc.PcTaskWalRecordMigrateTaskData_kPublishImages {
		err = task.publishImages(catalogItemByImage)
		if err != nil {
			return err
		}

		err = task.changePhase(marinaIfc.PcTaskWalRecordMigrateTaskData_kCreateSpecs)
		if err != nil {
			return err
		}
	}

	if wal.GetData().GetMigrate().GetTaskPhase() == marinaIfc.PcTaskWalRecordMigrateTaskData_kCreateSpecs {
		log.Infof("No intent specs are created on the PC")

		err = task.changePhase(marinaIfc.PcTaskWalRecordMigrateTaskData_kMigrationComplete)
		if err != nil {
			return err
		}
	}

	ret := &marinaIfc.CatalogMigratePcTaskRet{}
	retBytes, err := task.InternalInterfaces().ProtoIfc().Marshal(ret)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to serialize the return object: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	task.Complete(task, retBytes, nil)
	return nil
}

func (task *CatalogMigratePcTask) getCatalogMigratePcArg() (*marinaIfc.CatalogMigratePcArg, error) {
	embedded := task.Proto().GetRequest().GetArg().GetEmbedded()
	arg := &marinaIfc.CatalogMigratePcArg{}
	if err := proto.Unmarshal(embedded, arg); err != nil {
		errMsg := fmt.Sprintf("Failed to unmarshal the RPC arguments: %v", err)
		return nil, marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
	}
	return arg, nil
}

func (task *CatalogMigratePcTask) checkRemoteCatalogItemsExist() error {
	if len(task.catalogItemIds) == 0 {
		return nil
	}

	remoteCatalogItems, err := task.getAllRemoteCatalogItems()
	if err != nil {
		return err
	}

	var remoteCatalogItemUuids []*uuid4.Uuid
	for _, catalogItem := range remoteCatalogItems {
		remoteCatalogItemUuids = append(remoteCatalogItemUuids, uuid4.ToUuid4(catalogItem.GlobalCatalogItemUuid))
	}

	var catalogItemUuids []*uuid4.Uuid
	for _, catalogItemId := range task.catalogItemIds {
		catalogItemUuids = append(catalogItemUuids, uuid4.ToUuid4(catalogItemId.GlobalCatalogItemUuid))
	}

	nonExistentCatalogItems := task.InternalInterfaces().UuidIfc().Difference(catalogItemUuids, remoteCatalogItemUuids)
	if len(nonExistentCatalogItems) > 0 {
		images := task.InternalInterfaces().UuidIfc().StringFromUuidPointers(nonExistentCatalogItems)
		errMsg := fmt.Sprintf("Image(s): %s doesn't exist on cluster: %s", images, task.peClusterUuid.String())
		return marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}

func (task *CatalogMigratePcTask) getAllRemoteCatalogItems() ([]*marinaIfc.CatalogItemInfo, error) {
	arg := &marinaIfc.CatalogItemGetArg{
		CatalogItemTypeList: task.catalogItemTypes,
		CatalogItemIdList:   task.catalogItemIds,
		Latest:              proto.Bool(false),
	}
	ret := &marinaIfc.CatalogItemGetRet{}
	client := remoteCatalogService(task.ExternalInterfaces().ZkSession(), task.ExternalInterfaces().ZeusConfig(),
		task.peClusterUuid, nil, nil, nil)
	err := client.SendMsg("CatalogItemGet", arg, ret)
	if err != nil {
		endpointName := utils.GetUserVisibleId(task.ExternalInterfaces().ZeusConfig(), *task.peClusterUuid)
		log.Errorf("[%s] Failed to fetch catalog items from %s %s: %v", task.taskUuid.String(), endpointName,
			task.peClusterUuid.String(), err)
		errMsg := fmt.Sprintf("Failed to fetch catalog items from %s", endpointName)
		return nil, marinaError.ErrCatalogTaskForwardError.SetCauseAndLog(errors.New(errMsg))
	}

	return ret.CatalogItemList, nil
}

func (task *CatalogMigratePcTask) checkClusterReadiness() error {
	zeusConfig := task.ExternalInterfaces().ZeusConfig()
	catalogPeRegistered := zeusConfig.CatalogPeRegistered(task.peClusterUuid)
	if catalogPeRegistered == nil {
		errMsg := fmt.Sprintf("Failed to retrieve configuration for cluster %s", task.peClusterUuid.String())
		return marinaError.ErrRetry.SetCauseAndLog(errors.New(errMsg))
	}

	if !*catalogPeRegistered {
		errMsg := fmt.Sprintf("Catalog on the cluster %s not yet registered to prism central",
			task.peClusterUuid.String())
		return marinaError.ErrRetry.SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}

func (task *CatalogMigratePcTask) changePhase(phase marinaIfc.PcTaskWalRecordMigrateTaskData_TaskPhase) error {
	wal := task.Wal()
	wal.GetData().GetMigrate().TaskPhase = &phase
	err := task.SetWal(wal)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to set the task WAL: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	task.Save(task)
	return nil
}

func (task *CatalogMigratePcTask) changeRemoteOwnership() error {
	arg := &marinaIfc.CatalogMigratePeArg{
		ClusterUuid:         task.ExternalInterfaces().ZeusConfig().ClusterUuid().RawBytes(),
		CatalogItemTypeList: task.catalogItemTypes,
		CatalogItemIdList:   task.catalogItemIds,
	}

	taskObj, err := task.fanoutMigrateRequest(arg)
	if err != nil {
		log.Errorf("[%s]Failed to change ownership on remote cluster %s : %v", task.taskUuid.String(),
			task.peClusterUuid.String(), err)
		return err
	}

	ret := &marinaIfc.CatalogMigratePeTaskRet{}
	embedded := taskObj.GetResponse().GetRet().GetEmbedded()
	if err = proto.Unmarshal(embedded, ret); err != nil {
		errMsg := fmt.Sprintf("Unable to umarshal argument: %s", err)
		return marinaError.ErrInternal.SetCauseAndLog(errors.New(errMsg))
	}

	if ret.GetIntentSpecCount() > 0 {
		err = task.waitForIntentSpecSync(ret.GetIntentSpecCount(), "catalog:intent_spec_get")
		if err != nil {
			return err
		}
	}

	return nil
}

func (task *CatalogMigratePcTask) fanoutMigrateRequest(arg *marinaIfc.CatalogMigratePeArg) (*ergon.Task, error) {
	var taskUuid *uuid4.Uuid
	wal := task.Wal()
	tasks := wal.GetData().GetTaskList()
	for _, taskObj := range tasks {
		if clusterUuid := taskObj.GetClusterUuid(); clusterUuid != nil &&
			uuid4.ToUuid4(clusterUuid).Equals(task.peClusterUuid) {

			taskUuid = uuid4.ToUuid4(taskObj.GetTaskUuid())
		}
	}

	var err error
	if taskUuid == nil {
		taskUuid, err = task.InternalInterfaces().UuidIfc().New()
		if err != nil {
			return nil, err
		}

		endpointTask := &marinaIfc.EndpointTask{TaskUuid: taskUuid.RawBytes()}
		endpointTask.ClusterUuid = task.peClusterUuid.RawBytes()

		wal.GetData().TaskList = append(wal.GetData().TaskList, endpointTask)
		err = task.SetWal(wal)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to set the task WAL: %v", err)
			return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
		}
		task.Save(task)
	}

	arg.TaskUuid = taskUuid.RawBytes()
	arg.ParentTaskUuid = task.Proto().Uuid
	ret := &marinaIfc.CatalogMigratePeRet{}
	client := remoteCatalogService(task.ExternalInterfaces().ZkSession(), task.ExternalInterfaces().ZeusConfig(),
		task.peClusterUuid, nil, nil, nil)
	err = client.SendMsg("CatalogMigratePe", arg, ret)
	fanoutSuccess := true
	if err != nil {
		endpointName := utils.GetUserVisibleId(task.ExternalInterfaces().ZeusConfig(), *task.peClusterUuid)
		log.Errorf("[%s] Failed to fanout CatalogMigratePe request (task %s) to %s : %s", task.taskUuid.String(),
			taskUuid.String(), endpointName, err)
		fanoutSuccess = false
	}

	if !fanoutSuccess {
		errMsg := fmt.Sprintf("Failed to fanout CatalogMigratePe request to all remote clusters")
		return nil, marinaError.ErrCatalogTaskForwardError.SetCauseAndLog(errors.New(errMsg))
	}

	taskObj, err := task.pollTasks(taskUuid)
	if err != nil {
		return nil, err
	}
	return taskObj, nil
}

func (task *CatalogMigratePcTask) pollTasks(taskUuid *uuid4.Uuid) (*ergon.Task, error) {
	taskMap := make(map[string]bool)
	taskMap[taskUuid.String()] = true

	tasks, err := task.PollAll(task, taskMap)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to poll for tasks")
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	if len(tasks) != 1 {
		errMsg := fmt.Sprintf("Failed to poll completion for tasks")
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	return tasks[0], nil
}

func (task *CatalogMigratePcTask) waitForIntentSpecSync(intentSpecCount int64, queryName string) error {
	log.Infof("[%s] Intent specs to be synced from PE : %d", task.taskUuid.String(), intentSpecCount)
	checkKind := IN(COL(kindColumn), STR_LIST(imageKind))
	checkOwner := EQ(COL(masterClusterUuidColumn), STR(task.peClusterUuid.String()))
	query, err := QUERY(queryName).FROM(intentSpecTable).WHERE(AND(checkKind, checkOwner)).Proto()
	if err != nil {
		errMsg := fmt.Sprintf("Error encountered while building the IDF query %s: %v", queryName, err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}

	arg := &insights_interface.GetEntitiesWithMetricsArg{Query: query}
	backoff := miscUtil.NewConstantBackoff(time.Duration(*RegistrationTaskRetryDelayMs)*time.Millisecond,
		miscUtil.UnlimitedRetry)

	for {
		ret, err := task.ExternalInterfaces().CPDBIfc().GetEntitiesWithMetrics(arg)
		if err != nil {
			return err
		}

		if len(ret.GroupResultsList) > 0 {
			syncedIntentSpecCount := ret.GroupResultsList[0].TotalEntityCount
			log.Infof("[%s] Intent specs currently synced : %d", task.taskUuid.String(), syncedIntentSpecCount)
			if *syncedIntentSpecCount >= intentSpecCount {
				break
			}
		}

		if backoff.Backoff() == miscUtil.Stop {
			errMsg := fmt.Sprintf("Failed to get intent spec count query result. Retry limit reached!")
			return marinaError.ErrMarinaInternal.SetCauseAndLog(errors.New(errMsg))

		} else {
			log.Infof("[%s] Failed to get intent spec count query result. Retrying...")
		}
	}

	log.Infof("[%s] Intent specs from PE have synced", task.taskUuid.String())
	return nil
}

func (task *CatalogMigratePcTask) getCatalogItemsToMigrate() ([]*marinaIfc.CatalogItemInfo, error) {
	remoteCatalogItems, err := task.getAllRemoteCatalogItems()
	if err != nil {
		return nil, err
	}

	guuidToOwners := make(map[uuid4.Uuid]set.Set[uuid4.Uuid])
	for _, catalogItem := range remoteCatalogItems {
		guuid := *uuid4.ToUuid4(catalogItem.GetGlobalCatalogItemUuid())
		owner := *uuid4.ToUuid4(catalogItem.GetOwnerClusterUuid())
		if _, ok := guuidToOwners[guuid]; !ok {
			guuidToOwners[guuid] = set.NewSet[uuid4.Uuid]()
		}
		guuidToOwners[guuid].Add(owner)
	}

	pcCatalogItems, err := newCatalogItemImpl().GetAllCatalogItems(context.TODO(), task.ExternalInterfaces().CPDBIfc())
	if err != nil {
		return nil, err
	}

	pcGuuids := set.NewSet[uuid4.Uuid]()
	for _, catalogItem := range pcCatalogItems {
		guuid := *uuid4.ToUuid4(catalogItem.GetGlobalCatalogItemUuid())
		pcGuuids.Add(guuid)
	}

	guuidsToMigrate := set.NewSet[uuid4.Uuid]()
	for guuid, owners := range guuidToOwners {
		guuid := guuid
		if len(owners.ToSlice()) != 1 {
			log.Warningf("[%s] PE change owner incomplete for %s catalog items. Skipping migration",
				task.taskUuid.String(), guuid.String())
			continue
		}

		owner, _ := owners.Pop()
		if owner.Equals(task.peClusterUuid) {
			log.Warningf("[%s] %s owned by %s. Skipping migration", task.taskUuid.String(), guuid.String(),
				owner.String())
			continue
		}

		if pcGuuids.Contains(guuid) {
			log.Infof("[%s] The catalog item %s already exists on PC, attempting to reconcile with existing PC metadata",
				task.taskUuid.String(), guuid.String())
			continue
		}
		guuidsToMigrate.Add(guuid)
	}

	var catalogItemsToMigrate []*marinaIfc.CatalogItemInfo
	for _, catalogItem := range remoteCatalogItems {
		guuid := *uuid4.ToUuid4(catalogItem.GetGlobalCatalogItemUuid())
		if !guuidsToMigrate.Contains(guuid) {
			continue
		}

		if catalogItem.GetOpaque() != nil {
			imageInfo := &acropolis.ImageInfo{}
			if err = proto.Unmarshal(catalogItem.GetOpaque(), imageInfo); err != nil {
				errMsg := fmt.Sprintf("Failed to unmarshal the proto: %v", err)
				return nil, marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
			}

			if imageInfo.Hidden != nil && imageInfo.GetHidden() {
				log.Warningf("[%s] %s is hidden image. Skipping migration.", task.taskUuid.String(), guuid.String())
				continue
			}
		}
		catalogItemsToMigrate = append(catalogItemsToMigrate, catalogItem)
	}
	return catalogItemsToMigrate, nil
}

func (task *CatalogMigratePcTask) createLocalCatalogEntries(catalogItems []*marinaIfc.CatalogItemInfo) error {
	for _, catalogItem := range catalogItems {
		pcCatalogItems, err := newCatalogItemImpl().GetCatalogItems(context.TODO(), task.ExternalInterfaces().CPDBIfc(),
			task.InternalInterfaces().UuidIfc(), []*marinaIfc.CatalogItemId{{
				GlobalCatalogItemUuid: catalogItem.GetGlobalCatalogItemUuid(),
				Version:               catalogItem.Version,
			}}, []marinaIfc.CatalogItemInfo_CatalogItemType{}, false, "catalog:catalog_item_get")
		if err != nil {
			return err
		}

		wal := task.Wal()
		if len(pcCatalogItems) > 0 {
			wal.GetData().GetMigrate().CatalogItemUuidList = append(wal.GetData().GetMigrate().CatalogItemUuidList,
				pcCatalogItems[0].Uuid)
			err = task.SetWal(wal)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to set the task WAL: %v", err)
				return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
			}
			task.Save(task)
			continue
		}

		var pcSourceGroup *marinaIfc.SourceGroup
		// Not walling file uuid, it will be cleared eventually by scanner. Maybe will WAL it later
		for _, sourceGroup := range catalogItem.SourceGroupList {
			var sources []*marinaIfc.Source
			for _, source := range sourceGroup.SourceList {
				newFileUuid, err := task.InternalInterfaces().UuidIfc().New()
				if err != nil {
					return err
				}

				newSource := &marinaIfc.Source{FileUuid: newFileUuid.RawBytes()}
				sources = append(sources, newSource)
				peFileUuid := uuid4.ToUuid4(source.FileUuid)
				err = task.fileRepoIfc.CreateFile(context.TODO(), task.ExternalInterfaces().CPDBIfc(),
					task.InternalInterfaces().ProtoIfc(), newFileUuid, task.peClusterUuid, peFileUuid)
				if err != nil {
					return err
				}
			}

			sourceGroupUuid, err := task.InternalInterfaces().UuidIfc().New()
			if err != nil {
				return err
			}

			if len(sources) > 0 {
				pcSourceGroup = &marinaIfc.SourceGroup{Uuid: sourceGroupUuid.RawBytes()}
				for _, source := range sources {
					pcSourceGroup.SourceList = append(pcSourceGroup.SourceList, proto.Clone(source).(*marinaIfc.Source))
				}
			}
		}

		pcCatalogItemUuid, err := task.InternalInterfaces().UuidIfc().New()
		if err != nil {
			return err
		}

		err = newCatalogItemImpl().CreateCatalogItemFromProto(context.TODO(), task.ExternalInterfaces().CPDBIfc(),
			task.InternalInterfaces().ProtoIfc(), pcCatalogItemUuid, catalogItem, []uuid4.Uuid{*task.peClusterUuid},
			[]*marinaIfc.SourceGroup{pcSourceGroup})
		if err != nil {
			return err
		}

		wal.GetData().GetMigrate().CatalogItemUuidList = append(wal.GetData().GetMigrate().CatalogItemUuidList,
			pcCatalogItemUuid.RawBytes())
		err = task.SetWal(wal)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to set the task WAL: %v", err)
			return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
		}
		task.Save(task)
	}

	return nil
}

func (task *CatalogMigratePcTask) changeOwnershipForPe() (map[uuid4.Uuid]*marinaIfc.CatalogItemInfo, error) {
	wal := task.Wal()
	catalogItemTypes := set.NewSet[marinaIfc.CatalogItemInfo_CatalogItemType]()
	for _, catalogItemTypeStr := range strings.Split(*CatalogCreatSpecItemTypeList, ",") {
		catalogItemType := marinaIfc.CatalogItemInfo_CatalogItemType_value[catalogItemTypeStr]
		catalogItemTypes.Add(marinaIfc.CatalogItemInfo_CatalogItemType(catalogItemType))
	}

	createSpecUuids := set.NewSet[uuid4.Uuid]()
	catalogItemByImage := make(map[uuid4.Uuid]*marinaIfc.CatalogItemInfo)
	for _, catalogItemUuid := range wal.GetData().GetMigrate().CatalogItemUuidList {
		catalogItem, err := newCatalogItemImpl().GetCatalogItem(context.TODO(), task.ExternalInterfaces().CPDBIfc(),
			uuid4.ToUuid4(catalogItemUuid))
		if err != nil {
			return nil, err
		}

		if catalogItem.ItemType != nil {
			if catalogItem.GetItemType() == marinaIfc.CatalogItemInfo_kImage {
				imageUuid := *uuid4.ToUuid4(catalogItem.GetGlobalCatalogItemUuid())
				if imageCatalog, ok := catalogItemByImage[imageUuid]; ok {
					if imageCatalog.GetVersion() < catalogItem.GetVersion() {
						catalogItemByImage[imageUuid] = catalogItem
					}
				} else {
					catalogItemByImage[imageUuid] = catalogItem
				}
			}

			if catalogItemTypes.Contains(*catalogItem.ItemType) {
				createSpecUuids.Add(*uuid4.ToUuid4(catalogItem.GetGlobalCatalogItemUuid()))
			}
		}
	}

	for imageUuid := range catalogItemByImage {
		imageUuid := imageUuid
		wal.GetData().GetMigrate().ImageUuidList = append(wal.GetData().GetMigrate().ImageUuidList, imageUuid.RawBytes())
	}

	for _, createSpecUuid := range createSpecUuids.ToSlice() {
		createSpecUuid := createSpecUuid
		wal.GetData().GetMigrate().CreateSpecUuidList = append(wal.GetData().GetMigrate().CreateSpecUuidList,
			createSpecUuid.RawBytes())
	}

	err := task.SetWal(wal)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to set the task WAL: %v", err)
		return nil, marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	task.Save(task)
	return catalogItemByImage, nil
}

func (task *CatalogMigratePcTask) publishImages(catalogItemByImage map[uuid4.Uuid]*marinaIfc.CatalogItemInfo) error {
	wal := task.Wal()
	imageUuids := set.NewSet[uuid4.Uuid]()

	for _, imageUuid := range wal.GetData().GetMigrate().GetImageUuidList() {
		imageUuids.Add(*uuid4.ToUuid4(imageUuid))
	}

	for imageUuid, catalogItem := range catalogItemByImage {
		imageUuid := imageUuid
		err := task.publishImage(imageUuid, catalogItem, imageUuids)
		if err != nil {
			return err
		}
	}

	for _, imageUuid := range wal.GetData().GetMigrate().GetImageUuidList() {
		imageCatalog, err := newCatalogItemImpl().GetCatalogItems(context.TODO(), task.ExternalInterfaces().CPDBIfc(),
			task.InternalInterfaces().UuidIfc(), []*marinaIfc.CatalogItemId{{
				GlobalCatalogItemUuid: imageUuid,
			}}, []marinaIfc.CatalogItemInfo_CatalogItemType{}, true, "catalog:catalog_item_get")
		if err != nil {
			return err
		}

		if len(imageCatalog) == 0 {
			errMsg := fmt.Sprintf("Catalog item for image %s not found", uuid4.ToUuid4(imageUuid).String())
			return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
		}

		err = task.publishImage(*uuid4.ToUuid4(imageUuid), imageCatalog[0], imageUuids)
		if err != nil {
			return err
		}
	}

	return nil
}

func (task *CatalogMigratePcTask) publishImage(imageUuid uuid4.Uuid, catalogItem *marinaIfc.CatalogItemInfo,
	imageUuids set.Set[uuid4.Uuid]) error {
	imageInfo := &acropolis.ImageInfo{}
	if err := proto.Unmarshal(catalogItem.GetOpaque(), imageInfo); err != nil {
		errMsg := fmt.Sprintf("Failed to unmarshal the proto: %v", err)
		return marinaError.ErrInvalidArgument.SetCauseAndLog(errors.New(errMsg))
	}

	imageInfo.OwnerClusterUuid = catalogItem.OwnerClusterUuid
	err := task.imageIfc.CreateImage(context.TODO(), task.ExternalInterfaces().CPDBIfc(),
		task.InternalInterfaces().ProtoIfc(), uuid4.ToUuid4(imageUuid.RawBytes()), imageInfo)
	if err == insights_interface.ErrIncorrectCasValue {
		log.Errorf("[%s] Failed to publish image info to IDF %s", task.taskUuid.String(), imageUuid.String())
		return nil

	} else if err != nil {
		return err
	}

	log.V(2).Infof("Image %s is successfully published to IDF", imageUuid.String())
	imageUuids.Remove(imageUuid)
	wal := task.Wal()
	wal.GetData().GetMigrate().ImageUuidList = task.InternalInterfaces().UuidIfc().UuidsToUuidBytes(imageUuids.ToSlice())
	err = task.SetWal(wal)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to set the task WAL: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	task.Save(task)
	return nil
}
