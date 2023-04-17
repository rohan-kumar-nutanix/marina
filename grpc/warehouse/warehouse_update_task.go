/*
* Copyright (c) 2022 Nutanix Inc. All rights reserved.
*
* Authors: rajesh.battala@nutanix.com
*
* Implementation of Update Warehouse task.
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

type MarinaWarehouseUpdateTask struct {
	// Base Marina Task
	*MarinaBaseWarehouseTask
}

func NewMarinaWarehouseUpdateTask(baseWarehouseTask *MarinaBaseWarehouseTask) *MarinaWarehouseUpdateTask {
	return &MarinaWarehouseUpdateTask{
		MarinaBaseWarehouseTask: baseWarehouseTask,
	}
}

func (task *MarinaWarehouseUpdateTask) StartHook() error {
	arg := &warehousePB.UpdateWarehouseMetadataArg{}
	embedded := task.Proto().Request.Arg.Embedded
	if err := proto.Unmarshal(embedded, arg); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to UnMarshal UpdateWarehouseMetadataArg : %s", err))
	}
	log.Infof("Received UpdateMetadataWarehouseArg %v", arg)
	warehouseUuid, err := uuid4.StringToUuid4(*arg.ExtId)
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to convert Warehouse UUID: %s", err))
	}

	_, err = task.GetWarehouse(context.TODO(), task.ExternalInterfaces().CPDBIfc(), warehouseUuid)

	if err != nil {
		if err == marinaError.ErrNotFound {
			log.Infof("Warehouse UUID not present in IDF %v", warehouseUuid)
			return err
		} else {
			return marinaError.ErrMarinaInternal.SetCauseAndLog(
				fmt.Errorf("error occured while finding the Warehouse UUID %v in DB : %s", warehouseUuid, err))
		}
	}

	task.SetWarehouseUuid(warehouseUuid)

	wal := task.Wal()
	wal.WarehouseWal = &marinaPB.WarehouseWal{
		WarehouseUuid: task.GetWarehouseUuid().RawBytes(),
	}
	if err := task.SetWal(wal); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to set Warehouse WAL: %s", err))
	}
	log.Infof("Updating Warehouse with UUID :%s", task.WarehouseUuid)
	return nil
}

func (task *MarinaWarehouseUpdateTask) Run() error {
	ctx := task.TaskContext()
	task.SetWarehouseUuid(uuid4.ToUuid4(task.Wal().GetWarehouseWal().GetWarehouseUuid()))
	arg := task.getUpdateMetadataWarehouseArg()

	log.Infof("Running a Warehouse Update request with UUID : %s", task.GetWarehouseUuid())
	err := task.UpdateWarehouse(ctx, task.ExternalInterfaces().CPDBIfc(),
		task.InternalInterfaces().ProtoIfc(), task.GetWarehouseUuid(), arg.Body)
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

func (task *MarinaWarehouseUpdateTask) Enqueue() {
	task.ExternalInterfaces().SerialExecutor().SubmitJob(task)
}

func (task *MarinaWarehouseUpdateTask) SerializationID() string {
	return task.GetWarehouseUuid().String()
}

func (task *MarinaWarehouseUpdateTask) Execute() {
	task.Resume(task)
}

func (task *MarinaWarehouseUpdateTask) getUpdateMetadataWarehouseArg() *warehousePB.UpdateWarehouseMetadataArg {
	arg := &warehousePB.UpdateWarehouseMetadataArg{}
	proto.Unmarshal(task.Proto().Request.Arg.Embedded, arg)
	return arg
}
