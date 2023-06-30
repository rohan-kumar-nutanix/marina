/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Authors: rajesh.battala@nutanix.com
*
* Implementation of WarehouseItem Delete task.
 */

package warehouse

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	marinaError "github.com/nutanix-core/content-management-marina/errors"
	contentPB "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"
	marinaPB "github.com/nutanix-core/content-management-marina/protos/marina"
)

type MarinaWarehouseItemDeleteTask struct {
	// Base Marina Task
	*MarinaBaseWarehouseTask
}

func NewMarinaWarehouseItemDeleteTask(baseWarehouseTask *MarinaBaseWarehouseTask) *MarinaWarehouseItemDeleteTask {
	return &MarinaWarehouseItemDeleteTask{
		MarinaBaseWarehouseTask: baseWarehouseTask,
	}
}

func (task *MarinaWarehouseItemDeleteTask) StartHook() error {
	arg := &contentPB.DeleteWarehouseItemArg{}
	embedded := task.Proto().Request.Arg.Embedded
	if err := proto.Unmarshal(embedded, arg); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to UnMarshal DeleteWarehouseItemArg : %s", err))
	}
	log.Infof("Received DeleteWarehouseItemArg %v", arg)
	warehouseItemUuid, err := uuid4.StringToUuid4(*arg.ExtId)
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to create WarehouseItem UUID: %s", err))
	}
	warehouseUuid, _ := uuid4.StringToUuid4(*arg.WarehouseExtId)
	task.SetWarehouseUuid(warehouseUuid)
	task.SetWarehouseItemUuid(warehouseItemUuid)

	wal := task.Wal()
	wal.WarehouseWal = &marinaPB.WarehouseWal{
		WarehouseUuid:     task.WarehouseUuid.RawBytes(),
		WarehouseItemUuid: task.WarehouseItemUuid.RawBytes(),
	}
	if err := task.SetWal(wal); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to set Warehouse WAL: %s", err))
	}
	return nil
}

func (task *MarinaWarehouseItemDeleteTask) Run() error {
	ctx := task.TaskContext()

	task.SetWarehouseUuid(uuid4.ToUuid4(task.Wal().GetWarehouseWal().GetWarehouseUuid()))
	task.WarehouseItemUuid = uuid4.ToUuid4(task.Wal().GetWarehouseWal().GetWarehouseItemUuid())
	log.Infof("Running DeleteWarehouseItem UUID : %s", task.GetWarehouseUuid())
	err := task.DeleteWarehouseItem(ctx, task.ExternalInterfaces().IdfIfc(), task.ExternalInterfaces().CPDBIfc(),
		task.WarehouseUuid.String(), task.WarehouseItemUuid.String())

	if err != nil {
		return err
	}
	ret := &contentPB.DeleteWarehouseItemRet{}
	retBytes, err := task.InternalInterfaces().ProtoIfc().Marshal(ret)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to serialize the return object: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	task.CompleteWithContext(ctx, task, retBytes, nil)
	return nil
}

func (task *MarinaWarehouseItemDeleteTask) Enqueue() {
	task.ExternalInterfaces().SerialExecutor().SubmitJob(task)
}

func (task *MarinaWarehouseItemDeleteTask) SerializationID() string {
	return task.WarehouseItemUuid.String()
}

func (task *MarinaWarehouseItemDeleteTask) Execute() {
	task.Resume(task)
}

func (task *MarinaWarehouseItemDeleteTask) getDeleteWarehouseItemArg() *contentPB.DeleteWarehouseItemArg {
	arg := &contentPB.DeleteWarehouseItemArg{}
	proto.Unmarshal(task.Proto().Request.Arg.Embedded, arg)
	return arg
}
