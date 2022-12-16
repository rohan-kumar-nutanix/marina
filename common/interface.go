/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 * Marina common interface
 *
 */

package common

import (
	"strconv"
	"sync"
	"time"

	log "k8s.io/klog/v2"

	ergonClient "github.com/nutanix-core/acs-aos-go/ergon/client"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc/serial_executor"
	"github.com/nutanix-core/acs-aos-go/zeus"
	"github.com/nutanix-core/content-management-marina/db"
	"github.com/nutanix-core/content-management-marina/grpc/catalog/catalog_item"
	utils "github.com/nutanix-core/content-management-marina/util"
)

type MarinaCommonInterfaces interface {
	CatalogItemService() catalog_item.CatalogItemInterface
	ErgonService() ergonClient.Ergon
	IdfService() db.IdfClientInterface
	SerialExecutor() serial_executor.SerialExecutorIfc
	ZkSession() *zeus.ZookeeperSession
}

type singletonService struct {
	catalogItemService catalog_item.CatalogItemInterface
	ergonService       ergonClient.Ergon
	idfService         db.IdfClientInterface
	serialExecutor     serial_executor.SerialExecutorIfc
	zkSession          *zeus.ZookeeperSession
}

const (
	zkPort      = 9876
	zkTimeOut   = time.Duration(20) * time.Second
	serviceName = "marina"
)

var (
	singleton            MarinaCommonInterfaces
	singletonServiceOnce sync.Once
)

// InitSingletonService - Initialize a singleton Marina service.
func InitSingletonService() {
	singletonServiceOnce.Do(func() {
		singleton = &singletonService{
			catalogItemService: new(catalog_item.CatalogItemImpl),
			ergonService:       ergonClient.NewErgonService(utils.HostAddr, ergonClient.DefaultErgonPort),
			idfService:         db.IdfClientWithRetry(),
			serialExecutor:     serial_executor.NewSerialExecutor(),
			zkSession:          initZkSession(),
		}
	})
}

// Interfaces - Returns the singleton for MarinaCommonInterfaces
func Interfaces() MarinaCommonInterfaces {
	return singleton
}

// CatalogItemService - Returns the singleton for CatalogItemInterface
func (s *singletonService) CatalogItemService() catalog_item.CatalogItemInterface {
	return s.catalogItemService
}

// ErgonService returns the singleton Ergon service.
func (s *singletonService) ErgonService() ergonClient.Ergon {
	return s.ergonService
}

// IdfService - Returns the singleton for IdfClientInterface
func (s *singletonService) IdfService() db.IdfClientInterface {
	return s.idfService
}

// SerialExecutor returns the singleton Serial Executor.
func (s *singletonService) SerialExecutor() serial_executor.SerialExecutorIfc {
	return s.serialExecutor
}

// ZkSession returns the singleton ZK Session.
func (s *singletonService) ZkSession() *zeus.ZookeeperSession {
	return s.zkSession
}

// Initialize Zk session for Marina service.
func initZkSession() *zeus.ZookeeperSession {
	zkServers := []string{utils.HostAddr + ":" + strconv.Itoa(zkPort)}
	zkSession, err := zeus.NewZookeeperSession(zkServers, zkTimeOut)
	if err != nil {
		log.Fatalln("Failed to create a Zookeeper session: ", err)
	}
	_ = zkSession.WaitForConnection()
	log.Infof("Initialized a Zookeeper session, ID = %d.", zkSession.Conn.SessionID())
	zkSession.MayBeCloseOldSession(zkServers, zkTimeOut, serviceName)
	return zkSession
}
