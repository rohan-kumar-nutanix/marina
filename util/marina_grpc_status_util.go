/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Authors: rajesh.battala@nutanix.com
 *
 * gRPC status builder utils for Marina.
 *
 *
 */

package utils

import (
	statusPb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	marinaError "github.com/nutanix-core/content-management-marina/errors"
	"github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"
	ntnxApiCmsError "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/error"
)

const (
	EnglishLocale           = "en_US"
	MarinaInternalErrorCode = "10000"
)

type GrpcStatusUtilInterface interface {
	BuildGrpcErrorWithCode(codes.Code, marinaError.MarinaErrorInterface) error
	BuildGrpcError(marinaError.MarinaErrorInterface) error
}

type GrpcStatusUtil struct {
}

func NewGrpcStatusUtil() *GrpcStatusUtil {
	return &GrpcStatusUtil{}
}

func (e *GrpcStatusUtil) BuildGrpcErrorWithCode(grpcCode codes.Code,
	marinaErr marinaError.MarinaErrorInterface) error {
	return status.ErrorProto(
		&statusPb.Status{
			Code:    int32(grpcCode),
			Message: marinaErr.GetErrorDetail(),
		})
}

func (e *GrpcStatusUtil) BuildGrpcError(marinaErr marinaError.MarinaErrorInterface) error {
	return e.BuildGrpcErrorWithCode(marinaError.GetGrpcCodeFromMarinaError(marinaErr), marinaErr)
}

func SetErrorResponse(appMessageList []*ntnxApiCmsError.AppMessage) *content.ErrorResponseWrapper {
	return &content.ErrorResponseWrapper{
		Value: &ntnxApiCmsError.ErrorResponse{
			Error: &ntnxApiCmsError.ErrorResponse_AppMessageArrayError{
				AppMessageArrayError: &ntnxApiCmsError.AppMessageArrayWrapper{
					Value: appMessageList,
				},
			},
		},
	}
}
