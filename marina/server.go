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
	"net/http"
	"strconv"
	"sync"
	"syscall"

	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/net"

	"github.com/nutanix-core/content-management-marina/common"
	marinaGRPCServer "github.com/nutanix-core/content-management-marina/grpc"
	"github.com/nutanix-core/content-management-marina/proxy"
	"github.com/nutanix-core/content-management-marina/task"
	util "github.com/nutanix-core/content-management-marina/util"
)

// Marina service configurations.
type marinaServiceConfig struct {
	grpcServer  marinaGRPCServer.Server    // Marina gRPC server.
	rpcServer   *net.ProtobufRPCHTTPServer // Marina RPC server.
	taskManager *task.MarinaTaskManager    // Marina task manager.
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
	// Create a new HTTP server for Marina RPC server.
	httpServer := net.NewHTTPServer("0.0.0.0",
		strconv.FormatUint(*common.MarinaRPCPort, 10))
	httpServer.RegisterDebugHandlerFunc("exit", exitHTTPServerHandler)
	svcConfig.rpcServer = net.NewProtobufRPCHTTPServer(httpServer)

	// Export Catalog RPC service.
	catalogSvc := proxy.NewRpcService(util.CatalogServiceName)
	if !svcConfig.rpcServer.ExportService(catalogSvc) {
		log.Fatalf("Failed to register Catalog RPC services.", catalogSvc)
	}

	// Export Catalog RPC service.
	catalogOldSvc := proxy.NewRpcService(util.CatalogLegacyServiceName)
	if !svcConfig.rpcServer.ExportService(catalogOldSvc) {
		log.Fatalf("Failed to register Catalog RPC services.", catalogOldSvc)
	}

	// Export Catalog RPC service.
	catalogInternalSvc := proxy.NewRpcService(util.CatalogInternalServiceName)
	if !svcConfig.rpcServer.ExportService(catalogInternalSvc) {
		log.Fatalf("Failed to register Catalog RPC services.", catalogInternalSvc)
	}

	// TODO: Start TaskDispatcher on Leader node only.
	// Start a new task manager to dispatch Marina tasks on leader.
	log.Infof("Starting Marina TaskDispatcher...")
	svcConfig.taskManager = task.NewMarinaTaskManager()
	svcConfig.taskManager.StartTaskDispatcher()

	// Start the RPC server.
	go func() {
		log.Fatal(svcConfig.rpcServer.StartProtobufRPCHTTPServer())
	}()
}

func exitHTTPServerHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{\"acknowledged\":true}"))
	go terminate()
}

func terminate() {
	err := syscall.Kill(syscall.Getpid(), syscall.SIGQUIT)
	if err != nil {
		log.Errorf("unable to terminate the process/service.")
	}
}

func startGRPCServer() {
	log.Infof("Starting GRPC on port %v", *common.MarinaGRPCPort)
	svcConfig.grpcServer = marinaGRPCServer.NewServer(*common.MarinaGRPCPort)
	svcConfig.grpcServer.Start(&waitGroup)

	// TODO: Enable mTLS mode for gRPC server. Security Team suggestion.
	// Check with msp team for certs generation, istio or genesis .

}
