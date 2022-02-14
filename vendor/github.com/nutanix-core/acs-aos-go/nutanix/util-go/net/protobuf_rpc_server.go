// topyright (c) 2016 Nutanix Inc. All rights reserved.
//
// Author: akshay@nutanix.com

//Protobuf RPC server implementation.
package net

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"sync/atomic"

	glog "github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	proto_net "github.com/nutanix-core/acs-aos-go/nutanix/util-slbufs/util/sl_bufs/net"
)

var (
	disabledChunkedEncodingForResponses = flag.Bool("disable_http_response_chunked_encoding",
		true,
		"If true, disables chunked encoding while sending responses.")
)

// Interface to hold the RPC Metrics callback functions
type RPCServerMetrics interface {
	RequestReceived(requestId uint64, method string)
	SendingResponse(requestId uint64, method string, error bool)
}

// Struct to hold server side details.
type ProtobufRPCHTTPServer struct {
	HttpServer HTTPServer
	Services   map[string]*Service
	RpcMetrics RPCServerMetrics
	requestId  uint64
}

// This is function signature for RPC wrapper function to be registered by
// service exporting RPCs.
type ServiceMethodFn func(rpc *ProtobufRpc, handlers interface{}) error

// Description of Service.
// Name : Name of service, it should be used in url for rpc.
// Methods: Map of Method name to function implementing it.
type ServiceDesc struct {
	Name    string
	Methods map[string]ServiceMethodFn
}

// Holds all data required to call methods in a service.
type Service struct {
	Desc *ServiceDesc
	Impl interface{}
}

// RPC header size is 4 bytes.
const CRpcReqHeaderSize = 4

// Determine service name from URL path.
func (server *ProtobufRPCHTTPServer) getServiceName(path string) string {
	fullName := strings.Trim(path, "/")
	fullName = strings.ToLower(fullName)
	return strings.TrimLeft(fullName, "rpc/")
}

// Prepare Server response.
// Args:
//     requestHeader  : Standardized request header for protobuf rpcs.
//     responseHeader : Standardized response header for protobuf rpcs.
//     responseBuf    : Response buffer.
//
// Return:
//    Byte array : RPC response.
//    error      : Error occurred during execution, nil if succcessful.
func (server *ProtobufRPCHTTPServer) prepareRpcResponse(
	requestHeader proto_net.RpcRequestHeader,
	responseHeader proto_net.RpcResponseHeader,
	responseBuf []byte) ([]byte, error) {

	responseHeader.RpcId = proto.Int64(requestHeader.GetRpcId())
	// RequestBuf contains the marshalled result from the call.
	responseHeader.ProtobufSize = proto.Int32(int32(len(responseBuf)))
	responseHeaderBuf, err := proto.Marshal(&responseHeader)

	if err != nil {
		glog.Error("Protobuf marshaling failed: ", err)
		return nil, err
	}

	headerLenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(headerLenBuf, uint32(len(responseHeaderBuf)))

	return bytes.Join(
		[][]byte{headerLenBuf, responseHeaderBuf, responseBuf}, nil), nil
}

// Prepare Server error response.
// Args:
//     rpc            : Protobuf rpc object.
//     status         : RPC status
//     errorDetail    : Details of RPC error.
//
// Return:
//    Byte array : RPC response.
//    error      : Error occurred during execution, nil if succcessful.
func (server *ProtobufRPCHTTPServer) prepareRpcErrorResponse(
	rpc *ProtobufRpc, status proto_net.RpcResponseHeader_RpcStatus,
	errorDetail string) ([]byte, error) {
	rpc.ResponseHeader.RpcId = proto.Int64(rpc.RequestHeader.GetRpcId())
	rpc.ResponseHeader.RpcStatus = status.Enum()
	if errorDetail != "" {
		rpc.ResponseHeader.ErrorDetail = proto.String(errorDetail)
	}
	rpc.Status = uint32(status)

	rpc.ResponseHeader.ProtobufSize = proto.Int32(
		int32(len(rpc.ResponsePayload)))
	responseHeaderBuf, err := proto.Marshal(&rpc.ResponseHeader)
	if err != nil {
		glog.Error("Protobuf marshaling failed: ", err)
		return nil, err
	}

	headerLenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(headerLenBuf, uint32(len(responseHeaderBuf)))

	return bytes.Join(
		[][]byte{headerLenBuf, responseHeaderBuf,
			rpc.ResponsePayload}, nil), nil
}

// This method collects and sends RPC Response Metrics.
// Args:
//      requestId   : Server Request received Identifier
//      method      : RPC Service Method Requested
//      err         : RPC Service Method Error Status
func (server *ProtobufRPCHTTPServer) rpcMetricsResponse(requestId uint64, method string, err error) {
	var respErr bool
	if err != nil {
		respErr = true
	} else {
		respErr = false
	}
	server.RpcMetrics.SendingResponse(requestId, method, respErr)
}

// This method handles a single RPC request.
// Arg:
//    rpcSvc : Service of which RPC has to be served.
//    req : HTTP request.
//    res : HTTP response.
//
// Return:
//       error : Error occurred during execution, nil if successful.
func (server *ProtobufRPCHTTPServer) RPCRequestHandler(
	rpcSvc string, req *http.Request, res http.ResponseWriter) error {

	var errorDetail string
	var requestId uint64

	rpcService, ok := server.Services[rpcSvc]
	if !ok {
		errorDetail = fmt.Sprintf("Unknown service: '%s'", rpcSvc)
		glog.Error("RPC request handler: ", errorDetail)
		return errors.New(errorDetail)
	}

	rpc := new(ProtobufRpc)
	hdrSize := req.ContentLength

	rcvIobuf := make([]byte, hdrSize)
	defer req.Body.Close()
	_, err := req.Body.Read(rcvIobuf)
	if err != nil && err.Error() != "EOF" {
		errorDetail = "Error reading http_req.Body"
		glog.Error("Error in RPC request handler: ", errorDetail)
		return errors.New(errorDetail)
	}

	rpcbuflen := uint32(len(rcvIobuf))
	// No request in the body.
	if rpcbuflen < CRpcReqHeaderSize {
		errorDetail = "Invalid RPC request"
		glog.Error("Error in RPC request handler: ", errorDetail)
		return errors.New(errorDetail)
	}

	// Read RPC header size from first 4 bytes.
	headerLenBuf := rcvIobuf[:CRpcReqHeaderSize]
	rpcHeaderLen := binary.BigEndian.Uint32(headerLenBuf)

	if rpcHeaderLen < 0 || rpcHeaderLen+CRpcReqHeaderSize > rpcbuflen {
		errorDetail = fmt.Sprintf("Bad response header size (%v/%v)",
			rpcHeaderLen, rpcbuflen)
		glog.Error("Error in RPC request handler: ", errorDetail)
		return errors.New(errorDetail)
	}

	// Now read the RPCRequest Header into the struct
	// created from rpc.proto.
	var requestHeaderBuf []byte
	requestHeaderBuf =
		rcvIobuf[CRpcReqHeaderSize : CRpcReqHeaderSize+rpcHeaderLen]

	err = proto.Unmarshal(requestHeaderBuf, &rpc.RequestHeader)
	if err != nil {
		glog.Error("Unable to parse request header: ", err)
		return err
	}

	// Get RequestPayload from buffer.
	rpc.RequestPayload = rcvIobuf[CRpcReqHeaderSize+rpcHeaderLen:]

	// Lookup the method being invoked.
	method := rpc.RequestHeader.GetMethodName()

	if rpcService.Desc == nil {
		return errors.New("Service not supported")
	}

	if method == "" || rpcService.Desc.Methods[method] == nil {
		errorDetail = fmt.Sprintf("Method %v doesn't exist", method)

		data, err := server.prepareRpcErrorResponse(
			rpc, proto_net.RpcResponseHeader_kMethodError,
			errorDetail)
		if err != nil {
			glog.Errorf("Failed to prepare RPC error response %s", err)
		}
		res.Write(data)
		glog.Error(errorDetail)
		return errors.New(errorDetail)
	}

	// Call RPC Metrics Request Received, If Metrics enabled
	if server.RpcMetrics != nil {
		requestId = atomic.AddUint64(&server.requestId, 1)
		server.RpcMetrics.RequestReceived(requestId, method)
		defer func() { server.rpcMetricsResponse(requestId, method, err) }()
	}

	// Call the service from service map.
	err = rpcService.Desc.Methods[method](rpc, rpcService.Impl)
	if err != nil {
		data, err := server.prepareRpcErrorResponse(rpc,
			2, fmt.Sprint(err))
		if err != nil {
			glog.Error("Failed to prepare RPC response: ", err)
			return err
		}
		res.Write(data)
		return err
	}

	// Write response.
	data, err := server.prepareRpcResponse(
		rpc.RequestHeader, rpc.ResponseHeader, rpc.ResponsePayload)
	if err != nil {
		glog.Error("Failed to prepare RPC response: ", err)
		return err
	}

	// Golang http server sets the transfer-encoding to chunked (by default) if content-length
	// is not set explicitly and if we are using HTTP/1.1.
	// C++/Java based Nutanix RPC clients do not understand chunked encoded responses, so we
	// should disable chunked encoding by adding a Content-Length header.
	if *disabledChunkedEncodingForResponses {
		res.Header().Set("Content-Length", strconv.Itoa(len(data)))
	}

	// If the client sent a "Reorder-Request" header, we should just echo it back since we
	// do not perform any reordering ourselves. However, some clients (C++/Java) may use this
	// header from the response to manage their internal reordering state.
	reorder_hdr := req.Header.Get("Reorder-Request")
	if reorder_hdr != "" {
		res.Header().Set("Reorder-Request", reorder_hdr)
	}

	res.Write(data)

	return nil
}

// Create HTTP RPC handler method.
// Args:
//      Server: Instance of HTTP server.
// Return:
//      HTTP handler method.
func HandleProtobufRPC(server *ProtobufRPCHTTPServer) http.Handler {
	fn := func(res http.ResponseWriter, req *http.Request) {
		_, err := httputil.DumpRequest(req, true)
		if err != nil {
			http.Error(
				res, fmt.Sprint(err),
				http.StatusInternalServerError)
			return
		}

		if req.Method != "POST" {
			glog.Error("Only post methods are accepted")
			return
		}

		rpcSvc := server.getServiceName(req.URL.Path)
		server.RPCRequestHandler(rpcSvc, req, res)
	}
	return http.HandlerFunc(fn)
}

// Use this method to instantiate new protobuf rpc http server.
// Arg:
//    Server: Instance of HTTP server.
// Return:
//       Instance of Protobuf RPC HTTPServer
func NewProtobufRPCHTTPServer(Server *HTTPServer) *ProtobufRPCHTTPServer {
	glog.Infof("Initializing protobuf RPC HTTP server with %s:%s",
		Server.ServerAddress, Server.ServerPort)
	server := &ProtobufRPCHTTPServer{
		HttpServer: *Server,
		Services:   make(map[string]*Service),
		RpcMetrics: nil,
		requestId:  0,
	}
	return server
}

// Use this method to register RPC handlers.
func (server *ProtobufRPCHTTPServer) RegisterRPCHandlers() {
	th := HandleProtobufRPC(server)
	server.HttpServer.RegisterHandler("/rpc/", th)
}

// Use this method to enable RPC NuSights Metrics.
func (server *ProtobufRPCHTTPServer) RegisterRPCMetricsHandlers(metrics RPCServerMetrics) {
	server.RpcMetrics = metrics
}

// Use this method to export RPC service.
// Arg:
//    rpcSvc: struct of type Service which holds data required for a service.
// Return:
//    bool: true if successful, false otherwise.
func (server *ProtobufRPCHTTPServer) ExportService(rpcsvc *Service) bool {
	if rpcsvc == nil || rpcsvc.Desc.Name == "" {
		glog.Errorf("Invalid service name: %#v", *rpcsvc)
		return false
	}

	fullName := strings.ToLower(rpcsvc.Desc.Name)
	if _, ok := server.Services[fullName]; ok {
		glog.Error("Service is already registered")
		return false
	}

	glog.Info("Registering service: ", fullName)
	server.Services[fullName] = rpcsvc
	return true
}

// Use this method to unexport RPC service.
// Arg:
//    rpcSvc: struct of type Service which holds data required for a service.
// Return:
//    bool: true if successful, false otherwise.
func (server *ProtobufRPCHTTPServer) UnExportService(rpcsvc *Service) bool {
	if rpcsvc == nil || rpcsvc.Desc.Name == "" {
		glog.Errorf("Invalid service name to unexport: %#v", *rpcsvc)
		return false
	}

	fullName := strings.ToLower(rpcsvc.Desc.Name)
	glog.Info("Unregistering service: ", fullName)

	if _, ok := server.Services[fullName]; ok {
		delete(server.Services, fullName)
		return true
	}

	glog.Infof("Service %s was not present", rpcsvc.Desc.Name)
	return true
}

// Use this method to start HTTP Protobuf RPC server.
// Args:
//      Nothing.
// Return:
func (server *ProtobufRPCHTTPServer) StartProtobufRPCHTTPServer() error {
	glog.Infof("Starting protobuf RPC HTTP Server on %s:%s",
		server.HttpServer.ServerAddress, server.HttpServer.ServerPort)

	server.RegisterRPCHandlers()

	err := server.HttpServer.StartHTTPServer()
	if err != nil {
		glog.Error("Error in starting HTTP server: ", err)
		return err
	}
	return nil
}
