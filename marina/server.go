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
	"github.com/nutanix-core/content-management-marina/common"
	"net/http"

	"strconv"
	"sync"
	"syscall"

	"github.com/golang/glog"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/net"
	marinaGRPCServer "github.com/nutanix-core/content-management-marina/grpc"
	"github.com/nutanix-core/content-management-marina/proxy"
	util "github.com/nutanix-core/content-management-marina/util"
	log "k8s.io/klog/v2"
)

// Marina service configurations.
type marinaServiceConfig struct {
	grpcServer marinaGRPCServer.Server    // Marina gRPC server.
	rpcServer  *net.ProtobufRPCHTTPServer // Marina RPC server.
}

var (
	mutex     *sync.Mutex = &sync.Mutex{}
	waitGroup sync.WaitGroup
	svcConfig = &marinaServiceConfig{}
)

func runRPCServer() {
	mutex.Lock()
	waitGroup.Add(1)
	defer waitGroup.Done()

	// TODO: Handle gRPC server restart.
	/* 	if svcConfig.rpcServer != nil {
		restartRPCServer(isLeader)
		mutex.Unlock()
		return
	} */

	// Start the Prometheus handler.
	log.Infof("Starting Prometheus handler.")
	InitPrometheus(int(*common.PrometheusPort))
	startRPCServer(true)
	// Start the Marina gRPC server.
	startGRPCServer()
}

// startRPCServer starts a new RPC server on this node by registering supported
// services for PC Catalog. Note that this function is
// executed with the mutex lock, so no need to issue that lock again.
func startRPCServer(isLeader bool) {
	log.Info("Starting.. ProtobufRPCHTTP server... on Port : ", *common.MarinaRPCPort)
	// Create a new HTTP server for Narsil RPC server.
	httpServer := net.NewHTTPServer("0.0.0.0",
		strconv.FormatUint(*common.MarinaRPCPort, 10))
	httpServer.RegisterDebugHandlerFunc("exit", exitHTTPServerHandler)
	svcConfig.rpcServer = net.NewProtobufRPCHTTPServer(httpServer)

	// Export Catalog RPC service.
	catalogSvc := proxy.NewRpcService(util.CatalogServiceName)
	if !svcConfig.rpcServer.ExportService(catalogSvc) {
		glog.Fatalf("Failed to register Catalog RPC services.")
	}

	// Start the RPC server.
	go func() {
		glog.Fatal(svcConfig.rpcServer.StartProtobufRPCHTTPServer())
	}()
}

func exitHTTPServerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{\"acknowledged\":true}"))
	go terminate()
}

func terminate() {
	syscall.Kill(syscall.Getpid(), syscall.SIGQUIT)
}

func startGRPCServer() {
	log.Infof("Starting GRPC on port %v", *common.MarinaGRPCPort)
	svcConfig.grpcServer = marinaGRPCServer.NewServer(*common.MarinaGRPCPort)
	svcConfig.grpcServer.Start(&waitGroup)
}
