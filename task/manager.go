/*
* Copyright (c) 2022 Nutanix Inc. All rights reserved.
*
* Authors: rajesh.battala@nutanix.com
*
* This implements the Marina task manager. Basically it runs a task dispatcher
* that dispatches Marina tasks. When Marina is started, it queries Ergon for
* pending tasks, and dispatch them. The task dispatcher keeps polling Ergon for
* new Marina tasks to dispatch. So, when a new Ergon task is created for Marina,
* this task manager dispatcher will pick it up and execute it.
 */

package task

import (
	"sync"
	"time"

	log "k8s.io/klog/v2"

	"github.com/nutanix-core/acs-aos-go/ergon"
	ergonClient "github.com/nutanix-core/acs-aos-go/ergon/client"
	ergonTask "github.com/nutanix-core/acs-aos-go/ergon/task"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
	"github.com/nutanix-core/content-management-marina/grpc/catalog/catalog_item/tasks"
	"github.com/nutanix-core/content-management-marina/grpc/services"
	"github.com/nutanix-core/content-management-marina/interface/external"
	"github.com/nutanix-core/content-management-marina/proxy"
	"github.com/nutanix-core/content-management-marina/task/base"
	utils "github.com/nutanix-core/content-management-marina/util"
)

const (
	dispatcherRetryInitialDelaySecs = 1
	dispatcherRetryMaxDelaySecs     = 5
	dispatcherMaxRetries            = 20
)

var (
	taskManager      *MarinaTaskManager
	taskManagerOnce  sync.Once
	taskManagerMutex sync.Mutex
)

// MarinaTaskManager includes services to run the task manager.
type MarinaTaskManager struct {
	ergonTask.TaskMgrUtil // Implements StartDispatcher.
	stopDispatching       chan bool
}

// NewMarinaTaskManager creates a new Marina task manager singleton, and start
// task dispatcher.
func NewMarinaTaskManager() *MarinaTaskManager {
	taskManagerOnce.Do(func() {
		taskManager = &MarinaTaskManager{}
	})
	return taskManager
}

// StartTaskDispatcher start to run a task dispatcher. StartDispatcher() issues
// a TaskList RPC call to Ergon to get all pending tasks, and dispatch them.
// After that, it will keep on issuing TaskPoll RPC calls to Ergon to get newly
// created Marina tasks to dispatch.
func (m *MarinaTaskManager) StartTaskDispatcher() {
	taskManagerMutex.Lock()
	// Just make sure any previously still running task manager is stopped.
	if m.stopDispatching != nil {
		m.stopDispatching <- true
	}
	backoff := misc.NewExponentialBackoff(
		time.Duration(dispatcherRetryInitialDelaySecs)*time.Second,
		time.Duration(dispatcherRetryMaxDelaySecs)*time.Second,
		dispatcherMaxRetries)
	var stopChan chan bool
	var err error
	for {
		stopChan, err = m.StartDispatcher(m)
		if err == nil {
			log.Info("Started Marina task dispatcher successfully.")
			break
		}
		waited := backoff.Backoff()
		if waited == misc.Stop {
			log.Fatalf("Failed to start task dispatcher: %v.", err.Error())
		}
		log.Errorf("Failed to start task dispatcher: %v.", err.Error())
	}
	taskManagerMutex.Unlock()
	m.stopDispatching = stopChan
}

// StopTaskDispatcher stops the task dispatcher.
func (m *MarinaTaskManager) StopTaskDispatcher() {
	taskManagerMutex.Lock()
	defer taskManagerMutex.Unlock()

	if taskManager == nil {
		log.Error("Attempted to stop dispatcher when task manager has not been initialized.")
		return
	}
	taskManager.stopDispatching <- true
	taskManager.stopDispatching = nil
	log.Info("Stopped Marina task dispatcher.")
}

/*
 Implementation of the ergonTask.TaskManager interface: Ergon, Component, and
 Dispatch.
*/

// Ergon returns the ergon service.
func (m *MarinaTaskManager) Ergon() ergonClient.Ergon {
	return external.Interfaces().ErgonService()
}

// Component returns the Marina service name.
func (m *MarinaTaskManager) Component() string {
	return utils.ServiceName
}

// Dispatch is called with a list of dispatchable task. For each of those
// dispatchable tasks, we hydrate the task proto ("Task" as specified in
// ergon_type.proto), and dispatch it for execution by putting it in a task
// queue.
func (m *MarinaTaskManager) Dispatch(tasks []*ergon.Task) {
	for _, taskProto := range tasks {
		task := m.hydrateTask(taskProto)
		if task != nil {
			// Recover() calls RecoverHook() followed by Enqueue() to queue the task.
			task.Recover(task)
		}
	}
}

// hydrateTask returns a full Marina task based on the provided Ergon task proto.
func (m *MarinaTaskManager) hydrateTask(taskProto *ergon.Task) ergonTask.FullTask {
	if taskProto == nil {
		log.Error("Failed to hydrate a task with empty task proto.")
		return nil
	}
	methodName := taskProto.GetOperationType()
	log.Info("HydrateTask RPC methodName :", methodName)
	fullTask := services.GetTaskByRPC(tasks.NewCatalogItemBaseTask(base.NewMarinaBaseTask(taskProto)))
	if fullTask == nil {
		log.Warning("Proxying task to catalog service")
		fullTask = proxy.NewProxyTask(base.NewMarinaBaseTask(taskProto))
	}
	return fullTask
}
