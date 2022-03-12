/*
 Copyright (c) 2022 Nutanix Inc. All rights reserved.

 Authors: rajesh.battala@nutanix.com

 This file implements RPC services to serve Catalog RPC requests.
 The main entry point to here is through the call to NewRpcService(). The
 purpose of this file is for specifying the Catalog RPC service
 to be served by Marina.
*/

package proxy

import (
	"fmt"
	"github.com/nutanix-core/content-management-marina/common"
	"reflect"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/net"
	marinaError "github.com/nutanix-core/content-management-marina/error"
	marinaUtil "github.com/nutanix-core/content-management-marina/util"
	clientUtil "github.com/nutanix-core/content-management-marina/util/catalog/client"
	log "k8s.io/klog/v2"
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
	if svcName == marinaUtil.CatalogServiceName {
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

	// Handle the RPC request in sync.
	if err := syncRpcHandler(rpc, request, response); err != nil {
		return marinaError.ErrInvalidArgument.SetCauseAndLog(
			fmt.Errorf("Failed to handle sync RPC %s: %s.", rpcName, err))
	}
	// Serialize and sets the response into the RPC context.
	if err := marshalResponse(rpc, response); err != nil {
		return marinaError.ErrMarinaInternal.SetCauseAndLog(
			fmt.Errorf("Failed to serialize RPC response: %v.", err))
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

// unmarshalRequestPayload unmarshal a RPC request payload into the specified
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

// syncRpcHandler handles a sync RPC request by directly invoking the sync RPC
// call to Catalog.
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
