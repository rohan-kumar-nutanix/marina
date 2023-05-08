/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 * Security Policy Server for Managing the Warehouse Security Policies.
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
	"github.com/nutanix-core/content-management-marina/grpc/security_policy"
	internalIFC "github.com/nutanix-core/content-management-marina/interface/local"
	configPB "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/config"
	prismPB "github.com/nutanix-core/content-management-marina/protos/apis/prism/v4/config"
	util "github.com/nutanix-core/content-management-marina/util"
)

type SecurityPolicyServer struct {
	configPB.UnimplementedSecurityPolicyServiceServer
}

type SecurityPolicy interface {
	CreateSecurityPolicy(context.Context, *configPB.CreateSecurityPolicyArg) (*configPB.CreateSecurityPolicyRet, error)
	DeleteSecurityPolicyByExtId(context.Context, *configPB.DeleteSecurityPolicyByExtIdArg) (*configPB.DeleteSecurityPolicyByExtIdRet, error)
	GetSecurityPolicyByExtId(context.Context, *configPB.GetSecurityPolicyByExtIdArg) (*configPB.GetSecurityPolicyByExtIdRet, error)
	ListSecurityPolicies(context.Context, *configPB.ListSecurityPoliciesArg) (*configPB.ListSecurityPoliciesRet, error)
	UpdateSecurityPolicyByExtId(context.Context, *configPB.UpdateSecurityPolicyByExtIdArg) (*configPB.UpdateSecurityPolicyByExtIdRet, error)
}

// CreateSecurityPolicy gRPC Handler.
func (s *SecurityPolicyServer) CreateSecurityPolicy(ctx context.Context, arg *configPB.CreateSecurityPolicyArg) (
	*configPB.CreateSecurityPolicyRet, error) {
	taskUuid, err := s.asyncHandlerWithGrpcStatus(ctx, arg, CreateSecurityPolicy)
	if err != nil {
		ret := &configPB.CreateSecurityPolicyRet{
			Content: &configPB.CreateSecurityPolicyResponse{
				Data: &configPB.CreateSecurityPolicyResponse_ErrorResponseData{
					ErrorResponseData: util.GetConfigErrorResponse(nil),
				},
			},
		}
		return ret, err
	}
	ret := &configPB.CreateSecurityPolicyRet{
		Content: &configPB.CreateSecurityPolicyResponse{
			Data: &configPB.CreateSecurityPolicyResponse_TaskReferenceData{
				TaskReferenceData: s.getTaskData(taskUuid),
			},
		},
	}
	return ret, nil
}

// GetSecurityPolicyByExtId gRPC Handler
func (s *SecurityPolicyServer) GetSecurityPolicyByExtId(ctx context.Context, arg *configPB.GetSecurityPolicyByExtIdArg) (
	*configPB.GetSecurityPolicyByExtIdRet, error) {
	log.Infof("GetSecurityPolicyByExtId Arg received %v", arg)
	defer log.Infof("Finished GetSecurityPolicyByExtId RPC for SecurityPolicy UUID: %s", *arg.ExtId)
	span := tracer.GetSpanFromContext(ctx)
	span.SetTag("SecurityPolicyUUID", *arg.ExtId)

	return security_policy.SecurityPolicyGetWithGrpcStatus(ctx, arg)
}

// ListSecurityPolicies gRPC Handler
func (s *SecurityPolicyServer) ListSecurityPolicies(ctx context.Context, arg *configPB.ListSecurityPoliciesArg) (
	*configPB.ListSecurityPoliciesRet, error) {
	log.Infof("ListSecurityPolicies Arg received %v", arg)
	// defer log.Infof("Finished ListSecurityPolicies RPC ")
	span := tracer.GetSpanFromContext(ctx)
	span.SetTag("ListSecurityPolicies", "")

	return security_policy.SecurityPolicyListWithGrpcStatus(ctx, arg)
}

// DeleteSecurityPolicyByExtId gRPC Handler
func (s *SecurityPolicyServer) DeleteSecurityPolicyByExtId(ctx context.Context, arg *configPB.DeleteSecurityPolicyByExtIdArg) (
	*configPB.DeleteSecurityPolicyByExtIdRet, error) {
	taskUuid, err := s.asyncHandlerWithGrpcStatus(ctx, arg, DeleteSecurityPolicyByExtId)
	if err != nil {
		ret := &configPB.DeleteSecurityPolicyByExtIdRet{
			Content: &configPB.DeleteSecurityPolicyResponse{
				Data: &configPB.DeleteSecurityPolicyResponse_ErrorResponseData{
					ErrorResponseData: util.GetConfigErrorResponse(nil),
				},
			},
		}
		return ret, err
	}
	ret := &configPB.DeleteSecurityPolicyByExtIdRet{
		Content: &configPB.DeleteSecurityPolicyResponse{
			Data: &configPB.DeleteSecurityPolicyResponse_TaskReferenceData{
				TaskReferenceData: s.getTaskData(taskUuid),
			},
		},
	}
	return ret, nil
}

// UpdateSecurityPolicyByExtId gRPC Handler
func (s *SecurityPolicyServer) UpdateSecurityPolicyByExtId(ctx context.Context, arg *configPB.UpdateSecurityPolicyByExtIdArg) (
	*configPB.UpdateSecurityPolicyByExtIdRet, error) {
	taskUuid, err := s.asyncHandlerWithGrpcStatus(ctx, arg, UpdateSecurityPolicyByExtId)
	if err != nil {
		ret := &configPB.UpdateSecurityPolicyByExtIdRet{
			Content: &configPB.UpdateSecurityPolicyResponse{
				Data: &configPB.UpdateSecurityPolicyResponse_ErrorResponseData{
					ErrorResponseData: util.GetConfigErrorResponse(nil),
				},
			},
		}
		return ret, err
	}
	ret := &configPB.UpdateSecurityPolicyByExtIdRet{
		Content: &configPB.UpdateSecurityPolicyResponse{
			Data: &configPB.UpdateSecurityPolicyResponse_TaskReferenceData{
				TaskReferenceData: s.getTaskData(taskUuid),
			},
		},
	}
	return ret, nil
}

func (s *SecurityPolicyServer) asyncHandlerWithGrpcStatus(c context.Context, request proto.Message,
	operation string) (*uuid4.Uuid, error) {
	// TODO add marinaError.AppMessage for APIs
	taskUuid, err := s.asyncHandler(c, request, operation)
	if err != nil {
		// TODO: Create gRPC error from MarinaError and generate AppMessage Error and return to client.
		return taskUuid, internalIFC.Interfaces().ErrorIfc().BuildGrpcError(err)
	}
	return taskUuid, nil
}

func (s *SecurityPolicyServer) asyncHandler(c context.Context, request proto.Message,
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

func (s *SecurityPolicyServer) getTaskData(taskUuid *uuid4.Uuid) *configPB.TaskReferenceWrapper {
	return &configPB.TaskReferenceWrapper{
		Value: &prismPB.TaskReference{
			ExtId: proto.String(taskUuid.String()),
		},
	}
}
