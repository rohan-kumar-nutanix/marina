/*
* Copyright (c) 2022 Nutanix Inc. All rights reserved.
*
* Authors: rajesh.battala@nutanix.com
*
* Implementation of Creating Warehouse task.
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
	internal "github.com/nutanix-core/content-management-marina/interface/local"
	warehousePB "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"
	"github.com/nutanix-core/content-management-marina/protos/apis/common/v1/response"
	marinaPB "github.com/nutanix-core/content-management-marina/protos/marina"
)

type MarinaWarehouseCreateTask struct {
	// Base Marina Task
	*MarinaBaseWarehouseTask
}

func NewMarinaWarehouseCreateTask(baseWarehouseTask *MarinaBaseWarehouseTask) *MarinaWarehouseCreateTask {
	return &MarinaWarehouseCreateTask{
		MarinaBaseWarehouseTask: baseWarehouseTask,
	}
}

func (task *MarinaWarehouseCreateTask) StartHook() error {
	arg := &warehousePB.CreateWarehouseArg{}
	embedded := task.Proto().Request.Arg.Embedded
	if err := proto.Unmarshal(embedded, arg); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to UnMarshal CreateWarehouseArg : %s", err))
	}
	log.Infof("Received CreateWarehouseArg %v", arg)
	warehouseUuid, err := internal.Interfaces().UuidIfc().New()
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to create Warehouse UUID: %s", err))
	}
	task.SetWarehouseUuid(warehouseUuid)
	task.AddWarehouseEntity()
	log.Infof("Created a New Warehouse UUID %s", warehouseUuid.String())

	requestId, err := internal.Interfaces().UuidIfc().New()
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to generate a new Request ID : %s", err))
	}
	wal := task.Wal()
	wal.RequestUuid = requestId.RawBytes()
	wal.WarehouseWal = &marinaPB.WarehouseWal{
		WarehouseUuid: task.GetWarehouseUuid().RawBytes(),
	}
	if err := task.SetWal(wal); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to set Warehouse Create WAL: %s", err))
	}
	log.Infof("Creating Warehouse with UUID :%s", task.WarehouseUuid)
	return nil
}

func (task *MarinaWarehouseCreateTask) Run() error {
	ctx := task.TaskContext()
	task.SetWarehouseUuid(uuid4.ToUuid4(task.Wal().GetWarehouseWal().GetWarehouseUuid()))
	arg := task.getCreateWarehouseArg()

	log.Infof("Running a Warehouse Create with UUID : %s", task.GetWarehouseUuid())

	if arg.Body.Base == nil {
		arg.Body.Base = &response.ExternalizableAbstractModel{ExtId: proto.String(task.GetWarehouseUuid().String())}
	} else {
		arg.Body.Base.ExtId = proto.String(task.GetWarehouseUuid().String())
	}

	err := task.CreateWarehouse(ctx, task.ExternalInterfaces().CPDBIfc(), task.InternalInterfaces().ProtoIfc(),
		task.GetWarehouseUuid(), arg.Body)
	if err != nil {
		return err
	}

	// Create Warehouse Storage bucket
	task.CreateWarehouseBucket(context.TODO(), task.GetWarehouseUuid().String())
	ret := &warehousePB.CreateWarehouseRet{}
	retBytes, err := task.InternalInterfaces().ProtoIfc().Marshal(ret)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to serialize the return object: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	task.CompleteWithContext(ctx, task, retBytes, nil)
	return nil
}

func (task *MarinaWarehouseCreateTask) Enqueue() {
	task.ExternalInterfaces().SerialExecutor().SubmitJob(task)
}

func (task *MarinaWarehouseCreateTask) SerializationID() string {
	return task.GetWarehouseUuid().String()
}

func (task *MarinaWarehouseCreateTask) Execute() {
	task.Resume(task)
}

func (task *MarinaWarehouseCreateTask) getCreateWarehouseArg() *warehousePB.CreateWarehouseArg {
	arg := &warehousePB.CreateWarehouseArg{}
	proto.Unmarshal(task.Proto().Request.Arg.Embedded, arg)
	return arg
}
