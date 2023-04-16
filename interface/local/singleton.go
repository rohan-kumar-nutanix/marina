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
	"github.com/nutanix-core/content-management-marina/metadata"
	"github.com/nutanix-core/content-management-marina/odata"
	utils "github.com/nutanix-core/content-management-marina/util"
)

type MarinaInternalInterfaces interface {
	AuthzIfc() authz.AuthzInterface
	MetadataIfc() metadata.EntityMetadataInterface
	OdataIfc() odata.OdataInterface
	ProtoIfc() utils.ProtoUtilInterface
	UuidIfc() utils.UuidUtilInterface
	ErrorIfc() utils.GrpcStatusUtilInterface
}

type singletonObject struct {
	authzIfc    authz.AuthzInterface
	metadataIfc metadata.EntityMetadataInterface
	odataIfc    odata.OdataInterface
	protoIfc    utils.ProtoUtilInterface
	uuidIfc     utils.UuidUtilInterface
	errorIfc    utils.GrpcStatusUtilInterface
}

var (
	singleton            MarinaInternalInterfaces
	singletonServiceOnce sync.Once
)

// InitSingletonService - Initialize a singleton Marina service.
func InitSingletonService() {
	singletonServiceOnce.Do(func() {
		singleton = &singletonObject{
			authzIfc:    new(authz.AuthzUtil),
			metadataIfc: new(metadata.EntityMetadataUtil),
			odataIfc:    new(odata.OdataUtil),
			protoIfc:    new(utils.ProtoUtil),
			uuidIfc:     new(utils.UuidUtil),
			errorIfc:    new(utils.GrpcStatusUtil),
		}
	})
}

// GetSingletonServiceWithParams - Initialize a singleton Marina service with params. Should only be used in UTs
func GetSingletonServiceWithParams(authzIfc authz.AuthzInterface, metadataIfc metadata.EntityMetadataInterface,
	odataIfc odata.OdataInterface, protoIfc utils.ProtoUtilInterface, uuidIfc utils.UuidUtilInterface) *singletonObject {
	return &singletonObject{
		authzIfc:    authzIfc,
		metadataIfc: metadataIfc,
		odataIfc:    odataIfc,
		protoIfc:    protoIfc,
		uuidIfc:     uuidIfc,
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

// MetadataIfc - Returns the singleton for EntityMetadataInterface
func (s *singletonObject) MetadataIfc() metadata.EntityMetadataInterface {
	return s.metadataIfc
}

// OdataIfc - Returns the singleton for OdataInterface
func (s *singletonObject) OdataIfc() odata.OdataInterface {
	return s.odataIfc
}

// ProtoIfc - Returns the singleton for ProtoUtilInterface
func (s *singletonObject) ProtoIfc() utils.ProtoUtilInterface {
	return s.protoIfc
}

// UuidIfc returns the singleton for UuidUtilInterface
func (s *singletonObject) UuidIfc() utils.UuidUtilInterface {
	return s.uuidIfc
}

// ErrorIfc returns the singleton for GrpcStatusUtilInterface
func (s *singletonObject) ErrorIfc() utils.GrpcStatusUtilInterface {
	return s.errorIfc
}
