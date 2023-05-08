/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Authors: rajesh.battala@nutanix.com
*
* Implementation of Get and List of Warehouse Security Policy entity.
 */

package security_policy

import (
	"context"

	"google.golang.org/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/tracer"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"

	"github.com/nutanix-core/content-management-marina/errors"
	"github.com/nutanix-core/content-management-marina/interface/external"
	internal "github.com/nutanix-core/content-management-marina/interface/local"
	configPB "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/config"
	"github.com/nutanix-core/content-management-marina/protos/apis/common/v1/response"
	utils "github.com/nutanix-core/content-management-marina/util"
)

func SecurityPolicyGetWithGrpcStatus(ctx context.Context, arg *configPB.GetSecurityPolicyByExtIdArg) (
	*configPB.GetSecurityPolicyByExtIdRet, error) {
	span, ctx := tracer.StartSpan(ctx, "GetSecurityPolicyDetails")
	defer span.Finish()
	ret, err := SecurityPolicyGet(ctx, arg)
	log.Infof("Returning Security Policy Entity Details Ret %v", ret)

	if err != nil {
		ret = &configPB.GetSecurityPolicyByExtIdRet{
			Content: &configPB.GetSecurityPolicyResponse{
				Data: &configPB.GetSecurityPolicyResponse_ErrorResponseData{
					ErrorResponseData: utils.GetConfigErrorResponse(nil),
				},
			},
		}
		return nil, internal.Interfaces().ErrorIfc().BuildGrpcError(errors.ErrInternal)
	}

	return ret, nil

}

func SecurityPolicyGet(ctx context.Context, arg *configPB.GetSecurityPolicyByExtIdArg) (
	*configPB.GetSecurityPolicyByExtIdRet, error) {
	securityPolicyUUID, _ := uuid4.StringToUuid4(*arg.ExtId)
	// TODO: check if UUID has to be passed as uuid4 type or String.
	securityPolicy, err := newSecurityPolicyDBImpl().GetSecurityPolicy(ctx, external.Interfaces().CPDBIfc(),
		securityPolicyUUID)
	if err != nil {
		return nil, err
	}

	metadata, err := getSPGetMetadata(ctx, *arg.ExtId)

	if err != nil {
		log.Errorf("Error occurred while generating Metadata %s %v", metadata, err)
	}

	/*link := &response.ApiLink{
		Href: proto.String("http://localhost/href"),
		Rel:  proto.String("http://localhost/Rel"),
	}*/

	/*	base := &response.ExternalizableAbstractModel{
			TenantInfo: nil,
			ExtId:      proto.String(securityPolicyUUID.String()),
			Links:      []*response.ApiLink{link},
		}
		log.Infof("Base %v", base)*/

	ret := &configPB.GetSecurityPolicyByExtIdRet{
		Content: &configPB.GetSecurityPolicyResponse{
			Data: &configPB.GetSecurityPolicyResponse_SecurityPolicyData{
				SecurityPolicyData: &configPB.SecurityPolicyWrapper{
					Value: securityPolicy,
				},
			},
			// Metadata: metadata,
			/*Metadata: &response.ApiResponseMetadata{
				Flags:                 utils.GetMetadataFlags(),
				Links:                 []*response.ApiLink{link},
				TotalAvailableResults: proto.Int32(1),
				Messages:              nil,
				ExtraInfo:             nil,
			},*/
		},
	}
	return ret, nil

}

func SecurityPolicyList(ctx context.Context, arg *configPB.ListSecurityPoliciesArg) (
	*configPB.ListSecurityPoliciesRet, error) {
	securityPolices, err := newSecurityPolicyDBImpl().ListSecurityPolices(ctx, external.Interfaces().CPDBIfc())
	if err != nil {
		return nil, err
	}

	/*m, err := getSPGetMetadata(ctx, "04c2c4ba-7516-480e-4eb4-520b3e38214e")
	if err != nil {
		log.Errorf("Error occurred while generating Metadata %s %v", m, err)
	}*/

	ret := &configPB.ListSecurityPoliciesRet{
		Content: &configPB.ListSecurityPoliciesResponse{
			Data: &configPB.ListSecurityPoliciesResponse_SecurityPolicyArrayData{
				SecurityPolicyArrayData: &configPB.SecurityPolicyArrayWrapper{
					Value: securityPolices,
				},
			},
			// Metadata: m,
		},
	}

	// Metadata: &response.ApiResponseMetadata{
	// 	TotalAvailableResults: proto.Int32(totalEntities),
	// },
	// },
	// ResponseCode: proto.String("200"),
	// }

	return ret, nil
}

func SecurityPolicyListWithGrpcStatus(ctx context.Context, arg *configPB.ListSecurityPoliciesArg) (
	*configPB.ListSecurityPoliciesRet, error) {
	span, ctx := tracer.StartSpan(ctx, "ListSecurityPoliciesDetails")
	defer span.Finish()
	ret, err := SecurityPolicyList(ctx, arg)
	log.V(2).Infof("Returning List of Security Policies Entity Details Ret %v, %s", ret, err)

	if err != nil {
		ret = &configPB.ListSecurityPoliciesRet{
			Content: &configPB.ListSecurityPoliciesResponse{
				Data: &configPB.ListSecurityPoliciesResponse_ErrorResponseData{
					ErrorResponseData: utils.GetConfigErrorResponse(nil),
				},
			},
		}
		return ret, internal.Interfaces().ErrorIfc().BuildGrpcError(errors.ErrInternal)
	}
	return ret, nil
}

// Get metadata corresponding to the task get response.
func getSPGetMetadata(ctx context.Context,
	securityPolicyUUID string) (*response.ApiResponseMetadata, error) {
	metadata := new(response.ApiResponseMetadata)
	metadata.Links = new(response.ApiLinkArrayWrapper)
	metadata.Links.Value = make([]*response.ApiLink, 0, 2)

	// Populate the self link.
	selfLink := &response.ApiLink{}
	href, err := utils.GetLink(ctx, "security-policies", securityPolicyUUID)
	if err != nil {
		log.Errorf("Error occurred in getting links %v", err)
		return nil, err
	}
	log.Infof("href generated %v", href)
	selfLink.Href = proto.String(href) // proto.String("https://10.96.208.32:9440/api/cms/v4.0.a1/config/security-policies/6d377c72-388d-4582-5a1a-81a7970e8eb8") // href)
	selfLink.Rel = proto.String(utils.SELF_RELATION)
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
	log.Infof("SP Metadata :%v", metadata1)
	return metadata1, nil
}
