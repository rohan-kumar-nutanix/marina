//
// Copyright (c) 2019 Nutanix Inc. All rights reserved.
//
// Author: raghu.rapole@nutanix.com
//
// Implements JobQueue related functionality.

package serial_executor

import (
	"errors"
	"sync"

	"github.com/golang/glog"
)

type jobQueue struct {
	sID         string   // Serialization ID for this queue.
	pendingJobs []JobIfc // List of jobs to be executed.
	active      bool     // Denotes if Job Queue is active. Queue becomes inactive before terminating.
	sync.Mutex           // To protect pendingJobs in this queue.
}

// NewJobQueue creates new jobQueue and returns its address.
func NewJobQueue(serialID string) *jobQueue {
	return &jobQueue{
		sID:         serialID,
		pendingJobs: []JobIfc{},
		active:      true,
	}
}

// Adds job to the queue. If the queue is not active, this method
// will return an error.
func (jobQ *jobQueue) addJob(job JobIfc) error {
	jobQ.Lock()
	defer jobQ.Unlock()
	if !jobQ.active {
		return errors.New("Job queue is not in active state")
	}
	jobQ.pendingJobs = append(jobQ.pendingJobs, job)
	return nil
}

// Executes job in the pending queue one after another. This method runs in
// goroutine and terminates once the queue becomes empty.
func (jobQ *jobQueue) processJobQueue(completionCB func(string)) {
	for {
		jobQ.Lock()
		// Check if pending jobs exists.
		if len(jobQ.pendingJobs) == 0 {
			// No jobs to process further..
			jobQ.active = false
			jobQ.Unlock()
			break
		}

		// Fetch and execute the first job in queue.
		job := jobQ.pendingJobs[0]
		jobQ.pendingJobs = jobQ.pendingJobs[1:]
		jobQ.Unlock()

		// Execute the job.
		glog.Infof("Executing job from queue with SID: %s", jobQ.sID)
		job.Execute()
		glog.Infof("Completed executing job from queue with SID: %s", jobQ.sID)
	}
	glog.V(2).Infof("Finished all jobs from job queue with SID: %s", jobQ.sID)
	completionCB(jobQ.sID)
}
