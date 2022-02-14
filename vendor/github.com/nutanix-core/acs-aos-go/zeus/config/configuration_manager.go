/*
 * Copyright (c) 2019 Nutanix Inc. All rights reserved.
 *
 * Author: prateek.kajaria@nutanix.com
 *
 * This module provides wrapper functions to manage the cluster configuration
 * stored in Zookeeper.
 *
 */

package zeus_config

import (
	"flag"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/nutanix-core/acs-aos-go/go-zookeeper"
	"sync"
	"github.com/nutanix-core/acs-aos-go/zeus"
)

var (
	zkConfigPath = flag.String("zk_config_path", "/appliance/physical/configuration",
		"The path at which the configuration for the cluster is stored.")
)

// TODO(Prateek): Handle SymbolicLink field being set in ConfigProto.
// Reading of chunkified proto can be done by invoking zeus_config_printer
// binary on the cluster, much like how the zookeeper python client does.
// For writing we need to identify the chunk where the write is supposed to
// land. This needs to be handled in the python client as well.

type ConfigurationManager struct {
	zkSession          *zeus.ZookeeperSession // Zookeeper session
	configWatchStarted bool                   // Whether the config watch has started or not
	stopWatchCh        chan bool              // Notify the watch goroutine to stop executing
	watchStoppedCh     chan bool              // Notification informing that the watch goroutine has stopped
	rwMutex            sync.RWMutex           // Mutex to prevent two start/stop watch calls
}

// Returns a new configuration manager object. The configuration manager is
// responsible for managing the ConfigProto. It can read, update and also
// watch for changes on configuration. The manager can install only one watch.
// The library does not handle concurrent reads/writes. Any operation can be
// carried out first. The clients are supposed to handle it, or ensure that
// there are no concurrent requests. It is expected that the clients will
// take care of closing the zookeeper session.
func NewConfigurationManager(session *zeus.ZookeeperSession) *ConfigurationManager {

	if session == nil {
		glog.Error("Zookeeper session is nil")
		return nil
	}

	return &ConfigurationManager{
		zkSession:          session,
		configWatchStarted: false,
	}
}

// Get the current configuration from Zookeeper. If 'sync' is set to true, we'll
// request zookeeper server to sync the configuration with the zookeeper master
// so that the read sees the latest configuration, otherwise the local config from
// the zookeeper client is returned.
func (configMgr *ConfigurationManager) GetConfig(sync bool) (*ConfigurationProto, error) {

	glog.Info("Fetching configuration from Zookeeper")

	if sync {
		_, err := configMgr.zkSession.Conn.Sync(*zkConfigPath)
		if err != nil {
			glog.Fatal("Failed to sync zk node with error", err)
		}
	}

	data, _, err := configMgr.zkSession.Conn.Get(*zkConfigPath)
	if err != nil {
		glog.Fatal("Failed to get zk node with error", err)
	}

	configProto := &ConfigurationProto{}
	err = proto.Unmarshal(data, configProto)
	if err != nil {
		glog.Fatal("Failed to unmarshal zk config proto with error", err)
	}

	return configProto, nil
}

// Update Zookeeper with the configuration given in 'configProto'. Returns
// whether the update was successful or not. If there are multiple update
// calls only one succeed. The rest will result in a CAS error. In case of
// a CAS error, the client is expected to call GetConfig and try again with
// the updated logical timestamp.
func (configMgr *ConfigurationManager) UpdateConfig(
	configProto *ConfigurationProto) (bool, error) {

	// Zookeeper allows upto 1MB of data per node. If the size of the serialized
	// proto exceeds this fail for now. This will need to handled by using the
	// SymbolicLink field in configuration.proto
	if proto.Size(configProto) > (1000*1000 - 1) {
		glog.Fatal("Serialized proto size greater than zk node allowed size")
	}

	data, err := proto.Marshal(configProto)
	if err != nil {
		glog.Fatal("Failed to marshal zk config proto with error", err)
	}

	// version is the expected version of zknode. The Set call succeeds if the
	// current version is as expected. It results in a CAS failure otherwise.
	// The verion is automatically bumped up if the write request is accepted.
	version := configProto.GetLogicalTimestamp() - 1
	glog.V(2).Info("Updating with version ", version)
	_, err = configMgr.zkSession.Conn.Set(*zkConfigPath, data, int32(version))

	if err != nil {
		if err == zk.ErrBadVersion {
			glog.Error("Failed to update zk node with error ", err)
			return false, err
		}

		// TODO: Handle session expiry/disconnects gracefully when zookeeper
		// session starts differentiating between the events.
		glog.Fatal("Failed to update zk node with error ", err)
	}

	return true, nil
}

// Register the given callback to be invoked when the configuration changes
// in Zookeeper. The callback is expected to be lightweight since it can
// delay stopping the config watch. Note that this callback may be invoked
// even when the config is updated by calling UpdateConfiguration. The
// callback is invoked immediately after the watch is set to ensure that
// the caller has the latest configuration. At times the callback may be
// invoked redundantly even when the configuration may not have changed.
// Since the manager can install only one watch, attempts at registering
// multiple watches will lead to a FATAL.
func (configMgr *ConfigurationManager) StartConfigWatch(
	config_change_cb func(*ConfigurationProto)) {

	configMgr.rwMutex.Lock()

	// Return if watch has already started.
	if configMgr.configWatchStarted {
		configMgr.rwMutex.Unlock()
		glog.Fatal("Request to start watch when it's already started")
		return
	}

	configMgr.configWatchStarted = true

	configMgr.stopWatchCh = make(chan bool)
	configMgr.watchStoppedCh = make(chan bool)

	configMgr.rwMutex.Unlock()

	go configMgr.configWatch(config_change_cb)
}

// Remove the configuration watch set by an earlier call to
// StartConfigurationWatch. Configuration change callbacks can be invoked
// till the function returns. Multiple attempts at stopping the watch will
// be ignored.
func (configMgr *ConfigurationManager) StopConfigWatch() {
	configMgr.rwMutex.Lock()

	// Return if watch has already been stopped.
	if !configMgr.configWatchStarted {
		configMgr.rwMutex.Unlock()
		glog.Error("Request to stop watch when it's already stopped")
		return
	}

	configMgr.configWatchStarted = false

	// Notify the watch goroutine to stop executing.
	configMgr.stopWatchCh <- true

	// Wait till watch goroutine has ended.
	<-configMgr.watchStoppedCh

	close(configMgr.stopWatchCh)
	close(configMgr.watchStoppedCh)

	configMgr.rwMutex.Unlock()

	glog.Info("Watch stopped")
}

// Goroutine which notifies the consumer class of changes in configuration.
func (configMgr *ConfigurationManager) configWatch(
	config_change_cb func(*ConfigurationProto)) {

	for {
		data, _, zkCh, err := configMgr.zkSession.Conn.GetW(*zkConfigPath)
		if err != nil {
			glog.Fatal("Error while reading zk node", err)
		}

		newConfigProto := &ConfigurationProto{}
		err = proto.Unmarshal(data, newConfigProto)
		if err != nil {
			glog.Fatal("Could not umarshal config proto", err)
		}

		// The callback can potentially block StopConfigWatch. Please ensure that
		// the callback finishes as soon as possible. Launching a goroutine here
		// will ruin inorder delivery guarantees. (The callback can just add the
		// work to be done to a channel, and an executor goroutine can pick up the
		// work item one by one)
		config_change_cb(newConfigProto)

		select {
		// Zk node changed, read the new value and invoke the callback.
		case _, ok := <-zkCh:
			if !ok {
				glog.Fatal("Error reading zookeeper channel")
			}

		// Received notification to stop watch.
		case _, ok := <-configMgr.stopWatchCh:
			if !ok {
				glog.Fatal("Error reading stop watch channel")
			}
			glog.Info("Stopped config watch")
			// Notify of goroutine finish.
			configMgr.watchStoppedCh <- true
			return
		}
	}
}
