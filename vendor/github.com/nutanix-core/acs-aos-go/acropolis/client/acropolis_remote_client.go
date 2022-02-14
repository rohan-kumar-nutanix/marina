package acropolis_client

import (
	"time"

	glog "github.com/golang/glog"
	proto "github.com/golang/protobuf/proto"
	aplos_transport "github.com/nutanix-core/acs-aos-go/aplos/client/transport"
	util_misc "github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
	util_net "github.com/nutanix-core/acs-aos-go/nutanix/util-go/net"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/acs-aos-go/zeus"
)

//-----------------------------------------------------------------------------
// Use this function to instantiate new remote acropolis service.
// Args:
//   zkSession    : Zookeeper session for configuration.
//   clusterUUID  : UUID of the remote cluster.
//   userUUID     : User UUID for remote protobuf client.
//   tenantUUID   : Tenant UUID for remote protobuf client.
// Return:
//  Pointer to Remote anduril structure.

type AcropolisRemoteRpcClient struct {
	client util_net.ProtobufRPCClientIfc
}

func NewRemoteAcropolis(zkSession *zeus.ZookeeperSession,
	clusterUUID, userUUID, tenantUUID *uuid4.Uuid,
	userName *string) *AcropolisRemoteRpcClient {
	return &AcropolisRemoteRpcClient{
		client: aplos_transport.NewRemoteProtobufRPCClientIfc(zkSession,
			acropolisServiceName, defaultAcropolisPort, clusterUUID,
			userUUID, tenantUUID, 0 /* timeout */, userName, nil),
	}
}

func (svc *AcropolisRemoteRpcClient) SendMsg(
	service string, request, response proto.Message, timeoutSecs int64) error {

	backoff := util_misc.NewExponentialBackoffWithTimeout(
		time.Duration(clientRetryInitialDelayMilliSecs)*time.Millisecond,
		time.Duration(clientRetryMaxDelaySecs)*time.Second,
		time.Duration(clientRetryTimeoutSecs)*time.Second)
	retryCount := 0
	var err error
	for {
		err = svc.client.CallMethodSync(acropolisServiceName, service, request,
			response, timeoutSecs)
		if err == nil {
			return nil
		}
		appErr, ok := util_net.ExtractAppError(err)
		if ok {
			/* app error */
			err = Errors[appErr.GetErrorCode()]
			if !ErrRetry.Equals(err) {
				return err
			}
		} else if !util_net.IsRpcTransportError(err) {
			/* Not app error and not transport error */
			return err
		}
		/* error is retriable. Continue with retry */
		waited := backoff.Backoff()
		if waited == util_misc.Stop {
			glog.Errorf("Timed out retrying Acropolis RPC: %s. Request: %s",
				service, proto.MarshalTextString(request))
			break
		}
		retryCount += 1
		glog.Errorf("Error while executing %s: %s. Retry count: %d",
			service, err.Error(), retryCount)
	} // Continue the for loop.
	return err
}
