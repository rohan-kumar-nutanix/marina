/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 * Scanner Config Server for Managing the Warehouse Scanner Configurations.
 */

package services

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/ergon"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/tracer"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	marinaError "github.com/nutanix-core/content-management-marina/errors"
	"github.com/nutanix-core/content-management-marina/grpc/scanner_config"
	internalIFC "github.com/nutanix-core/content-management-marina/interface/local"
	configPB "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/config"
	prismPB "github.com/nutanix-core/content-management-marina/protos/apis/prism/v4/config"
	util "github.com/nutanix-core/content-management-marina/util"
)

type ScannerConfigurationServer struct {
	configPB.UnimplementedScannerToolsServiceServer
}

type ScannerConfiguration interface {
	CreateScannerToolConfig(context.Context, *configPB.CreateScannerToolConfigArg) (*configPB.CreateScannerToolConfigRet, error)
	DeleteScannerToolConfigByExtId(context.Context, *configPB.DeleteScannerToolConfigByExtIdArg) (*configPB.DeleteScannerToolConfigByExtIdRet, error)
	GetScannerToolConfigByExtId(context.Context, *configPB.GetScannerToolConfigByExtIdArg) (*configPB.GetScannerToolConfigByExtIdRet, error)
	ListScannerToolConfig(context.Context, *configPB.ListScannerToolConfigArg) (*configPB.ListScannerToolConfigRet, error)
	UpdateScannerToolConfigByExtId(context.Context, *configPB.UpdateScannerToolConfigByExtIdArg) (*configPB.UpdateScannerToolConfigByExtIdRet, error)
}

// CreateScannerToolConfig gRPC Handler.
func (s *ScannerConfigurationServer) CreateScannerToolConfig(ctx context.Context, arg *configPB.CreateScannerToolConfigArg) (
	*configPB.CreateScannerToolConfigRet, error) {
	taskUuid, err := s.asyncHandlerWithGrpcStatus(ctx, arg, CreateScannerToolConfig)
	if err != nil {
		ret := &configPB.CreateScannerToolConfigRet{
			Content: &configPB.CreateScannerToolConfigResponse{
				Data: &configPB.CreateScannerToolConfigResponse_ErrorResponseData{
					ErrorResponseData: util.GetConfigErrorResponse(nil),
				},
			},
		}
		return ret, err
	}
	ret := &configPB.CreateScannerToolConfigRet{
		Content: &configPB.CreateScannerToolConfigResponse{
			Data: &configPB.CreateScannerToolConfigResponse_TaskReferenceData{
				TaskReferenceData: s.getTaskData(taskUuid),
			},
		},
	}
	return ret, nil
}

// GetScannerToolConfigByExtId gRPC Handler

func (s *ScannerConfigurationServer) GetScannerToolConfigByExtId(ctx context.Context, arg *configPB.GetScannerToolConfigByExtIdArg) (
	*configPB.GetScannerToolConfigByExtIdRet, error) {
	log.Infof("GetScannerToolConfigByExtIdArg Arg received %v and RPC context %s", arg, ctx)
	defer log.Infof("Finished GetScannerToolConfigByExtIdArg RPC for ScannerConfiguration UUID: %s", *arg.ExtId)
	span := tracer.GetSpanFromContext(ctx)
	span.SetTag("ScannerConfigUUID", *arg.ExtId)

	return scanner_config.ScannerConfigGetWithGrpcStatus(ctx, arg)
}

// ListSecurityPolicies gRPC Handler

func (s *ScannerConfigurationServer) ListScannerToolConfig(ctx context.Context, arg *configPB.ListScannerToolConfigArg) (
	*configPB.ListScannerToolConfigRet, error) {
	log.Infof("ListScannerToolConfigArg Arg received %v RPC Context %s", arg, ctx)
	defer log.Infof("Finished ListScannerToolConfig RPC ")
	span := tracer.GetSpanFromContext(ctx)
	span.SetTag("ListScannerToolConfig", "")

	return scanner_config.ScannerConfigListWithGrpcStatus(ctx, arg)
}

// DeleteScannerToolConfigByExtId gRPC Handler
func (s *ScannerConfigurationServer) DeleteScannerToolConfigByExtId(ctx context.Context, arg *configPB.DeleteScannerToolConfigByExtIdArg) (
	*configPB.DeleteScannerToolConfigByExtIdRet, error) {
	taskUuid, err := s.asyncHandlerWithGrpcStatus(ctx, arg, DeleteScannerToolConfigByExtId)
	if err != nil {
		ret := &configPB.DeleteScannerToolConfigByExtIdRet{
			Content: &configPB.DeleteScannerToolConfigResponse{
				Data: &configPB.DeleteScannerToolConfigResponse_ErrorResponseData{
					ErrorResponseData: util.GetConfigErrorResponse(nil),
				},
			},
		}
		return ret, err
	}
	ret := &configPB.DeleteScannerToolConfigByExtIdRet{
		Content: &configPB.DeleteScannerToolConfigResponse{
			Data: &configPB.DeleteScannerToolConfigResponse_TaskReferenceData{
				TaskReferenceData: s.getTaskData(taskUuid),
			},
		},
	}
	return ret, nil
}

// UpdateSecurityPolicyByExtId gRPC Handler
func (s *ScannerConfigurationServer) UpdateScannerToolConfigByExtId(ctx context.Context,
	arg *configPB.UpdateScannerToolConfigByExtIdArg) (*configPB.UpdateScannerToolConfigByExtIdRet, error) {
	taskUuid, err := s.asyncHandlerWithGrpcStatus(ctx, arg, UpdateScannerToolConfigByExtId)
	if err != nil {
		ret := &configPB.UpdateScannerToolConfigByExtIdRet{
			Content: &configPB.UpdateScannerToolConfigResponse{
				Data: &configPB.UpdateScannerToolConfigResponse_ErrorResponseData{
					ErrorResponseData: util.GetConfigErrorResponse(nil),
				},
			},
		}
		return ret, err
	}
	ret := &configPB.UpdateScannerToolConfigByExtIdRet{
		Content: &configPB.UpdateScannerToolConfigResponse{
			Data: &configPB.UpdateScannerToolConfigResponse_TaskReferenceData{
				TaskReferenceData: s.getTaskData(taskUuid),
			},
		},
	}
	return ret, nil
}

func (s *ScannerConfigurationServer) asyncHandlerWithGrpcStatus(c context.Context, request proto.Message,
	operation string) (*uuid4.Uuid, error) {
	// TODO add marinaError.AppMessage for APIs
	taskUuid, err := s.asyncHandler(c, request, operation)
	if err != nil {
		// TODO: Create gRPC error from MarinaError and generate AppMessage Error and return to client.
		return taskUuid, internalIFC.Interfaces().ErrorIfc().BuildGrpcError(err)
	}
	return taskUuid, nil
}

func (s *ScannerConfigurationServer) asyncHandler(c context.Context, request proto.Message,
	operation string) (*uuid4.Uuid, marinaError.MarinaErrorInterface) {
	embeddedReq, err := internalIFC.Interfaces().ProtoIfc().Marshal(request)
	if err != nil {
		return nil, marinaError.ErrInternalError().SetCauseAndLog(fmt.Errorf("could not Marshal the request"))
	}
	taskProto := &ergon.Task{
		Request: &ergon.MetaRequest{
			MethodName: proto.String(operation),
			Arg:        &ergon.PayloadOrEmbeddedValue{Embedded: embeddedReq},
		},
	}
	task := GetErgonFullTaskByProto(taskProto)
	if task == nil {
		return nil, marinaError.ErrMarinaNotSupportedError(operation)
	}
	err = task.Start(task, proto.String(util.ServiceName), &operation)
	if err != nil {
		// Failed to start the Task.
		return nil, marinaError.ErrMarinaInternalError().SetCauseAndLog(
			fmt.Errorf("unable to start operation %s task: %s", operation, err))
	}
	taskUuid := uuid4.ToUuid4(task.Proto().GetUuid())
	log.Infof("Created a Marina %s task with UUID: %s", operation, taskUuid.String())
	return taskUuid, nil
}

func (s *ScannerConfigurationServer) getTaskData(taskUuid *uuid4.Uuid) *configPB.TaskReferenceWrapper {
	return &configPB.TaskReferenceWrapper{
		Value: &prismPB.TaskReference{
			ExtId: proto.String(taskUuid.String()),
		},
	}
}
