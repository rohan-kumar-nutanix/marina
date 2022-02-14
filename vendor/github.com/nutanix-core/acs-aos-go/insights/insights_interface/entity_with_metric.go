/*
 * Copyrigright (c) 2017 Nutanix Inc. All rights reserved.
 *
 * Adds helper methods to 'EntityWithMetric' protobuf generated object.
 */

package insights_interface

import (
  "errors"
  "strings"
  proto "github.com/golang/protobuf/proto"
  "github.com/golang/glog"
)

//-----------------------------------------------------------------------------

func (entity *EntityWithMetric) GetValue(attrName string) (
  interface{}, error) {

  attrMap := entity.GetMetricDataList()
  if attrMap == nil || len(attrMap) == 0 {
    return nil, ErrNotFound
  }

  found := false
  for i := 0; i < len(attrMap) && !found; i++ {
    if strings.Compare(*attrMap[i].Name, attrName) == 0 {
      found = true

      valueList := attrMap[i].GetValueList()
      if valueList == nil {
        return nil, ErrNotFound.SetCause(errors.New("ValueList is empty"))
      }
      // Get the first element
      value := valueList[0].GetValue()
      if value == nil {
        return nil, ErrNotFound.SetCause(errors.New("Value is nil"))
      }
      val, _, err := value.Value()
      return val, err
    }
  }
  return nil, ErrNotFound
}

//-----------------------------------------------------------------------------

func (entity *EntityWithMetric) GetString(attrName string) (string, error) {
  if entity == nil {
    return "", errors.New("Entity is nil, attrName: " + attrName)
  }
  val, err := entity.GetValue(attrName)
  if val == nil || err != nil {
    return "", nil
  }
  v := val.(string)
  return v, nil
}

//-----------------------------------------------------------------------------

func (entity *EntityWithMetric) GetStringList(attrName string) (
  []string, error) {
  if entity == nil {
    return nil, errors.New("Entity is nil, attrName: " + attrName)
  }

  val, err := entity.GetValue(attrName)
  if val == nil || err != nil {
    return nil, nil
  }
  switch v := val.(type) {
  case string:
    return []string{v}, nil
  case []string:
    return v, nil
  default:
    return nil, nil
  }
}

//-----------------------------------------------------------------------------

func (entity *EntityWithMetric) MustGetString(attrName string) string {
  str, err := entity.GetString(attrName)
  if err != nil {
    glog.Fatal(err)
  }
  return str
}

//-----------------------------------------------------------------------------

func (entity *EntityWithMetric) DeserializeEntity(
  protoMessage proto.Message) error {
  attrValue, err := entity.GetValue(COMPRESSED_PROTOBUF_ATTR)
  if err == ErrNotFound {
    return errors.New("Entity does not have compressed protobuf")
  }
  compressedProto := attrValue.([]byte)
  return DeserializeProto(compressedProto, protoMessage)
}
