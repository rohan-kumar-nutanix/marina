/*
 *
 * Marina service errors.
 *
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Author: Rajesh Battala <rajesh.battala@nutanix.com>
 *
 * Error Utilities
 */

package errors

import (
	"context"

	"github.com/nutanix-core/ntnx-api-utils-go/linkutils"

	"google.golang.org/protobuf/proto"

	"github.com/nutanix-core/content-management-marina/protos/apis/common/v1/config"
)

// Construct the hostPort and the V4 version for the link.
func constructLink(c context.Context) (string, string, error) {
	hostPort, err := linkutils.GetOriginHostPortFromGrpcContext(c)
	if err != nil {
		return "", "", err
	}
	version, err := linkutils.GetRequestFullVersionFromGrpcContext(c)
	if err != nil {
		return "", "", err
	}

	return hostPort, version, nil
}

/*func getAPILink(ctx context.Context) *response.ApiLink {

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
}*/

func getMetadataFlags() []*config.Flag {
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

/*func BuildGrpcErrorWithCode(httpCode int, msg string) error {
	code := securityError.CodeFromHTTPStatus(httpCode)
	return status.ErrorProto(
		&statusPb.Status{
			Code:    int32(code),
			Message: msg,
		})
}*/

// GetVariableFromGrpcContext read the metadata from grpc context and return
// the variable if exists in metadata.
/*func GetVariableFromGrpcContext(ctx context.Context, varName string) []string {
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
}*/

/*// GetLocaleFromGrpcContext returns the locale variant by reading the
// grpcgateway-accept-language value from grpc context metadata.
func GetLocaleFromGrpcContext(ctx context.Context) string {
	acceptLang := GetVariableFromGrpcContext(ctx, GRPCGatewayAcceptLanguage)
	if len(acceptLang) > 0 {
		return acceptLang[0]
	}
	return DefaultLocale
}*/

// GetPathFromGrpcContext returns the origin path by reading the
// x-envoy-original-path value from grpc context metadata.
// return : string ex. /api/security/v4.0.a1/dashboard/foo
/*func GetPathFromGrpcContext(ctx context.Context) string {
	uriPath := GetVariableFromGrpcContext(ctx, EnvoyOriginalPath)
	if len(uriPath) > 0 {
		return uriPath[0]
	}
	return ""
}*/
