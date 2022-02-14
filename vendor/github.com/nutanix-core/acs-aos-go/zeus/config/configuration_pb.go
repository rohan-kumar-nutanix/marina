/*
 * Copyright (c) 2017 Nutanix Inc. All rights reserved.
 *
 * Add helper routines to proto generated 'ConfigurationProto'
 *
 */

package zeus_config

import (
        "errors"
        "fmt"

        "github.com/golang/glog"

        util_config "github.com/nutanix-core/acs-aos-go/nutanix/util-go/config"
        ntnx_errors "github.com/nutanix-core/acs-aos-go/nutanix/util-go/errors"
        "github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
)

var (
        ErrNotFound = ntnx_errors.NewNtnxError("Not found", 1)
)

func (config *ConfigurationProto) GetNodeIp(
        nodeUuid *uuid4.Uuid) (string, error) {

        for _, node := range config.GetNodeList() {
                if node.GetUuid() == nodeUuid.String() {
                        return node.GetBackpaneIpFromNodeInfo(), nil
                }
        }
        return "", ErrNotFound
}

func (config *ConfigurationProto) MustGetIp() string {
        ip, err := config.GetThisNodeIp()
        if err != nil {
                glog.Fatal(err)
        }
        return ip
}

func (config *ConfigurationProto) GetThisNodeIp() (string, error) {
        uuid, err := util_config.GetUuidFromFactoryConfig()
        if err != nil {
                return "", err
        }
        return config.GetNodeIp(uuid)
}

func (config *ConfigurationProto) IsClusterUpgrading() bool {
        releaseVer := config.GetReleaseVersion()
        var upgrading bool = false
        for _, node := range config.GetNodeList() {
                nodeVer := node.GetSoftwareVersion()
                if nodeVer != releaseVer {
                        upgrading = true
                        break
                }
        }
        return upgrading
}

func (config *ConfigurationProto) MustGetNodeIndex() int {
        i, err := config.NodeIndex()
        if err != nil {
                glog.Fatal(err)
        }
        return i
}

// Return the node's position in the cluster, z-indexed.
//
// NOTE: uses a O(n) to find node index. Does not sort nodes. Sorting nodes
// not a good idea b/c it will mutate 'config', which prevents parallelism.
//
// Examples:
// * For a 8-node cluster, valid return values [0..7]
// * For a 11-node cluster, valid return values [0..10]
// * For a 4-node cluster, valid return values [0..3]
func (config *ConfigurationProto) NodeIndex() (int, error) {
        svmId, err := config.NodeSvmId()
        if err != nil {
                return -1, err
        }

        // Count the number of nodes below 'svmId'.
        numBelow := 0
        for _, node := range config.GetNodeList() {
                if node.GetServiceVmId() == svmId {
                        continue
                }
                if node.GetServiceVmId() < svmId {
                        numBelow += 1
                }
        }
        return numBelow, nil
}

func (config *ConfigurationProto) MustGetSvmId() int64 {
        svmId, err := config.NodeSvmId()
        if err != nil {
                glog.Fatal(err)
        }
        return svmId
}

func (config *ConfigurationProto) NodeSvmId() (int64, error) {
        uuid, err := util_config.GetUuidFromFactoryConfig()
        if err != nil {
                return -1, err
        }

        uuidStr := uuid.String()
        for _, node := range config.GetNodeList() {
                if node.GetUuid() == uuidStr {
                        return node.GetServiceVmId(), nil
                }
        }
        return -1, errors.New("Couldn't find service vm id. Node UUID is" + uuidStr)
}

func (config *ConfigurationProto) NumNodes() int {
        return len(config.GetNodeList())
}

// For a given container id, fetches container name and
// its uuid from the configuration.
//
func (config *ConfigurationProto) GetContainerInfo(containerId int64) (
        string, string, error) {

        containers := config.GetContainerList()
        for _, container := range containers {
                if container.GetContainerId() == containerId {
                        return container.GetContainerName(),
                                        container.GetContainerUuid(), nil
                }
        }

        return "", "", fmt.Errorf("Couldn't find container for id:%d", containerId)
}
