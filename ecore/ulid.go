package ecore

import (
	"reflect"

	"github.com/oklog/ulid/v2"
	"github.com/vmihailenco/msgpack/v5"
)

const ulidExtID int8 = 1

func init() {
	msgpack.RegisterExtEncoder(ulidExtID, ulid.ULID{}, ulidEncoder)
	msgpack.RegisterExtDecoder(ulidExtID, ulid.ULID{}, ulidDecoder)
}

func ulidEncoder(e *msgpack.Encoder, v reflect.Value) ([]byte, error) {
	id := v.Interface().(ulid.ULID)
	return id.MarshalBinary()
}

func ulidDecoder(d *msgpack.Decoder, v reflect.Value, extLen int) error {
	bytes := make([]byte, extLen)
	err := d.ReadFull(bytes)
	if err != nil {
		return err
	}
	id := v.Addr().Interface().(*ulid.ULID)
	return id.UnmarshalBinary(bytes)
}
