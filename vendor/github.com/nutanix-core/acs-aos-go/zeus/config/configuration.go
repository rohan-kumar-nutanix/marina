/*
 * Copyright (c) 2016 Nutanix Inc. All rights reserved.
 *
 * Author: mohammad.ahmad1@nutanix.com
 * Utils for the zeus configuration.
 *
 */

package zeus_config

import (
	"sort"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	pithos_ifc "github.com/nutanix-core/acs-aos-go/pithos"
	"github.com/nutanix-core/acs-aos-go/zeus"
)

const (
	ZeusConfigZKPath = "/appliance/physical/configuration"
	zkTimeoutSecs    = time.Second
)

type Configuration struct {
	ZKSession    *zeus.ZookeeperSession
	hostPortList []string
	ConfigProto  *ConfigurationProto
}

func MustGetNewConfiguration(
	zkSession *zeus.ZookeeperSession, hostPortList []string) *Configuration {

	config, err := NewConfiguration(zkSession, hostPortList)
	if err != nil {
		glog.Fatal(err)
	}
	return config
}

func NewConfiguration(
	zkSession *zeus.ZookeeperSession, hostPortList []string) (
	*Configuration, error) {

	// TODO return error if the zookeeper host port list is not set.
	if hostPortList == nil || len(hostPortList) == 0 {
		hostPortList = GetZkHostPortList()
	}

	var err error
	if zkSession == nil {
		zkSession, err = zeus.NewZookeeperSession(hostPortList, zkTimeoutSecs)
		if err != nil {
			glog.Errorf("Failed to create ZK session on host "+
				"port list: %s with: %s", hostPortList, err)
			return nil, err
		}
	}

	configProto, err := getConfigProto(zkSession)
	if err != nil {
		return nil, err
	}

	return &Configuration{
		ZKSession:    zkSession,
		hostPortList: hostPortList,
		ConfigProto:  configProto,
	}, nil
}

func MockConfiguration() *Configuration {
	return &Configuration{
		ConfigProto: MockConfig1(),
	}
}

func (config *Configuration) AddStorageContainer(
	ctrName string,
	storagePool *ConfigurationProto_StoragePool) (*ConfigurationProto, error) {

	containerUUID, err := uuid4.New()
	if err != nil {
		glog.Errorf("Failed to create new UUID: %s", err)
		return nil, err
	}

	containerID, err := zeus.Allocate_component_ids(config.ZKSession.Conn, 1)
	if err != nil {
		glog.Error("Failed to allocate component ID: ", err)
		return nil, err
	}

	storagePoolIDs := []int64{storagePool.GetStoragePoolId()}
	storagePoolUUIDs := []string{storagePool.GetStoragePoolUuid()}
	faultToleranceState := config.ConfigProto.ClusterFaultToleranceState
	rf := *faultToleranceState.DesiredMaxFaultTolerance + 1

	oplogParams := &pithos_ifc.VDiskConfig_Params_OplogParams{
		NumStripes:        proto.Int32(1),
		NeedSync:          proto.Bool(false),
		ReplicationFactor: proto.Int32(rf),
	}

	storageTierList := config.ConfigProto.StorageTierList
	sort.Sort(sort.Reverse(ByRandomPriority(storageTierList)))

	ilmDownMigrateTime := []int32{}
	randomIOTierPreference := []string{}
	for _, v := range storageTierList {
		if !v.GetDeleted() {
			randomIOTierPreference = append(randomIOTierPreference,
				v.GetStorageTierName())
			ilmDownMigrateTime = append(ilmDownMigrateTime, 1800)
		}
	}

	storageTierList = config.ConfigProto.StorageTierList
	sort.Sort(sort.Reverse(BySequentialPriority(storageTierList)))

	sequentialIOTierPreference := []string{}
	for _, v := range storageTierList {
		if !v.GetDeleted() {
			sequentialIOTierPreference = append(sequentialIOTierPreference,
				v.GetStorageTierName())
		}
	}

	params := &ConfigurationProto_ContainerParams{
		ReplicationFactor:          proto.Int32(rf),
		OplogParams:                oplogParams,
		FingerprintOnWrite:         proto.Bool(false),
		IlmDownMigrateTimeSecs:     ilmDownMigrateTime,
		RandomIoTierPreference:     randomIOTierPreference,
		SequentialIoTierPreference: sequentialIOTierPreference,
	}

	storageContainer := &ConfigurationProto_Container{
		ContainerName:   proto.String(ctrName),
		ContainerId:     proto.Int64(int64(containerID)),
		ContainerUuid:   proto.String(containerUUID.UuidToString()),
		StoragePoolId:   storagePoolIDs,
		StoragePoolUuid: storagePoolUUIDs,
		Params:          params,
	}

	configProto := config.ConfigProto
	configProto.ContainerList = append(configProto.ContainerList, storageContainer)

	vStore := &ConfigurationProto_VStore{
		VstoreId:      storageContainer.ContainerId,
		VstoreName:    storageContainer.ContainerName,
		ContainerId:   storageContainer.ContainerId,
		ContainerUuid: storageContainer.ContainerUuid,
		VstoreUuid:    storageContainer.ContainerUuid,
	}

	configProto.VstoreList = append(configProto.VstoreList, vStore)
	return configProto, nil
}

//-----------------------------------------------------------------------------

func getConfigProto(zkSession *zeus.ZookeeperSession) (*ConfigurationProto, error) {
	glog.Info("Reading the configuration proto")
	var data []byte
	var err error

	if zkSession.WaitForConnection() {
		data, _, err = zkSession.Conn.Get(ZeusConfigZKPath)
		if err != nil {
			glog.Errorf("Failed to get ZK node: %s with error: %s ",
				ZeusConfigZKPath, err)
			return nil, err
		}
	}

	configProto := &ConfigurationProto{}
	err = proto.Unmarshal(data, configProto)
	if err != nil {
		glog.Error("Failed to unmarshal ZK config proto with error: ", err)
		return nil, err
	}

	// TODO Get chunkified config if `SymbolicLink` field set.
	return configProto, nil
}

//-----------------------------------------------------------------------------

func (config *Configuration) Commit(configProto *ConfigurationProto) (
	*ConfigurationProto, bool, error) {

	oldLogicalTimestamp := configProto.GetLogicalTimestamp()
	*configProto.LogicalTimestamp = config.ConfigProto.GetLogicalTimestamp() + 1

	glog.Infof("Updating the zeus config timestamp from: %d to: %d",
		oldLogicalTimestamp, *configProto.LogicalTimestamp)

	var data []byte
	data, err := proto.Marshal(configProto)
	if err != nil {
		glog.Error("Failed to marshal config proto: ", err)
		return nil, false, err
	}

	ok, err := config.commit(data, int32(oldLogicalTimestamp))
	if err != nil {
		return nil, false, err
	}

	if ok {
		config.ConfigProto = configProto
		return nil, true, nil
	}

	newConfigProto, err := getConfigProto(config.ZKSession)
	if err != nil {
		return nil, false, err
	}

	glog.Warning("Commit failed due to a stale zeus config with timestamp: %d"+
		" new timestamp: %d", oldLogicalTimestamp, newConfigProto.LogicalTimestamp)

	return newConfigProto, false, nil
}

func (config *Configuration) WatchConfig() chan *ConfigurationProto {
	if config.ZKSession == nil {
		return nil
	}

	watch := make(chan *ConfigurationProto)
	go func() {
		for {
			data, _, ch, err := config.ZKSession.GetW(ZeusConfigZKPath, true /*retry*/)
			if err != nil {
				panic(err)
			}
			updatedConfig := &ConfigurationProto{}
			err = proto.Unmarshal(data, updatedConfig)
			if err != nil {
				panic(err)
			}
			watch <- updatedConfig
			<-ch
		}
	}()
	return watch
}

func (config *Configuration) commit(data []byte, version int32) (bool, error) {

	if config.ZKSession.WaitForConnection() {
		_, err := config.ZKSession.Conn.Set(ZeusConfigZKPath, data, int32(version))
		if err != nil {
			glog.Error("Failed to commit zeus configuration: ", err)
			return false, err
		}
	}

	return true, nil
}

func MarshalDataToConfigProto(data []byte) (*ConfigurationProto, error) {
	configProto := &ConfigurationProto{}
	err := proto.Unmarshal(data, configProto)
	if err != nil {
		glog.Error("Failed to unmarshal ZK config proto with error: ", err)
		return nil, err
	}
	return configProto, nil
}
