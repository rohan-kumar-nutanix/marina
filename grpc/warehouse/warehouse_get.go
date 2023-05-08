/*
* Copyright (c) 2022 Nutanix Inc. All rights reserved.
*
* Authors: rajesh.battala@nutanix.com
*
* Implementation of fetching Warehouse details.
 */

package warehouse

import (
	"context"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/tracer"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"google.golang.org/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/content-management-marina/errors"
	"github.com/nutanix-core/content-management-marina/interface/external"
	internal "github.com/nutanix-core/content-management-marina/interface/local"
	contentPB "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"
	apiResponse "github.com/nutanix-core/content-management-marina/protos/apis/common/v1/response"
	utils "github.com/nutanix-core/content-management-marina/util"
)

func WarehouseGetWithGrpcStatus(ctx context.Context, arg *contentPB.GetWarehouseArg) (*contentPB.GetWarehouseRet, error) {
	span, ctx := tracer.StartSpan(ctx, "GetWarehouseDetails")
	defer span.Finish()
	ret, err := WarehouseGet(ctx, arg)
	log.Infof("Returning Warehouse Entity Details Ret %v", ret)

	if err != nil {
		ret = &contentPB.GetWarehouseRet{
			Content: &contentPB.GetWarehouseResponse{
				Data: &contentPB.GetWarehouseResponse_ErrorResponseData{
					ErrorResponseData: utils.SetErrorResponse(nil),
				},
			},
		}
		return ret, internal.Interfaces().ErrorIfc().BuildGrpcError(errors.ErrInternal)
	}

	return ret, nil

}

func WarehouseGet(ctx context.Context, arg *contentPB.GetWarehouseArg) (*contentPB.GetWarehouseRet, error) {

	warehouseUUID, _ := uuid4.StringToUuid4(*arg.ExtId)
	// TODO: check if UUID has to be passed as uuid4 type or String.
	warehouse, err := newWarehouseDBImpl().GetWarehouse(ctx, external.Interfaces().CPDBIfc(), warehouseUUID)
	if err != nil {
		return nil, err
	}

	metadata, err := getWarehouseGetMetadata(ctx, *arg.ExtId)
	if err != nil {
		log.Errorf("error occurred in generating Warehouse Response Metadata %s", err)
	}

	log.Infof("Setting Warehouse Metadata %v", metadata)
	ret := &contentPB.GetWarehouseRet{
		Content: &contentPB.GetWarehouseResponse{
			Data: &contentPB.GetWarehouseResponse_WarehouseData{
				WarehouseData: &contentPB.WarehouseWrapper{
					Value: warehouse,
				},
			},
			Metadata: metadata,
		},
	}
	return ret, nil

}

func WarehouseList(ctx context.Context, arg *contentPB.ListWarehousesArg) (*contentPB.ListWarehousesRet, error) {
	warehouses, err := newWarehouseDBImpl().ListWarehouses(ctx, external.Interfaces().CPDBIfc())
	if err != nil {
		return nil, err
	}

	metadata, err := getWarehouseListMetadata(ctx, int32(len(warehouses)))
	if err != nil {
		log.Errorf("error occurred in generating Warehouse Response Metadata %s", err)
	}
	ret := &contentPB.ListWarehousesRet{
		Content: &contentPB.ListWarehousesResponse{
			Data: &contentPB.ListWarehousesResponse_WarehouseArrayData{
				WarehouseArrayData: &contentPB.WarehouseArrayWrapper{
					Value: warehouses,
				},
			},
			Metadata: metadata,
		},
	}
	return ret, nil
}

func WarehouseListWithGrpcStatus(ctx context.Context, arg *contentPB.ListWarehousesArg) (
	*contentPB.ListWarehousesRet, error) {
	span, ctx := tracer.StartSpan(ctx, "ListWarehouseDetails")
	defer span.Finish()
	ret, err := WarehouseList(ctx, arg)
	log.V(2).Infof("Returning Warehouse Entity Details Ret %v", ret)

	if err != nil {
		ret = &contentPB.ListWarehousesRet{
			Content: &contentPB.ListWarehousesResponse{
				Data: &contentPB.ListWarehousesResponse_ErrorResponseData{
					ErrorResponseData: utils.SetErrorResponse(nil),
				},
			},
		}
		return ret, internal.Interfaces().ErrorIfc().BuildGrpcError(errors.ErrInternal)
	}
	return ret, nil
}

// Get metadata corresponding to the task get response.
func getWarehouseGetMetadata(ctx context.Context,
	warehouseUUID string) (*apiResponse.ApiResponseMetadata, error) {
	metadata := new(apiResponse.ApiResponseMetadata)
	// metadata.Links = &apiResponse.ApiLinkArrayWrapper{}
	metadata.Links = new(apiResponse.ApiLinkArrayWrapper)
	metadata.Links.Value = make([]*apiResponse.ApiLink, 0, 1)

	// Populate the self link.
	selfLink := &apiResponse.ApiLink{}
	href, err := utils.GetLink(ctx, "warehouses", warehouseUUID)
	if err != nil {
		log.Errorf("Error occurred in getting links %v", err)
		// return nil, err
	}
	log.Infof("href generated %v", href)
	selfLink.Href = proto.String(href)
	selfLink.Rel = proto.String(utils.SELF_RELATION)

	metadata.Links.Value = append(metadata.Links.Value, selfLink)

	metadata1 := &apiResponse.ApiResponseMetadata{
		Flags:                 utils.GetMetadataFlagWrapper(),
		Links:                 metadata.Links,
		TotalAvailableResults: proto.Int32(1),
	}
	log.Infof("Warehouse Metadata1 :%v", metadata1)
	return metadata1, nil
}

func getWarehouseListMetadata(ctx context.Context,
	count int32) (*apiResponse.ApiResponseMetadata, error) {
	metadata := new(apiResponse.ApiResponseMetadata)
	// metadata.Links = &apiResponse.ApiLinkArrayWrapper{}
	metadata.Links = new(apiResponse.ApiLinkArrayWrapper)
	metadata.Links.Value = make([]*apiResponse.ApiLink, 0, 1)

	// Populate the self link.
	selfLink := &apiResponse.ApiLink{}
	href, _ := utils.GetLink(ctx, "warehouses", "")
	// if err != nil {
	// 	log.Errorf("Error occurred in getting links %v", err)
	// 	return nil, err
	// }
	log.Infof("href generated %v", href)
	selfLink.Href = proto.String(href)
	selfLink.Rel = proto.String(utils.SELF_RELATION)
	// metadata.Links = append(metadata.Links, selfLink)

	metadata.Links.Value = append(metadata.Links.Value, selfLink)

	metadata1 := &apiResponse.ApiResponseMetadata{
		Flags:                 utils.GetMetadataFlagWrapper(),
		Links:                 metadata.Links,
		TotalAvailableResults: proto.Int32(count),
	}
	log.Infof("Warehouse Metadata1 :%v", metadata1)
	return metadata1, nil
}
