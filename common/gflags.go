/*
Copyright (c) 2022 Nutanix Inc. All rights reserved.

Author: rajesh.battala@nutanix.com

The gflags management for Marina Server.
*/

package common

import (
	"flag"
)

// FlagsInit initializes flags.
func FlagsInit() {
	flag.Parse()
	flag.Set("logtostderr", "true")
}

var (
	MarinaGRPCPort = flag.Uint64("grpc_port", 9200,
		"Marina GRPC service port to server gRPC requests.")
	PrometheusPort = flag.Uint64("prometheus_port", 9201,
		"Port to listen Prometheus requests for Metrics collection.")
	MarinaRPCPort = flag.Uint64("proxyrpc_port", 9202,
		"Marina RPC service port to serve Proxy Requests for Legacy PC Catalog.")
	CatalogPort = flag.Uint64("legacycatalog_port", 2007,
		"Port on which legacy Catalog service is running.")
)
