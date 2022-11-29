/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 * The implementation for CatalogItemGet RPC
 */

package db

import (
	"context"
	"sync"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
	util "github.com/nutanix-core/content-management-marina/util"
)

var (
	idfClientWithRetryObj  IdfClientInterface
	idfClientWithRetryOnce sync.Once
	idfClientObj           IdfClientInterface
	idfClientOnce          sync.Once
)

type IdfClientInterface interface {
	Query(ctx context.Context, arg *insights_interface.GetEntitiesWithMetricsArg) (
		[]*insights_interface.EntityWithMetric, error)
	GetEntities(ctx context.Context, entityGuidList []*insights_interface.EntityGuid,
		metaDataOnly bool) ([]*insights_interface.Entity, error)
}

func newIdfClientWithRetry() *IdfClient {
	insightsSvc := insights_interface.NewInsightsServiceInterface(
		util.HostAddr,
		uint16(*insights_interface.InsightsPort))

	insightsSvc.SetRequestTimeout(util.IdfRpcTimeOut)
	idfClient := IdfClient{IdfSvc: insightsSvc}
	return &idfClient
}

func IdfClientWithRetry() IdfClientInterface {
	idfClientWithRetryOnce.Do(func() {
		idfClientWithRetryObj = newIdfClientWithRetry()
	})
	return idfClientWithRetryObj
}

func IdfClientWithoutRetry() IdfClientInterface {
	idfClientOnce.Do(func() {
		obj := newIdfClientWithRetry()
		obj.Retry = misc.NewExponentialBackoff(0, 0, 0)
		idfClientObj = obj
	})
	return idfClientObj
}
