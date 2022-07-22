/*
 Copyright (c) 2022 Nutanix Inc. All rights reserved.

 Authors: rajesh.battala@nutanix.com

 This implements the Marina base task, which includes the interface BaseTask as
 defined in ergon/task/task_util.go. The BaseTask interface is set with the
 pointer to TaskUtil (returned by ergonTask.NewTaskUtil()), which in turn
 implements the following methods requires by BaseTask:
 - Start (run StartHook() and create a new task in Ergon. This is usually called
by the Marina client to create a new Marina task.)
 - Save (update the task content in Ergon. This is usually called to save the
task WAL into Ergon)
 - Recover (run RecoverHook() and Enqueue() to enqueue the task)
 - Resume (move the task to running state, and call Run())
 - Complete (complete the task)

 TaskUtil also implements TaskUtilInterface which requires implementation of
 IterSubtasks(), RunSubtasks(), PollAll(), StartSubtask(), StartSubtaskRemote(),
 and AddEntity().

 A Marina task that embeds MarinaBaseTask must implement Enqueue() to enqueue
 the Marina task to an execution queue. It also can override the following
 methods that are defined in ergonTask.Task interface which is included in
 ergonTask.FullTask (which in turn used by MarinaTaskManager):
 - StartHook: to be invoked when the task is created
 - RecoverHook: to be invoked when the task is recovered
 - Run: to run the task

 Here is a brief description of the lifecycle of a task.
 - on the leader node, we run a task manager which runs a task dispatcher by
 getting tasks from Ergon to dispatch
 - to run a task, create a task that embeds NewMarinaBaseTask
 - call the method "Start()" for the above created task, which will in turn
 generate a new Ergon task, and save it to IDF through Ergon in queued state
 - the task manager periodically get queued Marina tasks from Ergon, and
 dispatch them by first hydrating a Marina task out from the Ergon task, and
 then enqueue the task by calling Enqueue()
 - the task is enqueued using a serial executor based on the Global Catalog UUID that is
 returned by SerializationID() for CatalogItem task serialization
 - once a task is ready to run, it executes Resume() which in turn invokes Run()
*/

package base

import (
	"fmt"
	"github.com/nutanix-core/acs-aos-go/ergon"
	ergonClient "github.com/nutanix-core/acs-aos-go/ergon/client"
	ergonTask "github.com/nutanix-core/acs-aos-go/ergon/task"
	marinaError "github.com/nutanix-core/content-management-marina/error"
	marinaProtos "github.com/nutanix-core/content-management-marina/protos/marina"
	util "github.com/nutanix-core/content-management-marina/util"

	"github.com/golang/protobuf/proto"
	glog "k8s.io/klog/v2"
)

type MarinaBaseTaskInterface interface {
	Proto() *ergon.Task
	Wal() *marinaProtos.PcTaskWalRecord
	SetWal(wal *marinaProtos.PcTaskWalRecord) error
}

// MarinaBaseTask defines a common interface for all Marina tasks. BaseTask is
// an interface that implements Start, Save, Recover, Resume, and Complete.
// Task implements Ergon, Proto, Enqueue, StartHook, RecoverHook, and Run.
type MarinaBaseTask struct {
	// BaseTask is included so as to access methods such as Start(), and
	// Complete().
	ergonTask.BaseTask
	// TaskUtilInterface is included so as to access methods such as PollAll()
	// and AddEntity().
	ergonTask.TaskUtilInterface
	// MarinaBaseTaskInterface is included for unit tests with SetWal.
	MarinaBaseTaskInterface
}

// NewMarinaBaseTask creates a new Marina base task.
func NewMarinaBaseTask(taskProto *ergon.Task) *MarinaBaseTask {
	baseTask := &MarinaBaseTask{
		// note that TaskUtil implements methods in both BaseTask and
		// TaskUtilInterface, and so we use it to initialize both interfaces.
		BaseTask:                ergonTask.NewTaskUtil(),
		TaskUtilInterface:       ergonTask.NewTaskUtil(),
		MarinaBaseTaskInterface: NewMarinaBaseTaskUtil(taskProto),
	}
	if err := proto.Unmarshal(taskProto.GetInternalOpaque(), baseTask.Wal()); err != nil {
		glog.Error("Failed to unmarshal WAL for a new base task.")
	}
	return baseTask
}

// Ergon returns an Ergon client.
func (t *MarinaBaseTask) Ergon() ergonClient.Ergon {
	return util.GetErgonService()
}

// StartHook initializes for task execution. A Marina task that embeds
// MarinaBaseTask can optionally override StartHook.
func (t *MarinaBaseTask) StartHook() error {
	return nil
}

// RecoverHook initializes for task recovery. A Marina task that embeds
// MarinaBaseTask can optionally override RecoverHook.
func (t *MarinaBaseTask) RecoverHook() error {
	return nil
}

type MarinaBaseTaskUtil struct {
	// Ergon task proto.
	proto *ergon.Task
	// Marina task WAL.
	wal *marinaProtos.PcTaskWalRecord
}

func NewMarinaBaseTaskUtil(taskProto *ergon.Task) *MarinaBaseTaskUtil {
	return &MarinaBaseTaskUtil{
		proto: taskProto,
		wal:   &marinaProtos.PcTaskWalRecord{},
	}
}

func (t *MarinaBaseTaskUtil) Proto() *ergon.Task {
	return t.proto
}

// Wal returns the Marina task WAL.
func (t *MarinaBaseTaskUtil) Wal() *marinaProtos.PcTaskWalRecord {
	return t.wal
}

// SetWal set the task WAL to ergon task proto.
func (t *MarinaBaseTaskUtil) SetWal(wal *marinaProtos.PcTaskWalRecord) error {
	data, err := util.Marshal(wal)
	if err != nil {
		return marinaError.ErrInternal.SetCause(
			fmt.Errorf("could not marshal WAL: %s", err))
	}
	t.proto.InternalOpaque = data
	return nil
}