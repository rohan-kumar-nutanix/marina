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
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/content-management-marina/errors"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
	"github.com/nutanix-core/content-management-marina/interface/external"
	internal "github.com/nutanix-core/content-management-marina/interface/local"
	contentPB "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"
	utils "github.com/nutanix-core/content-management-marina/util"
)

func WarehouseItemGetWithGrpcStatus(ctx context.Context, arg *contentPB.GetWarehouseItemByIdArg) (*contentPB.GetWarehouseItemByIdRet, error) {
	span, ctx := tracer.StartSpan(ctx, "GetWarehouseDetails")
	defer span.Finish()
	ret, err := WarehouseItemGet(ctx, arg)
	log.Infof("Returning WarehouseItem Entity Details Ret %v", ret)

	if err != nil {
		ret = &contentPB.GetWarehouseItemByIdRet{
			Content: &contentPB.GetWarehouseItemResponse{
				Data: &contentPB.GetWarehouseItemResponse_ErrorResponseData{
					ErrorResponseData: utils.SetErrorResponse(nil),
				},
			},
		}
		return ret, internal.Interfaces().ErrorIfc().BuildGrpcError(errors.ErrInternal)
	}

	return ret, nil

}

func WarehouseItemGet(ctx context.Context, arg *contentPB.GetWarehouseItemByIdArg) (
	*contentPB.GetWarehouseItemByIdRet, error) {

	warehouseItemUUID, _ := uuid4.StringToUuid4(*arg.ExtId)
	// TODO: check if UUID has to be passed as uuid4 type or String.
	warehouseItem, err := newWarehouseDBImpl().GetWarehouseItem(ctx, external.Interfaces().CPDBIfc(), warehouseItemUUID)
	if err != nil {
		return nil, err
	}

	ret := &contentPB.GetWarehouseItemByIdRet{
		Content: &contentPB.GetWarehouseItemResponse{
			Data: &contentPB.GetWarehouseItemResponse_WarehouseItemData{
				WarehouseItemData: &contentPB.WarehouseItemWrapper{
					Value: warehouseItem,
				},
			},
		},
	}
	return ret, nil
}

func WarehouseItemList(ctx context.Context, arg *contentPB.ListWarehouseItemsArg) (*contentPB.ListWarehouseItemsRet, error) {
	warehouseUuid, _ := uuid4.StringToUuid4(*arg.ExtId)
	warehousesItems, err := newWarehouseDBImpl().ListWarehouseItems(
		ctx, external.Interfaces().CPDBIfc(), warehouseUuid)

	if err != nil && err != marinaError.ErrNotFound {
		return nil, err
	}

	ret := &contentPB.ListWarehouseItemsRet{
		Content: &contentPB.ListWarehouseItemsResponse{
			Data: &contentPB.ListWarehouseItemsResponse_WarehouseItemArrayData{
				WarehouseItemArrayData: &contentPB.WarehouseItemArrayWrapper{
					Value: warehousesItems,
				},
			},
		},
	}
	return ret, nil
}

func WarehouseItemListWithGrpcStatus(ctx context.Context, arg *contentPB.ListWarehouseItemsArg) (
	*contentPB.ListWarehouseItemsRet, error) {
	span, ctx := tracer.StartSpan(ctx, "ListWarehouseDetails")
	defer span.Finish()
	ret, err := WarehouseItemList(ctx, arg)
	log.V(2).Infof("Returning WarehouseItem Entities Ret %v", ret)

	if err != nil {
		ret = &contentPB.ListWarehouseItemsRet{
			Content: &contentPB.ListWarehouseItemsResponse{
				Data: &contentPB.ListWarehouseItemsResponse_ErrorResponseData{
					ErrorResponseData: utils.SetErrorResponse(nil),
				},
			},
		}
		return ret, internal.Interfaces().ErrorIfc().BuildGrpcError(errors.ErrInternal)
	}
	return ret, nil
}
