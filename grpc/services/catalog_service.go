/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
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
	"github.com/nutanix-core/content-management-marina/grpc/catalog/catalog_item"
	"github.com/nutanix-core/content-management-marina/interface/external"
	internal "github.com/nutanix-core/content-management-marina/interface/local"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	util "github.com/nutanix-core/content-management-marina/util"
)

type MarinaServer struct {
	marinaIfc.UnimplementedMarinaServer
}

type MarinaServiceInterface interface {
	CatalogItemGet(ctx context.Context, arg *marinaIfc.CatalogItemGetArg) (*marinaIfc.CatalogItemGetRet, error)
	CatalogItemDelete(ctx context.Context, arg *marinaIfc.CatalogItemDeleteArg) (*marinaIfc.CatalogItemDeleteRet, error)
	CatalogItemCreate(ctx context.Context, arg *marinaIfc.CatalogItemCreateArg) (*marinaIfc.CatalogItemCreateRet, error)
	CatalogItemUpdate(ctx context.Context, arg *marinaIfc.CatalogItemUpdateArg) (*marinaIfc.CatalogItemUpdateRet, error)
}

func (s *MarinaServer) asyncHandler(ctx context.Context, request proto.Message, operation string) ([]byte, error) {
	embeddedReq, err := proto.Marshal(request)
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

	taskUuid := task.Proto().GetUuid()
	log.Infof("Created a Marina %s task with UUID: %s", operation, uuid4.ToUuid4(taskUuid).String())
	return taskUuid, nil
}

func (s *MarinaServer) CatalogItemGet(ctx context.Context, arg *marinaIfc.CatalogItemGetArg) (
	*marinaIfc.CatalogItemGetRet, error) {

	span, ctx := tracer.StartSpan(ctx, "catalogitem-get")
	defer span.Finish()

	return catalog_item.CatalogItemGet(ctx, arg,
		external.Interfaces().CPDBIfc(), internal.Interfaces().UuidIfc())
}

func (s *MarinaServer) CatalogItemDelete(ctx context.Context, arg *marinaIfc.CatalogItemDeleteArg) (
	*marinaIfc.CatalogItemDeleteRet, error) {

	taskUuid, err := s.asyncHandler(ctx, arg, CatalogItemDelete)
	if err != nil {
		return nil, err
	}

	ret := &marinaIfc.CatalogItemDeleteRet{
		TaskUuid: taskUuid,
	}
	return ret, nil
}

func (s *MarinaServer) CatalogItemCreate(ctx context.Context, arg *marinaIfc.CatalogItemCreateArg) (
	*marinaIfc.CatalogItemCreateRet, error) {

	taskUuid, err := s.asyncHandler(ctx, arg, CatalogItemCreate)
	if err != nil {
		return nil, err
	}

	ret := &marinaIfc.CatalogItemCreateRet{
		TaskUuid: taskUuid,
	}
	return ret, nil
}

func (s *MarinaServer) CatalogItemUpdate(ctx context.Context, arg *marinaIfc.CatalogItemUpdateArg) (
	*marinaIfc.CatalogItemUpdateRet, error) {

	taskUuid, err := s.asyncHandler(ctx, arg, CatalogItemUpdate)
	if err != nil {
		return nil, err
	}

	ret := &marinaIfc.CatalogItemUpdateRet{
		TaskUuid: taskUuid,
	}
	return ret, nil
}
