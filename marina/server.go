/*
* Copyright (c) 2021 Nutanix Inc. All rights reserved.
*
* Author: rajesh.battala@nutanix.com
*
* Run gRPC server to handle Marina RPC requests.
*
 */

package main

import (
	"flag"
	"net"
	"sync"

	marinaGRPCServer "github.com/nutanix-core/content-management-marina/grpc"
	util "github.com/nutanix-core/content-management-marina/util"
	log "k8s.io/klog/v2"
)

var GRPC_PORT, PROMETHEUS_PORT *uint64

// Marina service configurations.
type marinaServiceConfig struct {
	grpcServer marinaGRPCServer.Server // Marina gRPC server.
}

var (
	waitGroup sync.WaitGroup
	svcConfig = &marinaServiceConfig{}
)

func init() {
	// flag.Parse()
	GRPC_PORT = flag.Uint64("grpc_port", 9200, "GRPC port")
	PROMETHEUS_PORT = flag.Uint64("prometheus_port", 9201, "Prometheus port")

	// Fetch the host IP address for service operations.
	fqdn := "pcip"
	pcIP, err := net.LookupIP(fqdn)
	if err != nil || len(pcIP) == 0 {
		log.Fatal("Could not fetch host IP to initialise service")
	}

	util.HostAddr = pcIP[0].String()
	log.Info("HostAddr : ", util.HostAddr)
}

func runRPCServer() {
	waitGroup.Add(1)
	defer waitGroup.Done()

	// Start the Prometheus handler.
	log.Infof("Starting Prometheus handler.")
	metric := Init(int(*PROMETHEUS_PORT))
	log.Infof("Metric handler details %v", metric)

	// Start the Marina gRPC server.
	startGRPCServer()
}

func startGRPCServer() {
	log.Infof("Starting GRPC on port %v", *GRPC_PORT)
	svcConfig.grpcServer = marinaGRPCServer.NewServer(*GRPC_PORT)
	svcConfig.grpcServer.Start(&waitGroup)
}
