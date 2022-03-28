// Copyright (c) 2017 Nutanix Inc. All rights reserved.
//
// Author: igor.grobman@nutanix.com
//
// Implements generic task functionality
//
// Tasks are persistently stored in a WAL, from which they may be recovered
// during master failover. A basic task implementation must implement a method
// to enqueue itself to a deferred execution context ("Enqueue"), as well as a
// method to run from this context ("Run").
//
// The "Run" method is the single point of entry for initial execution as well
// as all future recovery attempts. This method must be idempotent. It may read
// and write to the task WAL entry to achieve this guarantee.
//
// Each service must implement a base task implementation that implements
// Proto() and Ergon() functions.  Each task can then implement the rest of
// Task interface.  Use composition (struct embedding) to construct each
// more specific type.  See metropolis.task and metropolis.vm packages for
// an example.
//
// Usage:
//   // Create the task, save it to the WAL, and enqueue it for deferred
//   // execution. The task is now safely persisted, and its UUID may be
//   // communicated to third parties for polling on task completion.
//   task := NewBaseTaskSubClass(taskProto, ergonSvc)
//   task.Start()
//
// Each service should also implement a TaskManager interface whose main job
// is to Dispatch the tasks for that service.  This package implements the
// rest of task WAL functionality which will call Dispatch() whenever new
// pending tasks show up.  On service restart, all tasks that are not
// complete will be dispatched as well.
//
//   // Rough outline of Dispatch() implementation
//   for _, TaskProto := range tasks {
//     // Hydrate task from the WAL,
//    // e.g. proto.Unmarshal(taskProto.InternalOpaque, task.Wal())
//     task := taskMgr.HydrateTask(taskProto)
//     task.Recover()
//   }
//
//
package task

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	glog "github.com/golang/glog"
	proto "github.com/golang/protobuf/proto"

	ergon_ifc "github.com/nutanix-core/acs-aos-go/ergon"
	ergon_client "github.com/nutanix-core/acs-aos-go/ergon/client"
	ntnx_errors "github.com/nutanix-core/acs-aos-go/nutanix/util-go/errors"
)

// task Error definition.
type taskError struct {
	*ntnx_errors.NtnxError
}

func (e *taskError) TypeOfError() int {
	return ntnx_errors.TaskErrorType
}

func TaskError(errStr string, errCode int) *taskError {
	return &taskError{ntnx_errors.NewNtnxError(errStr, errCode)}
}
func (e *taskError) SetCause(err error) *taskError {
	return &taskError{e.NtnxError.SetCause(err).(*ntnx_errors.NtnxError)}
}

var (
	ErrInternalError  = TaskError("Internal Error", 1)
	ErrNotSupported   = TaskError("Not Supported", 2)
	ErrInvalidRequest = TaskError("Invalid Request", 3)
)

const (
	taskPollTimeout   = 60
	taskCasMaxRetries = 5
)

// This interface must be implemented by structs embedding this package's
// struct
type Task interface {
	// Ergon returns an instance of ergon client.
	Ergon() ergon_client.Ergon

	// Proto of the task (ergon.Task protobuf)
	Proto() *ergon_ifc.Task

	// Enqueue Resume() for deferred execution (typically, start a goroutine)
	Enqueue()

	// StartHook Initializes the task for execution.
	//
	// This method is invoked once prior to creating the task (along with its
	// initial WAL record).
	//
	// A subclass may use this hook for synchronous argument validation, or any
	// other initial setup that can be done ahead of time. Abstract subclasses
	// should use the start() method to augment any additional funtionality
	// during this stage.
	// NOTE: This function does not have access to the same instance of
	// the object that will execute the task.  Use RecoverHook() to
	// set up internal struct members if any.
	StartHook() error

	// RecoverHook initializes the task for recovery.
	// This method is invoked every time the task starts recovery, before it is
	// enqueued for execution. A subclass may use this hook for initial setup
	// that can be done ahead of time
	RecoverHook() error

	// Run is the task implementation.

	// This method must be idempotent. It may write to and read from Wal()
	// to keep track of its progress. The implementation should use Save() to
	// save changes to WAL at critical points.
	Run() error

	// Each Task implementation must also implement Wal() that returns that
	// service's task Wal protobuf.  It cannot be specified here, since
	// the Wal type will differ for each service/component.
}

type BlockableTask interface {
	Wait(timeout time.Duration) (interface{}, error)
	Notify(change interface{})
}

// This interface must be implemented by structs embedding this package's
// struct.
type TaskManager interface {
	// Component is the name of the service implementing this dispatcher.
	Component() string

	// Ergon returns an instance of ergon client.
	Ergon() ergon_client.Ergon

	// Dispatches new pending tasks to their appropriate queues.
	Dispatch(tasks []*ergon_ifc.Task)
}

// Implemented by this package.
type BaseTask interface {
	Recover(Task)
	Complete(task Task, taskRet []byte, err error)
	Resume(task Task)
	Start(task Task, component *string, operation *string) error
	Save(Task)
}

// Implemented by this package.
type BaseTaskMgr interface {
	StartDispatcher(TaskManager) (chan bool, error)
}

type FullTask interface {
	BaseTask
	Task
}

type RemoteTask interface {
	// Gets the UUID of the task.
	GetUuid() (*uuid4.Uuid, error)

	// Call the Remote RPC.
	CallRpc(RpcName string, arg interface{}, SubtaskSequenceId uint64,
		ParentTaskUuid []byte) ([]byte, error)

	// Returns the argument to call the RPC with.
	GetArg() interface{}

	// Returns the name of the RPC (this should exactly match the method
	// name derived from the signature defined in the interface proto of
	// the service).
	RpcName() string
}

type FullTaskMgr interface {
	BaseTaskMgr
	TaskManager
}

type TaskUtilInterface interface {
	IterSubtasks(Task) ([]*ergon_ifc.Task, error)
	RunSubtasks(Task, []*uuid4.Uuid, []Task,
		[]RemoteTask) ([]*ergon_ifc.Task, error)
	PollAll(Task, map[string]bool) ([]*ergon_ifc.Task, error)
	StartSubtask(Task, []Task) error
	StartSubtaskRemote(Task, []RemoteTask, []uint64) error
	AddEntity(*ergon_ifc.Task, *ergon_ifc.EntityId)
}

type TaskUtil struct {
}

type TaskMgrUtil struct {
}

var componentErrors = map[string]int{
	"github.com/nutanix-core/acs-aos-go/metropolis":  ntnx_errors.MetropolisErrorType,
	"dpm":         ntnx_errors.DpmErrorType,
	"task":        ntnx_errors.TaskErrorType,
	"vcenter_lib": ntnx_errors.VCenterErrorType,
}

func NewTaskUtil() *TaskUtil {
	return &TaskUtil{}
}

func NewTaskMgrUtil() *TaskMgrUtil {
	return &TaskMgrUtil{}
}

// Utility to check whether a task is complete or not.
func TaskIsComplete(taskPr *ergon_ifc.Task) bool {
	return (*taskPr.Status == ergon_ifc.Task_kSucceeded ||
		*taskPr.Status == ergon_ifc.Task_kFailed ||
		*taskPr.Status == ergon_ifc.Task_kAborted)
}

// Utility to check whether a task is pending or not.
func TaskIsPending(taskPr *ergon_ifc.Task) bool {
	return !TaskIsComplete(taskPr)
}

// Utility to return enum list of "Final" Task States
func getTaskFinalStates() [3]ergon_ifc.Task_Status {
	return [3]ergon_ifc.Task_Status{ergon_ifc.Task_kSucceeded, ergon_ifc.Task_kFailed, ergon_ifc.Task_kAborted}
}

// Utility to return enum list of "Pending" Task States
func getTaskPendingStates() [3]ergon_ifc.Task_Status {
	return [3]ergon_ifc.Task_Status{ergon_ifc.Task_kQueued, ergon_ifc.Task_kRunning, ergon_ifc.Task_kSuspended}
}

// Utility to return list of "Final" Task States as string
func GetTaskFinalStatesStr() []string {
	taskFinalStates := getTaskFinalStates()
	var taskFinalStatesStr []string
	for _, item := range taskFinalStates {
		taskFinalStatesStr = append(taskFinalStatesStr, ergon_ifc.Task_Status_name[int32(item)])
	}
	return taskFinalStatesStr
}

// Utility to return list of "Pending" Task States as string
func GetTaskPendingStatesStr() []string {
	taskPendingStates := getTaskPendingStates()
	var taskPendingStatesStr []string
	for _, item := range taskPendingStates {
		taskPendingStatesStr = append(taskPendingStatesStr, ergon_ifc.Task_Status_name[int32(item)])
	}
	return taskPendingStatesStr
}

// Recover the task and enqueue it for asynchronous execution.
func (u *TaskUtil) Recover(task Task) {
	err := task.RecoverHook()
	if err != nil {
		u.Complete(task, nil, err)
		return
	}
	task.Enqueue()
}

// resume the task from the deferred execution context.
func (u *TaskUtil) Resume(task Task) {
	taskPr := task.Proto()
	taskPr.Status = ergon_ifc.Task_kRunning.Enum()
	taskPr.PercentageComplete = proto.Int32(0)
	update(task)
	err := task.Run()
	if err != nil {
		u.Complete(task, nil, err)
	}
}

// Complete a task.  One of ret or error must be specified.  ret encodes the
// protbuf return value of the task.
func (u *TaskUtil) Complete(task Task, ret []byte, err error) {
	taskPr := task.Proto()
	taskPr.Response = &ergon_ifc.MetaResponse{
		Ret: &ergon_ifc.PayloadOrEmbeddedValue{},
	}
	// Both ret and error are nil
	if ret == nil && err == nil {
		taskPr.Status = ergon_ifc.Task_kFailed.Enum()
		glog.Fatal(
			"Request to complete task, but neither return nor error specified")
	} else if ret == nil && err != nil { // Task failed with error and no ret
		taskPr.Status = ergon_ifc.Task_kFailed.Enum()
		errObj := getError(taskPr.Component, err)
		taskPr.Response.ErrorDetail = proto.String(errObj.Error())
		taskPr.Response.ErrorCode = proto.Int32(int32(errObj.GetErrorCode()))
	} else if ret != nil && err == nil { // Task Successful, no error
		taskPr.Status = ergon_ifc.Task_kSucceeded.Enum()
		taskPr.Response.Ret.Embedded = ret
	} else if ret != nil && err != nil { // Task failed with error and no ret
		taskPr.Status = ergon_ifc.Task_kFailed.Enum()
		taskPr.Response.Ret.Embedded = ret
		errObj := getError(taskPr.Component, err)
		taskPr.Response.ErrorDetail = proto.String(errObj.Error())
		taskPr.Response.ErrorCode = proto.Int32(int32(errObj.GetErrorCode()))
	}

	taskPr.CompleteTimeUsecs = proto.Uint64(uint64(
		time.Now().UnixNano() / int64(time.Microsecond)))
	update(task)
}

func getError(component *string, err error) ntnx_errors.INtnxError {
	errType := componentErrors[*component]
	if errType == 0 {
		glog.Fatal("Unknown component: " + *component)
	}
	errObj, ok := ntnx_errors.TypeAssert(err, errType)
	if !ok {
		errType := componentErrors["task"]
		errObj, ok = ntnx_errors.TypeAssert(err, errType)
		if !ok {
			glog.Info("err ", err)
			glog.Fatal(
				"Could not find error type for component " + *component)
		}
	}
	return errObj
}

func createUpdateArg(taskPr *ergon_ifc.Task,
	currState *ergon_ifc.Task) *ergon_ifc.TaskUpdateArg {
	arg := &ergon_ifc.TaskUpdateArg{
		Uuid:                 taskPr.Uuid,
		LogicalTimestamp:     currState.LogicalTimestamp,
		Response:             taskPr.Response,
		EntityList:           taskPr.EntityList,
		Message:              taskPr.Message,
		TotalSteps:           taskPr.TotalSteps,
		StepsCompleted:       taskPr.StepsCompleted,
		InternalOpaque:       taskPr.InternalOpaque,
		ParentTaskUuid:       taskPr.ParentTaskUuid,
		StartTimeUsecs:       taskPr.StartTimeUsecs,
		FailIfTaskIsCanceled: proto.Bool(false),
		SubtaskSequenceId:    taskPr.SubtaskSequenceId,
		InternalTask:         taskPr.InternalTask,
	}
	if taskPr.PercentageComplete != nil {
		arg.PercentageComplete = taskPr.PercentageComplete
	}
	if taskPr.Status != nil {
		arg.Status = taskPr.Status.Enum()
	}
	if taskPr.StepsUnit != nil {
		arg.StepsUnit = taskPr.StepsUnit.Enum()
	}
	return arg
}

// update synchronously saves the task and its WAL to persistent store.
func update(task Task) {
	taskPr := task.Proto()
	currState, err := task.Ergon().TaskGet(taskPr.Uuid)
	var numRetries int = 0
	if err != nil {
		glog.Fatal("ERGDWN Could not get task from Ergon: " + err.Error())
	}
	for {
		arg := createUpdateArg(taskPr, currState)
		ret, err := task.Ergon().TaskUpdateBase(arg)
		if err != nil && ergon_client.ErrLogicalTimestampMismatch.Equals(err) {
			currState, err = task.Ergon().TaskGet(taskPr.Uuid)
			if err != nil {
				glog.Fatal("ERGDWN Could not get task from Ergon: " + err.Error())
			}
			if *taskPr.PercentageComplete < *currState.PercentageComplete {
				taskPr.PercentageComplete = currState.PercentageComplete
			}
			numRetries += 1
			if numRetries > taskCasMaxRetries {
				glog.Fatalf("Logical Timestamp mismatch while updating task: %x", currState.Uuid)
			}
			glog.Warningf("Logical Timestamp mismatch while updating task: %x. Retrying again", currState.Uuid)

		} else if err != nil {
			glog.Fatal("ERGDWN Could not update task in Ergon: " + err.Error())
			return
		} else {
			taskPr.LogicalTimestamp = ret.Task.LogicalTimestamp
			taskPr.LastUpdatedTimeUsecs = ret.Task.LastUpdatedTimeUsecs
			return
		}
	}
}

// Start creates the task for component and operation specified,  and queues
// task for asynchronous execution.
func (u *TaskUtil) Start(
	task Task, component *string, operation *string) error {
	taskPr := task.Proto()
	taskPr.Status = ergon_ifc.Task_kQueued.Enum()
	taskPr.CreateTimeUsecs = proto.Uint64(uint64(
		time.Now().UnixNano() / int64(time.Microsecond)))
	err := task.StartHook()
	if err != nil {
		return err
	}
	create(task, component, operation)
	return nil
}

// create the task in the database.
func create(task Task, component *string, operation *string) {
	taskPr := task.Proto()
	arg := &ergon_ifc.TaskCreateArg{
		Uuid:               taskPr.GetUuid(),
		Component:          component,
		OperationType:      operation,
		CreateTimeUsecs:    taskPr.CreateTimeUsecs,
		Request:            taskPr.Request,
		Response:           taskPr.Response,
		EntityList:         taskPr.EntityList,
		Message:            taskPr.Message,
		Status:             taskPr.Status,
		TotalSteps:         taskPr.TotalSteps,
		StepsCompleted:     taskPr.StepsCompleted,
		StepsUnit:          taskPr.StepsUnit,
		InternalOpaque:     taskPr.InternalOpaque,
		ParentTaskUuid:     taskPr.ParentTaskUuid,
		StartTimeUsecs:     taskPr.StartTimeUsecs,
		SubtaskSequenceId:  taskPr.SubtaskSequenceId,
		InternalTask:       taskPr.InternalTask,
		Weight:             taskPr.Weight,
		Reason:             taskPr.Reason,
		PercentageComplete: proto.Int32(0),
		DisplayName:        taskPr.DisplayName,
		Capabilities:       taskPr.Capabilities,
		RequestContext:     taskPr.RequestContext,
	}
	ret, err := task.Ergon().TaskCreate(arg)
	if err != nil {
		glog.Fatal("ERGDWN Could not create task in Ergon: " + err.Error())
	}
	taskPr.Uuid = ret.GetUuid()
}

// Save the task WAL record to the database.  This is a synchronous operation.
func (u *TaskUtil) Save(task Task) {
	update(task)
}

// StartDispatcher starts the goroutine for dispatching tasks as they show up
// in the task database. mgr provides functionality to dispatch the task
// for the component in question.
func (u *TaskMgrUtil) StartDispatcher(mgr TaskManager) (chan bool, error) {
	// Get pending tasks.
	listArg := &ergon_ifc.TaskListArg{
		ComponentTypeList: []string{mgr.Component()},
		IncludeCompleted:  proto.Bool(false),
		LocalClusterOnly:  proto.Bool(true),
	}
	ret, err := mgr.Ergon().TaskList(listArg)
	if err != nil {
		return nil, ErrInternalError.SetCause(err)
	}
	if len(ret.MaxSequenceIdList) != 0 &&
		len(ret.MaxSequenceIdList) != 1 {
		glog.Fatal("Unexpected number of maximum sequence IDs: " +
			string(len(ret.MaxSequenceIdList)))
	}
	maxSequenceId := uint64(0)
	if len(ret.MaxSequenceIdList) != 0 {
		maxSequenceId = *ret.MaxSequenceIdList[0].SequenceId
		dispatchableTasks := []*ergon_ifc.Task{}
		for _, taskUuid := range ret.TaskUuidList {
			taskPr, err := mgr.Ergon().TaskGet(taskUuid)
			if err != nil {
				return nil, ErrInternalError.SetCause(err)
			}
			if TaskIsPending(taskPr) {
				dispatchableTasks = append(dispatchableTasks, taskPr)
			}
		}
		mgr.Dispatch(dispatchableTasks)
	}
	stop := make(chan bool)
	go pollForNewTasks(mgr, maxSequenceId, stop)
	return stop, nil
}

// pollForNewTasks waits for newly submitted tasks for the component,
// and submits them to be run by the component.
func pollForNewTasks(mgr TaskManager, lastKnownSeqId uint64,
	stop <-chan bool) {
	taskFilter := &ergon_ifc.TaskPollArg_NewTasksFilter{
		Component:           proto.String(mgr.Component()),
		LastKnownSequenceId: proto.Uint64(lastKnownSeqId),
	}
	taskFilters := []*ergon_ifc.TaskPollArg_NewTasksFilter{taskFilter}
	requestId := mgr.Component() + ":dispatch_poll"
	arg := &ergon_ifc.TaskPollArg{
		TimeoutSec:         proto.Uint64(taskPollTimeout),
		NewTasksFilterList: taskFilters,
		RequestId:          proto.String(requestId),
	}
	for {
		select {
		case <-stop:
			return
		default:
		}
		arg.NewTasksFilterList[0].LastKnownSequenceId =
			proto.Uint64(lastKnownSeqId)
		ret, err := mgr.Ergon().TaskPoll(arg)
		if err != nil {
			if ergon_client.IsRpcTimeoutError(err) {
				glog.Infof("Received client timeout during poll: Retrying.")
				continue
			}
			glog.Fatal("ERGDWN Could not poll for tasks in Ergon: " + err.Error())
		}
		if ret.Timedout != nil && *ret.Timedout {
			continue
		}
		if len(ret.NewTasks) <= 0 {
			glog.Fatal("TaskPoll returned no new tasks.")
		}
		lastKnownSeqId = *ret.NewTasks[len(ret.NewTasks)-1].SequenceId
		dispatchableTasks := []*ergon_ifc.Task{}
		for _, task := range ret.NewTasks {
			if TaskIsPending(task) {
				dispatchableTasks = append(dispatchableTasks, task)
			}
		}
		if len(dispatchableTasks) > 0 {
			mgr.Dispatch(dispatchableTasks)
		}
	}
}

func nowUsecs() uint64 {
	t := time.Now()
	secs := t.Unix() * 1e6
	return uint64(secs)
}

// To iterate over all subtasks of a given parent task.
func (u *TaskUtil) IterSubtasks(task Task) ([]*ergon_ifc.Task, error) {
	taskPr := task.Proto()
	subtasks := make([]*ergon_ifc.Task, 0)
	for _, subtaskUuid := range taskPr.SubtaskUuidList {
		subtask, err := task.Ergon().TaskGet(subtaskUuid)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Could not get task %s from Ergon: "+
				err.Error(), subtask.Uuid))
		}
		subtasks = append(subtasks, subtask)
	}
	return subtasks, nil
}

// Poll for all specified tasks and return them after they complete.
func (u *TaskUtil) PollAll(task Task, taskUuidSet map[string]bool) ([]*ergon_ifc.Task, error) {
	ergonClient := task.Ergon()
	taskList := make([]*ergon_ifc.Task, 0)
	for len(taskUuidSet) > 0 {
		completedTasksFilter := make([][]byte, 0)
		for taskUuid, _ := range taskUuidSet {
			uuid, _ := uuid4.StringToUuid4(taskUuid)
			taskUuidBytes := uuid.RawBytes()
			completedTasksFilter = append(completedTasksFilter, taskUuidBytes)
		}
		arg := &ergon_ifc.TaskPollArg{
			TimeoutSec:           proto.Uint64(taskPollTimeout),
			CompletedTasksFilter: completedTasksFilter,
		}
		ret, err := ergonClient.TaskPoll(arg)
		completedTasksFilter = nil
		if err != nil {
			if ergon_client.IsRpcTimeoutError(err) {
				glog.Infof("Received client timeout during poll: Retrying.")
				continue
			}
			return nil, errors.New("Could not poll for tasks in Ergon: " + err.Error())
		}
		if ret.Timedout != nil && *ret.Timedout {
			continue
		}
		for _, completedTask := range ret.CompletedTasks {
			delete(taskUuidSet, uuid4.ToUuid4(completedTask.Uuid).String())
			taskList = append(taskList, completedTask)
		}
	}
	return taskList, nil
}

// Called to start running local subtasks
func (u *TaskUtil) StartSubtask(task Task, subtaskList []Task) error {
	if len(subtaskList) == 0 {
		return nil
	}
	taskPr := task.Proto()
	if TaskIsComplete(taskPr) {
		return errors.New(fmt.Sprintf("Can't start the subtask of a parent task that "+
			"is completed, %s", taskPr.Status.String()))
	}

	taskPr.LastUpdatedTimeUsecs = proto.Uint64(nowUsecs())
	for _, subtask := range subtaskList {
		subtaskPr := subtask.Proto()
		subtaskPr.ParentTaskUuid = taskPr.Uuid
		taskPr.SubtaskUuidList = append(taskPr.SubtaskUuidList, subtaskPr.Uuid)
		create(subtask, subtaskPr.Component, subtaskPr.OperationType)
	}
	return nil
}

// Called to start running remote subtasks.
func (u *TaskUtil) StartSubtaskRemote(task Task, subtaskRemoteList []RemoteTask,
	remoteSubtaskSequenceList []uint64) error {
	if len(subtaskRemoteList) == 0 {
		return nil
	}
	taskPr := task.Proto()
	if TaskIsComplete(taskPr) {
		return errors.New(fmt.Sprintf("Can't start the subtask of a parent task that "+
			"is completed, %s", taskPr.Status.String()))
	}
	taskPr.LastUpdatedTimeUsecs = proto.Uint64(nowUsecs())
	for i, subtask := range subtaskRemoteList {
		arg := subtask.GetArg()
		rpc := subtask.RpcName()
		taskUuid, err := subtask.CallRpc(rpc, arg, remoteSubtaskSequenceList[i],
			taskPr.Uuid)
		if err != nil {
			return err
		}
		taskPr.SubtaskUuidList = append(taskPr.SubtaskUuidList, taskUuid)
	}
	return nil
}

// Called by the parent task to run the subtasks and return on their completion.
func (u *TaskUtil) RunSubtasks(task Task, subtaskUuidList []*uuid4.Uuid,
	subtaskList []Task, remoteSubtaskList []RemoteTask) ([]*ergon_ifc.Task, error) {
	allSubtasks, err := u.IterSubtasks(task)
	if err != nil {
		return nil, err
	}
	maxSubtaskSequenceId := uint64(0)
	if len(allSubtasks) > 0 {
		maxSubtaskSequenceId = *allSubtasks[len(allSubtasks)-1].SubtaskSequenceId
	}
	subtaskUuidSet := make(map[string]bool)
	for _, subtaskUuid := range subtaskUuidList {
		subtaskUuidSet[subtaskUuid.String()] = true
	}
	subtaskStartList := make([]Task, 0)
	subtaskRemoteStartList := make([]RemoteTask, 0)
	remoteSubtaskSequenceList := make([]uint64, 0)
	for _, subtask := range subtaskList {
		taskPr := subtask.Proto()
		maxSubtaskSequenceId += 1
		taskPr.SubtaskSequenceId = proto.Uint64(maxSubtaskSequenceId)
		subtaskStartList = append(subtaskStartList, subtask)
		subtaskUuidSet[uuid4.ToUuid4(taskPr.GetUuid()).String()] = true
	}
	for _, remoteSubtask := range remoteSubtaskList {
		taskUuid, err := remoteSubtask.GetUuid()
		if err != nil {
			return nil, err
		}
		maxSubtaskSequenceId += 1
		remoteSubtaskSequenceList = append(remoteSubtaskSequenceList, maxSubtaskSequenceId)
		subtaskRemoteStartList = append(subtaskRemoteStartList, remoteSubtask)
		subtaskUuidSet[taskUuid.String()] = true
	}
	err = u.StartSubtask(task, subtaskStartList)
	if err != nil {
		return nil, err
	}
	err = u.StartSubtaskRemote(task, subtaskRemoteStartList, remoteSubtaskSequenceList)
	if err != nil {
		return nil, err
	}
	returnSubtasks, err := u.PollAll(task, subtaskUuidSet)
	if err != nil {
		return nil, err
	}
	return returnSubtasks, nil
}

// Adds an entity to the specified task's entity list. This method is idempotent.
func (u *TaskUtil) AddEntity(task *ergon_ifc.Task, newEntity *ergon_ifc.EntityId) {
	for _, entity := range task.EntityList {
		if bytes.Equal(entity.EntityId, newEntity.EntityId) {
			return
		}
	}
	task.EntityList = append(task.EntityList, newEntity)
}
