/*
 * Copyright (c) 2017 Nutanix Inc. All rights reserved.
 *
 * This module defines enums for watches.
 */

package insights_interface

var (
  kEntityCreate = EntityWatchCondition_EntityWatchType_value["kEntityCreate"]
  kEntityUpdate = EntityWatchCondition_EntityWatchType_value["kEntityUpdate"]
  kEntityDelete = EntityWatchCondition_EntityWatchType_value["kEntityDelete"]
  kRegisterEntityType =
    EntitySchemaWatchCondition_SchemaWatchType_value["kRegisterEntityType"]
  kUpdateEntityType =
    EntitySchemaWatchCondition_SchemaWatchType_value["kUpdateEntityType"]
  kRegisterMetricType =
    EntitySchemaWatchCondition_SchemaWatchType_value["kRegisterMetricType"]
  kUpdateMetricType =
    EntitySchemaWatchCondition_SchemaWatchType_value["kUpdateMetricType"]
  kNoReset = WatchClientProto_ResetWatchClientOnOperations_value["kNoReset"]
  kResetOnClusterRegistration =
    WatchClientProto_ResetWatchClientOnOperations_value[
      "kResetOnClusterRegistration"]
  kResetOnClusterUnregistration =
    WatchClientProto_ResetWatchClientOnOperations_value[
      "kResetOnClusterUnregistration"]
)
