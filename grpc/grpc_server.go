/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 * The underlying grpc server that exports the various services. Services may
 * be added in the registerServices() implementation.
 *
 */

package grpc

import (
	"net"

	"strconv"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	services "github.com/nutanix-core/content-management-marina/grpc/services"
	marina_pb "github.com/nutanix-core/content-management-marina/protos/marina"
	"github.com/nutanix-core/ntnx-api-utils-go/tracer"

	"google.golang.org/grpc"
	log "k8s.io/klog/v2"
)

// Server encapsulates the grpc server.
type Server interface {
	// Start the server. waitGroup will be used to track execution of the server.
	Start(waitGroup *sync.WaitGroup)
	Stop()
}

type ServerImpl struct {
	port     uint64
	listener net.Listener
	gserver  *grpc.Server
}

// NewServer creates a new GRPC server that services can be exported with.
// The connections are secured by mTLS. Errors are fatal.
func NewServer(port uint64) (server Server) {
	var serverOpts []grpc.ServerOption
	s := &ServerImpl{port: port}
	// TODO: Enable support for mTLS based on usecase.

	// Enabling tracing OpenTelemetry.
	unaryTracer, streamTracer := tracer.GrpcServerTraceOptions()
	log.Infof("Tracer opts : %v %v", unaryTracer, streamTracer)

	// Enabling grpc Sever with Prometheus and OpenTelemetry.
	serverParams := []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(grpc_prometheus.UnaryServerInterceptor, unaryTracer),
		grpc_middleware.WithStreamServerChain(grpc_prometheus.StreamServerInterceptor, streamTracer),
	}

	log.Infof("Server opts %v", serverOpts)
	s.gserver = grpc.NewServer(serverParams...)
	s.registerServices(s.gserver)
	grpc_prometheus.Register(s.gserver)

	return s
}

// registerServices is a central place for the grpc services that need to be
// registered with the server before it is started.
func (server *ServerImpl) registerServices(s *grpc.Server) {
	log.Info("Registering services...")
	marina_pb.RegisterMarinaServer(server.gserver, &services.MarinaServer{})
}

// Start listening and serve. Errors are fatal (todo).
func (server *ServerImpl) Start(waitGroup *sync.WaitGroup) {
	addr := ":" + strconv.FormatUint(server.port, 10)
	log.Info("Starting Marina gRPC server on %s.", addr)
	var err error
	server.listener, err = net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v.", err)
	}
	log.Infof("Marina gRPC server listening on %s.", addr)

	waitGroup.Add(1)
	go func() {
		if err := server.gserver.Serve(server.listener); err != nil {
			log.Fatalf("Failed to serve: %v.", err)
		}
		waitGroup.Done()
	}()
}

// Stop the server.
func (server *ServerImpl) Stop() {
	if server.gserver != nil {
		server.gserver.GracefulStop()
		server.gserver = nil
		server.listener = nil
	}
}
