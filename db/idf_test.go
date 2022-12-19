/*
* Copyright (c) 2021 Nutanix Inc. All rights reserved.
*
* Author: rajesh.battala@nutanix.com
*
* IDF utilities for Marina Service.
 */

package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdfClientWithRetry(t *testing.T) {
	idfClient := IdfClientWithRetry()
	idfObj := idfClient.(*IdfClient)
	assert.NotNil(t, idfObj)
	assert.NotNil(t, idfObj.IdfSvc)
	assert.Nil(t, idfObj.Retry)
}

func TestIdfClientWithoutRetry(t *testing.T) {
	idfClient := IdfClientWithoutRetry()
	idfObj := idfClient.(*IdfClient)
	assert.NotNil(t, idfObj)
	assert.NotNil(t, idfObj.IdfSvc)
	assert.NotNil(t, idfObj.Retry)
}
