package ecore

import (
	"bytes"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack/v5"
)

func TestULID_Compare(t *testing.T) {
	id1 := ulid.Make()
	s := id1.String()
	id2, _ := ulid.Parse(s)
	require.False(t, &id1 == &id2)
	require.True(t, id1 == id2)
}

func TestULID_KeyMap(t *testing.T) {
	id1 := ulid.Make()
	s := id1.String()
	id2, _ := ulid.Parse(s)
	m := map[ulid.ULID]string{id1: "test"}
	v, isValue := m[id2]
	require.True(t, isValue)
	require.Equal(t, "test", v)
}

func TestULID_EncodeMsgPack(t *testing.T) {
	id, _ := ulid.Parse("01HSB2RH3YTNM4922XGR2R2N51")
	w := &bytes.Buffer{}
	e := msgpack.NewEncoder(w)
	require.NoError(t, e.Encode(id))
	require.Equal(t, []byte{216, 1, 1, 142, 86, 44, 68, 126, 213, 104, 68, 136, 93, 134, 5, 129, 84, 161}, w.Bytes())
}

func TestULID_DecodeMsgPack(t *testing.T) {
	expected, _ := ulid.Parse("01HSB2RH3YTNM4922XGR2R2N51")
	array := []byte{216, 1, 1, 142, 86, 44, 68, 126, 213, 104, 68, 136, 93, 134, 5, 129, 84, 161}
	d := msgpack.NewDecoder(bytes.NewReader(array))
	id, err := d.DecodeInterface()
	require.Nil(t, err)
	require.Equal(t, expected, id)
}

func TestULID_EncodeDecodeMsgPack(t *testing.T) {
	encoded := ulid.Make()

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
