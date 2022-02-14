//
// Copyright (c) 2016 Nutanix Inc. All rights reserved.
//
// Author: abhijit.paithankar@nutanix.com
//
// This module provides access to ZooKeeper client configuration file (zoo.cfg)

package zeus_config

import (
        "bufio"
        "net"
        "os"
        "regexp"
        "strconv"
        "strings"

        glog "github.com/golang/glog"
)

const (
        // Default path for zookeeper configuration.
        DEFAULT_ZOOKEEPER_CONFIG_PATH = "/home/nutanix/config/zookeeper/zoo.cfg"
)

var (
        // Regular expression to retrieve a zookeeper node id.
        SERVER_KEY_REGEX, _ = regexp.Compile(`^server\.([0-9]+)`)

        // Regular expression to retrieve a zookeeper node ip address.
        SERVER_VALUE_REGEX, _ = regexp.Compile(`^([0-9.]+|zk\d):[0-9]+(:[0-9]+)?$`)
)

type ZookeeperConfig struct {
        dict  map[string]string
        nodes map[string]int
}

// Initializes a new ZookeeperConfig with the given zookeeper configuration file.
func NewZookeeperConfig(zooCfg string) *ZookeeperConfig {
        zkDict := make(map[string]string)
        zkNodes := make(map[string]int)

        // Read zooCfg and populate maps above
        file, err := os.Open(zooCfg)
        if err != nil {
                glog.Fatal("Failed to open the zookeeper configuration file ", err)
                return nil
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                if err := scanner.Err(); err != nil {
                        glog.Error("Failed to read the zookeeper configuration file ", err)
                        continue
                }

                line := scanner.Text()
                if err := scanner.Err(); err != nil {
                        glog.Error("Failed to read zookeeper configuration file ", err)
                        continue
                }

                // Skip comments.
                line = strings.TrimSpace(strings.Split(line, "#")[0])
                if len(line) == 0 {
                        continue
                }

                words := strings.Split(line, "=")
                key := strings.TrimSpace(words[0])
                value := strings.TrimSpace(words[1])
                zkDict[key] = value

                // Parse out sever's node id and host name.
                matchedKey := SERVER_KEY_REGEX.MatchString(key)
                matchedVal := SERVER_VALUE_REGEX.MatchString(value)
                if !matchedKey || !matchedVal {
                        continue
                }
                nodeId := SERVER_KEY_REGEX.FindAllString(key, 1)
                host := SERVER_VALUE_REGEX.FindAllStringSubmatch(value, 1)

                // Find the server host ip address.
                hostName := host[0][1]
                hostIp, err := net.ResolveIPAddr("ip4", hostName)
                if err != nil {
                        glog.Errorf("Failed to get IP address for host: %s with: %s", hostName, err)
                        return nil
                }
                hostIpStr := hostIp.IP.String()

                id, err := strconv.Atoi(nodeId[0])
                if err != nil {
                        zkNodes[hostIpStr] = id
                } else {
                        glog.Error("Failed to parse zookeeper configuration file ", err)
                        return nil
                }
        }
        return &ZookeeperConfig{dict: zkDict, nodes: zkNodes}
}

// Returns a copy of the nodes ip to id map.
func (self *ZookeeperConfig) ZkIdMap() map[string]int {
        newMap := make(map[string]int)
        for k, v := range self.nodes {
                newMap[k] = v
        }
        return newMap
}

// Returns zookeeper configuration value for the key 'key'.
func (self *ZookeeperConfig) Get(key string) string {
        return self.dict[key]
}

// Returns zookeeper node id for a node with ip address 'ip'.
func (self *ZookeeperConfig) ZookeeperNodeId(ip string) int {
        return self.nodes[ip]
}
