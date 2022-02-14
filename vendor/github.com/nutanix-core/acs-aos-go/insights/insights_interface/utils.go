/*
 * Copyrigright (c) 2020 Nutanix Inc. All rights reserved.
 *
 */

package insights_interface

import (
  "bytes"
  "compress/zlib"
  proto "github.com/golang/protobuf/proto"
  "io"
)

// Deserializes the nested protobuf value in an IDF entry. On success this
// function will return nil, otherwise it will return an error.
func DeserializeProto(compressedProto []byte, protoMessage proto.Message) error {
  compressedProtoBuffer := bytes.NewBuffer(compressedProto)
  zlibReader, err := zlib.NewReader(compressedProtoBuffer)
  if err != nil {
    return err
  }
  uncompressedBuf := &bytes.Buffer{}
  _, err = io.Copy(uncompressedBuf, zlibReader)
  if err != nil {
    return err
  }
  return proto.Unmarshal(uncompressedBuf.Bytes(), protoMessage)
}
