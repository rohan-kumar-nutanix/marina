/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Authors: rajesh.battala@nutanix.com
*
* Implementation of Get and List of Warehouse ScannerConfiguration entity.
 */

package scanner_config

import (
	"context"
	"net/url"
	"strings"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/tracer"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/ntnx-api-utils-go/linkutils"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/content-management-marina/errors"
	"github.com/nutanix-core/content-management-marina/interface/external"
	internal "github.com/nutanix-core/content-management-marina/interface/local"
	configPB "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/config"
	"github.com/nutanix-core/content-management-marina/protos/apis/common/v1/response"
	utils "github.com/nutanix-core/content-management-marina/util"
)

const (
	EnvoyOriginalPath         = "x-envoy-original-path"
	GRPCGatewayAcceptLanguage = "grpcgateway-accept-language"
	DefaultLocale             = "en_US"
)

func ScannerConfigGetWithGrpcStatus(ctx context.Context, arg *configPB.GetScannerToolConfigByExtIdArg) (
	*configPB.GetScannerToolConfigByExtIdRet, error) {
	span, ctx := tracer.StartSpan(ctx, "GetScannerConfigDetails")
	defer span.Finish()
	ret, err := ScannerConfigGet(ctx, arg)
	log.Infof("Returning Scanner Config Entity Details Ret %v", ret)

	if err != nil {
		ret = &configPB.GetScannerToolConfigByExtIdRet{
			Content: &configPB.GetScannerToolConfigResponse{
				Data: &configPB.GetScannerToolConfigResponse_ErrorResponseData{
					ErrorResponseData: utils.GetConfigErrorResponse(nil),
				},
			},
		}
		return ret, internal.Interfaces().ErrorIfc().BuildGrpcError(errors.ErrNotFound)
	}

	return ret, nil

}

func ScannerConfigGet(ctx context.Context, arg *configPB.GetScannerToolConfigByExtIdArg) (
	*configPB.GetScannerToolConfigByExtIdRet, error) {
	scannerConfigUuid, _ := uuid4.StringToUuid4(*arg.ExtId)
	scannerConfig, err := newScannerConfigDBImpl().GetScannerConfig(ctx, external.Interfaces().CPDBIfc(),
		scannerConfigUuid)
	if err != nil {
		return nil, err
	}

	metadata, err := getScannerToolGetMetadata(ctx, scannerConfigUuid.String())
	if err != nil {
		log.Errorf("Error occurred while generating Metadata for GetScannerToolConfigByExtId %s", err)
	}
	ret := &configPB.GetScannerToolConfigByExtIdRet{
		Content: &configPB.GetScannerToolConfigResponse{
			Data: &configPB.GetScannerToolConfigResponse_ScannerConfigData{
				ScannerConfigData: &configPB.ScannerConfigWrapper{
					Value: scannerConfig,
				},
			},
			Metadata: metadata,
		},
	}
	return ret, nil

}

func ScannerConfigList(ctx context.Context, arg *configPB.ListScannerToolConfigArg) (
	*configPB.ListScannerToolConfigRet, error) {
	scannerConfigs, err := newScannerConfigDBImpl().ListScannerConfigurations(ctx, external.Interfaces().CPDBIfc())
	if err != nil {
		log.Errorf("error occurred while fetching the ScannerConfiguration List %s", err)
		return nil, err
	}

	// metadata, err := getScannerToolGetMetadata(ctx, "")

	ret := &configPB.ListScannerToolConfigRet{
		Content: &configPB.ListScannerToolConfigsResponse{
			Data: &configPB.ListScannerToolConfigsResponse_ScannerConfigArrayData{
				ScannerConfigArrayData: &configPB.ScannerConfigArrayWrapper{
					Value: scannerConfigs,
				},
			},
			// Metadata: metadata,
		},
	}
	return ret, nil
}

func ScannerConfigListWithGrpcStatus(ctx context.Context, arg *configPB.ListScannerToolConfigArg) (
	*configPB.ListScannerToolConfigRet, error) {
	span, ctx := tracer.StartSpan(ctx, "ListScannerConfigurations")
	defer span.Finish()
	ret, err := ScannerConfigList(ctx, arg)
	log.V(2).Infof("Returning List of Scanner Configuration Entity Details Ret %v", ret)
	log.Infof("Returning List of Scanner Configuration Entity Details Ret %v", ret)

	if err != nil {
		ret = &configPB.ListScannerToolConfigRet{
			Content: &configPB.ListScannerToolConfigsResponse{
				Data: &configPB.ListScannerToolConfigsResponse_ErrorResponseData{
					ErrorResponseData: utils.GetConfigErrorResponse(nil),
				},
			},
		}
		return ret, internal.Interfaces().ErrorIfc().BuildGrpcError(errors.ErrNotFound)
	}
	return ret, nil
}

// Get metadata corresponding to the task get response.
func getScannerToolGetMetadata(ctx context.Context,
	scannerToolConfigUUID string) (*response.ApiResponseMetadata, error) {
	metadata := new(response.ApiResponseMetadata)
	metadata.Links = new(response.ApiLinkArrayWrapper)
	metadata.Links.Value = make([]*response.ApiLink, 0, 2)

	// Populate the self link.
	/*selfLink := &response.ApiLink{}
	href, err := utils.GetLink(ctx, "scannertool-configuration", scannerToolConfigUUID)
	if err != nil {
		log.Errorf("Error occurred in getting links %v", err)
		return nil, err
	}
	log.Infof("href generated %v", href)
	selfLink.Href = proto.String(href) // proto.String("https://10.96.208.32:9440/api/cms/v4.0.a1/config/security-policies/6d377c72-388d-4582-5a1a-81a7970e8eb8") // href)
	selfLink.Rel = proto.String(utils.SELF_RELATION)*/
	log.Infof("Fetching Link from the context. %v", ctx)
	selfLink := getAPILink(ctx)
	metadata.Links.Value = append(metadata.Links.Value, selfLink)

	/*if shouldPopulateCancellationLink(warehouse) {
		// Populate the cancel link.
		cancelLink := new(response_ifc.ApiLink)
		href, err := getLink(ctx, warehouse.GetExtId(), common.CANCEL_LINK)
		if err != nil {
			return nil, err
		}
		cancelLink.Href = href
		cancelLink.Rel = common.CANCEL_RELATION
		metadata.Links = append(metadata.Links, cancelLink)
	}*/

	// flags := utils.GetMetadataFlags()

	metadata1 := &response.ApiResponseMetadata{
		Flags:                 utils.GetMetadataFlagWrapper(),
		Links:                 metadata.Links,
		TotalAvailableResults: proto.Int32(1),
	}
	log.Infof("Scanner Config Metadata :%v", metadata1)
	return metadata1, nil
}

// GetLocaleFromGrpcContext returns the locale variant by reading the
// grpcgateway-accept-language value from grpc context metadata.
func GetLocaleFromGrpcContext(ctx context.Context) string {
	acceptLang := GetVariableFromGrpcContext(ctx, GRPCGatewayAcceptLanguage)
	if len(acceptLang) > 0 {
		return acceptLang[0]
	}
	return DefaultLocale
}

// GetPathFromGrpcContext returns the origin path by reading the
// x-envoy-original-path value from grpc context metadata.
// return : string ex. /api/security/v4.0.a1/dashboard/foo
func GetPathFromGrpcContext(ctx context.Context) string {
	uriPath := GetVariableFromGrpcContext(ctx, EnvoyOriginalPath)
	if len(uriPath) > 0 {
		return uriPath[0]
	}
	return ""
}

// GetVariableFromGrpcContext read the metadata from grpc context and return
// the variable if exists in metadata.
func GetVariableFromGrpcContext(ctx context.Context, varName string) []string {
	if ctx == nil {
		log.Errorf("gRPC context is nil")
		return []string{}
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Errorf("gRPC context doesn't have metadata")
		return []string{}
	}
	retVal, ok := md[varName]
	if !ok {
		log.Errorf("gRPC context metadata doesn't have  %v", varName)
		return []string{}
	}
	return retVal
}

func getAPILink(ctx context.Context) *response.ApiLink {

	hostPort, err := linkutils.GetOriginHostPortFromGrpcContext(ctx)
	if err != nil {
		log.Errorf("Error in getting Host Port info from ctx %v", err)
	}
	if !strings.HasPrefix(hostPort, "https://") {
		hostPort = "https://" + hostPort
	}
	uriPath, err := url.Parse(GetPathFromGrpcContext(ctx))
	if err != nil {
		log.Errorf("Error in parsing URI path: %v", err)
	}
	uriPath.RawQuery = ""
	uriBase, err := url.Parse(hostPort)
	if err != nil {
		log.Errorf("Error in parsing Host Port for URI: %v", err)
	}
	selfUri := uriBase.ResolveReference(uriPath).String()

	apiLink := &response.ApiLink{
		Href: proto.String(selfUri),
		Rel:  proto.String(uriPath.String()),
	}
	return apiLink
}
