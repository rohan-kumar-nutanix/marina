/*
* Copyright (c) 2022 Nutanix Inc. All rights reserved.
*
* Authors: rajesh.battala@nutanix.com
*
* Implementation of Delete Warehouse task.
 */

package warehouse

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	marinaError "github.com/nutanix-core/content-management-marina/errors"
	warehousePB "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"
	marinaPB "github.com/nutanix-core/content-management-marina/protos/marina"
)

type MarinaWarehouseSyncTask struct {
	// Base Marina Task
	*MarinaBaseWarehouseTask
}

func NewMarinaWarehouseSyncTask(baseWarehouseTask *MarinaBaseWarehouseTask) *MarinaWarehouseSyncTask {
	return &MarinaWarehouseSyncTask{
		MarinaBaseWarehouseTask: baseWarehouseTask,
	}
}

func (task *MarinaWarehouseSyncTask) StartHook() error {
	arg := &warehousePB.SyncWarehouseMetadataArg{}
	embedded := task.Proto().Request.Arg.Embedded
	if err := proto.Unmarshal(embedded, arg); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to UnMarshal SyncWarehouseMetadataArg : %s", err))
	}
	log.Infof("Received SyncWarehouseMetadataArg %v", arg)
	warehouseUuid, err := uuid4.StringToUuid4(*arg.ExtId)
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to convert Warehouse UUID: %s", err))
	}
	task.SetWarehouseUuid(warehouseUuid)
	log.Infof("Warehouse UUID to be synced %s", warehouseUuid.String())

	wal := task.Wal()
	t := time.Date(1900, 56, 76, 43, 99, 101, 3444, time.UTC)
	encoding, _ := t.MarshalBinary()
	wal.SyncWal = &marinaPB.SyncWal{
		FileCompleted: encoding,
	}
	wal.WarehouseWal = &marinaPB.WarehouseWal{
		WarehouseUuid: task.GetWarehouseUuid().RawBytes(),
	}
	if err := task.SetWal(wal); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to set Sync WAL: %s", err))
	}
	return nil
}

func (task *MarinaWarehouseSyncTask) Run() error {
	ctx := task.TaskContext()
	task.SetWarehouseUuid(uuid4.ToUuid4(task.Wal().GetWarehouseWal().GetWarehouseUuid()))
	log.Infof("Running a Warehouse Sync request UUID : %s", task.GetWarehouseUuid())
	storageImpl := newWarehouseS3StorageImpl()
	fileList, fileModifiedTime, err := storageImpl.ListAllFilesInWarehouseBucket(ctx, task.WarehouseUuid.String(), "metadata/")
	if err != nil {
		return err
	}
	var file_completed time.Time
	err = file_completed.UnmarshalBinary(task.Wal().GetSyncWal().GetFileCompleted())
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to unamrshal time bytes: %s", err))
	}
	for i, fileName := range fileList {
		if !fileModifiedTime[i].After(file_completed) {
			continue
		}
		body, err := storageImpl.GetFileFromWarehouseBucket(ctx, task.WarehouseUuid.String(), fileName)
		if err != nil {
			return err
		}
		err = task.SyncWarehouse(ctx,
			task.ExternalInterfaces().CPDBIfc(), task.WarehouseUuid.String(), fileName, body)
		if err != nil {
			return err
		}
		wal := task.Wal()
		fileModifiedTimeBytes, err := fileModifiedTime[i].MarshalBinary()
		if err != nil {
			return marinaError.ErrMarinaInternal.SetCauseAndLog(
				fmt.Errorf("Could not marshal time to bytes: %s", err))
		}
		wal.SyncWal = &marinaPB.SyncWal{
			FileCompleted: fileModifiedTimeBytes,
		}
		if err := task.SetWal(wal); err != nil {
			return marinaError.ErrMarinaInternal.SetCauseAndLog(
				fmt.Errorf("failed to set Sync WAL: %s", err))
		}
		task.Save(task)
	}

	ret := &warehousePB.SyncWarehouseMetadataRet{}
	retBytes, err := task.InternalInterfaces().ProtoIfc().Marshal(ret)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to serialize the return object: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	task.CompleteWithContext(ctx, task, retBytes, nil)
	return nil
}

func (task *MarinaWarehouseSyncTask) Enqueue() {
	task.ExternalInterfaces().SerialExecutor().SubmitJob(task)
}

func (task *MarinaWarehouseSyncTask) SerializationID() string {
	return task.GetWarehouseUuid().String()
}

func (task *MarinaWarehouseSyncTask) Execute() {
	task.Resume(task)
}
