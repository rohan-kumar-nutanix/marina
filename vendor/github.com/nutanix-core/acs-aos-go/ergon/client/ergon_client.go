// Copyright (c) 2016 Nutanix Inc. All rights reserved.
//
// Author: akshay@nutanix.com

// This package implements Ergon client.
package ergon_client

import (
	"errors"
	"flag"
	"time"
	"github.com/nutanix-core/acs-aos-go/zeus"

	glog "github.com/golang/glog"
	"github.com/golang/protobuf/proto"

	ergon_pr "github.com/nutanix-core/acs-aos-go/ergon"

	aplos_transport "github.com/nutanix-core/acs-aos-go/aplos/client/transport"
	util_base "github.com/nutanix-core/acs-aos-go/nutanix/util-go/base"
	ntnx_errors "github.com/nutanix-core/acs-aos-go/nutanix/util-go/errors"
	util_misc "github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
	util_net "github.com/nutanix-core/acs-aos-go/nutanix/util-go/net"
	uuid4 "github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	net "github.com/nutanix-core/acs-aos-go/nutanix/util-slbufs/util/sl_bufs/net"
)

// Constants for Ergon service.
const (
	ErgonServiceName                 = "nutanix.ergon.ErgonRpcSvc"
	clientRetryInitialDelayMilliSecs = 250
	clientRetryMaxDelaySecs          = 2
	clientRetryTimeoutSecs           = 30
	clientMaxNumRetries              = 5
	taskPollBufferSecs               = 10
)

var (
	DefaultErgonAddr    string
	DefaultErgonPort    uint16 = 2090
	task_get_batch_size int64
)

func init() {
	flag.StringVar(&DefaultErgonAddr,
		"ergon_ip",
		"127.0.0.1",
		"IP address of Ergon service to connect.")

	flag.Var(&util_misc.PortNumber{&DefaultErgonPort},
		"ergon_port",
		"Port number of the Ergon service to connect.")

	flag.Int64Var(&task_get_batch_size,
		"task_get_batch_size",
		100,
		"Default batch size of TaskGet response.")
}

// Ergon Error definition.
type ErgonError_ struct {
	*ntnx_errors.NtnxError
}

func (e *ErgonError_) TypeOfError() int {
	return ntnx_errors.ErgonErrorType
}

func (e *ErgonError_) SetCause(err error) error {
	return ntnx_errors.NewNtnxErrorRef(err, e.GetErrorDetail(),
		e.GetErrorCode(), e)
}

func (e *ErgonError_) Equals(err interface{}) bool {
	if obj, ok := err.(*ErgonError_); ok {
		return e == obj
	}
	return false
}

func ErgonError(errMsg string, errCode int) *ErgonError_ {
	return &ErgonError_{ntnx_errors.NewNtnxError(errMsg, errCode)}
}

var (
	// Ergon Errors.
	ErrNoError                       = ErgonError("NoError", 0)
	ErrCanceled                      = ErgonError("Canceled", 1)
	ErrRetry                         = ErgonError("Retry", 2)
	ErrTimeout                       = ErgonError("Timeout", 3)
	ErrNotSupported                  = ErgonError("NotSupported", 4)
	ErrUncaughtException             = ErgonError("UncaughtException", 5)
	ErrInvalidArgument               = ErgonError("InvalidArgument", 6)
	ErrInvalidState                  = ErgonError("InvalidState", 7)
	ErrLogicalTimestampMismatch      = ErgonError("LogicalTimestampMismatch", 8)
	ErrExists                        = ErgonError("Exists", 9)
	ErrNotFound                      = ErgonError("NotFound", 10)
	ErrNotMaster                     = ErgonError("NotMaster", 11)
	ErrTransportError                = ErgonError("TransportError", 12)
	ErrTaskCanceled                  = ErgonError("TaskCanceled", 13)
	ErrTaskHasPendingSubtasks        = ErgonError("TaskHasPendingSubtasks", 14)
	ErrNotAuthorized                 = ErgonError("NotAuthorized", 15)
	ErrInvalidRequestOnCompletedTask = ErgonError("InvalidRequestOnCompletedTask", 16)

	// Ergon binding related errors
	ErrInvalidRpcResponse = ErgonError("Invalid RPC Response", 101)
)

//Ergon unmapped exception error
var (
	ErrUnmappedResponseCode = ErgonError("UnmappedResponseCode", 1001)
)

// Errors returned from Ergon RPC
var Errors = map[int]*ErgonError_{
	0:  ErrNoError,
	1:  ErrCanceled,
	2:  ErrRetry,
	3:  ErrTimeout,
	4:  ErrNotSupported,
	5:  ErrUncaughtException,
	6:  ErrInvalidArgument,
	7:  ErrInvalidState,
	8:  ErrLogicalTimestampMismatch,
	9:  ErrExists,
	10: ErrNotFound,
	11: ErrNotMaster,
	12: ErrTransportError,
	13: ErrTaskCanceled,
	14: ErrTaskHasPendingSubtasks,
	15: ErrNotAuthorized,
	16: ErrInvalidRequestOnCompletedTask,
}

type Ergon interface {
	TaskCreate(*ergon_pr.TaskCreateArg) (*ergon_pr.TaskCreateRet, error)
	TaskGet([]byte) (*ergon_pr.Task, error)
	TaskGetMultiple(*ergon_pr.TaskGetArg) (*ergon_pr.TaskGetRet, error)
	GetTaskError(*ergon_pr.Task) error
	TaskUpdateBase(*ergon_pr.TaskUpdateArg) (*ergon_pr.TaskUpdateRet, error)
	TaskUpdate(*ergon_pr.TaskUpdateArg) (*ergon_pr.TaskUpdateRet, error)
	TaskPollWait(taskUuid []byte, taskRet proto.Message, intervalSecs int64,
		timeoutSecs int64) (ergon_pr.Task_Status, *ergon_pr.MetaResponse, error)
	TaskPoll(*ergon_pr.TaskPollArg) (*ergon_pr.TaskPollRet, error)
	TaskList(*ergon_pr.TaskListArg) (*ergon_pr.TaskListRet, error)
	TaskCancel(*ergon_pr.TaskCancelArg) (*ergon_pr.TaskCancelRet, error)
	TaskCancelAny(*ergon_pr.TaskCancelAnyArg) (*ergon_pr.TaskCancelAnyRet, error)
}

// Structure to hold Ergon client side details.
type ergonClient struct {
	client util_net.ProtobufRPCClientIfc
}

// Use this function to instantiate new ergon service.
// Args:
//	serverIp   : IP address of server.
// 	serverPort : Port of ergon server.
// Return:
// 	Pointer to Ergon structure.
func NewErgonService(serverIp string, serverPort uint16) Ergon {
	glog.Info("Initializing Ergon Service")

	return &ergonClient{
		client: NewErgonClient(serverIp, serverPort),
	}
}

// Used to intantiate a remote ergon client via FanOut/RemoteConnection.
func NewRemoteErgonService(zkSession *zeus.ZookeeperSession,
	clusterUUID, userUUID, tenantUUID *uuid4.Uuid, userName *string,
	reqContext *net.RpcRequestContext) Ergon {
	return &ergonClient{
		client: aplos_transport.NewRemoteProtobufRPCClientIfc(zkSession,
			ErgonServiceName, DefaultErgonPort, clusterUUID,
			userUUID, tenantUUID, 0 /* timeout */, userName, reqContext),
	}
}

// Use this function to instantiate new ergon client.
// Args:
//	serverIp   : IP address of server.
// 	serverPort : Port of ergon server.
// Return:
// 	Pointer to ProtobufRPCClient structure.
func NewErgonClient(
	serverIp string, serverPort uint16) *util_net.ProtobufRPCClient {
	return util_net.NewProtobufRPCClient(serverIp, serverPort)
}

// This method reads payload/ embedded value.
// Args:
//	proto: Pointer to proto of PayloadOrEmbeddedValue
// Returns:
//	Embedded value if successful, nil otherwise.
func ErgonReadPayloadOrEmbeded(proto *ergon_pr.PayloadOrEmbeddedValue) []byte {
	embeded := proto.GetEmbedded()
	if embeded != nil {
		glog.Info("Task response has embedded value")
		return embeded
	}
	// Todo: Add code for payload.
	return nil
}

// This method creates a task.
// Args:
// 	taskCreateArg : Proto with all parameters for creating task.
// Returns:
//	Pointer to TaskCreateRet proto.
//	error: nil if successful, error otherwise.
func (svc *ergonClient) TaskCreate(taskCreateArg *ergon_pr.TaskCreateArg) (
	*ergon_pr.TaskCreateRet, error) {

	rpcret := &ergon_pr.TaskCreateRet{}
	err := svc.sendMsg("TaskCreate", taskCreateArg, rpcret, 0, false)
	if err != nil {
		return nil, err
	}
	return rpcret, nil
}

// NOTE: Please use TaskGetMultiple to call the TaskGet RPC.
// This function is being used by some clients and hence has not
// been removed. All new users of TaskGet RPC should be using
// TaskGetMultiple as it takes TaskGetArg as an argument.
func (svc *ergonClient) TaskGet(uuid []byte) (*ergon_pr.Task, error) {

	taskGetArg := &ergon_pr.TaskGetArg{
		TaskUuidList: [][]byte{uuid},
	}
	taskGetRet := &ergon_pr.TaskGetRet{}

	err := svc.sendMsg("TaskGet", taskGetArg, taskGetRet, 0, false)

	if err != nil {
		return nil, err
	}

	taskList := taskGetRet.GetTaskList()
	if len(taskList) == 1 {
		return taskList[0], nil
	}
	// no task or more than 1 task.
	return nil, ErrInvalidState
}

// TaskGetMultiple requires an additional client side handling to ease the
// load on the Ergon Service. We make use of TaskGetIterator to batch a
// TaskGet RPC to smaller requests which are then combined to be sent to
// the end client.
func (svc *ergonClient) TaskGetMultiple(taskGetArg *ergon_pr.TaskGetArg) (
	*ergon_pr.TaskGetRet, error) {

	taskGetRet := &ergon_pr.TaskGetRet{}
	if !requiresClientSideBatching(taskGetArg) {
		err := svc.sendMsg("TaskGet", taskGetArg, taskGetRet, 0, false)
		if err != nil {
			return nil, err
		}
	} else {
		iter := NewTaskGetIterator(svc, taskGetArg, task_get_batch_size)
		for {
			taskGetRetBatch, err := iter.Next()
			if err != nil {
				return nil, err
			}
			if taskGetRetBatch == nil {
				break
			}
			taskGetRet.TaskList = append(taskGetRet.TaskList,
				taskGetRetBatch.TaskList...)
		}
	}
	return taskGetRet, nil
}

// Checks if the given TaskGet request requires Client side handling
// to batch the request
func requiresClientSideBatching(taskGetArg *ergon_pr.TaskGetArg) bool {
	return int64(len(taskGetArg.TaskUuidList)) > task_get_batch_size
}

// Get task error from MetaResponse, if any.
func (svc *ergonClient) GetTaskError(task *ergon_pr.Task) error {
	if task.GetStatus() == ergon_pr.Task_kFailed {
		resp := task.GetResponse()
		if resp == nil {
			return ErrInvalidRpcResponse.SetCause(errors.New("Null response in task"))
		}
		errCode := resp.GetErrorCode()
		return Errors[int(errCode)]
	}
	return nil
}

// Send ergon RPC.
func (svc *ergonClient) sendMsg(service string, request,
	response proto.Message, timeoutSecs int64, useBackoffWithRetryCount bool) error {
	var backoff *util_misc.ExponentialBackoff
	if useBackoffWithRetryCount {
		backoff = util_misc.NewExponentialBackoff(
			time.Duration(clientRetryInitialDelayMilliSecs)*time.Millisecond,
			time.Duration(clientRetryMaxDelaySecs)*time.Second,
			clientMaxNumRetries)
	} else {
		backoff = util_misc.NewExponentialBackoffWithTimeout(
			time.Duration(clientRetryInitialDelayMilliSecs)*time.Millisecond,
			time.Duration(clientRetryMaxDelaySecs)*time.Second,
			time.Duration(clientRetryTimeoutSecs)*time.Second)
	}
	retryCount := 0
	var err error

	for {
		err = svc.client.CallMethodSync(ErgonServiceName, service, request,
			response, timeoutSecs)
		if err == nil {
			return nil
		}
		appErr, ok := util_net.ExtractAppError(err)
		if ok {
			/* app error */
			errCode := appErr.GetErrorCode()
			if val, ok := Errors[errCode]; ok {
				if !ErrRetry.Equals(val) {
					return val
				}
			} else if bd := util_base.BuildInfo(util_base.GetBuildVersion); bd.IsDebugBuild() {
				glog.Fatalf("Unmapped response code: %d", errCode)
			} else {
				glog.Errorf("Unmapped response code: %d", errCode)
				return ErrUnmappedResponseCode
			}
		} else if !util_net.IsRpcTransportError(err) {
			/* Not app error and not transport error */
			return err
		}
		/* error is retriable. Continue with retry */
		waited := backoff.Backoff()
		if waited == util_misc.Stop {
			glog.Errorf("Timed out retrying Ergon RPC: %s. Request: %s",
				service, proto.MarshalTextString(request))
			break
		}
		retryCount += 1
		glog.Errorf("Error while executing %s: %s. Retry count: %d",
			service, err.Error(), retryCount)
	} // Continue the for loop.
	return err
}

func (svc *ergonClient) TaskUpdateBase(taskUpdateArg *ergon_pr.TaskUpdateArg) (
	*ergon_pr.TaskUpdateRet, error) {

	taskUpdateRet := &ergon_pr.TaskUpdateRet{}
	err := svc.sendMsg("TaskUpdate", taskUpdateArg, taskUpdateRet, 0, false)
	if err != nil {
		return nil, err
	}
	return taskUpdateRet, nil
}

func (svc *ergonClient) TaskUpdate(taskUpdateArg *ergon_pr.TaskUpdateArg) (
	*ergon_pr.TaskUpdateRet, error) {

	taskUpdateRet, err := svc.TaskUpdateBase(taskUpdateArg)

	var taskErr error = nil
	if taskUpdateRet != nil {
		task := taskUpdateRet.GetTask()
		if task != nil {
			taskErr = svc.GetTaskError(task)
		}
	}
	if (err != nil && ErrLogicalTimestampMismatch.Equals(err)) ||
		(taskErr != nil && ErrLogicalTimestampMismatch.Equals(taskErr)) {
		// timestamp mismatch - retry task with updated timestamp.
		task, err := svc.TaskGet(taskUpdateArg.GetUuid())
		if err != nil {
			return nil, ErrLogicalTimestampMismatch
		}
		taskUpdateArg.LogicalTimestamp = task.LogicalTimestamp
		taskUpdateRet = &ergon_pr.TaskUpdateRet{}
		err = svc.sendMsg("TaskUpdate", taskUpdateArg, taskUpdateRet, 0, false)
		if err != nil {
			glog.Infof("Failure in TaskUpdate %#v", err)
			return nil, err
		}
	} else if taskErr != nil {
		glog.Infof("Task Failure in TaskUpdate %#v", taskErr)
		return nil, taskErr
	} else if err != nil {
		glog.Infof("Error in TaskUpdate %#v", err)
		return nil, err
	}

	return taskUpdateRet, nil
}

func (svc *ergonClient) TaskCancel(taskCancelArg *ergon_pr.TaskCancelArg) (
	*ergon_pr.TaskCancelRet, error) {

	taskCancelRet := &ergon_pr.TaskCancelRet{}
	err := svc.sendMsg("TaskCancel", taskCancelArg, taskCancelRet, 0, false)
	if err != nil {
		return nil, err
	}
	return taskCancelRet, nil
}

// Poll on a task with time intervals and rpc timeout.
// Args:
//	taskUuid: UUID of task.
//	taskRet: Each task returns a proto which has to be interpreted by client
//	intervalSecs: Polling time interval in seconds.
//	timeoutSecs: RPC timeout in seconds for each poll call.
// Return:
//	Task status  : Status of task defined in ergon_types.proto
//	Metaresponse : Metaresponse transparent to ergon.
//	Error : nil if successful, error otherwise.
func (svc *ergonClient) TaskPollWait(taskUuid []byte, taskRet proto.Message,
	intervalSecs int64, timeoutSecs int64) (
	ergon_pr.Task_Status, *ergon_pr.MetaResponse, error) {
	rpcarg := &ergon_pr.TaskPollArg{}
	rpcret := &ergon_pr.TaskPollRet{}

	rpcarg.CompletedTasksFilter = append(
		rpcarg.CompletedTasksFilter, taskUuid)
	// To avoid RPC timing out before task poll timeout.
	taskPollTimeout := uint64(timeoutSecs)
	rpcarg.TimeoutSec = &taskPollTimeout
	for {
		err := svc.sendMsg("TaskPoll", rpcarg, rpcret, timeoutSecs+taskPollBufferSecs, true)
		if err != nil {
			if IsRpcTimeoutError(err) {
				glog.Infof("Received client timeout during poll: Retrying.")
				time.Sleep(time.Duration(intervalSecs) * time.Second)
				continue
			}
			glog.Infof("Failure in TaskPoll %#v", err)
			return ergon_pr.Task_kFailed, nil, err
		}
		if rpcret.GetTimedout() == true {
			glog.Info("Task poll wait timed out, retrying..")
			time.Sleep(time.Duration(intervalSecs) * time.Second)
			continue

		}
		for _, task := range rpcret.GetCompletedTasks() {

			response := task.GetResponse()
			taskStatus := task.GetStatus()

			if taskStatus == ergon_pr.Task_kFailed {
				glog.Infof("Task %x failed", taskUuid)
				return taskStatus, response, nil
			}
			if taskStatus == ergon_pr.Task_kAborted {
				glog.Infof("Task %x Aborted", taskUuid)
				return taskStatus, response, nil
			}

			if task.GetStatus() != ergon_pr.Task_kSucceeded {
				glog.Infof("Task %s returned status %v",
					taskUuid, taskStatus)
				continue
			}

			glog.Infof("Task %x successfully completed", taskUuid)
			if response.GetErrorCode() != 0 {
				//TODO: Need to handle Metaresponse in better
				// way. Returning invalid state right now.
				glog.Infof(
					"Task %x returned error %d %s", taskUuid, response.GetErrorCode(), response.GetErrorDetail())
				return task.GetStatus(), response, nil

			}
			// Look for embedded data if taskRet is not nil
			if taskRet != nil {
				data := ErgonReadPayloadOrEmbeded(response.GetRet())
				if data == nil {
					glog.Fatal("Failed to read embedded" +
						" data or payload")
				}
				err := proto.Unmarshal(data, taskRet)
				if err != nil {
					glog.Info("Failed to unmarshal proto with error ", err)
					return task.GetStatus(), nil, ErrUncaughtException.SetCause(err)
				}
			}
			return taskStatus, response, nil
		}
		glog.Info("Task poll sleeping")
		time.Sleep(time.Duration(intervalSecs) * time.Second)
	}
}

// This method polls on tasks.
// Args:
// 	taskPollArg : Proto with all parameters for task polling.
// Returns:
//	Pointer to TaskPollRet proto.
//	error: nil if successful, error otherwise.
func (svc *ergonClient) TaskPoll(taskPollArg *ergon_pr.TaskPollArg) (
	*ergon_pr.TaskPollRet, error) {

	rpcret := &ergon_pr.TaskPollRet{}
	err := svc.sendMsg(
		"TaskPoll", taskPollArg, rpcret, int64(taskPollArg.GetTimeoutSec()+taskPollBufferSecs), true)
	if err != nil {
		return nil, err
	}
	return rpcret, nil
}

// This method lists tasks.
// Args:
// 	taskListArg : Proto with all parameters for listing tasks
// Returns:
//	Pointer to TaskListRet proto.
//	error: nil if successful, error otherwise.
func (svc *ergonClient) TaskList(taskListArg *ergon_pr.TaskListArg) (
	*ergon_pr.TaskListRet, error) {

	rpcret := &ergon_pr.TaskListRet{}
	err := svc.sendMsg("TaskList", taskListArg, rpcret, 0, false)
	if err != nil {
		return nil, err
	}
	return rpcret, nil
}

func (svc *ergonClient) TaskCancelAny(taskCancelAnyArg *ergon_pr.TaskCancelAnyArg) (
	*ergon_pr.TaskCancelAnyRet, error) {

	taskCancelAnyRet := &ergon_pr.TaskCancelAnyRet{}
	err := svc.sendMsg("TaskCancelAny", taskCancelAnyArg, taskCancelAnyRet, 0, false)
	if err != nil {
		return nil, err
	}
	return taskCancelAnyRet, nil
}
