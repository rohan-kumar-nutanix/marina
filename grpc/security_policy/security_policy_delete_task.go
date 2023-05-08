/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Authors: rajesh.battala@nutanix.com
*
* Implementation of Delete Security Policy task.
 */

package security_policy

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

type MarinaSecurityPolicyDeleteTask struct {
	// Base Marina Task
	*MarinaBaseSecurityPolicyTask
}

func NewSecurityPolicyDeleteTask(baseSecurityPolicyTask *MarinaBaseSecurityPolicyTask) *MarinaSecurityPolicyDeleteTask {
	return &MarinaSecurityPolicyDeleteTask{
		MarinaBaseSecurityPolicyTask: baseSecurityPolicyTask,
	}
}

func (task *MarinaSecurityPolicyDeleteTask) StartHook() error {
	arg := &configPB.DeleteSecurityPolicyByExtIdArg{}
	embedded := task.Proto().Request.Arg.Embedded
	if err := proto.Unmarshal(embedded, arg); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to UnMarshal DeleteSecurityPolicyByExtIdArg : %s", err))
	}
	log.Infof("Received DeleteSecurityPolicyByExtIdArg %v", arg)
	securityPolicyUuid, err := uuid4.StringToUuid4(*arg.ExtId)
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to convert SecurityPolicy UUID: %s", err))
	}
	task.SetSecurityPolicyUuid(securityPolicyUuid)
	log.Infof("SecurityPolicy UUID to be deleted %s", securityPolicyUuid.String())

	wal := task.Wal()
	wal.SecurityPolicyWal = &marinaPB.SecurityPolicyWal{
		SecurityPolicyUuid: task.GetSecurityPolicyUuid().RawBytes(),
	}
	if err := task.SetWal(wal); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to set SecurityPolicy WAL: %s", err))
	}
	return nil
}

func (task *MarinaSecurityPolicyDeleteTask) Run() error {
	ctx := task.TaskContext()
	task.SetSecurityPolicyUuid(uuid4.ToUuid4(task.Wal().GetSecurityPolicyWal().GetSecurityPolicyUuid()))
	log.Infof("Running a Security Policy Delete request UUID : %s", task.SecurityPolicyUuid)
	err := task.DeleteSecurityPolicy(ctx, task.ExternalInterfaces().IdfIfc(),
		task.ExternalInterfaces().CPDBIfc(), task.SecurityPolicyUuid.String())
	if err != nil {
		return err
	}
	ret := &configPB.DeleteSecurityPolicyByExtIdRet{}
	retBytes, err := task.InternalInterfaces().ProtoIfc().Marshal(ret)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to serialize the return object: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	task.CompleteWithContext(ctx, task, retBytes, nil)
	return nil
}

func (task *MarinaSecurityPolicyDeleteTask) Enqueue() {
	task.ExternalInterfaces().SerialExecutor().SubmitJob(task)
}

func (task *MarinaSecurityPolicyDeleteTask) SerializationID() string {
	return task.SecurityPolicyUuid.String()
}

func (task *MarinaSecurityPolicyDeleteTask) Execute() {
	task.Resume(task)
}
