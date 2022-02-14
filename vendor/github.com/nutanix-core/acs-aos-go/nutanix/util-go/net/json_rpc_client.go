//
// Copyright (c) 2016 Nutanix Inc. All rights reserved.
//
// Author: akshay@nutanix.com
//
// RPC server implementation.

package net

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

type JsonRpcClient struct {
	Port   int
	RpcUrl string
}

// Times out in 30 seconds.
// Synchronous.
func (client *JsonRpcClient) Call(
	oid, method string, params interface{}) (string, error) {

	params_bytes, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	oid_bytes, err := json.Marshal(oid)
	if err != nil {
		return "", err
	}
	oid_msg := json.RawMessage(oid_bytes)

	params_msg := json.RawMessage(params_bytes)
	arg := ServerRequest{
		Id:     &oid_msg,
		Method: method,
		Params: &params_msg,
	}

	payload, err := json.Marshal(arg)
	if err != nil {
		return "", err
	}

	rpcUrl := "http://localhost:" + strconv.Itoa(client.Port) + client.RpcUrl
	req, err := http.NewRequest("POST", rpcUrl, bytes.NewBuffer(payload))

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
