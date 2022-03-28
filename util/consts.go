/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 * Marina service specific constants.
 */

package utils

const (
	DefaultLogDir              = "/home/nutanix/data/logs"
	ServiceName                = "Marina"
	LocalhostAddr              = "127.0.0.1"
	RpcTimeoutSec              = 60
	ZkHostPort                 = 9876
	ZkTimeOut                  = 10
	IdfRpcTimeOut              = 10
	CatalogServiceName         = "nutanix.catalog_pc.catalogpcexternalrpcservice" //"nutanix.catalog.CatalogRpcService"
	CatalogLegacyServiceName   = "nutanix.catalog.CatalogRpcService"
	CatalogInternalServiceName = "nutanix.catalog_pc.CatalogPcRpcService" //nutanix.catalog_pc.catalogpcrpcservice
	PcFQDN                     = "pcip"
)

var HostAddr string // init it dynamically with PC IP during init.
