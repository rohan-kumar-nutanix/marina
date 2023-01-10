/*
* Copyright (c) 2023 Nutanix Inc. All rights reserved.
*
* Author: rishabh.gupta@nutanix.com
 */

package tasks

import (
	"errors"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/nutanix-core/acs-aos-go/ergon"
	mockTask "github.com/nutanix-core/acs-aos-go/ergon/task/mocks"
	"github.com/nutanix-core/acs-aos-go/insights/insights_interface"
	"github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
	"github.com/nutanix-core/acs-aos-go/zeus"
	marinaError "github.com/nutanix-core/content-management-marina/errors"
	mockBase "github.com/nutanix-core/content-management-marina/mocks/task/base"
	mockClient "github.com/nutanix-core/content-management-marina/mocks/util/catalog/client"
	marinaIfc "github.com/nutanix-core/content-management-marina/protos/marina"
	"github.com/nutanix-core/content-management-marina/task/base"
	utils "github.com/nutanix-core/content-management-marina/util"
	catalogClient "github.com/nutanix-core/content-management-marina/util/catalog/client"
	marinaZeus "github.com/nutanix-core/content-management-marina/zeus"
)

type MockRemoteCatalogClient struct {
	client *mockClient.RemoteCatalogInterface
}

func (mock *MockRemoteCatalogClient) NewRemoteCatalogService(zkSession *zeus.ZookeeperSession,
	zeusConfig marinaZeus.ConfigCache, clusterUUID, userUUID, tenantUUID *uuid4.Uuid,
	userName *string) catalogClient.RemoteCatalogInterface {
	return mock.client
}

func TestNewCatalogItemDeleteTask(t *testing.T) {
	catalogItemBaseTask := &CatalogItemBaseTask{}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	assert.Equal(t, catalogItemBaseTask, catalogItemDeleteTask.CatalogItemBaseTask)
}

func TestCatalogItemDeleteStartHook(t *testing.T) {
	mockBaseTask := &mockBase.MarinaBaseTaskInterface{}
	baseTask := &base.MarinaBaseTask{MarinaBaseTaskInterface: mockBaseTask}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	arg := &marinaIfc.CatalogItemDeleteArg{
		CatalogItemId: &marinaIfc.CatalogItemId{
			GlobalCatalogItemUuid: testCatalogItemUuid.RawBytes(),
		},
	}
	embedded, _ := proto.Marshal(arg)
	payload := &ergon.PayloadOrEmbeddedValue{Embedded: embedded}
	req := &ergon.MetaRequest{Arg: payload}
	task := &ergon.Task{Request: req}
	mockBaseTask.On("Proto").Return(task).Once()
	wal := &marinaIfc.PcTaskWalRecord{}
	mockBaseTask.On("Wal").Return(wal).Once()
	mockBaseTask.On("SetWal", mock.Anything).Return(nil).Once()

	err := catalogItemDeleteTask.StartHook()

	assert.NoError(t, err)
	assert.True(t, proto.Equal(arg, catalogItemDeleteTask.arg))
	assert.Equal(t, testCatalogItemUuid, catalogItemDeleteTask.globalCatalogItemUuid)
	assert.NotNil(t, wal.GetData().GetCatalogItem())
	mockBaseTask.AssertExpectations(t)
}

func TestCatalogItemDeleteStartHookArgError(t *testing.T) {
	mockBaseTask := &mockBase.MarinaBaseTaskInterface{}
	baseTask := &base.MarinaBaseTask{MarinaBaseTaskInterface: mockBaseTask}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	payload := &ergon.PayloadOrEmbeddedValue{Embedded: []byte("{foobar")}
	req := &ergon.MetaRequest{Arg: payload}
	task := &ergon.Task{Request: req}
	mockBaseTask.On("Proto").Return(task).Once()

	err := catalogItemDeleteTask.StartHook()

	assert.Error(t, err)
	assert.IsType(t, marinaError.ErrInvalidArgument, err)
}

func TestCatalogItemDeleteStartHookUuidError(t *testing.T) {
	mockBaseTask := &mockBase.MarinaBaseTaskInterface{}
	baseTask := &base.MarinaBaseTask{MarinaBaseTaskInterface: mockBaseTask}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	arg := &marinaIfc.CatalogItemDeleteArg{}
	embedded, _ := proto.Marshal(arg)
	payload := &ergon.PayloadOrEmbeddedValue{Embedded: embedded}
	req := &ergon.MetaRequest{Arg: payload}
	task := &ergon.Task{Request: req}
	mockBaseTask.On("Proto").Return(task).Once()

	err := catalogItemDeleteTask.StartHook()

	assert.Error(t, err)
	assert.IsType(t, marinaError.ErrInvalidArgument, err)
}

func TestCatalogItemDeleteRecoverHook(t *testing.T) {
	mockBaseTask := &mockBase.MarinaBaseTaskInterface{}
	baseTask := &base.MarinaBaseTask{MarinaBaseTaskInterface: mockBaseTask}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	arg := &marinaIfc.CatalogItemDeleteArg{
		CatalogItemId: &marinaIfc.CatalogItemId{
			GlobalCatalogItemUuid: testCatalogItemUuid.RawBytes(),
		},
	}
	embedded, _ := proto.Marshal(arg)
	payload := &ergon.PayloadOrEmbeddedValue{Embedded: embedded}
	req := &ergon.MetaRequest{Arg: payload}
	task := &ergon.Task{Request: req}
	mockBaseTask.On("Proto").Return(task).Once()

	err := catalogItemDeleteTask.RecoverHook()

	assert.NoError(t, err)
	assert.True(t, proto.Equal(arg, catalogItemDeleteTask.arg))
	assert.Equal(t, testCatalogItemUuid, catalogItemDeleteTask.globalCatalogItemUuid)
	mockBaseTask.AssertExpectations(t)
}

func TestCatalogItemDeleteRecoverHookArgError(t *testing.T) {
	mockBaseTask := &mockBase.MarinaBaseTaskInterface{}
	baseTask := &base.MarinaBaseTask{MarinaBaseTaskInterface: mockBaseTask}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	payload := &ergon.PayloadOrEmbeddedValue{Embedded: []byte("{foobar")}
	req := &ergon.MetaRequest{Arg: payload}
	task := &ergon.Task{Request: req}
	mockBaseTask.On("Proto").Return(task).Once()

	err := catalogItemDeleteTask.RecoverHook()

	assert.Error(t, err)
	assert.IsType(t, marinaError.ErrInvalidArgument, err)
}

func TestCatalogItemDeleteRecoverHookUuidError(t *testing.T) {
	mockBaseTask := &mockBase.MarinaBaseTaskInterface{}
	baseTask := &base.MarinaBaseTask{MarinaBaseTaskInterface: mockBaseTask}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	arg := &marinaIfc.CatalogItemDeleteArg{}
	embedded, _ := proto.Marshal(arg)
	payload := &ergon.PayloadOrEmbeddedValue{Embedded: embedded}
	req := &ergon.MetaRequest{Arg: payload}
	task := &ergon.Task{Request: req}
	mockBaseTask.On("Proto").Return(task).Once()

	err := catalogItemDeleteTask.RecoverHook()

	assert.Error(t, err)
	assert.IsType(t, marinaError.ErrInvalidArgument, err)
}

func TestCatalogItemDeleteEnqueue(t *testing.T) {
	mockExtInterfaces := mockExternalInterfaces()
	//mockIntInterfaces := mockInternalInterfaces()
	baseTask := &base.MarinaBaseTask{
		ExternalSingletonInterface: mockExtInterfaces,
		//InternalSingletonInterface: mockIntInterfaces,
	}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	serialExecutorIfc.On("SubmitJob", mock.Anything).Return().Once()

	catalogItemDeleteTask.Enqueue()

	serialExecutorIfc.AssertExpectations(t)
}

func TestCatalogItemDeleteExecute(t *testing.T) {
	mockBaseTask := &mockTask.BaseTask{}
	baseTask := &base.MarinaBaseTask{BaseTask: mockBaseTask}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	mockBaseTask.On("Resume", mock.Anything).Return().Once()

	catalogItemDeleteTask.Execute()

	mockBaseTask.AssertExpectations(t)
}

func TestCatalogItemDeleteSerializationID(t *testing.T) {
	catalogItemBaseTask := &CatalogItemBaseTask{globalCatalogItemUuid: testCatalogItemUuid}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)

	serialisationId := catalogItemDeleteTask.SerializationID()

	assert.Equal(t, testCatalogItemUuid.String(), serialisationId)
}

func TestCatalogItemDeleteRun(t *testing.T) {
	mockErgonBaseTask := &mockTask.BaseTask{}
	mockBaseTask := &mockBase.MarinaBaseTaskInterface{}
	mockExtInterfaces := mockExternalInterfaces()
	mockIntInterfaces := mockInternalInterfaces()
	baseTask := &base.MarinaBaseTask{
		BaseTask:                   mockErgonBaseTask,
		MarinaBaseTaskInterface:    mockBaseTask,
		ExternalSingletonInterface: mockExtInterfaces,
		InternalSingletonInterface: mockIntInterfaces,
	}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	catalogItemDeleteTask.arg = &marinaIfc.CatalogItemDeleteArg{
		CatalogItemId: &marinaIfc.CatalogItemId{
			Version: &testCatalogItemVersion,
		},
	}
	pcUuid, _ := uuid4.New()
	peUuid, _ := uuid4.New()
	remoteTaskUuid, _ := uuid4.New()
	wal := &marinaIfc.PcTaskWalRecord{
		Data: &marinaIfc.PcTaskWalRecordData{
			CatalogItem: &marinaIfc.PcTaskWalRecordCatalogItemData{},
			TaskList: []*marinaIfc.EndpointTask{{
				ClusterUuid: pcUuid.RawBytes(),
				TaskUuid:    remoteTaskUuid.RawBytes(),
			}},
		},
	}
	mockBaseTask.On("Wal").Return(wal).Once()
	configIfc.On("PeClusterUuids").Return([]*uuid4.Uuid{peUuid}).Once()
	mockBaseTask.On("SetWal", mock.Anything).Return(nil).Once()
	mockErgonBaseTask.On("Save", mock.Anything).Return(nil).Once()
	configIfc.On("ClusterUuid").Return(pcUuid).Once()
	catalogItemUuid := testCatalogItemUuid.String()
	entities := []*insights_interface.EntityWithMetric{
		{
			EntityGuid: &insights_interface.EntityGuid{
				EntityId: &catalogItemUuid,
			},
		},
	}
	mockBaseTask.On("Wal").Return(wal).Once()
	uuidIfc.On("New").Return(remoteTaskUuid, nil).Once()
	mockBaseTask.On("SetWal", mock.Anything).Return(nil).Once()
	mockErgonBaseTask.On("Save", mock.Anything).Return(nil).Once()
	taskUuid, _ := uuid4.New()
	mockBaseTask.On("Proto").Return(&ergon.Task{Uuid: taskUuid.RawBytes()}).Once()
	mockRemoteClient := &mockClient.RemoteCatalogInterface{}
	mockRemoteClient.On("SendMsg", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	mockRemoteCatalogClient := &MockRemoteCatalogClient{client: mockRemoteClient}
	remoteCatalogService = mockRemoteCatalogClient.NewRemoteCatalogService
	taskMap := make(map[utils.RemoteEndpoint]*ergon.Task)
	fanoutIfc.On("PollAllRemoteTasks", mock.Anything, mock.Anything).Return(&taskMap).Once()
	cpdbIfc.On("Query", mock.Anything).Return(entities, nil).Once()
	idfIfc.On("DeleteEntities", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()
	marshal, _ := proto.Marshal(&marinaIfc.CatalogItemDeleteRet{})
	protoIfc.On("Marshal", mock.Anything).Return(marshal, nil).Once()
	mockErgonBaseTask.On("Complete", mock.Anything, mock.Anything, mock.Anything).Return().Once()

	err := catalogItemDeleteTask.Run()

	assert.NoError(t, err)
	assert.Equal(t, [][]byte{peUuid.RawBytes()}, wal.GetData().GetCatalogItem().GetClusterUuidList())
	mockErgonBaseTask.AssertExpectations(t)
	mockBaseTask.AssertExpectations(t)
	configIfc.AssertExpectations(t)
	cpdbIfc.AssertExpectations(t)
	idfIfc.AssertExpectations(t)
	uuidIfc.AssertExpectations(t)
	mockRemoteClient.AssertExpectations(t)
	fanoutIfc.AssertExpectations(t)
	protoIfc.AssertExpectations(t)
}

func TestCatalogItemDeleteRunWalError(t *testing.T) {
	mockBaseTask := &mockBase.MarinaBaseTaskInterface{}
	mockExtInterfaces := mockExternalInterfaces()
	mockIntInterfaces := mockInternalInterfaces()
	baseTask := &base.MarinaBaseTask{
		MarinaBaseTaskInterface:    mockBaseTask,
		ExternalSingletonInterface: mockExtInterfaces,
		InternalSingletonInterface: mockIntInterfaces,
	}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	catalogItemDeleteTask.arg = &marinaIfc.CatalogItemDeleteArg{
		CatalogItemId: &marinaIfc.CatalogItemId{
			Version: &testCatalogItemVersion,
		},
	}
	wal := &marinaIfc.PcTaskWalRecord{
		Data: &marinaIfc.PcTaskWalRecordData{
			CatalogItem: &marinaIfc.PcTaskWalRecordCatalogItemData{},
		},
	}
	mockBaseTask.On("Wal").Return(wal).Once()
	peUuid, _ := uuid4.New()
	configIfc.On("PeClusterUuids").Return([]*uuid4.Uuid{peUuid}).Once()
	mockBaseTask.On("SetWal", mock.Anything).Return(errors.New("oh no")).Once()

	err := catalogItemDeleteTask.Run()

	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
	assert.Equal(t, [][]byte{peUuid.RawBytes()}, wal.GetData().GetCatalogItem().GetClusterUuidList())
	mockBaseTask.AssertExpectations(t)
	configIfc.AssertExpectations(t)
}

func TestCatalogItemDeleteRunIdfError(t *testing.T) {
	mockErgonBaseTask := &mockTask.BaseTask{}
	mockBaseTask := &mockBase.MarinaBaseTaskInterface{}
	mockExtInterfaces := mockExternalInterfaces()
	mockIntInterfaces := mockInternalInterfaces()
	baseTask := &base.MarinaBaseTask{
		BaseTask:                   mockErgonBaseTask,
		MarinaBaseTaskInterface:    mockBaseTask,
		ExternalSingletonInterface: mockExtInterfaces,
		InternalSingletonInterface: mockIntInterfaces,
	}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	catalogItemDeleteTask.arg = &marinaIfc.CatalogItemDeleteArg{
		CatalogItemId: &marinaIfc.CatalogItemId{
			Version: &testCatalogItemVersion,
		},
	}
	pcUuid, _ := uuid4.New()
	peUuid, _ := uuid4.New()
	remoteTaskUuid, _ := uuid4.New()
	wal := &marinaIfc.PcTaskWalRecord{
		Data: &marinaIfc.PcTaskWalRecordData{
			CatalogItem: &marinaIfc.PcTaskWalRecordCatalogItemData{
				ClusterUuidList: [][]byte{peUuid.RawBytes()},
			},
			TaskList: []*marinaIfc.EndpointTask{{
				ClusterUuid: pcUuid.RawBytes(),
				TaskUuid:    remoteTaskUuid.RawBytes(),
			}},
		},
	}
	mockBaseTask.On("Wal").Return(wal).Once()
	configIfc.On("ClusterUuid").Return(pcUuid).Once()
	mockBaseTask.On("Wal").Return(wal).Once()
	uuidIfc.On("New").Return(remoteTaskUuid, nil).Once()
	mockBaseTask.On("SetWal", mock.Anything).Return(nil).Once()
	mockErgonBaseTask.On("Save", mock.Anything).Return(nil).Once()
	taskUuid, _ := uuid4.New()
	mockBaseTask.On("Proto").Return(&ergon.Task{Uuid: taskUuid.RawBytes()}).Once()
	mockRemoteClient := &mockClient.RemoteCatalogInterface{}
	mockRemoteClient.On("SendMsg", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	mockRemoteCatalogClient := &MockRemoteCatalogClient{client: mockRemoteClient}
	remoteCatalogService = mockRemoteCatalogClient.NewRemoteCatalogService
	taskMap := make(map[utils.RemoteEndpoint]*ergon.Task)
	fanoutIfc.On("PollAllRemoteTasks", mock.Anything, mock.Anything).Return(&taskMap).Once()
	cpdbIfc.On("Query", mock.Anything).Return(nil, errors.New("oh no")).Once()

	err := catalogItemDeleteTask.Run()

	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
	assert.Equal(t, [][]byte{peUuid.RawBytes()}, wal.GetData().GetCatalogItem().GetClusterUuidList())
	mockBaseTask.AssertExpectations(t)
	configIfc.AssertExpectations(t)
	cpdbIfc.AssertExpectations(t)
	uuidIfc.AssertExpectations(t)
	mockRemoteClient.AssertExpectations(t)
	fanoutIfc.AssertExpectations(t)
	mockErgonBaseTask.AssertExpectations(t)
}

func TestCatalogItemDeleteRunUuidError(t *testing.T) {
	mockErgonBaseTask := &mockTask.BaseTask{}
	mockBaseTask := &mockBase.MarinaBaseTaskInterface{}
	mockExtInterfaces := mockExternalInterfaces()
	mockIntInterfaces := mockInternalInterfaces()
	baseTask := &base.MarinaBaseTask{
		BaseTask:                   mockErgonBaseTask,
		MarinaBaseTaskInterface:    mockBaseTask,
		ExternalSingletonInterface: mockExtInterfaces,
		InternalSingletonInterface: mockIntInterfaces,
	}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	catalogItemDeleteTask.arg = &marinaIfc.CatalogItemDeleteArg{
		CatalogItemId: &marinaIfc.CatalogItemId{
			Version: &testCatalogItemVersion,
		},
	}
	pcUuid, _ := uuid4.New()
	peUuid, _ := uuid4.New()
	remoteTaskUuid, _ := uuid4.New()
	wal := &marinaIfc.PcTaskWalRecord{
		Data: &marinaIfc.PcTaskWalRecordData{
			CatalogItem: &marinaIfc.PcTaskWalRecordCatalogItemData{},
			TaskList: []*marinaIfc.EndpointTask{{
				ClusterUuid: pcUuid.RawBytes(),
				TaskUuid:    remoteTaskUuid.RawBytes(),
			}},
		},
	}
	mockBaseTask.On("Wal").Return(wal).Once()
	configIfc.On("PeClusterUuids").Return([]*uuid4.Uuid{peUuid}).Once()
	mockBaseTask.On("SetWal", mock.Anything).Return(nil).Once()
	mockErgonBaseTask.On("Save", mock.Anything).Return(nil).Once()
	configIfc.On("ClusterUuid").Return(pcUuid).Once()
	mockBaseTask.On("Wal").Return(wal).Once()
	uuidIfc.On("New").Return(nil, marinaError.ErrInternalError()).Once()
	mockRemoteClient := &mockClient.RemoteCatalogInterface{}
	mockRemoteCatalogClient := &MockRemoteCatalogClient{client: mockRemoteClient}
	remoteCatalogService = mockRemoteCatalogClient.NewRemoteCatalogService

	err := catalogItemDeleteTask.Run()

	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
	assert.Equal(t, [][]byte{peUuid.RawBytes()}, wal.GetData().GetCatalogItem().GetClusterUuidList())
	mockErgonBaseTask.AssertExpectations(t)
	mockBaseTask.AssertExpectations(t)
	configIfc.AssertExpectations(t)
	uuidIfc.AssertExpectations(t)
	mockRemoteClient.AssertExpectations(t)
}

func TestCatalogItemDeleteRunFanoutError(t *testing.T) {
	mockErgonBaseTask := &mockTask.BaseTask{}
	mockBaseTask := &mockBase.MarinaBaseTaskInterface{}
	mockExtInterfaces := mockExternalInterfaces()
	mockIntInterfaces := mockInternalInterfaces()
	baseTask := &base.MarinaBaseTask{
		BaseTask:                   mockErgonBaseTask,
		MarinaBaseTaskInterface:    mockBaseTask,
		ExternalSingletonInterface: mockExtInterfaces,
		InternalSingletonInterface: mockIntInterfaces,
	}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	catalogItemDeleteTask.arg = &marinaIfc.CatalogItemDeleteArg{
		CatalogItemId: &marinaIfc.CatalogItemId{
			Version: &testCatalogItemVersion,
		},
	}
	pcUuid, _ := uuid4.New()
	peUuid, _ := uuid4.New()
	remoteTaskUuid, _ := uuid4.New()
	wal := &marinaIfc.PcTaskWalRecord{
		Data: &marinaIfc.PcTaskWalRecordData{
			CatalogItem: &marinaIfc.PcTaskWalRecordCatalogItemData{},
			TaskList: []*marinaIfc.EndpointTask{{
				ClusterUuid: pcUuid.RawBytes(),
				TaskUuid:    remoteTaskUuid.RawBytes(),
			}},
		},
	}
	mockBaseTask.On("Wal").Return(wal).Once()
	configIfc.On("PeClusterUuids").Return([]*uuid4.Uuid{peUuid}).Once()
	mockBaseTask.On("SetWal", mock.Anything).Return(nil).Once()
	mockErgonBaseTask.On("Save", mock.Anything).Return(nil).Once()
	configIfc.On("ClusterUuid").Return(pcUuid).Once()
	mockBaseTask.On("Wal").Return(wal).Once()
	uuidIfc.On("New").Return(remoteTaskUuid, nil).Once()
	mockBaseTask.On("SetWal", mock.Anything).Return(marinaError.ErrInternalError()).Once()
	mockRemoteClient := &mockClient.RemoteCatalogInterface{}
	mockRemoteCatalogClient := &MockRemoteCatalogClient{client: mockRemoteClient}
	remoteCatalogService = mockRemoteCatalogClient.NewRemoteCatalogService

	err := catalogItemDeleteTask.Run()

	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
	assert.Equal(t, [][]byte{peUuid.RawBytes()}, wal.GetData().GetCatalogItem().GetClusterUuidList())
	mockErgonBaseTask.AssertExpectations(t)
	mockBaseTask.AssertExpectations(t)
	configIfc.AssertExpectations(t)
	uuidIfc.AssertExpectations(t)
	mockRemoteClient.AssertExpectations(t)
}

func TestCatalogItemDeleteRunProtoError(t *testing.T) {

	mockErgonBaseTask := &mockTask.BaseTask{}
	mockBaseTask := &mockBase.MarinaBaseTaskInterface{}
	mockExtInterfaces := mockExternalInterfaces()
	mockIntInterfaces := mockInternalInterfaces()
	baseTask := &base.MarinaBaseTask{
		BaseTask:                   mockErgonBaseTask,
		MarinaBaseTaskInterface:    mockBaseTask,
		ExternalSingletonInterface: mockExtInterfaces,
		InternalSingletonInterface: mockIntInterfaces,
	}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	catalogItemDeleteTask.arg = &marinaIfc.CatalogItemDeleteArg{
		CatalogItemId: &marinaIfc.CatalogItemId{
			Version: &testCatalogItemVersion,
		},
	}
	pcUuid, _ := uuid4.New()
	peUuid, _ := uuid4.New()
	remoteTaskUuid, _ := uuid4.New()
	wal := &marinaIfc.PcTaskWalRecord{
		Data: &marinaIfc.PcTaskWalRecordData{
			CatalogItem: &marinaIfc.PcTaskWalRecordCatalogItemData{},
			TaskList: []*marinaIfc.EndpointTask{{
				ClusterUuid: pcUuid.RawBytes(),
				TaskUuid:    remoteTaskUuid.RawBytes(),
			}},
		},
	}
	mockBaseTask.On("Wal").Return(wal).Once()
	configIfc.On("PeClusterUuids").Return([]*uuid4.Uuid{peUuid}).Once()
	mockBaseTask.On("SetWal", mock.Anything).Return(nil).Once()
	mockErgonBaseTask.On("Save", mock.Anything).Return(nil).Once()
	configIfc.On("ClusterUuid").Return(pcUuid).Once()
	catalogItemUuid := testCatalogItemUuid.String()
	entities := []*insights_interface.EntityWithMetric{
		{
			EntityGuid: &insights_interface.EntityGuid{
				EntityId: &catalogItemUuid,
			},
		},
	}
	mockBaseTask.On("Wal").Return(wal).Once()
	uuidIfc.On("New").Return(remoteTaskUuid, nil).Once()
	mockBaseTask.On("SetWal", mock.Anything).Return(nil).Once()
	mockErgonBaseTask.On("Save", mock.Anything).Return(nil).Once()
	taskUuid, _ := uuid4.New()
	mockBaseTask.On("Proto").Return(&ergon.Task{Uuid: taskUuid.RawBytes()}).Once()
	mockRemoteClient := &mockClient.RemoteCatalogInterface{}
	mockRemoteClient.On("SendMsg", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	mockRemoteCatalogClient := &MockRemoteCatalogClient{client: mockRemoteClient}
	remoteCatalogService = mockRemoteCatalogClient.NewRemoteCatalogService
	taskMap := make(map[utils.RemoteEndpoint]*ergon.Task)
	fanoutIfc.On("PollAllRemoteTasks", mock.Anything, mock.Anything).Return(&taskMap).Once()
	cpdbIfc.On("Query", mock.Anything).Return(entities, nil).Once()
	idfIfc.On("DeleteEntities", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()
	protoIfc.On("Marshal", mock.Anything).Return(nil, marinaError.ErrInternalError()).Once()

	err := catalogItemDeleteTask.Run()

	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
	assert.Equal(t, [][]byte{peUuid.RawBytes()}, wal.GetData().GetCatalogItem().GetClusterUuidList())
	mockErgonBaseTask.AssertExpectations(t)
	mockBaseTask.AssertExpectations(t)
	configIfc.AssertExpectations(t)
	cpdbIfc.AssertExpectations(t)
	idfIfc.AssertExpectations(t)
	uuidIfc.AssertExpectations(t)
	mockRemoteClient.AssertExpectations(t)
	fanoutIfc.AssertExpectations(t)
	protoIfc.AssertExpectations(t)
}

func TestDeleteCatalogItem(t *testing.T) {
	mockExtInterfaces := mockExternalInterfaces()
	mockIntInterfaces := mockInternalInterfaces()
	baseTask := &base.MarinaBaseTask{
		ExternalSingletonInterface: mockExtInterfaces,
		InternalSingletonInterface: mockIntInterfaces,
	}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	catalogItemDeleteTask.arg = &marinaIfc.CatalogItemDeleteArg{
		CatalogItemId: &marinaIfc.CatalogItemId{
			Version: &testCatalogItemVersion,
		},
	}
	catalogItemUuid := testCatalogItemUuid.String()
	entities := []*insights_interface.EntityWithMetric{
		{
			EntityGuid: &insights_interface.EntityGuid{
				EntityId: &catalogItemUuid,
			},
		},
	}
	cpdbIfc.On("Query", mock.Anything).Return(entities, nil).Once()
	idfIfc.On("DeleteEntities", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()

	err := catalogItemDeleteTask.deleteCatalogItem("test_query")

	assert.NoError(t, err)
	cpdbIfc.AssertExpectations(t)
	idfIfc.AssertExpectations(t)
}

func TestDeleteCatalogItemCpdbNotFound(t *testing.T) {
	mockExtInterfaces := mockExternalInterfaces()
	mockIntInterfaces := mockInternalInterfaces()
	baseTask := &base.MarinaBaseTask{
		ExternalSingletonInterface: mockExtInterfaces,
		InternalSingletonInterface: mockIntInterfaces,
	}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	catalogItemDeleteTask.arg = &marinaIfc.CatalogItemDeleteArg{CatalogItemId: &marinaIfc.CatalogItemId{}}
	cpdbIfc.On("Query", mock.Anything).Return(nil, insights_interface.ErrNotFound).Once()

	err := catalogItemDeleteTask.deleteCatalogItem("test_query")

	assert.NoError(t, err)
	cpdbIfc.AssertExpectations(t)
}

func TestDeleteCatalogItemIdfNotFound(t *testing.T) {
	mockExtInterfaces := mockExternalInterfaces()
	mockIntInterfaces := mockInternalInterfaces()
	baseTask := &base.MarinaBaseTask{
		ExternalSingletonInterface: mockExtInterfaces,
		InternalSingletonInterface: mockIntInterfaces,
	}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	catalogItemDeleteTask.arg = &marinaIfc.CatalogItemDeleteArg{CatalogItemId: &marinaIfc.CatalogItemId{}}
	catalogItemUuid := testCatalogItemUuid.String()
	entities := []*insights_interface.EntityWithMetric{
		{
			EntityGuid: &insights_interface.EntityGuid{
				EntityId: &catalogItemUuid,
			},
		},
	}
	cpdbIfc.On("Query", mock.Anything).Return(entities, nil).Once()
	idfIfc.On("DeleteEntities", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(insights_interface.ErrNotFound).Once()

	err := catalogItemDeleteTask.deleteCatalogItem("test_query")

	assert.NoError(t, err)
	cpdbIfc.AssertExpectations(t)
	idfIfc.AssertExpectations(t)
}

func TestDeleteCatalogItemQueryError(t *testing.T) {
	catalogItemBaseTask := &CatalogItemBaseTask{}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	catalogItemDeleteTask.arg = &marinaIfc.CatalogItemDeleteArg{CatalogItemId: &marinaIfc.CatalogItemId{}}

	err := catalogItemDeleteTask.deleteCatalogItem("")

	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
}

func TestDeleteCatalogItemCpdbError(t *testing.T) {
	mockExtInterfaces := mockExternalInterfaces()
	mockIntInterfaces := mockInternalInterfaces()
	baseTask := &base.MarinaBaseTask{
		ExternalSingletonInterface: mockExtInterfaces,
		InternalSingletonInterface: mockIntInterfaces,
	}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	catalogItemDeleteTask.arg = &marinaIfc.CatalogItemDeleteArg{CatalogItemId: &marinaIfc.CatalogItemId{}}
	cpdbIfc.On("Query", mock.Anything).Return(nil, errors.New("oh no")).Once()

	err := catalogItemDeleteTask.deleteCatalogItem("test_query")

	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
	cpdbIfc.AssertExpectations(t)
}

func TestDeleteCatalogItemIdfError(t *testing.T) {
	mockExtInterfaces := mockExternalInterfaces()
	mockIntInterfaces := mockInternalInterfaces()
	baseTask := &base.MarinaBaseTask{
		ExternalSingletonInterface: mockExtInterfaces,
		InternalSingletonInterface: mockIntInterfaces,
	}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	catalogItemDeleteTask.arg = &marinaIfc.CatalogItemDeleteArg{CatalogItemId: &marinaIfc.CatalogItemId{}}
	catalogItemUuid := testCatalogItemUuid.String()
	entities := []*insights_interface.EntityWithMetric{
		{
			EntityGuid: &insights_interface.EntityGuid{
				EntityId: &catalogItemUuid,
			},
		},
	}
	cpdbIfc.On("Query", mock.Anything).Return(entities, nil).Once()
	idfIfc.On("DeleteEntities", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("oh no")).Once()

	err := catalogItemDeleteTask.deleteCatalogItem("test_query")

	assert.Error(t, err)
	assert.IsType(t, new(marinaError.InternalError), err)
	cpdbIfc.AssertExpectations(t)
	idfIfc.AssertExpectations(t)
}

func TestFanoutCatalogRequestsWithEndpointNameError(t *testing.T) {
	mockErgonBaseTask := &mockTask.BaseTask{}
	mockBaseTask := &mockBase.MarinaBaseTaskInterface{}
	mockExtInterfaces := mockExternalInterfaces()
	mockIntInterfaces := mockInternalInterfaces()
	baseTask := &base.MarinaBaseTask{
		BaseTask:                   mockErgonBaseTask,
		MarinaBaseTaskInterface:    mockBaseTask,
		ExternalSingletonInterface: mockExtInterfaces,
		InternalSingletonInterface: mockIntInterfaces,
	}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	arg := &marinaIfc.CatalogItemDeleteArg{CatalogItemId: &marinaIfc.CatalogItemId{}}
	catalogItemDeleteTask.arg = arg
	pcUuid, _ := uuid4.New()
	peUuid, _ := uuid4.New()
	remoteTaskUuid, _ := uuid4.New()
	wal := &marinaIfc.PcTaskWalRecord{
		Data: &marinaIfc.PcTaskWalRecordData{
			CatalogItem: &marinaIfc.PcTaskWalRecordCatalogItemData{},
			TaskList: []*marinaIfc.EndpointTask{{
				ClusterUuid: peUuid.RawBytes(),
				TaskUuid:    remoteTaskUuid.RawBytes(),
			}},
		},
	}
	mockBaseTask.On("Wal").Return(wal).Once()
	embedded, _ := proto.Marshal(arg)
	payload := &ergon.PayloadOrEmbeddedValue{Embedded: embedded}
	req := &ergon.MetaRequest{Arg: payload}
	task := &ergon.Task{Request: req}
	mockBaseTask.On("Proto").Return(task).Once()
	mockRemoteClient := &mockClient.RemoteCatalogInterface{}
	mockRemoteClient.On("SendMsg", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("oh no")).Once()
	mockRemoteCatalogClient := &MockRemoteCatalogClient{client: mockRemoteClient}
	remoteCatalogService = mockRemoteCatalogClient.NewRemoteCatalogService
	name := "TestName"
	configIfc.On("PeClusterName", mock.Anything).Return(&name).Once()

	remoteEndpoint := utils.RemoteEndpoint{RemoteClusterUuid: *peUuid, SourceClusterUuid: *pcUuid}
	err := catalogItemDeleteTask.fanoutCatalogRequests([]utils.RemoteEndpoint{remoteEndpoint}, pcUuid)

	assert.Error(t, err)
	assert.IsType(t, marinaError.ErrCatalogTaskForwardError, err)
	mockBaseTask.AssertExpectations(t)
	mockRemoteClient.AssertExpectations(t)
	configIfc.AssertExpectations(t)
}

func TestFanoutCatalogRequestsWithoutEndpointNameError(t *testing.T) {
	mockErgonBaseTask := &mockTask.BaseTask{}
	mockBaseTask := &mockBase.MarinaBaseTaskInterface{}
	mockExtInterfaces := mockExternalInterfaces()
	mockIntInterfaces := mockInternalInterfaces()
	baseTask := &base.MarinaBaseTask{
		BaseTask:                   mockErgonBaseTask,
		MarinaBaseTaskInterface:    mockBaseTask,
		ExternalSingletonInterface: mockExtInterfaces,
		InternalSingletonInterface: mockIntInterfaces,
	}
	catalogItemBaseTask := &CatalogItemBaseTask{MarinaBaseTask: baseTask}
	catalogItemDeleteTask := NewCatalogItemDeleteTask(catalogItemBaseTask)
	arg := &marinaIfc.CatalogItemDeleteArg{CatalogItemId: &marinaIfc.CatalogItemId{}}
	catalogItemDeleteTask.arg = arg
	pcUuid, _ := uuid4.New()
	peUuid, _ := uuid4.New()
	remoteTaskUuid, _ := uuid4.New()
	wal := &marinaIfc.PcTaskWalRecord{
		Data: &marinaIfc.PcTaskWalRecordData{
			CatalogItem: &marinaIfc.PcTaskWalRecordCatalogItemData{},
			TaskList: []*marinaIfc.EndpointTask{{
				ClusterUuid: peUuid.RawBytes(),
				TaskUuid:    remoteTaskUuid.RawBytes(),
			}},
		},
	}
	mockBaseTask.On("Wal").Return(wal).Once()
	embedded, _ := proto.Marshal(arg)
	payload := &ergon.PayloadOrEmbeddedValue{Embedded: embedded}
	req := &ergon.MetaRequest{Arg: payload}
	task := &ergon.Task{Request: req}
	mockBaseTask.On("Proto").Return(task).Once()

	uuidIfc.On("New").Return(remoteTaskUuid, nil).Once()
	mockBaseTask.On("SetWal", mock.Anything).Return(nil).Once()
	mockErgonBaseTask.On("Save", mock.Anything).Return(nil).Once()

	mockRemoteClient := &mockClient.RemoteCatalogInterface{}
	mockRemoteClient.On("SendMsg", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("oh no")).Once()
	mockRemoteCatalogClient := &MockRemoteCatalogClient{client: mockRemoteClient}
	remoteCatalogService = mockRemoteCatalogClient.NewRemoteCatalogService

	remoteEndpoint := utils.RemoteEndpoint{}
	err := catalogItemDeleteTask.fanoutCatalogRequests([]utils.RemoteEndpoint{remoteEndpoint}, pcUuid)

	assert.Error(t, err)
	assert.IsType(t, marinaError.ErrCatalogTaskForwardError, err)
	mockBaseTask.AssertExpectations(t)
	mockErgonBaseTask.AssertExpectations(t)
	mockRemoteClient.AssertExpectations(t)
}
