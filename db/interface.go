/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 * The IDF interface
 */

package db

import (
	"context"

	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"

	utils "github.com/nutanix-core/content-management-marina/util"
)

type IdfClientInterface interface {
	DeleteEntities(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, entityType EntityType,
		entityUuids []string, isCasEnabled bool) error
	GetEntitiesWithMetrics(ctx context.Context, arg *insights_interface.GetEntitiesWithMetricsArg,
		ret *insights_interface.GetEntitiesWithMetricsRet) error
}

func newIdfClientWithRetry(insightsIfc insights_interface.InsightsServiceInterface) *IdfClient {
	idfClient := IdfClient{IdfSvc: insightsIfc}
	return &idfClient
}

func InsightsServiceInterface() insights_interface.InsightsServiceInterface {
	insightsIfc := insights_interface.NewInsightsServiceInterface(
		utils.HostAddr,
		uint16(*insights_interface.InsightsPort))

	err := insightsIfc.SetRequestTimeout(utils.IdfRpcTimeOut)
	if err != nil {
		log.Fatalf("Error occurred while setting IdfRpcTimeOut : %v", err)
	}
	return insightsIfc
}

func IdfClientWithRetry(insightsIfc insights_interface.InsightsServiceInterface) IdfClientInterface {
	return newIdfClientWithRetry(insightsIfc)
}

func IdfClientWithoutRetry(insightsIfc insights_interface.InsightsServiceInterface) IdfClientInterface {
	idfClientObj := newIdfClientWithRetry(insightsIfc)
	idfClientObj.Retry = misc.NewExponentialBackoff(0, 0, 0)
	return idfClientObj
}
