/*
 * Copyright (c) 2017 Nutanix Inc. All rights reserved.
 *
 * Helper routines added to proto-generated 'ConfigurationProto_Node'.
 */

package zeus_config

type NodeList []*ConfigurationProto_Node

func (node *ConfigurationProto_Node) GetBackpaneIpFromNodeInfo() string {
        if node.ControllerVmBackplaneIp != nil {
                return node.GetControllerVmBackplaneIp()
        }
        return node.GetServiceVmExternalIp()
}

// Len is the number of elements in the collection.
func (n NodeList) Len() int {
        return len(n)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (n NodeList) Less(i, j int) bool {
        return n[i].GetServiceVmId() < n[j].GetServiceVmId()
}

// Swap swaps the elements with indexes i and j.
func (n NodeList) Swap(i, j int) {
        tmp := n[i]
        n[i] = n[j]
        n[j] = tmp
}
