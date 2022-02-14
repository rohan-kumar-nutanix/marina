package net

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"
	"sync"
	"time"

	"github.com/golang/glog"
	"golang.org/x/crypto/ssh"
)

// Flags
var (
	// The path to SSH private key file.
	sshPvtKeyPath = flag.String(
		"ssh_pvt_key_path",
		"/home/nutanix/.ssh/id_rsa",
		"The file from where ssh private key is read.")

	// The system user-name where this client is running.
	userName = flag.String(
		"ssh_user_name",
		"nutanix",
		"The system user-name where this client is running.")
)

const (
	ioTimeOutSec = 15
	maxRetries   = 5
)

type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

type SSHtunnel struct {
	local  *Endpoint
	server *Endpoint
	remote *Endpoint

	connChannel chan net.Conn
	callback    tunnelStoppedCallback
	stopped     bool
	listener    net.Listener
	client      *ssh.Client
	wg          sync.WaitGroup
}

type tunnelStoppedCallback func(error) error

func NewSshClient(localEndpoint *Endpoint, serverEndpoint *Endpoint,
	remoteEndpoint *Endpoint, callback tunnelStoppedCallback) *SSHtunnel {
	tunnel := &SSHtunnel{
		local:       localEndpoint,
		server:      serverEndpoint,
		remote:      remoteEndpoint,
		connChannel: make(chan net.Conn),
		callback:    callback,
		stopped:     false,
		listener:    nil,
	}

	return tunnel
}

func (tunnel *SSHtunnel) Start() error {
	if tunnel.stopped {
		return errors.New("Cannot (re)start a stopped SSH tunnel")
	}

	err := tunnel.startListener()
	for retry := 1; err != nil && retry <= maxRetries; retry++ {
		time.Sleep(time.Duration(math.Exp2(float64(retry))) * time.Second)
		err = tunnel.startListener()
	}

	if err == nil {
		go tunnel.proxyIncomingConnections()
		tunnel.wg.Add(1)
	} else {
		glog.Error("Error in starting listener after retries: ", err)
	}

	return err
}

func (tunnel *SSHtunnel) Stop() {
	tunnel.stopped = true
	if tunnel.listener != nil {
		if err := tunnel.listener.Close(); err != nil {
			glog.Error("Error closing listener: ", err)
		}
		tunnel.listener = nil
	}
	// wait for proxyIncomingConnections to be done
	tunnel.wg.Wait()
}

func (tunnel *SSHtunnel) proxyIncomingConnections() {
	defer tunnel.wg.Done()

	go tunnel.forward()

	for !tunnel.stopped {
		conn, err := tunnel.listener.Accept()
		if err != nil {
			if tunnel.stopped {
				continue
			}
			glog.Error("Error in accepting incoming connection: ", err)
			if nerr, ok := err.(net.Error); ok && (nerr.Temporary() || nerr.Timeout()) {
				continue
			}
			glog.Info("Attempting to recreate listener")
			if err = tunnel.listener.Close(); err != nil {
				glog.Error("Error closing listener to recreate: ", err)
			} else {
				err = tunnel.startListener()
				for retry := 1; !tunnel.stopped && err != nil && retry <= maxRetries; retry++ {
					time.Sleep(time.Duration(math.Exp2(float64(retry))) * time.Second)
					err = tunnel.startListener()
				}
			}
			if !tunnel.stopped && err != nil {
				tunnel.internalStop(err)
			}
		} else if conn != nil {
			glog.Infof("Accepted connection to forward %v", conn)
			tunnel.connChannel <- conn
		}
	}

	close(tunnel.connChannel)
	if tunnel.listener != nil {
		if err := tunnel.listener.Close(); err != nil {
			glog.Error("Error closing listener: ", err)
		}
	}
}

func (tunnel *SSHtunnel) internalStop(tunnelError error) {
	tunnel.Stop()
	if tunnel.callback != nil {
		go func(callback tunnelStoppedCallback, tunnelError error) {
			if err := callback(tunnelError); err != nil {
				glog.Error("Error making tunnel stopped callback: ", err)
			}
		}(tunnel.callback, tunnelError)
	} else {
		glog.Fatal("Stopping tunnel internally due to unrecoverable error: ", tunnelError)
	}
}

func (tunnel *SSHtunnel) startListener() error {
	var err error
	tunnel.listener, err = net.Listen("tcp", tunnel.local.String())
	if err != nil {
		glog.Error("Error creating listener: ", err)
		if tunnel.listener != nil {
			if err = tunnel.listener.Close(); err != nil {
				glog.Error("Error closing listener: ", err)
			}
		}
		tunnel.listener = nil
	}
	return err
}

func (tunnel *SSHtunnel) maybeCreateClient() error {
	var err error = nil
	if tunnel.client == nil {
		sshConfig, configErr := getSSHConfig()
		if configErr != nil {
			return configErr
		}

		tunnel.client, err = ssh.Dial("tcp", tunnel.server.String(),
			sshConfig)
		if err != nil {
			glog.Error("Error creating SSH client: ", err)
			tunnel.closeClient()
		} else {
			glog.Info("Created new SSH client to remote ",
				tunnel.remote.String())
		}
	}
	return err
}

func (tunnel *SSHtunnel) closeClient() error {
	var err error = nil
	if tunnel.client != nil {
		if cerr := tunnel.client.Close(); cerr != nil {
			glog.Error("Error closing SSH client: ", cerr)
		}
		tunnel.client = nil
	}
	return err
}

func (tunnel *SSHtunnel) connectRemote() (net.Conn, error) {
	err := tunnel.maybeCreateClient()
	if err != nil {
		return nil, err
	}

	remoteConn, err := tunnel.client.Dial("tcp", tunnel.remote.String())
	if err != nil {
		glog.Error("Error connecting to remote server: ", err)
		if remoteConn != nil {
			if cerr := remoteConn.Close(); cerr != nil {
				glog.Error("Error closing remote server: ", cerr)
			}
			remoteConn = nil
		}
		tunnel.closeClient()
	}

	return remoteConn, err
}

type IdleTimeoutConn struct {
	hasTimeOut bool
	Connection net.Conn
}

func (self IdleTimeoutConn) Read(buf []byte) (int, error) {
	if self.hasTimeOut {
		self.Connection.SetDeadline(time.Now().Add(
			ioTimeOutSec * time.Second))
	} else {
		self.hasTimeOut = true
	}
	return self.Connection.Read(buf)
}

func (self IdleTimeoutConn) Write(buf []byte) (int, error) {
	if self.hasTimeOut {
		self.Connection.SetDeadline(time.Now().Add(
			ioTimeOutSec * time.Second))
	} else {
		self.hasTimeOut = true
	}
	return self.Connection.Write(buf)
}

func waitCopyToFinish(closeWg *sync.WaitGroup, conn1, conn2 net.Conn) {
	closeWg.Wait()
	conn1.Close()
	conn2.Close()
}

func copyConn(writer, reader net.Conn, closeWg *sync.WaitGroup) {
	defer closeWg.Done()

	wcon := IdleTimeoutConn{
		hasTimeOut: false,
		Connection: writer,
	}
	rcon := IdleTimeoutConn{
		hasTimeOut: false,
		Connection: reader,
	}

	_, err := io.Copy(wcon, rcon)
	if err != nil {
		glog.Error("io.Copy error: ", err)
	}
}

func (tunnel *SSHtunnel) forward() {
	for localConn := range tunnel.connChannel {
		if localConn == nil {
			continue
		}

		if tunnel.stopped {
			if err := localConn.Close(); err != nil {
				glog.Error("Error closing local conn: ", err)
			}
			continue
		}

		glog.Info("establishing proxy connection for %v", localConn)
		remoteConn, err := tunnel.connectRemote()
		for retry := 1; !tunnel.stopped && err != nil && retry <= maxRetries; retry++ {
			time.Sleep(time.Duration(math.Exp2(float64(retry))) * time.Second)
			remoteConn, err = tunnel.connectRemote()
		}

		if err != nil {
			glog.Error("Giving up connect to remote after retries: ", err)
			if err = localConn.Close(); err != nil {
				glog.Error("Error closing local conn: ", err)
			}
			continue
		}

		if tunnel.stopped {
			glog.Error("Closing unforwarded connection due to tunnel stop %v",
				localConn)
			if err := localConn.Close(); err != nil {
				glog.Error("Error closing local conn: ", err)
			}
			if remoteConn != nil {
				if err = remoteConn.Close(); err != nil {
					glog.Error("Error closing remote conn: ", err)
				}
			}
			continue
		}

		var closeWg sync.WaitGroup
		closeWg.Add(1)
		go copyConn(localConn, remoteConn, &closeWg)
		closeWg.Add(1)
		go copyConn(remoteConn, localConn, &closeWg)
		go waitCopyToFinish(&closeWg, localConn, remoteConn)

	}
	tunnel.closeClient()
	glog.Info("Connection forwarding ended")
}

func getSSHConfig() (*ssh.ClientConfig, error) {
	encoded, err := ioutil.ReadFile(*sshPvtKeyPath)
	if err != nil {
		glog.Error("unable to read private key: ", err)
		return nil, err
	}

	key, _ := pem.Decode(encoded)
	if key == nil {
		errMsg := "No pem block found"
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(key.Bytes)
	if err != nil {
		privateKey, err = x509.ParsePKCS1PrivateKey(key.Bytes)
	}
	if err != nil {
		glog.Error("unable to parse private key: ", err)
		return nil, err
	}

	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		glog.Error("unable to create signer: ", err)
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: *userName,
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr,
			key ssh.PublicKey) error {
			return nil
		},
	}

	return config, nil
}
