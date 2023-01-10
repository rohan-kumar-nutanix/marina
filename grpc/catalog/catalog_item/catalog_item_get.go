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

	log "k8s.io/klog/v2"

	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	utils "github.com/nutanix-core/content-management-marina/util"
	"github.com/nutanix-core/ntnx-api-utils-go/tracer"
)

// CatalogItemGet implements the CatalogItemGet RPC.
func CatalogItemGet(ctx context.Context, arg *marinaIfc.CatalogItemGetArg, catalogItemIfc CatalogItemInterface,
	cpdbIfc cpdb.CPDBClientInterface, uuidIfc utils.UuidUtilInterface) (*marinaIfc.CatalogItemGetRet, error) {
	log.V(2).Info("CatalogItemGet RPC started.")
	span, ctx := tracer.StartSpan(ctx, "CatalogItemGet")
	defer span.Finish()

	var catalogItemList []*marinaIfc.CatalogItemInfo
	catalogItemIdList := arg.GetCatalogItemIdList()
	catalogItemIdListSize := len(catalogItemIdList)
	catalogItemChan := make(chan []*marinaIfc.CatalogItemInfo)
	catalogItemErrChan := make(chan error)
	if catalogItemIdListSize <= *CatalogIdfQueryChunkSize {
		go catalogItemIfc.GetCatalogItemsChan(ctx, cpdbIfc, uuidIfc, catalogItemIdList, arg.GetCatalogItemTypeList(),
			arg.GetLatest(), catalogItemChan, catalogItemErrChan)

		catalogItemList = <-catalogItemChan
		err := <-catalogItemErrChan

		if err != nil {
			log.Errorf("Error while fetching Catalog Item list: %v", err)
			return nil, err
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
			go catalogItemIfc.GetCatalogItemsChan(ctx, cpdbIfc, uuidIfc, catalogItemIdList[start:end],
				arg.GetCatalogItemTypeList(), arg.GetLatest(), catalogItemChan, catalogItemErrChan)
		}

		for i := 0; i < count; i++ {
			catalogItemList = append(catalogItemList, <-catalogItemChan...)
			err := <-catalogItemErrChan

			if err != nil {
				log.Errorf("Error while fetching CatalogItem list: %v", err)
				return nil, err
			}
		}
	}

	ret := &marinaIfc.CatalogItemGetRet{CatalogItemList: catalogItemList}
	log.V(2).Info("CatalogItemGet RPC finished.")
	return ret, nil
}
