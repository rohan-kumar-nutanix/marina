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

	"google.golang.org/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
	util "github.com/nutanix-core/content-management-marina/util"
)

type IdfClient struct {
	IdfSvc insights_interface.InsightsServiceInterface
	Retry  *misc.ExponentialBackoff
}

func (idf *IdfClient) IdfService() insights_interface.InsightsServiceInterface {
	return idf.IdfSvc
}

// Query IDF and returns entities which match the query.
// If none matches, return ErrNotFound.
func (idf *IdfClient) Query(
	ctx context.Context,
	arg *insights_interface.GetEntitiesWithMetricsArg) (
	[]*insights_interface.EntityWithMetric, error) {

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

func QueryIDF(ctx context.Context, idfService insights_interface.InsightsServiceInterface,
	query *insights_interface.Query) (
	*insights_interface.GetEntitiesWithMetricsRet,
	error) {

	idfArg := &insights_interface.GetEntitiesWithMetricsArg{Query: query}
	idfResponse := &insights_interface.GetEntitiesWithMetricsRet{}

	err := idfService.SendMsgWithTimeout(
		"GetEntitiesWithMetrics",
		idfArg,
		idfResponse,
		nil,                /* backoff */
		util.IdfRpcTimeOut, /* timeoutSecs */
	)

	if err != nil {
		log.Errorf("IDF query failed. Error - %s", err)
		return nil, err
	}
	return idfResponse, nil
}

func (idf *IdfClient) GetEntities(
	ctx context.Context,
	entityGuidList []*insights_interface.EntityGuid, metaDataOnly bool) (
	[]*insights_interface.Entity, error) {
	// Get entities specified in entity guid list.

	entityArg := &insights_interface.GetEntitiesArg{
		EntityGuidList:         entityGuidList,
		MetaDataOnly:           proto.Bool(metaDataOnly),
		IncludeDeletedEntities: proto.Bool(false),
	}
	entityRet := &insights_interface.GetEntitiesRet{}
	err := idf.IdfSvc.SendMsg("GetEntities", entityArg, entityRet, idf.Retry)
	if err != nil {
		return nil, err
	}
	entities := entityRet.GetEntity()
	return entities, err
}
