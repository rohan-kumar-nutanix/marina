/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Author: rishabh.gupta@nutanix.com
 */

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nutanix-core/acs-aos-go/ergon"
	ergonClient "github.com/nutanix-core/acs-aos-go/ergon/client"
	mockErgon "github.com/nutanix-core/acs-aos-go/ergon/client/mocks"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-slbufs/util/sl_bufs/net"
	"github.com/nutanix-core/acs-aos-go/zeus"
)

type MockErgonClient struct {
	client *mockErgon.Ergon
}

func (mock *MockErgonClient) NewRemoteErgonService(zkSession *zeus.ZookeeperSession, clusterUUID, userUUID,
	tenantUUID *uuid4.Uuid, userName *string, reqContext *net.RpcRequestContext) ergonClient.Ergon {
	return mock.client
}

func TestRemoteEndpointPoll(t *testing.T) {
	task := &ergon.Task{}
	ret := &ergon.TaskPollRet{CompletedTasks: []*ergon.Task{task}}
	mockErgonObj := &mockErgon.Ergon{}
	mockErgonObj.On("TaskPoll", mock.Anything).Return(ret, nil).Once()
	mockErgonClient := &MockErgonClient{client: mockErgonObj}
	ergonService = mockErgonClient.NewRemoteErgonService

	taskuuid, _ := uuid4.New()
	retTask := remoteEndpointPoll(&zeus.ZookeeperSession{}, RemoteEndpoint{}, taskuuid)

	assert.Equal(t, task, retTask)
	mockErgonObj.AssertExpectations(t)
}

func TestRemoteEndpointPollError(t *testing.T) {
	mockErgonObj := &mockErgon.Ergon{}
	mockErgonObj.On("TaskPoll", mock.Anything).Return(nil, ergonClient.ErrNotFound).Once()
	mockErgonObj.On("TaskPoll", mock.Anything).Return(nil, ergonClient.ErrCanceled).Once()
	mockErgonClient := &MockErgonClient{client: mockErgonObj}
	ergonService = mockErgonClient.NewRemoteErgonService

	taskuuid, _ := uuid4.New()
	retTask := remoteEndpointPoll(&zeus.ZookeeperSession{}, RemoteEndpoint{}, taskuuid)

	assert.Nil(t, retTask)
	mockErgonObj.AssertExpectations(t)
}

func TestPollAllRemoteTasks(t *testing.T) {
	taskUuid, _ := uuid4.New()
	clusterUuid, _ := uuid4.New()
	remoteEndpoint := RemoteEndpoint{RemoteClusterUuid: *clusterUuid}
	taskUuidByEndpoint := make(map[RemoteEndpoint]*uuid4.Uuid)
	taskUuidByEndpoint[remoteEndpoint] = taskUuid
	fanoutUtil := &FanoutTaskPollerUtil{}
	task := &ergon.Task{}
	ret := &ergon.TaskPollRet{CompletedTasks: []*ergon.Task{task}}
	mockErgonObj := &mockErgon.Ergon{}
	mockErgonObj.On("TaskPoll", mock.Anything).Return(ret, nil).Once()
	mockErgonClient := &MockErgonClient{client: mockErgonObj}
	ergonService = mockErgonClient.NewRemoteErgonService

	taskByRemoteEndpoint := *fanoutUtil.PollAllRemoteTasks(&zeus.ZookeeperSession{}, &taskUuidByEndpoint)

	assert.Equal(t, 1, len(taskByRemoteEndpoint))
	assert.Equal(t, task, taskByRemoteEndpoint[remoteEndpoint])
	mockErgonObj.AssertExpectations(t)
}
