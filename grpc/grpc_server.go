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

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/tracer"

	"google.golang.org/grpc"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"

	"github.com/nutanix-core/content-management-marina/grpc/services"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
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
	var unaryInterceptors []grpc.UnaryServerInterceptor
	var streamInterceptors []grpc.StreamServerInterceptor

	s := &ServerImpl{port: port}
	// TODO: Enable support for mTLS based on usecase.

	// Enabling tracing OpenTelemetry.
	unaryTracer, streamTracer := tracer.GrpcServerTraceOptions()

	unaryInterceptors = append(unaryInterceptors, grpcPrometheus.UnaryServerInterceptor)
	streamInterceptors = append(streamInterceptors, grpcPrometheus.StreamServerInterceptor)
	if unaryTracer != nil {
		log.Infof("Appending Unary Tracer opts : %v", unaryTracer)
		unaryInterceptors = append(unaryInterceptors, unaryTracer)
	}
	if streamTracer != nil {
		log.Infof("Appending Stream Tracer opts : %v", streamTracer)
		streamInterceptors = append(streamInterceptors, streamTracer)
	}
	// Enabling grpc Sever with Prometheus and OpenTelemetry.
	serverParams := []grpc.ServerOption{
		grpcMiddleware.WithUnaryServerChain(unaryInterceptors...),
		grpcMiddleware.WithStreamServerChain(streamInterceptors...),
	}

	log.Infof("Server opts %v", serverOpts)
	s.gserver = grpc.NewServer(serverParams...)
	s.registerServices(s.gserver)
	grpcPrometheus.Register(s.gserver)

	return s
}

// registerServices is a central place for the grpc services that need to be
// registered with the server before it is started.
func (server *ServerImpl) registerServices(s *grpc.Server) {
	log.Info("Registering services...")
	marinaIfc.RegisterMarinaServer(server.gserver, &services.MarinaServer{})
	content.RegisterWarehouseServiceServer(server.gserver, &services.WarehouseServer{})
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
