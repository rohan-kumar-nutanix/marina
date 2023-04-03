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
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/tracer"

	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	utils "github.com/nutanix-core/content-management-marina/util"
)

// CatalogItemGet implements the CatalogItemGet RPC.
func CatalogItemGet(ctx context.Context, arg *marinaIfc.CatalogItemGetArg, cpdbIfc cpdb.CPDBClientInterface,
	uuidIfc utils.UuidUtilInterface) (*marinaIfc.CatalogItemGetRet, error) {

	log.V(2).Info("CatalogItemGet RPC started.")
	span, ctx := tracer.StartSpan(ctx, "CatalogItemGet")
	defer span.Finish()
	defer log.V(2).Info("CatalogItemGet RPC finished.")

	var catalogItems []*marinaIfc.CatalogItemInfo
	catalogItemIds := arg.GetCatalogItemIdList()
	catalogItemIdsSize := len(catalogItemIds)
	catalogItemChan := make(chan []*marinaIfc.CatalogItemInfo)
	catalogItemErrChan := make(chan error)
	if catalogItemIdsSize <= *CatalogIdfQueryChunkSize {
		go newCatalogItemImpl().GetCatalogItemsChan(ctx, cpdbIfc, uuidIfc, catalogItemIds, arg.GetCatalogItemTypeList(),
			arg.GetLatest(), catalogItemChan, catalogItemErrChan)

		catalogItems = <-catalogItemChan
		err := <-catalogItemErrChan

		if err != nil {
			log.Errorf("Error while fetching Catalog Item list: %v", err)
			return nil, err
		}

	} else {
		count := 0
		for start := 0; start < catalogItemIdsSize; start += *CatalogIdfQueryChunkSize {
			count++
			end := start + *CatalogIdfQueryChunkSize
			if end > catalogItemIdsSize {
				end = catalogItemIdsSize
			}

			log.Infof("Fetching catalog items from index %v to %v", start, end)
			go newCatalogItemImpl().GetCatalogItemsChan(ctx, cpdbIfc, uuidIfc, catalogItemIds[start:end],
				arg.GetCatalogItemTypeList(), arg.GetLatest(), catalogItemChan, catalogItemErrChan)
		}

		for i := 0; i < count; i++ {
			catalogItems = append(catalogItems, <-catalogItemChan...)
			err := <-catalogItemErrChan

			if err != nil {
				log.Errorf("Error while fetching CatalogItem list: %v", err)
				return nil, err
			}
		}
	}

	ret := &marinaIfc.CatalogItemGetRet{CatalogItemList: catalogItems}
	return ret, nil
}
