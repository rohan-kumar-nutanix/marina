/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 * The helper file for Zeus config
 */

package zeus

import (
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	log "k8s.io/klog/v2"

	zk "github.com/nutanix-core/acs-aos-go/go-zookeeper"
	miscUtil "github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/acs-aos-go/zeus"
	zeusConfig "github.com/nutanix-core/acs-aos-go/zeus/config"
)

const (
	peZeusConfigPath = "/appliance/physical/zeusconfig"
)

var (
	logFatalln = log.Fatalln
)

// Configuration proto cache for PE
type peConfigCache struct {
	clusterExternalIps map[uuid4.Uuid][]string
	clusterNames       map[uuid4.Uuid]string
	sspContainers      map[uuid4.Uuid]*zeusConfig.ConfigurationProto_Container
}

// Configuration proto cache for PC
type configCache struct {
	ConfigCache
	sync.RWMutex
	configProto *zeusConfig.ConfigurationProto

	// Fields of interest to cache.
	clusterUUID  *uuid4.Uuid
	clusterName  string
	localNodeIp  string
	peConfig     peConfigCache
	zkSession    *zeus.ZookeeperSession
	peConfigOnce sync.Once
}

type ZookeeperConnInterface interface {
	Get(string) ([]byte, *zk.Stat, error)
	ChildrenWatchWithBackoff(string) ([]string, *zk.Stat, <-chan zk.Event, error)
}

// InitConfigCache initializes global config cache, returns it and watches for changes
func InitConfigCache(zkSession *zeus.ZookeeperSession) ConfigCache {
	log.Info("Initializing Marina Zeus config cache")
	interval := time.Duration(1000) * time.Millisecond
	maxRetries := 5
	backoff := miscUtil.NewConstantBackoff(interval, maxRetries)

	var configProto *zeusConfig.ConfigurationProto
	for {
		newZeusConfig, err := zeusConfig.NewConfiguration(zkSession, nil)
		if err == nil {
			configProto = newZeusConfig.ConfigProto
			break
		}

		if backoff.Backoff() == miscUtil.Stop {
			logFatalln("Failed to initialize Zeus config: ", err)
		}
	}

	configObj := &configCache{}
	configObj.zkSession = zkSession
	configObj.processConfigProto(configProto)

	// Start zeus config watcher thread.
	go zeus.KeepWatchingContentsOf(zkSession.Conn, zeusConfig.ZeusConfigZKPath, configObj.zeusConfigWatcher)

	log.Info("Initialized Marina Zeus config cache")
	return configObj
}

// Helper method to update cache.
func (config *configCache) processConfigProto(configProto *zeusConfig.ConfigurationProto) {
	log.Info("Updating Marina Zeus config cache")

	config.Lock()
	defer config.Unlock()

	config.configProto = configProto
	config.clusterName = configProto.GetClusterName()
	config.localNodeIp = configProto.MustGetIp()

	var err error
	config.clusterUUID, err = uuid4.StringToUuid4(configProto.GetClusterUuid())
	if err != nil {
		logFatalln("Failed to read cluster UUID: ", err)
	}
}

// Zookeeper callback handler for changes in configuration Zk watcher.
// Also, further triggers updateCb function if passed in Init.
func (config *configCache) zeusConfigWatcher(conn *zk.Conn, data []byte) {
	log.Info("Zeus config watcher callback triggered")
	updatedConfig := &zeusConfig.ConfigurationProto{}
	err := proto.Unmarshal(data, updatedConfig)
	if err != nil {
		logFatalln("Failed to unmarshal Zeus config proto: ", err)
		return
	}
	config.processConfigProto(updatedConfig)
}

func (config *configCache) updatePeConfig(conn ZookeeperConnInterface, clusterUuids []string) {
	log.Info("Updating PE config")
	config.Lock()
	defer config.Unlock()

	clusterExternalIps := make(map[uuid4.Uuid][]string)
	clusterNames := make(map[uuid4.Uuid]string)
	sspContainers := make(map[uuid4.Uuid]*zeusConfig.ConfigurationProto_Container)

	for _, clusterUuidStr := range clusterUuids {
		clusterUuid, err := uuid4.StringToUuid4(clusterUuidStr)
		if err != nil {
			log.Errorf("Invalid cluster uuid in zk node : %s", clusterUuidStr)
			continue
		}

		zkData, _, err := conn.Get(peZeusConfigPath + "/" + clusterUuidStr)
		if err != nil {
			log.Errorf("Failed to get zk node %s for cluster %s : %s", peZeusConfigPath, clusterUuidStr, err)
			continue
		}

		configProto := &zeusConfig.ConfigurationProto{}
		if err = proto.Unmarshal(zkData, configProto); err != nil {
			log.Errorf("Failed to parse zookeeper cluster config from zookeeper for cluster with UUID %s: %s.",
				clusterUuidStr, err)
			continue
		}

		var externalIps []string
		for _, node := range configProto.GetNodeList() {
			if node.GetServiceVmExternalIp() != "" {
				externalIps = append(externalIps, node.GetServiceVmExternalIp())
			}
		}
		clusterExternalIps[*clusterUuid] = externalIps
		clusterNames[*clusterUuid] = configProto.GetClusterName()
		if configProto.DefaultSspContainerId != nil {
			for _, container := range configProto.GetContainerList() {
				if container.ContainerId != nil && *container.ContainerId == *configProto.DefaultSspContainerId {
					sspContainers[*clusterUuid] = container
				}
			}
		}
	}
	config.peConfig.clusterExternalIps = clusterExternalIps
	config.peConfig.clusterNames = clusterNames
	config.peConfig.sspContainers = sspContainers
}

func (config *configCache) initPeConfigWatcher(conn ZookeeperConnInterface, firstPublishWg *sync.WaitGroup) {
	log.Info("Initialising PE config")
	firstPublishDone := false
	for {
		clusterUuids, _, eventChan, err := conn.ChildrenWatchWithBackoff(peZeusConfigPath)
		if err != nil {
			logFatalln("Failed to get children of %s zk node: ", peZeusConfigPath, err)
			return
		}

		config.updatePeConfig(conn, clusterUuids)
		if !firstPublishDone {
			firstPublishWg.Done()
			firstPublishDone = true
		}

		event := <-eventChan
		log.Infof("PE config updated, received zk event: %v.", event)
	}
}

func (config *configCache) maybeInitPeConfigWatcher() {
	config.peConfigOnce.Do(func() {
		var firstPublishWg sync.WaitGroup
		firstPublishWg.Add(1)
		go config.initPeConfigWatcher(config.zkSession.Conn, &firstPublishWg)
		firstPublishWg.Wait()
	})
}

// ClusterName returns the name of the cluster.
func (config *configCache) ClusterName() string {
	config.RLock()
	defer config.RUnlock()
	return config.clusterName
}

// ClusterUuid returns the UUID of the cluster.
func (config *configCache) ClusterUuid() *uuid4.Uuid {
	config.RLock()
	defer config.RUnlock()
	return config.clusterUUID
}

// LocalNodeIp returns the backplane IP address of Local Node.
func (config *configCache) LocalNodeIp() string {
	config.RLock()
	defer config.RUnlock()
	return config.localNodeIp
}

// ContainerName gets container name for a container id
func (config *configCache) ContainerName(containerId int64) (string, error) {
	config.RLock()
	defer config.RUnlock()
	name, _, err := config.configProto.GetContainerInfo(containerId)
	return name, err
}

// PeClusterUuids returns the UUIDs of the PE cluster registered to PC.
func (config *configCache) PeClusterUuids() []*uuid4.Uuid {
	config.maybeInitPeConfigWatcher()
	config.RLock()
	defer config.RUnlock()
	pes := make([]*uuid4.Uuid, 0)
	for peUuid := range config.peConfig.clusterNames {
		pes = append(pes, uuid4.ToUuid4(peUuid.RawBytes()))
	}
	return pes
}

// ClusterExternalIps returns the CVM IPs for the PE with given UUID.
func (config *configCache) ClusterExternalIps(peUuid *uuid4.Uuid) []string {
	config.maybeInitPeConfigWatcher()
	config.RLock()
	defer config.RUnlock()
	if ips, ok := config.peConfig.clusterExternalIps[*peUuid]; ok {
		return ips
	}
	return nil
}

// PeClusterName returns the Cluster Name for the PE with given UUID.
func (config *configCache) PeClusterName(peUuid *uuid4.Uuid) *string {
	config.maybeInitPeConfigWatcher()
	config.RLock()
	defer config.RUnlock()
	if name, ok := config.peConfig.clusterNames[*peUuid]; ok {
		return &name
	}
	return nil
}

// ClusterSSPContainerUuidMap returns the map of Cluster UUIDs and SSP containers.
func (config *configCache) ClusterSSPContainerUuidMap() map[uuid4.Uuid]*zeusConfig.ConfigurationProto_Container {
	config.maybeInitPeConfigWatcher()
	config.RLock()
	defer config.RUnlock()
	return config.peConfig.sspContainers
}
