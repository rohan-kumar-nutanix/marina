/*
 * Copyright (c) 2018 Nutanix Inc. All rights reserved.
 *
 * Author: yiran.deng@nutanix.com
 *
 * WinShell enables clients to talk to and execute powershell commands on local
 * Hyper-V host.
 *
 * Requires the NutanixHostAgent service to be running on the host.
 */
package net

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ShellState int

const (
	HostAgentHandshake = "Hello Nutanix\n"
	NewLine            = byte('\n')

	New ShellState = iota
	Connected
	Closed
)

var (
	underlyingCommand = flag.String(
		"underlying_command",
		"powershell",
		"the underlying command to run with when sending request to host",
	)

	invalidResponseError = fmt.Errorf("server sends invalid response")
)

type NutanixHostAgentResponse struct {
	HadErrors bool
	Output    string
	Error     string
	Warning   string
	Verbose   string
	Debug     string
}

type WinShell struct {
	state       ShellState
	certFile    string
	keyFile     string
	hostIp      string
	hostPort    int
	timeoutSecs int
	conn        *tls.Conn
	lock        *sync.Mutex
}

func NewWinShell(certFile, keyFile, ip string,
	port, connectTimeoutSecs int) *WinShell {
	return &WinShell{
		state:       New,
		certFile:    certFile,
		keyFile:     keyFile,
		hostIp:      ip,
		hostPort:    port,
		timeoutSecs: connectTimeoutSecs,
		conn:        nil,
		lock:        &sync.Mutex{},
	}
}

func (shell *WinShell) IsConnected() bool {
	return shell.state == Connected
}

// Open a WinShell connection.
func (shell *WinShell) Open() error {
	if shell.state != New {
		return fmt.Errorf("WinShell objects can be opened at most once")
	}

	cert, err := tls.LoadX509KeyPair(shell.certFile, shell.keyFile)
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		// This flag means TLS will accept any certificate presented by the
		// server and any host name in that certificate.
		// TODO: can we remove it?
	}

	dialer := &net.Dialer{
		Timeout: time.Second * time.Duration(shell.timeoutSecs),
	}
	addr := fmt.Sprintf("%s:%d", shell.hostIp, shell.hostPort)

	shell.conn, err = tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	if err != nil {
		return err
	}

	if _, err = io.WriteString(shell.conn, HostAgentHandshake); err != nil {
		return err
	}

	shell.state = Connected
	return nil
}
func (shell *WinShell) Close() error {
	if shell.IsConnected() {
		shell.conn.Close()
		shell.state = Closed
		return nil
	} else {
		return fmt.Errorf("can only close connected WinShell objects")
	}
}

// Execute command in WinShell.
func (shell *WinShell) Execute(command string,
	timeoutSecs int64) (*NutanixHostAgentResponse, error) {
	// Never execute more than one command at the same time,
	// otherwise we will get mixed response from the server.
	shell.lock.Lock()
	defer shell.lock.Unlock()

	if !shell.IsConnected() {
		return nil, fmt.Errorf("WinShell is not already connected to server")
	}

	if err := shell.sendCommandRequest(command, timeoutSecs); err != nil {
		return nil, fmt.Errorf("failed to send request: %s", err)
	}
	response := &NutanixHostAgentResponse{
		Output: "", Error: "", Warning: "", Verbose: "", Debug: "",
	}

	for {
		dataLengthHeader, err := shell.readline()
		if err != nil {
			return nil, fmt.Errorf("error reading response: %s", err)
		} else {
			length, err := strconv.Atoi(dataLengthHeader)
			if err != nil {
				return nil, invalidResponseError
			}

			buff := bytes.NewBuffer(make([]byte, 0, length))
			_, err = io.CopyN(buff, shell.conn, int64(length))
			if err != nil {
				return nil, fmt.Errorf("error reading response: %s", err)
			}

			// Use a map to capture output/error.
			var responseMap map[string]interface{}
			json.Unmarshal(buff.Bytes(), &responseMap)

			v, ok := responseMap["had_errors"]
			if ok {
				// The response is {"had_errors": true/false}
				//
				// We don't rely on the "had_errors" value returned by the
				// host agent because that value is set to true if any error
				// occurred regardless of whether the error was caught using a
				// try-catch block or ignored using the "SilentlyContinue"
				// primitive. Instead, we set the return value to non-zero if
				// the error output is not empty.
				response.HadErrors = v.(bool)
				return response, nil
			} else {
				// The response is {"type": "...", "data": "..."}
				t, ok := responseMap["type"]
				if !ok {
					return nil, invalidResponseError
				}

				typeString, ok := t.(string)
				if !ok {
					return nil, invalidResponseError
				}

				d, ok := responseMap["data"]
				if !ok {
					return nil, invalidResponseError
				}

				data, ok := d.(string)
				if !ok {
					return nil, invalidResponseError
				}

				switch typeString {
				case "output":
					response.Output += data
				case "error":
					response.Error += data
				case "warning":
					response.Warning += data
				case "verbose":
					response.Verbose += data
				case "debug":
					response.Debug += data
				case "progress":
					// We do not maintain the progress data.
				default:
					return nil, invalidResponseError
				}
			}
		}
	}
}

// Read one line from server's response, this is useful to capture the data
// length header.
func (shell *WinShell) readline() (string, error) {
	var line []byte
	var buff [1]byte

	for {
		if _, err := shell.conn.Read(buff[:]); err != nil {
			return "", err
		}

		if buff[0] == NewLine {
			break
		}
		line = append(line, buff[0])
	}

	// Trim "\r" since Windows uses "\r\n" as line break.
	return strings.TrimRight(string(line), "\r"), nil
}

func (shell *WinShell) sendCommandRequest(command string, timeout int64) error {
	// Prepare JSON packet for command.
	jsonMap := make(map[string]string)
	jsonMap["command"] = *underlyingCommand
	jsonMap["script"] = command
	jsonMap["timeout"] = fmt.Sprintf("%d", timeout*1000)

	cmdJSON, err := json.Marshal(jsonMap)
	if err != nil {
		return err
	}

	dataLengthHeader := fmt.Sprintf("%d\n", len(cmdJSON))
	// Send the JSON packet.
	if _, err = io.WriteString(shell.conn, dataLengthHeader); err != nil {
		return err
	}
	if _, err = io.WriteString(shell.conn, string(cmdJSON)); err != nil {
		return err
	}
	return nil
}

// DownloadFile copies a file from the host to the CVM
// Args:
//   src: The Windows-style path of the source file
//   dst: The destination writer to which the source file will be written to
//
// Returns:
//   int: The number of bytes written to dst
//   error: Error if any
func (shell *WinShell) DownloadFile(src string, dst io.Writer) (int64, error) {
	shell.lock.Lock()
	defer shell.lock.Unlock()
	if !shell.IsConnected() {
		return 0, fmt.Errorf("WinShell is not connected to the server")
	}
	jsonCmdMap := map[string]string{
		"command": "download",
		"path":    src,
	}

	cmdJSON, err := json.Marshal(jsonCmdMap)
	if err != nil {
		return 0, err
	}
	dataLengthHeader := fmt.Sprintf("%d\n", len(cmdJSON))

	// Send download command.
	if _, err = io.WriteString(shell.conn, dataLengthHeader); err != nil {
		return 0, err
	}
	if _, err = io.WriteString(shell.conn, string(cmdJSON)); err != nil {
		return 0, err
	}

	// Start reading file
	dataLengthHeader, err = shell.readline()
	if err != nil {
		return 0, fmt.Errorf("error reading response: %s", err)
	}
	length, err := strconv.Atoi(dataLengthHeader)
	if err != nil {
		return 0, fmt.Errorf("unexpected response. Expected length of following data. "+
			"Got: %v, Error: %v", dataLengthHeader, err)
	}

	return io.CopyN(dst, shell.conn, int64(length))
}
