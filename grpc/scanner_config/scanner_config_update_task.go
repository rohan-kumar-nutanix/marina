/*
 * Copyright (c) 2022-2023 Nutanix Inc. All rights reserved.
 *
 * Authors: rajesh.battala@nutanix.com
 *
 * Implementation of Update on ScannerConfiguration.
 */

package scanner_config

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	marinaError "github.com/nutanix-core/content-management-marina/errors"
	"github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/config"
	marinaPB "github.com/nutanix-core/content-management-marina/protos/marina"
)

type MarinaScannerConfigUpdateTask struct {
	// Base Marina Task
	*MarinaBaseScannerConfigTask
}

func NewMarinaScannerConfigUpdateTask(baseScannerConfigTask *MarinaBaseScannerConfigTask) *MarinaScannerConfigUpdateTask {
	return &MarinaScannerConfigUpdateTask{
		MarinaBaseScannerConfigTask: baseScannerConfigTask,
	}
}

func (task *MarinaScannerConfigUpdateTask) StartHook() error {
	arg := &config.UpdateScannerToolConfigByExtIdArg{}
	embedded := task.Proto().Request.Arg.Embedded
	if err := proto.Unmarshal(embedded, arg); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to UnMarshal UpdateScannerToolConfigByExtIdArg : %s", err))
	}
	log.Infof("Received UpdateScannerToolConfigByExtIdArg %v", arg)
	scannerConfigUuid, err := uuid4.StringToUuid4(*arg.ExtId)
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to convert ScannerConfig UUID: %s", err))
	}
	task.SetScannerConfigUuid(scannerConfigUuid)

	wal := task.Wal()
	wal.ScannerConfigWal = &marinaPB.ScannerConfigWal{
		ScannerConfigUuid: task.GetScannerConfigUuid().RawBytes(),
	}
	if err := task.SetWal(wal); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to set Scanner Config Create WAL: %s", err))
	}
	return nil
}

func (task *MarinaScannerConfigUpdateTask) Run() error {
	ctx := task.TaskContext()
	task.SetScannerConfigUuid(uuid4.ToUuid4(task.Wal().GetScannerConfigWal().GetScannerConfigUuid()))

	arg := task.getUpdateScannerToolConfigByExtIdArg()

	log.Infof("Running a Security Policy Update with UUID : %s", task.GetScannerConfigUuid())
	err := task.UpdateScannerConfig(ctx, task.ExternalInterfaces().CPDBIfc(), task.InternalInterfaces().ProtoIfc(),
		task.GetScannerConfigUuid(), arg.Body)
	if err != nil {
		return err
	}
	ret := &config.UpdateScannerToolConfigByExtIdRet{}
	retBytes, err := task.InternalInterfaces().ProtoIfc().Marshal(ret)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to serialize the return object: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	task.CompleteWithContext(ctx, task, retBytes, nil)
	return nil
}

func (task *MarinaScannerConfigUpdateTask) Enqueue() {
	task.ExternalInterfaces().SerialExecutor().SubmitJob(task)
}

func (task *MarinaScannerConfigUpdateTask) SerializationID() string {
	return task.ScannerConfigUuid.String()
}

func (task *MarinaScannerConfigUpdateTask) Execute() {
	task.Resume(task)
}

func (task *MarinaScannerConfigUpdateTask) getUpdateScannerToolConfigByExtIdArg() *config.UpdateScannerToolConfigByExtIdArg {
	arg := &config.UpdateScannerToolConfigByExtIdArg{}
	proto.Unmarshal(task.Proto().Request.Arg.Embedded, arg)
	return arg
}
