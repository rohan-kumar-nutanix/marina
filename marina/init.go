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
	"io"
	"net"

	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/tracer"
	"github.com/nutanix-core/content-management-marina/common"
	"github.com/nutanix-core/content-management-marina/interface/external"
	internal "github.com/nutanix-core/content-management-marina/interface/local"
	utils "github.com/nutanix-core/content-management-marina/util"
)

var traceProvider io.Closer

func initMarina() {
	external.InitSingletonService()
	internal.InitSingletonService()
}

// initFlags initialize gflags for Marina.
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

// TODO: ENG-538562: Use OpenTelemetry instead of OpenTracing
func initOpenTelemetryTracing() {
	var err error
	traceProvider = tracer.InitTracer(utils.ServiceName)
	if err != nil {
		log.Errorf("Error while initializing tracer: %v ", err.Error())

	} else if traceProvider != nil {
		log.Infof("OpenTelemetry Tracer got initialized %v", traceProvider)

	} else {
		log.Info("OpenTelemetry Tracer is disabled")
	}
}
