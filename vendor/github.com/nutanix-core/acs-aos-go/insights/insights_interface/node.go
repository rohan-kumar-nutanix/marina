/*
 * Copyrigright (c) 2017 Nutanix Inc. All rights reserved.
 *
 * Adds helper methods to 'Entity' protobuf generated object.
 */

package insights_interface

import (
  "regexp"
)

type Node struct {
  entity EntityIfc
}

func NewNode(entity EntityIfc) *Node {
  return &Node{
    entity: entity,
  }
}

const (
  NumBlock     = "(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
  RegexPattern = NumBlock + "\\." + NumBlock + "\\." + NumBlock + "\\." + NumBlock
)

var (
  ipRegex = regexp.MustCompile(RegexPattern)
)

// Get the host IPs of 'node'.
func (node *Node) GetHostIPs() ([]string, error) {
  result := make([]string, 0)
  ips, err := node.entity.GetStringList("ipv4_addresses")
  if err != nil {
    return nil, err
  }
  for _, ip := range ips {
    if ip := ipRegex.FindString(ip); ip != "" {
      result = append(result, ip)
    }
  }
  return result, nil
}
