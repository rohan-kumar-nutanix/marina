/*
* Copyright (c) 2022 Nutanix Inc. All rights reserved.
*
* Authors: rajesh.battala@nutanix.com
*
* Implementation of AddWarehouseItem task.
 */

package warehouse

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	marinaError "github.com/nutanix-core/content-management-marina/errors"
	internal "github.com/nutanix-core/content-management-marina/interface/local"
	contentPB "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"
	"github.com/nutanix-core/content-management-marina/protos/apis/common/v1/response"
	marinaPB "github.com/nutanix-core/content-management-marina/protos/marina"
)

type MarinaWarehouseItemCreateTask struct {
	// Base Marina Task
	*MarinaBaseWarehouseTask
}

func NewMarinaWarehouseItemCreateTask(baseWarehouseTask *MarinaBaseWarehouseTask) *MarinaWarehouseItemCreateTask {
	return &MarinaWarehouseItemCreateTask{
		MarinaBaseWarehouseTask: baseWarehouseTask,
	}
}

func (task *MarinaWarehouseItemCreateTask) StartHook() error {
	arg := &contentPB.AddItemToWarehouseArg{}
	embedded := task.Proto().Request.Arg.Embedded
	if err := proto.Unmarshal(embedded, arg); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to UnMarshal CreateWarehouseArg : %s", err))
	}
	log.Infof("Received AddItemWarehouseArg %v", arg)
	warehouseItemUuid, err := internal.Interfaces().UuidIfc().New()
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to create Warehouse UUID: %s", err))
	}
	warehouseUuid, _ := uuid4.StringToUuid4(*arg.ExtId)
	task.SetWarehouseUuid(warehouseUuid)
	task.SetWarehouseItemUuid(warehouseItemUuid)
	task.AddWarehouseItemEntity()
	log.Infof("Created a New Warehouse Item with UUID %s", warehouseUuid.String())

	/*requestId, err := internal.Interfaces().UuidIfc().New()
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to generate a new Request ID : %s", err))
	}
	wal.RequestUuid = requestId.RawBytes()*/
	wal := task.Wal()
	wal.WarehouseWal = &marinaPB.WarehouseWal{
		WarehouseUuid:     task.WarehouseUuid.RawBytes(),
		WarehouseItemUuid: task.WarehouseItemUuid.RawBytes(),
	}
	if err := task.SetWal(wal); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to set AddWarehouseItem WAL: %s", err))
	}
	return nil
}

func (task *MarinaWarehouseItemCreateTask) Run() error {
	ctx := task.TaskContext()
	task.SetWarehouseUuid(uuid4.ToUuid4(task.Wal().GetWarehouseWal().GetWarehouseUuid()))
	task.WarehouseItemUuid = uuid4.ToUuid4(task.Wal().GetWarehouseWal().GetWarehouseItemUuid())
	log.Infof("Running AddWarehouseItem with UUID : %s", task.GetWarehouseUuid())
	arg := task.getAddItemWarehouseArg()
	if arg.GetBody().GetBase() == nil {
		arg.Body.Base = &response.ExternalizableAbstractModel{ExtId: proto.String(task.WarehouseItemUuid.String())}
	} else {
		arg.Body.Base.ExtId = proto.String(task.WarehouseItemUuid.String())
	}
	err := task.CreateWarehouseItem(ctx, task.ExternalInterfaces().CPDBIfc(), task.InternalInterfaces().ProtoIfc(),
		task.WarehouseUuid, task.WarehouseItemUuid, arg.Body)
	if err != nil {
		return err
	}
	ret := &contentPB.AddItemToWarehouseRet{}
	retBytes, err := task.InternalInterfaces().ProtoIfc().Marshal(ret)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to serialize the return object: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	task.CompleteWithContext(ctx, task, retBytes, nil)
	return nil
}

func (task *MarinaWarehouseItemCreateTask) Enqueue() {
	task.ExternalInterfaces().SerialExecutor().SubmitJob(task)
}

func (task *MarinaWarehouseItemCreateTask) SerializationID() string {
	return task.WarehouseItemUuid.String()
}

func (task *MarinaWarehouseItemCreateTask) Execute() {
	task.Resume(task)
}

func (task *MarinaWarehouseItemCreateTask) getAddItemWarehouseArg() *contentPB.AddItemToWarehouseArg {
	arg := &contentPB.AddItemToWarehouseArg{}
	proto.Unmarshal(task.Proto().Request.Arg.Embedded, arg)
	return arg
}
