package common

import (
	"encoding/hex"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"reflect"
	"strings"
)

// The main purpose of HexBytes is to enable HEX-encoding for json/encoding.
type HexBytes []byte

// EncodeValue  for mongodb encode
func (bz *HexBytes) EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if val.IsNil() {
		return errors.New("HexByte value is nil")
	}
	return vw.WriteString(fmt.Sprintf("%X", val.Bytes()))
}

// DecodeValue negates the value of ID when reading
func (bz *HexBytes) DecodeValue(ectx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	i, err := vr.ReadInt64()
	if err != nil {
		return err
	}
	val.SetInt(i * -1)
	return nil
}

// Marshal needed for protobuf compatibility
func (bz HexBytes) Marshal() ([]byte, error) {
	return bz, nil
}

// Unmarshal needed for protobuf compatibility
func (bz *HexBytes) Unmarshal(data []byte) error {
	*bz = data
	return nil
}

// This is the point of Bytes.
func (bz HexBytes) MarshalJSON() ([]byte, error) {
	s := strings.ToUpper(hex.EncodeToString(bz))
	jbz := make([]byte, len(s)+2)
	jbz[0] = '"'
	copy(jbz[1:], []byte(s))
	jbz[len(jbz)-1] = '"'
	return jbz, nil
}

// This is the point of Bytes.
func (bz *HexBytes) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("Invalid hex string: %s", data)
	}
	bz2, err := hex.DecodeString(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}
	*bz = bz2
	return nil
}

// Allow it to fulfill various interfaces in light-client, etc...
func (bz HexBytes) Bytes() []byte {
	return bz
}

func (bz HexBytes) String() string {
	return strings.ToUpper(hex.EncodeToString(bz))
}

func (bz HexBytes) Format(s fmt.State, verb rune) {
	switch verb {
	case 'p':
		s.Write([]byte(fmt.Sprintf("%p", bz)))
	default:
		s.Write([]byte(fmt.Sprintf("%X", []byte(bz))))
	}
}
