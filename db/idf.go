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
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"

	marinaError "github.com/nutanix-core/content-management-marina/errors"
)

// EntityType - Definitions for Entity Type enum
type EntityType int64

const (
	CatalogItem EntityType = iota
	File
	Image
	RateLimit
	EntityCapability
)

func (entityType EntityType) ToString() string {
	switch entityType {
	case CatalogItem:
		return "catalog_item_info"
	case File:
		return "file_info"
	case Image:
		return "image_info"
	case RateLimit:
		return "catalog_rate_limit_info"
	case EntityCapability:
		return "abac_entity_capability"
	}
	return "unknown"
}

var (
	idfInitialInterval = 5
	idfMaxInterval     = 10
	idfMaxRetries      = 30
)

type IdfClient struct {
	IdfSvc insights_interface.InsightsServiceInterface
	Retry  *misc.ExponentialBackoff
}

func (idf *IdfClient) DeleteEntities(ctx context.Context, cpdbIfc cpdb.CPDBClientInterface, entityType EntityType,
	entityUuids []string, isCasEnabled bool) error {

	typeName := entityType.ToString()
	var guids []*insights_interface.EntityGuid
	for _, uuid := range entityUuids {
		guids = append(guids, &insights_interface.EntityGuid{
			EntityTypeName: proto.String(typeName),
			EntityId:       proto.String(uuid),
		})
	}

	entities, err := cpdbIfc.GetEntities(guids, false)
	if err != nil {
		log.Errorf("Getting entities for %v failed with error: %v", entityUuids, err)
		return err
	}

	if len(entities) == 0 {
		return nil
	}

	retryBackoff := misc.NewExponentialBackoff(time.Duration(idfInitialInterval)*time.Second,
		time.Duration(idfMaxInterval)*time.Second, idfMaxRetries)
	entitiesToDelete := entities
	for {
		arg := &insights_interface.BatchDeleteEntitiesArg{}
		for _, entity := range entitiesToDelete {
			entityArg := &insights_interface.DeleteEntityArg{EntityGuid: entity.GetEntityGuid()}
			if isCasEnabled {
				cas := entity.GetCasValue() + 1
				entityArg.CasValue = &cas
			}
			arg.EntityList = append(arg.EntityList, entityArg)
		}

		ret := &insights_interface.BatchDeleteEntitiesRet{}
		err = idf.IdfSvc.SendMsg("BatchDeleteEntities", arg, ret, idf.Retry)
		if err == nil {
			return nil
		}

		if err != insights_interface.ErrPartial {
			errMsg := fmt.Sprintf("Deleting batch %v entities failed with errType: %T error: %v",
				typeName, err, err.Error())
			return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
		}

		var failedEntities []*insights_interface.Entity
		for _, ret := range ret.GetRetList() {
			status := ret.GetStatus().GetErrorType()
			if status != insights_interface.InsightsErrorProto_kNoError {
				failedEntities = append(failedEntities, ret.GetEntity())
			}
		}

		if len(failedEntities) == 0 {
			errMsg := "Partial failure occurred but could not find entities to delete"
			return marinaError.ErrInternalError().SetCauseAndLog(errors.New(errMsg))
		}

		entitiesToDelete = failedEntities
		msg := fmt.Sprintf("%v entities for %v entity type: %v", len(entitiesToDelete), typeName, entitiesToDelete)
		if retryBackoff.Backoff() == misc.Stop {
			errMsg := fmt.Sprintf("Failed to delete %s. Retry limit reached!", msg)
			return marinaError.ErrRetry.SetCauseAndLog(errors.New(errMsg))
		}
		log.Warningf("Could not delete %s. Retrying...", msg)
	}
}

func (idf *IdfClient) GetEntitiesWithMetrics(ctx context.Context, arg *insights_interface.GetEntitiesWithMetricsArg,
	ret *insights_interface.GetEntitiesWithMetricsRet) error {

	err := idf.IdfSvc.SendMsg("GetEntitiesWithMetrics", arg, ret, idf.Retry)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get entities from IDF: %v", err)
		return marinaError.ErrRetry.SetCauseAndLog(errors.New(errMsg))
	}
	return nil
}
