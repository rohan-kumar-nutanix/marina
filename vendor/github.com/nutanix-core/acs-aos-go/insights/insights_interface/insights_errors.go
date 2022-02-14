package insights_interface

import (
  ntnx_errors "github.com/nutanix-core/acs-aos-go/nutanix/util-go/errors"
)

// Insights interface error def.
type InsightsInterfaceError_ struct {
  *ntnx_errors.NtnxError
}

func (e *InsightsInterfaceError_) TypeOfError() int {
  return ntnx_errors.InsightErrorType
}

func InsightsError(errMsg string, errCode int) *InsightsInterfaceError_ {
  return &InsightsInterfaceError_{ntnx_errors.NewNtnxError(errMsg, errCode)}
}

func (e *InsightsInterfaceError_) SetCause(err error) error {
  return ntnx_errors.NewNtnxErrorRef(err, e.GetErrorDetail(),
    e.GetErrorCode(), e)
}

func (e *InsightsInterfaceError_) Equals(err interface{}) bool {
  if obj, ok := err.(*InsightsInterfaceError_); ok {
    return e == obj
  }
  return false
}

// Errors defined in protobuf.
var (
  ErrNoError                            = InsightsError("NoError", 0)
  ErrUnknown                            = InsightsError("Unknown", 1)
  ErrCanceled                           = InsightsError("Canceled", 2)
  ErrTimeout                            = InsightsError("Timeout", 3)
  ErrTransportError                     = InsightsError("TransportError", 4)
  ErrRetry                              = InsightsError("Retry", 5)
  ErrInternalError                      = InsightsError("InternalError", 6)
  ErrUnavailable                        = InsightsError("Unavailable", 7)
  ErrNonCasUpdateForCasEntity           = InsightsError("NonCasUpdateForCasEntity", 8)
  ErrCasUpdateForNonCasEntity           = InsightsError("CasUpdateForNonCasEntity", 9)
  ErrIncorrectCasValue                  = InsightsError("IncorrectCasValue", 10)
  ErrIncorrectCreatedTimestamp          = InsightsError("IncorrectCreatedTimestamp", 11)
  ErrPartial                            = InsightsError("Partial", 12)
  ErrEntityTypeNotRegistered            = InsightsError("EntityTypeNotRegistered", 13)
  ErrParentEntityTypeNotRegistered      = InsightsError("ParentEntityTypeNotRegistered", 14)
  ErrShardEntityTypeNotRegistered       = InsightsError("ShardEntityTypeNotRegistered", 15)
  ErrMetricTypeNotRegistered            = InsightsError("MetricTypeNotRegistered", 16)
  ErrEntityNotPresent                   = InsightsError("EntityNotPresent", 17)
  ErrNotFound                           = InsightsError("NotFound", 18)
  ErrAttributeMultipleDefinition        = InsightsError("AttributeMultipleDefinition", 19)
  ErrResetWatchClient                   = InsightsError("ResetWatchClient", 20)
  ErrAttributeInPutMetricDataRPC        = InsightsError("AttributeInPutMetricDataRPC", 21)
  ErrMetricInUpdateEntityRPC            = InsightsError("MetricInUpdateEntityRPC", 22)
  ErrInvalidTimestamp                   = InsightsError("InvalidTimestamp", 23)
  ErrWatchAlreadyRegistered             = InsightsError("WatchAlreadyRegistered", 24)
  ErrEntityReadOnly                     = InsightsError("EntityReadOnly", 25)
  ErrWatchClientAlreadyRegistered       = InsightsError("WatchClientAlreadyRegistered", 26)
  ErrWatchNotRegistered                 = InsightsError("WatchNotRegistered", 27)
  ErrWatchClientUnexpectedState         = InsightsError("WatchClientUnexpectedState", 28)
  ErrDerivedMetricInPutMetricDataRPC    = InsightsError("DerivedMetricInPutMetricDataRPC", 29)
  ErrInvalidWatchIdForUnregistration    = InsightsError("InvalidWatchIdForUnregistration", 30)
  ErrMetricDataDroppedByRetentionPolicy = InsightsError("MetricDataDroppedByRetentionPolicy", 31)
  ErrIncorrectIncarnationId             = InsightsError("IncorrectIncarnationId", 32)
  ErrEntityAlreadyExists                = InsightsError("EntityAlreadyExists", 33)
  ErrNestedFieldDirectUpdatesNotAllowed = InsightsError("NestedFieldDirectUpdatesNotAllowed", 34)
  ErrWatchClientNotRegistered           = InsightsError("WatchClientNotRegistered", 35)
  ErrTenantIdChanged                    = InsightsError("TenantIdChanged", 36)
  ErrEntityUpdateOutOfOrder             = InsightsError("EntityUpdateOutOfOrder", 37)
  ErrCurrentStateNotAllowedForEvictableEntityTypeWatch =
    InsightsError("CurrentStateNotAllowedForEvictableEntityTypeWatch", 38)
  ErrAmbiguousAncestorTree              = InsightsError("AmbiguousAncestorTree", 39)
  ErrFeatureDisabled                    = InsightsError("FeatureDisabled", 40)
  ErrRetryCursorQueryFromStart          = InsightsError("RetryCursorQueryFromStart", 41)
  ErrInvalidRequest                     = InsightsError("InvalidRequest", 1000)
  ErrMultipleDefinition                 = InsightsError("MultipleDefinition", 1001)
  ErrNoMetricName                       = InsightsError("NoMetricName", 1002)
  ErrNoEntityTypeName                   = InsightsError("NoEntityTypeName", 1003)
  ErrAttributeIsBucketized              = InsightsError("AttributeIsBucketized", 1004)
  ErrInvalidMetricTypeName              = InsightsError("InvalidMetricTypeName", 1005)
  ErrInvalidMetricTypeDefinition        = InsightsError("InvalidMetricTypeDefinition", 1006)
  ErrInvalidDownSamplingInterval        = InsightsError("InvalidDownSamplingInterval", 1007)
  ErrAttributeIsDownsampled             = InsightsError("AttributeIsDownsampled", 1008)
  ErrInvalidSerialisedProto             = InsightsError("InvalidSerialisedProto", 1009)
  ErrInvalidNestedFieldInfo             = InsightsError("InvalidNestedFieldInfo", 1010)
  ErrInvalidBaseAttribute               = InsightsError("InvalidBaseAttribute", 1011)
  ErrCannotUpdateNestedFieldInfo        = InsightsError("CannotUpdateNestedFieldInfo", 1012)
  ErrInvalidPutMetricData               = InsightsError("InvalidPutMetricData", 1013)
  ErrCannotClearBaseAttributeSchema     = InsightsError("CannotClearBaseAttributeSchema", 1014)
  ErrInvalidSerialisedProtoCompressionType =
    InsightsError("InvalidSerialisedProtoCompressionType", 1015)
  ErrInvalidEvictableEntityType         = InsightsError("InvalidEvictableEntityType", 1016)
  ErrInvalidIndexColumn                 = InsightsError("InvalidIndexColumn", 1017)
  ErrGetAllEvictableEntitiesUnsupported = InsightsError("GetAllEvictableEntitiesUnsupported", 1018)
  ErrDimentsionFactTableInEntityWithMetric =
    InsightsError("DimentsionFactTableInEntityWithMetric", 1019)
  ErrEntityInMetricDataSample           = InsightsError("EntityInMetricDataSample", 1020)
  ErrMissingPrimaryKeyAttribute         = InsightsError("MissingPrimaryKeyAttribute", 1021)
  ErrPrimaryKeyIsNotAttribute           = InsightsError("PrimaryKeyIsNotAttribute", 1022)
  ErrMetricTypeDoesNotMatchIndexType    = InsightsError("MetricTypeDoesNotMatchIndexType", 1033)
  ErrMissingIndexType                   = InsightsError("MissingIndexType", 1034)
  CannotModifyParentTypeName            = InsightsError("CannotModifyParentTypeName", 1035)
  ErrDBUpdateInProgress                 = InsightsError("DBUpdateInProgress", 2001)
  ErrCacheSyncWithDBInProgress          = InsightsError("CacheSyncWithDBInProgress", 2002)
  ErrCacheSyncWithDBFailed              = InsightsError("CacheSyncWithDBFailed", 2003)
  ErrCassandraMutateFailed              = InsightsError("CassandraMutateFailed", 2004)
  ErrCassandraMutateEpochError          = InsightsError("CassandraMutateEpochError", 2005)
  ErrCassandraReadFailed                = InsightsError("CassandraReadFailed", 2006)
)

var Errors = map[int]error{
  0:    ErrNoError,
  1:    ErrUnknown,
  2:    ErrCanceled,
  3:    ErrTimeout,
  4:    ErrTransportError,
  5:    ErrRetry,
  6:    ErrInternalError,
  7:    ErrUnavailable,
  8:    ErrNonCasUpdateForCasEntity,
  9:    ErrCasUpdateForNonCasEntity,
  10:   ErrIncorrectCasValue,
  11:   ErrIncorrectCreatedTimestamp,
  12:   ErrPartial,
  13:   ErrEntityTypeNotRegistered,
  14:   ErrParentEntityTypeNotRegistered,
  15:   ErrShardEntityTypeNotRegistered,
  16:   ErrMetricTypeNotRegistered,
  17:   ErrEntityNotPresent,
  18:   ErrNotFound,
  19:   ErrAttributeMultipleDefinition,
  20:   ErrResetWatchClient,
  21:   ErrAttributeInPutMetricDataRPC,
  22:   ErrMetricInUpdateEntityRPC,
  23:   ErrInvalidTimestamp,
  24:   ErrWatchAlreadyRegistered,
  25:   ErrEntityReadOnly,
  26:   ErrWatchClientAlreadyRegistered,
  27:   ErrWatchNotRegistered,
  28:   ErrWatchClientUnexpectedState,
  29:   ErrDerivedMetricInPutMetricDataRPC,
  30:   ErrInvalidWatchIdForUnregistration,
  31:   ErrMetricDataDroppedByRetentionPolicy,
  32:   ErrIncorrectIncarnationId,
  33:   ErrEntityAlreadyExists,
  34:   ErrNestedFieldDirectUpdatesNotAllowed,
  35:   ErrWatchClientNotRegistered,
  36:   ErrTenantIdChanged,
  37:   ErrEntityUpdateOutOfOrder,
  38:   ErrCurrentStateNotAllowedForEvictableEntityTypeWatch,
  39:   ErrAmbiguousAncestorTree,
  40:   ErrFeatureDisabled,
  41:   ErrRetryCursorQueryFromStart,
  1000: ErrInvalidRequest,
  1001: ErrMultipleDefinition,
  1002: ErrNoMetricName,
  1003: ErrNoEntityTypeName,
  1004: ErrAttributeIsBucketized,
  1005: ErrInvalidMetricTypeName,
  1006: ErrInvalidMetricTypeDefinition,
  1007: ErrInvalidDownSamplingInterval,
  1008: ErrAttributeIsDownsampled,
  1009: ErrInvalidSerialisedProto,
  1010: ErrInvalidNestedFieldInfo,
  1011: ErrInvalidBaseAttribute,
  1012: ErrCannotUpdateNestedFieldInfo,
  1013: ErrInvalidPutMetricData,
  1014: ErrCannotClearBaseAttributeSchema,
  1015: ErrInvalidSerialisedProtoCompressionType,
  1016: ErrInvalidEvictableEntityType,
  1017: ErrInvalidIndexColumn,
  1018: ErrGetAllEvictableEntitiesUnsupported,
  1019: ErrDimentsionFactTableInEntityWithMetric,
  1020: ErrEntityInMetricDataSample,
  1021: ErrMissingPrimaryKeyAttribute,
  1022: ErrPrimaryKeyIsNotAttribute,
  1033: ErrMetricTypeDoesNotMatchIndexType,
  1034: ErrMissingIndexType,
  1035: CannotModifyParentTypeName,
  2001: ErrDBUpdateInProgress,
  2002: ErrCacheSyncWithDBInProgress,
  2003: ErrCacheSyncWithDBFailed,
  2004: ErrCassandraMutateFailed,
  2005: ErrCassandraMutateEpochError,
  2006: ErrCassandraReadFailed,
}
