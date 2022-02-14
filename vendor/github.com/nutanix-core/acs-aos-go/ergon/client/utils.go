/*
 * Copyright (c) 2019 Nutanix Inc. All rights reserved.
 *
 * Author: vistaar.juneja@nutanix.com
 *
 * Ergon go client utilities.
 *
 */

package ergon_client

import (
	"encoding/base64"
	"github.com/nutanix-core/acs-aos-go/ergon"
	util_net "github.com/nutanix-core/acs-aos-go/nutanix/util-go/net"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"strconv"
	"time"
)

type Entity struct {
	EntityId         string `json:"entity_id,omitempty"`
	EntityType       string `json:"entity_type,omitempty"`
	LogicalTimestamp int64  `json:"logical_timestamp,omitempty"`
}

type Response struct {
	ErrorCode   int32  `json:"error_code,omitempty"`
	ErrorDetail string `json:"error_detail,omitempty"`
	Ret         string `json:"ret,omitempty"`
}

type RequestContext struct {
	UserUuid string `json:"user_uuid,omitempty"`
	Username string `json:"user_name,omitempty"`
	UserIp   string `json:"user_ip,omitempty"`
}

type ErgonTaskOutput struct {
	ParentTaskUuid            string         `json:"parent_task_uuid,omitempty"`
	InternalOpaque            string         `json:"internal_opaque,omitempty"`
	SubtaskSequenceId         uint64         `json:"subtask_sequence_id,omitempty"`
	DisableAutoProgressUpdate bool           `json:"disable_auto_progress_update,omitempty"`
	CapabilityList            []string       `json:"capability_list,omitempty"`
	DebugInfo                 string         `json:"debug_info,omitempty"`
	Reason                    string         `json:"reason,omitempty"`
	UseSyncBarrier            bool           `json:"use_sync_barrier,omitempty"`
	RequestContext            RequestContext `json:"request_context,omitempty"`
	InternalTask              bool           `json:"internal_task,omitempty"`
	LogicalTimestamp          int64          `json:"logical_timestamp,omitempty"`
	Weight                    uint64         `json:"weight,omitempty"`
	Canceled                  bool           `json:"canceled,omitempty"`
	ClusterUuid               string         `json:"cluster_uuid,omitempty"`
	Message                   string         `json:"message,omitempty"`
	Uuid                      string         `json:"uuid,omitempty"`
	StartTimeUsecs            time.Time      `json:"start_time_usecs,omitempty"`
	Response                  Response       `json:"response,omitempty"`
	LastUpdatedTimeUsecs      time.Time      `json:"last_updated_time_usecs,omitempty"`
	CreateTimeUsecs           time.Time      `json:"create_time_usecs,omitempty"`
	Status                    string         `json:"status,omitempty"`
	LocalRootTaskUuid         string         `json:"local_root_task_uuid,omitempty"`
	EntityList                []Entity       `json:"entity_list,omitempty"`
	Component                 string         `json:"component,omitempty"`
	SequenceId                uint64         `json:"sequence_id,omitempty"`
	RequestedStateTransition  string         `json:"requested_state_transition,omitempty"`
	CompleteTimeUsecs         time.Time      `json:"complete_time_usecs,omitempty"`
	PercentageComplete        int32          `json:"percentage_complete,omitempty"`
	Request                   string         `json:"request,omitempty"`
	OperationType             string         `json:"operation_type,omitempty"`
}

// Formats time as a human-friendly timestamp.
// Takes Unix epoch time as input. Returns nil if an error is encountered.
func FormatTime(epochTime int) time.Time {
	b := strconv.Itoa(epochTime)
	convertedTime, _ := strconv.ParseInt(b, 10, 64)
	return time.Unix(convertedTime, 0)
}

// Returns structure with type ErgonTaskOutput. This structure populates its
// values in human-friendly format from the Ergon task proto.
func SerializeTaskProto(task *ergon.Task) *ErgonTaskOutput {

	res := &ErgonTaskOutput{}

	// Populate ErgonTaskOutput from the task proto.

	res.InternalOpaque = base64.StdEncoding.EncodeToString(task.GetInternalOpaque())
	res.SubtaskSequenceId = task.GetSubtaskSequenceId()
	res.DisableAutoProgressUpdate = task.GetDisableAutoProgressUpdate()
	res.DebugInfo = task.GetDebugInfo()
	res.Reason = task.GetReason()
	res.UseSyncBarrier = task.GetUseSyncBarrier()
	res.Status = task.GetStatus().String()
	res.Component = task.GetComponent()
	res.SequenceId = task.GetSequenceId()
	res.RequestedStateTransition = task.GetRequestedStateTransition().String()
	res.InternalTask = task.GetInternalTask()
	res.LogicalTimestamp = task.GetLogicalTimestamp()
	res.Weight = task.GetWeight()
	res.Canceled = task.GetCanceled()
	res.Message = task.GetMessage()
	res.PercentageComplete = task.GetPercentageComplete()
	res.Request = base64.StdEncoding.EncodeToString(task.GetRequest().GetArg().GetEmbedded())
	res.OperationType = task.GetOperationType()

	res.ParentTaskUuid = uuid4.ToUuid4(task.GetParentTaskUuid()).String()
	res.Uuid = uuid4.ToUuid4(task.GetUuid()).String()
	res.ClusterUuid = uuid4.ToUuid4(task.GetClusterUuid()).String()
	res.LocalRootTaskUuid = uuid4.ToUuid4(task.GetLocalRootTaskUuid()).String()

	res.Response.ErrorCode = task.GetResponse().GetErrorCode()
	res.Response.ErrorDetail = task.GetResponse().GetErrorDetail()
	res.Response.Ret = base64.StdEncoding.EncodeToString(task.GetResponse().GetRet().GetEmbedded())

	res.RequestContext.Username = task.GetRequestContext().GetUserName()
	res.RequestContext.UserUuid = uuid4.ToUuid4(task.GetRequestContext().GetUserUuid()).String()
	res.RequestContext.UserIp = task.GetRequestContext().GetUserIp()

	res.StartTimeUsecs = FormatTime(int(task.GetStartTimeUsecs()) / 1000000)
	res.LastUpdatedTimeUsecs = FormatTime(int(task.GetLastUpdatedTimeUsecs()) / 1000000)
	res.CreateTimeUsecs = FormatTime(int(task.GetCreateTimeUsecs()) / 1000000)
	res.CompleteTimeUsecs = FormatTime(int(task.GetCompleteTimeUsecs()) / 1000000)

	for _, capability := range task.GetCapabilities() {
		res.CapabilityList = append(res.CapabilityList, string(capability))
	}

	for _, entity := range task.GetEntityList() {
		ent := Entity{}
		ent.EntityId = uuid4.ToUuid4(entity.GetEntityId()).String()
		ent.EntityType = entity.GetEntityType().String()
		ent.LogicalTimestamp = entity.GetLogicalTimestamp()
		res.EntityList = append(res.EntityList, ent)
	}

	return res

}

/*
	Provides an iterator based look up of Ergon Tasks.
	Example usage:
	import(
		ergonClient "github.com/nutanix-core/acs-aos-go/ergon/client"
		ergon_pr "github.com/nutanix-core/acs-aos-go/ergon"
	)
	taskgetarg := &ergon_pr.TaskGetArg{
		TaskUuidList: [][]byte{uuid1, uuid2, uuid3},
		... other fields as needed,
	}
	...
	iter := NewTaskGetIterator(ergon_client, taskGetArg, 100)
	for{
			taskGetRetBatch,_ := iter.Value()
			if taskGetRetBatch == nil{
				break
			}
			...
		}
		...

	Args:
	ergon_client: (ergonClient)
		Ergon client used to invoke RPCs on Ergon.
	taskGetArg: (TaskGetArg)
		Provides a sample TaskGetArg with the desired fields for the
		batch responses
	batch_size: (int64)
		Batch size of each response. The TaskGetArg provided in the arguments
		is split into batches of given size and response is returned.
		if batch_size == 1 : Returns 1 Task per Next call
*/
type TaskGetIterator struct {
	ergon_client  Ergon
	taskGetArg    *ergon.TaskGetArg
	taskUuidList  [][]byte
	next_iter_beg int64
	next_iter_end int64
	batch_size    int64
}

// NewTaskGetIterator creates a new TaskGet iterator
func NewTaskGetIterator(ergon_client Ergon, taskGetArg *ergon.TaskGetArg,
	batch_size int64) *TaskGetIterator {
	return &TaskGetIterator{
		ergon_client:  ergon_client,
		taskGetArg:    taskGetArg,
		taskUuidList:  taskGetArg.TaskUuidList,
		next_iter_beg: int64(0),
		next_iter_end: int64(0),
		batch_size:    batch_size,
	}
}

// Next advances to next batch.
// Returns the TaskGet response for the batch
// Returns nil on end of iteration.
func (iter *TaskGetIterator) Next() (*ergon.TaskGetRet, error) {
	// min returns the smaller of x and y of type int64
	min := func(x, y int64) int64 {
		if x > y {
			return y
		}
		return x
	}

	iter.next_iter_beg = iter.next_iter_end
	iter.next_iter_end = min(int64(len(iter.taskUuidList)),
		iter.next_iter_beg+iter.batch_size)
	iterable := (iter.next_iter_beg < int64(len(iter.taskUuidList)) &&
		iter.next_iter_end <= int64(len(iter.taskUuidList)))
	if !iterable {
		return nil, nil
	}

	taskuuidlist := iter.taskUuidList[iter.next_iter_beg:iter.next_iter_end]
	iter.taskGetArg.TaskUuidList = taskuuidlist
	taskGetRetBatch, err := iter.ergon_client.TaskGetMultiple(iter.taskGetArg)
	if err != nil {
		return nil, err
	}
	return taskGetRetBatch, nil
}

func IsRpcTimeoutError(err error) bool {
	rpcErr, ok := util_net.ExtractRpcError(err)
	return ok && util_net.ErrTimeout.Equals(rpcErr)
}
