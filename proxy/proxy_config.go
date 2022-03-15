/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 * ProxyConfig for handling ProxyService requests to Catalog Service.
 */

package proxy

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/nutanix-core/acs-aos-go/catalog"
	marinaError "github.com/nutanix-core/content-management-marina/error"
	log "k8s.io/klog/v2"
)

// serviceConfig specifies fields of Catalog RPCs that are needed for
// proxying. We only include Catalog RPCs that are supported to be
// proxied by Marina Server.
type serviceConfig struct {
	// isSync indicates if the RPC is a sync RPC or an async RPC.
	isSync bool
	// requestTaskUuid refers to the field name for the task UUID as found in
	// the request. This is only applicable for async RPCs.
	requestTaskUuid string
	// responseTaskUuid refers to the field name for the task UUID as found in
	// the RPC return message. This is only applicable for async RPCs.
	responseTaskUuid string
	// requestType refers to the reflect type of the request argument, e.g.
	// CatalogItemCreateArg. This is auto filled in InitRpcServiceConfig().
	requestType reflect.Type
	// responseType refers to the reflect type of the response argument e.g.
	// CatalogItemCreateRet. This is auto filled in InitRpcServiceConfig().
	responseType reflect.Type
	// method that can be called to fetch the serialization token
	extractSerializationToken func(reflect.Value) (string, error)
}

const (
	// The field name in response message that stores parent task UUID.
	responseParentTaskUuid = "ParentTaskUuid"
	// The field name in response message that store the task UUID.
	responseTaskUuid = "TaskUuid"
)

var (
	// Fetch global_catalog_item_uuid from request arg to serialize CatalogItem requests based on
	// global_catalog_item_uuid. If global_catalog_item_uuid is not present, return an empty string.
	catalogItemSerializationToken = serializationTokenFunc("global_catalog_item_uuid")
	// Sync RPC configuration where only "isSync" is needed.
	syncRpc = serviceConfig{isSync: true}
	// CatalogItem async RPCs serialized by global_catalog_item_uuid, if exists.
	asyncCatalogItemRpc  = asyncRpc("TaskUuid", catalogItemSerializationToken)
	supportedCatalogRpcs = map[string]serviceConfig{
		// Catalog sync RPCs.
		"CatalogItemGet":                  syncRpc,
		"FileGet":                         syncRpc,
		"ImageGet":                        syncRpc,
		"CatalogPlacementPolicyGet":       syncRpc,
		"CatalogPlacementPolicyStatusGet": syncRpc,
		"VmTemplatesGet":                  syncRpc,
		"VmTemplateVersionsGet":           syncRpc,
		"CatalogRateLimitGet":             syncRpc,
		// TODO: As of now only 5.19 release Catalog RPC's are available in github repo.
		//"ResumableUploadStatusGet":        syncRpc,
		//"ImageViewGet":                    syncRpc,

		// Catalog async RPCs.
		"CatalogItemCreate": asyncCatalogItemRpc,
	}

	rpcServiceConfigOnce sync.Once
	// CatalogRpcNames is the list of supported PC Catalog RPC names.
	CatalogRpcNames = make([]string, len(supportedCatalogRpcs))
)

// serializationTokenFunc Construct a function that returns the first found field in an argument with a
// given list of potential argument fields.
func serializationTokenFunc(params ...string) func(reflect.Value) (string, error) {
	return func(arg reflect.Value) (string, error) {
		return "", nil
	}
}

func asyncRpc(requestTaskUuid string, serializationFunc func(reflect.Value) (
	string, error)) serviceConfig {
	return serviceConfig{
		isSync:                    false,
		requestTaskUuid:           requestTaskUuid,
		responseTaskUuid:          responseTaskUuid,
		extractSerializationToken: serializationFunc,
	}
}

func init() {
	log.Info("Initializing rpcServiceConfig to build supportedCatalogRpcs config map")
	InitRpcServiceConfig()
}

func InitRpcServiceConfig() {
	rpcServiceConfigOnce.Do(func() {
		// Initialize the map of method request and response types for all
		// supported Catalog RPC methods.
		var rpcNum int
		for rpcName, rpcService := range supportedCatalogRpcs {
			CatalogRpcNames = append(CatalogRpcNames, rpcName)
			// Replace RPC service specification with autofilled values.
			handlerMethod := reflect.ValueOf(&catalog.CatalogExternalRpcClient{}).MethodByName(rpcName)
			requestType := handlerMethod.Type().In(0).Elem()
			responseType := handlerMethod.Type().Out(0).Elem()
			validateFields(requestType, false, []string{rpcService.requestTaskUuid})
			validateFields(responseType, true, []string{rpcService.responseTaskUuid})
			supportedCatalogRpcs[rpcName] = serviceConfig{
				isSync:                    rpcService.isSync,
				requestTaskUuid:           rpcService.requestTaskUuid,
				responseTaskUuid:          rpcService.responseTaskUuid,
				requestType:               requestType,
				responseType:              responseType,
				extractSerializationToken: rpcService.extractSerializationToken,
			}
			rpcNum++
		}
		log.Infof("Initialized %d supported PC Catalog RPCs for proxying.", rpcNum)
	})
}

// validateFields validates if a reflect type has all the specified fields.
func validateFields(reflectType reflect.Type, isSettable bool,
	fieldNames []string) bool {
	if isSettable {
		// For response type, also validate the required field "parent_task_uuid".
		fieldNames = append(fieldNames, responseParentTaskUuid)
	}
	for _, fieldName := range fieldNames {
		reflectValue := reflect.New(reflectType).Elem()
		field := reflectValue.FieldByName(fieldName)
		if !field.IsValid() {
			return false
		}
		if isSettable && !field.CanSet() {
			return false
		}
	}
	return true
}

// CatalogProxyConfigInterface Helper methods for getting proxy configurations.
// and its is mainly for mocking out methods in unit tests.
type CatalogProxyConfigInterface interface {
	getRequestResponseValues(rpcName string) (reflect.Value, reflect.Value, error)
}

type catalogProxyConfigUtil struct{}

// A singleton to access proxy config utility functions.
var proxyConfig = new(catalogProxyConfigUtil)

// getRequestResponseValues returns the RPC method request and response reflect
// values associated with the given RPC method name.
func (*catalogProxyConfigUtil) getRequestResponseValues(rpcName string) (
	reflect.Value, reflect.Value, error) {
	if rpcService, found := supportedCatalogRpcs[rpcName]; found {
		if rpcService.requestType == nil {
			return reflect.Value{}, reflect.Value{}, marinaError.ErrMarinaInternal.SetCause(
				fmt.Errorf("RPC method %s has nil request type", rpcName))
		}
		requestValue := reflect.New(rpcService.requestType)
		if rpcService.responseType == nil {
			return reflect.Value{}, reflect.Value{}, marinaError.ErrMarinaInternal.SetCause(
				fmt.Errorf("RPC method %s has nil response type", rpcName))
		}
		responseValue := reflect.New(rpcService.responseType)
		return requestValue, responseValue, nil
	}
	return reflect.Value{}, reflect.Value{}, marinaError.ErrMarinaInternal.SetCause(
		fmt.Errorf("unsupported RPC method: %s", rpcName))
}
