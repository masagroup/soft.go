package ecore

import (
	"reflect"

	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
)

const uuidExtID int8 = 2

func init() {
	msgpack.RegisterExtEncoder(uuidExtID, uuid.UUID{}, uuidEncoder)
	msgpack.RegisterExtDecoder(uuidExtID, uuid.UUID{}, uuidDecoder)
}

func uuidEncoder(e *msgpack.Encoder, v reflect.Value) ([]byte, error) {
	id := v.Interface().(uuid.UUID)
	return id.MarshalBinary()
}

func uuidDecoder(d *msgpack.Decoder, v reflect.Value, extLen int) error {
	bytes := make([]byte, extLen)
	err := d.ReadFull(bytes)
	if err != nil {
		return err
	}
	id := v.Addr().Interface().(*uuid.UUID)
	return id.UnmarshalBinary(bytes)
}
