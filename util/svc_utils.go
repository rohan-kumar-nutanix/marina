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
	"github.com/golang/protobuf/proto"
	ergon "github.com/nutanix-core/acs-aos-go/ergon/client"
	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	serialExecutor "github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc/serial_executor"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/acs-aos-go/zeus"
	log "k8s.io/klog/v2"
	"strconv"
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

// GetInsightsService Initialise instance of IDF service.
func GetInsightsService() insights_interface.InsightsServiceInterface {
	return insights_interface.NewInsightsServiceInterface(HostAddr,
		uint16(*insights_interface.DefaultInsightPort))
}

// GetErgonService Initialise instance of Ergon service.
func GetErgonService() ergon.Ergon {
	return ergon.NewErgonService(HostAddr, ergon.DefaultErgonPort)
}

type MarshalInterface interface {
	Marshal(m proto.Message) ([]byte, error)
}

type marshalUtil struct {
}

func Marshal(m proto.Message) ([]byte, error) {
	return proto.Marshal(m)
}

type UuidInterface interface {
	New() (*uuid4.Uuid, error)
}

type UuidUtil struct {
}

func NewUuid() (*uuid4.Uuid, error) {
	uuid, err := uuid4.New()
	return uuid, err
}

// SerialExecutor returns the singleton serial executor.
func SerialExecutor() serialExecutor.SerialExecutorIfc {
	// TODO: init with supported Max Tasks
	return serialExecutor.NewSerialExecutor()
}
