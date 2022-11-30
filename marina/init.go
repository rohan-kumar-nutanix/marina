/*
Copyright (c) 2022 Nutanix Inc. All rights reserved.

Author: rajesh.battala@nutanix.com

Initialize Marian, including:
 - gflags
 - OpenTelemetry Tracing.
 - HostAddress - Resolving PCIP.

*/

package main

import (
	"github.com/nutanix-core/content-management-marina/common"
	utils "github.com/nutanix-core/content-management-marina/util"
	"github.com/nutanix-core/ntnx-api-utils-go/tracer"
	otelSdkTrace "go.opentelemetry.io/otel/sdk/trace"
	log "k8s.io/klog/v2"
	"net"
)

var traceProvider *otelSdkTrace.TracerProvider

// initGflags initialize gflags for Marina.
func initFlags() {
	common.FlagsInit()
}

func initHostIP() {
	pcIP, err := net.LookupIP(utils.PcFQDN)
	if err != nil || len(pcIP) == 0 {
		log.Fatal("Could not fetch host IP to initialise service")
	}

	utils.HostAddr = pcIP[0].String()
	log.Info("Setting HostAddr to : ", utils.HostAddr)
}
func initOpenTelemetryTracing() {
	traceProvider, err := tracer.InitTracer(utils.ServiceName)
	if traceProvider == nil {
		log.Errorf("Error while initializing tracer: %v ", err.Error())
	} else {
		log.Infof("OpenTelemetry Tracer got initialized %v", traceProvider)
	}
}
