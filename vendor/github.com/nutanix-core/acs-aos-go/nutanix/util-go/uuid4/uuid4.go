/*
 * Copyright (c) 2016 Nutanix Inc. All rights reserved.
 * This module creates Version 4 UUIDs as specified in RFC 4122
 *
 */

package uuid4

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/golang/glog"
	"github.com/google/uuid"
)

const ZeroUuid = "00000000-0000-0000-0000-000000000000"

// Hashable and comparable.
type Uuid struct {
	uuidBytes [16]byte
}

// Hashable and comparable.
// Use Equals() to compare two UUIDs explicitly.
func EmptyUuid() *Uuid {
	return &Uuid{[16]byte{}}
}

// Create a new instance of Uuid4.
// Use Equals() to compare the instance returned.
func New() (*Uuid, error) {
	b := [16]byte{}

	_, err := rand.Read(b[:])
	if err != nil {
		return nil, err
	}

	// Set the two most significant bits (bits 6 and 7) of the
	// clock_seq_hi_and_reserved to zero and one, respectively.
	b[8] = (b[8] | 0x40) & 0x7F

	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	b[6] = (b[6] & 0xF) | (4 << 4)

	uuid := &Uuid{
		uuidBytes: b,
	}

	// Return unparsed version of the generated UUID sequence.
	//return fmt.Sprintf("%x-%x-%x-%x-%x",
	//	b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil

	return uuid, nil
}

// Utility function to get Uuid4 string formed from input list of strings.
func NewFromStrings(inputs []string) string {
	var buffer bytes.Buffer
	for _, str := range inputs {
		buffer.WriteString(str)
	}

	md5Sum := md5.Sum(buffer.Bytes())

	// Convert this md5sum to cluster uuid format
	uuid := ToUuid4(md5Sum[0:])
	return uuid.UuidToString()
}

// utility function to get an Uuid4 object from raw bytes.
// Expectation is bytes is in uuid4 format.
func ToUuid4(b []byte) *Uuid {
	if len(b) != 16 {
		glog.Errorf("Invalid uuid bytes: %x", b)
		return nil
	}
	u := Uuid{}
	copy(u.uuidBytes[:], b)
	return &u
}

// Get raw bytes from  uuid4 instance
func (u *Uuid) RawBytes() []byte {
	return u.uuidBytes[:]
}

// convert 16 bytes of UUID to a standard 36 byte formatted string
// Insight takes in UUID as 36 bytes of formatted string
func (u *Uuid) UuidToString() string {
	if u == nil || len(u.uuidBytes) == 0 {
		return ""
	}
	uuid := u.uuidBytes
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:16])
}

// Convert 36 bytes of formatted UUID string to 16 raw bytes of
// representing the uuid.
// Most infrastructure components - ergon/acropolis/uhura etc take
// in 16 byte uuid
func StringToUuid4(str string) (*Uuid, error) {
	uuid4 := EmptyUuid()
	if err := uuid4.Set(str); err != nil {
		return nil, err
	}
	return uuid4, nil
}

// Define Flag.Value interface for Uuid
// This makes it possible for uuid to be used as an arg/opt directly
// while using flag package (https://golang.org/pkg/flag/)
func (u *Uuid) Set(v string) error {
	if _, err := uuid.Parse(v); err != nil {
		return err
	}
	str := fmt.Sprintf("%s%s%s%s%s", v[0:8], v[9:13], v[14:18],
		v[19:23], v[24:36])
	uuidBytes, err := hex.DecodeString(str)
	if err != nil {
		return err
	}
	copy(u.uuidBytes[:], uuidBytes)
	return nil
}

func (u *Uuid) String() string {
	return u.UuidToString()
}

func (u *Uuid) Equals(uuid *Uuid) bool {
	return bytes.Equal(u.RawBytes(), uuid.RawBytes())
}

func (u *Uuid) CopyFrom(uuid *Uuid) {
	u.uuidBytes = uuid.uuidBytes
}

func (u *Uuid) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

func (u *Uuid) UnmarshalJSON(b []byte) error {
	var uuidStr string
	if err := json.Unmarshal(b, &uuidStr); err != nil {
		return err
	}
	uuidPtr, err := StringToUuid4(uuidStr)
	if err != nil {
		return err
	}
	u.CopyFrom(uuidPtr)
	return nil
}

// SkipEmpytyUnmarshal ensures *Uuid field is nil in an IDF unmarshalled struct
// if the corressponding field value in IDF was an empty string.
func (u *Uuid) SkipEmpytyUnmarshal() {
}
