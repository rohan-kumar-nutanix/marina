//
// Copyright (c) 2017 Nutanix Inc. All rights reserved.
//
// Author: sumanth.suresh@nutanix.com
//
// This provides functionalities for a remote rpc client.
package aplos_transport

import (
        "bytes"
        "errors"
        util_net "github.com/nutanix-core/acs-aos-go/nutanix/util-go/net"
        "github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
        "io/ioutil"
        "math/rand"
        "net/http"
        "net/url"
        net "github.com/nutanix-core/acs-aos-go/nutanix/util-slbufs/util/sl_bufs/net"
        "reflect"
        "strconv"
        "time"
        "github.com/nutanix-core/acs-aos-go/zeus"

        "github.com/golang/glog"
        "github.com/golang/protobuf/proto"
)

// Struct to hold all client side details.
type RemoteProtobufRPCClient struct {
        *util_net.ProtobufRPCClient
        serverName string
}

const (
        RPCHeaderSizeBits       = 4
        DefaultRemoteRPCTimeout = 40
)

func NewRemoteProtobufRPCClientIfc(
        zkSession *zeus.ZookeeperSession, serverName string,
        serverPort uint16, clusterUuid *uuid4.Uuid, userUuid *uuid4.Uuid,
        tenantUuid *uuid4.Uuid, timeout int, userName *string,
        reqContext *net.RpcRequestContext) util_net.ProtobufRPCClientIfc {
        client, err := NewRemoteProtobufRPCClient(zkSession, serverName, serverPort,
                clusterUuid, userUuid, tenantUuid, timeout, userName, reqContext)
        if err != nil {
                glog.Error("Failed to create a remote RPC client with error: ", err)
                return nil
        }
        return client
}

// This method should be used to instantiate new protobuf rpc client.
func NewRemoteProtobufRPCClient(zkSession *zeus.ZookeeperSession,
        serverName string, serverPort uint16, clusterUuid *uuid4.Uuid,
        userUuid *uuid4.Uuid, tenantUuid *uuid4.Uuid,
        timeout int, userName *string,
        reqContext *net.RpcRequestContext) (*RemoteProtobufRPCClient, error) {
        transport, err := NewRcHttpTransport(zkSession,
                clusterUuid, userUuid, tenantUuid, DefaultRemoteRPCTimeout, userName,
                reqContext)
        if err != nil {
                return nil, err
        }
        return &RemoteProtobufRPCClient{
                ProtobufRPCClient: &util_net.ProtobufRPCClient{
                        HttpClient: &http.Client{Transport: transport},
                        Id:         rand.Int31(),
                        ServerPort: serverPort,
                },
                serverName: serverName,
        }, nil
}

// This method constructs HTTP post data request and sends to RPC server.
// Args:
//        serviceName  : Name of service.
//        methodName   : Method to call on server belonging to service.
//        requestProto : Input Protobuf object for method.
//        responseProto: Return Protobuf object for method.
//        timeoutSecs  : timeout for each request.
// Return:
//        error: nil if successful, error otherwise.
func (client *RemoteProtobufRPCClient) callMethod(
        serviceName string, methodName string, rpc *util_net.ProtobufRpc,
        requestProto proto.Message, responseProto proto.Message) {

        // Fill out request header.
        rpcHeader := &rpc.RequestHeader
        rpcId := rand.Int63()
        rpcHeader.RpcId = &rpcId
        rpcHeader.MethodName = &methodName
        epochTime := time.Now().UnixNano() / int64(time.Millisecond)
        rpcHeader.SendTimeUsecs = &epochTime
        httpClient := client.HttpClient

        if rpc.TimeoutSecs == 0 {
                rpc.TimeoutSecs = DefaultRemoteRPCTimeout
                if httpClient.Timeout > 0 {
                        // If the caller has set a timeout on HTTP, use it.
                        rpc.TimeoutSecs = int64(httpClient.Timeout / time.Second)
                }
        }

        deadlineUsecs := epochTime + (rpc.TimeoutSecs * int64(time.Microsecond))
        rpcHeader.DeadlineUsecs = &deadlineUsecs

        var rpc_error error = nil
        defer func() { util_net.RpcFinish(rpc, rpc_error) }()

        // Serialize RPC request.
        rpcRequest, err := client.SerializeRPCRequest(
                rpc.RequestHeader, requestProto, rpc.RequestPayload)

        if err != nil {
                rpc_error = util_net.ErrSendingRpc.SetCause(err)
                return
        }

        // Make rpc request.
        bodyReader := bytes.NewReader(rpcRequest)

        request, err := http.NewRequest(
                "POST", client.getUrlParameterList(serviceName, client.ServerPort,
                        rpc.TimeoutSecs*1000), bodyReader)

        if err != nil {
                rpc_error = util_net.ErrSendingRpc.SetCause(err)
                return
        }

        request.Header.Add("Content-Type", "application/octet-stream")
        httpClient.Timeout = time.Duration(rpc.TimeoutSecs) * time.Second
        response, err := httpClient.Do(request)
        if err != nil {
                // The request Body, if non-nil, will be closed by the
                // underlying Transport, even on errors.
                glog.Errorf("RPC request error: %s type: %s", err, reflect.TypeOf(err))

                now := time.Now().UnixNano() / int64(time.Millisecond)
                glog.Infof("Now: %d -- epoch time: %d", now, epochTime)
                if (now - epochTime) >= ((rpc.TimeoutSecs) * 1000) {
                        rpc_error = util_net.ErrTimeout.SetCause(err)
                        return
                }
                rpc_error = util_net.ErrRpcTransport.SetCause(err)
                return
        }

        // Get data.
        defer response.Body.Close()
        data, err := ioutil.ReadAll(response.Body)
        if err != nil {
                rpc_error = util_net.ErrRpcTransport.SetCause(err)
                return
        } else if len(data) == 0 {
                rpc_error = util_net.ErrRpcTransport.SetCause(
                        errors.New("Data not found in RPC response."))
                return
        }

        err = client.DeserializeRPCRequest(rpc, data, responseProto)
        if err != nil {
                glog.Error("Failed to deserialize RPC with error: ", err)
                rpc_error = err
        }
}

// This is wrapper around the `CallMethodSync` function with the addition of
// the response and request payloads.
// Args:
//        serviceName  : Name of service.
//        methodName   : Method to call on server belonging to service.
//        requestProto : Input Protobuf object for method.
//        responseProto: Return Protobuf object for method.
//        requestPayload: Payload to send with the RPC request.
//        responsePayload: Payload received with the RPC response.
// Return:
//       RPC status and error details
func (client *RemoteProtobufRPCClient) CallMethodSyncWithPayload(
        serviceName, methodName string, requestProto,
        responseProto proto.Message, requestPayload, responsePayload []byte,
        timeoutSecs int64) error {

        rpc := util_net.NewProtobufRPC()
        rpc.RequestPayload = requestPayload
        rpc.ResponsePayload = responsePayload
        rpc.TimeoutSecs = timeoutSecs

        go client.callMethod(
                serviceName, methodName, rpc, requestProto, responseProto)

        <-rpc.Done
        return rpc.RpcError
}

// This method constructs HTTP post data request and sends to RPC server.
// The method is executed in goroutine,
// client has to wait on rpc.Done channel if synchronous operation is desired.
// Args:
//        serviceName  : Name of service.
//        methodName   : Method to call on server belonging to service.
//        requestProto : Input Protobuf object for method.
//        responseProto: Return Protobuf object for method.
// Return:
//        RPC status and error details.
func (client *RemoteProtobufRPCClient) CallMethodSync(serviceName, methodName string,
        requestProto, responseProto proto.Message, timeoutSecs int64) error {

        return client.CallMethodSyncWithPayload(serviceName, methodName,
                requestProto, responseProto, nil, nil, timeoutSecs)
}

// Prepare URL for service.
func (client *RemoteProtobufRPCClient) getUrlParameterList(serviceName string, servicePort uint16,
        timeoutMs int64) string {
        Url, _ := url.Parse("")
        parameters := url.Values{
                "service_name": {serviceName},
                "port":         {strconv.Itoa(int(servicePort))},
                "timeout_ms":   {strconv.FormatInt(timeoutMs, 10)},
                "base_url":     {"/rpc"},
        }
        Url.RawQuery = parameters.Encode()
        urlParameterList := "remote_rpc_request" + Url.String()
        glog.Errorf("Remote rpc url parameters: %s", urlParameterList)
        return urlParameterList
}
