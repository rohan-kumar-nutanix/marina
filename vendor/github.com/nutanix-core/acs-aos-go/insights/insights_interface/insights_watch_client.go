/*
 * Copyrfunc (insights_svc *InsightsService) ght (c) 2016 Nutanix Inc. All
 * rights reserved.
 *
 * This module defines Insights Watch Semantics
 *
 */

package insights_interface

import (
  "errors"
  "flag"
  "fmt"
  util_misc "github.com/nutanix-core/acs-aos-go/nutanix/util-go/misc"
  uuid4 "github.com/nutanix-core/acs-aos-go/nutanix/util-go/uuid4"
  "strings"
  "sync"
  "time"

  glog "github.com/golang/glog"
  proto "github.com/golang/protobuf/proto"
)

var (
  WatchClientTimeout = flag.Uint(
    "watch_client_default_timeout",
    60,
    "Default Watch Client Timeout period.")

  WatchClientFiredQSize = flag.Int(
    "watch_client_default_fired_queue_size",
    1000,
    "Default Watch Client Fired Queue size.")
)

// struct holding relevant information for a watch callback
type InsightsWatch struct {
  WatchName              string
  EntityType             string
  EntityUuid             string
  FiredWatchInitialSeqNo uint64
  WatchId                uint64
}

// Insights watch interface
// To implement a watch callback - define a struct that implements
// InsightsWatchCb interface. The only method that requires to be
// implemented is WatchCallback. This method is invoked when the
// watch callback is fired.
type InsightsWatchCb interface {
  WatchCallback(fired *FiredWatch) error

  setupCallback(watchName string, entityType string, Uuid string,
    FiredInitialSeqNo uint64, watchId uint64)

  GetWatchName() string

  GetEntityType() string

  GetEntityUuid() string

  GetFiredWatchInitialSeqNo() uint64

  GetWatchId() uint64
}

type WatchCbInfo struct {
  WatchCb    InsightsWatchCb
  FiredWatch *FiredWatch
  CloseQ     bool
}

type RegisterWatchResult struct {
  Error      error            // error if any on registration
  EntityList []*Entity        // Entity state on registration of watch
  WatchInfo  *EntityWatchInfo // watch info used for this registration
}

// Watch client definition
type WatchClient struct {
  ClientId            string
  FiredWatchLastSeqNo uint64
  FiredQSize          int32 //size of fired Q.
  WatchService        InsightsServiceInterface
  RegisteredWatches   map[string]*EntityWatchInfo // All registered watches.
  WatchClientProto    *WatchClientProto
  FiredWatchQ         chan *WatchCbInfo
  stateMu             sync.Mutex
  state               int32
  mu                  sync.Mutex
  done                bool
  mapMu               sync.Mutex
  entityWatches       map[string]InsightsWatchCb // Watch <-> Callback
  entityUnregWatches  []string                   // list of watches to be unregistered
  watchSync           chan bool                  // synchronize RegisterWatch/UnregisterWatch calls
  isChannelOpen       bool                       // is the channel open
  resetOnClusterRegistration bool                // reset in case of cluster registration
  resetOnClusterUnregistration bool              // reset in case of cluster unregistration
}

const (
  InsightsEvent_kInit           int32 = 1
  InsightsEvent_kInitRegister   int32 = 2
  InsightsEvent_kRegisterDone   int32 = 3
  InsightsEvent_kStart          int32 = 4
  InsightsEvent_kInitUnregister int32 = 5
  InsightsEvent_kUnregistered   int32 = 6
  InsightsEvent_kReset          int32 = 7
  InsightsEvent_kInitStop       int32 = 8
  InsightsEvent_kStopped        int32 = 9
  InsightsEvent_kInitReregister int32 = 10
  InsightsEvent_kReregisterDone int32 = 11
  InsightsEvent_kError          int32 = 12
)

var InsightsEventString = map[int32]string{
  1:  "kInit",
  2:  "kInitRegister",
  3:  "kRegisterDone",
  4:  "kStart",
  5:  "kUnregistering",
  6:  "kUnregistered",
  7:  "kReset",
  8:  "kInitStop",
  9:  "kStopped",
  10: "kInitReregister",
  11: "kReregisterDone",
  12: "kError",
}

const (
  InsightsWatchState_kUnInitialized int32 = 0
  InsightsWatchState_kInitialized   int32 = 1
  InsightsWatchState_kRegistering   int32 = 2
  InsightsWatchState_kRegistered    int32 = 3
  InsightsWatchState_kStarted       int32 = 4
  InsightsWatchState_kStopping      int32 = 5
  InsightsWatchState_kStopped       int32 = 6
  InsightsWatchState_kReset         int32 = 7
  InsightsWatchState_kUnregistering int32 = 8
  InsightsWatchState_kUnregistered  int32 = 9
  InsightsWatchState_kReregistering int32 = 10
  InsightsWatchState_kError         int32 = 11
)

var InsightsStateString = map[int32]string{
  0:  "kUninit",
  1:  "kInitialized",
  2:  "kRegistering",
  3:  "kRegistered",
  4:  "kStarted",
  5:  "kStopping",
  6:  "kStopped",
  7:  "kReset",
  8:  "kUnregistering",
  9:  "kUnregistered",
  10: "kReregistering",
  11: "kError",
}

var wcStateFSM [12][14]int32

// For creating a new instance of EntityWatchInfo user should use these helper
// methods:- NewEntityWatchInfo(), NewEntityWatchInfoForEvictableEntities(),
// NewEntitySchemaWatchInfo(), NewEntityWatchInfoWithProjectionList() or
// NewEvictableEntityWatchInfoWithProjectionList().
type EntityWatchInfo struct {
  // Use 'EntityGuid' only for entity watches. For Schema watches use
  // 'EntityTypeName'.
  EntityGuid *EntityGuid
  // WatchType can be Create/update/delete for entity watches and
  // RegisterEntityType/UpdateEntityType/RegisterMetricType/UpdateMetricType
  // for schema watches.
  // This is an internal field, user don't have to set it.
  WatchType                      int32
  WatchName                      string // watch name
  WatchCategory                  Watch_WatchCategory
  // Metric name. Used only for attribute watches.
  WatchMetric                    string
  ReturnPreviousEntityState      bool
  GetCurrentState                bool
  FilterExpr                     *BooleanExpression
  WatchCb                        InsightsWatchCb
  ApplyProjectionOnFiredWatch    bool
  FiredWatchProjectionList       []string
  // Use 'EntityTypeName' only for schema watches. For entity watches use
  // 'EntityGuid'.
  EntityTypeName string
}

// This struct defines all possible watch types for schema watches.
type SchemaWatchTypes struct {
  RegisterEntityType bool
  UpdateEntityType   bool
  RegisterMetricType bool
  UpdateMetricType   bool
}

func (cb *InsightsWatch) WatchCallback(fired *FiredWatch) error {
  glog.Info("InsightsWatch-Cb - Not Implemented ")
  return errors.New("InsightsWatchCb: Not implemented")
}

func (cb *InsightsWatch) setupCallback(watchName string, entityType string,
  uuid string, FiredInitialSeqNo uint64, watchId uint64) {

  if cb == nil {
    glog.Error("Nil CB")
    return
  }
  cb.WatchName = watchName
  cb.EntityType = entityType
  cb.EntityUuid = uuid
  cb.FiredWatchInitialSeqNo = FiredInitialSeqNo
  cb.WatchId = watchId
}

func (cb *InsightsWatch) GetWatchName() string {
  return cb.WatchName
}

func (cb *InsightsWatch) GetEntityType() string {
  return cb.EntityType
}

func (cb *InsightsWatch) GetEntityUuid() string {
  return cb.EntityUuid
}

func (cb *InsightsWatch) GetFiredWatchInitialSeqNo() uint64 {
  return cb.FiredWatchInitialSeqNo
}

func (cb *InsightsWatch) GetWatchId() uint64 {
  return cb.WatchId
}

var (
  ErrChannelClosed      = InsightsError("Channel Closed", 3100)
  ErrInvalidClientState = InsightsError("Invalid client state error",
    3101)
  ErrClientAlreadyStarted  = InsightsError("Client already started", 3102)
  ErrWatchClientInProgress = InsightsError(
    "Watch client registration in progress ",
    3103)
  ErrInvalidWatchClientSession = InsightsError("Invalid watch client session",
    3103)
)

// Create an instance of NewWatchClient
func NewWatchClient(clientId string,
  watchService InsightsServiceInterface) *WatchClient {

  wc := &WatchClient{
    ClientId:     clientId,
    FiredQSize:   int32(*WatchClientFiredQSize),
    WatchService: watchService,
    done:         false,
  }

  wc.entityWatches = make(map[string]InsightsWatchCb)
  // Q size determined by the caller.
  wc.FiredWatchQ = make(chan *WatchCbInfo, wc.FiredQSize)
  wc.RegisteredWatches = make(map[string]*EntityWatchInfo)

  // raise event init
  wc.transitionWatchState(InsightsEvent_kInit)
  return wc

}

// Done with this watch.
func (wc *WatchClient) Stop() {
  wc.mu.Lock()
  wc.done = true
  wc.mu.Unlock()
  wc.transitionWatchState(InsightsEvent_kInitStop)
}

// Helper method for creating an instance of EntityWatchInfo for entity watches
// for evictable entities.
func (wc *WatchClient) NewEntityWatchInfoForEvictableEntities(
  guid *EntityGuid, name string, return_previous_entity_state bool,
  cb InsightsWatchCb, filterExpression *BooleanExpression,
  watchMetric string) *EntityWatchInfo {

  curr_state := false
  if guid.EntityId != nil {
    curr_state = true
  }
  return wc.newEntityWatchInfoInternal(
    guid, name, curr_state, return_previous_entity_state, cb,
    Watch_kEntityWatch /* WatchCategory */, filterExpression, watchMetric,
    false /* ApplyProjectionOnFiredWatch */, nil /* FiredWatchProjectionList */)
}

// Helper method for creating an instance of EntityWatchInfo for entity watches
// with projection list for evictable entities.
func (wc *WatchClient) NewEvictableEntityWatchInfoWithProjectionList(
  guid *EntityGuid, name string, return_previous_entity_state bool,
  cb InsightsWatchCb, filterExpression *BooleanExpression,
  watchMetric string, enableProjectionOnFiredWatch bool,
  firedWatchProjectionList []string) *EntityWatchInfo {

  curr_state := false
  if guid.EntityId != nil {
    curr_state = true
  }
  return wc.newEntityWatchInfoInternal(
    guid, name, curr_state, return_previous_entity_state, cb,
    Watch_kEntityWatch /* WatchCategory */, filterExpression, watchMetric,
    enableProjectionOnFiredWatch, firedWatchProjectionList)
}

// Helper method for creating an instance of EntityWatchInfo for entity watches.
func (wc *WatchClient) NewEntityWatchInfo(
  guid *EntityGuid, name string, get_current_state bool,
  return_previous_entity_state bool, cb InsightsWatchCb,
  filterExpression *BooleanExpression, watchMetric string) *EntityWatchInfo {

  return wc.newEntityWatchInfoInternal(
    guid, name, get_current_state, return_previous_entity_state, cb,
    Watch_kEntityWatch /* WatchCategory */, filterExpression, watchMetric,
    false /* ApplyProjectionOnFiredWatch */, nil /* FiredWatchProjectionList */)
}

// Helper method for creating an instance of EntityWatchInfo for entity watches
// with projection list.
func (wc *WatchClient) NewEntityWatchInfoWithProjectionList(
  guid *EntityGuid, name string, get_current_state bool,
  return_previous_entity_state bool, cb InsightsWatchCb,
  filterExpression *BooleanExpression, watchMetric string,
  enableProjectionOnFiredWatch bool,
  firedWatchProjectionList []string) *EntityWatchInfo {

  return wc.newEntityWatchInfoInternal(
    guid, name, get_current_state, return_previous_entity_state, cb,
    Watch_kEntityWatch /* WatchCategory */, filterExpression, watchMetric,
    enableProjectionOnFiredWatch, firedWatchProjectionList)
}

// Helper method for creating an instance of EntityWatchInfo for Schema Watches.
func (wc *WatchClient) NewEntitySchemaWatchInfo(
  name string, getCurrentState bool,
  cb InsightsWatchCb) *EntityWatchInfo {

  // Creating a dummy guid instance, as guid is not used for schema watches
  // but GoLang doesn't support default parmas and method overloading,
  // hence passing this dummy guid to newEntityWatchInfoInternal() method.
  guid := &EntityGuid{}
  return wc.newEntityWatchInfoInternal(
    guid, name, getCurrentState, false /* ReturnPreviousEntityState */, cb,
    Watch_kEntitySchemaWatch, /* WatchCategory */
    nil /* filterExpression */, "" /* WatchMetric */,
    false /* ApplyProjectionOnFiredWatch */, nil /* FiredWatchProjectionList */)
}

// Internal method to get a new instance of Watch Info.
func (wc *WatchClient) newEntityWatchInfoInternal(
  guid *EntityGuid, name string, getCurrentState bool,
  returnPreviousEntityState bool, cb InsightsWatchCb,
  watchCategory Watch_WatchCategory, filterExpression *BooleanExpression,
  watchMetric string, enableProjectionOnFiredWatch bool,
  firedWatchProjectionList []string) *EntityWatchInfo {

  return &EntityWatchInfo{
    EntityGuid:                   guid,
    WatchName:                    name,
    WatchCategory:                watchCategory,
    ReturnPreviousEntityState:    returnPreviousEntityState,
    GetCurrentState:              getCurrentState,
    WatchCb:                      cb,
    FilterExpr:                   filterExpression,
    WatchMetric:                  watchMetric,
    ApplyProjectionOnFiredWatch:  enableProjectionOnFiredWatch,
    FiredWatchProjectionList:     firedWatchProjectionList,
  }
}

func (wc *WatchClient) finishRegistration(err error) {
  if err != nil {
    wc.transitionWatchState(InsightsEvent_kReset)
    //close the channel
    if wc.isChannelOpen {
      close(wc.watchSync)
      wc.isChannelOpen = false
    }
  } else {
    wc.transitionWatchState(InsightsEvent_kRegisterDone)
    // Enable RegisterWatch/UnregisterWatch to enter call block.
    wc.watchSync <- true
  }
}

// Gets the int value of ResetWatchClient to be set in WatchClientProto
// from boolean flags resetOnClusterRegistration and
// resetOnClusterUnregistration.
func resetWatchClientFromBoolToInt(resetOnClusterRegistration bool,
  resetOnClusterUnregistration bool) int32 {

  var resetWatchClientOnOperations int32 = kNoReset
  if resetOnClusterRegistration {
    resetWatchClientOnOperations |= kResetOnClusterRegistration
  }
  if resetOnClusterUnregistration {
    resetWatchClientOnOperations |= kResetOnClusterUnregistration
  }
  return resetWatchClientOnOperations
}

// Register a watch client
func (wc *WatchClient) Register() error {
  var err error = nil
  var state int32
  state, err = wc.transitionWatchState(InsightsEvent_kInitRegister)
  if err != nil {
    switch state {
    case InsightsWatchState_kRegistering:
      err = ErrWatchClientInProgress
    case InsightsWatchState_kRegistered:
      err = ErrWatchClientAlreadyRegistered
    default:
    }
    return err
  }

  defer func() {
    wc.finishRegistration(err)
  }()

  var sessionId *uuid4.Uuid
  sessionId, err = uuid4.New()
  // The reset will be received after cluster registration/unregistration after
  // IDF has loaded/unloaded all the entities of that cluster so that client can
  // fetch the current state and get the new entities in current state.
  var resetWatchClientOnOperations int32 =  resetWatchClientFromBoolToInt(
    wc.resetOnClusterRegistration, wc.resetOnClusterUnregistration)
  wc.WatchClientProto = &WatchClientProto{
    ClientId:  proto.String(wc.ClientId),
    SessionId: proto.String(sessionId.UuidToString()),
    ResetWatchClientOnOperations: proto.Int32(resetWatchClientOnOperations),
  }
  regArg := &RegisterWatchClientArg{
    WatchClient: wc.WatchClientProto,
  }
  regRet := &RegisterWatchClientRet{}
  err = wc.WatchService.SendMsg("RegisterWatchClient", regArg, regRet, nil)
  if err != nil {
    return err
  }

  // Registration done - open boolean channel to synchronize
  // RegisterWatch/UnreigsterWatch enabled only after the client
  // has been registered
  wc.watchSync = make(chan bool, 1)

  // For the first GetFiredWatchList call - FiredWatchLastSeqNum should be 0
  // reset it.
  wc.FiredWatchLastSeqNo = 0
  wc.mu.Lock()
  wc.done = false
  wc.mu.Unlock()

  glog.Infof("Registered Watch for client %s-%s\n", wc.ClientId, sessionId)
  // Client is registered -> raise the Registered event

  return nil
}

// Unregister watch client
func (wc *WatchClient) Unregister() error {
  //close the call block channel
  wc.isChannelOpen = <-wc.watchSync
  if wc.isChannelOpen {
    close(wc.watchSync)
    wc.isChannelOpen = false
  }

  var err error = nil

  _, err = wc.transitionWatchState(InsightsEvent_kInitUnregister)
  if err != nil {
    return err
  }

  glog.Infof("Unregistering watch client with session-id %s",
    wc.WatchClientProto.GetSessionId())

  unregArg := &UnregisterWatchClientArg{
    WatchClient: wc.WatchClientProto,
  }

  unregRet := &UnregisterWatchClientRet{}

  err = wc.WatchService.SendMsg("UnregisterWatchClient", unregArg, unregRet, nil)
  if err != nil {
    glog.Warning(err)
  }
  _, err = wc.transitionWatchState(InsightsEvent_kUnregistered)

  wc.mapMu.Lock()
  // cleanup the callback list.
  wc.entityWatches = make(map[string]InsightsWatchCb)
  wc.mapMu.Unlock()

  _, err = wc.transitionWatchState(InsightsEvent_kInit)
  return err
}

func (wc *WatchClient) isWatchUnregistered(name string) bool {
  isUnregistered := false
  for i := 0; i < len(wc.entityUnregWatches); i++ {
    if strings.Compare(wc.entityUnregWatches[i], name) == 0 {
      isUnregistered = true
    }
  }
  return isUnregistered
}

func (wc *WatchClient) removeFromUnregisterList(name string) bool {
  removed := false
  index := 0
  for i := 0; i < len(wc.entityUnregWatches); i++ {
    if strings.Compare(wc.entityUnregWatches[i], name) == 0 {
      removed = true
      index = i
    }
  }
  if removed {
    a := wc.entityUnregWatches
    // this is just a search list - not a ordered list.
    // exchange with the last element and shrink it.
    a[index] = a[len(a)-1]
    a[len(a)-1] = "" // zero value for string
    a = a[:len(a)-1]
    wc.entityUnregWatches = a
  }

  return removed
}

// This simply adds to the Unreg Q.
func (wc *WatchClient) addToUnregisterList(name string) {
  if !wc.isWatchUnregistered(name) {
    wc.entityUnregWatches = append(wc.entityUnregWatches, name)
  }
}

// Unregister a specific watch by name.
func (wc *WatchClient) UnregisterWatch(name string) error {
  // Wait here till we can proceed
  wc.isChannelOpen = <-wc.watchSync
  if !wc.isChannelOpen {
    // Channel has been closed -> new watch client. Return error
    // to the client
    return ErrInvalidWatchClientSession
  }

  defer func() {
    // Release the map mutex lock.
    wc.mapMu.Unlock()
    // Release when done - so that other callers can get in.
    wc.watchSync <- true
  }()

  glog.Infof("Unregistering watch - %s", name)

  // get the cached callback for this watch.
  wc.mapMu.Lock()
  watchCb := wc.entityWatches[name]
  if watchCb == nil {
    return ErrNotFound.SetCause(errors.New("Watch callback not found"))
  }

  var idList []uint64
  var nameList []string

  idList = append(idList, watchCb.GetWatchId())
  nameList = append(nameList, name)

  unregArg := &UnregisterWatchArg{
    WatchClient:   wc.WatchClientProto,
    WatchNameList: nameList,
    WatchIdList:   idList,
  }
  unregRet := &UnregisterWatchRet{}

  err := wc.WatchService.SendMsg("UnregisterWatch", unregArg, unregRet, nil)
  if err != nil {
    return err
  }

  // Remember the unregistered watch - this can be used to verify fired watches.
  // It is cleaned up after each iteration of fired watches.
  wc.addToUnregisterList(name)
  delete(wc.entityWatches, name)
  delete(wc.RegisteredWatches, name)
  return nil
}

// Converts a given WatchType struct to WatchType int.
func entityWatchOperationsToWatchType(createWatch bool, updateWatch bool,
  deleteWatch bool) int32 {

  var watchType int32 = 0
  if createWatch {
    watchType |= kEntityCreate
  }
  if updateWatch {
    watchType |= kEntityUpdate
  }
  if deleteWatch {
    watchType |= kEntityDelete
  }
  return watchType
}

// Wrapper method for consumers to call to watch creation of entities of type
// 'entity_type_name' specified in 'entity_guid' in watchInfo.
// For creating an instance of EntityWatchInfo use helper methods
// NewEntityWatchInfo(), NewEntityWatchInfoForEvictableEntities(),
// NewEntityWatchInfoWithProjectionList() or
// NewEvictableEntityWatchInfoWithProjectionList()
func (wc *WatchClient) WatchEntityCreateOfType(watchInfo *EntityWatchInfo) (
  []*Entity, error) {

  watchInfo.WatchType = entityWatchOperationsToWatchType(
    true /* createWatch*/, false /* updateWatch*/, false /* deleteWatch*/)
  return wc.watchEntityType(watchInfo)
}

// Wrapper method for consumers to call to watch update of entities of type
// 'entity_type_name' specified in 'entity_guid' in watchInfo.
// For creating an instance of EntityWatchInfo use helper methods
// NewEntityWatchInfo(), NewEntityWatchInfoForEvictableEntities(),
// NewEntityWatchInfoWithProjectionList() or
// NewEvictableEntityWatchInfoWithProjectionList()
func (wc *WatchClient) WatchEntityUpdateOfType(watchInfo *EntityWatchInfo) (
  []*Entity, error) {

  watchInfo.WatchType = entityWatchOperationsToWatchType(
    false /* createWatch*/, true /* updateWatch*/, false /*deleteWatch*/)
  return wc.watchEntityType(watchInfo)
}

// Wrapper method for consumers to call to watch deletion of entities of type
// 'entity_type_name' specified in 'entity_guid' in watchInfo.
// For creating an instance of EntityWatchInfo use helper methods
// NewEntityWatchInfo(), NewEntityWatchInfoForEvictableEntities(),
// NewEntityWatchInfoWithProjectionList() or
// NewEvictableEntityWatchInfoWithProjectionList()
func (wc *WatchClient) WatchEntityDeleteOfType(watchInfo *EntityWatchInfo) (
  []*Entity, error) {

  watchInfo.WatchType = entityWatchOperationsToWatchType(
    false /* createWatch*/, false /* updateWatch*/, true /*deleteWatch*/)
  return wc.watchEntityType(watchInfo)
}

// Wrapper method for consumers to call to watch operations on entities of
// type 'entity_type_name' specified in 'entity_guid' in watchInfo.
// For creating an instance of EntityWatchInfo use helper methods
// NewEntityWatchInfo(), NewEntityWatchInfoForEvictableEntities(),
// NewEntityWatchInfoWithProjectionList() or
// NewEvictableEntityWatchInfoWithProjectionList()
func (wc *WatchClient) CompositeWatchOnEntitiesOfType(
  watchInfo *EntityWatchInfo, createWatch bool, updateWatch bool,
  deleteWatch bool) ([]*Entity, error) {

  watchInfo.WatchType = entityWatchOperationsToWatchType(
    createWatch, updateWatch, deleteWatch)
  return wc.watchEntityType(watchInfo)
}

// Note that there are various 'watchTypes'. See 'entity_watch_type.go'.
func (wc *WatchClient) watchEntityType(watchInfo *EntityWatchInfo) (
  []*Entity, error) {

  if watchInfo.EntityGuid.EntityId != nil {
    return nil, errors.New("Entity Id should not be set for Entity type watch")
  }
  var entList []*Entity
  dataList, err := wc.RegisterWatch(watchInfo)
  if err != nil {
    return nil, err
  }
  for i := 0; i < len(dataList); i++ {
    entList = append(entList, dataList[i].GetChangedEntityData())
  }
  return entList, nil
}

// Wrapper method for consumers to call to watch creation of entity described
// by 'entity_guid' in watchInfo.
// For creating an instance of EntityWatchInfo use helper methods:-
// NewEntityWatchInfo() or NewEntityWatchInfoForEvictableEntities().
func (wc *WatchClient) WatchEntityCreate(watchInfo *EntityWatchInfo) (
  *Entity, error) {

  if watchInfo.EntityGuid.EntityId == nil {
    return nil, errors.New("Entity Id should be set for Entity watch")
  }
  watchInfo.WatchType = entityWatchOperationsToWatchType(
    true /* createWatch */, false /* updateWatch */, false /* deleteWatch */)

  entList, err := wc.RegisterWatch(watchInfo)
  var current *Entity = nil
  if err == nil && len(entList) > 0 {
    current = entList[0].GetChangedEntityData()
  }
  return current, err
}

// Wrapper method for consumers to call to watch updates of entity described
// by 'entity_guid' in watchInfo.
// For creating an instance of EntityWatchInfo use helper methods
// NewEntityWatchInfo(), NewEntityWatchInfoForEvictableEntities(),
// NewEntityWatchInfoWithProjectionList() or
// NewEvictableEntityWatchInfoWithProjectionList().
func (wc *WatchClient) WatchEntityUpdate(watchInfo *EntityWatchInfo) (
  *Entity, error) {

  if watchInfo.EntityGuid.EntityId == nil {
    return nil, errors.New("Entity Id should be set for Entity watch")
  }
  watchInfo.WatchType = entityWatchOperationsToWatchType(
    false /* createWatch */, true /* updateWatch */, false /* deleteWatch */)
  entList, err := wc.RegisterWatch(watchInfo)

  var current *Entity = nil
  if err == nil && len(entList) > 0 {
    current = entList[0].GetChangedEntityData()
  }
  return current, err
}

// Wrapper method for consumers to call to watch deletion of entity described
// by 'entity_guid' in watchInfo.
// For creating an instance of EntityWatchInfo use helper methods
// NewEntityWatchInfo(), NewEntityWatchInfoForEvictableEntities(),
// NewEntityWatchInfoWithProjectionList() or
// NewEvictableEntityWatchInfoWithProjectionList.
func (wc *WatchClient) WatchEntityDelete(watchInfo *EntityWatchInfo) (
  *Entity, error) {

  if watchInfo.EntityGuid.EntityId == nil {
    return nil, errors.New("Entity Id should be set for Entity watch")
  }
  watchInfo.WatchType = entityWatchOperationsToWatchType(
    false /* createWatch */, false /* updateWatch */, true /* deleteWatch */)
  entList, err := wc.RegisterWatch(watchInfo)

  var current *Entity = nil
  if err == nil && len(entList) > 0 {
    current = entList[0].GetChangedEntityData()
  }
  return current, err
}

// Wrapper method for consumers to call to watch the given operations on
// entity described by 'entity_guid' in watchInfo.
// For creating an instance of EntityWatchInfo use helper methods
// NewEntityWatchInfo(), NewEntityWatchInfoForEvictableEntities(),
// NewEntityWatchInfoWithProjectionList() or
// NewEvictableEntityWatchInfoWithProjectionList.
func (wc *WatchClient) CompositeWatchOnEntity(
  watchInfo *EntityWatchInfo, createWatch bool, updateWatch bool,
  deleteWatch bool) (*Entity, error) {

  if watchInfo.EntityGuid.EntityId == nil {
    return nil, errors.New("Entity Id should be set for Entity watch")
  }
  watchInfo.WatchType =
    entityWatchOperationsToWatchType(createWatch, updateWatch, deleteWatch)
  entList, err := wc.RegisterWatch(watchInfo)

  var current *Entity = nil
  if err == nil && len(entList) > 0 {
    current = entList[0].GetChangedEntityData()
  }
  return current, err
}

func schemaWatchOperationsToWatchType(
  watchOnOpers *SchemaWatchTypes) int32 {

  var watchType int32 = 0
  if watchOnOpers.RegisterEntityType {
    watchType |= kRegisterEntityType
  }
  if watchOnOpers.UpdateEntityType {
    watchType |= kUpdateEntityType
  }
  if watchOnOpers.RegisterMetricType {
    watchType |= kRegisterMetricType
  }
  if watchOnOpers.UpdateMetricType {
    watchType |= kUpdateMetricType
  }
  return watchType
}

// Method for registering schema watches for all entity types.
// For creating an instance of EntityWatchInfo use helper methods
// NewEntitySchemaWatchInfo().
func (wc *WatchClient) RegisterCompositeSchemaWatchForAllEntityTypes(
  watchInfo *EntityWatchInfo, watchOnOpers *SchemaWatchTypes) (
  []*EntitySchemaChange, error) {

  return wc.RegisterCompositeSchemaWatchForEntityType("", watchInfo,
    watchOnOpers)
}

// Method for registering schema watches for given entity type.
// For creating an instance of EntityWatchInfo use helper methods
// NewEntitySchemaWatchInfo().
func (wc *WatchClient) RegisterCompositeSchemaWatchForEntityType(
  entityType string, watchInfo *EntityWatchInfo,
  watchOnOpers *SchemaWatchTypes) ([]*EntitySchemaChange, error) {

  watchInfo.EntityTypeName = entityType
  watchInfo.WatchType = schemaWatchOperationsToWatchType(watchOnOpers)
  currentData, err := wc.RegisterWatch(watchInfo)

  var current []*EntitySchemaChange
  if err != nil {
    return nil, err
  }

  for i := 0; i < len(currentData); i++ {
    current = append(current, currentData[i].GetEntitySchemaChange())
  }

  return current, nil
}

// Internal method used by wrapper methods for registering a watch based on
// params specified in 'watchInfo'.
func (wc *WatchClient) RegisterWatch(watchInfo *EntityWatchInfo) (
  []*FiredWatch_ChangedData, error) {

  if watchInfo.ApplyProjectionOnFiredWatch == false &&
     watchInfo.FiredWatchProjectionList != nil {

    return nil ,errors.New("FiredWatchProjectionList has to be nil if " +
      "ApplyProjectionOnFiredWatch is false")
  } else if watchInfo.ApplyProjectionOnFiredWatch == true &&
            watchInfo.FiredWatchProjectionList == nil {

    return nil, errors.New("FiredWatchProjectionList should not be nil if " +
      "ApplyProjectionOnFiredWatch is true")
  }

  watchType := watchInfo.WatchType
  entityGuid := watchInfo.EntityGuid

  nowUsecs := (uint64)(time.Now().UnixNano() * 1000)

  watch := &Watch{
    WatchCategory:               watchInfo.WatchCategory.Enum(),
    WatchName:                   proto.String(watchInfo.WatchName),
    WatchPriority:               Watch_kNormal.Enum(),
    ReturnPreviousEntityState:   proto.Bool(watchInfo.ReturnPreviousEntityState),
  }

  if watchInfo.WatchCategory == Watch_kEntityWatch {
    watch.ApplyProjectionOnFiredWatch =
      proto.Bool(watchInfo.ApplyProjectionOnFiredWatch)
    watch.FiredWatchProjectionList = watchInfo.FiredWatchProjectionList
  }

  if watchInfo.WatchCategory == Watch_kEntitySchemaWatch {
    watch.EntitySchemaWatchCondition = &EntitySchemaWatchCondition{
      WatchType: proto.Uint32(uint32(watchType)),
    }
    if watchInfo.EntityTypeName != "" {
      watch.WatchSubject = &WatchSubject{
        EntityTypeName: proto.String(watchInfo.EntityTypeName),
      }
    } else {
      watch.WatchSubject = &WatchSubject{}
    }
  } else {
    watch.WatchSubject = &WatchSubject{
      EntityGuid: entityGuid,
    }
  }

  if watchInfo.WatchCategory == Watch_kEntityWatch {
    if len(watchInfo.WatchMetric) > 0 {
      glog.Infof("Watch metric %s", watchInfo.WatchMetric)
      watch.WatchSubject.MetricName = proto.String(watchInfo.WatchMetric)
    }
    entityWatchCondition := &EntityWatchCondition{
      EntityWatchType: proto.Int32(watchType),
      FilterExpr:      watchInfo.FilterExpr,
    }
    watch.EntityWatchCondition = entityWatchCondition
  }

  var watchList []*Watch
  watchList = append(watchList, watch)

  // If this watch was just Unregisterd - after the last watch fired, it
  // will still be there in unregistered list - Remove it from there.
  wc.mapMu.Lock()
  removed := wc.removeFromUnregisterList(watchInfo.WatchName)
  wc.mapMu.Unlock()
  if removed {
    glog.Infof("Removed watch %s from unregistered list", watchInfo.WatchName)
  }
  watchArg := &RegisterWatchArg{
    WatchClient:     wc.WatchClientProto,
    WatchList:       watchList,
    TimestampUsecs:  proto.Uint64(nowUsecs),
    GetCurrentState: proto.Bool(watchInfo.GetCurrentState),
  }
  watchRet := &RegisterWatchRet{}

  retryWait := util_misc.NewExponentialBackoff(4*time.Second, 8*time.Second, 3)
  for {
    // Wait here till we can proceed
    wc.isChannelOpen = <-wc.watchSync
    if wc.isChannelOpen {
      if _, ok := wc.RegisteredWatches[watchInfo.WatchName]; ok {
        wc.watchSync <- true
        return nil, ErrWatchAlreadyRegistered
      }
      err := wc.WatchService.SendMsg("RegisterWatch", watchArg, watchRet, nil)
      wc.watchSync <- true
      // We got a ResetWatchClient so lets retry after sometime and hope the
      // watch client is reconnected by then.
      if err == ErrResetWatchClient && retryWait.Backoff() != util_misc.Stop {
        continue
      } else if err != nil {
        return nil, err
      } else {
        break
      }
    } else {
      // Channel has been closed -> new watch client. If retries are exhausted
      // return error to the client
      if retryWait.Backoff() == util_misc.Stop {
        return nil, ErrInvalidWatchClientSession
      }
    }
  }
  glog.Infof("Registered watch - %s", watchInfo.WatchName)

  // maintain list of all successful watches registered for this client
  wc.RegisteredWatches[watchInfo.WatchName] = watchInfo

  respList := watchRet.GetResponseList()
  if len(respList) == 0 {
    glog.Warning("Error in reteriving the response list")
    return nil, ErrResetWatchClient
  }
  resp := respList[0]
  if resp.GetResponseStatus().GetErrorType() !=
    InsightsErrorProto_kNoError {
    // Check for Error registration
    glog.Warning("Error in registering watch %s", watchInfo.WatchName)
    return nil, ErrWatchNotRegistered
  }

  //Cache the callback info for this watch.
  wc.FiredWatchLastSeqNo = resp.GetFiredWatchLastSequenceNum()
  watchName := watchInfo.WatchName
  watchCb := watchInfo.WatchCb
  watchCb.setupCallback(watchName, entityGuid.GetEntityTypeName(), "",
    resp.GetFiredWatchLastSequenceNum(),
    resp.GetWatchId())

  wc.mapMu.Lock()
  wc.entityWatches[watchName] = watchCb
  wc.mapMu.Unlock()

  currentStateList := resp.GetCurrentStateList()

  return currentStateList, nil
}

func (wc *WatchClient) invokeFiredWatchesFromQ() error {
  for {
    fireCb := <-wc.FiredWatchQ

    if fireCb.CloseQ {
      // Done - return
      return ErrChannelClosed
    }
    if fireCb != nil {
      fireCb.WatchCb.WatchCallback(fireCb.FiredWatch)
    }
  }
}

// Re-registers existing registered watches
// Returns a list of watches that could not be re-registered
func (wc *WatchClient) Reregister() ([]*RegisterWatchResult, error) {
  var err error = nil

  _, err = wc.transitionWatchState(InsightsEvent_kInitReregister)
  if err != nil {
    return nil, err
  }

  //First unregister with the server.
  err = wc.Unregister()
  if err != nil {
    glog.Warning(err)
  }

  //Re-register for a new session.
  err = wc.Register()
  if err != nil {
    glog.Warningf("Error in re-registering watch client - %s:%s",
      wc.ClientId, err)
    _, _ = wc.transitionWatchState(InsightsEvent_kReset)
    return nil, err
  }

  registeredWatches := make(map[string]*EntityWatchInfo)
  // Copying the registered watches map and clearing it.
  registeredWatches, wc.RegisteredWatches =
    wc.RegisteredWatches, registeredWatches

  glog.Infof("Registered watch count %d", len(registeredWatches))
  // Register all the watches again.
  var watchResult []*RegisterWatchResult
  for _, watchInfo := range registeredWatches {
    changeList, err := wc.RegisterWatch(watchInfo)

    if err != nil {
      // Reset the watch client and let call re-register again
      _, err1 := wc.transitionWatchState(InsightsEvent_kReset)
      glog.Warning(err1)
      // Copy back the registered watches so that they can be reregistered on
      // next retry.
      registeredWatches, wc.RegisteredWatches =
        wc.RegisteredWatches, registeredWatches
      return nil, err
    }

    var entList []*Entity
    for j := 0; j < len(changeList); j++ {
      entList = append(entList, changeList[j].GetChangedEntityData())
    }

    retInfo := &RegisterWatchResult{
      Error:      err,
      EntityList: entList,
      WatchInfo:  watchInfo,
    }
    watchResult = append(watchResult, retInfo)
  }
  _, err = wc.transitionWatchState(InsightsEvent_kReregisterDone)
  if err != nil {
    return nil, err
  }
  return watchResult, nil
}

// This function returns the type of operation that happened on a given entity.
// If returned error is not nil that means there is some error and in this case
// returned operation type is not a reliable result. User must make sure error
// is nil before using the returned operation type.
func (entity *Entity) WatchEntityOperation() (
  EntityWatchCondition_EntityWatchType, error) {

  if entity.DeletedTimestampUsecs != nil {
    return EntityWatchCondition_kEntityDelete, nil
  } else if *entity.CreatedTimestampUsecs == *entity.ModifiedTimestampUsecs {
    return EntityWatchCondition_kEntityCreate, nil
  } else if *entity.CreatedTimestampUsecs < *entity.ModifiedTimestampUsecs {
    return EntityWatchCondition_kEntityUpdate, nil
  }
  glog.Error("Not able to identify operation type- %s", *entity)
  return EntityWatchCondition_kEntityUpdate,
    errors.New("Not able to identify operation type.")
}

func init() {
  // initialize the state transition map.
  // the map contains only valid transitions.
  wcStateFSM[InsightsWatchState_kUnInitialized][InsightsEvent_kInit] =
    InsightsWatchState_kInitialized

  // transitions from Initialized on init-register, init-reregister
  wcStateFSM[InsightsWatchState_kInitialized][InsightsEvent_kInitRegister] =
    InsightsWatchState_kRegistering
  wcStateFSM[InsightsWatchState_kInitialized][InsightsEvent_kInitReregister] =
    InsightsWatchState_kReregistering

  // transition from registering on registerdone
  wcStateFSM[InsightsWatchState_kRegistering][InsightsEvent_kRegisterDone] =
    InsightsWatchState_kRegistered

  wcStateFSM[InsightsWatchState_kRegistering][InsightsEvent_kError] =
    InsightsWatchState_kError

  // transitions from registered on
  // start, init-unregister, reset
  wcStateFSM[InsightsWatchState_kRegistered][InsightsEvent_kStart] =
    InsightsWatchState_kStarted
  wcStateFSM[InsightsWatchState_kRegistered][InsightsEvent_kInitUnregister] =
    InsightsWatchState_kUnregistering
  wcStateFSM[InsightsWatchState_kRegistered][InsightsEvent_kReset] =
    InsightsWatchState_kReset

  // transitions from started on reset, init-unregister, init-stop
  wcStateFSM[InsightsWatchState_kStarted][InsightsEvent_kReset] =
    InsightsWatchState_kReset
  wcStateFSM[InsightsWatchState_kStarted][InsightsEvent_kInitUnregister] =
    InsightsWatchState_kUnregistering
  wcStateFSM[InsightsWatchState_kStarted][InsightsEvent_kInitStop] =
    InsightsWatchState_kStopping

  // transition from unregistering on unregistered
  wcStateFSM[InsightsWatchState_kUnregistering][InsightsEvent_kUnregistered] =
    InsightsWatchState_kUnregistered

  // transition from unregistered on init
  wcStateFSM[InsightsWatchState_kUnregistered][InsightsEvent_kInit] =
    InsightsWatchState_kInitialized

  // transitions from reset on unregistering, init-reregister and stopped
  wcStateFSM[InsightsWatchState_kReset][InsightsEvent_kInitUnregister] =
    InsightsWatchState_kUnregistering
  wcStateFSM[InsightsWatchState_kReset][InsightsEvent_kInitReregister] =
    InsightsWatchState_kReregistering
  wcStateFSM[InsightsWatchState_kReset][InsightsEvent_kInitStop] =
    InsightsWatchState_kStopping

  // transitions from Stopping on Stopped
  wcStateFSM[InsightsWatchState_kStopping][InsightsEvent_kStopped] =
    InsightsWatchState_kStopped

  // transitions from Re-registering to registered on
  // valid events
  wcStateFSM[InsightsWatchState_kReregistering][InsightsEvent_kInitUnregister] =
    InsightsWatchState_kReregistering
  wcStateFSM[InsightsWatchState_kReregistering][InsightsEvent_kUnregistered] =
    InsightsWatchState_kReregistering
  wcStateFSM[InsightsWatchState_kReregistering][InsightsEvent_kInit] =
    InsightsWatchState_kReregistering
  wcStateFSM[InsightsWatchState_kReregistering][InsightsEvent_kInitRegister] =
    InsightsWatchState_kReregistering
  wcStateFSM[InsightsWatchState_kReregistering][InsightsEvent_kRegisterDone] =
    InsightsWatchState_kReregistering
  wcStateFSM[InsightsWatchState_kReregistering][InsightsEvent_kReregisterDone] =
    InsightsWatchState_kRegistered
  wcStateFSM[InsightsWatchState_kReregistering][InsightsEvent_kReset] =
    InsightsWatchState_kReset

  // transitions from Stop on init
  wcStateFSM[InsightsWatchState_kStopped][InsightsEvent_kInitStop] =
    InsightsWatchState_kStopped
  wcStateFSM[InsightsWatchState_kStopped][InsightsEvent_kStart] =
    InsightsWatchState_kStarted
  wcStateFSM[InsightsWatchState_kStopped][InsightsEvent_kInitReregister] =
    InsightsWatchState_kReregistering
}

// transition watch client on event to valid state.
// returns previous state after transition and err if it was an invalid
// transition
func (wc *WatchClient) transitionWatchState(event int32) (int32, error) {
  var err error = nil
  wc.stateMu.Lock()

  nextState := wcStateFSM[wc.state][event]
  previousState := wc.state
  if nextState != 0 {
    wc.state = nextState
    glog.Infof("Transition from [%s][%s] -> [%s]",
      InsightsStateString[previousState], InsightsEventString[event],
      InsightsStateString[wc.state])
  } else {
    err = ErrInvalidClientState
    glog.Infof("Invalid transition from [%s][%s] -> [%s]",
      InsightsStateString[previousState], InsightsEventString[event],
      InsightsStateString[nextState])
  }
  wc.stateMu.Unlock()

  return previousState, err
}

// Client to invoke this as a goroutine
// Depending upon the error returned, Client can take apporopriate
// actions - including reseting of watch client and re-registering.
func (wc *WatchClient) Start() error {
  var err error = nil

  // raise the start event
  _, err = wc.transitionWatchState(InsightsEvent_kStart)
  if err != nil {
    switch wc.state {
    case InsightsWatchState_kStarted:
      err = ErrClientAlreadyStarted
    default:
      err = ErrInvalidClientState
    }
    return err
  }

  arg := &GetFiredWatchListArg{
    WatchClient: wc.WatchClientProto,
  }

  go func() {
    err := wc.invokeFiredWatchesFromQ()
    if err != nil && err != ErrChannelClosed {
      glog.Warning("Error from Invoking Fired watches")
      glog.Warningln(err)
    }
  }()

  // Loop till not done.
  for !wc.done {
    arg.RpcTimeoutInSecs = proto.Uint32(uint32(*WatchClientTimeout))
    arg.FiredWatchLastSequenceNum = proto.Uint64(wc.FiredWatchLastSeqNo)

    // cleanup the unregistered watches
    wc.mapMu.Lock()
    wc.entityUnregWatches = wc.entityUnregWatches[0:0]
    wc.mapMu.Unlock()

    ret := &GetFiredWatchListRet{}

    err := wc.WatchService.SendMsgWithTimeout("GetFiredWatchList", arg, ret,
                                              nil, int64(*WatchClientTimeout))
    if err == ErrCanceled {
      glog.Warning(err)
    } else if err != nil {
      msg := fmt.Sprintf("Reset watch: Fired Watch returned with %s", err)
      glog.Warning(msg)
      _, err = wc.transitionWatchState(InsightsEvent_kReset)
      if err != nil {
        glog.Warning("Incorrect State transition")
      }
      return ErrResetWatchClient
    }
    // Proceed on no error.
    firedList := ret.GetFiredWatchList()
    for i := 0; i < len(firedList); i++ {
      fired := firedList[i]
      seqNo := fired.GetSequenceNum()
      wc.FiredWatchLastSeqNo = seqNo
      watchName := fired.GetWatchName()
      // Lookup the callback

      wc.mapMu.Lock()
      cb := wc.entityWatches[watchName]
      wc.mapMu.Unlock()
      if cb == nil {
        if wc.isWatchUnregistered(watchName) {
          glog.Info("Received callback - but watch is unregistered")
        }
        glog.Warning("Received watch with no callback")
        continue
      }
      if seqNo < cb.GetFiredWatchInitialSeqNo() {
        glog.Warningf(`Warning - Recd fired watch = %s with a seq. number
          = %d, lower than %d.\n`, watchName, seqNo,
          cb.GetFiredWatchInitialSeqNo())
        continue
      }
      // Queue the fired event so that it can be executed in
      // sequence.
      fireCb := &WatchCbInfo{
        WatchCb:    cb,
        FiredWatch: fired,
        CloseQ:     false,
      }
      wc.FiredWatchQ <- fireCb
    }

    // check if we are done before going back to listen for
    // next round of fired watches
    done := false
    wc.mu.Lock()
    done = wc.done
    wc.mu.Unlock()
    if done {
      glog.Info("Watches Done\n")
      fireCb := &WatchCbInfo{
        CloseQ: true,
      }
      wc.FiredWatchQ <- fireCb
    }
  }

  // client has been stopped
  wc.transitionWatchState(InsightsWatchState_kStopped)

  return nil
}
