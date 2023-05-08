/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Authors: rajesh.battala@nutanix.com
*
* Implementation of Delete Scanner Configuration  task.
 */

package scanner_config

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	marinaError "github.com/nutanix-core/content-management-marina/errors"
	configPB "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/config"
	marinaPB "github.com/nutanix-core/content-management-marina/protos/marina"
)

type MarinaScannerConfigDeleteTask struct {
	// Base Marina Task
	*MarinaBaseScannerConfigTask
}

func NewMarinaScannerConfigDeleteTask(baseScannerConfigTask *MarinaBaseScannerConfigTask) *MarinaScannerConfigDeleteTask {
	return &MarinaScannerConfigDeleteTask{
		MarinaBaseScannerConfigTask: baseScannerConfigTask,
	}
}

func (task *MarinaScannerConfigDeleteTask) StartHook() error {
	arg := &configPB.DeleteScannerToolConfigByExtIdArg{}
	embedded := task.Proto().Request.Arg.Embedded
	if err := proto.Unmarshal(embedded, arg); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to UnMarshal DeleteScannerToolConfigByExtIdArg : %s", err))
	}
	log.Infof("Received DeleteScannerToolConfigByExtIdArg %v", arg)
	scannerConfigUuid, err := uuid4.StringToUuid4(*arg.ExtId)
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to convert Scanner Config UUID: %s", err))
	}
	task.SetScannerConfigUuid(scannerConfigUuid)
	log.Infof("ScannerConfig UUID to be deleted %s", scannerConfigUuid.String())

	wal := task.Wal()
	wal.ScannerConfigWal = &marinaPB.ScannerConfigWal{
		ScannerConfigUuid: task.GetScannerConfigUuid().RawBytes(),
	}
	if err := task.SetWal(wal); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to set Scanner Config WAL: %s", err))
	}
	return nil
}

func (task *MarinaScannerConfigDeleteTask) Run() error {
	ctx := task.TaskContext()
	task.SetScannerConfigUuid(uuid4.ToUuid4(task.Wal().GetScannerConfigWal().GetScannerConfigUuid()))
	log.Infof("Running a Scanner Config Delete request UUID : %s", task.ScannerConfigUuid)
	err := task.DeleteScannerConfig(ctx, task.ExternalInterfaces().IdfIfc(),
		task.ExternalInterfaces().CPDBIfc(), task.ScannerConfigUuid.String())
	if err != nil {
		return err
	}
	ret := &configPB.DeleteScannerToolConfigByExtIdRet{}
	retBytes, err := task.InternalInterfaces().ProtoIfc().Marshal(ret)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to serialize the return object: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	task.CompleteWithContext(ctx, task, retBytes, nil)
	return nil
}

func (task *MarinaScannerConfigDeleteTask) Enqueue() {
	task.ExternalInterfaces().SerialExecutor().SubmitJob(task)
}

func (task *MarinaScannerConfigDeleteTask) SerializationID() string {
	return task.ScannerConfigUuid.String()
}

func (task *MarinaScannerConfigDeleteTask) Execute() {
	task.Resume(task)
}
