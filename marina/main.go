/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 * Marina service.
 */

package main

import (
	"flag"
	"os"

	"github.com/nutanix-core/ntnx-api-utils-go/tracer"
	log "k8s.io/klog/v2"
)

func main() {
	flag.Parse()
	trace_provider, err := tracer.InitTracer("marina")
	if err != nil {
		log.Errorf("Error while initilizing tracer: %v ", err.Error())
	} else {
		log.Info("Tracer got initilized")
	}

	defer trace_provider.Shutdown(nil)
	runRPCServer()
	log.SetOutput(os.Stdout)
	waitGroup.Wait()
}
