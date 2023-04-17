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
	"github.com/nutanix-core/content-management-marina/interface/external"
	internal "github.com/nutanix-core/content-management-marina/interface/local"
	contentPB "github.com/nutanix-core/content-management-marina/protos/apis/cms/v4/content"
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

	ret := &contentPB.GetWarehouseRet{
		Content: &contentPB.GetWarehouseResponse{
			Data: &contentPB.GetWarehouseResponse_WarehouseData{
				WarehouseData: &contentPB.WarehouseWrapper{
					Value: warehouse,
				},
			},
		},
	}
	return ret, nil

}

func WarehouseList(ctx context.Context, arg *contentPB.ListWarehousesArg) (*contentPB.ListWarehousesRet, error) {
	warehouses, err := newWarehouseDBImpl().ListWarehouses(ctx, external.Interfaces().CPDBIfc())
	if err != nil {
		return nil, err
	}

	ret := &contentPB.ListWarehousesRet{
		Content: &contentPB.ListWarehousesResponse{
			Data: &contentPB.ListWarehousesResponse_WarehouseArrayData{
				WarehouseArrayData: &contentPB.WarehouseArrayWrapper{
					Value: warehouses,
				},
			},
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
