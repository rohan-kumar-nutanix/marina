/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Author: rishabh.gupta@nutanix.com
*
* The implementation for Fanout Task Poller
 */

package utils

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	log "k8s.io/klog/v2"

	ergonIfc "github.com/nutanix-core/acs-aos-go/ergon"
	ergonClient "github.com/nutanix-core/acs-aos-go/ergon/client"
	util_net "github.com/nutanix-core/acs-aos-go/nutanix/util-go/net"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/acs-aos-go/zeus"
)

type FanoutTaskPollerUtil struct {
}

type FanoutTaskPollerInterface interface {
	PollAllRemoteTasks(zkSession *zeus.ZookeeperSession,
		taskUuidByEndpoint *map[RemoteEndpoint]*uuid4.Uuid) *map[RemoteEndpoint]*ergonIfc.Task
}

type EndpointTaskPair struct {
	RemoteEndpoint
	*ergonIfc.Task
}

var (
	ergonService = ergonClient.NewRemoteErgonService
)

const progressIntervalSecs = uint64(60)

func remoteEndpointPoll(zkSession *zeus.ZookeeperSession, remoteEndpoint RemoteEndpoint, taskUuid *uuid4.Uuid) *ergonIfc.Task {

	requestId := fmt.Sprintf("PollTask:%s", taskUuid.String())
	arg := &ergonIfc.TaskPollArg{
		RequestId:            &requestId,
		CompletedTasksFilter: [][]byte{taskUuid.RawBytes()},
		TimeoutSec:           proto.Uint64(progressIntervalSecs),
	}

	log.Infof("Started Polling for task : %s", taskUuid.String())
	client := ergonService(zkSession, &remoteEndpoint.RemoteClusterUuid, nil, nil, nil, nil)

	var task *ergonIfc.Task
	var ret *ergonIfc.TaskPollRet
	var err error
	for {
		ret, err = client.TaskPoll(arg)
		if util_net.IsRpcTransportError(err) || err == ergonClient.ErrRetry || err == ergonClient.ErrNotFound {
			log.Warningf("Retrying polling on remote endpoint %s for task %s",
				remoteEndpoint.RemoteClusterUuid.String(), taskUuid.String())
			continue

		} else if err != nil {
			log.Errorf("Polling remote endpoint %s for task %s failed: %v",
				remoteEndpoint.RemoteClusterUuid.String(), taskUuid.String(), err)
			break
		}

		if len(ret.CompletedTasks) == 1 {
			log.Infof("Polling remote endpoint %s for task %s completed",
				remoteEndpoint.RemoteClusterUuid.String(), taskUuid.String())
			task = ret.CompletedTasks[0]
			break
		}
	}
	return task
}

func (poller *FanoutTaskPollerUtil) PollAllRemoteTasks(zkSession *zeus.ZookeeperSession,
	taskUuidByEndpoint *map[RemoteEndpoint]*uuid4.Uuid) *map[RemoteEndpoint]*ergonIfc.Task {

	endpointTaskChan := make(chan EndpointTaskPair)
	count := 0
	for remoteEndpoint, taskUuid := range *taskUuidByEndpoint {
		count++
		go func(remoteEndpoint RemoteEndpoint, taskUuid *uuid4.Uuid, endpointTaskChan chan EndpointTaskPair) {
			task := remoteEndpointPoll(zkSession, remoteEndpoint, taskUuid)
			endpointTaskChan <- EndpointTaskPair{remoteEndpoint, task}
		}(remoteEndpoint, taskUuid, endpointTaskChan)
	}

	taskByRemoteEndpoint := make(map[RemoteEndpoint]*ergonIfc.Task)
	for i := 0; i < count; i++ {
		endpointTask := <-endpointTaskChan
		taskByRemoteEndpoint[endpointTask.RemoteEndpoint] = endpointTask.Task
	}
	return &taskByRemoteEndpoint
}
