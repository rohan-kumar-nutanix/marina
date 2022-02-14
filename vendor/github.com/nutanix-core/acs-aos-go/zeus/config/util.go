/*
 * Copyright (c) 2017 Nutanix Inc. All rights reserved.
 *
 * Stateless functions around 'zeus/config'. Very convenient, particularly if
 * you don't intend to save the configuration proto, zksession, etc. But
 * expensive b/c getting a zksession each time is expensive.
 */

package zeus_config

import (
        "os"
        "strings"

        "github.com/golang/glog"
)

const (
        zkHostPortList = "ZOOKEEPER_HOST_PORT_LIST"
)

//-----------------------------------------------------------------------------

func NumNodes() int {
        var _n int
        GetZeusConfigWrapper(func(config *Configuration) {
                _n = len(config.ConfigProto.GetNodeList())
        })
        return _n
}

//-----------------------------------------------------------------------------

func GetThisNodeIp() string {
        var _ip string
        GetZeusConfigWrapper(func(config *Configuration) {
                ip, err := config.ConfigProto.GetThisNodeIp()
                if err != nil || ip == "" {
                        glog.Fatal(err, ip)
                }
                _ip = ip
        })
        return _ip
}

//-----------------------------------------------------------------------------

func GetThisNodeIndex() int {
        var _idx int
        GetZeusConfigWrapper(func(config *Configuration) {
                idx, err := config.ConfigProto.NodeIndex()
                if err != nil || idx < 0 {
                        glog.Fatal(err, idx)
                }
                _idx = idx
        })
        return _idx
}

//-----------------------------------------------------------------------------

func GetNodeSvmId() int64 {
        var svmid int64
        var err error
        err = nil
        GetZeusConfigWrapper(func(config *Configuration) {
                svmid, err = config.ConfigProto.NodeSvmId()
                if err != nil {
                        glog.Fatal(err)
                }
        })
        return svmid
}

//-----------------------------------------------------------------------------

func GetZeusConfigWrapper(callback func(*Configuration)) {
        zeusConfig, err :=
                NewConfiguration(nil /* zkSession */, nil /* hostPortList */)
        if err != nil {
                glog.Fatal(err)
        }

        callback(zeusConfig)

        // This is the reason why this function call is expensive. If you are going
        // to make this call many times, then use something else.
        zeusConfig.ZKSession.Close()
}

//-----------------------------------------------------------------------------

func GetZkHostPortList() []string {
        return strings.Split(os.Getenv(zkHostPortList), ",")
}

//-----------------------------------------------------------------------------

func GetZkPort() string {
        hostportvec := GetZkHostPortList()
        if len(hostportvec) == 0 || strings.TrimSpace(hostportvec[0]) == "" {
                glog.Fatal(
                        "Should have been able to obtain ZkHostPortList from environment ",
                        "variable 'ZOOKEEPER_HOST_PORT_LIST'")
                return ""
        }
        // Assumes the port on each hostport is the same.
        hostport := hostportvec[0]
        pieces := strings.Split(hostport, ":")
        port := pieces[len(pieces)-1]
        return port
}
