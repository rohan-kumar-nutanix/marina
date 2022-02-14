// Copyright (c) 2016 Nutanix Inc. All rights reserved.
//
// Author: akshay@nutanix.com

// HTTPServer with an API that allows handlers to be registered.
package net

import (
	"errors"
	"net"
	"net/http"

	"github.com/golang/glog"
)

const DebugRequestsBaseUrl = "/h/"

// Struct to hold http server details.
type HTTPServer struct {
	ServerAddress string
	ServerPort    string
	Listener      net.Listener
	Instance      *http.Server
}

// Use this method to instantiate new http server.
// Arg:
//    ServerAddress : IP address of server.
//    ServerPort    : Port number of server.
// Return :
//        Instance of HTTP server.
func NewHTTPServer(ServerAddress string, ServerPort string) *HTTPServer {

	glog.Infof("Initializing HTTP server with address %s:%s",
		ServerAddress, ServerPort)
	server := &HTTPServer{
		ServerAddress: ServerAddress,
		ServerPort:    ServerPort,
		Listener:      nil,
		Instance:      nil,
	}
	return server
}

// Use this method to register HTTP handlers.
func (server *HTTPServer) RegisterHandler(
	url string, handler http.Handler) {

	glog.Info("Register a handler for URL: ", url)
	http.Handle(url, handler)
}

func (server *HTTPServer) RegisterHandlerFunc(
	pattern string, handler func(http.ResponseWriter, *http.Request)) {

	glog.Info("Register a func handler for URL: ", pattern)
	http.HandleFunc(pattern, handler)
}

func (server *HTTPServer) RegisterDebugHandlerFunc(pattern string,
	handler func(http.ResponseWriter, *http.Request)) {
	glog.Info("Registering debug handler for URL:", pattern)
	http.HandleFunc(DebugRequestsBaseUrl+pattern, handler)
}

// Use this method to start HTTP server.
// Arg:
//    None.
// Return:
//       error : Standard error.
func (server *HTTPServer) StartHTTPServer() error {
	glog.Info("Starting HTTP Server\n")
	addr := server.ServerAddress + ":" + server.ServerPort
	server.Instance = &http.Server{Addr: addr}
	err := server.Instance.ListenAndServe()
	if err != nil {
		glog.Infof("Error in starting HTTP server ", err)
		return err
	}
	return nil
}

func (server *HTTPServer) StartListening() error {
	var err error
	addr := server.ServerAddress + ":" + server.ServerPort

	glog.Infof("Listening for TCP connections at: %s", addr)
	server.Listener, err = net.Listen("tcp", addr)
	return err
}

func (server *HTTPServer) StartServing() error {
	if server.Listener == nil {
		return errors.New("Listener is nil")
	}
	server.Instance = &http.Server{}
	err := server.Instance.Serve(server.Listener)
	return err
}
