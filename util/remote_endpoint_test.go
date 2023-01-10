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

	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	mockZeus "github.com/nutanix-core/content-management-marina/mocks/zeus"
)

func TestIsPE(t *testing.T) {
	endpoint := RemoteEndpoint{}
	assert.False(t, endpoint.IsPE())

	uuid, _ := uuid4.New()
	endpoint = RemoteEndpoint{RemoteClusterUuid: *uuid}
	assert.True(t, endpoint.IsPE())
}

func TestGetUserVisibleIdClusterName(t *testing.T) {
	uuid, _ := uuid4.New()
	endpoint := RemoteEndpoint{RemoteClusterUuid: *uuid}
	clusterName := "Foo"
	mockCache := &mockZeus.ConfigCache{}
	mockCache.On("PeClusterName", mock.Anything).Return(&clusterName).Once()

	id := endpoint.GetUserVisibleId(mockCache)

	assert.Equal(t, *id, clusterName)
	mockCache.AssertExpectations(t)
}

func TestGetUserVisibleIdIP(t *testing.T) {
	uuid, _ := uuid4.New()
	endpoint := RemoteEndpoint{RemoteClusterUuid: *uuid}
	ips := []string{"127.0.0.1"}
	mockCache := &mockZeus.ConfigCache{}
	mockCache.On("PeClusterName", mock.Anything).Return(nil).Once()
	mockCache.On("ClusterExternalIps", mock.Anything).Return(ips).Once()

	id := endpoint.GetUserVisibleId(mockCache)

	assert.Equal(t, *id, ips[0])
	mockCache.AssertExpectations(t)
}

func TestGetUserVisibleIdUuid(t *testing.T) {
	uuid, _ := uuid4.New()
	endpoint := RemoteEndpoint{RemoteClusterUuid: *uuid}
	mockCache := &mockZeus.ConfigCache{}
	mockCache.On("PeClusterName", mock.Anything).Return(nil).Once()
	mockCache.On("ClusterExternalIps", mock.Anything).Return(nil).Once()

	id := endpoint.GetUserVisibleId(mockCache)

	assert.Equal(t, *id, uuid.String())
	mockCache.AssertExpectations(t)
}

func TestGetUserVisibleIdNotPe(t *testing.T) {
	endpoint := RemoteEndpoint{}
	mockCache := &mockZeus.ConfigCache{}

	id := endpoint.GetUserVisibleId(mockCache)

	assert.Nil(t, id)
}
