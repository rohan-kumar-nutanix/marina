//
// Copyright (c) 2016 Nutanix Inc. All rights reserved.
//
// Author: akshay@nutanix.com
//
// RPC server implementation.

package net

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/golang/glog"
)

var null = json.RawMessage([]byte("null"))

// Struct to hold RPC server data.
type RPCHTTPServer struct {
	HttpServer HTTPServer
	Services   map[string]*service
}

// CodecRequest decodes and encodes a single request.
type CodecRequest struct {
	request *ServerRequest
	err     error
}

// ServerRequest represents a JSON-RPC request received by the server.
type ServerRequest struct {
	// A String containing the name of the method to be invoked.
	Method string `json:".method"`
	// An Array of objects to pass as arguments to the method.
	Params *json.RawMessage `json:".kwargs"`
	// The request id. This can be of any type. It is used to match the
	// response with the request that it is replying to.
	Id *json.RawMessage `json:".oid"`
}

// serverResponse represents a JSON-RPC response returned by the server.
type serverResponse struct {
	// The Object that was returned by the invoked method. This must be null
	// in case there was an error invoking the method.
	Result interface{} `json:"result"`
	// An Error object if there was an error invoking the method. It must be
	// null if there was no error.
	Error interface{} `json:"error"`
	// This must be the same id as the request it is responding to.
	Id *json.RawMessage `json:"id"`
}

// Service destribes name and methods for service exposed via RPC.
type service struct {
	name     string                    // name of service
	rcvr     reflect.Value             // receiver of methods for service
	rcvrType reflect.Type              // type of the receiver
	methods  map[string]*serviceMethod // registered methods
}

// serviceMethod describes a single method in a service.
type serviceMethod struct {
	method    reflect.Method // receiver method
	argsType  reflect.Type   // type of the request argument
	replyType reflect.Type   // type of the response argument
}

// newCodecRequest returns a new CodecRequest.
func newCodecRequest(r *http.Request) *CodecRequest {
	// Decode the request body and check if RPC method is valid.
	req := new(ServerRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	r.Body.Close()
	return &CodecRequest{request: req, err: err}
}

// Method returns the RPC method for the current request.
// The method uses a dotted notation as in "Service.Method".
func (c *CodecRequest) Method() (string, error) {
	if c.err == nil {
		return c.request.Method, nil
	}
	return "", c.err
}

// ReadRequest reads params , decodes and returns params.
func (c *CodecRequest) ReadRequest(args interface{}) error {
	if c.err == nil {
		if c.request.Params != nil {
			// JSON params is array value. RPC params is struct.
			// Unmarshal into array containing the request struct.
			params := &args
			c.err = json.Unmarshal(*c.request.Params, &params)
		} else {
			c.err = errors.New("rpc: method request ill-formed: " +
				"missing params field")
		}
	}
	return c.err
}

// WriteResponse encodes the response and writes it to the ResponseWriter.
func (c *CodecRequest) WriteResponse(w http.ResponseWriter, reply interface{}) {
	if c.request.Id != nil {
		// Id is null for notifications and they don't have a response.
		res := &serverResponse{
			Result: reply,
			Error:  &null,
			Id:     c.request.Id,
		}
		glog.Info("Write Response %+v", *res.Result.(*int))
		c.writeServerResponse(w, 200, res)
	}
}

// Write server respose.
func (c *CodecRequest) writeServerResponse(
	w http.ResponseWriter, status int, res *serverResponse) {
	b, err := json.Marshal(res)
	if err == nil {
		w.WriteHeader(status)
		w.Header().Set(
			"Content-Type", "application/json; charset=utf-8")
		w.Write(b)
	}
}

// The method name uses a dotted notation as in "Service.Method".
func (server *RPCHTTPServer) get(
	method string) (*service, *serviceMethod, error) {
	parts := strings.Split(method, ".")
	if len(parts) != 2 {
		err := fmt.Errorf(
			"rpc: service/method request ill-formed: %q", method)
		return nil, nil, err
	}
	service := server.Services[parts[0]]
	if service == nil {
		err := fmt.Errorf("rpc: can't find service %q", method)
		return nil, nil, err
	}
	serviceMethod := service.methods[parts[1]]
	if serviceMethod == nil {
		err := fmt.Errorf("rpc: can't find method %q", method)
		return nil, nil, err
	}
	return service, serviceMethod, nil
}

// This method handles JSON RPCs.
func HandleJsonRPC(server *RPCHTTPServer) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		codecReq := newCodecRequest(r)

		// Get service method to be called.
		method, errMethod := codecReq.Method()
		if errMethod != nil {
			glog.Error("ERROR in finding method: ", errMethod)
			return
		}

		serviceSpec, methodSpec, _ := server.get(method)

		args := reflect.New(methodSpec.argsType)
		errRead := codecReq.ReadRequest(args.Interface())
		if errRead != nil {
			glog.Errorf("Error while reading request: %#v", errRead)
		}

		reply := reflect.New(methodSpec.replyType)
		_ = methodSpec.method.Func.Call(
			[]reflect.Value{serviceSpec.rcvr, args, reply})

		codecReq.WriteResponse(w, reply.Interface())
	}
	return http.HandlerFunc(fn)
}

// Use this method to instantiate JSON RPC server.
func NewRPCHTTPServer(Server *HTTPServer) *RPCHTTPServer {
	glog.Info("Initializing RPC HTTP server with address %s, port %s",
		Server.ServerAddress, Server.ServerPort)
	server := &RPCHTTPServer{
		HttpServer: *Server,
		Services:   make(map[string]*service),
	}
	return server
}

// Register HTTP hanlers.
func (server *RPCHTTPServer) RPCRegisterHandlers() {
	th := HandleJsonRPC(server)
	server.HttpServer.RegisterHandler("/jsonrpc", th)
}

// Register RPC service.
func (server *RPCHTTPServer) RPCRegisterService(
	rcvr interface{}, name string) error {
	s := &service{
		name:     name,
		rcvr:     reflect.ValueOf(rcvr),
		rcvrType: reflect.TypeOf(rcvr),
		methods:  make(map[string]*serviceMethod),
	}

	paramOffset := 0
	for i := 0; i < s.rcvrType.NumMethod(); i++ {
		method := s.rcvrType.Method(i)
		mtype := method.Type
		args := mtype.In(1 + paramOffset)
		reply := mtype.In(2 + paramOffset)

		s.methods[method.Name] = &serviceMethod{
			method:    method,
			argsType:  args.Elem(),
			replyType: reply.Elem(),
		}
	}

	server.Services[s.name] = s

	return nil
}

// Start JSON RPC Server.
func (server *RPCHTTPServer) StartRPCHTTPServer() error {
	glog.Info("Starting RPC HTTP Server on address %s, port %s",
		server.HttpServer.ServerAddress, server.HttpServer.ServerPort)

	server.RPCRegisterHandlers()

	err := server.HttpServer.StartHTTPServer()
	if err != nil {
		glog.Error("Failed to start the HTTP server")
		return err
	}
	return nil
}
