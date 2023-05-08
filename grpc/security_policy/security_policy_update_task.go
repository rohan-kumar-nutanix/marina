/*
 * Copyright (c) 2022-2023 Nutanix Inc. All rights reserved.
 *
 * Authors: rajesh.battala@nutanix.com
 *
 * Updates a SecurityPolicy.
 */

package security_policy

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

type MarinaSecurityPolicyUpdateTask struct {
	// Base Marina Task
	*MarinaBaseSecurityPolicyTask
}

func NewMarinaSecurityPolicyUpdateTask(baseSecurityPolicyTask *MarinaBaseSecurityPolicyTask) *MarinaSecurityPolicyUpdateTask {
	return &MarinaSecurityPolicyUpdateTask{
		MarinaBaseSecurityPolicyTask: baseSecurityPolicyTask,
	}
}

func (task *MarinaSecurityPolicyUpdateTask) StartHook() error {
	arg := &config.UpdateSecurityPolicyByExtIdArg{}
	embedded := task.Proto().Request.Arg.Embedded
	if err := proto.Unmarshal(embedded, arg); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to UnMarshal UpdateSecurityPolicyByExtIdArg : %s", err))
	}
	log.Infof("Received UpdateSecurityPolicyByExtIdArg %v", arg)
	securityPolicyUuid, err := uuid4.StringToUuid4(*arg.ExtId)
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to convert SecurityPolicy UUID: %s", err))
	}
	task.SetSecurityPolicyUuid(securityPolicyUuid)

	wal := task.Wal()
	wal.SecurityPolicyWal = &marinaPB.SecurityPolicyWal{
		SecurityPolicyUuid: task.GetSecurityPolicyUuid().RawBytes(),
	}
	if err := task.SetWal(wal); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to set SecurityPolicy Create WAL: %s", err))
	}
	log.Infof("Updating SecurityPolicy with UUID :%s", task.SecurityPolicyUuid)
	return nil
}

func (task *MarinaSecurityPolicyUpdateTask) Run() error {
	ctx := task.TaskContext()
	task.SetSecurityPolicyUuid(uuid4.ToUuid4(task.Wal().GetSecurityPolicyWal().GetSecurityPolicyUuid()))

	arg := task.getUpdateSecurityPolicyByExtIdArg()

	log.Infof("Running a Security Policy Update with UUID : %s", task.GetSecurityPolicyUuid())
	err := task.UpdateSecurityPolicy(ctx, task.ExternalInterfaces().CPDBIfc(), task.InternalInterfaces().ProtoIfc(),
		task.GetSecurityPolicyUuid(), arg.Body)
	if err != nil {
		return err
	}
	ret := &config.UpdateSecurityPolicyByExtIdRet{}
	retBytes, err := task.InternalInterfaces().ProtoIfc().Marshal(ret)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to serialize the return object: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	task.CompleteWithContext(ctx, task, retBytes, nil)
	return nil
}

func (task *MarinaSecurityPolicyUpdateTask) Enqueue() {
	task.ExternalInterfaces().SerialExecutor().SubmitJob(task)
}

func (task *MarinaSecurityPolicyUpdateTask) SerializationID() string {
	return task.SecurityPolicyUuid.String()
}

func (task *MarinaSecurityPolicyUpdateTask) Execute() {
	task.Resume(task)
}

func (task *MarinaSecurityPolicyUpdateTask) getUpdateSecurityPolicyByExtIdArg() *config.UpdateSecurityPolicyByExtIdArg {
	arg := &config.UpdateSecurityPolicyByExtIdArg{}
	proto.Unmarshal(task.Proto().Request.Arg.Embedded, arg)
	return arg
}
