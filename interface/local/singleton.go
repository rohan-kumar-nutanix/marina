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
	AuthzService() authz.AuthzInterface
	CatalogItemService() catalog_item.CatalogItemInterface
	FanoutTaskPollerService() utils.FanoutTaskPollerInterface
	MetadataService() metadata.EntityMetadataInterface
	ProtoService() utils.ProtoUtilInterface
	UuidService() utils.UuidUtilInterface
}

type singletonObject struct {
	authzService            authz.AuthzInterface
	catalogItemService      catalog_item.CatalogItemInterface
	fanoutTaskPollerService utils.FanoutTaskPollerInterface
	metadataService         metadata.EntityMetadataInterface
	protoService            utils.ProtoUtilInterface
	uuidService             utils.UuidUtilInterface
}

var (
	singleton            MarinaInternalInterfaces
	singletonServiceOnce sync.Once
)

// InitSingletonService - Initialize a singleton Marina service.
func InitSingletonService() {
	singletonServiceOnce.Do(func() {
		singleton = &singletonObject{
			authzService:            new(authz.AuthzUtil),
			catalogItemService:      new(catalog_item.CatalogItemImpl),
			fanoutTaskPollerService: new(utils.FanoutTaskPollerUtil),
			metadataService:         new(metadata.EntityMetadataUtil),
			protoService:            new(utils.ProtoUtil),
			uuidService:             new(utils.UuidUtil),
		}
	})
}

// GetSingletonServiceWithParams - Initialize a singleton Marina service with params. Should only be used in UTs
func GetSingletonServiceWithParams(authzService authz.AuthzInterface,
	catalogItemService catalog_item.CatalogItemInterface, metadataService metadata.EntityMetadataInterface,
	poller utils.FanoutTaskPollerInterface, protoService utils.ProtoUtilInterface,
	uuidService utils.UuidUtilInterface) *singletonObject {

	return &singletonObject{
		authzService:            authzService,
		catalogItemService:      catalogItemService,
		fanoutTaskPollerService: poller,
		metadataService:         metadataService,
		protoService:            protoService,
		uuidService:             uuidService,
	}
}

// Interfaces - Returns the singleton for MarinaExternalInterfaces
func Interfaces() MarinaInternalInterfaces {
	return singleton
}

// AuthzService - Returns the singleton for Authz interface.
func (s *singletonObject) AuthzService() authz.AuthzInterface {
	return s.authzService
}

// CatalogItemService - Returns the singleton for CatalogItemInterface
func (s *singletonObject) CatalogItemService() catalog_item.CatalogItemInterface {
	return s.catalogItemService
}

// FanoutTaskPollerService returns the singleton for FanoutTaskPollerInterface.
func (s *singletonObject) FanoutTaskPollerService() utils.FanoutTaskPollerInterface {
	return s.fanoutTaskPollerService
}

// MetadataService - Returns the singleton for Metadata Interface.
func (s *singletonObject) MetadataService() metadata.EntityMetadataInterface {
	return s.metadataService
}

// ProtoService - Returns the singleton for ProtoUtilInterface
func (s *singletonObject) ProtoService() utils.ProtoUtilInterface {
	return s.protoService
}

// UuidService returns the singleton for UuidUtilInterface.
func (s *singletonObject) UuidService() utils.UuidUtilInterface {
	return s.uuidService
}
