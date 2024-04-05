package ecore

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack/v5"
)

func TestUUID_Compare(t *testing.T) {
	id1 := uuid.New()
	s := id1.String()
	id2, _ := uuid.Parse(s)
	require.False(t, &id1 == &id2)
	require.True(t, id1 == id2)
}

func TestUUID_KeyMap(t *testing.T) {
	id1 := uuid.New()
	s := id1.String()
	id2, _ := uuid.Parse(s)
	m := map[uuid.UUID]string{id1: "test"}
	v, isValue := m[id2]
	require.True(t, isValue)
	require.Equal(t, "test", v)
}

func TestUUID_EncodeMsgPack(t *testing.T) {
	id, _ := uuid.Parse("b4c6c281-69b6-48e6-9847-850e52cafe2e")
	w := &bytes.Buffer{}
	e := msgpack.NewEncoder(w)
	require.NoError(t, e.Encode(id))
	require.Equal(t, []byte{216, 2, 180, 198, 194, 129, 105, 182, 72, 230, 152, 71, 133, 14, 82, 202, 254, 46}, w.Bytes())
}

func TestUUID_DecodeMsgPack(t *testing.T) {
	expected, _ := uuid.Parse("b4c6c281-69b6-48e6-9847-850e52cafe2e")
	array := []byte{216, 2, 180, 198, 194, 129, 105, 182, 72, 230, 152, 71, 133, 14, 82, 202, 254, 46}
	d := msgpack.NewDecoder(bytes.NewReader(array))
	id, err := d.DecodeInterface()
	require.Nil(t, err)
	require.Equal(t, expected, id)
}

func TestUUID_EncodeDecodeMsgPack(t *testing.T) {
	encoded := uuid.New()

	// encode
	w := &bytes.Buffer{}
	e := msgpack.NewEncoder(w)
	require.NoError(t, e.Encode(encoded))

	// decode
	d := msgpack.NewDecoder(bytes.NewReader(w.Bytes()))
	decoded, err := d.DecodeInterface()
	require.Nil(t, err)
	require.Equal(t, encoded, decoded)
}
