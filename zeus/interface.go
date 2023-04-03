/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 * The interfaces for Zeus config
 */

package zeus

import (
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	zeusConfig "github.com/nutanix-core/acs-aos-go/zeus/config"
)

type ConfigCache interface {
	ClusterName() string
	ClusterUuid() *uuid4.Uuid
	LocalNodeIp() string
	ContainerName(containerId int64) (string, error)
	PeClusterUuids() []*uuid4.Uuid
	ClusterExternalIps(peUuid *uuid4.Uuid) []string
	PeClusterName(peUuid *uuid4.Uuid) *string
	ClusterSSPContainerUuidMap() map[uuid4.Uuid]*zeusConfig.ConfigurationProto_Container
	CatalogPeRegistered(peUuid *uuid4.Uuid) *bool
	IsRateLimitSupported(peUuid *uuid4.Uuid) *bool
}
