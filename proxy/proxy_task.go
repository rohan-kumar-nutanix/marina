/*
* Copyright (c) 2022 Nutanix Inc. All rights reserved.
*
* Authors: rajesh.battala@nutanix.com
*
* This is the implementation of the RPC proxy task for proxying Catalog RPC requests.
 */

package proxy

import (
	"context"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	log "k8s.io/klog/v2"

	catalogClient "github.com/nutanix-core/acs-aos-go/catalog/client"
	"github.com/nutanix-core/acs-aos-go/ergon"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/content-management-marina/common"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
	marinaProtos "github.com/nutanix-core/content-management-marina/protos/marina"
	"github.com/nutanix-core/content-management-marina/task/base"
	util "github.com/nutanix-core/content-management-marina/util"
	"github.com/nutanix-core/ntnx-api-utils-go/tracer"
)

// ProxyTask is the task to handle Catalog RPC proxy.
type ProxyTask struct {
	*base.MarinaBaseTask
}

// NewProxyTask creates a new proxy task for handling RPC proxy of Catalog requests.
func NewProxyTask(marinaBaseTask *base.MarinaBaseTask) *ProxyTask {
	return &ProxyTask{
		MarinaBaseTask: marinaBaseTask,
	}
}

func (t *ProxyTask) StartHook() error {
	rpcName := t.Proto().GetRequest().GetMethodName()
	log.Info("RPC method :", rpcName)
	t.Proto().InternalTask = proto.Bool(false)
	wal := t.Wal()
	serializationToken, err := t.serializationToken()
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to create proxy task serialization token: %s", err))
	}
	taskUuid, err := uuid4.New()
	if err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to create a new proxy task UUID: %s", err))
	}
	wal.ProxyTaskWal = &marinaProtos.ProxyTaskWal{
		SerializationToken: serializationToken,
		TaskUuid:           taskUuid.RawBytes(),
	}
	if err := t.SetWal(wal); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to set proxy task WAL: %s", err))
	}
	log.Infof("Created a new Marina proxy task for %s with Catalog "+
		"task UUID %s, and serialization ID %s.",
		t.Proto().GetRequest().GetMethodName(), taskUuid.String(), *serializationToken)
	return nil
}

// Run initializes for task recovery.
func (t *ProxyTask) Run() error {
	taskProto := t.Proto()
	// Get the RPC argument and return value based on the given RPC name.
	rpcName := t.Proto().GetRequest().GetMethodName()
	var ctxBytes []byte
	requestContext := t.Proto().GetRequestContext()
	if requestContext != nil {
		// TODO: Update the requestcontext package which has TraceContext.
		log.Infof("Setting ctxByes with requestContext TraceContent value")
	}
	// Start span for execution of proxy task
	span, ctx := tracer.StartSpanFromBytesCtx(ctxBytes,
		rpcName, context.Background())
	defer span.Finish()
	span.SetTag("task-uuid", uuid4.ToUuid4(t.Proto().GetUuid()).String())

	request, response, err := proxyConfig.getRequestResponseValues(rpcName)
	if err != nil {
		return t.taskError(ctx, marinaError.ErrInternal.SetCauseAndLog(
			fmt.Errorf("failed to get argument types for %s: %s", rpcName, err)))
	}
	// Set "parent_task_uuid" and "task_uuid" in the request.
	if err := t.updateRequestArgument(taskProto, request, rpcName); err != nil {
		return t.taskError(ctx, marinaError.ErrInternal.SetCauseAndLog(
			fmt.Errorf("failed to process request for %s: %s", rpcName, err)))
	}
	// Setup the RPC handler, and the arguments to the handler.
	reflectRpcName := reflect.ValueOf(rpcName)
	var reflectRpcHandler reflect.Value
	var args []reflect.Value
	reflectRpcHandler = reflect.ValueOf(catalogService.SendMsg)
	args = []reflect.Value{reflectRpcName, request, response}
	// Issue the RPC call to Catalog RPC server.
	if err := reflectRpcHandler.Call(args)[0]; !err.IsNil() {
		errObj := err.Elem().Interface().(error)
		log.Errorf("Failed to execute request %s: %v.", rpcName, errObj)
		return t.taskError(ctx, errObj)
	}
	// Get Catalog task UUID.
	taskUuid, err := t.getTaskUuid(response, rpcName)
	if err != nil {
		return t.taskError(ctx, marinaError.ErrInternal.SetCauseAndLog(
			fmt.Errorf("failed to process response for %s: %s", rpcName, err)))
	}

	// Poll for task completion.
	taskReturnValue, err := t.pollTaskCompletion(ctx, taskUuid, rpcName)
	if err != nil {
		return t.taskError(ctx, err)
	}
	log.Infof("Handled a %s proxy request with Catalog task UUID %s.",
		rpcName, taskUuid.String())

	t.Complete(t, taskReturnValue, nil)
	return nil
}

// updateRequestArgument updates the Catalog request argument by
// setting the parent task UUID to be the Marina task UUID, and set the
// "task_uuid" field of the Catalog request with the already created
// and WALed task UUID.
func (t *ProxyTask) updateRequestArgument(taskProto *ergon.Task,
	request reflect.Value, rpcName string) error {
	// Get the request and response arguments based on the given RPC name.
	// Unmarshal the RPC embedded value into RPC request argument.
	embeddedValue := taskProto.GetRequest().GetArg().GetEmbedded()
	if err := util.UnmarshalReflectBytes(embeddedValue, request); err != nil {
		return marinaError.ErrInternal.SetCause(
			fmt.Errorf("failed to unmarshal request for %s: %s", rpcName, err))
	}
	// Set request's parent task UUID to be Marina task UUID.
	taskUuidBytes := uuid4.ToUuid4(t.Proto().GetUuid()).RawBytes()
	// Note that all supported RPCs must have "parent_task_uuid" as a field in
	// the request.
	if err := util.SetReflectBytes(request, "ParentTaskUuid", taskUuidBytes); err != nil {
		return err
	}
	// If request has "task_uuid" field, set it with a new task UUID for
	// idempotence.
	requestTaskUuidField, err := proxyConfig.getRequestTaskUuidFieldName(rpcName)
	if err != nil {
		return marinaError.ErrInternal.SetCause(
			fmt.Errorf("failed to get request task UUID field for %s: %s", rpcName, err))
	}
	if requestTaskUuidField != "" {
		taskUuid := uuid4.ToUuid4(t.Wal().GetProxyTaskWal().GetTaskUuid())
		if taskUuid != nil {
			if err := util.SetReflectBytes(request, requestTaskUuidField,
				taskUuid.RawBytes()); err != nil {
				return err
			}
		}
	}
	return nil
}

// getTaskUuid gets the Catalog task UUID.
func (t *ProxyTask) getTaskUuid(response reflect.Value, rpcName string) (*uuid4.Uuid, error) {
	// Get the response task UUID.
	taskUuidFieldName, err := proxyConfig.getResponseTaskUuidFieldName(rpcName)
	if err != nil {
		return nil, marinaError.ErrInternal.SetCause(
			fmt.Errorf("failed to get RPC response task UUID field for %s: %s",
				rpcName, err))
	}
	taskUuid, err := util.GetReflectBytes(response, taskUuidFieldName)
	if err != nil {
		return nil, marinaError.ErrInternal.SetCause(
			fmt.Errorf("failed to get response task UUID for %s: %s",
				rpcName, err))
	}
	retTaskUuid := uuid4.ToUuid4(taskUuid)
	if retTaskUuid == nil {
		return nil, marinaError.ErrInternal.SetCause(
			fmt.Errorf("failed to parse response task UUID for %s", rpcName))
	}
	return retTaskUuid, nil
}

func (t *ProxyTask) taskError(ctx context.Context, err error) error {
	t.Complete(t, nil, err)
	return err
}

// pollTaskCompletion polls for a task completion. If the task failed, return Catalog RPC error.
func (t *ProxyTask) pollTaskCompletion(ctx context.Context, taskUuid *uuid4.Uuid,
	rpcName string) ([]byte, error) {
	taskUuidSet := map[string]bool{taskUuid.String(): true}
	retTaskList, err := t.PollAll(t, taskUuidSet)
	if err != nil {
		return nil, marinaError.ErrInternal.SetCause(
			fmt.Errorf("failed to poll for %s task %s: %s",
				rpcName, taskUuid.String(), err))
	}
	if len(retTaskList) != 1 {
		return nil, marinaError.ErrInternal.SetCause(
			fmt.Errorf("failed to poll for task completion for %s with UUID %s",
				rpcName, taskUuid.String()))
	}
	retTask := retTaskList[0]
	t.updateMarinaTask(ctx, retTask, rpcName)
	response := retTask.GetResponse()
	if retTask.GetStatus() != ergon.Task_kSucceeded {
		errCode := response.GetErrorCode()
		errDetail := response.GetErrorDetail()

		if proxyConfig.isSupportedRpc(rpcName) {
			log.Errorf("Catalog request %s failed with error code %d: %s.",
				rpcName, errCode, errDetail)
			return nil, catalogClient.CatalogError(errDetail, int(errCode))
		} else {
			return nil, marinaError.ErrInternal.SetCause(
				fmt.Errorf("request %s failed with error code %d: %s",
					rpcName, errCode, errDetail))
		}
	}
	return response.GetRet().GetEmbedded(), nil
}

// updateMarinaTask updates the Marina task to set the following fields based
// on the corresponding values as found in the Catalog task:
// "entity_list" and "start_time_usecs"
func (t *ProxyTask) updateMarinaTask(ctx context.Context, retTask *ergon.Task, rpcName string) {
	// Update the Marina task "entity_list" based on proxied task's entity list.
	taskProto := t.Proto()
	for _, entityID := range retTask.GetEntityList() {
		t.AddEntity(taskProto, entityID)
	}
	// Also update the Marina task start time.
	taskProto.StartTimeUsecs = proto.Uint64(retTask.GetStartTimeUsecs())
	t.Save(t)
}

// Enqueue added the task to serial executor.
func (t *ProxyTask) Enqueue() {
	se := common.Interfaces().SerialExecutor()
	se.SubmitJob(t)
}

// Execute runs a task in serial executor.
func (t *ProxyTask) Execute() {
	t.Resume(t)
}

// SerializationID identifies the field for request serialization.
func (t *ProxyTask) SerializationID() string {
	return t.Wal().GetProxyTaskWal().GetSerializationToken()
}

func (t *ProxyTask) serializationToken() (*string, error) {
	taskProto := t.Proto()
	rpcName := taskProto.GetRequest().GetMethodName()
	request, _, err := proxyConfig.getRequestResponseValues(rpcName)
	if err != nil {
		return nil, marinaError.ErrInternal.SetCauseAndLog(
			fmt.Errorf("failed to get argument types for %s: %s", rpcName, err))
	}
	embeddedValue := taskProto.GetRequest().Arg.Embedded
	if err := util.UnmarshalReflectBytes(embeddedValue, request); err != nil {
		return nil, marinaError.ErrInternal.SetCauseAndLog(
			fmt.Errorf("failed to unmarshal request for %s: %s", rpcName, err))
	}
	// The call to getRequestResponseValues() above has already verified that the
	// given RPC name is supported, and hence this call to
	// getSerializationTokenFunc() will not return any error.
	serializationTokenFunc, _ := proxyConfig.getSerializationTokenFunc(rpcName)
	id, err := serializationTokenFunc(request)
	if err != nil {
		return nil, err
	}
	// If there is no serialization ID, use random UUID as the request need
	// not be serialized
	if id == "" {
		log.Infof("%s request is not getting serialized.", rpcName)
		uuid, err := uuid4.New()
		if err != nil {
			return nil, marinaError.ErrInternal.SetCauseAndLog(
				fmt.Errorf("failed to create new UUID"))
		}
		id = uuid.String()
	}
	log.Infof("%s request serialization ID: %s.", rpcName, id)
	return &id, nil
}
