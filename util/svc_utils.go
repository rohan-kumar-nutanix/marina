/*
 * Copyright (c) 2021 Nutanix Inc. All rights reserved.
 *
 * Author: rajesh.battala@nutanix.com
 *
 * This file implements utility functions to do localhost operations when the
 * service is running inside a docker container.
 */

package utils

import (
	"errors"
	"fmt"
	"strconv"

	log "k8s.io/klog/v2"

	ergon "github.com/nutanix-core/acs-aos-go/ergon/client"
	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	zeus "github.com/nutanix-core/acs-aos-go/zeus"
)

// Initialise zookeeper connection.
func GetZkSession() (*zeus.ZookeeperSession, error) {
	zkIP := []string{HostAddr + ":" + strconv.Itoa(ZkHostPort)}
	zkSession, err := zeus.NewZookeeperSession(zkIP, ZkTimeOut)
	if err != nil {
		log.Error(err)
		return nil, errors.New(fmt.Sprintf("Failed to establish " +
			"zookeeper connection"))
	}
	return zkSession, nil
}

// Initialise instance of idf service.
func GetInsightsService() insights_interface.InsightsServiceInterface {
	return insights_interface.NewInsightsServiceInterface(HostAddr,
		uint16(*insights_interface.DefaultInsightPort))
}

// Initialise instance of ergon service.
func GetErgonService() ergon.Ergon {
	return ergon.NewErgonService(HostAddr, ergon.DefaultErgonPort)
}
