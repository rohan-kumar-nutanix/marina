/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 * Marina service.
 */

package main

import (
	"os"

	log "k8s.io/klog/v2"
)

func main() {
	initFlags()
	initMarina()
	initHostIP()
	initOpenTelemetryTracing()
	if traceProvider != nil {
		defer traceProvider.Shutdown(nil)
	}
	runRPCServer()
	log.SetOutput(os.Stdout)
	waitGroup.Wait()
}
