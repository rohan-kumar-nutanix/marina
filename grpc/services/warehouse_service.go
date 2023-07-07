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
	contentPB.UnimplementedWarehouseItemsServiceServer
}

type WarehouseInterface interface {
	// Warehouse Methods
	CreateWarehouse(c context.Context, arg *contentPB.CreateWarehouseArg) (*contentPB.CreateWarehouseRet, error)
	ListWarehouses(context.Context, *contentPB.ListWarehousesArg) (*contentPB.ListWarehousesRet, error)
	GetWarehouse(context.Context, *contentPB.GetWarehouseArg) (*contentPB.GetWarehouseRet, error)
	UpdateWarehouseMetadata(context.Context, *contentPB.UpdateWarehouseMetadataArg) (*contentPB.UpdateWarehouseMetadataRet, error)
	DeleteWarehouse(context.Context, *contentPB.DeleteWarehouseArg) (*contentPB.DeleteWarehouseRet, error)

	// WarehouseItems Methods
	AddItemToWarehouse(context.Context, *contentPB.AddItemToWarehouseArg) (*contentPB.AddItemToWarehouseRet, error)
	GetWarehouseItemById(context.Context, *contentPB.GetWarehouseItemByIdArg) (*contentPB.GetWarehouseItemByIdRet, error)
	ListWarehouseItems(context.Context, *contentPB.ListWarehouseItemsArg) (*contentPB.ListWarehouseItemsRet, error)
	DeleteWarehouseItem(context.Context, *contentPB.DeleteWarehouseItemArg) (*contentPB.DeleteWarehouseItemRet, error)
	UpdateWarehouseItemMetadata(context.Context, *contentPB.UpdateWarehouseItemMetadataArg) (*contentPB.UpdateWarehouseItemMetadataRet, error)
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

// SyncWarehouseMetadata gRPC handler
func (s *WarehouseServer) SyncWarehouseMetadata(c context.Context,
	arg *contentPB.SyncWarehouseMetadataArg) (*contentPB.SyncWarehouseMetadataRet, error) {
	taskUuid, err := s.asyncHandlerWithGrpcStatus(c, arg, SyncWarehouse)
	if err != nil {
		ret := &contentPB.SyncWarehouseMetadataRet{
			Content: &contentPB.SyncWarehouseResponse{
				Data: &contentPB.SyncWarehouseResponse_ErrorResponseData{
					ErrorResponseData: util.SetErrorResponse(nil),
				},
			},
		}
		return ret, err
	}
	ret := &contentPB.SyncWarehouseMetadataRet{
		Content: &contentPB.SyncWarehouseResponse{
			Data: &contentPB.SyncWarehouseResponse_TaskReferenceData{
				TaskReferenceData: setTaskData(taskUuid),
			},
		},
	}
	return ret, nil
}

// AddItemToWarehouse Handler, to add WarehouseItem to the Warehouse
func (s *WarehouseServer) AddItemToWarehouse(c context.Context,
	arg *contentPB.AddItemToWarehouseArg) (*contentPB.AddItemToWarehouseRet, error) {
	taskUuid, err := s.asyncHandlerWithGrpcStatus(c, arg, AddItemToWarehouse)
	if err != nil {
		ret := &contentPB.AddItemToWarehouseRet{
			Content: &contentPB.AddItemToWarehouseResponse{
				Data: &contentPB.AddItemToWarehouseResponse_ErrorResponseData{
					ErrorResponseData: nil,
				},
			},
		}
		return ret, err
	}
	ret := &contentPB.AddItemToWarehouseRet{
		Content: &contentPB.AddItemToWarehouseResponse{
			Data: &contentPB.AddItemToWarehouseResponse_TaskReferenceData{
				TaskReferenceData: setTaskData(taskUuid),
			},
		},
	}
	return ret, nil
}

// GetWarehouseItemById gRPC Handler.
func (s *WarehouseServer) GetWarehouseItemById(ctx context.Context, arg *contentPB.GetWarehouseItemByIdArg) (
	*contentPB.GetWarehouseItemByIdRet, error) {
	log.Infof("Arg received %v", arg)
	defer log.Infof("Finished GetWarehouseItemById RPC for WarehouseUUID: %s WarehouseItemUUID: ",
		*arg.ExtId, arg.ExtId)
	span := tracer.GetSpanFromContext(ctx)
	span.SetTag("WarehouseUUID", arg.ExtId)

	return warehouse.WarehouseItemGetWithGrpcStatus(ctx, arg)
}

// ListWarehouseItems gRPC Handler.
func (s *WarehouseServer) ListWarehouseItems(ctx context.Context, arg *contentPB.ListWarehouseItemsArg) (
	*contentPB.ListWarehouseItemsRet, error) {
	span := tracer.GetSpanFromContext(ctx)
	span.SetTag("ListWarehouseItems", "")
	return warehouse.WarehouseItemListWithGrpcStatus(ctx, arg)
}

// DeleteWarehouseItem gRPC Handler
func (s *WarehouseServer) DeleteWarehouseItem(ctx context.Context, arg *contentPB.DeleteWarehouseItemArg) (
	*contentPB.DeleteWarehouseItemRet, error) {
	taskUuid, err := s.asyncHandlerWithGrpcStatus(ctx, arg, DeleteWarehouseItem)
	if err != nil {
		ret := &contentPB.DeleteWarehouseItemRet{
			Content: &contentPB.DeleteWarehouseItemResponse{
				Data: &contentPB.DeleteWarehouseItemResponse_ErrorResponseData{
					ErrorResponseData: util.SetErrorResponse(nil),
				},
			},
		}
		return ret, err
	}
	ret := &contentPB.DeleteWarehouseItemRet{
		Content: &contentPB.DeleteWarehouseItemResponse{
			Data: &contentPB.DeleteWarehouseItemResponse_TaskReferenceData{
				TaskReferenceData: setTaskData(taskUuid),
			},
		},
	}
	return ret, nil

}

func (s *WarehouseServer) UpdateWarehouseItemMetadata(ctx context.Context, arg *contentPB.UpdateWarehouseItemMetadataArg) (
	*contentPB.UpdateWarehouseItemMetadataRet, error) {
	taskUuid, err := s.asyncHandlerWithGrpcStatus(ctx, arg, UpdateWarehouseItemMetadata)
	if err != nil {
		ret := &contentPB.UpdateWarehouseItemMetadataRet{
			Content: &contentPB.UpdateWarehouseItemResponse{
				Data: &contentPB.UpdateWarehouseItemResponse_ErrorResponseData{
					ErrorResponseData: util.SetErrorResponse(nil),
				},
			},
		}
		return ret, err
	}
	ret := &contentPB.UpdateWarehouseItemMetadataRet{
		Content: &contentPB.UpdateWarehouseItemResponse{
			Data: &contentPB.UpdateWarehouseItemResponse_TaskReferenceData{
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
