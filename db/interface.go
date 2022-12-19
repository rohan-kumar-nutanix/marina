/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 * The IDF interface
 */

package db

import (
	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
	utils "github.com/nutanix-core/content-management-marina/util"
)

type IdfClientInterface interface {
}

func newIdfClientWithRetry() *IdfClient {
	insightsSvc := insights_interface.NewInsightsServiceInterface(
		utils.HostAddr,
		uint16(*insights_interface.InsightsPort))

	insightsSvc.SetRequestTimeout(utils.IdfRpcTimeOut)
	idfClient := IdfClient{IdfSvc: insightsSvc}
	return &idfClient
}

func IdfClientWithRetry() IdfClientInterface {
	return newIdfClientWithRetry()
}

func IdfClientWithoutRetry() IdfClientInterface {
	idfClientObj := newIdfClientWithRetry()
	idfClientObj.Retry = misc.NewExponentialBackoff(0, 0, 0)
	return idfClientObj
}
