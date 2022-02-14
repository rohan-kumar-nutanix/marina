package zeus

import (
        "github.com/nutanix-core/acs-aos-go/go-zookeeper"
        "time"

        glog "github.com/golang/glog"
)

func Connect(hostPortList []string) *zk.Conn {
        conn, _, err := zk.Connect(hostPortList, time.Second*10)
        if err != nil {
                panic(err)
        }
        return conn
}

func KeepWatchingChildrenOf(conn *zk.Conn, path string, f func(*zk.Conn, []string)) {
        for {
                children, _, ch, err := conn.ChildrenWatchWithBackoff(path)
                if err == zk.ErrConnectionClosed {
                        continue
                }
                if err != nil {
                        panic(err)
                }
                f(conn, children)
                <-ch
        }
}

func KeepWatchingContentsOf(conn *zk.Conn, path string, f func(*zk.Conn, []byte)) {
        for {
                data, _, ch, err := conn.GetWatchWithBackoff(path)
                if err == zk.ErrConnectionClosed {
                        continue
                }
                if err != nil {
                        panic(err)
                }
                f(conn, data)
                <-ch
        }
}

func KeepWatchingContentsTillNodeIsDeleted(conn *zk.Conn, path string, f func(*zk.Conn, []byte)) {
        for {
                data, _, ch, err := conn.GetWatchWithBackoff(path)
                if err == zk.ErrConnectionClosed {
                        continue
                }
                if err != nil {
                        glog.Infof("GetWatchWithBackoff for %s returned error: %s", path, err)
                        if err == zk.ErrNoNode {
                                glog.Infof("Zk path %s got deleted. So exiting the go routine that "+
                                        "was watching its changes.", path)
                                break
                        }
                }
                f(conn, data)
                <-ch
        }
}
