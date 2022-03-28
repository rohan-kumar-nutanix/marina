/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 */

package services

import (
	"context"
	"github.com/nutanix-core/content-management-marina/grpc/catalog/catalog_item"
	marinapb "github.com/nutanix-core/content-management-marina/protos/marina"
	"github.com/nutanix-core/ntnx-api-utils-go/tracer"
	log "k8s.io/klog/v2"
)

type MarinaServer struct {
	marinapb.UnimplementedMarinaServer
}

type CatalogItemInterface interface {
	CatalogItemGet(ctx context.Context, arg *marinapb.CatalogItemGetArg) (*marinapb.CatalogItemGetRet, error)
}

func (s *MarinaServer) CatalogItemGet(ctx context.Context, arg *marinapb.CatalogItemGetArg) (*marinapb.CatalogItemGetRet, error) {
	span, ctx := tracer.StartSpan(ctx, "catalogitem-get")
	defer span.Finish()

	log.Infof("Arg received %v", arg)
	return catalog_item.GetCatalogItems(ctx, arg)
}
