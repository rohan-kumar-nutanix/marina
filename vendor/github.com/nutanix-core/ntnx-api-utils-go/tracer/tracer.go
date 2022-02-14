//
// Copyright (c) 2021 Nutanix Inc. All rights reserved.
//
// Author: deepanshu.singhal@nutanix.com
//
// Initializes tracer for distributed tracing.
//
// Tracer will be initialized if file at path in gflag
// ObservabilityConfigFilePath is present and is in appropriate json format. In
// case the file is not present, tracer will not be initialized.
// Supported format for config file:
// {
// 	"Observability": {
// 		"Distributed tracing": {
// 			"Enabled": true,
// 			"Jaeger": {
// 				"Agent endpoint": "localhost:5775",
// 				"Agent Container endpoint": "jaeger-agent:5775",
// 				"Collector endpoint": "localhost:14250"
// 			},
// 			"Opentelemetry": {
//                 "Agent endpoint": "http://localhost:6831",
//                 "Agent Container endpoint": "jaeger-agent:5775",
//                 "Collector endpoint": "http://localhost:14268/api/traces"
//             }
// 		}
// 	}
// }

package tracer

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/golang/glog"
	otel "go.opentelemetry.io/otel"
	otel_jaeger "go.opentelemetry.io/otel/exporters/jaeger"
	otel_propagation "go.opentelemetry.io/otel/propagation"
	otel_resource "go.opentelemetry.io/otel/sdk/resource"
	otel_sdktrace "go.opentelemetry.io/otel/sdk/trace"
	otel_semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var tracingEnabled = false

// Initializes tracer if file at path in gflag ObservabilityConfigFilePath is in
// appropriate format.
// Usage:
//  func main() {
//    ...
//    closer := InitTracer("metropolis")
//    if closer != nil {
//      defer closer.Close()
//    }
//    ...
//  }
//
// Args:
//  serviceName: Name of service that is initializing tracer.
//
// Returns: io.Closer interface. Used to close tracer istance.

func InitTracer(serviceName string) (*otel_sdktrace.TracerProvider, error) {
	enabled, collector, err := checkTracingEnabled()
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, nil
	}

	exporter, _ := otel_jaeger.New(
		otel_jaeger.WithCollectorEndpoint(otel_jaeger.WithEndpoint(collector)),
	)
	resources := otel_resource.NewWithAttributes(
		otel_semconv.SchemaURL,
		otel_semconv.ServiceNameKey.String(serviceName),
	)
	traceProvider := otel_sdktrace.NewTracerProvider(otel_sdktrace.WithBatcher(exporter), otel_sdktrace.WithResource(resources))
	propagator := otel_propagation.NewCompositeTextMapPropagator(otel_propagation.Baggage{}, otel_propagation.TraceContext{})

	tracingEnabled = true
	otel.SetTextMapPropagator(propagator)
	otel.SetTracerProvider(traceProvider)
	return traceProvider, nil
}

func checkTracingEnabled() (bool, string, error) {
	data, err := os.ReadFile(*observabilityConfigFilePath)
	if err == nil {
		enabled, collector, err := getTraceDataFromConfig(data)
		if err != nil {
			return false, "", err
		}

		return enabled, collector, nil
	} else if os.IsNotExist(err) {
		// Config file does not exist. Disable tracing.
		return false, "", nil
	} else {
		// Error while reading config file. Disable tracing
		glog.Error("Error while reading config file: ", err)
		return false, "", nil
	}
}

func getTraceDataFromConfig(config []byte) (bool, string, error) {

	if config == nil {
		glog.Error("Invalid config provided: ")
		return false, "", errors.New("Invalid config")
	}

	var data map[string]interface{}
	err := json.Unmarshal(config, &data)
	if err != nil {
		glog.Error("Error while unmarshalling data from config file: ", err)
		return false, "", err
	}

	if data["Observability"] == nil {
		errMsg := fmt.Sprintf("Config file does not contain key Observability: %v",
			data)
		glog.Error(errMsg)
		return false, "", errors.New(errMsg)
	}

	obsv, ok := data["Observability"].(map[string]interface{})
	if !ok {
		errMsg := fmt.Sprintf("Config file contains invalid value for "+
			"Observability: %v", data)
		glog.Error(errMsg)
		return false, "", errors.New(errMsg)
	}

	if obsv["Distributed tracing"] == nil {
		errMsg := fmt.Sprintf("Config file does not contain key Distributed "+
			"tracing: %v", data)
		glog.Error(errMsg)
		return false, "", errors.New(errMsg)
	}

	traceData, ok := obsv["Distributed tracing"].(map[string]interface{})
	if !ok {
		errMsg := fmt.Sprintf("Config file contains invalid value for "+
			"Distributed tracing: %v", data)
		glog.Error(errMsg)
		return false, "", errors.New(errMsg)
	}

	if traceData["Enabled"] == nil {
		errMsg := fmt.Sprintf("Config file does not contain key Enabled: %v", data)
		glog.Error(errMsg)
		return false, "", errors.New(errMsg)
	}

	enabled, ok := traceData["Enabled"].(bool)
	if !ok {
		errMsg := fmt.Sprintf("Config file contains invalid value for Enabled: %v",
			data)
		glog.Error(errMsg)
		return false, "", errors.New(errMsg)
	}

	if !enabled {
		return enabled, "", nil
	}

	/*
		ToDo: put logic to initilize only one type of tracer
	*/
	// if traceData["Jaeger"] != nil {
	// 	jaegerData, ok := traceData["Jaeger"].(map[string]interface{})
	// 	if !ok {
	// 		errMsg := fmt.Sprintf("Config file contains invalid value for Jaeger: %v",
	// 			data)
	// 		glog.Error(errMsg)
	// 		return enabled, "", "", errors.New(errMsg)
	// 	}

	// 	if jaegerData["Agent endpoint"] != nil {
	// 		endpoint, ok := jaegerData["Agent endpoint"].(string)
	// 		if !ok {
	// 			errMsg := fmt.Sprintf("Config file contains invalid value for Agent "+
	// 				"endpoint: %v", data)
	// 			glog.Error(errMsg)
	// 			return enabled, "", "", errors.New(errMsg)
	// 		}

	// 		return enabled, endpoint, "", nil
	// 	}
	// } else

	if traceData["Opentelemetry"] != nil {
		jaegerData, ok := traceData["Opentelemetry"].(map[string]interface{})
		if !ok {
			errMsg := fmt.Sprintf("Config file contains invalid value for Opentelemetry: %v",
				data)
			glog.Error(errMsg)
			return enabled, "", errors.New(errMsg)
		}

		//Checking collector endpoint config for opentelemetry
		if jaegerData["Collector endpoint"] != nil {
			collector, ok := jaegerData["Collector endpoint"].(string)
			if !ok {
				errMsg := fmt.Sprintf("Config file contains invalid value for Collector "+
					"endpoint: %v", data)
				glog.Error(errMsg)
				return enabled, "", errors.New(errMsg)
			}

			return enabled, collector, nil
		}
	}

	return enabled, "", nil
}
