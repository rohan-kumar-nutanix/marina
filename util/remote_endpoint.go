/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Author: rishabh.gupta@nutanix.com
*
 */

package utils

import (
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	marinaZeus "github.com/nutanix-core/content-management-marina/zeus"
)

type RemoteEndpoint struct {
	SourceClusterUuid uuid4.Uuid
	RemoteClusterUuid uuid4.Uuid
}

type RemoteEndpointInterface interface {
	IsPE() bool
	GetUserVisibleId(zeusConfig marinaZeus.ConfigCache) *string
}

func (remoteEndpoint *RemoteEndpoint) IsPE() bool {
	return remoteEndpoint.RemoteClusterUuid != NilUuid
}

func (remoteEndpoint *RemoteEndpoint) GetUserVisibleId(zeusConfig marinaZeus.ConfigCache) *string {
	if remoteEndpoint.IsPE() {
		peUuid := remoteEndpoint.RemoteClusterUuid
		userVisibleId := zeusConfig.PeClusterName(&peUuid)
		if userVisibleId != nil {
			return userVisibleId
		}

		ip := zeusConfig.ClusterExternalIps(&peUuid)
		if len(ip) > 0 {
			return &ip[0]
		}

		peUuidStr := peUuid.String()
		return &peUuidStr
	}
	return nil
}
