/*
 * Copyrigright (c) 2017 Nutanix Inc. All rights reserved.
 *
 * Adds helper methods to 'DataValue' protobuf generated object.
 */

package insights_interface

import "encoding/json"

func (m *DataValue) Value() (interface{}, string, error) {
  if m == nil {
    return nil, "", ErrNotFound
  }

  switch val := m.GetValueType().(type) {
  case *DataValue_StrValue:
    // return string value
    return val.StrValue, "str_value", nil
  case *DataValue_Int64Value:
    return val.Int64Value, "int64_value", nil
  case *DataValue_Uint64Value:
    return val.Uint64Value, "uint64_value", nil
  case *DataValue_FloatValue:
    return val.FloatValue, "float_value", nil
  case *DataValue_DoubleValue:
    return val.DoubleValue, "double_value", nil
  case *DataValue_BoolValue:
    return val.BoolValue, "bool_value", nil
  case *DataValue_BytesValue:
    return val.BytesValue, "bytes_value", nil
  case *DataValue_JsonObj:
    var jsonDict map[string]interface{}
    err := json.Unmarshal(val.JsonObj, &jsonDict)
    return jsonDict, "json_obj", err
  case *DataValue_StrList_:
    // return [] string
    strList := val.StrList
    DataValue_type := "str_list"
    if strList == nil {
      // this should be some other more specific
      return nil, DataValue_type, ErrNotFound
    }
    return strList.GetValueList(), DataValue_type, nil
  case *DataValue_Int64List_:
    // return [] int64
    intList := val.Int64List
    if intList == nil {
      // this should be some other more specific
      return nil, "int64_list", ErrNotFound
    }
    return intList.GetValueList(), "int64_list", nil
  case *DataValue_Uint64List:
    // return [] uint64
    uintList := val.Uint64List
    if uintList == nil {
      // this should be some other more specific
      return nil, "uint64_list", ErrNotFound
    }
    return uintList.GetValueList(), "uint64_list", nil
  case *DataValue_BoolList_:
    // return [] bool
    boolList := val.BoolList
    if boolList == nil {
      // this should be some other more specific
      return nil, "bool_list", ErrNotFound
    }
    return boolList.GetValueList(), "bool_list", nil
  case *DataValue_HistogramObj:
    return val.HistogramObj, "histogram", nil
  default:
  }
  return nil, "", ErrNotFound
}
