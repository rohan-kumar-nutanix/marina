/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
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
	"github.com/nutanix-core/content-management-marina/grpc/warehouse"
	internalIFC "github.com/nutanix-core/content-management-marina/interface/local"
	contentPB "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"
	prismPB "github.com/nutanix-core/content-management-marina/protos/apis/prism/v4/config"
	util "github.com/nutanix-core/content-management-marina/util"
)

type WarehouseServer struct {
	contentPB.UnimplementedWarehouseServiceServer
}

type WarehouseInterface interface {
	CreateWarehouse(c context.Context, arg *contentPB.CreateWarehouseArg) (*contentPB.CreateWarehouseRet, error)
	ListWarehouses(context.Context, *contentPB.ListWarehousesArg) (*contentPB.ListWarehousesRet, error)
	GetWarehouse(context.Context, *contentPB.GetWarehouseArg) (*contentPB.GetWarehouseRet, error)
	UpdateWarehouseMetadata(context.Context, *contentPB.UpdateWarehouseMetadataArg) (*contentPB.UpdateWarehouseMetadataRet, error)
	DeleteWarehouse(context.Context, *contentPB.DeleteWarehouseArg) (*contentPB.DeleteWarehouseRet, error)
}

func (s *WarehouseServer) GetWarehouse(ctx context.Context, arg *contentPB.GetWarehouseArg) (*contentPB.GetWarehouseRet, error) {
	log.Infof("Arg received %v", arg)
	defer log.Infof("Finished GetWarehouse RPC for WarehouseUUID: %s", *arg.ExtId)
	span := tracer.GetSpanFromContext(ctx)
	span.SetTag("WarehouseUUID", arg.ExtId)

	return warehouse.WarehouseGetWithGrpcStatus(ctx, arg)
}

// ListWarehouses gRPC Handler.
func (s *WarehouseServer) ListWarehouses(ctx context.Context, arg *contentPB.ListWarehousesArg) (
	*contentPB.ListWarehousesRet, error) {
	span := tracer.GetSpanFromContext(ctx)
	span.SetTag("ListWarehouse", "")
	return warehouse.WarehouseListWithGrpcStatus(ctx, arg)
}

// CreateWarehouse gRPC Handler.
func (s *WarehouseServer) CreateWarehouse(c context.Context,
	arg *contentPB.CreateWarehouseArg) (*contentPB.CreateWarehouseRet, error) {
	// TODO get operation in same method.
	taskUuid, err := s.asyncHandlerWithGrpcStatus(c, arg, CreateWarehouse)
	if err != nil {
		ret := &contentPB.CreateWarehouseRet{
			Content: &contentPB.CreateWarehouseResponse{
				Data: &contentPB.CreateWarehouseResponse_ErrorResponseData{
					ErrorResponseData: util.SetErrorResponse(nil),
				},
			},
		}
		return ret, err
	}
	ret := &contentPB.CreateWarehouseRet{
		Content: &contentPB.CreateWarehouseResponse{
			Data: &contentPB.CreateWarehouseResponse_TaskReferenceData{
				TaskReferenceData: setTaskData(taskUuid),
			},
		},
	}
	return ret, nil
}

// DeleteWarehouse gRPC Handler.
func (s *WarehouseServer) DeleteWarehouse(c context.Context,
	arg *contentPB.DeleteWarehouseArg) (*contentPB.DeleteWarehouseRet, error) {
	taskUuid, err := s.asyncHandlerWithGrpcStatus(c, arg, DeleteWarehouse)
	if err != nil {
		ret := &contentPB.DeleteWarehouseRet{
			Content: &contentPB.DeleteWarehouseResponse{
				Data: &contentPB.DeleteWarehouseResponse_ErrorResponseData{
					ErrorResponseData: util.SetErrorResponse(nil),
				},
			},
		}
		return ret, err
	}
	ret := &contentPB.DeleteWarehouseRet{
		Content: &contentPB.DeleteWarehouseResponse{
			Data: &contentPB.DeleteWarehouseResponse_TaskReferenceData{
				TaskReferenceData: setTaskData(taskUuid),
			},
		},
	}
	return ret, nil
}

// UpdateWarehouseMetadata  gRPC Handler.
func (s *WarehouseServer) UpdateWarehouseMetadata(c context.Context,
	arg *contentPB.UpdateWarehouseMetadataArg) (*contentPB.UpdateWarehouseMetadataRet, error) {
	// TODO get operation in same method.
	taskUuid, err := s.asyncHandlerWithGrpcStatus(c, arg, UpdateWarehouse)
	if err != nil {
		ret := &contentPB.UpdateWarehouseMetadataRet{
			Content: &contentPB.UpdateWarehouseResponse{
				Data: &contentPB.UpdateWarehouseResponse_ErrorResponseData{
					ErrorResponseData: util.SetErrorResponse(nil),
				},
			},
		}
		return ret, err
	}
	ret := &contentPB.UpdateWarehouseMetadataRet{
		Content: &contentPB.UpdateWarehouseResponse{
			Data: &contentPB.UpdateWarehouseResponse_TaskReferenceData{
				TaskReferenceData: setTaskData(taskUuid),
			},
		},
	}
	return ret, nil
}

func (s *WarehouseServer) asyncHandlerWithGrpcStatus(c context.Context, request proto.Message,
	operation string) (*uuid4.Uuid, error) {
	// TODO add marinaError.AppMessage for APIs
	taskUuid, err := s.asyncHandler(c, request, operation)
	if err != nil {
		// TODO: Create gRPC error from MarinaError and generate AppMessage Error and return to client.
		return taskUuid, internalIFC.Interfaces().ErrorIfc().BuildGrpcError(err)
	}
	return taskUuid, nil
}

func (s *WarehouseServer) asyncHandler(c context.Context, request proto.Message,
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

func setTaskData(taskUuid *uuid4.Uuid) *contentPB.TaskReferenceWrapper {
	return &contentPB.TaskReferenceWrapper{
		Value: &prismPB.TaskReference{
			ExtId: proto.String(taskUuid.String()),
		},
	}
}
