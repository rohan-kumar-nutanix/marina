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
	"context"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
)

type IdfClient struct {
	IdfSvc insights_interface.InsightsServiceInterface
	Retry  *misc.ExponentialBackoff
}

// Query IDF and returns entities which match the query.
// If none matches, return ErrNotFound.
func (idf *IdfClient) Query(ctx context.Context, query *insights_interface.Query) (
	[]*insights_interface.EntityWithMetric, error) {
	arg := &insights_interface.GetEntitiesWithMetricsArg{Query: query}
	ret := &insights_interface.GetEntitiesWithMetricsRet{}
	err := idf.IdfSvc.SendMsg("GetEntitiesWithMetrics", arg, ret, idf.Retry)
	if err != nil {
		return nil, err
	}

	grpResults := ret.GetGroupResultsList()
	if len(grpResults) == 0 {
		return nil, insights_interface.ErrNotFound
	}

	entities := grpResults[0].GetRawResults()
	if len(entities) == 0 {
		return nil, insights_interface.ErrNotFound
	}

	return entities, err
}
