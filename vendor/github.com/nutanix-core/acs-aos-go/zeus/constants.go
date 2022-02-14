/*
 * Copyright (c) 2019 Nutanix Inc. All rights reserved.
 *
 * Authors: noufal.muhammed@nutanix.com
 *
 * This file contains zeus global constants.
 */

package zeus

const (
	GO_LEADERSHIP_ROOT     = "/appliance/logical/goleaders"
	PY_LEADERSHIP_ROOT     = "/appliance/logical/pyleaders"
	COMMON_LEADERSHIP_ROOT = "/appliance/logical/leaders"

	HEALTH_MONITOR_ROOT = "/appliance/logical/health-monitor"

	// Common path name prefix for Golang leader node children zknodes.
	// This need to be updated if the go-zookeeper client changes it.
	PROTECTED_PREFIX = "_c_"
)
