/*
 * Copyrigright (c) 2016 Nutanix Inc. All rights reserved.
 *
 * Defines client interface and utility functions for Insight Service
 */

package insights_interface

import (
  "flag"
  "time"

  ntnx_errors "github.com/nutanix-core/acs-aos-go/nutanix/util-go/errors"
  util_misc "github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
  util_net "github.com/nutanix-core/acs-aos-go/nutanix/util-go/net"

  "github.com/golang/glog"
  "github.com/golang/protobuf/proto"
)

const (
  insightServiceName = "nutanix.insights.interface.InsightsRpcSvc"
  // Special Incarnation ID that should be used for doing the force attach of
  // entities while migrating from one cluster to another. For more details
  // please look into 'insights_interface.proto'.
  KForceAttachIncarnationID uint64 = (1<<64 -1)
)

var (
  DefaultInsightAddr = flag.String(
    "insights_default_ip",
    "127.0.0.1",
    "Default insights IP.")

  DefaultInsightPort = flag.Uint(
    "insights_default_port",
    2027,
    "Default insights port.")

  retryInitialDelaySecs = flag.Int64(
    "insights_retry_initial_delay_secs",
    1,
    "Initial wait time before retrying a failed insights request.")

  retryMaxDelaySecs = flag.Int64(
    "insights_retry_max_delay_secs",
    30,
    "Max wait time between failed insights request retries.")

  maxRetries = flag.Int(
    "insights_max_retries",
    10,
    "Number of retries for a failed insights request. It does not include "+
      "the first call.")

  InsightsPort = flag.Uint(
    "insights_port",
    *DefaultInsightPort,
    "Port number of Insights Service to connect.")

  InsightsServerIP = flag.String(
    "insights_server_ip",
    *DefaultInsightAddr,
    "IP address of Insights Service to connect.")
)

// Errors generated out of this module
var (
  ErrMissingSupportInsights = InsightsError("Missing support in Insights",
    3001)
)

type InsightsServiceInterface interface {
  SetRequestTimeout(timeout int64) error
  SendMsg(service string, request,
    response proto.Message, backoff *util_misc.ExponentialBackoff) error
  SendMsgWithTimeout(service string,
    request, response proto.Message, backoff *util_misc.ExponentialBackoff,
    timeoutSecs int64) error
  IsRetryError(err error) (bool, error)
  GetMasterLocation(
    backoff *util_misc.ExponentialBackoff) (string, error)
}

func DefaultInsightsServiceInterface() InsightsServiceInterface {
  return DefaultInsightsService()
}

func NewInsightsServiceInterface(
  serverIp string, serverPort uint16) InsightsServiceInterface {
  return NewInsightsService(serverIp, serverPort)
}

type InsightsService struct {
  client *util_net.ProtobufRPCClient
}

func DefaultInsightsService() *InsightsService {
  return NewInsightsService(*DefaultInsightAddr, uint16(*DefaultInsightPort))
}

func NewInsightsServiceFromGflags() *InsightsService {
  return NewInsightsService(*InsightsServerIP, uint16(*InsightsPort))
}

func NewInsightsService(serverIp string, serverPort uint16) *InsightsService {
  glog.Infoln("Initializing Insights Service")
  return &InsightsService{
    client: util_net.NewProtobufRPCClient(serverIp, serverPort),
  }
}

func (insights_svc *InsightsService) SetRequestTimeout(timeout int64) error {
  if insights_svc.client == nil {
    return ErrInternalError
  }
  insights_svc.client.SetRequestTimeout(timeout)
  return nil
}

func (insights_svc *InsightsService) SendMsg(service string, request,
  response proto.Message, backoff *util_misc.ExponentialBackoff) error {

  return insights_svc.SendMsgWithTimeout(service, request, response, backoff,
                                         0 /* timeoutSecs */)
}

func (insights_svc *InsightsService) SendMsgWithTimeout(service string,
  request, response proto.Message, backoff *util_misc.ExponentialBackoff,
  timeoutSecs int64) error {

  // Use insights_svc default backoff if client doesn't supply backoff
  // mechanism for this request.
  if backoff == nil {
    backoff = util_misc.NewExponentialBackoff(
      time.Duration(*retryInitialDelaySecs)*time.Second,
      time.Duration(*retryMaxDelaySecs)*time.Second,
      *maxRetries)
  }

  for {
    err := insights_svc.client.CallMethodSync(insightServiceName, service,
      request, response, timeoutSecs)

    if err == nil {
      return nil
    }

    glog.Error("Insights RPC Error: ", err)

    isRetry, newErr := insights_svc.IsRetryError(err)
    if !isRetry {
      return newErr
    }

    // Retry attempt
    waited := backoff.Backoff()
    if waited == util_misc.Stop {
      glog.Error(
        "Done with maximum retries, Giving up connecting to insights_svc.")
      return err
    }
    glog.Infof("Failed to send msg. Retry msg %s", service)
  }
}

// Check if the error is retry error. Return true if it is and false
// otherwise. Also return the new error.
func (insights_svc *InsightsService) IsRetryError(err error) (bool, error) {

  // Focus on App error and RPC error, only certain App errors and transport
  // error in RPC error need retry.
  if obj, ok := ntnx_errors.TypeAssert(err, ntnx_errors.AppErrorType); ok {
    appErr, ok := obj.(*util_net.AppError_)
    if !ok {
      return false, ErrInternalError
    }
    errCode := appErr.GetErrorCode()
    errRet := Errors[errCode]
    insightsErr, ok := errRet.(*InsightsInterfaceError_)
    if !ok || !insightsErr.IsRetryError() {
      return false, errRet
    }
  } else if obj, ok := ntnx_errors.TypeAssert(err, ntnx_errors.RpcErrorType); ok {
    rpcErr, ok := obj.(*util_net.RpcError_)
    if !ok || !rpcErr.IsRetryError() {
      return false, err
    }
  } else {
    return false, err
  }

  return true, err
}

func (e *InsightsInterfaceError_) IsRetryError() bool {
  if e == ErrRetry || e == ErrTransportError || e == ErrUnavailable ||
    e == ErrTimeout || (2000 <= e.GetErrorCode() && e.GetErrorCode() <= 2999) {
    return true
  }
  return false
}

func (insights_svc *InsightsService) GetMasterLocation(
  backoff *util_misc.ExponentialBackoff) (string, error) {
  arg := &GetMasterLocationArg{}
  ret := &GetMasterLocationRet{}

  err := insights_svc.SendMsg("GetMasterLocation", arg, ret, backoff)
  if err != nil {
    return "", err
  }
  return ret.GetMasterHandle(), nil
}
