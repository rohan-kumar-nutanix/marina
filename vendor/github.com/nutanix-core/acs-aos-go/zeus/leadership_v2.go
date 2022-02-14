/*
* Copyright (c) 2019 Nutanix Inc. All rights reserved.
*
* Authors: noufal.muhammed@nutanix.com
*
* This is the Golang implementation of 'leader election functionality using
* zookeeper'.
*
* Race with other language (C++, Python, etc) leader election implementations
* is avoided by creating a different namespace (election path) for Golang
* leader election.
*
* Example usage:
*
*	(1) Leader election in a distributed Golang service.
*
*		Each instance can choose to create a local leadership object.
*			leadership = NewGoLeadershipV2(<service name>, zkSession)
*
*		Each instance can choose to volunteer for leadership.
*			leadership.Volunteer(myInstanceId)
*
*		Each instance can choose to start leader change watch, which will notify
*			the leader change events.
*			leadership.StartLeaderChangeWatch( < old_leader = "" >, leaderChangeCb)
*
*		...
*
*		Each instance can choose to stop the leader-change-watch and remove the
*		leadership intention path from zookeeper as part of its graceful termination.
*			leadership.StopLeaderChangeWatch()
*			leadership.Relinquish()
*
*
*	(2) From a Golang service 'service-A', fetch the leader value of another
*		service 'service-B'. Here 'service-B' is not necessarily a Golang
*		service.
*
*		Create a local leadership object.
*			leadership = NewLeadershipV2("service-B", zkSession)
*			                     OR
*			leadership = NewPyLeadershipV2("service-B", zkSession)
*			                     OR
*			leadership = NewGoLeadershipV2("service-B", zkSession)
*
*		Fetch the current leader value.
*			leaderValue, err = leadership.FetchLeaderValue()
*			                     OR
*		Start a leader change watch, which will notify the leader change events.
*			leadership.StartLeaderChangeWatch( < old_leader = "" >, leaderChangeCb)
*
*		...
*
*		While graceful termination of the instance, stop the leader-change-watch if
*		it is running.
*			leadership.StopLeaderChangeWatch()
*
*
*	(3) Leader election in a distributed Golang service where instances will wait
*		for leadership.
*
*		Each instance will create a local leadership object.
*			leadership = NewGoLeadershipV2(<service name>, zkSession)
*
*		Each instance will volunteer for leadership.
*			leadership.Volunteer(myInstanceId)
*
*		Each instance will wait for leadership.
*			if leadership.WaitForLeadership() { ... }
*
*		...
*
*		Each instance will remove the leadership intention path from zookeeper as
*		part of its graceful termination.
*			leadership.Relinquish()
*
 */

package zeus

import (
	"flag"
	"fmt"
	"path"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/nutanix-core/acs-aos-go/go-zookeeper"
)

var (
	leadershipMaxNumRetries = flag.Uint("zeus_leadership_max_num_retries", 5,
		"Maximum number of retries in zeus leadership")

	leadershipMaxBackoffTime = flag.Uint("zeus_leadership_max_backoff_time", 32,
		"Maximum backoff time (seconds) in zeus leadership")
)

// A new leadership class needs to be created for every lock name/variable
// name which needs master election functionality from zookeeper.
type LeadershipV2 struct {
	ZkSession       *ZookeeperSession
	electionPath    string     // Eg: 'GO_LEADERSHIP_ROOT/castor'.
	intentionPath   string     // Eg: 'GO_LEADERSHIP_ROOT/castor/<ephemeral-sequential volunteer node>'.
	goLeadership    bool       // True if leadership zknode is being created on go leadership path.
	leaderWatchLock sync.Mutex // Lock for synchronization between leader change watch start and stop.
	watchStarted    bool       // True if leader-change-watch is running. False otherwise.
	stopWatchCh     chan bool  // Channel used to signal stop leader-change-watch request.
	ackStopWatchCh  chan bool  // Channel used to acknowledge that leader-change-watch is stopped.
}

// GO zookeeper library prepends all intention paths for leadership with
// "_c_" followed by GUID. This creates a problem when sorting paths to
// figure out a leader as guid do not follow strictly incrementing order.
// Solution is filter out this additional information before sorting.
// One peculiar thing is that different libraries behave differently.
// e.g. python library does not prepend guid, so code has to check if guid is
// present and then filter it out.
type LeadershipWithGUIDV2 []string

func (s LeadershipWithGUIDV2) Len() int {
	return len(s)
}
func (s LeadershipWithGUIDV2) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s LeadershipWithGUIDV2) Less(i, j int) bool {
	// Path is prefixed in Go zookeeper library as follows:
	// path = "protectedPrefix (_c_)" + guid (slice of 16) + proposed path
	first := s[i]
	second := s[j]

	// Logic here is : remove guid, protectedPrefix and then compare strings
	if strings.HasPrefix(s[i], PROTECTED_PREFIX) {
		first = s[i][len(PROTECTED_PREFIX)+33:]
	}
	if strings.HasPrefix(s[j], PROTECTED_PREFIX) {
		second = s[j][len(PROTECTED_PREFIX)+33:]
	}
	status := strings.Compare(first, second)
	if status > 0 {
		return false
	}
	return true
}

// This function will start a new Go leadership object and returns it.
func NewGoLeadershipV2(lockName string, zkSession *ZookeeperSession) (*LeadershipV2, error) {
	return newLeadershipWithPathPrefixV2(lockName, zkSession, GO_LEADERSHIP_ROOT)
}

// This function will start a new Python leadership object and returns it.
func NewPyLeadershipV2(lockName string, zkSession *ZookeeperSession) (*LeadershipV2, error) {
	return newLeadershipWithPathPrefixV2(lockName, zkSession, PY_LEADERSHIP_ROOT)
}

// This function will start a new common leadership object and returns it.
func NewLeadershipV2(lockName string, zkSession *ZookeeperSession) (*LeadershipV2, error) {
	return newLeadershipWithPathPrefixV2(lockName, zkSession, COMMON_LEADERSHIP_ROOT)
}

// This function will start a new leadership object on the given
// electionPathPrefix and returns it.
// If the electionPathPrefix is not GO_LEADERSHIP_ROOT, it is expected that
// electionPath zknode already exist in zookeeper.
func newLeadershipWithPathPrefixV2(lockName string, zkSession *ZookeeperSession, electionPathPrefix string) (*LeadershipV2, error) {

	if !zkSession.WaitForConnection() {
		return nil, ErrZkNotConnected
	}

	goLeadership := false
	if electionPathPrefix == GO_LEADERSHIP_ROOT {
		goLeadership = true
	}

	// Optimised for fast path: in most cases, electionPath zknode will be already
	// there in zookeeper. So check whether electionPath exist or not.
	electionPath := fmt.Sprintf("%s/%s", electionPathPrefix, lockName)
	exist, _, err := zkSession.Exist(electionPath, true /* retry */)
	if err != nil {
		glog.Fatalf("Unable to check whether zknode '%s' exist or not: %s",
			electionPath, err)
	}
	if exist {
		return &LeadershipV2{ZkSession: zkSession, electionPath: electionPath,
			goLeadership: goLeadership}, nil
	}

	// We are here means electionPath doesn't exist.
	// It is not supposed to create electionPath for other language services
	// from a Golang service. So if the leadership is not on GO_LEADERSHIP_ROOT,
	// fatal the service.
	if !goLeadership {
		glog.Fatalf("The zknode '%s' doesn't exist in zookeeper", electionPath)
	}

	// We are here means electionPath doesn't exist and it is a Golang service.

	_, err = zkSession.Create(electionPath, []byte{}, 0, WORLD_ACL, true /* retry */)
	// It is possible that by now some other client created the electionPath.
	// In that case error will be zk.ErrNodeExists.
	if err == nil || err == zk.ErrNodeExists {
		return &LeadershipV2{ZkSession: zkSession, electionPath: electionPath,
			goLeadership: goLeadership}, nil
	}
	// This error might be because GO_LEADERSHIP_ROOT zknode doesn't exist in
	// zookeeper. The GO_LEADERSHIP_ROOT zknode supposed to be created as part
	// of zookeeper initialization.
	glog.Fatalf("Unable to create election path zknode '%s' in zookeeper: %s",
		electionPath, err)

	return nil, err
}

// Volunteer for leadership on leader election path 'l.electionPath/<value>'.
// Note that it is not allowed to volunteer for leadership from a Go service
// to a non Go service.
func (l *LeadershipV2) Volunteer(value string) (bool, error) {
	if !l.goLeadership {
		glog.Fatalf("Golang service attempts to volunteer for leadership of "+
			"a non-Golang service '%s//%s'", l.electionPath, value)
		return false, ErrCouldNotVolunteer
	}

	if !l.ZkSession.WaitForConnection() {
		return false, ErrZkNotConnected
	}

	retry := uint(0)
	backoff := uint(1)
	var err error
	prefix := fmt.Sprintf("%s/n_", l.electionPath)
	for {
		l.intentionPath, err =
			l.ZkSession.Conn.CreateProtectedEphemeralSequential(prefix, []byte(value), WORLD_ACL)
		if err == nil {
			break
		}
		if err != zk.ErrNoNode {
			glog.Errorf("Unable to create leadership intention znode %s: %s",
				l.electionPath, err)
			return false, ErrCouldNotVolunteer
		}
		// We reached here means electionPath doesn't exist in zookeeper (ErrNoNode
		// error). This may be because electionPath is removed after creating it
		// (by calling CleanElectionPath()). Here create the election path again
		// and retry.
		glog.Warningf("Leadership election path '%s' doesn't exist in zookeeper", l.electionPath)
		_, err := l.ZkSession.Create(l.electionPath, []byte{}, 0, WORLD_ACL, true /* retry */)
		if err != nil && err != zk.ErrNodeExists {
			glog.Errorf("Unable to create leadership election path zknode %s: %s", l.electionPath, err)
			return false, ErrCouldNotVolunteer
		}
		maybeRetryAndFatal(&retry, &backoff, *leadershipMaxNumRetries, *leadershipMaxBackoffTime)
	}

	return true, nil
}

// This function will return the previous contender.
// If this node is the first contender, returns nil.
// Returns error in case of failure.
func (l *LeadershipV2) previousContender() (*string, error) {
	mySequence := path.Base(l.intentionPath)
	_, contenders, err := l.contendersW(false)
	if err != nil {
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

// This function will wait for leadership and returns true if leadership
// acquired. Returns false in case of failure.
func (l *LeadershipV2) WaitForLeadership() bool {
	if !l.goLeadership {
		glog.Fatalf("Golang service attempts to wait for leadership of a "+
			"non-Golang service '%s'", l.electionPath)
		return false
	}

	for {
		previous, err := l.previousContender()
		if err != nil {
			glog.Errorf("Error getting previous contender during leadership wait: %s", err)
			return false
		}
		if previous == nil {
			glog.Info("Claiming Leadership")
			return true
		}
		previousContenderPath := fmt.Sprintf("%s/%s", l.electionPath, *previous)
		exists, _, ch, err := l.ZkSession.ExistW(previousContenderPath, true /* retry */)
		if err != nil {
			glog.Fatal("Failed to set up watch on the previous contender")
		}
		if exists {
			<-ch
		}
	}
}

// This function will try to delete leadership intention path.
// Returns error in case of failure.
func (l *LeadershipV2) Relinquish() error {
	if !l.goLeadership {
		glog.Fatalf("Golang service attempts to relinquish leadership of a "+
			"non-Golang service '%s'", l.electionPath)
		return ErrNotAGoLeadership
	}

	err := l.ZkSession.Delete(l.intentionPath, -1, true /* retry */)
	if err != nil {
		glog.Errorf("Unable to delete zknode '%s' : %s", l.intentionPath, err)
		return err
	}
	return nil
}

// CleanElectionPath cleans up the path created for the lockName,
// if there are no intentionPaths or Volunteers for it.
// Returns error in case of cleanup failure.
func (l *LeadershipV2) CleanElectionPath() error {
	if !l.goLeadership {
		glog.Fatalf("Golang service attempts to remove election path of a "+
			"non-Golang service '%s'", l.electionPath)
		return ErrNotAGoLeadership
	}

	_, contenders, err := l.contendersW(false)
	if err != nil {
		glog.Errorf("Error getting contenders.ElectionPath '%s' cleanup "+
			"failed with %s", l.electionPath, err)
		return err
	}
	if len(contenders) == 0 {
		err = l.ZkSession.Delete(l.electionPath, -1, true /* retry */)
		if err != nil {
			glog.Errorf("ElectionPath %s cleanup failed with %s", l.electionPath, err)
			return err
		}
		glog.Infof("Election Path '%s' deleted", l.electionPath)
	}
	return nil
}

// This function will keep track of number of retries, backoff time and finally
// terminate the process if number of retries reaches the maximum limit.
// Arguments -
//	retry: Address of retry count variable.
//	backoff: Address of backoff time variable.
func maybeRetryAndFatal(retry *uint, backoff *uint, maxNumRetries uint, maxBackoffTime uint) {
	// Fatal if the number of retries exceeds the maximum.
	if *retry >= maxNumRetries {
		glog.Fatal("Number of retries exceeded the maximum value ",
			maxNumRetries)
	}
	*retry++

	// Retry with exponential backoff.
	time.Sleep(time.Duration(*backoff) * time.Second)
	*backoff = *backoff * 2
	if *backoff > maxBackoffTime {
		*backoff = maxBackoffTime
	}
}

// Function returns the sorted list of zknodes present at the electionPath.
// If wait is set to true, then this function will get blocked until there are
// contenders in the electionPath.
// Arguments -
//	Wait: whether to wait for the leader to be present.
// Returns -
//	Contenders list change watch.
//	Contenders list.
//	Error.
func (l *LeadershipV2) contendersW(wait bool) (<-chan zk.Event, []string, error) {
	for {
		children, _, ch, err := l.ZkSession.Conn.ChildrenW(l.electionPath)
		if err != nil {
			return nil, nil, err
		}
		if len(children) > 0 {
			sort.Sort(LeadershipWithGUIDV2(children))
			return ch, children, nil
		}
		if wait {
			<-ch
		} else {
			return ch, children, nil
		}
	}
}

// Function to fetch leader value from Zookeeper election path.
// Arguments -
// Returns -
//	Leader value.
//	Error.
func (l *LeadershipV2) FetchLeaderValue() (string, error) {
	// Read the sorted list of children nodes present at election path. As per
	// leader election algorithm, the first children in the list is the leader.
	// The value stored in the leader znode is the leader value.

	retry := uint(0)
	backoff := uint(1)
	for {
		_, contenders, err := l.contendersW(false /* wait */)
		if err != nil {
			glog.Error("Reading contenders from Zookeeper election path ",
				l.electionPath, " returned error ", err)
			return "", err
		}
		if len(contenders) == 0 {
			return "", nil
		}

		currentLeaderPath := fmt.Sprintf("%s/%s", l.electionPath, contenders[0])
		bytes, _, err := l.ZkSession.Get(currentLeaderPath, false /* retry */)
		if err == zk.ErrNoNode {
			glog.Infof("Reading leader value from Zookeeper leader path '%s' returned error %s",
				currentLeaderPath, err)
			// It is possible that 'currentLeaderPath' node went down after reading the contenders list.
			// To address this case, here we re-read contenders list and try to get leader value again.
			// Theoretically it is possible that this scenario continue happening infinitely. But in
			// practice, if nodes are repeatedly going down this fast, then it is safe to assume that
			// there have some other serious issue in the cluster. So here we do fatal once retry count
			// exceeds its maximum.
			maybeRetryAndFatal(&retry, &backoff, *leadershipMaxNumRetries, *leadershipMaxBackoffTime)
			continue
		}
		if err != nil {
			glog.Error("Reading leader value from Zookeeper leader path ", currentLeaderPath,
				" returned error ", err)
			return "", err
		}
		val := string(bytes[:])
		return val, nil
	}
}

// This function runs an infinite loop which will set a zookeeper watch on the
// first contender and then wait for the watch event. Upon receiving the event,
// the callback routine will get called if the current leader is different from
// the older one. Simultaneously, this infinite loop will also wait for stop
// leader-change-watch request from the consumer. Upon receiving the stop
// request, this function will return from the execution.
// Arguments -
//	oldLeaderValue: The current known leader value to the consumer class. This
//					can be an empty string.
//                  Callback will be invoked if the new leader differs from the
//                  current leader.
//	changeCb: The callback invoked on a leader change.
func (l *LeadershipV2) leaderChangeWatch(oldLeaderValue string, changeCb func(string)) {

	retry := uint(0)
	backoff := uint(1)
	for {
		contendersWaitCh, contenders, err := l.contendersW(false /* wait */)
		if err != nil {
			glog.Error("Reading contenders from zookeeper election path ", l.electionPath,
				" returned error ", err)
			maybeRetryAndFatal(&retry, &backoff, *leadershipMaxNumRetries, *leadershipMaxBackoffTime)
			continue
		}
		if len(contenders) == 0 {
			// There is no leader now. Notify consumer about this, if not done already.
			if oldLeaderValue != "" {
				oldLeaderValue = ""
				changeCb("" /* leaderValue */)
			}
			// Simultaneously wait for contender list change event and user stop event.
			// For contender list change event, re-read the contenders list and try again.
			select {
			case _, ok := <-contendersWaitCh:
				if !ok {
					glog.Fatal("Zookeeper leader path watch channel closed unexpecteadily")
				}
			case _, ok := <-l.stopWatchCh:
				if !ok {
					glog.Fatal("Stop watch channel of Zk leader change watch closed")
				}
				glog.Info("Zeus leader change watch is ending")
				l.ackStopWatchCh <- true
				return
			}
			// Re-read contenders list and try again.
			continue
		}

		currentLeaderPath := fmt.Sprintf("%s/%s", l.electionPath, contenders[0])
		glog.V(2).Info("Current leader path is ", currentLeaderPath)

		// Get the value of zknode at 'currentLeaderPath' and set a watch on
		// content change.
		bytes, _, zkWatchCh, err := l.ZkSession.GetW(currentLeaderPath, false /* retry */)
		if err != nil {
			glog.Errorf("Unable to read values from zknode '%s' : '%s' ", currentLeaderPath, err)
			maybeRetryAndFatal(&retry, &backoff, *leadershipMaxNumRetries, *leadershipMaxBackoffTime)
			continue
		}

		leaderValue := string(bytes[:])
		// The leader has change from the last known leader.
		if oldLeaderValue != leaderValue {
			oldLeaderValue = leaderValue
			glog.V(2).Info("Invoking leader change callback with value ", leaderValue)
			// Multiple leader change events might have happened while this
			// goroutine is executing the callback function. In such cases we
			// will only return latest leader value. That is okay since only the
			// latest leader value is important in leader change notification
			// use case.
			changeCb(leaderValue)
		}

		// Waiting for events. Watch event from Zk or stop request event from consumer.
		select {

		// Here it is possible that both zkWatchCh and stopWatchCh have channel
		// message receive event simultaneously. In such scenario 'select' will
		// randomly choose either one of the event first. So here it is possible
		// that leader-change-watch get terminated without notifying pending
		// leader change event.

		case _, ok := <-zkWatchCh:
			if !ok {
				glog.Fatal("Watch channel of Zk leadership change watch closed")
			}
		case _, ok := <-l.stopWatchCh:
			if !ok {
				// We shouldn't have reached here. This channel will not get
				// closed at this point.
				glog.Fatal("Stop watch channel of Zk leader change watch closed")
			}
			glog.Info("Zk leader change watch is ending")
			l.ackStopWatchCh <- true
			return
		}
		// Reset retry count and backoff time to initial values.
		retry = 0
		backoff = 1
	}
}

// Function to start a watch over leadership changes.
// The leadership change events will be notified by invoking callback function
// 'changeCb' with leader value as argument. If the leadership changes multiple
// times in quick succession, callbacks may not be raised for all intermediate
// changes. However, the last leader change will be notified always.
// If there is no contender for leadership, an empty string will be returned
// as leader value.
// Arguments -
//	oldLeaderValue: The current known leader value to the consumer class. This
//					can be an empty string.
//                  Callback will be invoked if the new leader differs from the
//                  current leader.
//	changeCb: The callback invoked on a leader change event.
func (l *LeadershipV2) StartLeaderChangeWatch(oldLeaderValue string, changeCb func(string)) {

	if changeCb == nil {
		glog.Fatal("Callback function pointer passed for start leader-change-watch API is nil")
	}

	l.leaderWatchLock.Lock()

	// Only one leader-change-watch in a leadership session.
	if l.watchStarted {
		l.leaderWatchLock.Unlock()
		glog.Fatal("Ignoring the request to start Zk leader change watch when it's already started")
		return
	}

	l.stopWatchCh = make(chan bool)
	l.ackStopWatchCh = make(chan bool)
	l.watchStarted = true

	l.leaderWatchLock.Unlock()

	go l.leaderChangeWatch(oldLeaderValue, changeCb)
}

// Function to stop leadership change watch.
// Users may keep getting leader-change callback until this function returns.
func (l *LeadershipV2) StopLeaderChangeWatch() {
	glog.Info("Request received to stop Zk leader-change-watch")

	l.leaderWatchLock.Lock()

	if !l.watchStarted {
		glog.Warning("Ignoring the request to stop Zk leader change watch when it's not started")
		l.leaderWatchLock.Unlock()
		return
	}
	l.stopWatchCh <- true
	glog.Info("Sent channel message to stop Zk leader-change-watch")
	<-l.ackStopWatchCh
	close(l.stopWatchCh)
	close(l.ackStopWatchCh)
	l.watchStarted = false
	glog.Info("Zk leader-change-watch stopped successfully")

	// To reduce the length of critical section, if we do unlock in the beginning
	// of this function (by also moving assignment statement watchStarted=false
	// to the beginning), then a start watch request from user can lead to start
	// of a new watch before terminating the existing watch. Statements like
	// <-l.ackStopWatchCh and channel close might get executed after starting
	// second watch, and that can lead to program crash. To avoid that, we do
	// unlock only after closing the watch completely.
	l.leaderWatchLock.Unlock()
}
