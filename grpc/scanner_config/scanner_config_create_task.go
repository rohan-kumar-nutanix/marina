/*
 * Copyright (c) 2022-2023 Nutanix Inc. All rights reserved.
 *
 * Authors: rajesh.battala@nutanix.com
 *
 * Creates a SecurityPolicy.
 */

package scanner_config

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	marinaError "github.com/nutanix-core/content-management-marina/errors"
	internal "github.com/nutanix-core/content-management-marina/interface/local"
	"github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/config"
	response "github.com/nutanix-core/content-management-marina/protos/apis/common/v1/response"
	marinaPB "github.com/nutanix-core/content-management-marina/protos/marina"
)

type MarinaScannerConfigCreateTask struct {
	// Base Marina Task
	*MarinaBaseScannerConfigTask
}

func NewMarinaScannerConfigCreateTask(baseScannerConfigTask *MarinaBaseScannerConfigTask) *MarinaScannerConfigCreateTask {
	return &MarinaScannerConfigCreateTask{
		MarinaBaseScannerConfigTask: baseScannerConfigTask,
	}
}

func (task *MarinaScannerConfigCreateTask) StartHook() error {
	arg := &config.CreateScannerToolConfigArg{}
	embedded := task.Proto().Request.Arg.Embedded
	if err := proto.Unmarshal(embedded, arg); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to UnMarshal CreateScannerToolConfigArg : %s", err))
	}
	log.Infof("Received CreateScannerToolConfigArg %v", arg)
	scannerConfigUuid, err := internal.Interfaces().UuidIfc().New()
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to create ScannerConfig UUID: %s", err))
	}
	task.SetScannerConfigUuid(scannerConfigUuid)
	task.AddScannerConfigEntity()

	requestId, err := internal.Interfaces().UuidIfc().New()
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to generate a new Request ID : %s", err))
	}
	wal := task.Wal()
	wal.RequestUuid = requestId.RawBytes()
	wal.ScannerConfigWal = &marinaPB.ScannerConfigWal{
		ScannerConfigUuid: task.GetScannerConfigUuid().RawBytes(),
	}
	if err := task.SetWal(wal); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to set ScannerConfig Create WAL: %s", err))
	}
	return nil
}

func (task *MarinaScannerConfigCreateTask) Run() error {
	ctx := task.TaskContext()
	task.SetScannerConfigUuid(uuid4.ToUuid4(task.Wal().GetScannerConfigWal().GetScannerConfigUuid()))

	arg := task.getCreateScannerToolConfigArg()

	log.Infof("Running Task for Create ScannerToolConfig for UUID : %s", task.GetScannerConfigUuid())

	if arg.GetBody().GetBase() == nil {
		arg.Body.Base = &response.ExternalizableAbstractModel{ExtId: proto.String(task.ScannerConfigUuid.String())}
	} else {
		arg.Body.Base.ExtId = proto.String(task.ScannerConfigUuid.String())
	}

	err := task.CreateScannerConfig(ctx, task.ExternalInterfaces().CPDBIfc(), task.InternalInterfaces().ProtoIfc(),
		task.GetScannerConfigUuid(), arg.Body)
	if err != nil {
		return err
	}
	ret := &config.CreateScannerToolConfigRet{}
	retBytes, err := task.InternalInterfaces().ProtoIfc().Marshal(ret)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to serialize the return object: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	task.CompleteWithContext(ctx, task, retBytes, nil)
	return nil
}

func (task *MarinaScannerConfigCreateTask) Enqueue() {
	task.ExternalInterfaces().SerialExecutor().SubmitJob(task)
}

func (task *MarinaScannerConfigCreateTask) SerializationID() string {
	return task.ScannerConfigUuid.String()
}

func (task *MarinaScannerConfigCreateTask) Execute() {
	task.Resume(task)
}

func (task *MarinaScannerConfigCreateTask) getCreateScannerToolConfigArg() *config.CreateScannerToolConfigArg {
	arg := &config.CreateScannerToolConfigArg{}
	proto.Unmarshal(task.Proto().Request.Arg.Embedded, arg)
	return arg
}
