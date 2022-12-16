/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 * Marina service specific constants.
 */

package utils

const (
	ServiceName                = "Marina"
	RpcTimeoutSec              = 60
	IdfRpcTimeOut              = 10
	CatalogServiceName         = "nutanix.catalog_pc.catalogpcexternalrpcservice" //"nutanix.catalog.CatalogRpcService"
	CatalogLegacyServiceName   = "nutanix.catalog.CatalogRpcService"
	CatalogInternalServiceName = "nutanix.catalog_pc.CatalogPcRpcService" //nutanix.catalog_pc.catalogpcrpcservice
	PcFQDN                     = "pcip"
)

var HostAddr string // will be initialised dynamically with PC IP during init.
