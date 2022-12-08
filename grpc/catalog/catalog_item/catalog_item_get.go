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

	log "k8s.io/klog/v2"

	"github.com/nutanix-core/content-management-marina/db"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	"github.com/nutanix-core/ntnx-api-utils-go/tracer"
)

// CatalogItemGet implements the CatalogItemGet RPC.
func CatalogItemGet(ctx context.Context, arg *marinaIfc.CatalogItemGetArg, idfIfc db.IdfClientInterface) (*marinaIfc.CatalogItemGetRet, error) {
	log.V(2).Info("CatalogItemGet RPC started.")
	span, ctx := tracer.StartSpan(ctx, "CatalogItemGet")
	defer span.Finish()

	var catalogItemList []*marinaIfc.CatalogItemInfo
	catalogItemIdList := arg.GetCatalogItemIdList()
	catalogItemIdListSize := len(catalogItemIdList)
	catalogItemChan := make(chan []*marinaIfc.CatalogItemInfo)
	catalogItemErrChan := make(chan error)
	if catalogItemIdListSize <= *CatalogIdfQueryChunkSize {

		go GetCatalogItemsChan(ctx, idfIfc, catalogItemIdList, arg.GetCatalogItemTypeList(), arg.GetLatest(),
			catalogItemChan, catalogItemErrChan)

		catalogItemList = <-catalogItemChan
		err := <-catalogItemErrChan

		if err != nil {
			log.Error("Error occurred : ", err)
			errMsg := fmt.Sprintf("Error while fetching CatalogItem list: %v", err)
			return nil, marinaError.ErrInternal.SetCause(errors.New(errMsg))
		}

	} else {
		count := 0
		for start := 0; start < catalogItemIdListSize; start += *CatalogIdfQueryChunkSize {
			count++
			end := start + *CatalogIdfQueryChunkSize
			if end > catalogItemIdListSize {
				end = catalogItemIdListSize
			}

			log.Infof("Fetching catalog items from index %v to %v", start, end)
			go GetCatalogItemsChan(ctx, idfIfc, catalogItemIdList[start:end], arg.GetCatalogItemTypeList(), arg.GetLatest(),
				catalogItemChan, catalogItemErrChan)
		}

		for i := 0; i < count; i++ {
			catalogItemList = append(catalogItemList, <-catalogItemChan...)
			err := <-catalogItemErrChan

			if err != nil {
				log.Error("Error occurred : ", err)
				errMsg := fmt.Sprintf("Error while fetching CatalogItem list: %v", err)
				return nil, marinaError.ErrInternal.SetCause(errors.New(errMsg))
			}
		}
	}

	ret := &marinaIfc.CatalogItemGetRet{CatalogItemList: catalogItemList}
	log.V(2).Info("CatalogItemGet RPC finished.")
	return ret, nil
}
