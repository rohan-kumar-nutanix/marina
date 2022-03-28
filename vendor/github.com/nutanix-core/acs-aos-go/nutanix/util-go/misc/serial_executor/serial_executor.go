//
// Copyright (c) 2019 Nutanix Inc. All rights reserved.
//
// Author: raghu.rapole@nutanix.com
//
// Implements jobs executor which serializes based on the serialisation identifier.
// The jobs are scheduled one after other with specific serialisation identifier
// at any given point. The order of execution for given serialisation identifier
// will be in the order of submission of jobs.
//
// Jobs executed by this serial executor must implement interface JobIfc.
//
// Usage:
//   se := NewSerialExecutor()
//   for _, job := range jobs {
//      se.SubmitJob(job)
//   }
//   se.WaitToFinish() // optional, if want to wait for completion of all jobs.
//
////////////////////////////////////////////////////////////////////////////////
// EXAMPLE
////////////////////////////////////////////////////////////////////////////////
// package test
//
// import (
// 	"fmt"
// 	"testing"
// 	"time"
// )
//
// type Task struct {
// 	uniqueId string
// 	<actual task>
// }
//
// func (t *Task) Execute() {
//  // Perform the operation.
// }
//
// func (t *Task) SerializationID() string {
// 	return t.uniqueId
// }
//
// func main() {
// 	se := NewSerialExecutor()
//
// 	for i := 0; i < 2; i++ {
// 		jobId := fmt.Sprintf("u1_j%d", i)
// 		se.SubmitJob(&Task{"u1", jobId})
// 	}
//
// 	for i := 0; i < 5; i++ {
// 		jobId := fmt.Sprintf("u2_j%d", i)
// 		se.SubmitJob(&Task{"u2", jobId})
// 	}
//
// 	se.WaitToFinish()
// }

package serial_executor

import (
	"sync"

	"github.com/golang/glog"
)

var gFatalf = glog.Fatalf

// JobIfc interface needs to be implemented by the jobs which are being processed
// using this serial executor.
type JobIfc interface {
	Execute()                // Perform the Job which is called after it got scheduled
	SerializationID() string // Should return unique Id based on which serialization is done
}

// SerialExecutorIfc interface implemented by serial executor.
type SerialExecutorIfc interface {
	SubmitJob(job JobIfc)
	WaitToFinish()
}

type serialExecutor struct {
	name         string               // Name of the executor, for debugging
	jobQueueMap  map[string]*jobQueue // serialId -> [job, job2..]
	completionWG *sync.WaitGroup      // To wait for completion of all jobs
	sync.RWMutex                      // To protect jobQueueMap
}

// NewSerialExecutor creates a new serial executor.
func NewSerialExecutor() *serialExecutor {
	return &serialExecutor{
		jobQueueMap:  make(map[string]*jobQueue),
		completionWG: &sync.WaitGroup{}}
}

// SubmitJob submits job for execution. This method will not actually execute the job,
// but adds job to queue which will be processed based on pending jobs in queue.
//
// If jobQueue for the SID is not present in `jobQueueMap`, it will create a new
// queue, add job to queue and spawns goroutine to process this queue.
//
// If SID for the job exists, it will append the new job to the corresponding job
// queue pointed by `jobsQueueMap`.
func (se *serialExecutor) SubmitJob(job JobIfc) {
	var createdNewJobQueue bool
	//var jobQ *jobQueue
	sID := job.SerializationID()

	for { // Loop until job is added to the queue.
		var jobQ *jobQueue
		var ok bool

		se.RLock()
		jobQ, ok = se.jobQueueMap[sID]
		se.RUnlock()

		if !ok {
			// Job queue not found for SID, take lock and create a new job queue
			// if not found.
			newJobQ := NewJobQueue(sID)

			se.Lock()
			jobQ, ok = se.jobQueueMap[sID]
			if !ok {
				jobQ = newJobQ
				se.jobQueueMap[sID] = jobQ
				createdNewJobQueue = true
			}
			se.Unlock()
		}

		// Add job to job queue.
		err := jobQ.addJob(job)
		if err != nil {
			// Failed to add job.
			continue
		}

		// If new Job Queue is created, process the same.
		var logMsg string
		if createdNewJobQueue {
			logMsg = "Created new job queue. "
			se.completionWG.Add(1)
			go jobQ.processJobQueue(se.jobQueueCompletionCB)
		}

		glog.Infof("[%s] %sAdded to the job queue.", jobQ.sID, logMsg)
		break
	}
}

func (se *serialExecutor) jobQueueCompletionCB(sID string) {
	se.Lock()
	if _, ok := se.jobQueueMap[sID]; !ok {
		// This should never happen.
		gFatalf("SID: %s is not found in job queue map", sID)
	}
	delete(se.jobQueueMap, sID)
	se.Unlock()

	se.completionWG.Done()
	glog.Infof("Done with job queue with SID: %s", sID)
}

// Waits until all the jobs in the queue are finished.
func (se *serialExecutor) WaitToFinish() {
	// Wait for the completion of all jobs.
	glog.Infof("Waiting for executor to process all jobs")
	se.completionWG.Wait()
	glog.Infof("Finished execution of all jobs in queue..")
}
