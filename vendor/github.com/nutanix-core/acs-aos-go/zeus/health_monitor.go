/*
 * Copyright (c) 2020 Nutanix Inc. All rights reserved.
 *
 * Author: hitesh.bhagchandani@nutanix.com
 *
 * This module provides the basic functionality to create an ephemeral node of
 * a component in Zookeeper for health-monitoring. Components can monitor their
 * health and that of other components by registering watches on these nodes
 * which are fired when the status of these nodes change.
 *
 * Example usage:
 *    Each process can choose to create a local health-monitor object.
 *        healthMonitor = NewHealthMonitor(zkSession)
 *    And register its health with Zookeeper, its handle being of the form
 *    ip:port.
 *        healthMonitor.Register(<process_name>, <ip:port>, <peer_id>,
 *                               <incarnation_id>)
 */

package zeus

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/nutanix-core/acs-aos-go/go-zookeeper"
)

var (
	healthMonitorMaxNumRetries = flag.Uint("zeus_health_register_max_num_retries", 5,
		"Maximum number of retries in zeus health-registering")

	healthMonitorMaxBackoffTime = flag.Uint("zeus_health_register_max_backoff_time", 32,
		"Maximum backoff time (seconds) in zeus health-registering")
)

type HealthMonitor struct {
	ZkSession *ZookeeperSession
}

// Create a HealthMonitor object that verifies that the root node:
// HEALTH_MONITOR_ROOT is present in Zookeeper already.
func NewHealthMonitor(zkSession *ZookeeperSession) *HealthMonitor {

	if !zkSession.WaitForConnection() {
		return nil
	}

	// zkpath 'HEALTH_MONITOR_ROOT' is always created when the Zookeeper
	// server starts. So we expect the path to be present beforehand.
	exist, _, err := zkSession.Exist(HEALTH_MONITOR_ROOT, true /* retry */)
	if err != nil {
		glog.Fatalf("Unable to check whether zknode '%s' exist or not: %s", HEALTH_MONITOR_ROOT, err)
	}
	if !exist {
		glog.Fatalf("The zknode '%s' doesn't exist in zookeeper", HEALTH_MONITOR_ROOT)
	}

	return &HealthMonitor{ZkSession: zkSession}
}

// Register a component's health with Zookeeper by creating an ephemeral
// node in the HEALTH_MONITOR_ROOT/healthNamespace/peerHandle
// namespace.
// Arguments -
//  healthNamespace: The namespace used to limit the scope of health
//                   monitoring, i.e., rather than a flat universe comprising
//                   of every process, one can watch and monitor a related
//                   set of processes that are in a common namespace.
//  peerHandle: Of the form ip:port.
//  peerId: A unique id assigned to this component; it remains the same
//          irrespective of the number of restarts it goes through.
//  incarnationId: A unique integer that identifies an incarnation of the given
//                 component.
func (healthMonitor *HealthMonitor) Register(healthNamespace string, peerHandle string, peerId uint64, incarnationId uint64) error {
	healthNamespacePath := fmt.Sprintf("%s/%s", HEALTH_MONITOR_ROOT, healthNamespace)
	healthNodePath := fmt.Sprintf("%s/%s", healthNamespacePath, peerHandle)
	retry := uint(0)
	backoff := uint(1)
	var err error
	data := []byte(fmt.Sprintf("peer_id: %d incarnation_id: %d", peerId, incarnationId))
	for {
		_, err = healthMonitor.ZkSession.Create(healthNodePath, data, 1 /* flags */, WORLD_ACL, true /* retry */)
		if err == nil {
			break
		}
		// The create request is expected to return with no errors other than
		// 'ErrNoNode' which implies that the parent node heirarchy doesn't
		// exist. Incase of an unexpected error, return the error signifying
		// failure to register the health.
		if err == zk.ErrNodeExists {
			glog.Fatal("The node should not exist already as we expect all the previous sessions to be closed gracefully")
		}
		if err != zk.ErrNoNode {
			glog.Errorf("Unable to create ephemeral health zknode %s: %s", healthNodePath, err)
			return err
		}
		glog.Infof("Ephemeral health node '%s' doesn't exist in zookeeper", healthNamespacePath)
		_, err := healthMonitor.ZkSession.Create(healthNamespacePath, []byte{} /* data */, 0 /* flags */, WORLD_ACL, true /* retry */)
		if err != nil && err != zk.ErrNodeExists {
			glog.Errorf("Unable to create zknode '%s': %s", healthNamespacePath, err)
			return err
		}
		// If we are here, it implies that the node at the path
		// 'healthNamespacePath' exists and we need to retry creating the node
		// at the path 'healthNodePath'.
		maybeRetryAndFatal(&retry, &backoff, *healthMonitorMaxNumRetries, *healthMonitorMaxBackoffTime)
	}

	glog.Infof("Successfully registered health in health namespace: %s with handle: %s", healthNamespacePath, peerHandle)
	return nil
}

// TODO(Hitesh) : Add functionality to support components to set watches on other components' health.
func (healthMonitor *HealthMonitor) StartHealthWatch(healthNamespace string, changeCb func(string)) {
	glog.Fatal("Not implemented")
}
