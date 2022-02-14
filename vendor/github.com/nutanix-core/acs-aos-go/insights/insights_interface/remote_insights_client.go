package insights_interface

import (
	"fmt"
	"time"

	glog "github.com/golang/glog"
	proto "github.com/golang/protobuf/proto"
	client "github.com/nutanix-core/acs-aos-go/acropolis/client"
	aplos_transport "github.com/nutanix-core/acs-aos-go/aplos/client/transport"
	util_misc "github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
	util_net "github.com/nutanix-core/acs-aos-go/nutanix/util-go/net"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/acs-aos-go/zeus"
)

const (
	clientRetryInitialDelayMilliSecs = 250
	clientRetryMaxDelaySecs          = 2
	clientRetryTimeoutSecs           = 30
)

type InsightsRemoteRpcClient struct {
	client util_net.ProtobufRPCClientIfc
}

func NewRemoteInsights(zkSession *zeus.ZookeeperSession,
	clusterUUID, userUUID, tenantUUID *uuid4.Uuid,
	userName *string) *InsightsRemoteRpcClient {

	return &InsightsRemoteRpcClient{
		client: aplos_transport.NewRemoteProtobufRPCClientIfc(zkSession, insightServiceName,
			uint16(*InsightsPort), clusterUUID, userUUID,
			tenantUUID, 0 /*  timeout */, userName, nil),
	}
}

func InsightsRemoteClient(nosClusterUUID string,
	zkSession *zeus.ZookeeperSession) (*InsightsRemoteRpcClient, error) {

	clusterUUID, err := uuid4.StringToUuid4(nosClusterUUID)
	if err != nil {
		fmt.Printf("Unable to convert %s to UUID4: %q", nosClusterUUID,
			err.Error())
		return nil, err
	}

	// Setting to tenantUUID as below works. Need to revisit as what is expected.
	tenantUUIDStr := "00000000-0000-0000-0000-000000000000"
	tUUID, err := uuid4.StringToUuid4(tenantUUIDStr)
	if err != nil {
		fmt.Printf("Unable to convert %s to UUID4: %q", nosClusterUUID,
			err.Error())
		return nil, err
	}

	remoteClient := NewRemoteInsights(zkSession,
		clusterUUID, nil, tUUID, nil)
	return remoteClient, nil
}

func (svc *InsightsRemoteRpcClient) SendMsg(
	service string, request, response proto.Message, timeoutSecs int64) error {

	backoff := util_misc.NewExponentialBackoffWithTimeout(
		time.Duration(clientRetryInitialDelayMilliSecs)*time.Millisecond,
		time.Duration(clientRetryMaxDelaySecs)*time.Second,
		time.Duration(clientRetryTimeoutSecs)*time.Second)
	retryCount := 0
	var err error
	for {
		err = svc.client.CallMethodSync(insightServiceName, service, request,
			response, timeoutSecs)
		if err == nil {
			return nil
		}
		appErr, ok := util_net.ExtractAppError(err)
		if ok {
			/* app error */
			err = client.Errors[appErr.GetErrorCode()]
			if !client.ErrRetry.Equals(err) {
				return err
			}
		} else if !util_net.IsRpcTransportError(err) {
			/* Not app error and not transport error */
			return err
		}
		/* error is retriable. Continue with retry */
		waited := backoff.Backoff()
		if waited == util_misc.Stop {
			glog.Errorf("Timed out retrying Insights RPC: %s. Request: %s",
				service, proto.MarshalTextString(request))
			break
		}
		retryCount += 1
		glog.Errorf("Error while executing %s: %s. Retry count: %d",
			service, err.Error(), retryCount)
	} // Continue the for loop.
	return err
}
