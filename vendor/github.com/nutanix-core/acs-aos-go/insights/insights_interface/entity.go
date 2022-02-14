/*
 * Copyrigright (c) 2017 Nutanix Inc. All rights reserved.
 *
 * Adds helper methods to 'Entity' protobuf generated object.
 */

package insights_interface

import (
        "bytes"
        "compress/zlib"
        "errors"
        proto "github.com/golang/protobuf/proto"
)

// Attribute that is used to store any compressed protobuf message that might
// be kept with the entity. Note that the serialized protobuf can be kept under
// any attribute, but this is the most commonly used location.
var COMPRESSED_PROTOBUF_ATTR = "__zprotobuf__"

func (entity *Entity) GetString(attrName string) (string, error) {
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

func (entity *Entity) GetStringList(attrName string) ([]string, error) {
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

func (entity *Entity) GetBool(attrName string) (bool, error) {
        val, err := entity.GetValue(attrName)
        if val == nil || err != nil {
                return false, nil
        }
        v := val.(bool)
        return v, nil
}

func (entity *Entity) GetValue(attrName string) (interface{}, error) {
        attrMap := entity.GetAttributeDataMap()
        if attrMap == nil || len(attrMap) == 0 {
                return nil, ErrNotFound
        }

        for i := 0; i < len(attrMap); i++ {
                attr := attrMap[i]
                if attr.GetName() != attrName {
                        continue
                }

                val, _, err := attr.GetValue().Value()
                return val, err
        }
        return nil, ErrNotFound
}

// Utility function to get the values for list of attributes from an entity
// This allows the caller to deal with only language (GO) types and not
// Insight data types
// Example:
//   map[Name]value, err := GetAttributeValues(entity, {"app_container_name"})
//
func (entity *Entity) GetValueMap(attrNames []string) (
        map[string]interface{}, error) {

        attributeValueMap := make(map[string]interface{})
        for _, attrName := range attrNames {
                value, _ := entity.GetValue(attrName)
                attributeValueMap[attrName] = value
        }
        return attributeValueMap, nil
}

// Deserializes the nested protobuf value in an IDF entry. On success this
// function will return nil, otherwise it will return an error. Note this
// method assumes that the serialized protobuf is kept in the attribute
// being pointed at by COMPRESSED_PROTOBUF_ATTR.
func (entity *Entity) DeserializeEntity(protoMessage proto.Message) error {
        attrValue, err := entity.GetValue(COMPRESSED_PROTOBUF_ATTR)
        if err == ErrNotFound {
                return errors.New("Entity does not have compressed protobuf")
        }
        compressedProto := attrValue.([]byte)
        return DeserializeProto(compressedProto, protoMessage)
}

// Serializes and compresses a protobuf object and stores it into
// COMPRESSED_PROTOBUF_ATTR. This method takes a protobuf object to serialize
// and an UpdateEntityArg, which is where the attribute will be inserted. If
// there is already a serialized entity present in COMPRESSED_PROTOBUF_ATTR
// this method will overwrite the value of that attribute. Returns nil on
// success and an error otherwise.
func AddSerializedProto(protoMessage proto.Message,
        updateArg *UpdateEntityArg) error {
        serializedProto, err := proto.Marshal(protoMessage)
        if err != nil {
                return err
        }
        compressedProtoBuf := &bytes.Buffer{}
        zlibWriter := zlib.NewWriter(compressedProtoBuf)
        zlibWriter.Write(serializedProto)
        zlibWriter.Close()

        compressedProtoValue := &DataValue{
                ValueType: &DataValue_BytesValue{
                        BytesValue: compressedProtoBuf.Bytes(),
                },
        }
        // If there is an existing attribute with the serialized protobuf already
        // present then we will overwrite it, otherwise we will add a new attribute.
        var protoAttrDataArg *AttributeDataArg
        for _, attrDataArg := range updateArg.GetAttributeDataArgList() {
                if attrDataArg.GetAttributeData().GetName() == COMPRESSED_PROTOBUF_ATTR {
                        protoAttrDataArg = attrDataArg
                        break
                }
        }

        if protoAttrDataArg != nil {
                protoAttrDataArg.AttributeData.Value = compressedProtoValue
                return nil
        }
        protoAttrDataArg = &AttributeDataArg{
                Operation: AttributeDataArg_kSET.Enum(),
                AttributeData: &AttributeData{
                        Name:  proto.String(COMPRESSED_PROTOBUF_ATTR),
                        Value: compressedProtoValue,
                },
        }
        updateArg.AttributeDataArgList = append(updateArg.AttributeDataArgList,
                protoAttrDataArg)
        return nil
}
