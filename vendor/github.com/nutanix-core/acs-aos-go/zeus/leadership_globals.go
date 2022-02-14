/*
* Copyright (c) 2019 Nutanix Inc. All rights reserved.
*
* Authors: noufal.muhammed@nutanix.com
*
* This file contains leader election module related global variables.
*/

package zeus

import (
	"errors"

	"github.com/nutanix-core/acs-aos-go/go-zookeeper"
)

var (
	WORLD_ACL = zk.WorldACL(zk.PermAll)
)

var (
	ErrCouldNotVolunteer   = errors.New("Could not create volunteer node")
	ErrGetContendersFailed = errors.New("Could not get list of contenders")
	ErrSelfNotFound        = errors.New("Could not find self entry in list of contenders")
	ErrZkNotConnected      = errors.New("Zookeeper session not connected")
	ErrNoContenders        = errors.New("No contenders")
	ErrRetryExceeded       = errors.New("Number of retries exceeded the maximum limit")
	ErrNotAGoLeadership    = errors.New("Not a Go leadership")
)
