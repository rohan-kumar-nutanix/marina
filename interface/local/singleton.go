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
	utils "github.com/nutanix-core/content-management-marina/util"
)

type MarinaInternalInterfaces interface {
	AuthzIfc() authz.AuthzInterface
	MetadataIfc() metadata.EntityMetadataInterface
	ProtoIfc() utils.ProtoUtilInterface
	UuidIfc() utils.UuidUtilInterface
}

type singletonObject struct {
	authzIfc    authz.AuthzInterface
	metadataIfc metadata.EntityMetadataInterface
	protoIfc    utils.ProtoUtilInterface
	uuidIfc     utils.UuidUtilInterface
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
			protoIfc:    new(utils.ProtoUtil),
			uuidIfc:     new(utils.UuidUtil),
		}
	})
}

// GetSingletonServiceWithParams - Initialize a singleton Marina service with params. Should only be used in UTs
func GetSingletonServiceWithParams(authzIfc authz.AuthzInterface, metadataIfc metadata.EntityMetadataInterface,
	protoIfc utils.ProtoUtilInterface, uuidIfc utils.UuidUtilInterface) *singletonObject {

	return &singletonObject{
		authzIfc:    authzIfc,
		metadataIfc: metadataIfc,
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

// ProtoIfc - Returns the singleton for ProtoUtilInterface
func (s *singletonObject) ProtoIfc() utils.ProtoUtilInterface {
	return s.protoIfc
}

// UuidIfc returns the singleton for UuidUtilInterface
func (s *singletonObject) UuidIfc() utils.UuidUtilInterface {
	return s.uuidIfc
}
