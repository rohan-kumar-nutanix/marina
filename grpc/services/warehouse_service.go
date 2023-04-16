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
	internalIFC "github.com/nutanix-core/content-management-marina/interface/local"
	warehousePB "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"
	prismPB "github.com/nutanix-core/content-management-marina/protos/apis/prism/v4/config"
	util "github.com/nutanix-core/content-management-marina/util"
)

type WarehouseServer struct {
	warehousePB.UnimplementedWarehouseServiceServer
}

type WarehouseInterface interface {
	GetWarehouse(ctx context.Context, arg *warehousePB.GetWarehouseArg) (*warehousePB.GetWarehouseRet, error)
	CreateWarehouse(c context.Context, arg *warehousePB.CreateWarehouseArg) (*warehousePB.CreateWarehouseRet, error)
	ListWarehouse(context.Context, *warehousePB.ListWarehousesArg) (*warehousePB.ListWarehousesRet, error)
	AddItemToWarehouse(context.Context, *warehousePB.AddItemToWarehouseArg) (*warehousePB.AddItemToWarehouseRet, error)
	GetWarehouseItemById(context.Context, *warehousePB.GetWarehouseItemByIdArg) (*warehousePB.GetWarehouseItemByIdRet, error)
	ListAllWarehouseItems(context.Context, *warehousePB.ListAllWarehouseItemsArg) (*warehousePB.ListAllWarehouseItemsRet, error)
}

func (s *WarehouseServer) ListAllWarehouseItems(ctx context.Context, arg *warehousePB.ListAllWarehouseItemsArg) (
	*warehousePB.ListAllWarehouseItemsRet, error) {
	span := tracer.GetSpanFromContext(ctx)
	span.SetTag("ListWarehouseItems", arg.ExtId)
	log.Infof("ListAll WarehouseItems Arg received %v", arg)
	// return warehouse.ListWarehouseItems(ctx, arg)
	return nil, nil
}

func (s *WarehouseServer) GetWarehouse(ctx context.Context, arg *warehousePB.GetWarehouseArg) (*warehousePB.GetWarehouseRet, error) {
	span := tracer.GetSpanFromContext(ctx)
	span.SetTag("WarehouseUUID", arg.ExtId)

	log.Infof("Arg received %v", arg)
	// return warehouse.WarehouseGet(ctx, arg)
	return nil, nil
}

func (s *WarehouseServer) ListWarehouse(ctx context.Context, arg *warehousePB.ListWarehousesArg) (*warehousePB.ListWarehousesRet, error) {
	span := tracer.GetSpanFromContext(ctx)
	span.SetTag("ListWarehouse", "")
	log.Infof("Arg received %v", arg)
	// return warehouse.WarehouseGet(ctx, arg)
	// return warehouse.ListWarehouse(ctx, arg)
	return nil, nil
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
	// TODO write util code to get Task proto.
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

func setTaskData(taskUuid *uuid4.Uuid) *warehousePB.TaskReferenceWrapper {
	return &warehousePB.TaskReferenceWrapper{
		Value: &prismPB.TaskReference{
			ExtId: proto.String(taskUuid.String()),
		},
	}
}

// CreateWarehouse Handler.
func (s *WarehouseServer) CreateWarehouse(c context.Context,
	arg *warehousePB.CreateWarehouseArg) (*warehousePB.CreateWarehouseRet, error) {
	// TODO get operation in same method.
	taskUuid, err := s.asyncHandlerWithGrpcStatus(c, arg, CreateWarehouse)
	if err != nil {
		ret := &warehousePB.CreateWarehouseRet{
			Content: &warehousePB.CreateWarehouseResponse{
				Data: &warehousePB.CreateWarehouseResponse_ErrorResponseData{
					ErrorResponseData: nil,
				},
			},
		}
		return ret, err
	}
	ret := &warehousePB.CreateWarehouseRet{
		Content: &warehousePB.CreateWarehouseResponse{
			Data: &warehousePB.CreateWarehouseResponse_TaskReferenceData{
				TaskReferenceData: setTaskData(taskUuid),
			},
		},
	}
	return ret, nil
}

// AddItemToWarehouse Handler, to add warehouse item.
func (s *WarehouseServer) AddItemToWarehouse(c context.Context,
	arg *warehousePB.AddItemToWarehouseArg) (*warehousePB.AddItemToWarehouseRet, error) {
	// TODO get operation in same method.
	taskUuid, err := s.asyncHandlerWithGrpcStatus(c, arg, AddItemToWarehouse)
	if err != nil {
		ret := &warehousePB.AddItemToWarehouseRet{
			Content: &warehousePB.AddItemToWarehouseResponse{
				Data: &warehousePB.AddItemToWarehouseResponse_ErrorResponseData{
					ErrorResponseData: nil,
				},
			},
		}
		return ret, err
	}
	ret := &warehousePB.AddItemToWarehouseRet{
		Content: &warehousePB.AddItemToWarehouseResponse{
			Data: &warehousePB.AddItemToWarehouseResponse_TaskReferenceData{
				TaskReferenceData: setTaskData(taskUuid),
			},
		},
	}
	return ret, nil
}
