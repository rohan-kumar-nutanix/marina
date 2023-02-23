/*
 * Copyright (c) 2022 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 * Marina external interfaces singleton
 *
 */

package external

import (
	"strconv"
	"sync"
	"time"

	log "k8s.io/klog/v2"

	ergonClient "github.com/nutanix-core/acs-aos-go/ergon/client"
	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	cpdb "github.com/nutanix-core/acs-aos-go/nusights/util/db"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/authz"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/authz/authz_cache"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc/serial_executor"
	"github.com/nutanix-core/acs-aos-go/zeus"
	"github.com/nutanix-core/content-management-marina/common"
	"github.com/nutanix-core/content-management-marina/db"
	utils "github.com/nutanix-core/content-management-marina/util"
	marinaZeus "github.com/nutanix-core/content-management-marina/zeus"
)

type MarinaExternalInterfaces interface {
	CPDBIfc() cpdb.CPDBClientInterface
	ErgonIfc() ergonClient.Ergon
	IamIfc() authz_cache.IamClientIfc
	IdfIfc() db.IdfClientInterface
	SerialExecutor() serial_executor.SerialExecutorIfc
	ZeusConfig() marinaZeus.ConfigCache
	ZkSession() *zeus.ZookeeperSession
}

type singletonObject struct {
	cpdbIfc        cpdb.CPDBClientInterface
	ergonIfc       ergonClient.Ergon
	iamIfc         authz_cache.IamClientIfc
	idfIfc         db.IdfClientInterface
	serialExecutor serial_executor.SerialExecutorIfc
	zeusConfig     marinaZeus.ConfigCache
	zkSession      *zeus.ZookeeperSession
}

const (
	zkPort      = 9876
	zkTimeOut   = time.Duration(20) * time.Second
	serviceName = "marina"
)

var (
	singleton            MarinaExternalInterfaces
	singletonServiceOnce sync.Once
)

// InitSingletonService - Initialize a singleton Marina service
func InitSingletonService() {
	singletonServiceOnce.Do(func() {
		zkSession := initZkSession()
		singleton = &singletonObject{
			cpdbIfc:        cpdb.NewCPDBService(utils.HostAddr, uint16(*insights_interface.InsightsPort)),
			ergonIfc:       ergonClient.NewErgonService(utils.HostAddr, ergonClient.DefaultErgonPort),
			iamIfc:         newIamClient(),
			idfIfc:         db.IdfClientWithRetry(),
			serialExecutor: serial_executor.NewSerialExecutor(),
			zeusConfig:     marinaZeus.InitConfigCache(zkSession),
			zkSession:      zkSession,
		}
	})
}

// GetSingletonServiceWithParams - Initialize a singleton Marina service with params. Should only be used in UTs
func GetSingletonServiceWithParams(cpdbService cpdb.CPDBClientInterface, ergonService ergonClient.Ergon,
	idfService db.IdfClientInterface, serialExecutor serial_executor.SerialExecutorIfc,
	zeusConfig marinaZeus.ConfigCache, zkSession *zeus.ZookeeperSession) *singletonObject {

	return &singletonObject{
		cpdbIfc:        cpdbService,
		ergonIfc:       ergonService,
		idfIfc:         idfService,
		serialExecutor: serialExecutor,
		zeusConfig:     zeusConfig,
		zkSession:      zkSession,
	}
}

// Interfaces - Returns the singleton for MarinaExternalInterfaces
func Interfaces() MarinaExternalInterfaces {
	return singleton
}

// CPDBIfc - Returns the singleton for CPDBClientInterface
func (s *singletonObject) CPDBIfc() cpdb.CPDBClientInterface {
	return s.cpdbIfc
}

// ErgonIfc returns the singleton for Ergon Interface
func (s *singletonObject) ErgonIfc() ergonClient.Ergon {
	return s.ergonIfc
}

// IamIfc returns the singleton for IamClientIfc
func (s *singletonObject) IamIfc() authz_cache.IamClientIfc {
	return s.iamIfc
}

// IdfIfc - Returns the singleton for IdfClientInterface
func (s *singletonObject) IdfIfc() db.IdfClientInterface {
	return s.idfIfc
}

// SerialExecutor returns the singleton for Serial Executor
func (s *singletonObject) SerialExecutor() serial_executor.SerialExecutorIfc {
	return s.serialExecutor
}

// ZeusConfig returns the singleton for Zeus config
func (s *singletonObject) ZeusConfig() marinaZeus.ConfigCache {
	return s.zeusConfig
}

// ZkSession returns the singleton for ZK Session
func (s *singletonObject) ZkSession() *zeus.ZookeeperSession {
	return s.zkSession
}

// newIamClient - Initializes the IAM Interface
func newIamClient() authz_cache.IamClientIfc {
	authOptions := &authz.Options{Retry: false}
	iamClient, err := authz.NewIamClient(
		common.MarinaServiceCaChainPath,
		common.MarinaServiceCertPath,
		common.MarinaServiceKeyPath,
		[]*string{common.MarinaServiceIcaPath},
		authOptions)
	if err != nil {
		log.Fatalln("Failed to initialize Iam Client: ", err)
	}
	return iamClient
}

// Initialize Zk session for Marina service
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
