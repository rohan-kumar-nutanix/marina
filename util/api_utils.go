/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Authors: rajesh.battala@nutanix.com
 *
 */

package utils

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/ntnx-api-utils-go/linkutils"

	"github.com/nutanix-core/content-management-marina/protos/apis/common/v1/config"
)

const (
	DefaultContextPath = "api"
	WarehouseUri       = "https://%s/%s/cms/%s/content/warehouses/%s"
	CmsURI             = "https://%s/%s/cms/%s/%s/%s/%s"
	SELF_RELATION      = "self"
)

// Get link to fill in the api response depending upon the link type.
func GetLink(c context.Context, entityName, extId string) (
	string, error) {
	apiModule := map[string]string{
		"warehouses":                 "content",
		"warehouseitems":             "content",
		"security-policies":          "config",
		"scannertool-configurations": "config",
	}

	hostPort, version, err := constructLink(c)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(CmsURI, hostPort, DefaultContextPath, version, apiModule["entityName"], entityName, extId), nil
	// return GetWarehouseUri(hostPort, version, extId), nil

}

// GetWarehouseUri Gets the task URI given a host-port, task API version and task id.
// Example: ("192.168.0.1:9440", "v4.0.a1", "e370710c-093f-48f7-b65d-3c610e7fb3d6s") will return
// "https://192.168.0.1:9440/api/cms/v4.0.a1/warehouses/e370710c-093f-48f7-b65d-3c610e7fb3d6"
func GetWarehouseUri(hostPort string, version string, extId string) string {
	return fmt.Sprintf(WarehouseUri, hostPort, DefaultContextPath, version, extId)

}

func GetMetadataFlags() []*config.Flag {
	return []*config.Flag{
		{
			Name:  proto.String("hasError"),
			Value: proto.Bool(false),
		},
		{
			Name:  proto.String("isPaginated"),
			Value: proto.Bool(false),
		},
	}
}

func GetMetadataFlagWrapper() *config.FlagArrayWrapper {
	return &config.FlagArrayWrapper{Value: GetMetadataFlags()}

}

// Construct the hostPort and the V4 version for the link.
func constructLink(c context.Context) (string, string, error) {
	log.Info("constructLink context values %s", c)
	hostPort, err := linkutils.GetOriginHostPortFromGrpcContext(c)
	if err != nil {
		return "", "", err
	}
	version, err := linkutils.GetRequestFullVersionFromGrpcContext(c)
	if err != nil {
		return "", "", err
	}
	log.Info("constructLink hostport %s and version %s", hostPort, version)
	return hostPort, version, nil
}
