/*
 * Copyright (c) 2022-2023 Nutanix Inc. All rights reserved.
 *
 * Authors: rajesh.battala@nutanix.com
 *
 * Creates a SecurityPolicy.
 */

package security_policy

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	marinaError "github.com/nutanix-core/content-management-marina/errors"
	internal "github.com/nutanix-core/content-management-marina/interface/local"
	"github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/config"
	marinaPB "github.com/nutanix-core/content-management-marina/protos/marina"
)

type MarinaSecurityPolicyCreateTask struct {
	// Base Marina Task
	*MarinaBaseSecurityPolicyTask
}

func NewMarinaSecurityPolicyCreateTask(baseSecurityPolicyTask *MarinaBaseSecurityPolicyTask) *MarinaSecurityPolicyCreateTask {
	return &MarinaSecurityPolicyCreateTask{
		MarinaBaseSecurityPolicyTask: baseSecurityPolicyTask,
	}
}

func (task *MarinaSecurityPolicyCreateTask) StartHook() error {
	arg := &config.CreateSecurityPolicyArg{}
	embedded := task.Proto().Request.Arg.Embedded
	if err := proto.Unmarshal(embedded, arg); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to UnMarshal CreateSecurityPolicyArg : %s", err))
	}
	log.Infof("Received CreateSecurityPolicyArg %v", arg)
	securityPolicyUuid, err := internal.Interfaces().UuidIfc().New()
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to create SecurityPolicy UUID: %s", err))
	}
	task.SetSecurityPolicyUuid(securityPolicyUuid)
	task.AddSecurityPolicyEntity()
	log.Infof("Created a New Security Policy UUID %s", securityPolicyUuid.String())

	requestId, err := internal.Interfaces().UuidIfc().New()
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to generate a new Request ID : %s", err))
	}
	wal := task.Wal()
	wal.RequestUuid = requestId.RawBytes()
	wal.SecurityPolicyWal = &marinaPB.SecurityPolicyWal{
		SecurityPolicyUuid: task.GetSecurityPolicyUuid().RawBytes(),
	}
	if err := task.SetWal(wal); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to set SecurityPolicy Create WAL: %s", err))
	}
	log.Infof("Creating SecurityPolicy with UUID :%s", task.SecurityPolicyUuid)
	return nil
}

func (task *MarinaSecurityPolicyCreateTask) Run() error {
	ctx := task.TaskContext()
	task.SetSecurityPolicyUuid(uuid4.ToUuid4(task.Wal().GetSecurityPolicyWal().GetSecurityPolicyUuid()))

	arg := task.getCreateSecurityPolicyArg()

	log.Infof("Running a Security Policy Create with UUID : %s", task.GetSecurityPolicyUuid())
	/*if arg.GetBody().GetBase() == nil {
		arg.Body.Base = &response.ExternalizableAbstractModel{ExtId: proto.String(task.GetSecurityPolicyUuid().String())}
	} else {
		arg.Body.Base.ExtId = proto.String(task.GetSecurityPolicyUuid().String())
	}*/

	err := task.CreateSecurityPolicy(ctx, task.ExternalInterfaces().CPDBIfc(), task.InternalInterfaces().ProtoIfc(),
		task.GetSecurityPolicyUuid(), arg.Body)
	if err != nil {
		return err
	}
	ret := &config.CreateSecurityPolicyRet{}
	retBytes, err := task.InternalInterfaces().ProtoIfc().Marshal(ret)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to serialize the return object: %v", err)
		return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
	}
	task.CompleteWithContext(ctx, task, retBytes, nil)
	return nil
}

func (task *MarinaSecurityPolicyCreateTask) Enqueue() {
	task.ExternalInterfaces().SerialExecutor().SubmitJob(task)
}

func (task *MarinaSecurityPolicyCreateTask) SerializationID() string {
	return task.SecurityPolicyUuid.String()
}

func (task *MarinaSecurityPolicyCreateTask) Execute() {
	task.Resume(task)
}

func (task *MarinaSecurityPolicyCreateTask) getCreateSecurityPolicyArg() *config.CreateSecurityPolicyArg {
	arg := &config.CreateSecurityPolicyArg{}
	proto.Unmarshal(task.Proto().Request.Arg.Embedded, arg)
	return arg
}
