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

type RemoteEndpointInterface interface {
	IsPE() bool
	GetUserVisibleId(zeusConfig marinaZeus.ConfigCache) *string
}

func IsPE(remoteClusterUuid uuid4.Uuid) bool {
	return remoteClusterUuid != NilUuid
}

func GetUserVisibleId(zeusConfig marinaZeus.ConfigCache, remoteClusterUuid uuid4.Uuid) string {
	if IsPE(remoteClusterUuid) {
		peUuid := remoteClusterUuid
		userVisibleId := zeusConfig.PeClusterName(&peUuid)
		if userVisibleId != nil {
			return *userVisibleId
		}

		ip := zeusConfig.ClusterExternalIps(&peUuid)
		if len(ip) > 0 {
			return ip[0]
		}

		peUuidStr := peUuid.String()
		return peUuidStr
	}
	return ""
}
