//
// Copyright (c) 2017 Nutanix Inc. All rights reserved.

// Author: zxie@nutanix.com

// Package that provide transport using the Aplos Remote Connection
// to send REST API from one cluster to another (PE <-> PC)

// Usage:
//     transport, err := NewRcHttpTransport("some cluster uuid", 5)
//     if err != nil {
//        // handle error here
//     }
//     httpClient := http.Client{
//         Transport: transport,
//     }
//     httpClient.Get("https://example.com")

package aplos_transport

import (
        "bytes"
        "crypto/tls"
        "flag"
        "fmt"
        "github.com/golang/glog"
        util_misc "github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
        "github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
        net "github.com/nutanix-core/acs-aos-go/nutanix/util-slbufs/util/sl_bufs/net"
        "io/ioutil"
        "net/http"
        "strings"
        "time"
        "github.com/nutanix-core/acs-aos-go/zeus"
)

const (
        retryInitialDelaySecs = 1
        retryMaxDelaySecs     = 30
        maxRetries            = 6
)

var (
        fanoutPort             uint16 = 9440
        useHttpTransport       bool   = false
        disableChunkedEncoding bool   = false
        defaultServerAddress   string = "127.0.0.1"
)

func init() {
        flag.Var(&util_misc.PortNumber{&fanoutPort}, "fanout_proxy_port",
                "Port number of the fanout proxy.")
        flag.BoolVar(&useHttpTransport, "fanout_use_http_transport", false,
                "If true, we will use HTTP instead of HTTPS for sending fanout requests.")
        flag.BoolVar(&disableChunkedEncoding, "fanout_disable_chuncked_encoding",
                true, "Whether to disable chunked transfer encoding for fanout requests.")
        flag.StringVar(&defaultServerAddress, "default_server_address",
                "127.0.0.1", "The default server address")
}

type RcHttpTransport struct {
        clusterUuid    *uuid4.Uuid
        userUuid       *uuid4.Uuid
        tenantUuid     *uuid4.Uuid
        timeoutSeconds int
        zkSession      *zeus.ZookeeperSession
        waitForTask    bool
        userName       *string
        reqContext     *net.RpcRequestContext
}

func NewRcHttpTransport(zkSession *zeus.ZookeeperSession,
        clusterUuid *uuid4.Uuid, userUuid *uuid4.Uuid,
        tenantUuid *uuid4.Uuid, timeout int, userName *string,
        reqContext *net.RpcRequestContext) (*RcHttpTransport, error) {
        transport := &RcHttpTransport{
                clusterUuid:    clusterUuid,
                userUuid:       userUuid,
                tenantUuid:     tenantUuid,
                timeoutSeconds: timeout,
                zkSession:      zkSession,
                waitForTask:    false,
                userName:       userName,
                reqContext:     reqContext,
        }
        return transport, nil
}

func NewRcHttpTransportWithChildTask(zkSession *zeus.ZookeeperSession,
        clusterUuid *uuid4.Uuid, userUuid *uuid4.Uuid,
        tenantUuid *uuid4.Uuid, timeout int, userName *string,
        reqContext *net.RpcRequestContext) (*RcHttpTransport, error) {
        transport, err := NewRcHttpTransport(
                zkSession, clusterUuid, userUuid, tenantUuid, timeout, userName, reqContext)
        transport.waitForTask = true
        return transport, err
}

func (transport *RcHttpTransport) RoundTrip(
        request *http.Request) (*http.Response, error) {
        uri := strings.TrimPrefix(request.URL.RequestURI()[1:], "api/nutanix/v3/")

        var bodyBytes []byte
        if request.Body != nil {
                var readErr error
                bodyBytes, readErr = ioutil.ReadAll(request.Body)
                if readErr != nil {
                        glog.Warning("Unable to read request body: ", readErr)
                        return nil, readErr
                }
        }

        fanoutUrl := fmt.Sprintf("https://" + defaultServerAddress +
                ":%d/api/nutanix/v3/fanout_proxy", fanoutPort)
        if useHttpTransport {
                fanoutUrl = fmt.Sprintf("http://" + defaultServerAddress +
                        ":%d/v3/fanout_proxy", fanoutPort)
        }
        proxyReq, err := http.NewRequest(http.MethodPost, fanoutUrl,
                ioutil.NopCloser(bytes.NewReader(bodyBytes)))
        if err != nil {
                glog.Warning("Failed to prepare fanout proxy request: ", err)
                return nil, err
        }

        if disableChunkedEncoding {
                // If content length is set, then the transfer encoding type will
                // default to 'identity'.
                proxyReq.ContentLength = int64(len(bodyBytes))
        }

        proxyReq.Header.Set("Content-Type", "application/octet-stream")

        query := proxyReq.URL.Query()
        query.Set("remote_cluster_uuid", transport.clusterUuid.String())
        query.Set("method", request.Method)
        if transport.waitForTask {
                query.Set("wait_for_task", "true")
        }

        // The only scneario when the v3 prefix does not get trimmed above is when
        // that prefix does not exist. That only happens in the case of PC-PE RPCs
        // before the call is made because of the way the remote protobuf client
        // formats the URL. In case of PC-PE RPCs, relevant additional query
        // parameters are set for the remote RPC to work.
        if uri == request.URL.RequestURI()[1:] {
                query.Set("tenant_uuid", transport.tenantUuid.String())
                query.Set("url_path", request.URL.RequestURI())
        } else {
                query.Set("url_path", uri)
        }
        if bodyBytes != nil {
                query.Set("content_type", request.Header.Get("Content-Type"))
        }
        proxyReq.URL.RawQuery = query.Encode()

        useAplosKey := false
        token, err := GetJWT(transport.zkSession, transport.userUuid,
                transport.tenantUuid, useAplosKey, transport.userName, transport.reqContext)
        if err != nil {
                // Server might be in ECDSA key so use aplos key
                useAplosKey = true
                token, err = GetJWT(transport.zkSession, transport.userUuid,
                        transport.tenantUuid, useAplosKey, transport.userName, transport.reqContext)
                if err != nil {
                        glog.Warning("Failed to prepare JWT: ", err)
                        return nil, err
                }
        }

        proxyReq.Header.Add("Authorization", "Bearer "+token)

        tr := &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
        eb := util_misc.NewExponentialBackoff(
                time.Duration(retryInitialDelaySecs)*time.Second,
                time.Duration(retryMaxDelaySecs)*time.Second,
                maxRetries)

        var response *http.Response
        for {
                client := http.Client{
                        Transport: tr,
                        Timeout: time.Duration(transport.timeoutSeconds) *
                                time.Second,
                }
                response, err = client.Do(proxyReq)
                if err == nil && (response.StatusCode == 403 ||
                        (200 <= response.StatusCode &&
                                response.StatusCode <= 499)) {
                        break
                }
                duration := eb.Backoff()
                if duration == util_misc.Stop {
                        break
                }
                proxyReq.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
        }
        if err != nil {
                glog.Warning("Failed to read response from aplos: ", err)
                return nil, err
        }
        // incorrect JWT reslults in a 403 access denied instead of 401
        if response.StatusCode == 403 && !useAplosKey {
                glog.Info("Using aplos specific key for JWT, if not tried with aplos key")
                token, tokenErr := GetJWT(transport.zkSession, transport.userUuid,
                        transport.tenantUuid, true, transport.userName, transport.reqContext)
                if tokenErr != nil {
                        glog.Warning("Failed to prepare JWT: ", tokenErr)
                        return nil, tokenErr
                }

                proxyReq.Header.Set("Authorization", "Bearer "+token)

                eb = util_misc.NewExponentialBackoff(
                        time.Duration(retryInitialDelaySecs)*time.Second,
                        time.Duration(retryMaxDelaySecs)*time.Second,
                        maxRetries)
                for {
                        client := http.Client{
                                Transport: tr,
                                Timeout: time.Duration(
                                        transport.timeoutSeconds) *
                                        time.Second,
                        }
                        proxyReq.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
                        response, err = client.Do(proxyReq)
                        if err == nil && 200 <= response.StatusCode &&
                                response.StatusCode <= 499 {
                                break
                        }
                        duration := eb.Backoff()
                        if duration == util_misc.Stop {
                                break
                        }
                }
                if err != nil {
                        glog.Warning("Failed to read response from aplos: ", err)
                }
        }
        return response, err
}
