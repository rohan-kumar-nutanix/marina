/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 * Marina internal interfaces singleton
 *
 */

package local

import (
	"sync"

	"github.com/nutanix-core/content-management-marina/authz"
	"github.com/nutanix-core/content-management-marina/grpc/catalog/catalog_item"
	"github.com/nutanix-core/content-management-marina/metadata"
	utils "github.com/nutanix-core/content-management-marina/util"
)

type MarinaInternalInterfaces interface {
	AuthzIfc() authz.AuthzInterface
	CatalogItemIfc() catalog_item.CatalogItemInterface
	FanoutTaskPollerIfc() utils.FanoutTaskPollerInterface
	MetadataIfc() metadata.EntityMetadataInterface
	ProtoIfc() utils.ProtoUtilInterface
	UuidIfc() utils.UuidUtilInterface
}

type singletonObject struct {
	authzIfc            authz.AuthzInterface
	catalogItemIfc      catalog_item.CatalogItemInterface
	fanoutTaskPollerIfc utils.FanoutTaskPollerInterface
	metadataIfc         metadata.EntityMetadataInterface
	protoIfc            utils.ProtoUtilInterface
	uuidIfc             utils.UuidUtilInterface
}

var (
	singleton            MarinaInternalInterfaces
	singletonServiceOnce sync.Once
)

// InitSingletonService - Initialize a singleton Marina service.
func InitSingletonService() {
	singletonServiceOnce.Do(func() {
		singleton = &singletonObject{
			authzIfc:            new(authz.AuthzUtil),
			catalogItemIfc:      new(catalog_item.CatalogItemImpl),
			fanoutTaskPollerIfc: new(utils.FanoutTaskPollerUtil),
			metadataIfc:         new(metadata.EntityMetadataUtil),
			protoIfc:            new(utils.ProtoUtil),
			uuidIfc:             new(utils.UuidUtil),
		}
	})
}

// GetSingletonServiceWithParams - Initialize a singleton Marina service with params. Should only be used in UTs
func GetSingletonServiceWithParams(authzService authz.AuthzInterface,
	catalogItemService catalog_item.CatalogItemInterface, metadataService metadata.EntityMetadataInterface,
	poller utils.FanoutTaskPollerInterface, protoService utils.ProtoUtilInterface,
	uuidService utils.UuidUtilInterface) *singletonObject {

	return &singletonObject{
		authzIfc:            authzService,
		catalogItemIfc:      catalogItemService,
		fanoutTaskPollerIfc: poller,
		metadataIfc:         metadataService,
		protoIfc:            protoService,
		uuidIfc:             uuidService,
	}
}

// Interfaces - Returns the singleton for MarinaExternalInterfaces
func Interfaces() MarinaInternalInterfaces {
	return singleton
}

// AuthzIfc - Returns the singleton for AuthzInterface
func (s *singletonObject) AuthzIfc() authz.AuthzInterface {
	return s.authzIfc
}

// CatalogItemIfc - Returns the singleton for CatalogItemInterface
func (s *singletonObject) CatalogItemIfc() catalog_item.CatalogItemInterface {
	return s.catalogItemIfc
}

// FanoutTaskPollerIfc returns the singleton for FanoutTaskPollerInterface
func (s *singletonObject) FanoutTaskPollerIfc() utils.FanoutTaskPollerInterface {
	return s.fanoutTaskPollerIfc
}

// MetadataIfc - Returns the singleton for EntityMetadataInterface
func (s *singletonObject) MetadataIfc() metadata.EntityMetadataInterface {
	return s.metadataIfc
}

// ProtoIfc - Returns the singleton for ProtoUtilInterface
func (s *singletonObject) ProtoIfc() utils.ProtoUtilInterface {
	return s.protoIfc
}

// UuidIfc returns the singleton for UuidUtilInterface
func (s *singletonObject) UuidIfc() utils.UuidUtilInterface {
	return s.uuidIfc
}
