/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 * IDF utilities for Marina Service.
 *
 */

package idf

import (
	"context"

	"google.golang.org/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
	util "github.com/nutanix-core/content-management-marina/util"
)

const (
	CatalogDB             = "catalog_item_info"
	DefaultSortAttribute  = "created_time_usecs"
	Annotation            = "annotation"
	GlobalCatalogItemUuid = "global_catalog_item_uuid"
	CatalogItemType       = "item_type"
	CatalogName           = "name"
	CatalogItemUuid       = "uuid"
	CatalogVersion        = "version"
)

var CatalogItemAttributes []interface{} = []interface{}{
	Annotation,
	GlobalCatalogItemUuid,
	CatalogItemType,
	CatalogName,
	CatalogItemUuid,
	CatalogVersion,
}

type IdfClientInterface interface {
	Query(ctx context.Context, arg *insights_interface.GetEntitiesWithMetricsArg) (
		[]*insights_interface.EntityWithMetric, error)
	GetEntities(ctx context.Context, entityGuidList []*insights_interface.EntityGuid,
		metaDataOnly bool) ([]*insights_interface.Entity, error)
}

type IdfClient struct {
	IdfSvc insights_interface.InsightsServiceInterface
	Retry  *misc.ExponentialBackoff
}

func NewIdfClient() *IdfClient {
	insightsSvc := insights_interface.NewInsightsServiceInterface(
		util.HostAddr,
		uint16(*insights_interface.InsightsPort))
	// TODO check and move IdfRpcTimeout to gflag.
	insightsSvc.SetRequestTimeout(util.IdfRpcTimeOut)
	idfClient := IdfClient{IdfSvc: insightsSvc}
	return &idfClient
}

func NewIdfClientWithoutRetry() *IdfClient {
	idfClient := NewIdfClient()
	idfClient.Retry = misc.NewExponentialBackoff(0, 0, 0)
	return idfClient
}

func (idf *IdfClient) IdfService() insights_interface.InsightsServiceInterface {
	return idf.IdfSvc
}

func (idf *IdfClient) Query(
	ctx context.Context,
	arg *insights_interface.GetEntitiesWithMetricsArg) (
	[]*insights_interface.EntityWithMetric, error) {
	/*
	   Queries the IDF and returns entities which matches the following query
	   If none matches, it will return ErrNotFound
	*/

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
		util.IdfRpcTimeOut, /* timeoutSecs TODO: move it to Utils */
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
