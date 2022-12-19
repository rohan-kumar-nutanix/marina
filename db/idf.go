/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 * IDF utilities for Marina Service.
 *
 */

package db

import (
	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
)

type IdfClient struct {
	IdfSvc insights_interface.InsightsServiceInterface
	Retry  *misc.ExponentialBackoff
}
