/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Authors: rajesh.battala@nutanix.com
*
* Implementation of Update WarehouseItem task.
 */

package warehouse

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	marinaError "github.com/nutanix-core/content-management-marina/errors"
	warehousePB "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"
	marinaPB "github.com/nutanix-core/content-management-marina/protos/marina"
)

type MarinaWarehouseItemUpdateTask struct {
	// Base Marina Task
	*MarinaBaseWarehouseTask
}

func NewMarinaWarehouseItemUpdateTask(baseWarehouseTask *MarinaBaseWarehouseTask) *MarinaWarehouseItemUpdateTask {
	return &MarinaWarehouseItemUpdateTask{
		MarinaBaseWarehouseTask: baseWarehouseTask,
	}
}

func (task *MarinaWarehouseItemUpdateTask) StartHook() error {
	arg := &warehousePB.UpdateWarehouseItemMetadataArg{}
	embedded := task.Proto().Request.Arg.Embedded
	if err := proto.Unmarshal(embedded, arg); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to UnMarshal UpdateWarehouseItemMetadataArg : %s", err))
	}
	log.Infof("Received UpdateWarehouseItemMetadataArg %v", arg)
	warehouseUuid, _ := uuid4.StringToUuid4(*arg.WarehouseExtId)
	warehouseItemUuid, err := uuid4.StringToUuid4(*arg.ExtId)
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to convert Warehouse UUID: %s", err))
	}

	_, err = task.GetWarehouseItem(context.TODO(), task.ExternalInterfaces().CPDBIfc(), warehouseItemUuid)

	if err != nil {
		if err == marinaError.ErrNotFound {
			log.Infof("WarehouseItem UUID not present in IDF %v", warehouseItemUuid)
			return err
		} else {
			return marinaError.ErrMarinaInternal.SetCauseAndLog(
				fmt.Errorf("error occured while finding the WarehouseItem UUID %v in DB : %s", warehouseUuid, err))
		}
	}

	task.SetWarehouseUuid(warehouseUuid)
	task.SetWarehouseItemUuid(warehouseItemUuid)

	wal := task.Wal()
	wal.WarehouseWal = &marinaPB.WarehouseWal{
		WarehouseUuid:     task.GetWarehouseUuid().RawBytes(),
		WarehouseItemUuid: task.WarehouseItemUuid.RawBytes(),
	}
	if err := task.SetWal(wal); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to set Warehouse WAL: %s", err))
	}
	log.Infof("Updating WarehouseItem with UUID :%s", task.WarehouseItemUuid)
	return nil
}

func (task *MarinaWarehouseItemUpdateTask) Run() error {
	ctx := task.TaskContext()
	task.SetWarehouseUuid(uuid4.ToUuid4(task.Wal().GetWarehouseWal().GetWarehouseUuid()))
	task.SetWarehouseItemUuid(uuid4.ToUuid4(task.Wal().GetWarehouseWal().GetWarehouseItemUuid()))
	arg := task.getUpdateWarehouseItemMetadataArg()

	log.Infof("Running a WarehouseItem Update request for UUID : %s", task.GetWarehouseUuid())
	err := task.UpdateWarehouseItem(ctx, task.ExternalInterfaces().CPDBIfc(),
		task.InternalInterfaces().ProtoIfc(), task.GetWarehouseUuid(), task.WarehouseItemUuid, arg.Body)
	if err != nil {
		return err
	}
	ret := &warehousePB.UpdateWarehouseMetadataRet{}
	retBytes, err := task.InternalInterfaces().ProtoIfc().Marshal(ret)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to serialize the return object: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	task.CompleteWithContext(ctx, task, retBytes, nil)
	return nil
}

func (task *MarinaWarehouseItemUpdateTask) Enqueue() {
	task.ExternalInterfaces().SerialExecutor().SubmitJob(task)
}

func (task *MarinaWarehouseItemUpdateTask) SerializationID() string {
	return task.GetWarehouseUuid().String()
}

func (task *MarinaWarehouseItemUpdateTask) Execute() {
	task.Resume(task)
}

func (task *MarinaWarehouseItemUpdateTask) getUpdateWarehouseItemMetadataArg() *warehousePB.UpdateWarehouseItemMetadataArg {
	arg := &warehousePB.UpdateWarehouseItemMetadataArg{}
	proto.Unmarshal(task.Proto().Request.Arg.Embedded, arg)
	return arg
}
