/*
 * Copyright (c) 2023 Nutanix Inc. All rights reserved.
 *
 * Author: rishabh.gupta@nutanix.com
 *
 * This file is to make RPC requests to Remote Catalog RPC Server
 *
 */

package catalog_client

import (
	"time"

	"github.com/golang/protobuf/proto"
	log "k8s.io/klog/v2"

	aplos_transport "github.com/nutanix-core/acs-aos-go/aplos/client/transport"
	utilMisc "github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
	utilNet "github.com/nutanix-core/acs-aos-go/nutanix/util-go/net"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/acs-aos-go/zeus"
	"github.com/nutanix-core/content-management-marina/common"
	utils "github.com/nutanix-core/content-management-marina/util"
	marinaZeus "github.com/nutanix-core/content-management-marina/zeus"
)

const (
	clientRetryInitialDelayMilliSecs = 250
	clientRetryMaxDelaySecs          = 2
	clientRetryTimeoutSecs           = 30
)

type RemoteCatalogInterface interface {
	SendMsg(service string, request, response proto.Message) error
}

type RemoteCatalog struct {
	client utilNet.ProtobufRPCClientIfc
}

func NewRemoteCatalogService(zkSession *zeus.ZookeeperSession, zeusConfig marinaZeus.ConfigCache,
	clusterUUID, userUUID, tenantUUID *uuid4.Uuid, userName *string) RemoteCatalogInterface {

	return &RemoteCatalog{
		client: aplos_transport.NewRemoteProtobufRPCClientIfcWithSource(zkSession, catalogServiceName,
			uint16(*common.CatalogPort), clusterUUID, userUUID, tenantUUID,
			0 /* timeout */, userName, nil /* request context */, zeusConfig.ClusterUuid()),
	}
}

func (svc *RemoteCatalog) SendMsg(service string, request, response proto.Message) error {
	log.Info("Sending RPC request to Remote Catalog Service for method : ", service)
	return svc.sendMsg(service, request, response)
}

func (svc *RemoteCatalog) sendMsg(service string, request, response proto.Message) error {
	backoff := utilMisc.NewExponentialBackoffWithTimeout(
		time.Duration(clientRetryInitialDelayMilliSecs)*time.Millisecond,
		time.Duration(clientRetryMaxDelaySecs)*time.Second,
		time.Duration(clientRetryTimeoutSecs)*time.Second)
	retryCount := 0
	var err error
	for {
		err = svc.client.CallMethodSync(catalogServiceName, service, request, response, utils.RpcTimeoutSec)
		if err == nil {
			return nil
		}
		appErr, ok := utilNet.ExtractAppError(err)
		if ok {
			err = Errors[appErr.GetErrorCode()]
			if !ErrRetry.Equals(err) {
				return err
			}
		} else if !utilNet.IsRpcTransportError(err) {
			return err
		}
		waited := backoff.Backoff()
		if waited == utilMisc.Stop {
			log.Errorf("Timed out retrying Catalog RPC: %s. Request: %s",
				service, proto.MarshalTextString(request))
			break
		}
		retryCount += 1
		log.Errorf("Error while executing %s: %s. Retry count: %d",
			service, err.Error(), retryCount)
	}
	return err
}
