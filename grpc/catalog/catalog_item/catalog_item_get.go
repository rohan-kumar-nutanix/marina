/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 * The implementation for CatalogItemGet RPC
 */

package catalog_item

import (
	"context"
	"errors"
	"fmt"

	marinaError "github.com/nutanix-core/content-management-marina/error"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	"github.com/nutanix-core/ntnx-api-utils-go/tracer"
	log "k8s.io/klog/v2"
)

func CatalogItemGet(ctx context.Context, arg *marinaIfc.CatalogItemGetArg) (*marinaIfc.CatalogItemGetRet, error) {
	span, ctx := tracer.StartSpan(ctx, "GetCatalogItems")
	defer span.Finish()

	catalogItemChan := make(chan []*marinaIfc.CatalogItemInfo)
	catalogItemErrChan := make(chan error)

	go GetCatalogItemsChan(ctx, arg.GetCatalogItemIdList(), arg.GetCatalogItemTypeList(), arg.GetLatest(),
		catalogItemChan, catalogItemErrChan)

	catalogItemList := <-catalogItemChan
	err := <-catalogItemErrChan

	if err != nil {
		log.Error("Error occurred : ", err)
		errMsg := fmt.Sprintf("Error while fetching CatalogItem list: %v", err)
		return nil, marinaError.ErrInternal.SetCause(errors.New(errMsg))
	}

	ret := &marinaIfc.CatalogItemGetRet{CatalogItemList: catalogItemList}
	return ret, nil
}
