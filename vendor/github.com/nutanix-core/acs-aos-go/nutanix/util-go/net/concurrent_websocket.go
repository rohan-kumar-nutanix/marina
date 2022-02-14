/*
 * Copyright (c) 2017 Nutanix Inc. All rights reserved.
 *
 * We need a concurrent websocket design that can handle concurrent writes.
 */

package net

import (
	"github.com/golang/glog"
)

// We'd like an interface around 'Conn' so that we can use different websocket
// implementations if we need to. For example, we might use
// leavengood/websocket, or we might use gorilla/websocket.
type Conn interface {
	WriteJSON(v interface{}) error
}

type ConcurrentWebsocket struct {
	ws      Conn
	newjson chan interface{}
	close   chan bool
}

func NewConcurrentWebsocket(ws Conn) *ConcurrentWebsocket {
	cws := new(ConcurrentWebsocket)
	cws.ws = ws
	cws.newjson = make(chan interface{})
	cws.close = make(chan bool)
	go cws.waitForEvents()
	return cws
}

//-----------------------------------------------------------------------------
// Interface
//-----------------------------------------------------------------------------

func (cws *ConcurrentWebsocket) WriteJSON(v interface{}) {
	cws.WriteJSONChan() <- v
}

//-----------------------------------------------------------------------------

func (cws *ConcurrentWebsocket) WriteJSONChan() chan interface{} {
	return cws.newjson
}

//-----------------------------------------------------------------------------

func (cws *ConcurrentWebsocket) Close() {
	cws.close <- true
}

//-----------------------------------------------------------------------------
// Internals.
//-----------------------------------------------------------------------------

func (cws *ConcurrentWebsocket) waitForEvents() {
outer:
	for {
		select {
		case v := <-cws.newjson:
			cws.writeJSON(v)
		case <-cws.close:
			break outer
		}
	}
}

//-----------------------------------------------------------------------------

func (cws *ConcurrentWebsocket) writeJSON(v interface{}) {
	if err := cws.ws.WriteJSON(v); err != nil {
		glog.Errorln(err)
	}
}
