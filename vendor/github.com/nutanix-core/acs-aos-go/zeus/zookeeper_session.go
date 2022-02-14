package zeus

import (
	"bytes"
	"encoding/binary"
	"flag"
	fmt "fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	glog "github.com/golang/glog"
	zk "github.com/nutanix-core/acs-aos-go/go-zookeeper"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
)

type ZookeeperSession struct {
	lock chan bool
	Conn *zk.Conn
}

type ClientID struct {
	SessionID int64
	Passwd    [16]byte
}

var (
	emptyPassword = [16]byte{0}

	zeusSessionFilePath = flag.String(
		"zeus_session_file_path",
		"/home/nutanix/data/zeus/gosessions",
		"Path to the zeus sessions file")
)

const SessionTimeout = 20 * time.Second

func NewZookeeperSessionFromEnvVar() (*ZookeeperSession, error) {
	zk_servers := strings.Split(os.Getenv("ZOOKEEPER_HOST_PORT_LIST"), ",")
	return NewZookeeperSession(zk_servers, SessionTimeout)
}

func NewZookeeperSession(servers []string, sessionTimeout time.Duration) (
	*ZookeeperSession, error) {
	return NewZookeeperSessionWithSessionId(servers, sessionTimeout,
		ClientID{SessionID: 0, Passwd: emptyPassword})
}

func NewZookeeperSessionWithSessionId(servers []string, sessionTimeout time.Duration,
	clId ClientID) (*ZookeeperSession, error) {

	glog.Infof("Creating zookeeper session with client id: 0x%x", clId.SessionID)
	conn, eventChan, err := zk.ConnectWithSessionIdPassword(
		servers, sessionTimeout, clId.SessionID, clId.Passwd)
	if err != nil {
		return nil, err
	}
	ch := make(chan bool, 1)
	zk_session := &ZookeeperSession{ch, conn}
	if clId.SessionID == 0 {
		go zk_session.monitor_session(eventChan, false)
	} else {
		go zk_session.monitor_session(eventChan, true)
	}

	return zk_session, nil
}

// If a service restarts, then there'll be a fork. 'sessionIsOld' is used to
// differentiate between an old and new session.
//
// We need to know whether a session is new or old b/c in the event of an
// expiration event, we FATAL new sessions, but ignore it for old sessions.
func (zk_session *ZookeeperSession) monitor_session(
	ch <-chan zk.Event, sessionIsOld bool) {

	for {
		ev, ok := <-ch
		if !ok {
			glog.Warningf("The session event channel is closed, "+
				"exiting the session monitor for: 0x%x",
				zk_session.Conn.SessionID())
			return
		}

		if ev.Type == zk.EventSession {
			switch state := ev.State; state {
			case zk.StateHasSession:
				select {
				case zk_session.lock <- true:
					glog.Infof("ZK: session 0x%x reconnected",
						zk_session.Conn.SessionID())
				default:
				}
			case zk.StateExpired:
				if !sessionIsOld {
					glog.Fatalf("ZK: session: 0x%x has expired",
						zk_session.Conn.SessionID())
				}
				zk_session.lock <- true
				return
			case zk.StateClosed:
				glog.Infof("ZK: session: 0x%x is closed",
					zk_session.Conn.SessionID())
				return
			default:
				select {
				case <-zk_session.lock:
					glog.Infof("ZK: session 0x%x disconnected, "+
						"waiting for it to reconnect",
						zk_session.Conn.SessionID())
				default:
				}
			}
		}
	}
}

// This function checks if a value at "path" exists and returns true if it does and false otherwise.
func (zk_session *ZookeeperSession) Exist(path string, retry bool) (bool, *zk.Stat, error) {
	if retry {
		return zk_session.Conn.ExistsWithBackoff(path)
	}
	return zk_session.Conn.Exists(path)
}

// ExistW sets a watch on the data and raises an event zk.Event if the value of that data changes.
func (zk_session *ZookeeperSession) ExistW(path string, retry bool) (bool, *zk.Stat, <-chan zk.Event, error) {
	if retry {
		return zk_session.Conn.ExistsWatchWithBackoff(path)
	}
	return zk_session.Conn.ExistsW(path)
}

// This function gets data at "path" and returns the data
func (zk_session *ZookeeperSession) Get(path string, retry bool) ([]byte, *zk.Stat, error) {
	if retry {
		return zk_session.Conn.GetWithBackoff(path)
	}
	return zk_session.Conn.Get(path)
}

// GetW sets a watch on the data and raises an event zk.Event if the value of that data changes.
func (zk_session *ZookeeperSession) GetW(path string, retry bool) ([]byte, *zk.Stat, <-chan zk.Event, error) {
	if retry {
		return zk_session.Conn.GetWatchWithBackoff(path)
	}
	return zk_session.Conn.GetW(path)
}

// This function lists nodes at "path"
func (zk_session *ZookeeperSession) Children(path string, retry bool) ([]string, *zk.Stat, error) {
    if retry {
        return zk_session.Conn.ChildrenWithBackoff(path)
    }
    return zk_session.Conn.Children(path)
}

// ChildrenW sets a watch on the path and raises an event zk.Event if the value at that path changes.
func (zk_session *ZookeeperSession) ChildrenW(path string, retry bool) ([]string, *zk.Stat, <-chan zk.Event, error) {
    if retry {
        return zk_session.Conn.ChildrenWatchWithBackoff(path)
    }
    return zk_session.Conn.ChildrenW(path)
}

// This function creates data at "path" and returns the path of response
func (zk_session *ZookeeperSession) Create(path string, data []byte,
	flags int32, acl []zk.ACL, retry bool) (string, error) {
	if retry {
		return zk_session.Conn.CreateWithBackoff(path, data, flags, acl)
	}
	return zk_session.Conn.Create(path, data, flags, acl)
}

// This function deletes data at "path" and returns the error if exists
func (zk_session *ZookeeperSession) Delete(path string, version int32, retry bool) error {
	if retry {
		return zk_session.Conn.DeleteWithBackoff(path, version)
	}
	return zk_session.Conn.Delete(path, version)
}

// TODO: Add Timeout
func (zk_session *ZookeeperSession) WaitForConnection() bool {
	<-zk_session.lock
	zk_session.lock <- true
	return true
}

func (zk_session *ZookeeperSession) Close() {
	zk_session.Conn.Close()
}

// Read session file.
func ReadSessionFile(name string) ([]byte, error) {
	path := filepath.Join(*zeusSessionFilePath, name)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Errorf("Error %s reading %s", err.Error(), path)
		return nil, err
	}
	return data, err
}

// Write session file.
func WriteSessionFile(name string, sessionID int64, passwd []byte) error {
	path := *zeusSessionFilePath
	var password [16]byte
	copy(password[:], passwd)
	clID := ClientID{SessionID: sessionID, Passwd: password}
	glog.Infof("Saving current zookeeper session: 0x%x ", clID.SessionID)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		glog.Errorf("Failed to create %s with err %s", path, err)
		return err
	}
	path = filepath.Join(path, name)
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, clID)
	if err != nil {
		glog.Error("Failed to convert client id to bytes with ", err)
		return err
	}
	err = ioutil.WriteFile(path, buf.Bytes(), 0755)
	if err != nil {
		glog.Errorf("Error %s while writing %s", err.Error(), path)
		return err
	}
	return nil
}

// Close old session and persist current one.
func (zk_session *ZookeeperSession) MayBeCloseOldSession(servers []string,
	sessionTimeout time.Duration, sessionName string) {
	if sessionName == "" {
		glog.Error("Zookeeper old session close needs session name")
		return
	}
	data, err := ReadSessionFile(sessionName)
	if err == nil {
		clID := ClientID{}
		buf := bytes.NewBuffer(data)
		err = binary.Read(buf, binary.BigEndian, &clID)
		if err != nil {
			glog.Errorf("Failed to convert %s session to struct ",
				sessionName)
		} else {
			glog.Infof("Connecting to old client id: 0x%x ", clID.SessionID)

			zkSession, err := NewZookeeperSessionWithSessionId(servers, sessionTimeout,
				clID)
			if err != nil {
				glog.Error("Failed to establish zookeeper session with error ",
					err.Error())
			}

			if zkSession.WaitForConnection() {
				glog.Infof("Closing old zk session for session id %x", clID.SessionID)
				zkSession.Close()
			}
		}
	} else {
		glog.Error("Failed to read session id for ", sessionName)

	}
	_ = WriteSessionFile(sessionName, zk_session.Conn.SessionID(),
		zk_session.Conn.Password())
}

// CreateZkPath ...
func (zk_session *ZookeeperSession) CreateZkPath(path string, data []byte,
	flags int32, acl []zk.ACL) (string, error) {

	conn := zk_session.Conn

	for {
		ret, errC := conn.Create(path, data, flags, acl)
		if errC != nil && errC == zk.ErrNoNode {
			retP, errP := zk_session.CreateZkPath(filepath.Dir(path), []byte{}, flags,
				acl)
			if errP != nil {
				fmt.Printf("Zk create %s failed: %q\n", filepath.Dir(path), errC)
				return retP, errP
			} else {
				fmt.Printf("Zk create of parent path %s suceeded\n",
					filepath.Dir(path))
				continue
			}
		}
		return ret, errC
	}
	return "", nil
}

// DeleteZkPath ...
func (zk_session *ZookeeperSession) DeleteZkPath(path string,
	version int32) error {

	errD := zk_session.Conn.Delete(path, version)
	if errD != zk.ErrNoNode {
		return nil
	}
	return errD
}

// SetWithBackOff ...
func (zk_session *ZookeeperSession) SetWithBackOff(path string,
	data []byte, version int32) (*zk.Stat, error) {

	var stat *zk.Stat
	var err error
	connRetrySecs := zk.GetZkConnRetryInterval()

	// Keep retrying till session timeout if Server is not reachable.
	cb := misc.NewConstantBackoff(
		time.Duration(connRetrySecs)*time.Second,
		misc.UnlimitedRetry)

	conn := zk_session.Conn
	timeAtBeginning := time.Now()
	for {
		stat, err = conn.Set(path, data, version)

		// Keep retrying only for the following 2 errors.
		if !(err == zk.ErrNoServer || err == zk.ErrConnectionClosed) {
			break
		}

		timeLapsed := time.Since(timeAtBeginning)
		if timeLapsed > time.Duration(conn.SessionTimeoutMs()/1000)*time.Second {
			return stat, err
		}

		_ = cb.Backoff()
	}
	return stat, err
}

// GetWithBackOff ...
func (zk_session *ZookeeperSession) GetWithBackOff(path string) ([]byte,
	*zk.Stat, error) {

	var stat *zk.Stat
	var err error
	var data []byte

	// Keep retrying till session timeout if Server is not reachable.
	cb := misc.NewConstantBackoff(
		time.Duration(zk.GetZkConnRetryInterval())*time.Second,
		misc.UnlimitedRetry)

	conn := zk_session.Conn
	timeAtBeginning := time.Now()
	for {
		data, stat, err = conn.Get(path)

		// Keep retrying only for the following 2 errors.
		if !(err == zk.ErrNoServer || err == zk.ErrConnectionClosed) {
			break
		}

		timeLapsed := time.Since(timeAtBeginning)
		if timeLapsed > time.Duration(conn.SessionTimeoutMs()/1000)*time.Second {
			return data, stat, err
		}

		_ = cb.Backoff()
	}
	return data, stat, err
}
