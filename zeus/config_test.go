/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Author: rishabh.gupta@nutanix.com
 */

package zeus

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	log "k8s.io/klog/v2"

	zk "github.com/nutanix-core/acs-aos-go/go-zookeeper"
	utilConfig "github.com/nutanix-core/acs-aos-go/nutanix/util-go/config"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	zeusConfig "github.com/nutanix-core/acs-aos-go/zeus/config"
	mockZeus "github.com/nutanix-core/content-management-marina/mocks/zeus"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
)

var (
	// Node UUID matching with one on factory config path.
	testNodeUuid       = "f81bfc7b-377c-440c-927d-d20cdf31e50a"
	testNodeIp         = "127.0.0.2"
	testClusterName    = "TestCluster"
	testClusterUuid, _ = uuid4.New()
)

func TestProcessConfigProto(t *testing.T) {
	// We mock factory_config.json because of MustGetIp call. The node UUID set
	// in this function is f81bfc7b-377c-440c-927d-d20cdf31e50a. Changing it on
	// util-master branch would require changing the nodeUuidTestString
	// variable. Today, there's no utility function to pass the desired config
	// to TestSetupMockFactoryConfig, so we will stick to this.
	utilConfig.TestSetupMockFactoryConfig()
	configProto := &zeusConfig.ConfigurationProto{
		ClusterName: proto.String(testClusterName),
		ClusterUuid: proto.String(testClusterUuid.String()),
		NodeList: []*zeusConfig.ConfigurationProto_Node{
			{
				Uuid:                    &testNodeUuid,
				ControllerVmBackplaneIp: &testNodeIp,
			},
		},
	}

	config := &configCache{}
	config.processConfigProto(configProto)

	assert.Equal(t, testClusterName, config.ClusterName())
	assert.Equal(t, testClusterUuid, config.ClusterUuid())
	assert.Equal(t, testNodeIp, config.LocalNodeIp())
}

func TestProcessConfigProtoUuidError(t *testing.T) {
	utilConfig.TestSetupMockFactoryConfig()
	configProto := &zeusConfig.ConfigurationProto{
		ClusterName: proto.String(testClusterName),
		ClusterUuid: proto.String("foobar"),
		NodeList: []*zeusConfig.ConfigurationProto_Node{
			{
				Uuid:                    &testNodeUuid,
				ControllerVmBackplaneIp: &testNodeIp,
			},
		},
	}

	// Mocking log.Fatalln call
	var fatalErrors []string
	logFatalln = func(args ...interface{}) {
		fatalErrors = append(fatalErrors, fmt.Sprint(args...))
	}

	config := &configCache{}
	config.processConfigProto(configProto)

	// Restoring log.Fatalln call
	logFatalln = log.Fatalln

	assert.Equal(t, 1, len(fatalErrors))
}

func TestZeusConfigWatcher(t *testing.T) {
	utilConfig.TestSetupMockFactoryConfig()
	var zkConn *zk.Conn
	data := &zeusConfig.ConfigurationProto{
		LogicalTimestamp: proto.Int64(2),
		ClusterName:      proto.String(testClusterName),
		ClusterUuid:      proto.String(testClusterUuid.String()),
		NodeList: []*zeusConfig.ConfigurationProto_Node{
			{
				// ServiceVmId is necessary to UnMarshal
				ServiceVmId:             proto.Int64(1),
				Uuid:                    &testNodeUuid,
				ControllerVmBackplaneIp: &testNodeIp,
			},
		},
	}
	protoBytes, _ := proto.Marshal(data)

	config := &configCache{}
	config.zeusConfigWatcher(zkConn, protoBytes)

	assert.Equal(t, testClusterName, config.ClusterName())
	assert.Equal(t, testClusterUuid, config.ClusterUuid())
	assert.Equal(t, testNodeIp, config.LocalNodeIp())
}

func TestZeusConfigWatcherError(t *testing.T) {
	var zkConn *zk.Conn
	message := &marinaIfc.CatalogItemInfo{}
	protoBytes, _ := proto.Marshal(message)

	var fatalErrors []string
	logFatalln = func(args ...interface{}) {
		fatalErrors = append(fatalErrors, fmt.Sprint(args...))
	}

	config := &configCache{}
	config.zeusConfigWatcher(zkConn, protoBytes)

	logFatalln = log.Fatalln

	assert.Equal(t, 1, len(fatalErrors))
}

func TestUpdatePeConfigUuidError(t *testing.T) {
	var zkConn *zk.Conn
	config := &configCache{}
	config.updatePeConfig(zkConn, []string{"Invalid Uuid"})

	assert.Equal(t, config.configsByPe, map[uuid4.Uuid]peConfigCache{})
}

func TestUpdatePeConfigZkError(t *testing.T) {
	mockConn := &mockZeus.ZookeeperConnInterface{}
	mockConn.On("Get", peZeusConfigPath+"/"+testClusterUuid.String()).
		Return(nil, nil, errors.New("oh no")).Once()

	config := &configCache{}
	config.updatePeConfig(mockConn, []string{testClusterUuid.String()})

	assert.Equal(t, config.configsByPe, map[uuid4.Uuid]peConfigCache{})
	mockConn.AssertExpectations(t)
}

func TestUpdatePeConfigProtoError(t *testing.T) {
	mockConn := &mockZeus.ZookeeperConnInterface{}
	mockConn.On("Get", peZeusConfigPath+"/"+testClusterUuid.String()).
		Return([]byte{255}, nil, nil).Once()

	config := &configCache{}
	config.updatePeConfig(mockConn, []string{testClusterUuid.String()})

	assert.Equal(t, config.configsByPe, map[uuid4.Uuid]peConfigCache{})
	mockConn.AssertExpectations(t)
}

func TestUpdatePeConfigNilZkDataError(t *testing.T) {
	mockConn := &mockZeus.ZookeeperConnInterface{}
	mockConn.On("Get", peZeusConfigPath+"/"+testClusterUuid.String()).
		Return(nil, nil, nil).Once()

	config := &configCache{}
	config.updatePeConfig(mockConn, []string{testClusterUuid.String()})

	assert.Equal(t, config.configsByPe, map[uuid4.Uuid]peConfigCache{})
	mockConn.AssertExpectations(t)
}

func TestUpdatePeConfig(t *testing.T) {
	mockConn := &mockZeus.ZookeeperConnInterface{}
	configProto := &zeusConfig.ConfigurationProto{
		ClusterName: proto.String(testClusterName),
		NodeList: zeusConfig.NodeList{
			{ServiceVmId: proto.Int64(1)},
			{ServiceVmExternalIp: proto.String(testNodeIp), ServiceVmId: proto.Int64(1)},
		},
		LogicalTimestamp: proto.Int64(1),
	}
	zkData, _ := proto.Marshal(configProto)
	mockConn.On("Get", peZeusConfigPath+"/"+testClusterUuid.String()).
		Return(zkData, nil, nil).Once()

	config := &configCache{}
	config.updatePeConfig(mockConn, []string{testClusterUuid.String()})

	assert.Equal(t, testClusterName, config.configsByPe[*testClusterUuid].clusterName)
	assert.Equal(t, []string{testNodeIp}, config.configsByPe[*testClusterUuid].clusterExternalIps)
	assert.Equal(t, map[uuid4.Uuid]*zeusConfig.ConfigurationProto_Container{}, config.sspContainers)
	mockConn.AssertExpectations(t)
}

func TestUpdatePeConfigWithContainers(t *testing.T) {
	mockConn := &mockZeus.ZookeeperConnInterface{}
	testContainer := &zeusConfig.ConfigurationProto_Container{
		ContainerId:   proto.Int64(420),
		ContainerName: proto.String("TestContainer"),
	}
	configProto := &zeusConfig.ConfigurationProto{
		ClusterName: proto.String(testClusterName),
		NodeList: zeusConfig.NodeList{
			{ServiceVmId: proto.Int64(1)},
			{ServiceVmExternalIp: proto.String(testNodeIp), ServiceVmId: proto.Int64(1)},
		},
		DefaultSspContainerId: proto.Int64(420),
		ContainerList: []*zeusConfig.ConfigurationProto_Container{
			{ContainerId: proto.Int64(421), ContainerName: proto.String("TestContainer")},
			testContainer,
		},
		LogicalTimestamp: proto.Int64(1),
	}
	zkData, _ := proto.Marshal(configProto)
	mockConn.On("Get", peZeusConfigPath+"/"+testClusterUuid.String()).
		Return(zkData, nil, nil).Once()

	config := &configCache{}
	config.updatePeConfig(mockConn, []string{testClusterUuid.String()})

	assert.Equal(t, testClusterName, config.configsByPe[*testClusterUuid].clusterName)
	assert.Equal(t, []string{testNodeIp}, config.configsByPe[*testClusterUuid].clusterExternalIps)
	assert.Equal(t, 1, len(config.sspContainers))
	mockConn.AssertExpectations(t)
}

func TestInitPeConfigWatcher(t *testing.T) {
	mockConn := &mockZeus.ZookeeperConnInterface{}
	configProto := &zeusConfig.ConfigurationProto{
		ClusterName: proto.String(testClusterName),
		NodeList: zeusConfig.NodeList{
			{ServiceVmId: proto.Int64(1)},
			{ServiceVmExternalIp: proto.String(testNodeIp), ServiceVmId: proto.Int64(1)},
		},
		LogicalTimestamp: proto.Int64(1),
	}
	zkData, _ := proto.Marshal(configProto)
	mockConn.On("Get", peZeusConfigPath+"/"+testClusterUuid.String()).
		Return(zkData, nil, nil).Once()

	getChannel := func() <-chan zk.Event {
		eventChan := make(chan zk.Event)
		go func() { eventChan <- zk.Event{} }()
		return eventChan
	}
	mockConn.On("ChildrenWatchWithBackoff", mock.Anything).
		Return([]string{testClusterUuid.String()}, nil, getChannel(), nil).Once()
	mockConn.On("ChildrenWatchWithBackoff", mock.Anything).
		Return(nil, nil, nil, errors.New("oh no")).Once()

	config := &configCache{}
	wg := sync.WaitGroup{}
	wg.Add(1)
	var fatalErrors []string
	logFatalln = func(args ...interface{}) {
		fatalErrors = append(fatalErrors, fmt.Sprint(args...))
	}

	config.initPeConfigWatcher(mockConn, &wg)

	logFatalln = log.Fatalln

	assert.Equal(t, testClusterName, config.configsByPe[*testClusterUuid].clusterName)
	assert.Equal(t, []string{testNodeIp}, config.configsByPe[*testClusterUuid].clusterExternalIps)
	assert.Equal(t, map[uuid4.Uuid]*zeusConfig.ConfigurationProto_Container{}, config.sspContainers)
	assert.Equal(t, 1, len(fatalErrors))
	mockConn.AssertExpectations(t)
}

func TestClusterName(t *testing.T) {
	utilConfig.TestSetupMockFactoryConfig()
	config := &configCache{clusterName: testClusterName}
	assert.Equal(t, testClusterName, config.ClusterName())
}

func TestClusterUuid(t *testing.T) {
	utilConfig.TestSetupMockFactoryConfig()
	config := &configCache{clusterUUID: testClusterUuid}
	assert.Equal(t, testClusterUuid, config.ClusterUuid())
}

func TestLocalNodeIp(t *testing.T) {
	utilConfig.TestSetupMockFactoryConfig()
	config := &configCache{localNodeIp: testNodeIp}
	assert.Equal(t, testNodeIp, config.LocalNodeIp())
}

func TestContainerName(t *testing.T) {
	containerId := int64(420)
	containerName := "TestContainer"
	containerUuid, _ := uuid4.New()
	confProto := &zeusConfig.ConfigurationProto{
		ContainerList: []*zeusConfig.ConfigurationProto_Container{
			{
				ContainerId:   &containerId,
				ContainerName: &containerName,
				ContainerUuid: proto.String(containerUuid.String()),
			},
		},
	}

	config := &configCache{configProto: confProto}
	name, err := config.ContainerName(containerId)

	assert.NoError(t, err)
	assert.Equal(t, containerName, name)
}

func TestPeClusterUuids(t *testing.T) {
	config := &configCache{}
	config.peConfigOnce.Do(func() {})
	config.configsByPe = make(map[uuid4.Uuid]peConfigCache)
	config.configsByPe[*testClusterUuid] = peConfigCache{clusterName: "0.0.0.1"}

	peList := config.PeClusterUuids()

	assert.Equal(t, peList, []*uuid4.Uuid{testClusterUuid})
}

func TestClusterExternalIpsExist(t *testing.T) {
	config := &configCache{}
	config.peConfigOnce.Do(func() {})
	config.configsByPe = make(map[uuid4.Uuid]peConfigCache)
	config.configsByPe[*testClusterUuid] = peConfigCache{clusterExternalIps: []string{"0.0.0.1", "1.2.3.4"}}

	externalIps := config.ClusterExternalIps(testClusterUuid)

	assert.Equal(t, config.configsByPe[*testClusterUuid].clusterExternalIps, externalIps)
}

func TestClusterExternalIpsNotExist(t *testing.T) {
	config := &configCache{}
	config.peConfigOnce.Do(func() {})
	config.configsByPe = make(map[uuid4.Uuid]peConfigCache)

	externalIps := config.ClusterExternalIps(testClusterUuid)

	assert.Nil(t, externalIps)
}

func TestPeClusterName(t *testing.T) {
	config := &configCache{}
	config.peConfigOnce.Do(func() {})
	config.configsByPe = make(map[uuid4.Uuid]peConfigCache)
	config.configsByPe[*testClusterUuid] = peConfigCache{clusterName: "testName"}

	clusterName := config.PeClusterName(testClusterUuid)

	assert.Equal(t, config.configsByPe[*testClusterUuid].clusterName, *clusterName)
}

func TestPeClusterNameNotExist(t *testing.T) {
	config := &configCache{}
	config.peConfigOnce.Do(func() {})
	config.configsByPe = make(map[uuid4.Uuid]peConfigCache)

	clusterName := config.PeClusterName(testClusterUuid)

	assert.Nil(t, clusterName)
}

func TestClusterSSPContainerUuidMap(t *testing.T) {
	config := &configCache{}
	config.peConfigOnce.Do(func() {})
	config.configsByPe = make(map[uuid4.Uuid]peConfigCache)
	config.sspContainers = map[uuid4.Uuid]*zeusConfig.ConfigurationProto_Container{}

	sspContainers := config.ClusterSSPContainerUuidMap()

	assert.Equal(t, config.sspContainers, sspContainers)
}

func TestCatalogPeRegisteredExist(t *testing.T) {
	config := &configCache{}
	config.peConfigOnce.Do(func() {})
	config.configsByPe = make(map[uuid4.Uuid]peConfigCache)
	config.configsByPe[*testClusterUuid] = peConfigCache{catalogPeRegistered: proto.Bool(true)}

	catalogPeRegistered := config.CatalogPeRegistered(testClusterUuid)

	assert.True(t, *catalogPeRegistered)
}

func TestCatalogPeRegisteredNotExist(t *testing.T) {
	config := &configCache{}
	config.peConfigOnce.Do(func() {})
	config.configsByPe = make(map[uuid4.Uuid]peConfigCache)

	catalogPeRegistered := config.CatalogPeRegistered(testClusterUuid)

	assert.Nil(t, catalogPeRegistered)
}
