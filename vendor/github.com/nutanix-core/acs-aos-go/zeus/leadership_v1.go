//
// Copyright (c) 2016 Nutanix Inc. All rights reserved.
//
// This module provides a distributed lock implemented using Zookeeper
// election algorithm as described here:
//   http://zookeeper.apache.org/doc/trunk/recipes.html#sc_leaderElection
//
// Racing with Zeus Python module is avoided by creating the leaders in a
// different namespace.
//

package zeus

import (
        "errors"
        "fmt"
        "path"
        "sort"
        "strings"
        "time"

        "github.com/golang/glog"
        "github.com/nutanix-core/acs-aos-go/go-zookeeper"
)

type Leadership struct {
        ZkSession     *ZookeeperSession
        electionPath  string
        intentionPath string
}

// GO zookeeper library prepends all intention paths for leadership with
// "_c_" followed by GUID. This creates a problem when sorting paths to
// figure out a leader as guid do not follow strictly incrementing order.
// Solution is filter out this additional information before sorting.
// One peculiar thing is that different libraries behave differently.
// e.g. python library does not prepend guid, so code has to check if guid is
// present and then filter it out.
type LeadershipWithGUID []string

var protectedPrefix = "_c_"

func (s LeadershipWithGUID) Len() int {
        return len(s)
}
func (s LeadershipWithGUID) Swap(i, j int) {
        s[i], s[j] = s[j], s[i]
}
func (s LeadershipWithGUID) Less(i, j int) bool {
        // Path is prefixed in zookeeper library as follows:
        // path = "protectedPrefix (_c_)" + guid (slice of 16) + proposed path
        first := s[i]
        second := s[j]

        // Logic here is : remove guid, protectedPrefix and then compare strings
        if strings.HasPrefix(s[i], protectedPrefix) {
                first = s[i][len(protectedPrefix)+33:]
        }
        if strings.HasPrefix(s[j], protectedPrefix) {
                second = s[j][len(protectedPrefix)+33:]
        }
        status := strings.Compare(first, second)
        if status > 0 {
                return false
        }
        return true
}

func NewGoLeadership(lockName string, zkSession *ZookeeperSession) (*Leadership, error) {
        if zkSession.WaitForConnection() {
                glog.Infof("Ensuring zeus path '%s' exists", GO_LEADERSHIP_ROOT)
                _, err := zkSession.Conn.CreateWithBackoff(GO_LEADERSHIP_ROOT, []byte{}, 0, WORLD_ACL)
                if err != nil && err != zk.ErrNodeExists {
                        return nil, err
                }
        } else {
                return nil, ErrZkNotConnected
        }

        electionPath := fmt.Sprintf("%s/%s", GO_LEADERSHIP_ROOT, lockName)
        if zkSession.WaitForConnection() {
                glog.Infof("Ensuring zeus path '%s' exists", electionPath)
                _, err := zkSession.Conn.CreateWithBackoff(electionPath, []byte{}, 0, WORLD_ACL)
                if err != nil && err != zk.ErrNodeExists {
                        return nil, err
                }
        }
        return NewLeadershipWithPathPrefix(lockName, zkSession, GO_LEADERSHIP_ROOT)
}

func NewPyLeadership(lockName string, zkSession *ZookeeperSession) (*Leadership, error) {
        return NewLeadershipWithPathPrefix(lockName, zkSession, PY_LEADERSHIP_ROOT)
}

func NewLeadership(lockName string, zkSession *ZookeeperSession) (*Leadership, error) {
        return NewLeadershipWithPathPrefix(lockName, zkSession, COMMON_LEADERSHIP_ROOT)
}

func NewLeadershipWithPathPrefix(lockName string, zkSession *ZookeeperSession, electionPathPrefix string) (*Leadership, error) {
        electionPath := fmt.Sprintf("%s/%s", electionPathPrefix, lockName)
        _, _, err := zkSession.Conn.GetWithBackoff(electionPath)
        if err != nil {
                return nil, err
        }
        return &Leadership{ZkSession: zkSession, electionPath: electionPath}, nil
}

func (l *Leadership) createElectionPath() {
        if l.ZkSession.WaitForConnection() {
                _, _ = l.ZkSession.Conn.CreateWithBackoff(l.electionPath, []byte{}, 0, WORLD_ACL)
        }
}

func (l *Leadership) Volunteer(value string) (bool, error) {
        var err error
        l.createElectionPath()
        prefix := fmt.Sprintf("%s/n_", l.electionPath)

        l.intentionPath, err = l.ZkSession.Conn.CreateProtectedEphemeralSequential(prefix, []byte(value), WORLD_ACL)
        if err != nil {
                glog.Errorf("Unable to create leadership intention znode %s: %s", l.electionPath, err)
                return false, ErrCouldNotVolunteer
        }
        return true, nil
}

func (l *Leadership) Contenders(wait bool) ([]string, error) {
        for {
                children, _, ch, err := l.ZkSession.Conn.ChildrenWatchWithBackoff(l.electionPath)
                if err == zk.ErrNoNode {
                        return nil, err
                }
                if err != nil {
                        return nil, ErrGetContendersFailed
                }
                if len(children) > 0 {
                        sort.Sort(LeadershipWithGUID(children))
                        return children, nil
                }
                if wait {
                        <-ch
                } else {
                        return children, nil
                }
        }
}

func (l *Leadership) previousContender() (*string, error) {
        mySequence := path.Base(l.intentionPath)
        contenders, err := l.Contenders(false)
        if err != nil {
                glog.Errorf("Unable to get previous contenders: %s", err)
                return nil, err
        }
        if len(contenders) == 0 {
                glog.Fatal("No contenders found")
        }
        if contenders[0] == mySequence {
                return nil, nil
        }
        for i, s := range contenders {
                if s == mySequence {
                        return &contenders[i-1], nil
                }
        }
        return nil, ErrSelfNotFound
}

func (l *Leadership) HasLeadership() bool {
        str, err := l.previousContender()
        return err == nil && str == nil
}

func (l *Leadership) Wait() bool {
        for {
                previous, err := l.previousContender()
                if err != nil {
                        glog.Infof("Error getting previousContender during leadership wait: %s", err)
                        return false
                }
                if previous == nil {
                        glog.Info("Claiming Leadership")
                        return true
                }
                previousContenderPath := fmt.Sprintf("%s/%s", l.electionPath, *previous)
                exists, _, ch, err := l.ZkSession.Conn.ExistsWatchWithBackoff(previousContenderPath)
                if err != nil {
                        glog.Fatal("Failed to set up watch on the previous contender")
                }
                if exists {
                        <-ch
                }
        }
}

func (l *Leadership) Relinquish() error {
	err := l.ZkSession.Conn.DeleteWithBackoff(l.intentionPath, -1)
	if err != nil {
		return err
	}
	glog.Info("Relinqiushed leadership")
	return nil
}

// CleanElectionPath cleans up the path created for the lockName,
// if there are no intentionPaths or volunteers for it.
// Returns error in case of cleanup failure
func (l *Leadership) CleanElectionPath() error {
        contenders, err := l.Contenders(false)
        if err != nil {
                glog.Errorf("Error getting contenders.ElectionPath" +
                "%s cleanup failed with %s", l.electionPath, err)
                return err
        }
        if len(contenders) == 0 {
                err = l.ZkSession.Conn.DeleteWithBackoff(l.electionPath, -1)
                if err != nil {
                        glog.Errorf("ElectionPath %s cleanup failed with %s", l.electionPath, err)
                        return err
                }
        }
        glog.Infof("Election Path '%s' deleted", l.electionPath)
        return nil
}

func (l *Leadership) LeaderValue(wait bool, watch_cb func(string, string)) (string, string, error) {
        //
        // Returns the path and value proposed by current leader or Nil if there
        // is no leader and wait is False.
        // Arguments:
        //    wait - optional, whether to wait for the leader to be present
        //    watch_cb - optional, one time callback to get fired when the
        //    leader is changed or deleted.
        //    Returns:  (leader path, leader value)
        backoff := time.Duration(1) * time.Second
        for {
                contenders, err := l.Contenders(wait)
                if err == zk.ErrNoNode || (len(contenders) == 0 && wait == false) {
                        return "", "", errors.New("No contenders found")
                } else {
                        if len(contenders) > 0 {
                                currentLeaderPath := fmt.Sprintf("%s/%s", l.electionPath, contenders[0])
                                bytes, _, ch, err := l.ZkSession.Conn.GetWatchWithBackoff(currentLeaderPath)
                                if err != nil {
                                        panic(err)
                                }
                                val := string(bytes[:])
                                go func() {
                                        <-ch
                                        if watch_cb != nil {
                                                watch_cb(currentLeaderPath, val)
                                        }
                                }()
                                return currentLeaderPath, val, nil
                        }
                }
                // If we are here means that there were no contenders and wait is True
                // Backoff exponentially till 32 seconds
                time.Sleep(backoff)
                if backoff < time.Duration(32)*time.Second {
                        backoff = backoff * 2
                }
        }
}

func (l *Leadership) StartLeaderChangeWatch(
        old_leader_path string, change_cb func(string)) {

        leader_path, leader_value, err := l.LeaderValue(true,
                func(leader_path, _ string) {
                        l.StartLeaderChangeWatch(leader_path, change_cb)
                })

        if err != nil {
                glog.Error("Could not get leader")
        }

        if leader_path != old_leader_path {
                change_cb(leader_value)
        }
}

func (l *Leadership) StartContendersChangeWatch(lockName string,
        changeCb func([]string)) {
        electionPath := fmt.Sprintf("%s/%s", GO_LEADERSHIP_ROOT, lockName)
        contendersChangeCb := func(conn *zk.Conn, data []string) {
                changeCb(data)
        }
        KeepWatchingChildrenOf(l.ZkSession.Conn, electionPath, contendersChangeCb)
}
