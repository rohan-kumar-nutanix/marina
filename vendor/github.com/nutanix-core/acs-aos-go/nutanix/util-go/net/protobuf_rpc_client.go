//
// Copyright (c) 2016 Nutanix Inc. All rights reserved.

// Author: akshay@nutanix.com

// Package net provides various network related functionalities, including
// Protobuf RPC client, server and JSON RPC client, server.
package net

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	proto_net "github.com/nutanix-core/acs-aos-go/nutanix/util-slbufs/util/sl_bufs/net"
)

// Interface to hold the RPC Metrics callback functions
type RPCClientMetrics interface {
	BeforeSendingRequest(requestId uint64, method string, svcName string)
	AfterReceivingResponse(requestId uint64, method string, svcName string, error bool)
}

// Struct to hold all client side details.
type ProtobufRPCClient struct {
	HttpClient *http.Client
	Id         int32
	serverIp   string
	ServerPort uint16
	RpcMetrics RPCClientMetrics
	requestId  uint64
}

const (
	RPCHeaderSizeBits = 4
	DefaultRPCTimeout = 10
)

type ProtobufRPCClientIfc interface {
	CallMethodSync(serviceName, methodName string,
		requestProto, responseProto proto.Message, timeoutSecs int64) error

	SetRequestTimeout(timeoutSecs int64) error
}

func NewProtobufRPCClientIfc(
	serverIp string, serverPort uint16) ProtobufRPCClientIfc {

	return NewProtobufRPCClient(serverIp, serverPort)
}

// This method should be used to instantiate new protobuf rpc client.
func NewProtobufRPCClient(
	serverIp string, serverPort uint16) *ProtobufRPCClient {

	return &ProtobufRPCClient{
		HttpClient: &http.Client{},
		Id:         rand.Int31(),
		serverIp:   serverIp,
		ServerPort: serverPort,
		RpcMetrics: nil,
		requestId:  0,
	}
}

// Use this method to enable RPC NuSights Metrics.
func (client *ProtobufRPCClient) RegisterRPCMetricsHandlers(metrics RPCClientMetrics) {
	client.RpcMetrics = metrics
}

func (client *ProtobufRPCClient) SetRequestTimeout(timeoutSecs int64) error {
	if client.HttpClient == nil {
		return errors.New("Internal error: nil HttpClient")
	}
	client.HttpClient.Timeout = time.Second * time.Duration(timeoutSecs)
	return nil
}

// This is function signature of callback function to be registered with RPCs.
// It is executed when RPC is done.
// This signature is not final yet, it will be revised.
// Args:
//	rpc : Pointer to ProtobufRpc (RPC control data structure)
//	req : Request protobuf for RPC method.
//	res : Response protobuf for RPC method.
type RPCDoneMethodFn func(
	rpc *ProtobufRpc, req proto.Message, res proto.Message)

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
// 	 RPC status and error details
func (client *ProtobufRPCClient) CallMethodSyncWithPayload(
	serviceName, methodName string, requestProto,
	responseProto proto.Message, requestPayload []byte, responsePayload *[]byte,
	timeoutSecs int64) error {

	rpc := NewProtobufRPC()
	rpc.RequestPayload = requestPayload
	rpc.TimeoutSecs = timeoutSecs

	go client.callMethod(
		serviceName, methodName, rpc, requestProto, responseProto, false)

	<-rpc.Done
	// Populate the response payload appropriately.
	if responsePayload != nil {
		*responsePayload = rpc.ResponsePayload
	}
	return rpc.RpcError
}

// This is wrapper around the `CallMethod` function with the addition of
// the response and request payloads with explicit connection close flag set.
// Args:
//        serviceName  : Name of service.
//        methodName   : Method to call on server belonging to service.
//        requestProto : Input Protobuf object for method.
//        responseProto: Return Protobuf object for method.
//        requestPayload: Payload to send with the RPC request.
//        responsePayload: Payload received with the RPC response.
// Return:
// 	 RPC status and error details
func (client *ProtobufRPCClient) CallMethodSyncWithPayloadAndConnClosed(
	serviceName, methodName string, requestProto,
	responseProto proto.Message, requestPayload []byte, responsePayload *[]byte,
	timeoutSecs int64) error {

	rpc := NewProtobufRPC()
	rpc.RequestPayload = requestPayload
	rpc.TimeoutSecs = timeoutSecs

	go client.callMethod(
		serviceName, methodName, rpc, requestProto, responseProto, true)

	<-rpc.Done
	// Populate the response payload appropriately.
	if responsePayload != nil {
		*responsePayload = rpc.ResponsePayload
	}
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
func (client *ProtobufRPCClient) CallMethodSync(serviceName, methodName string,
	requestProto, responseProto proto.Message, timeoutSecs int64) error {

	return client.CallMethodSyncWithPayload(serviceName, methodName,
		requestProto, responseProto, nil, nil, timeoutSecs)
}

// This method collects and sends RPC Server Request Response Metrics.
// Args:
//      requestId   : Client Request sent Identifier
//      method      : RPC Service Method Requested
//      svcName     : Name of service.
//      err         : RPC Service Method Error Status
func (client *ProtobufRPCClient) rpcMetricsResponse(requestId uint64, method string, svcName string, err error) {
	var respErr bool
	if err != nil {
		respErr = true
	} else {
		respErr = false
	}
	client.RpcMetrics.AfterReceivingResponse(requestId, method, svcName, respErr)
}

// This method constructs HTTP post data request and sends to RPC server
// with explicit connection close flag set.
// Args:
//        serviceName  : Name of service.
//        methodName   : Method to call on server belonging to service.
//        requestProto : Input Protobuf object for method.
//        responseProto: Return Protobuf object for method.
//        timeoutSecs  : timeout for each request.
//        closeConnection : Flag to explicitly close the RPC connection.
// Return:
//        error: nil if successful, error otherwise.
func (client *ProtobufRPCClient) callMethod(
	serviceName string, methodName string, rpc *ProtobufRpc,
	requestProto proto.Message, responseProto proto.Message, closeConnection bool) {

	var requestId uint64

	// Fill out request header.
	rpcHeader := &rpc.RequestHeader
	rpcId := rand.Int63()
	rpcHeader.RpcId = &rpcId
	rpcHeader.MethodName = &methodName
	epochTime := time.Now().UnixNano() / int64(time.Millisecond)
	rpcHeader.SendTimeUsecs = &epochTime
	httpClient := client.HttpClient

	if rpc.TimeoutSecs == 0 {
		rpc.TimeoutSecs = DefaultRPCTimeout
		if httpClient.Timeout > 0 {
			// If the caller has set a timeout on HTTP, use it.
			rpc.TimeoutSecs = int64(httpClient.Timeout / time.Second)
		}
	}

	deadlineUsecs := epochTime + (rpc.TimeoutSecs * int64(time.Microsecond))
	rpcHeader.DeadlineUsecs = &deadlineUsecs

	var rpc_error error = nil
	defer func() { RpcFinish(rpc, rpc_error) }()

	// Serialize RPC request.
	rpcRequest, err := client.SerializeRPCRequest(
		rpc.RequestHeader, requestProto, rpc.RequestPayload)

	if err != nil {
		rpc_error = ErrSendingRpc.SetCause(err)
		return
	}

	// Make RPC request.
	bodyReader := bytes.NewReader(rpcRequest)

	request, err := http.NewRequest(
		"POST", client.getserviceUrl(serviceName), bodyReader)

	if err != nil {
		rpc_error = ErrSendingRpc.SetCause(err)
		return
	}

	// Record RPC metrics before sending request, If option is set.
	if client.RpcMetrics != nil {
		requestId = atomic.AddUint64(&client.requestId, 1)
		client.RpcMetrics.BeforeSendingRequest(requestId, methodName, serviceName)
		defer func() { client.rpcMetricsResponse(requestId, methodName, serviceName, err) }()
	}

	// This flag sets the connection close flag to true.
	// TODO: Set this flag to true by default in future (couple of months down the lane to have
	// a proper soak time)
	if closeConnection {
		request.Close = true
	}
	request.Header.Add("Content-Type", "application/x-rpc")
	httpClient.Timeout = time.Duration(rpc.TimeoutSecs) * time.Second
	response, err := httpClient.Do(request)
	if err != nil {
		// The request body, if non-nil, will be closed by the
		// underlying transport, even on errors.
		glog.Errorf("RPC request error: %s type: %s", err, reflect.TypeOf(err))

		now := time.Now().UnixNano() / int64(time.Millisecond)
		glog.Infof("Now: %d -- epoch time: %d", now, epochTime)
		if (now - epochTime) >= ((rpc.TimeoutSecs) * 1000) {
			rpc_error = ErrTimeout.SetCause(err)
			return
		}
		rpc_error = ErrRpcTransport.SetCause(err)
		return
	}

	// Get data.
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		rpc_error = ErrRpcTransport.SetCause(err)
		return
	}

	err = client.DeserializeRPCRequest(rpc, data, responseProto)
	if err != nil {
		glog.Error("Failed to deserialize RPC with error: ", err)
		rpc_error = err
	}
}

// Prepare URL for service.
func (client *ProtobufRPCClient) getserviceUrl(serviceName string) string {
	return fmt.Sprintf("http://%v:%v/rpc/%v/",
		client.serverIp, client.ServerPort, serviceName)
}

// Serialize protobuf for rpc request, add headers and
// make data ready for request.
// Args:
//        rpcHeader    : Header protobuf of RPC.
//        requestProto : Input protobuf object for method.
// Return:
//        byte[] : serialized data ready for RPC.
//        error  : nil if successful, error otherwise.
//
//    RPC requests are encoded on the wire as follows:
//    (1) First four bytes indicating the size of the RPC header in network byte
//        order.
//    (2) RPC header for the request/response.
//    (3) RPC request protobuf.
//    (4) RPC request payload.
//    The RPC header itself contains the size of (3) and (4).
func (client *ProtobufRPCClient) SerializeRPCRequest(
	rpcHeader proto_net.RpcRequestHeader, requestProto proto.Message,
	requestPayload []byte) ([]byte, error) {

	// Serialize request protobuf object.
	var err error
	requestBuf, err := proto.Marshal(requestProto)
	if err != nil {
		glog.Error("Protobuf marshaling failed with err: ", err)
		return nil, err
	}

	rpcHeader.ProtobufSize = proto.Int32(int32(len(requestBuf)))
	if requestPayload != nil {
		rpcHeader.PayloadSize = proto.Int32(int32(len(requestPayload)))
	}

	// Serialize header.
	requestHeaderBuf, err := proto.Marshal(&rpcHeader)
	if err != nil {
		glog.Error("Protobuf marshaling failed with err: ", err)
		return nil, err
	}

	headerLenBuf := make([]byte, RPCHeaderSizeBits)
	binary.BigEndian.PutUint32(headerLenBuf, uint32(len(requestHeaderBuf)))

	// Combine header length, header, actual request and the payload.
	rpcRequest := bytes.Join([][]byte{
		headerLenBuf, requestHeaderBuf, requestBuf, requestPayload},
		nil)

	return rpcRequest, nil
}

//  Interpret data sent by server in response.
//  Args:
//        data          : Complete RPC response to be deserialized.
//        responseProto : Protobuf filled with deserialized data.
//  Return:
//        error : nil if successful, error otherwise.
//
//    RPC responses are encoded on the wire as follows:
//    (1) First four bytes indicating the size of the RPC header in network byte
//        order.
//    (2) RPC header for the response.
//    (3) RPC response protobuf.
//    (4) RPC response payload.
//    The RPC header itself contains the size of (3) and (4).
func (client *ProtobufRPCClient) DeserializeRPCRequest(rpc *ProtobufRpc,
	data []byte, responseProto proto.Message) error {

	if len(data) < RPCHeaderSizeBits {
		return ErrInvalidRpcResponse.SetCause(errors.New(
			"received buffer too short"))
	}

	headerSizeBuf := data[:RPCHeaderSizeBits]
	headerSize := int(binary.BigEndian.Uint32(headerSizeBuf))
	if headerSize > (len(data) + RPCHeaderSizeBits) {
		return ErrInvalidRpcResponse.SetCause(errors.New("header size is > len(data)"))
	}
	headerBuf := data[RPCHeaderSizeBits : RPCHeaderSizeBits+headerSize]

	// Deserialize response header.
	var errorDetail string

	err := proto.Unmarshal(headerBuf, &rpc.ResponseHeader)
	if err != nil {
		return ErrInvalidRpcResponse.SetCause(errors.New(err.Error()))
	}

	errorDetail = rpc.ResponseHeader.GetErrorDetail()

	switch rpc.ResponseHeader.GetRpcStatus() {

	case proto_net.RpcResponseHeader_kNoError:
		break

	case proto_net.RpcResponseHeader_kMethodError:
		return ErrMethod.SetCause(errors.New(errorDetail))

	case proto_net.RpcResponseHeader_kAppError:
		return AppError(errorDetail,
			int(rpc.ResponseHeader.GetAppError()))

	case proto_net.RpcResponseHeader_kCanceled:
		return ErrCanceled.SetCause(errors.New(errorDetail))

	case proto_net.RpcResponseHeader_kTimeout:
		return ErrTimeout.SetCause(errors.New(errorDetail))

	case proto_net.RpcResponseHeader_kTransportError:
		return ErrTransport.SetCause(errors.New(errorDetail))

	default:
		return ErrInvalidRpcResponse.SetCause(
			errors.New("Unknown RPC status returned. " + string(errorDetail)))
	}

	responseBufSize := int(rpc.ResponseHeader.GetProtobufSize())
	if responseBufSize > 1 {
		if len(data) < RPCHeaderSizeBits+headerSize+responseBufSize {
			errorDetail = "size of protobuf and payload do not " +
				"correspond to the size of the received buffer."
			return ErrInvalidRpcResponse.SetCause(errors.New(errorDetail))
		}

		start := RPCHeaderSizeBits + headerSize
		end := RPCHeaderSizeBits + int(headerSize) + int(responseBufSize)
		responseBuf := data[start:end]

		// Deserialize RPC response.
		err = proto.Unmarshal(responseBuf, responseProto)
		if err != nil {
			errorDetail = "unable to parse RPC response"
			return ErrInvalidRpcResponse.SetCause(errors.New(errorDetail))
		}
	}

	payloadSize := int(rpc.ResponseHeader.GetPayloadSize())
	if payloadSize > 0 {
		if len(data) < RPCHeaderSizeBits+headerSize+responseBufSize+payloadSize {
			errorDetail = "size of protobuf and payload do not " +
				"correspond to the size of the received buffer."
			return ErrInvalidRpcResponse.SetCause(errors.New(errorDetail))
		}
		start := RPCHeaderSizeBits + headerSize + responseBufSize
		end := RPCHeaderSizeBits + headerSize + responseBufSize + payloadSize
		rpc.ResponsePayload = make([]byte, payloadSize)
		copy(rpc.ResponsePayload, data[start:end])
	}

	return nil
}

// This method sets correct status and error details in rpc and sends
// response to client accordingly.
// It also sets rpc as done.
func RpcFinish(rpc *ProtobufRpc, rpc_error error) {
	if rpc_error != nil {
		rpc.RpcError = rpc_error
	}
	rpc.Done <- true
}
