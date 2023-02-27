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
	assert.False(t, IsPE(NilUuid))

	uuid, _ := uuid4.New()
	assert.True(t, IsPE(*uuid))
}

func TestGetUserVisibleIdClusterName(t *testing.T) {
	uuid, _ := uuid4.New()
	clusterName := "Foo"
	mockCache := &mockZeus.ConfigCache{}
	mockCache.On("PeClusterName", mock.Anything).Return(&clusterName).Once()

	id := GetUserVisibleId(mockCache, *uuid)

	assert.Equal(t, id, clusterName)
	mockCache.AssertExpectations(t)
}

func TestGetUserVisibleIdIP(t *testing.T) {
	uuid, _ := uuid4.New()
	ips := []string{"127.0.0.1"}
	mockCache := &mockZeus.ConfigCache{}
	mockCache.On("PeClusterName", mock.Anything).Return(nil).Once()
	mockCache.On("ClusterExternalIps", mock.Anything).Return(ips).Once()

	id := GetUserVisibleId(mockCache, *uuid)

	assert.Equal(t, id, ips[0])
	mockCache.AssertExpectations(t)
}

func TestGetUserVisibleIdUuid(t *testing.T) {
	uuid, _ := uuid4.New()
	mockCache := &mockZeus.ConfigCache{}
	mockCache.On("PeClusterName", mock.Anything).Return(nil).Once()
	mockCache.On("ClusterExternalIps", mock.Anything).Return(nil).Once()

	id := GetUserVisibleId(mockCache, *uuid)

	assert.Equal(t, id, uuid.String())
	mockCache.AssertExpectations(t)
}

func TestGetUserVisibleIdNotPe(t *testing.T) {
	mockCache := &mockZeus.ConfigCache{}

	id := GetUserVisibleId(mockCache, NilUuid)

	assert.Equal(t, id, "")
}
