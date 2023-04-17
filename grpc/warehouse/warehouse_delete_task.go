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

	"google.golang.org/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	marinaError "github.com/nutanix-core/content-management-marina/errors"
	warehousePB "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"
	marinaPB "github.com/nutanix-core/content-management-marina/protos/marina"
)

type MarinaWarehouseDeleteTask struct {
	// Base Marina Task
	*MarinaBaseWarehouseTask
}

func NewMarinaWarehouseDeleteTask(baseWarehouseTask *MarinaBaseWarehouseTask) *MarinaWarehouseDeleteTask {
	return &MarinaWarehouseDeleteTask{
		MarinaBaseWarehouseTask: baseWarehouseTask,
	}
}

func (task *MarinaWarehouseDeleteTask) StartHook() error {
	arg := &warehousePB.DeleteWarehouseArg{}
	embedded := task.Proto().Request.Arg.Embedded
	if err := proto.Unmarshal(embedded, arg); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to UnMarshal DeleteWarehouseArg : %s", err))
	}
	log.Infof("Received DeleteWarehouseArg %v", arg)
	warehouseUuid, err := uuid4.StringToUuid4(*arg.ExtId)
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to convert Warehouse UUID: %s", err))
	}
	task.SetWarehouseUuid(warehouseUuid)
	log.Infof("Warehouse UUID to be deleted %s", warehouseUuid.String())

	wal := task.Wal()
	wal.WarehouseWal = &marinaPB.WarehouseWal{
		WarehouseUuid: task.GetWarehouseUuid().RawBytes(),
	}
	if err := task.SetWal(wal); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to set Warehouse WAL: %s", err))
	}
	return nil
}

func (task *MarinaWarehouseDeleteTask) Run() error {
	ctx := task.TaskContext()
	task.SetWarehouseUuid(uuid4.ToUuid4(task.Wal().GetWarehouseWal().GetWarehouseUuid()))
	log.Infof("Running a Warehouse Delete request UUID : %s", task.GetWarehouseUuid())
	err := task.DeleteWarehouse(ctx, task.ExternalInterfaces().IdfIfc(),
		task.ExternalInterfaces().CPDBIfc(), task.WarehouseUuid.String())
	if err != nil {
		return err
	}
	ret := &warehousePB.DeleteWarehouseRet{}
	retBytes, err := task.InternalInterfaces().ProtoIfc().Marshal(ret)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to serialize the return object: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	task.CompleteWithContext(ctx, task, retBytes, nil)
	return nil
}

func (task *MarinaWarehouseDeleteTask) Enqueue() {
	task.ExternalInterfaces().SerialExecutor().SubmitJob(task)
}

func (task *MarinaWarehouseDeleteTask) SerializationID() string {
	return task.GetWarehouseUuid().String()
}

func (task *MarinaWarehouseDeleteTask) Execute() {
	task.Resume(task)
}
