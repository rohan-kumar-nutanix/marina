/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Authors: rajesh.battala@nutanix.com

 * This file implements RPC services to serve Catalog RPC requests.
 * The main entry point to here is through the call to NewRpcService(). The
 * purpose of this file is for specifying the Catalog RPC service
 * to be served by Marina.
 */

package proxy

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/ergon"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/net"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	slbufsNet "github.com/nutanix-core/acs-aos-go/nutanix/util-slbufs/util/sl_bufs/net"
	"github.com/nutanix-core/content-management-marina/common"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
	"github.com/nutanix-core/content-management-marina/task/base"
	marinaUtil "github.com/nutanix-core/content-management-marina/util"
	clientUtil "github.com/nutanix-core/content-management-marina/util/catalog/client"
)

var catalogService *clientUtil.Catalog

// NewRpcService creates a new RPC service for Catalog. This
// service is used for configuring Marina RPC server to serve Catalog RPC requests.
// "Desc" is a map from RPC names to RPC handlers
// (func(rpc *ProtobufRpc, handlers interface{}) error). Upon receiving a RPC
// request, a ProtobufRpc is constructed, and the RPC handler as specified in
// the above map is invoked (using RPCRequestHandler() in net) with the
// first argument being the constructed RPC request, and the second argument is
// the function as specified in "Impl". Note that we don't need to set Impl
// which is passed into rpcHandler(), as we are able to get RPC request and
// response argument types statically.
// Init Catalog Service Handler Instance.
func NewRpcService(svcName string) *net.Service {
	// TODO move it to singleton.
	catalogService = clientUtil.NewCatalogService(marinaUtil.HostAddr, uint16(*common.CatalogPort))
	// Set up the RPC services.
	return &net.Service{
		Desc: getRpcSvcDesc(svcName),
	}
}

// getRpcSvcDesc returns a service description for Catalog RPCs
// served by Marina. It is registering rpcHandler() as the RPC handler for each
// RPC exported by Marina.
func getRpcSvcDesc(svcName string) *net.ServiceDesc {
	svcDesc := net.ServiceDesc{
		Name:    svcName,
		Methods: map[string]net.ServiceMethodFn{},
	}
	var rpcNames []string
	if svcName == marinaUtil.CatalogServiceName ||
		svcName == marinaUtil.CatalogLegacyServiceName ||
		svcName == marinaUtil.CatalogInternalServiceName {
		rpcNames = CatalogRpcNames
	} else {
		log.Fatalln("Unsupported RPC service: ", svcName)
	}
	// Set the method to handle Catalog RPC requests.
	for _, rpcName := range rpcNames {
		svcDesc.Methods[rpcName] = rpcHandler
	}
	log.Infof("Created a new RPC service for %s.", svcName)
	return &svcDesc
}

// rpcHandler handles a RPC request.
func rpcHandler(rpc *net.ProtobufRpc, handlers interface{}) error {
	// Retrieve the request and response argument type based on request name.
	rpcName := *rpc.RequestHeader.MethodName
	request, response, err := proxyConfig.getRequestResponseValues(rpcName)
	if err != nil {
		return marinaError.ErrInternal.SetCauseAndLog(
			fmt.Errorf("failed to get request/response type: %v", err))
	}

	if proxyConfig.isSyncRpc(rpcName) {
		// Handle the RPC request in sync.
		if err := syncRpcHandler(rpc, request, response); err != nil {
			return marinaError.ErrInvalidArgument.SetCauseAndLog(
				fmt.Errorf("failed to handle sync RPC %s: %s", rpcName, err))
		}
	} else {
		// Handle the RPC request in ASync.
		if err := asyncRpcHandler(rpc, request, response); err != nil {
			return marinaError.ErrInvalidArgument.SetCauseAndLog(fmt.Errorf(
				"failed to create proxy task for %s: %s", rpcName, err))
		}
	}
	// Serialize and sets the response into the RPC context.
	if err := marshalResponse(rpc, response); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("failed to serialize RPC response: %v", err))
	}

	return nil
}

// processRpcError handles an error as returned by RPC call to Catalog.
func processRpcError(rpc *net.ProtobufRpc, err reflect.Value,
	marinaErr *marinaError.MarinaError) error {
	// Set app error if available.
	errObj := err.Elem().Interface().(error)
	if appErr, ok := net.ExtractAppError(errObj); ok {
		rpc.SetAppError(int32(appErr.GetErrorCode()), appErr.GetErrorDetail())
		return appErr
	}

	if marinaErr == nil {
		return errObj
	}
	// Set the provided Catalog error as the RPC error.
	rpc.SetAppError(int32(marinaErr.GetErrorCode()), marinaErr.GetErrorDetail())
	return marinaErr
}

// unmarshalRequestPayload unmarshal an RPC request payload into the specified
// request argument for the given RPC.
func unmarshalRequestPayload(rpc *net.ProtobufRpc, request reflect.Value,
	rpcName string) error {
	if rpc.RequestPayload == nil {
		return marinaError.ErrInvalidArgument.SetCause(
			fmt.Errorf("request %s must contain a payload", rpcName))
	}
	if err := marinaUtil.UnmarshalReflectBytes(rpc.RequestPayload, request); err != nil {
		return marinaError.ErrInvalidArgument.SetCause(
			fmt.Errorf("failed to unmarshal request %s", rpcName))
	}
	return nil
}

// marshalResponse marshals and sets the response into the RPC context.
func marshalResponse(rpc *net.ProtobufRpc, response reflect.Value) error {
	responseBytes, err := marinaUtil.MarshalReflectBytes(response)
	if err != nil {
		return marinaError.ErrInternal.SetCause(
			fmt.Errorf("failed to marshal response for %s", *rpc.RequestHeader.MethodName))
	}
	rpc.ResponsePayload = responseBytes
	return nil
}

// syncRpcHandler handles a sync RPC request by directly invoking the sync RPC call to Catalog.
func syncRpcHandler(rpc *net.ProtobufRpc, request reflect.Value,
	response reflect.Value) error {

	rpcName := rpc.RequestHeader.GetMethodName()
	// Unmarshal the RPC payload into request argument for the given RPC name.
	if err := unmarshalRequestPayload(rpc, request, rpcName); err != nil {
		return err
	}
	// Get the RPC handler (e.g. CatalogItemGet) Catalog service, and set up
	// args (e.g. CatalogItemGetArg and CatalogItemGetRet).
	reflectName := reflect.ValueOf(rpcName)
	var rpcHandler reflect.Value
	var args []reflect.Value

	// Make request to Catalog service
	rpcHandler = reflect.ValueOf(catalogService.SendMsg)
	args = []reflect.Value{reflectName, request, response}

	// Make the sync RPC call to Catalog.
	if err := rpcHandler.Call(args)[0]; !err.IsNil() {
		return processRpcError(rpc, err, nil)
	}
	return nil
}

// asyncRpcHandler creates a new proxy task for handling async Catalog RPCs.
func asyncRpcHandler(rpc *net.ProtobufRpc, request reflect.Value,
	response reflect.Value) error {
	//ctx := context.TODO()
	rpcName := rpc.RequestHeader.GetMethodName()
	// Create a new task proto for the task to proxy RPCs.
	taskProto := &ergon.Task{
		Request: &ergon.MetaRequest{
			MethodName: proto.String(rpcName),
			Arg: &ergon.PayloadOrEmbeddedValue{
				Embedded: rpc.RequestPayload,
			},
		},
	}
	taskProto.RequestContext = &slbufsNet.RpcRequestContext{}
	// If request specifies a task UUID for idempotent, set it in Marina task.
	if err := maybeSetRequestTaskUuid(rpcName, rpc.RequestPayload, taskProto); err != nil {
		return marinaError.ErrInternal.SetCause(
			fmt.Errorf("failed to set request task UUID for %s: %s", rpcName, err))
	}
	// Create a new Marina proxy task. Start() will save the task to Ergon.
	proxyTask := NewProxyTask(base.NewMarinaBaseTask(taskProto))
	var operationType string
	operationType = rpcName
	err := proxyTask.Start(proxyTask, proto.String(marinaUtil.ServiceName), &operationType)
	if err != nil {
		return marinaError.ErrInternal.SetCause(
			fmt.Errorf("failed to start task for %s: %s", rpcName, err))
	}
	// Get Ergon task UUID and set it in response.
	taskUUID := uuid4.ToUuid4(proxyTask.Proto().GetUuid())
	if taskUUID == nil {
		return marinaError.ErrInternal.SetCause(
			fmt.Errorf("invalid task UUID for %s", rpcName))
	}
	taskUUIDFieldName, err := proxyConfig.getResponseTaskUuidFieldName(rpcName)
	if err != nil {
		return marinaError.ErrInternal.SetCause(
			fmt.Errorf("failed to get RPC response task UUID field for %s: %s", rpcName, err))
	}
	if err := marinaUtil.SetReflectBytes(response, taskUUIDFieldName, taskUUID.RawBytes()); err != nil {
		return marinaError.ErrInternal.SetCause(
			fmt.Errorf("%s response missing task UUID field: %s", rpcName, err))
	}
	return nil
}

func maybeSetRequestTaskUuid(rpcName string, requestPayload []byte,
	taskProto *ergon.Task) error {
	// Retrieve the task UUID field, if exists, in the request.
	requestTaskUuidField, err := proxyConfig.getRequestTaskUuidFieldName(rpcName)
	if err != nil {
		return marinaError.ErrInternal.SetCause(
			fmt.Errorf("failed to get request task UUID field for %s: %s", rpcName, err))
	}
	if requestTaskUuidField == "" {
		// Request doesn't allow task UUID, so no need to set it.
		return nil
	}
	// Get request reflect value to unmarshal the request payload.
	request, _, err := proxyConfig.getRequestResponseValues(rpcName)
	if err != nil {
		return marinaError.ErrInternal.SetCause(
			fmt.Errorf("failed to get argument types for %s: %s", rpcName, err))
	}
	// Unmarshal the RPC request payload into request argument.
	if err := marinaUtil.UnmarshalReflectBytes(requestPayload, request); err != nil {
		return marinaError.ErrInternal.SetCause(
			fmt.Errorf("failed to unmarshal request for %s: %s", rpcName, err))
	}
	requestTaskUuid, err := marinaUtil.GetReflectBytes(request, requestTaskUuidField)
	if err != nil {
		return marinaError.ErrInternal.SetCause(
			fmt.Errorf("failed to get response task UUID for %s: %s", rpcName, err))
	}
	// Set request task UUID as Marina task UUID, if specified.
	if len(requestTaskUuid) > 0 {
		taskUUID := uuid4.ToUuid4(requestTaskUuid)
		taskProto.Uuid = taskUUID.RawBytes()
	}
	return nil
}
